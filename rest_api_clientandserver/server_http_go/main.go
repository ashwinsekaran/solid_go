// Command server_http_go is a minimal REST API for a User resource, built with the
// standard library only (net/http, no framework).
//
// It demonstrates several idioms:
//   - Go 1.22 method+path routing patterns on http.ServeMux ("POST /users",
//     "GET /users/{id}").
//   - An in-memory store guarded by a sync.RWMutex for concurrent-safe access.
//   - Bearer-token authentication as wrapping middleware around the whole mux.
//   - Graceful shutdown on SIGINT/SIGTERM via http.Server.Shutdown with a timeout.
//
// Routes (all require the header "Authorization: Bearer abcxyz"):
//
//	POST /users        create a user from a JSON body (id required)
//	GET  /users/{id}   fetch a single user by id
//	GET  /users?limit= list users, optionally capped by ?limit=N (default 100)
//
// The paired client lives in ../client_http_go.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// User is the resource stored and served by the API. Created is stamped
// server-side at save time and returned in JSON responses.
type User struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

// Store is a concurrent-safe, in-memory user store. The RWMutex lets many
// readers (get/list) proceed together while writes (save) take exclusive access.
type Store struct {
	mu   sync.RWMutex
	data map[string]User
}

// NewStore returns a Store with its backing map initialised and ready to use.
func NewStore() *Store {
	return &Store{
		data: make(map[string]User),
	}
}

// main wires the routes, wraps them in auth middleware, starts the server in a
// goroutine, and blocks until an interrupt signal triggers a graceful shutdown.
func main() {
	store := NewStore()
	r := http.NewServeMux()

	r.HandleFunc("GET /users/{id}", store.get)
	r.HandleFunc("GET /users", store.list)
	r.HandleFunc("POST /users", store.save)
	var handler http.Handler = r
	handler = auth(handler) // wrap the whole mux so every route is authenticated

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Listen for Ctrl-C / SIGTERM so we can shut down without dropping requests.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	<-shutdown
	log.Printf("shutting down...")

	// Give in-flight requests up to 5s to finish before forcing the close.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown %v", err)
	}
}

// auth is middleware that rejects any request lacking the expected bearer token
// with 401 Unauthorized, and otherwise passes it through to the next handler.
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer abcxyz" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// ctx := context.WithValue(r.Context(), "user", "app")
		// next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)

	})
}

// save handles POST /users: it decodes the JSON body into a User, requires a
// non-empty id, stamps Created, stores it, and returns 201 Created. Malformed
// bodies or a missing id yield 400 Bad Request.
func (s *Store) save(w http.ResponseWriter, r *http.Request) {
	var u User
	var id string
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "incorrect data", http.StatusBadRequest)
		return
	}

	if u.Id != "" {
		id = u.Id
	} else {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	u.Created = time.Now()
	s.data[id] = u

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created successfully"))

}

// get handles GET /users/{id}: it looks up the user under a read lock and
// returns it as JSON, or 404 Not Found if the id is unknown.
func (s *Store) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// v := r.URL.Query().Get("limit")
	s.mu.RLock()
	u, ok := s.data[id]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "user not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)

}

// list handles GET /users: it returns up to `limit` users as a JSON array
// (default 100, overridable via ?limit=N). Note map iteration order is
// non-deterministic, so which users appear when the limit truncates may vary.
func (s *Store) list(w http.ResponseWriter, r *http.Request) {
	limit := 100
	v := r.URL.Query().Get("limit")
	if v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			limit = n
		}
	}

	s.mu.RLock()
	res := make([]User, 0, len(s.data))

	for _, u := range s.data {
		if len(res) >= limit {
			break
		}
		res = append(res, u)

	}

	s.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
