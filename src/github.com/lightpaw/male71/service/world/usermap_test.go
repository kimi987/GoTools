package world

import (
	"github.com/lightpaw/male7/gen/iface"
	"testing"
)

var userMap = Newusermap()
var hashMap = map[int64]iface.ConnectedUser{}
var array = make([]int64, 10000)

func init() {
	for i := int64(0); i < 10000; i++ {
		userMap.Set(i, nil)
		hashMap[i] = nil
		array[i] = i
	}
}

func getHashMap(key int64) iface.ConnectedUser {
	return hashMap[key]
}

func getArray(key int) int64 {
	if key >= 0 && key < len(array) {
		return array[key]
	}
	return 0
}

func BenchmarkGetUserMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		userMap.Get(int64(i))
	}
}

func BenchmarkGetHashMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHashMap(int64(i))
	}
}

func BenchmarkGetArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getArray(i % 10000)
	}
}
