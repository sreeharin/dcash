package main

import (
	"errors"
	"slices"
	"sync"
)

type Cache struct {
	mu      sync.RWMutex
	entries map[string][]byte
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string][]byte)}
}

func (c *Cache) Set(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = slices.Clone(value)
	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	val, exists := c.entries[key]
	if !exists {
		return nil, errors.New("KEY doesn't exists.")
	}

	return val, nil
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.entries[key]
	return exists
}

func (c *Cache) Delete(key string) error {
	if c.Exists(key) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return nil
	}

	return errors.New("KEY doesn't exists.")
}
