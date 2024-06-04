# shardmap

> A performant, highly concurrent and simple sharded hashmap implementation using generics.

This package contains a `ShardedMap` and a `FIFOMap`.

## `ShardedMap`
A `ShardedMap` is a simple map that uses a sharded design. A sharded map splits the map into
buckets or shards according to the key, where each shard has its own `RWMutex`. This ensures that lock contention
is heavily minimized compared to using one mutex for the whole map, making it very high throughput and low latency
in highly concurrent situations.

It has the following interface:
```go
type ShardedMapInterface interface {
    Get(key K) (val V, ok bool)
    Put(key K, val V)
    Has(key K) ok bool
    Del(key K)
    Keys() []K
    Iter() <-chan KVPair[K, V]
    Len() int
}
```
### Example
```go
import "github.com/chainbound/shardmap"

func main() {
    // Initialize a new sharded int -> string map with size 1000, and 10 shards.
    // We need to provide the hash function for our key type, the defaults being contained
    // in this package. You can also provide your own.
    sm := shardmap.NewShardedMap[int, string](1000, 10, shardmap.HashInt)
    sm.Put(1, "josh")

    fmt.Println(sm.Get(1))
}
```

## `FIFOMap`
The `FIFOMap` is a map with a FIFO eviction policy, meaning that the oldest values get removed once your map
reaches a certain `size`. Internally, it uses the sharded map above, and shares the same interface.
```go
import "github.com/chainbound/shardmap"

func main() {
    // Initialize a new sharded int -> string map with size 1000, and 10 shards.
    // We need to provide the hash function for our key type, the defaults being contained
    // in this package. You can also provide your own.
    sm := shardmap.NewFIFOMap[int, string](1000, 10, shardmap.HashInt)

    // Once the size is reached, the next put will remove the oldest inserted KV pair.
    sm.Put(1, "josh")

    fmt.Println(sm.Get(1))
}
```