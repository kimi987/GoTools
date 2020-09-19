package weight

import (
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"github.com/pkg/errors"
	"math/rand"
	"sort"
)

func NewWeightRandomer(weight []uint64) (*WeightRandomer, error) {

	if len(weight) == 0 {
		return nil, errors.Errorf("创建权重随机器，权重列表长度为0")
	}

	array := make([]uint64, len(weight))
	maxIndex := len(array) - 1

	var totalWeight uint64
	for i, x := range weight {
		if x == 0 {
			return nil, errors.Errorf("创建权重随机器失败，权重值不能为0，rankArray[%v] == 0", i)
		}

		array[maxIndex-i] = totalWeight // reverse
		totalWeight += x
	}

	r := &WeightRandomer{
		totalWeight: u64.Int(totalWeight),
		maxIndex:    maxIndex,
		n:           len(array),
		rankArray:   array,
	}

	return r, nil
}

func NewRankRandomer(rankArray []uint64) (*WeightRandomer, error) {

	if len(rankArray) == 0 {
		return nil, errors.Errorf("创建权重随机器，rank列表长度为0")
	}

	array := make([]uint64, len(rankArray))
	maxIndex := len(array) - 1

	var prev uint64
	for i, x := range rankArray {
		if x <= prev {
			return nil, errors.Errorf("创建权重随机器，rank列表必须是从小到大的排列顺序")
		}

		array[maxIndex-i] = x
	}

	r := &WeightRandomer{
		totalWeight: u64.Int(array[0]),
		maxIndex:    maxIndex,
		n:           len(array),
		rankArray:   array,
	}

	return r, nil
}

type WeightRandomer struct {
	totalWeight int
	maxIndex    int
	n           int
	rankArray   []uint64
}

func (r *WeightRandomer) RandomIndex() int {
	return r.Index(uint64(rand.Intn(r.totalWeight)))
}

func (r *WeightRandomer) Index(x uint64) int {
	return imath.Max(0, r.maxIndex-sort.Search(r.n, func(i int) bool { return r.rankArray[i] <= x }))
}
