package main

import (
	"errors"
	"hash/fnv"
	"log"
	"slices"
	"sync"
	"time"
)

var ErrKeyDoesntExists = errors.New("KEY doesn't exists.")

type Cache struct {
	shards []*shard

	// shardCount is the number of shards
	shardCount int
}

type entry struct {
	value []byte

	// 0 means no expiry
	expiriesAt int64
}

func (e entry) isExpired() bool {
	now := time.Now().UnixNano()
	return e.expiriesAt > 0 && now >= e.expiriesAt
}

type shard struct {
	mu      sync.RWMutex
	entries map[string]entry
}

// NewCache creates a map with 'shardCount' shards
func NewCache() *Cache {
	// TODO: accept user input for shardCount
	shardCount := 64
	shards := make([]*shard, shardCount)

	for i := range shardCount {
		shards[i] = &shard{
			entries: make(map[string]entry),
		}
	}

	return &Cache{shards: shards, shardCount: shardCount}
}

// getShard is used to find the shard to store the key:value
// it uses the 'key' string to find which shard to store to
func (c *Cache) getShard(key string) *shard {
	hash := fnv.New32()
	_, err := hash.Write([]byte(key))

	if err != nil {
		log.Fatal("Error while writing hash. ", err)
	}

	idx := int(hash.Sum32()) % c.shardCount
	return c.shards[idx]
}

// Set stores the key:value to a shard
func (c *Cache) Set(key string, value []byte, expiresAt int64) error {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.entries[key] = entry{value: value, expiriesAt: expiresAt}
	return nil
}

// Get fetches 'value' for 'key'
// if key is expired, it gets deleted and nil is returned
func (c *Cache) Get(key string) ([]byte, error) {
	shard := c.getShard(key)

	shard.mu.RLock()
	// defer shard.mu.RUnlock()

	entry, exists := shard.entries[key]
	if !exists {
		shard.mu.RUnlock()
		return nil, ErrKeyDoesntExists
	}

	// check if entry is expired
	// delete entry if expired
	if entry.isExpired() {
		shard.mu.RUnlock()

		shard.mu.Lock()
		delete(shard.entries, key)
		shard.mu.Unlock()

		return nil, ErrKeyDoesntExists
	}

	shard.mu.RUnlock()

	return slices.Clone(entry.value), nil
}

// Exists checks whether 'value' exists for 'key'
func (c *Cache) Exists(key string) bool {
	shard := c.getShard(key)

	shard.mu.RLock()

	_, exists := shard.entries[key]

	// check for expiry
	// delete entry if expired
	if exists && shard.entries[key].isExpired() {
		shard.mu.RUnlock()

		shard.mu.Lock()
		delete(shard.entries, key)
		shard.mu.Unlock()

		return false
	}

	shard.mu.RUnlock()

	return exists
}

// Delete a 'key:value' pair
func (c *Cache) Delete(key string) error {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.entries[key]; exists {
		delete(shard.entries, key)

		return nil
	}

	return ErrKeyDoesntExists
}
