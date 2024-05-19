package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	CacheEntry map[string]Entry
	Mu         sync.Mutex
}

type Entry struct {
	createdAt time.Time
	val       []byte
}

// func NewCache(interval time.Duration) Cache {

// 	c := Cache{
// 		cacheEntry: map[string]Entry{},
// 		mu:         &sync.Mutex{},
// 	}

// 	return c

// }

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	fmt.Printf("Adding to cache for  url %s\n", key)
	currentTime := time.Now().UTC()

	// ??
	// if currentTime.Sub(c.CacheEntry[key].createdAt) > time.Second * 30 {
	c.CacheEntry[key] = Entry{currentTime, val}
	// }
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	fmt.Printf("Cache hit for url %s", key)
	currentTime := time.Now().UTC()
	// check if the cache was created 10 within current time
	// t1 - created < 10
	if c.CacheEntry[key].val != nil && currentTime.Sub(c.CacheEntry[key].createdAt) < time.Second*300 {
		fmt.Printf(" Succesful\n")
		return c.CacheEntry[key].val, true
	}
	fmt.Printf(" Failed\n")
	return nil, false
}
