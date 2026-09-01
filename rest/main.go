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
// HTTP_ADDRESS (split_words converts HttpAddress → HTTP_ADDRESS) defaults to ":8080".
type Config struct {
	HttpAddress string `split_words:"true" default:":8080"`
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

	// Two middleware, mounted at different levels:
	//   - auth.Auth wraps the entire API router (like rest_api_clientandserver), so every
	//     business route requires a bearer token.
	//   - /metrics is served from an outer http.ServeMux, deliberately OUTSIDE auth so
	//     Prometheus can scrape without a token. It also can't live on the httprouter itself:
	//     a static "/metrics" route conflicts with the root wildcard "/:id" and httprouter
	//     would panic. The ServeMux "/" pattern forwards everything else to the authed router.
	root := http.NewServeMux()
	root.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	root.Handle("/", auth.Auth(r))

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

	// Give in-flight requests up to 5 seconds to complete before the process exits.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
