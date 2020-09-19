package u64

import (
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type StringU64Map [32]*stringU64MapShard

// A "thread" safe string to anything map.
type stringU64MapShard struct {
	items        map[string]uint64
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewStringU64Map() *StringU64Map {
	var m StringU64Map
	for i := 0; i < 32; i++ {
		m[i] = &stringU64MapShard{items: make(map[string]uint64)}
	}
	return &m
}

// Returns shard under given key
func (m *StringU64Map) getShard(key string) *stringU64MapShard {

	var out uint64
	ln := uint64(len(key) - 1)
	if ln > 1023 || ln < 0 {
		return m[0]
	}
	for i := uint64(0); i < ln; i++ {
		out += uint64(key[i]) + 1024*i
	}
	return m[out&31]

}

// Sets the given value under the specified key.
func (m *StringU64Map) Set(key string, value uint64) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCbStringU64Map func(exist bool, valueInMap uint64, newValue uint64) uint64

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *StringU64Map) Upsert(key string, value uint64, cb UpsertCbStringU64Map) (res uint64) {
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
func (m *StringU64Map) SetIfAbsent(key string, value uint64) (uint64, bool) {
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
func (m *StringU64Map) Get(key string) (uint64, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *StringU64Map) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *StringU64Map) Has(key string) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *StringU64Map) Remove(key string) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *StringU64Map) Pop(key string) (v uint64, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *StringU64Map) RemoveIfSame(key string, value uint64) (ok bool) {
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
func (m *StringU64Map) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleStringU64Map struct {
	Key string
	Val uint64
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *StringU64Map) Iter() chan []MapTupleStringU64Map {
	ch := make(chan []MapTupleStringU64Map, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *stringU64MapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleStringU64Map, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleStringU64Map{key, val}
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
func (m *StringU64Map) Items() map[string]uint64 {
	tmp := make(map[string]uint64, 32)

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
type IterCbStringU64Map func(key string, v uint64)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *StringU64Map) IterCb(fn IterCbStringU64Map) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *StringU64Map) Keys() chan []string {
	ch := make(chan []string, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *stringU64MapShard) {
				// Foreach key, value pair.
				shard.RLock()
				array := make([]string, len(shard.items))
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
