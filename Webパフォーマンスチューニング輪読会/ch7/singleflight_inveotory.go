package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
	"time"
)

var group singleflight.Group

type Cache struct {
	mu    sync.Mutex
	items map[int]int
}

func NewCache() *Cache {
	m := make(map[int]int)
	c := &Cache{
		items: m,
	}
	return c
}

func (c *Cache) Set(key int, value int) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()
	v, ok := c.items[key]
	c.mu.Unlock()

	if ok {
		return v
	}

	// Use singleflight to avoid multiple HeavyGet calls
	vv, err, _ := group.Do(fmt.Sprintf("cacheGet_%d", key), func() (interface{}, error) {
		value := HeavyGet(key)
		c.Set(key, value)
		return value, nil
	})

	if err != nil {
		panic(err)
	}

	return vv.(int)
}

func HeavyGet(key int) int {
	log.Printf("HeavyGet(%d)\n", key)
	time.Sleep(time.Second)
	return key * 2
}

func main() {
	mCache := NewCache()

	for i := 0; i < 100; i++ {
		go func(i int) {
			mCache.Get(i % 10)
		}(i)
	}

	time.Sleep(time.Second * 2)

	for i := 0; i < 10; i++ {
		log.Println(mCache.Get(i))
	}
}
