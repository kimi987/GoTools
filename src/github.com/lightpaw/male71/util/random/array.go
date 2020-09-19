package random

import "math/rand"

func NewArray() *Array {
	return &Array{}
}

type Array struct {
	arr []interface{}
}

func (a *Array) Len() int {
	return len(a.arr)
}

func (a *Array) Add(toAdd interface{}) {
	a.arr = append(a.arr, toAdd)
}

func (a *Array) Random() (int, interface{}) {
	n := len(a.arr)
	if n == 0 {
		return -1, nil
	}

	idx := rand.Intn(n)
	return idx, a.arr[idx]
}

func (a *Array) Remove(idx int) {
	n := len(a.arr)
	if idx >= n {
		return
	}

	if idx < n-1 {
		a.arr[idx], a.arr[n-1] = a.arr[n-1], a.arr[idx]
	}

	a.arr = a.arr[:n-1]
}
