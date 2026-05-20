package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/singleflight"
)

// Store is a concurrent-safe in-memory map of short-code → full URL.
// RWMutex lets multiple goroutines read simultaneously (high-traffic redirect path)
// while writes (cache population) are exclusive.
type Store struct {
	Data map[string]string
	mu   sync.RWMutex
}

// group deduplicates concurrent DB lookups for the same short code.
// If 100 goroutines miss the cache for "/abc" at the same time, only one
// GetFromDB call is made; the other 99 block and share its result.
var group singleflight.Group

func NewStore() *Store {
	return &Store{Data: make(map[string]string)}
}

func main() {
	data := NewStore()

	r := httprouter.New()
	r.Handle("GET", "/r/:url", redirect(data))

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// redirect resolves a short code to its full URL and issues an HTTP 302.
// Cache hit path: RLock → return immediately (no DB call, no singleflight overhead).
// Cache miss path: singleflight → one DB call → populate cache → redirect.
func redirect(data *Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		url := p.ByName("url") // short code from the path parameter

		// Fast path: return from cache without touching the DB.
		val, ok := data.Get(url)
		if ok {
			http.Redirect(w, r, val, http.StatusFound)
			return
		}

		// Slow path: on a cache miss, use singleflight to coalesce concurrent DB lookups.
		// The key is the short code — all goroutines racing for the same code share one call.
		res, err, _ := group.Do(url, func() (interface{}, error) {
			return GetFromDB(url), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := res.(string)
		data.Set(url, result) // populate the cache so future requests take the fast path

		http.Redirect(w, r, result, http.StatusFound) // 302 Found
	}
}

// Get acquires a read lock — allows many concurrent redirects to read simultaneously.
func (c *Store) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok := c.Data[key]
	return data, ok
}

// Set acquires an exclusive write lock — blocks all readers and other writers.
func (c *Store) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data[key] = value
}

// GetFromDB simulates a slow lookup from a persistent store (database, Redis, etc.).
// In production, query the DB for the full URL associated with the short code.
func GetFromDB(key string) string {
	return "https://ashwinsekaran.com"
}
