package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type FileProcessor struct {
	File     []byte
	Response chan Result
	Ctx      context.Context
}

type Result struct {
	Success string
	Err     error
}

func main() {
	wg := sync.WaitGroup{}

	jobCh := make(chan FileProcessor, 3)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				processFile(job)
			}
		}()
	}

	r := httprouter.New()
	r.Handle("POST", "/post", file(jobCh))

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, os.Interrupt)
	<-shutdown

	close(jobCh)
	wg.Wait()

}

func file(jobCh chan<- FileProcessor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resCh := make(chan Result, 1)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var f FileProcessor
		f.Response = resCh
		f.Ctx = ctx

		//err := json.NewDecoder(r.Body).Decode(&f.File)
		fi, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f.File = fi

		jobCh <- f
		result := <-f.Response
		if result.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.Success))

	}
}

func processFile(f FileProcessor) {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("file processed")
		f.Response <- Result{Success: "file processed"}
	case <-f.Ctx.Done():
		f.Response <- Result{Err: f.Ctx.Err()}

	}
}
