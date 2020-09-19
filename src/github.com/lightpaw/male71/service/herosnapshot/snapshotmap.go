package herosnapshot

import (
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type snapshotmap [32]*snapshotmapShard

// A "thread" safe string to anything map.
type snapshotmapShard struct {
	items        map[int64]*snapshotdata.HeroSnapshot
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func Newsnapshotmap() *snapshotmap {
	var m snapshotmap
	for i := 0; i < 32; i++ {
		m[i] = &snapshotmapShard{items: make(map[int64]*snapshotdata.HeroSnapshot)}
	}
	return &m
}

// Returns shard under given key
func (m *snapshotmap) getShard(key int64) *snapshotmapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *snapshotmap) Set(key int64, value *snapshotdata.HeroSnapshot) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()

	if old, has := shard.items[key]; has && old.Version() >= value.Version() {
		shard.Unlock()
		return
	}

	shard.items[key] = value
	shard.Unlock()
}

// Only update the value if value is already present, and value has newer version
// Return true if value is present, no matter the version
func (m *snapshotmap) SetIfPresent(key int64, value *snapshotdata.HeroSnapshot) (currentValue *snapshotdata.HeroSnapshot, oldHas bool) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()

	if old, has := shard.items[key]; has {
		if old.Version() < value.Version() {
			shard.items[key] = value
			currentValue = value
		} else {
			currentValue = old
		}

		shard.Unlock()
		oldHas = true
		return
	}

	shard.Unlock()
	return nil, false
}

// Sets the given value under the specified key if no value was associated with it.
// Returns old value if present. if ok, old is nil, new value is set. if !ok, old is the previous value, new value is not set.
func (m *snapshotmap) SetIfAbsent(key int64, value *snapshotdata.HeroSnapshot) (*snapshotdata.HeroSnapshot, bool) {
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
func (m *snapshotmap) Get(key int64) (*snapshotdata.HeroSnapshot, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *snapshotmap) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *snapshotmap) Has(key int64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *snapshotmap) Remove(key int64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *snapshotmap) Pop(key int64) (v *snapshotdata.HeroSnapshot, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *snapshotmap) RemoveIfSame(key int64, value *snapshotdata.HeroSnapshot) (ok bool) {
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
func (m *snapshotmap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTuplesnapshotmap struct {
	Key int64
	Val *snapshotdata.HeroSnapshot
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *snapshotmap) Iter() chan []MapTuplesnapshotmap {
	ch := make(chan []MapTuplesnapshotmap, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *snapshotmapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTuplesnapshotmap, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTuplesnapshotmap{key, val}
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
func (m *snapshotmap) Items() map[int64]*snapshotdata.HeroSnapshot {
	tmp := make(map[int64]*snapshotdata.HeroSnapshot, 32)

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
type IterCbsnapshotmap func(key int64, v *snapshotdata.HeroSnapshot)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *snapshotmap) IterCb(fn IterCbsnapshotmap) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *snapshotmap) Keys() chan []int64 {
	ch := make(chan []int64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *snapshotmapShard) {
				// Foreach key, value pair.
				shard.RLock()
				array := make([]int64, len(shard.items))
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
