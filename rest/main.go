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

type Config struct {
	HttpAddress string `split_words:"true" default:":8080"`
}

func main() {
	var config Config

	envconfig.Process("", &config)

	rep := repo.NewRepo()

	r := httprouter.New()

	getUc := uc.MakeGetUc(rep)
	SaveUc := uc.MakeSaveUc(rep)

	r.Handle("GET", "/:id", handlers.GetHandler(getUc))
	r.Handle("POST", "/post", handlers.PostHandler(SaveUc))

	server := http.Server{
		Addr:    config.HttpAddress,
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on %s", config.HttpAddress)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve error, %v", err)
		}
	}()
	<-shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

}
