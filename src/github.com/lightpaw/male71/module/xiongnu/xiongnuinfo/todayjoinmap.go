package xiongnuinfo

import (
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type TodayJoinMap [32]*todayJoinMapShard

// A "thread" safe string to anything map.
type todayJoinMapShard struct {
	items        map[int64]int64
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewTodayJoinMap() *TodayJoinMap {
	var m TodayJoinMap
	for i := 0; i < 32; i++ {
		m[i] = &todayJoinMapShard{items: make(map[int64]int64)}
	}
	return &m
}

// Returns shard under given key
func (m *TodayJoinMap) getShard(key int64) *todayJoinMapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *TodayJoinMap) Set(key int64, value int64) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCbTodayJoinMap func(exist bool, valueInMap int64, newValue int64) int64

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *TodayJoinMap) Upsert(key int64, value int64, cb UpsertCbTodayJoinMap) (res int64) {
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
func (m *TodayJoinMap) SetIfAbsent(key int64, value int64) (int64, bool) {
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
func (m *TodayJoinMap) Get(key int64) (int64, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *TodayJoinMap) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *TodayJoinMap) Has(key int64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *TodayJoinMap) Remove(key int64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *TodayJoinMap) Pop(key int64) (v int64, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *TodayJoinMap) RemoveIfSame(key int64, value int64) (ok bool) {
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
func (m *TodayJoinMap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleTodayJoinMap struct {
	Key int64
	Val int64
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *TodayJoinMap) Iter() chan []MapTupleTodayJoinMap {
	ch := make(chan []MapTupleTodayJoinMap, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *todayJoinMapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleTodayJoinMap, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleTodayJoinMap{key, val}
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
func (m *TodayJoinMap) Items() map[int64]int64 {
	tmp := make(map[int64]int64, 32)

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
type IterCbTodayJoinMap func(key int64, v int64)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *TodayJoinMap) IterCb(fn IterCbTodayJoinMap) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *TodayJoinMap) Keys() chan []int64 {
	ch := make(chan []int64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *todayJoinMapShard) {
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
