package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Batch is a mutex-protected in-memory buffer.
// HTTP handlers write into it; the flusher goroutine drains it on a timer.
// The mutex prevents a concurrent flush and an append from racing on Data.
type Batch struct {
	mu   sync.Mutex
	Data []string
}

func main() {
	wg := sync.WaitGroup{}
	batch := newBatch()
	r := httprouter.New()

	r.Handle("POST", "/post", post(batch))

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// quit is a done-channel: closing it signals the flusher to do one final flush and exit.
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		startFlusher(batch, quit)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown // block until SIGINT / SIGTERM

	// Closing quit triggers the flusher's final drain before the process exits.
	close(quit)
	wg.Wait() // don't exit until the flusher has finished its last flush
}

// startFlusher runs a background loop that flushes the batch on a fixed interval.
// When quit is closed (graceful shutdown) it flushes one final time before returning.
func startFlusher(batch *Batch, quit <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			flush(batch) // periodic flush: drain accumulated events to the DB
		case <-quit:
			flush(batch) // final flush: ensure no events are lost on shutdown
			return
		}
	}
}

// flush holds the lock for the shortest possible time:
// it snapshots Data, clears the slice, then writes to the DB.
func flush(batch *Batch) {
	batch.mu.Lock()
	defer batch.mu.Unlock()
	saveToDB(batch.Data)
	batch.Data = []string{} // reset the buffer after successfully draining it
}

// saveToDB simulates persisting a batch of events.
// In production, replace with a bulk INSERT or a message-queue publish.
func saveToDB(data []string) {
	for _, d := range data {
		fmt.Printf("data %s saved to DB\n", d)
	}
}

func newBatch() *Batch {
	return &Batch{Data: []string{}}
}

// post decodes an event from the request body and appends it to the shared batch.
// The HTTP handler never calls the DB directly — it just buffers; the flusher writes.
func post(batch *Batch) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var event string

		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Lock only while appending — keeps the critical section tiny.
		batch.mu.Lock()
		defer batch.mu.Unlock()
		batch.Data = append(batch.Data, event)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("event created"))
	}
}
