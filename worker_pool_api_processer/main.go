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

// Job carries everything a worker needs: the raw payload, the request context
// (so workers respect client disconnects), and a reply channel sized 1 so the
// HTTP handler can block waiting for exactly one result.
type Job struct {
	Image []byte
	Ctx   context.Context
	ResCh chan Result // buffered(1): worker sends once, handler receives once
}

// Result is the worker's output — either processed data or an error.
type Result struct {
	Data []byte
	Err  error
}

func main() {
	wg := sync.WaitGroup{}
	jobCh := make(chan Job, 10) // buffer absorbs short bursts without blocking HTTP handlers
	workers := 100

	r := httprouter.New()
	r.Handle("POST", "/post", post(jobCh))

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the fixed-size worker pool.
	// Each worker goroutine lives for the lifetime of the server, pulling jobs from jobCh.
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh { // exits when jobCh is closed (graceful shutdown)
				// Honour client disconnects before doing expensive work.
				select {
				case <-job.Ctx.Done():
					job.ResCh <- Result{Err: job.Ctx.Err()}
					continue
				default:
				}

				data, err := processImage(job.Image)
				job.ResCh <- Result{Data: data, Err: err} // unblocks the waiting HTTP handler
			}
		}()
	}

	// quit channel is the shutdown trigger.
	// Closing jobCh signals all workers to drain and exit; wg.Wait() ensures clean shutdown.
	quit := make(chan struct{})
	go func() {
		<-quit
		close(jobCh) // no more jobs accepted; workers finish in-flight work and return
		wg.Wait()
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// post returns an httprouter.Handle that enqueues a job and synchronously waits for its result.
// The handler blocks on job.ResCh, so the HTTP response is sent only after the worker finishes.
func post(jobCh chan<- Job) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var job Job
		job.ResCh = make(chan Result, 1) // size-1 buffer: worker never blocks on send
		job.Ctx = r.Context()            // propagates client cancellation into the worker

		image, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		job.Image = image

		jobCh <- job // hand off to the worker pool (blocks if jobCh buffer is full)

		// Block until the worker sends its result back.
		result := <-job.ResCh
		if result.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result.Data)
	}
}

// processImage is the CPU/IO-bound work done by workers.
// Replace this with real image processing; kept as a passthrough for demonstration.
func processImage(image []byte) ([]byte, error) {
	return image, nil
}
