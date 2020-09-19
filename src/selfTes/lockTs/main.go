package main

import "sync"

var nums map[int]int
var wg sync.WaitGroup
var rwLock sync.RWMutex

func main() {
	nums = make(map[int]int)
	wg.Add(2)
	go WriteMap()
	go ReadMap()

	wg.Wait()
}
func WriteMap() {
	rwLock.RLock()
	defer rwLock.RUnlock()
	for i := 0; i < 10000; i++ {
		nums[i] = 1
	}
	wg.Done()
}
func ReadMap() int {
	rwLock.Lock()
	defer rwLock.Unlock()
	num := 0
	for index := 0; index < 10000; index++ {
		num += nums[index]
	}

	wg.Done()
	return num
}
