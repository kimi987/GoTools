package i64

import "math/rand"

// 打乱 int64 数组
func Shuffle(array []int64) {
	l := len(array)
	for i := 0; i < l; i++ {
		swapIndex := rand.Intn(l)
		Swap(array, i, swapIndex)
	}
}

// 交换 int64 数组中数据的位置，不确保i，j越界，谁调用谁检查
func Swap(array []int64, i, j int) {
	array[i], array[j] = array[j], array[i]
}
