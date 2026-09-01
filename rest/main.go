package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"solid_go/rest/auth"
	"solid_go/rest/handlers"
	"solid_go/rest/health"
	"solid_go/rest/metrics"
	repo "solid_go/rest/repo"
	"solid_go/rest/uc"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config is populated from environment variables by envconfig.
//   - HTTP_ADDRESS (split_words converts HttpAddress → HTTP_ADDRESS) defaults to ":8080".
//   - SHUTDOWN_DELAY is the drain window after readiness flips to not-ready and before
//     the listener closes, so the load balancer observes the not-ready state and stops
//     routing new traffic first. In k8s this should be ≥ the readiness probe period.
type Config struct {
	HttpAddress   string        `split_words:"true" default:":8080"`
	ShutdownDelay time.Duration `split_words:"true" default:"5s"`
}

func main() {
	var config Config
	envconfig.Process("", &config) // reads env vars into config fields

	// Build the three-layer stack: repo → use-case → handler.
	// Each layer only depends on the interface of the layer below it (DIP in practice).
	rep := repo.NewRepo()

	// Metrics collectors + their registry (served at /metrics via promhttp below).
	m := metrics.NewMetrics()
	reg := m.Register()

	r := httprouter.New()

	// Use-cases are plain functions (closures) that close over the repo.
	// Handlers receive use-case functions, keeping HTTP logic separate from business logic.
	getUc := uc.MakeGetUc(rep)
	SaveUc := uc.MakeSaveUc(rep)

	// Metrics middleware wraps each business handler per-route (like worker_pool_api_processer),
	// so the /metrics endpoint itself is never instrumented by it.
	r.Handle("GET", "/:id", metrics.Middleware(m, handlers.GetHandler(getUc)))
	r.Handle("POST", "/post", metrics.Middleware(m, handlers.PostHandler(SaveUc)))

	// Readiness flag: starts not-ready, set true once wiring is complete below.
	ready := &health.Readiness{}

	// Outer http.ServeMux carries the infra endpoints that Kubernetes and Prometheus
	// hit WITHOUT the app's bearer token, and forwards everything else ("/") to the
	// auth-wrapped API router:
	//   - /metrics scrape endpoint (also avoids a static-vs-wildcard clash with "/:id").
	//   - /.well-known/live  and /.well-known/ready k8s probes.
	//   - auth.Auth wraps the entire API router (like rest_api_clientandserver), so every
	//     business route requires a bearer token.
	root := http.NewServeMux()
	root.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	root.HandleFunc("/.well-known/live", health.Live)
	root.HandleFunc("/.well-known/ready", ready.Handler)
	root.Handle("/", auth.Auth(r))

	// Everything is wired — start accepting traffic.
	ready.SetReady(true)

	server := http.Server{
		Addr:    config.HttpAddress,
		Handler: root,
	}

	// Graceful shutdown: listen for SIGINT / SIGTERM on a buffered channel.
	// Buffered(1) ensures the signal is not dropped if we're not yet blocking on <-shutdown.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine so main can block on the shutdown signal below.
	go func() {
		log.Printf("Server started on %s", config.HttpAddress)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve error, %v", err)
		}
	}()

	<-shutdown // block until OS sends a termination signal

	log.Println("Shutting down server...")

	// Fail readiness first so Kubernetes removes this pod from Service endpoints and
	// stops routing new traffic here, while liveness stays green (no restart).
	ready.SetReady(false)

	// Wait out the drain window so the load balancer actually observes the not-ready
	// state before we close the listener; otherwise new requests would race in until
	// the socket closes. Liveness keeps returning 200 throughout.
	if config.ShutdownDelay > 0 {
		log.Printf("readiness not-ready; draining for %s before shutdown...", config.ShutdownDelay)
		time.Sleep(config.ShutdownDelay)
	}

	// Give in-flight requests up to 5 seconds to complete before the process exits.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
