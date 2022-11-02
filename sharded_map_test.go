package shardmap

import "testing"

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
