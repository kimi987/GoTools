package random

import "math/rand"

var intCount1 = []int{0}

func NewIntIndexArray(count int) []int {
	switch count {
	case 0:
		return nil
	case 1:
		return intCount1
	}

	a := make([]int, count)
	for i := 0; i < count; i++ {
		a[i] = i
	}

	MixIntArray(a)
	return a
}

func MixIntArray(list []int) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		list[idx], list[swap] = list[swap], list[idx]
	}
}

var u64Count1 = []uint64{0}

func NewU64IndexArray(count uint64) []uint64 {
	switch count {
	case 0:
		return nil
	case 1:
		return u64Count1
	}

	a := make([]uint64, count)
	for i := uint64(0); i < count; i++ {
		a[i] = i
	}

	MixU64Array(a)
	return a
}

func MixU64Array(list []uint64) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		list[idx], list[swap] = list[swap], list[idx]
	}
}

func MixI64Array(list []int64) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		list[idx], list[swap] = list[swap], list[idx]
	}
}

func MixStrArray(list []string) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		list[idx], list[swap] = list[swap], list[idx]
	}
}

// 从M中选N个
func NewMNIntIndexArray(m, n int) []int {
	if m <= n {
		return NewIntIndexArray(m)
	}

	if m < n*2 {
		// 要取一半以上的值，直接创建一个新的值出来
		return NewIntIndexArray(m)[:n]
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		x := rand.Intn(m)

		if !isIntExist(a, i, x) {
			a[i] = x
			continue
		}

		for j := 1; j <= i; j++ {
			// 向前向后找空位
			nx := x - j
			if nx < 0 {
				nx += m
			}

			if !isIntExist(a, i, nx) {
				a[i] = nx
				break
			}

			nx = x + j
			if nx >= m {
				nx -= m
			}
			if !isIntExist(a, i, nx) {
				a[i] = nx
				break
			}
		}
	}

	return a
}

func isIntExist(a []int, n, x int) bool {
	for i := 0; i < n; i++ {
		if x == a[i] {
			return true
		}
	}
	return false
}
