package concurrent

import (
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type u64_cache_map [32]*u64_cache_mapShard

// A "thread" safe string to anything map.
type u64_cache_mapShard struct {
	items        map[uint64]*OnceBufferCache
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func Newu64_cache_map() *u64_cache_map {
	var m u64_cache_map
	for i := 0; i < 32; i++ {
		m[i] = &u64_cache_mapShard{items: make(map[uint64]*OnceBufferCache)}
	}
	return &m
}

// Returns shard under given key
func (m *u64_cache_map) getShard(key uint64) *u64_cache_mapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *u64_cache_map) Set(key uint64, value *OnceBufferCache) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCbu64_cache_map func(exist bool, valueInMap *OnceBufferCache, newValue *OnceBufferCache) *OnceBufferCache

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *u64_cache_map) Upsert(key uint64, value *OnceBufferCache, cb UpsertCbu64_cache_map) (res *OnceBufferCache) {
	shard := m.getShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

// Sets the given value under the specified key if no value was associated with it.
// Returns old value if present. if ok, old is nil, new value is set. if !ok, old is the previous value, new value is not set.
func (m *u64_cache_map) SetIfAbsent(key uint64, value *OnceBufferCache) (*OnceBufferCache, bool) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	old, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return old, !ok
}

// Retrieves an element from map under given key.
func (m *u64_cache_map) Get(key uint64) (*OnceBufferCache, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *u64_cache_map) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *u64_cache_map) Has(key uint64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *u64_cache_map) Remove(key uint64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *u64_cache_map) Pop(key uint64) (v *OnceBufferCache, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *u64_cache_map) RemoveIfSame(key uint64, value *OnceBufferCache) (ok bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists := shard.items[key]
	if exists && v == value {
		delete(shard.items, key)
		ok = true
	}
	shard.Unlock()
	return
}

// Checks if map is empty.
func (m *u64_cache_map) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleu64_cache_map struct {
	Key uint64
	Val *OnceBufferCache
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *u64_cache_map) Iter() chan []MapTupleu64_cache_map {
	ch := make(chan []MapTupleu64_cache_map, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *u64_cache_mapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleu64_cache_map, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleu64_cache_map{key, val}
					idx++
				}
				shard.RUnlock()
				ch <- tuples
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}

// Returns all items as map[string]interface{}
func (m *u64_cache_map) Items() map[uint64]*OnceBufferCache {
	tmp := make(map[uint64]*OnceBufferCache, 32)

	// Insert items to temporary map.
	for tp := range m.Iter() {
		for _, item := range tp {
			tmp[item.Key] = item.Val
		}
	}

	return tmp
}

// Iterator callback,called for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCbu64_cache_map func(key uint64, v *OnceBufferCache)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *u64_cache_map) IterCb(fn IterCbu64_cache_map) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *u64_cache_map) Keys() chan []uint64 {
	ch := make(chan []uint64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *u64_cache_mapShard) {
				// Foreach key, value pair.
				shard.RLock()
				array := make([]uint64, len(shard.items))
				idx := 0
				for key := range shard.items {
					array[idx] = key
					idx++
				}
				shard.RUnlock()
				ch <- array
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}
