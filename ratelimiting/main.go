package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket implements the token-bucket rate-limiting algorithm.
// Tokens accumulate up to Bucket capacity at RefillRate per RefillDuration.
// Each allowed request consumes one token. mu serialises concurrent access.
type TokenBucket struct {
	Bucket         int           // maximum tokens the bucket can hold
	Tokens         int           // current available tokens
	RefillRate     int           // how many tokens to add per refill period
	RefillDuration time.Duration // how often a refill happens
	LastRefill     time.Time     // timestamp of the last refill — used to compute elapsed periods
	mu             sync.Mutex    // guards all mutable fields
}

// NewTokenBucket creates a full bucket (Tokens == capacity) ready to accept requests.
func NewTokenBucket(capacity, refillRate int, refillEvery time.Duration) *TokenBucket {
	return &TokenBucket{
		Bucket:         capacity,
		Tokens:         capacity, // start full so the first burst is immediately available
		RefillRate:     refillRate,
		RefillDuration: refillEvery,
		LastRefill:     time.Now(),
	}
}

// Allow returns true if a request is permitted, false if the bucket is empty.
// It is safe to call from multiple goroutines concurrently.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	elapsed := time.Since(tb.LastRefill)

	// Refill: if enough time has passed, add tokens proportional to elapsed periods.
	if elapsed >= tb.RefillDuration {
		periods := int(elapsed / tb.RefillDuration) // how many full refill windows elapsed
		tb.Tokens += periods * tb.RefillRate
		if tb.Tokens > tb.Bucket {
			tb.Tokens = tb.Bucket // cap at maximum capacity — no burst credit beyond the bucket size
		}
		tb.LastRefill = time.Now()
	}

	// Consume one token if available.
	if tb.Tokens > 0 {
		tb.Tokens--
		return true // request is allowed
	}
	return false // bucket empty — request is throttled
}

func main() {
	NoOfReq := 10
	// Bucket capacity 5, refills 5 tokens every 10 seconds.
	// With 10 simultaneous requests only the first 5 will be allowed.
	tb := NewTokenBucket(5, 5, 10*time.Second)

	wg := sync.WaitGroup{}

	for i := 1; i <= NoOfReq; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			isAllowed := tb.Allow()
			fmt.Printf("Req %d is %v\n", i, isAllowed)
		}(i)
	}

	wg.Wait()
}
