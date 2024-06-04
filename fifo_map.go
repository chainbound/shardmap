package shardmap

type FIFOMap[K hashable, V any] struct {
	internalMap *ShardedMap[K, V]
	queue       chan K
	maxSize     int
	currentSize int
}

// NewFIFOMap returns a FIFO map that internally uses the sharded map. It keeps count of all the items
// inserted, and when we exceed the size, values will be evicted using the FIFO (first in, first out) policy.
func NewFIFOMap[K hashable, V any](size, shards int, hashFn HashFn[K]) *FIFOMap[K, V] {
	return &FIFOMap[K, V]{
		internalMap: NewShardedMap[K, V](size, shards, hashFn),
		queue:       make(chan K, size*2),
		maxSize:     size,
		currentSize: 0,
	}
}

func (m *FIFOMap[K, V]) Get(key K) (V, bool) {
	return m.internalMap.Get(key)
}

func (m *FIFOMap[K, V]) Put(key K, val V) {
	if !m.internalMap.Has(key) {
		// If we're about to exceed max size, remove first value from the map
		if m.currentSize >= m.maxSize {
			f := <-m.queue
			m.internalMap.Del(f)
			m.currentSize--
		}

		m.internalMap.Put(key, val)
		m.queue <- key
		m.currentSize++
	} else {
		m.internalMap.Put(key, val)
	}
}

func (m *FIFOMap[K, V]) Has(key K) bool {
	return m.internalMap.Has(key)
}

func (m *FIFOMap[K, V]) Del(key K) {
	if m.internalMap.Has(key) {
		m.internalMap.Del(key)
		m.currentSize--
	}
}

func (m *FIFOMap[K, V]) Len() int {
	return m.currentSize
}

func (m *FIFOMap[K, V]) Keys() []K {
	return m.internalMap.Keys()
}

func (m *FIFOMap[K, V]) Iter() <-chan KVPair[K, V] {
	return m.internalMap.Iter()
}
