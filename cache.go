package main

import (
	"sync"
	"slices"
)

type Cache struct {
	mu 	sync.RWMutex
	entries map[string][]byte
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string][]byte)}
}

func (c *Cache) Set (key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = slices.Clone(value)
	return nil
}

func (c *Cache) Exists (key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.entries[key]	
	return exists
}

