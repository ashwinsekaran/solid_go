// Different shape from the pool — this one is REST handlers + a concurrent map, which is the area Cresta named that you've practised least live.

// Build a small HTTP service, net/http only, no frameworks, no database

// Requirements

// 1 — POST /shorten accepts JSON {"url": "https://example.com/some/long/path"} and returns {"key": "a1b2c3",
// "short": "http://localhost:8080/a1b2c3"}.
// Generate a short random key. Reject a missing or malformed body with 400.

// 2 — GET /{key} looks up the key and issues an HTTP redirect to the original URL. Unknown key → 404.

// 3 — GET /stats/{key} returns {"url": "...", "hits": 42, "created": "..."}. Every successful redirect increments that key's hit count.

// 4 — Concurrency. The store is an in-memory map behind a mutex. Every handler runs in its own goroutine, so all access must be safe.
// Must be clean under go run -race while being hammered concurrently.

// 5 — Expiry. Entries live 10 minutes. A background goroutine sweeps expired entries every 30 seconds and must stop cleanly when the
// server shuts down — no leaked goroutine.

// Constraints
// No sync.Map
// Graceful shutdown on SIGINT: stop accepting, drain in-flight requests, stop the sweeper
// Include a main that starts the server on :8080

// Command rest_urlshortner_gonative is an in-memory URL shortener built with the
// standard library only (net/http, no framework, no database).
//
// It ties together several concurrency and HTTP idioms:
//   - An in-memory map guarded by a sync.RWMutex (no sync.Map, by constraint).
//   - Go 1.22 method+path routing patterns ("POST /url", "GET /{key}").
//   - Cryptographically random short keys via crypto/rand + base64 URL encoding.
//   - A background sweeper goroutine that evicts entries older than a TTL and is
//     stopped cleanly through context cancellation (no leaked goroutine).
//   - Graceful shutdown on SIGINT/SIGTERM: stop the sweeper, then drain in-flight
//     requests with http.Server.Shutdown under a timeout.
//
// Routes:
//
//	POST /url          body {"url":"https://..."} → {"key","short"}; 400 on a bad URL
//	GET  /{key}        302 redirect to the original URL and increment its hit count
//	GET  /stats/{key}  {"url","hits","created"} for the key; 404 if unknown
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Url is a stored short-link record: the original target, its hit count, and
// the creation time used by the TTL sweeper. It doubles as the /stats response.
type Url struct {
	Url     string    `json:"url"`
	Hits    int       `json:"hits"`
	Created time.Time `json:"created"`
}

// Shortened is the response to POST /url: the generated key and the full short URL.
type Shortened struct {
	Key   string `json:"key"`
	Short string `json:"short"`
}

// Store is the concurrent-safe, in-memory key→Url map. Reads (stats) take an
// RLock; writes and the redirect's hit-count bump take the exclusive Lock.
type Store struct {
	mu   sync.RWMutex
	data map[string]Url
}

// NewStore returns a Store with its backing map initialised.
func NewStore() *Store {
	return &Store{data: make(map[string]Url)}
}

// main starts the sweeper and HTTP server, then blocks until an interrupt
// signal triggers a graceful shutdown of both.
func main() {
	s := NewStore()
	r := http.NewServeMux()

	r.HandleFunc("POST /url", s.postUrl)
	r.HandleFunc("GET /stats/{key}", s.getStats)
	r.HandleFunc("GET /{key}", s.getUrl)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	sweepCtx, stopSweep := context.WithCancel(context.Background())
	defer stopSweep()
	go s.sweep(sweepCtx, 30*time.Second, 10*time.Minute)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-shutdown
	log.Printf("shutting down..")
	stopSweep()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

}

// postUrl handles POST /url: it validates the submitted URL, reuses the existing
// key if that URL was already shortened, otherwise generates a new random key and
// stores it. Responds 201 with the {key, short} JSON, or 400 on a malformed URL.
func (s *Store) postUrl(w http.ResponseWriter, r *http.Request) {

	var url Url
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err = validURL(url.Url)
	if err != nil {
		http.Error(w, "invalid/malformed url", http.StatusBadRequest)
		return
	}
	var shortend Shortened
	key, isExists := s.isUrlExistsAndKey(url.Url)
	if isExists {
		shortend = Shortened{
			Key:   key,
			Short: shortenedUrl(key),
		}

	} else {
		key, err = newKey(6)
		if err != nil {
			http.Error(w, "server error, pls try again later", http.StatusInternalServerError)
			return
		}
		s.mu.Lock()
		url.Created = time.Now()
		s.data[key] = url
		s.mu.Unlock()
		shortend = Shortened{
			Key:   key,
			Short: shortenedUrl(key),
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&shortend)

}

// getUrl handles GET /{key}: it looks up the key, increments its hit count, and
// issues a 302 redirect to the original URL. Unknown keys return 404. The write
// lock is held (not RLock) because the hit count is mutated in place.
func (s *Store) getUrl(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	s.mu.Lock()
	e, ok := s.data[key]
	if !ok {
		s.mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	e.Hits++
	s.data[key] = e
	target := e.Url
	s.mu.Unlock()

	http.Redirect(w, r, target, http.StatusFound)
}

// getStats handles GET /stats/{key}: it returns the stored record (url, hits,
// created) as JSON under a read lock, or 404 if the key is unknown.
func (s *Store) getStats(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	s.mu.RLock()
	e, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

// sweep runs until ctx is cancelled, deleting entries older than ttl on every
// tick. Cancelling the context (on shutdown) makes it return so no goroutine leaks.
func (s *Store) sweep(ctx context.Context, every, ttl time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			cutoff := time.Now().Add(-ttl)
			s.mu.Lock()
			for k, e := range s.data {
				if e.Created.Before(cutoff) {
					delete(s.data, k)
				}
			}
			s.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// newKey returns a URL-safe base64 string derived from n cryptographically
// random bytes (so the encoded key is longer than n characters).
func newKey(n int) (string, error) {
	s := make([]byte, n)
	if _, err := rand.Read(s); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(s), nil
}

// validURL rejects anything that isn't a well-formed absolute http/https URL
// with a host, returning a descriptive error for the 400 response.
func validURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("malformed url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("missing host")
	}
	return nil
}

// isUrlExistsAndKey scans the store under a read lock for an entry whose target
// matches url, returning its key so the same short link is reused. This is an
// O(n) linear scan — fine for a demo, but a reverse index would scale better.
func (s *Store) isUrlExistsAndKey(url string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		if v.Url == url {
			return k, true
		}
	}
	return "", false
}

// shortenedUrl builds the full short link for a key against the local server.
func shortenedUrl(s string) string {
	return fmt.Sprintf("http://localhost:8080/%s", s)
}
