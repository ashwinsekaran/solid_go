package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"solid_go/rest/handlers"
	repo "solid_go/rest/repo"
	"solid_go/rest/uc"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
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

	r := httprouter.New()

	// Use-cases are plain functions (closures) that close over the repo.
	// Handlers receive use-case functions, keeping HTTP logic separate from business logic.
	getUc := uc.MakeGetUc(rep)
	SaveUc := uc.MakeSaveUc(rep)

	r.Handle("GET", "/:id", handlers.GetHandler(getUc))
	r.Handle("POST", "/post", handlers.PostHandler(SaveUc))

	server := http.Server{
		Addr:    config.HttpAddress,
		Handler: r,
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
