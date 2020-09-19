package concurrent

import (
	"github.com/lightpaw/pbutil"
)

type U64BufferMap interface {
	GetBuffer(key uint64) (buffer pbutil.Buffer, err error)
	GetVersionBuffer(key uint64) (version uint64, buffer pbutil.Buffer, err error)
	GetI32VersionBuffer(key uint64) (version int32, buffer pbutil.Buffer, err error)
	Clear(key uint64)
}
