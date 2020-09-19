package xionnuservice

import (
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type infoMap [32]*infoMapShard

// A "thread" safe string to anything map.
type infoMapShard struct {
	items        map[int64]xiongnuface.ResistXiongNuInfo
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewinfoMap() *infoMap {
	var m infoMap
	for i := 0; i < 32; i++ {
		m[i] = &infoMapShard{items: make(map[int64]xiongnuface.ResistXiongNuInfo)}
	}
	return &m
}

// Returns shard under given key
func (m *infoMap) getShard(key int64) *infoMapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *infoMap) Set(key int64, value xiongnuface.ResistXiongNuInfo) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCbinfoMap func(exist bool, valueInMap xiongnuface.ResistXiongNuInfo, newValue xiongnuface.ResistXiongNuInfo) xiongnuface.ResistXiongNuInfo

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *infoMap) Upsert(key int64, value xiongnuface.ResistXiongNuInfo, cb UpsertCbinfoMap) (res xiongnuface.ResistXiongNuInfo) {
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
func (m *infoMap) SetIfAbsent(key int64, value xiongnuface.ResistXiongNuInfo) (xiongnuface.ResistXiongNuInfo, bool) {
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
func (m *infoMap) Get(key int64) (xiongnuface.ResistXiongNuInfo, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *infoMap) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *infoMap) Has(key int64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *infoMap) Remove(key int64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *infoMap) Pop(key int64) (v xiongnuface.ResistXiongNuInfo, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *infoMap) RemoveIfSame(key int64, value xiongnuface.ResistXiongNuInfo) (ok bool) {
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
func (m *infoMap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleinfoMap struct {
	Key int64
	Val xiongnuface.ResistXiongNuInfo
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *infoMap) Iter() chan []MapTupleinfoMap {
	ch := make(chan []MapTupleinfoMap, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *infoMapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleinfoMap, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleinfoMap{key, val}
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
func (m *infoMap) Items() map[int64]xiongnuface.ResistXiongNuInfo {
	tmp := make(map[int64]xiongnuface.ResistXiongNuInfo, 32)

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
type IterCbinfoMap func(key int64, v xiongnuface.ResistXiongNuInfo)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *infoMap) IterCb(fn IterCbinfoMap) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *infoMap) Keys() chan []int64 {
	ch := make(chan []int64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *infoMapShard) {
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
