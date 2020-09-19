package idbytes

import (
	"github.com/lightpaw/male7/util/i64"
)

func ToBytes(id int64) []byte {
	return i64.ToBytes(id)
}

func ToBytesArray(ids []int64) [][]byte {
	bufs := make([][]byte, len(ids))
	for i, id := range ids {
		bufs[i] = i64.ToBytes(id)
	}
	return bufs
}

func ToId(buf []byte) (id int64, ok bool) {
	return i64.FromBytes(buf)
}

func ToIds(bufs [][]byte) ([]int64, bool) {
	ids := make([]int64, len(bufs))

	for i, buf := range bufs {
		if id, ok := i64.FromBytes(buf); ok {
			ids[i] = id
			continue
		}
		return nil, false
	}

	return ids, true
}

func DefId(buf []byte, def int64) int64 {
	if id, ok := ToId(buf); ok {
		return id
	}

	return def
}

func IsPositive(buf []byte) bool {
	if len(buf) <= 0 {
		return true
	}

	return buf[0]&1 == 0
}
