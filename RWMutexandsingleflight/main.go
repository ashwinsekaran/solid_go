package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group
var dbCallCount atomic.Int64

type Cache struct {
	mu    sync.RWMutex
	Store map[string]Data
}

type Data struct {
	Id   int
	Name string
}

func main() {
	cache := NewCache()
	wg := sync.WaitGroup{}

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

func GetUser(userId string, cache *Cache) (Data, error) {
	data, ok := cache.Get(userId)
	if ok {
		return data, nil
	}

	result, err, shared := group.Do(userId, func() (interface{}, error) {
		return GetFromDB(userId), nil
	})
	if err != nil {
		return Data{}, err
	}

	d := result.(Data)
	cache.Set(userId, d)

	if shared {
		fmt.Println("singleflight deduped a call")
	}

	return d, nil
}

func NewCache() *Cache {
	return &Cache{
		Store: make(map[string]Data),
	}
}

func (s *Cache) Get(key string) (Data, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, i := s.Store[key]
	return data, i

}

func (s *Cache) Set(key string, value Data) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Store[key] = value
}

func GetFromDB(userId string) Data {
	dbCallCount.Add(1)
	time.Sleep(100 * time.Millisecond)

	return Data{Id: 1, Name: userId}
}
