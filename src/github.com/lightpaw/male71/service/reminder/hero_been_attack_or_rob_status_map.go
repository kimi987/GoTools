package reminder

import (
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type hero_been_attack_or_rob_status_map [32]*hero_been_attack_or_rob_status_mapShard

// A "thread" safe string to anything map.
type hero_been_attack_or_rob_status_mapShard struct {
	items        map[int64]*hero_been_attack_or_rob_status
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func Newhero_been_attack_or_rob_status_map() *hero_been_attack_or_rob_status_map {
	var m hero_been_attack_or_rob_status_map
	for i := 0; i < 32; i++ {
		m[i] = &hero_been_attack_or_rob_status_mapShard{items: make(map[int64]*hero_been_attack_or_rob_status)}
	}
	return &m
}

// Returns shard under given key
func (m *hero_been_attack_or_rob_status_map) getShard(key int64) *hero_been_attack_or_rob_status_mapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *hero_been_attack_or_rob_status_map) Set(key int64, value *hero_been_attack_or_rob_status) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCbhero_been_attack_or_rob_status_map func(exist bool, valueInMap *hero_been_attack_or_rob_status, newValue *hero_been_attack_or_rob_status) *hero_been_attack_or_rob_status

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *hero_been_attack_or_rob_status_map) Upsert(key int64, value *hero_been_attack_or_rob_status, cb UpsertCbhero_been_attack_or_rob_status_map) (res *hero_been_attack_or_rob_status) {
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
func (m *hero_been_attack_or_rob_status_map) SetIfAbsent(key int64, value *hero_been_attack_or_rob_status) (*hero_been_attack_or_rob_status, bool) {
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
func (m *hero_been_attack_or_rob_status_map) Get(key int64) (*hero_been_attack_or_rob_status, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *hero_been_attack_or_rob_status_map) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *hero_been_attack_or_rob_status_map) Has(key int64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *hero_been_attack_or_rob_status_map) Remove(key int64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *hero_been_attack_or_rob_status_map) Pop(key int64) (v *hero_been_attack_or_rob_status, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *hero_been_attack_or_rob_status_map) RemoveIfSame(key int64, value *hero_been_attack_or_rob_status) (ok bool) {
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
func (m *hero_been_attack_or_rob_status_map) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTuplehero_been_attack_or_rob_status_map struct {
	Key int64
	Val *hero_been_attack_or_rob_status
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *hero_been_attack_or_rob_status_map) Iter() chan []MapTuplehero_been_attack_or_rob_status_map {
	ch := make(chan []MapTuplehero_been_attack_or_rob_status_map, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *hero_been_attack_or_rob_status_mapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTuplehero_been_attack_or_rob_status_map, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTuplehero_been_attack_or_rob_status_map{key, val}
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
func (m *hero_been_attack_or_rob_status_map) Items() map[int64]*hero_been_attack_or_rob_status {
	tmp := make(map[int64]*hero_been_attack_or_rob_status, 32)

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
type IterCbhero_been_attack_or_rob_status_map func(key int64, v *hero_been_attack_or_rob_status)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *hero_been_attack_or_rob_status_map) IterCb(fn IterCbhero_been_attack_or_rob_status_map) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *hero_been_attack_or_rob_status_map) Keys() chan []int64 {
	ch := make(chan []int64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *hero_been_attack_or_rob_status_mapShard) {
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
