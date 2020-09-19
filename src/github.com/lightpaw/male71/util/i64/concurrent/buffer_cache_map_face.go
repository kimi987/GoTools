package concurrent

import (
	"github.com/lightpaw/pbutil"
)

type I64BufferMap interface {
	GetBuffer(key int64) (buffer pbutil.Buffer, err error)
	GetVersionBuffer(key int64) (version uint64, buffer pbutil.Buffer, err error)
	GetI32VersionBuffer(key int64) (version int32, buffer pbutil.Buffer, err error)
	Clear(key int64)
}
