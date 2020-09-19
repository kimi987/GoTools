package guildsnapshot

import (
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"sync"
)

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.
type SnapshotMap [32]*snapshotMapShard

// A "thread" safe string to anything map.
type snapshotMapShard struct {
	items        map[int64]*guildsnapshotdata.GuildSnapshot
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewSnapshotMap() *SnapshotMap {
	var m SnapshotMap
	for i := 0; i < 32; i++ {
		m[i] = &snapshotMapShard{items: make(map[int64]*guildsnapshotdata.GuildSnapshot)}
	}
	return &m
}

// Returns shard under given key
func (m *SnapshotMap) getShard(key int64) *snapshotMapShard {
	return m[key&31]
}

// Sets the given value under the specified key.
func (m *SnapshotMap) Set(key int64, value *guildsnapshotdata.GuildSnapshot) *guildsnapshotdata.GuildSnapshot {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	old := shard.items[key]
	shard.items[key] = value
	shard.Unlock()

	return old
}

type UpsertCbSnapshotMap func(exist bool, valueInMap *guildsnapshotdata.GuildSnapshot, newValue *guildsnapshotdata.GuildSnapshot) *guildsnapshotdata.GuildSnapshot

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *SnapshotMap) Upsert(key int64, value *guildsnapshotdata.GuildSnapshot, cb UpsertCbSnapshotMap) (res *guildsnapshotdata.GuildSnapshot) {
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
func (m *SnapshotMap) SetIfAbsent(key int64, value *guildsnapshotdata.GuildSnapshot) (*guildsnapshotdata.GuildSnapshot, bool) {
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
func (m *SnapshotMap) Get(key int64) (*guildsnapshotdata.GuildSnapshot, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *SnapshotMap) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Looks up an item under specified key
func (m *SnapshotMap) Has(key int64) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *SnapshotMap) Remove(key int64) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Removes an element from the map and returns it
func (m *SnapshotMap) Pop(key int64) (v *guildsnapshotdata.GuildSnapshot, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// Only removes an element if it exists and the value is the same
func (m *SnapshotMap) RemoveIfSame(key int64, value *guildsnapshotdata.GuildSnapshot) (ok bool) {
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
func (m *SnapshotMap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleSnapshotMap struct {
	Key int64
	Val *guildsnapshotdata.GuildSnapshot
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *SnapshotMap) Iter() chan []MapTupleSnapshotMap {
	ch := make(chan []MapTupleSnapshotMap, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *snapshotMapShard) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleSnapshotMap, len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleSnapshotMap{key, val}
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
func (m *SnapshotMap) Items() map[int64]*guildsnapshotdata.GuildSnapshot {
	tmp := make(map[int64]*guildsnapshotdata.GuildSnapshot, 32)

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
type IterCbSnapshotMap func(key int64, v *guildsnapshotdata.GuildSnapshot)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *SnapshotMap) IterCb(fn IterCbSnapshotMap) {
	for _, shard := range m {
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys
func (m *SnapshotMap) Keys() chan []int64 {
	ch := make(chan []int64, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *snapshotMapShard) {
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
