package cache

import (
	"testing"
)

var newCache *Cache = NewCache()

func TestCacheAdd(t *testing.T) {
	//newCache := NewCache()
	key := "Hello"
	value := "World"
	check := newCache.AddToCache(key, value)

	if !check {
		t.Errorf("Check failed expected %v got %v", true, check)
	}

}

func TestCacheFind(t *testing.T) {
	key := "Hello"
	response, ok := newCache.FindResponse(key)

	if !ok {
		t.Errorf("Find failed expected %v got %v", true, ok)
	}

	if response.(string) != "World" {
		t.Errorf("Find failed expected %v got %v", "World", response.(string))
	}

}
