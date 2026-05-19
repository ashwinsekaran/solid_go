package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

// curl -X POST http://localhost:8080/post \                                                                                                                                    ✔  16:11:37
//  -d '"hello world"' -v

type Job struct {
	Image []byte
	Ctx   context.Context
	ResCh chan Result
}

type Result struct {
	Data []byte
	Err  error
}

func main() {
	wg := sync.WaitGroup{}
	jobCh := make(chan Job, 10)
	workers := 100

	r := httprouter.New()

	r.Handle("POST", "/post", post(jobCh))

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				select {
				case <-job.Ctx.Done():
					job.ResCh <- Result{Err: job.Ctx.Err()}
					continue
				default:
				}

				data, err := processImage(job.Image)

				job.ResCh <- Result{Data: data, Err: err}
			}
		}()
	}

	quit := make(chan struct{})

	go func() {
		<-quit
		close(jobCh)
		wg.Wait()
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

}

func post(jobCh chan<- Job) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//var image []byte
		var job Job
		job.ResCh = make(chan Result, 1)
		job.Ctx = r.Context()

		image, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		job.Image = image
		//err := json.NewDecoder(r.Body).Decode(&job.Image)
		//if err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		jobCh <- job

		result := <-job.ResCh
		if result.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result.Data)
		return
	}
}

func processImage(image []byte) ([]byte, error) {
	return image, nil
}
