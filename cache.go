package main

import (
	"errors"
	"hash/fnv"
	"log"
	"slices"
	"sync"
)

type Cache struct {
	shards []*shard

	// shardCount is the number of shards
	shardCount int
}

type shard struct {
	mu      sync.RWMutex
	entries map[string][]byte
}

// NewCache creates a map with 'shardCount' shards
func NewCache() *Cache {
	// TODO: accept user input for shardCount
	shardCount := 64
	shards := make([]*shard, shardCount)

	for i := range shardCount {
		shards[i] = &shard{
			entries: make(map[string][]byte),
		}
	}

	return &Cache{shards: shards, shardCount: shardCount}
}

func (c *Cache) getShard(key string) *shard {
	hash := fnv.New32()
	_, err := hash.Write([]byte(key))

	if err != nil {
		log.Fatal("Error while writing hash. ", err)
	}

	idx := int(hash.Sum32()) % c.shardCount
	return c.shards[idx]
}

func (c *Cache) Set(key string, value []byte) error {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.entries[key] = slices.Clone(value)
	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	shard := c.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	val, exists := shard.entries[key]
	if !exists {
		return nil, errors.New("KEY doesn't exists.")
	}

	return slices.Clone(val), nil
}

func (c *Cache) Exists(key string) bool {
	shard := c.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	_, exists := shard.entries[key]
	return exists
}

func (c *Cache) Delete(key string) error {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.entries, key)

	return errors.New("KEY doesn't exists.")
}
