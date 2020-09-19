package sortkeys

import (
	"testing"
	"time"
	"math/rand"
	"sort"
	. "github.com/onsi/gomega"
)

func TestIntKv(t *testing.T) {
	RegisterTestingT(t)

	keys := []int{1, 4, 5, 2, 3, 7}
	sorted := make([]int, len(keys))
	copy(sorted, keys)
	sort.Ints(sorted)

	kvs := make([]*IntKV, len(keys))
	for i, k := range keys {
		kvs[i] = &IntKV{K: k}
	}

	sort.Sort(IntKVSlice(kvs))

	for i, k := range sorted {
		Ω(kvs[i].K).Should(Equal(k))
	}
}

func TestKv(t *testing.T) {
	RegisterTestingT(t)

	keys := []int{1, 4, 5, 2, 3, 7}
	sorted := make([]int, len(keys))
	copy(sorted, keys)
	sort.Ints(sorted)

	kvs := make([]*KV, len(keys))
	for i, k := range keys {
		kvs[i] = &KV{K: k}
	}

	sort.Sort(KVIntSlice(kvs))

	for i, k := range sorted {
		Ω(kvs[i].K).Should(Equal(k))
	}
}

var seed = time.Now().UnixNano()

func randKeys(count int) []int {
	rand.Seed(seed)

	var keys []int
	for i := 0; i < count; i++ {
		keys = append(keys, rand.Int())
	}
	return keys
}

func BenchmarkIntKV(b *testing.B) {

	keys := randKeys(10000)
	kvs := make([]*IntKV, len(keys))
	for i, k := range keys {
		kvs[i] = &IntKV{K: k}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort.Sort(IntKVSlice(kvs))
		for i, k := range keys {
			kvs[i].K = k
		}
	}
}

func BenchmarkCKV(b *testing.B) {

	keys := randKeys(10000)
	kvs := make([]*KV, len(keys))
	for i, k := range keys {
		kvs[i] = &KV{K: k}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort.Sort(KVIntSlice(kvs))
		for i, k := range keys {
			kvs[i].K = k
		}
	}
}

type KVIntSlice []*KV

func (a KVIntSlice) Len() int           { return len(a) }
func (a KVIntSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a KVIntSlice) Less(i, j int) bool { return a[i].IntKey() < a[j].IntKey() }
