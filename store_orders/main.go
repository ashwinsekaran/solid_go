package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type OrderStore struct {
	mu      sync.Mutex
	Data    map[int]Orders
	Counter atomic.Int64
}

type ProcessOrders struct {
	Orders   Orders
	Response chan Result
	Ctx      context.Context
}

type Result struct {
	Success string
	Err     error
}

type Orders struct {
	Product  string
	Quantity int
}

func main() {
	wg := sync.WaitGroup{}

	jobCh := make(chan ProcessOrders, 3)
	store := NewStore()

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				processOrders(job, store)
			}
		}()
	}

	r := httprouter.New()
	r.Handle("POST", "/post", orders(jobCh))
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Printf("server starting at port: 8080")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, os.Interrupt)
	<-shutdown
	log.Printf("shutting down..")
	close(jobCh)
	wg.Wait()

}

func processOrders(o ProcessOrders, store *OrderStore) {
	select {
	case <-time.After(2 * time.Second):
		store.Store(o.Orders)
		o.Response <- Result{Success: fmt.Sprintf("Order processed for : %s", o.Orders.Product)}
	case <-o.Ctx.Done():
		o.Response <- Result{Err: o.Ctx.Err()}
	}

}

func orders(jobCh chan<- ProcessOrders) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var orders ProcessOrders
		resCh := make(chan Result, 1)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		orders.Ctx = ctx

		orders.Response = resCh

		err := json.NewDecoder(r.Body).Decode(&orders.Orders)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jobCh <- orders

		result := <-orders.Response

		if result.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(result.Success))
	}
}

func NewStore() *OrderStore {
	return &OrderStore{
		Data: make(map[int]Orders),
	}
}

func (s *OrderStore) Store(order Orders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := int(s.Counter.Add(1))
	s.Data[id] = order
}
