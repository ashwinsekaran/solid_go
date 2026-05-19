package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// group deduplicates concurrent calls that share the same key.
// Only one real DB call runs at a time per key; all other callers block and share the result.
var group singleflight.Group

// dbCallCount tracks how many times GetFromDB was actually invoked (atomic = no mutex needed).
var dbCallCount atomic.Int64

// Cache is a concurrent-safe in-memory store backed by sync.RWMutex.
// RWMutex allows many readers simultaneously but only one writer at a time.
type Cache struct {
	mu    sync.RWMutex   // guards Store
	Store map[string]Data
}

type Data struct {
	Id   int
	Name string
}

func main() {
	cache := NewCache()
	wg := sync.WaitGroup{}

	// Simulate 1000 concurrent requests for the same user.
	// Without singleflight, this would fire 1000 DB calls on a cold cache.
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, err := GetUser("Ashwin", cache)
			if err != nil {
				log.Fatalf("error getting user: %v", err)
				return
			}
			_ = d
		}()
	}

	wg.Wait()
	fmt.Printf("Total DB calls made: %d\n", dbCallCount.Load())
	fmt.Printf("Requests served:     1000\n")
	fmt.Printf("DB calls saved:      %d\n", 1000-dbCallCount.Load())
}

// GetUser first checks the cache (read lock, cheap).
// On a miss, singleflight ensures only one goroutine calls GetFromDB —
// all other concurrent callers for the same userId wait and share that one result.
func GetUser(userId string, cache *Cache) (Data, error) {
	// Fast path: return from cache without touching the DB.
	data, ok := cache.Get(userId)
	if ok {
		return data, nil
	}

	// Slow path: singleflight.Do collapses duplicate in-flight calls.
	// The key is userId — any goroutine with the same key that arrives while
	// the first call is running will block here and receive the same result.
	result, err, shared := group.Do(userId, func() (interface{}, error) {
		return GetFromDB(userId), nil
	})
	if err != nil {
		return Data{}, err
	}

	d := result.(Data)
	cache.Set(userId, d) // populate cache so future reads skip the DB entirely

	if shared {
		// shared == true means this goroutine received a result produced by another goroutine.
		fmt.Println("singleflight deduped a call")
	}

	return d, nil
}

func NewCache() *Cache {
	return &Cache{Store: make(map[string]Data)}
}

// Get acquires a read lock — multiple goroutines can call Get concurrently.
func (s *Cache) Get(key string) (Data, bool) {
	s.mu.RLock() // allows other readers in simultaneously
	defer s.mu.RUnlock()
	data, i := s.Store[key]
	return data, i
}

// Set acquires an exclusive write lock — all readers and writers block until this returns.
func (s *Cache) Set(key string, value Data) {
	s.mu.Lock() // exclusive: no reads or writes allowed while this runs
	defer s.mu.Unlock()
	s.Store[key] = value
}

// GetFromDB simulates a slow database call.
// atomic.Add is safe without a mutex because it uses a CPU-level compare-and-swap.
func GetFromDB(userId string) Data {
	dbCallCount.Add(1)
	time.Sleep(100 * time.Millisecond) // simulate latency
	return Data{Id: 1, Name: userId}
}
