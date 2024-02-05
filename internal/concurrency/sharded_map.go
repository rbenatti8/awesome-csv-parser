package concurrency

import (
	"crypto/sha1"
	"sync"
)

type Shard[T any] struct {
	m map[string]T
	sync.RWMutex
}

type ShardedMap[T any] []*Shard[T]

// NewShardedMap creates a new ShardedMap with the given number of shards.
func NewShardedMap[T any](numShards int) ShardedMap[T] {
	shards := make([]*Shard[T], numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &Shard[T]{m: make(map[string]T)}
	}
	return shards
}

func (sm ShardedMap[T]) getShardIndex(key string) int {
	hash := sha1.Sum([]byte(key))

	// Grab an arbitrary byte from the hash and mod it by the number of shards.
	// This is a simple way to get a random-ish shard index.
	return int(hash[17]) % len(sm)
}

func (sm ShardedMap[T]) getShard(key string) *Shard[T] {
	return sm[sm.getShardIndex(key)]
}

// Delete deletes the given key from the map.
func (sm ShardedMap[T]) Delete(key string) {
	shard := sm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

// Get returns the value associated with the given key.
func (sm ShardedMap[T]) Get(key string) (T, bool) {
	shard := sm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	val, ok := shard.m[key]
	return val, ok
}

// Set sets the given key to the given value.
func (sm ShardedMap[T]) Set(key string, val T) {
	shard := sm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = val
}
