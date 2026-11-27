package main

import "testing"

func TestGetWithTTL(t *testing.T) {
	// test for lazy expiry
	t.Fail()
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
