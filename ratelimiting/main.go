package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	Bucket         int
	Tokens         int
	RefillRate     int
	RefillDuration time.Duration
	LastRefill     time.Time
	mu             sync.Mutex
}

func NewTokenBucket(capacity, refillRate int, refillEvery time.Duration) *TokenBucket {
	return &TokenBucket{
		Bucket:         capacity,
		Tokens:         capacity,
		RefillRate:     refillRate,
		RefillDuration: refillEvery,
		LastRefill:     time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	elapsed := time.Since(tb.LastRefill)

	if elapsed >= tb.RefillDuration {
		periods := int(elapsed / tb.RefillDuration)
		tb.Tokens += periods + tb.RefillRate
		if tb.Tokens > tb.Bucket {
			tb.Tokens = tb.Bucket
		}
		tb.LastRefill = time.Now()
		//tb.LastRefill = tb.LastRefill.Add(time.Duration(periods) * tb.RefillDuration)
	}
	if tb.Tokens > 0 {
		tb.Tokens--
		return true
	}
	return false
	//if tb.Tokens == 0 {
	//	if elapsed >= tb.RefillDuration && tb.RefillRate <= tb.Bucket {
	//		tb.Tokens += tb.RefillRate
	//		return true
	//	}
	//	return false
	//} else {
	//	if elapsed >= tb.RefillDuration && tb.RefillRate <= tb.Bucket {
	//		tb.Tokens += tb.RefillRate
	//		return true
	//	}
	//	return true
	//
	//}

	// your code
}

func main() {
	NoOfReq := 10
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
