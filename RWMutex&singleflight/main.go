package main

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

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
	d, err := GetUser("Ashwin", cache)
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}
	fmt.Println(d)
}

func GetUser(userId string, cache *Cache) (Data, error) {
	data, ok := cache.Get(userId)
	if ok {
		return data, nil
	}

	result, err, _ := group.Do(userId, func() (interface{}, error) {
		return GetFromDB(userId), nil
	})
	if err != nil {
		return Data{}, err
	}

	d := result.(Data)
	cache.Set(userId, d)

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
	return Data{}
}
