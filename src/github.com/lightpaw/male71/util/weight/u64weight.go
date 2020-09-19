package weight

import "github.com/pkg/errors"

func NewU64WeightRandomer(weight, values []uint64) (*U64WeightRandomer, error) {

	if len(weight) != len(values) {
		return nil, errors.Errorf("创建U64权重随机器，权重个数[%v]跟配置个数[%v]不相同，Weight: %v, Value: %v", len(weight), len(values), weight, values)
	}

	w, err := NewWeightRandomer(weight)
	if err != nil {
		return nil, errors.Wrapf(err, "创建U64权重随机器，Value: %v", values)
	}

	return &U64WeightRandomer{weight: w, values: values}, nil
}

func NewU64RankRandomer(rankArray, values []uint64) (*U64WeightRandomer, error) {

	if len(rankArray) != len(values) {
		return nil, errors.Errorf("创建U64权重随机器，权重个数[%v]跟配置个数[%v]不相同，Weight: %v, Value: %v", len(rankArray), len(values), rankArray, values)
	}

	w, err := NewRankRandomer(rankArray)
	if err != nil {
		return nil, errors.Wrapf(err, "创建U64权重随机器，Value: %v", values)
	}

	return &U64WeightRandomer{weight: w, values: values}, nil
}

type U64WeightRandomer struct {
	weight *WeightRandomer

	values []uint64
}

func (r *U64WeightRandomer) Random() uint64 {
	return r.values[r.weight.RandomIndex()]
}

func (r *U64WeightRandomer) Get(x uint64) uint64 {
	return r.values[r.weight.Index(x)]
}
