package u64

import "sort"

func Sort(l []uint64) {
	sort.Sort(Uint64Slice(l))
}

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Contain(poses []uint64, pos uint64) bool {
	for _, p := range poses {
		if p == pos {
			return true
		}
	}

	return false
}
