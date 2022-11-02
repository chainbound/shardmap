package shardmap

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	sm := NewShardedMap[int, int](1000, 10, HashInt)
	sm.Put(1, 10)

	if !sm.Has(1) {
		t.Fail()
	}

	if v, ok := sm.Get(1); ok {
		if v != 10 {
			t.Fail()
		}
	}

}

func BenchmarkSyncMapPut(b *testing.B) {
	var m sync.Map
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
}

func BenchmarkSyncMapGet(b *testing.B) {
	var m sync.Map
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()
	var v int
	for i := 0; i < b.N; i++ {
		if val, ok := m.Load(i); ok {
			v = val.(int)
		}
	}

	_ = v
}

func BenchmarkShardMapPut(b *testing.B) {
	sm := NewShardedMap[int, int](b.N, 1, HashInt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Put(i, i)
	}
}

func BenchmarkShardMapGet(b *testing.B) {
	sm := NewShardedMap[int, int](b.N, 1, HashInt)
	for i := 0; i < b.N; i++ {
		sm.Put(i, i)
	}
	b.ResetTimer()
	var v int
	for i := 0; i < b.N; i++ {
		v, _ = sm.Get(i)
	}

	_ = v
}
