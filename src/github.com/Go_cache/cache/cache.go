package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

//Cache ...
type Cache struct {
	*cache
}

type cache struct {
	mu       sync.RWMutex
	cacheMap map[interface{}]interface{}
}

// func (c *cache) AddToCache(url string, response string) bool {

// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.cacheMap[url] = response
// 	return true
// }

func (c *cache) AddToCache(key interface{}, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = value
	return true
}

func (c *cache) FindResponse(request interface{}) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, isPresent := c.cacheMap[request]
	if isPresent {
		return val, true
	}

	return "", false
}

func (c *cache) SaveToFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(fp)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	err = enc.Encode(&c.cacheMap)

	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := make(map[interface{}]interface{})
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			_, found := c.cacheMap[k]
			if !found {
				c.cacheMap[k] = v
			}
		}
	}
	return err
}

func (c *cache) LoadFromFile(fname string) error {

	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

//NewCache ...
func NewCache() *Cache {

	items := make(map[interface{}]interface{})
	c := &cache{
		cacheMap: items,
	}
	C := &Cache{c}
	return C
}
