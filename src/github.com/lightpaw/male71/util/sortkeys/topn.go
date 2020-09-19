package sortkeys

import "sort"

func NewU64TopN(n uint64) *U64TopN {
	return &U64TopN{
		n:     n,
		array: make([]*U64KV, 0, n),
	}
}

type U64TopN struct {
	n uint64

	array []*U64KV

	minKey      uint64
	minKeyIndex uint64
}

func (t *U64TopN) Array() []*U64KV {
	return t.array
}

func (t *U64TopN) CopyArray() []*U64KV {
	array := make([]*U64KV, len(t.array))
	copy(array, t.array)
	return array
}

func (t *U64TopN) SortAsc() []*U64KV {
	array := t.CopyArray()
	sort.Sort(U64KVSlice(array))
	return array
}

func (t *U64TopN) SortDesc() []*U64KV {
	array := t.CopyArray()
	sort.Sort(sort.Reverse(U64KVSlice(array)))
	return array
}

func (t *U64TopN) Size() int {
	return len(t.array)
}

func (t *U64TopN) Add(k uint64, v interface{}) {

	oldLen := uint64(len(t.array))
	if oldLen < t.n {

		kv := NewU64KV(k, v)
		t.array = append(t.array, kv)

		if oldLen == 0 || k < t.minKey {
			t.minKey = k
			t.minKeyIndex = oldLen
		}
		return
	}

	// topN已经满了，如果超过最小值，那么跟最小值进行交换，然后重新获取最小值
	if t.minKey < k {
		t.array[t.minKeyIndex] = NewU64KV(k, v)

		for i, v := range t.array {

			if i == 0 || v.K < t.minKey {
				t.minKey = v.K
				t.minKeyIndex = uint64(i)
			}
		}
	}

}
