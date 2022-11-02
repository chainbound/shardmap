package shardmap

import (
	"sync"
)

type ShardedMap[K hashable, V any] struct {
	numShards int
	shards    []*shard[K, V]

	hashFn HashFn[K]
}

type shard[K hashable, V any] struct {
	sync.RWMutex
	internalMap map[K]V
}

type hashable interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

// NewShardedMap returns a new sharded map with `numShards` shards. Each of the shards are pre-allocated
// with a length of `size` / `numShards`. `size` is not the max size by any means, but just an estimation.
// hashFn is used to hash the key.
func NewShardedMap[K hashable, V any](size, numShards int, hashFn HashFn[K]) *ShardedMap[K, V] {
	if numShards < 1 {
		numShards = 0
	}

	m := &ShardedMap[K, V]{
		numShards: numShards,
		shards:    make([]*shard[K, V], numShards),
		hashFn:    hashFn,
	}

	for i := 0; i < numShards; i++ {
		m.shards[i] = &shard[K, V]{
			internalMap: make(map[K]V, size/numShards),
		}
	}

	return m
}

// Get returns the value and true if the value is present, otherwise it returns the default value
// and false.
func (m *ShardedMap[K, V]) Get(key K) (v V, ok bool) {
	shard := m.hashFn(key) & uint64(m.numShards-1)
	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].RLock()
	defer m.shards[shard].RUnlock()

	if v, ok = m.shards[shard].internalMap[key]; ok {
		return
	}

	return
}

// Put puts the key value pair in the map.
func (m *ShardedMap[K, V]) Put(key K, val V) {
	shard := m.hashFn(key) & uint64(m.numShards-1)
	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].Lock()
	defer m.shards[shard].Unlock()

	m.shards[shard].internalMap[key] = val
}

// Has returns true if the key is present.
func (m *ShardedMap[K, V]) Has(key K) bool {
	shard := m.hashFn(key) & uint64(m.numShards-1)
	if m.shards[shard] == nil {
		return false
	}

	m.shards[shard].RLock()
	defer m.shards[shard].RUnlock()

	if _, ok := m.shards[shard].internalMap[key]; ok {
		return true
	}

	return false
}

// Del deletes the value from the map.
func (m *ShardedMap[K, V]) Del(key K) {
	shard := m.hashFn(key) & uint64(m.numShards-1)
	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].Lock()
	defer m.shards[shard].Unlock()

	delete(m.shards[shard].internalMap, key)
}

// Len returns the count of all the items in the sharded map.
// It will RLock every one of the shards so use it scarcely.
func (m *ShardedMap[K, V]) Len() int {
	total := 0

	for _, s := range m.shards {
		s.RLock()
		total += len(s.internalMap)
		s.RUnlock()
	}

	return total
}
