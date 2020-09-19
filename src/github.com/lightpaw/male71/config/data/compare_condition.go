package data

import (
	"strings"
	"strconv"
	"github.com/pkg/errors"
	"sort"
)

func ParseCompareCondition(data string) (*CompareCondition, error) {

	c := &CompareCondition{}

	original := data
	data, c.Greater = removeSuccess(data, ">")
	data, c.Less = removeSuccess(data, "<")
	data, c.Equal = removeSuccess(data, "=")

	if !c.Greater && !c.Less {
		c.Equal = true
	}

	amount, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "解析CompareCondition失败，%s", original)
	}
	c.Amount = amount

	return c, nil
}

func removeSuccess(str, toRemove string) (string, bool) {
	n := len(str)
	str = strings.Replace(str, toRemove, "", -1)
	return str, n != len(str)
}

//gogen:config
type CompareCondition struct {
	Greater bool

	Less bool

	Equal bool

	Amount uint64 `validator:"uint"`
}

func (c *CompareCondition) Compare(amt uint64) bool {
	if amt > c.Amount {
		return c.Greater
	} else if amt < c.Amount {
		return c.Less
	}

	return c.Equal
}

func NewCompareConditionKV(k *CompareCondition, v interface{}) *CompareConditionKV {
	return &CompareConditionKV{
		K: k,
		V: v,
	}
}

type CompareConditionKV struct {
	K *CompareCondition

	V interface{}
}

type CompareConditionKVSlice []*CompareConditionKV

func (p CompareConditionKVSlice) Len() int           { return len(p) }
func (p CompareConditionKVSlice) Less(i, j int) bool { return p[i].K.Amount < p[j].K.Amount }
func (p CompareConditionKVSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortCompareCondition(array []*CompareConditionKV) error {
	if len(array) <= 0 {
		return nil
	}

	isGreater := array[0].K.Greater
	isLess := array[0].K.Less
	if isGreater == isLess {
		return errors.Errorf("CompareCondition数组不能同时配置 > 和 <，也不能只配置 =")
	}

	for i, v := range array {
		if isGreater != v.K.Greater {
			return errors.Errorf("CompareCondition数组不能同时配置 > 和 <，也不能只配置 =")
		}

		for j := i + 1; j < len(array); j++ {
			if v.K.Amount == array[j].K.Amount {
				return errors.Errorf("CompareCondition数组不能同时配置 > 和 <，也不能只配置 =")
			}
		}
	}

	if isGreater {
		// 倒序排序，满足条件中最大的
		sortCompareCondition(array, false)
	} else {
		// 比小，顺序排序，满足条件中最小的
		sortCompareCondition(array, true)
	}
	return nil
}

func sortCompareCondition(kvs []*CompareConditionKV, asc bool) {
	if asc {
		sort.Sort(CompareConditionKVSlice(kvs))
	} else {
		sort.Sort(sort.Reverse(CompareConditionKVSlice(kvs)))
	}
}
