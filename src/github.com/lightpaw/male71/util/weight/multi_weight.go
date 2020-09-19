package weight

import (
	"github.com/pkg/errors"
	"math/rand"
)

var (
	ErrChoiceNotEnough = errors.Errorf("weight.len < n")
	ErrZeroWeight      = errors.Errorf("weight is 0")
	ErrLogic           = errors.Errorf("error logic")
)

func RandomN(weight []uint64, n int) ([]int, error) {

	if n <= 0 {
		return nil, nil
	}

	count := len(weight)
	if count <= 0 || count < n {
		return nil, ErrChoiceNotEnough
	}

	if count == n {
		a := make([]int, count)
		for i := 0; i < count; i++ {
			a[i] = i
		}
		return a, nil
	}

	// 循环n次，每次随机一个出来，然后
	var totalWeight uint64
	for _, w := range weight {
		if w <= 0 {
			return nil, ErrZeroWeight
		}
		totalWeight += w
	}

	copyWeight := make([]uint64, count)
	copy(copyWeight, weight)

	index := make([]int, 0, n)

out:
	for i := 0; i < n; i++ {
		x := rand.Uint64() % totalWeight

		// 从第一个开始找起，找到第一个 cw <= x < cw+w
		var cur uint64
		for i := 0; i < count; i++ {

			w := copyWeight[i]

			if w > 0 && x < cur+w {
				totalWeight -= w
				copyWeight[i] = 0

				index = append(index, i)
				continue out
			}

			cur += w
		}

		return nil, ErrLogic
	}

	return index, nil
}
