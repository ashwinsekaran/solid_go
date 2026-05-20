package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

// RateLimiterStore holds a per-user rate limiter keyed by user ID.
// A single mutex guards the entire map — fine for moderate concurrency.
// For very high traffic, shard by user ID to reduce contention.
type RateLimiterStore struct {
	rateLimiter map[string]*rateLimiter
	mu          sync.Mutex
}

// rateLimiter tracks token-bucket state for one user.
// It is not safe to use concurrently on its own — the parent store's mu protects it.
type rateLimiter struct {
	capacity   int       // maximum tokens the bucket can hold
	tokens     int       // current available tokens
	refillRate int       // tokens to restore on each refill (currently full restore)
	lastRefill time.Time // timestamp used to decide when to refill
}

func newRateLimiterStore() *RateLimiterStore {
	return &RateLimiterStore{
		rateLimiter: make(map[string]*rateLimiter),
	}
}

// Allow checks whether the given user is within their rate limit.
// On first request, a new limiter is created and the first token is consumed immediately.
// On subsequent requests, the bucket is refilled if ≥1 second has elapsed since the last refill.
func (rl *RateLimiterStore) Allow(userId string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	data, ok := rl.rateLimiter[userId]
	if ok {
		// Existing user: check whether a refill window has elapsed.
		elapsed := time.Since(data.lastRefill)
		if elapsed.Seconds() >= 1 {
			// Refill fully (refillRate == capacity here), reset the clock.
			data.tokens = data.capacity
			data.lastRefill = time.Now()
		}

		if data.tokens > 0 {
			data.tokens--
			return true // request is allowed
		}
		// Bucket empty, refill window not reached — throttle.
		return false
	}

	// First request from this user: create a limiter and consume the first token.
	r := rateLimiter{
		capacity:   10,
		tokens:     9, // 10 - 1 for this first request
		refillRate: 10,
		lastRefill: time.Now(),
	}
	rl.rateLimiter[userId] = &r
	return true
}

func main() {
	r := httprouter.New()
	rl := newRateLimiterStore()

	// Wrap the get handler with the rate-limiting middleware.
	// Middleware reads X-User-ID from the request header and delegates to Allow().
	r.Handle("GET", "/", Middleware(rl, get))

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error %v", err)
	}
}

// Middleware is a higher-order function: it wraps any httprouter.Handle with rate limiting.
// Returning an httprouter.Handle (rather than using http.Handler) keeps it compatible with
// httprouter's route registration and preserves access to path parameters.
func Middleware(rl *RateLimiterStore, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user := r.Header.Get("X-User-ID") // identify the caller by a request header
		if rl.Allow(user) {
			next(w, r, p) // within limit — pass through to the real handler
		} else {
			// 429 Too Many Requests is the standard status code for rate limiting.
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
}

func get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
