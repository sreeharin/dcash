package main

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestGetWithTTL(t *testing.T) {
	// test for lazy expiry
	t.Run("lazy expiry", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			cache := NewCache()

			var (
				user      string = "user:100"
				name      []byte = []byte("pavanayi")
				expiresAt int64  = time.Now().UnixNano() + (5 * time.Second).Nanoseconds()
			)

			cache.Set(user, name, expiresAt)

			if _, err := cache.Get(user); err != nil {
				t.Error("Get failed. Key expired before duration.")
			}

			// advance time
			synctest.Sleep(6 * time.Second)

			if _, err := cache.Get(user); err == nil {
				t.Error("Get failed. Key didn't expire after duration.")
			}
		})

	})

	// test for 0 expiry
}

func TestExistsWithTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cache := NewCache()

		var (
			user      string = "user:1001:name"
			value     []byte = []byte("prabhakaran")
			expiresAt int64  = time.Now().UnixNano() + (5 * time.Second).Nanoseconds()
		)

		cache.Set(user, value, expiresAt)

		synctest.Sleep(3 * time.Second)

		if !cache.Exists(user) {
			t.Error("Exists failed. KEY expired before duration.")
		}

		synctest.Sleep(2 * time.Second)

		if cache.Exists(user) {
			t.Error("Exists failed. KEY found after expiry duration.")
		}
	})
}

func TestSETWithTTL(t *testing.T) {
	// SET cmd should update TTL if key already exists
	t.Fail()
}

func TestCulling(t *testing.T) {
	// test to make sure 'culling' removes expired keys from shards
	// check it reaches every shard
	t.Fail()
}
