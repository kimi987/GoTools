package entity

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/pb/shared_proto"
	"math"
)

func TestHeroBools(t *testing.T) {
	RegisterTestingT(t)

	n := len(shared_proto.HeroBoolType_name)
	for k := range shared_proto.HeroBoolType_name {
		Ω(int(k) < n).Should(BeTrue())
	}

	//if n >= 32 {
	//	// 客户端那边使用了1个int32来发送，因此，数据最多31位，超出需要改代码
	//	t.Fatal("客户端那边使用了1个int32来发送，因此，数据最多31位，超出需要改代码")
	//}

	bools := newBools()
	Ω(bools.serverBools).Should(BeEmpty())
	Ω(bools.encodeClient()).Should(Equal(int32(0)))

	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())

	bools.SetFalse(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())
	Ω(bools.serverBools).Should(Equal([]uint64{0}))
	Ω(bools.encodeClient()).Should(Equal(int32(0)))

	bools.SetTrue(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeTrue())
	Ω(bools.serverBools).Should(Equal([]uint64{1}))
	Ω(bools.encodeClient()).Should(Equal(int32(1)))

	bools.SetFalse(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())
	Ω(bools.serverBools).Should(Equal([]uint64{0}))

	bools.SetFalse(shared_proto.HeroBoolType(1))
	Ω(bools.serverBools).Should(Equal([]uint64{0}))

	bools.SetFalse(shared_proto.HeroBoolType(65))
	Ω(bools.serverBools).Should(Equal([]uint64{0}))

	bools.serverBools = []uint64{64, 20}
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())

	bools.SetFalse(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())
	Ω(bools.serverBools).Should(Equal([]uint64{64, 20}))

	bools.SetTrue(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeTrue())
	Ω(bools.serverBools).Should(Equal([]uint64{65, 20}))

	bools.SetFalse(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
	Ω(bools.Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)).Should(BeFalse())
	Ω(bools.serverBools).Should(Equal([]uint64{64, 20}))

	bools.SetFalse(shared_proto.HeroBoolType(1))
	Ω(bools.serverBools).Should(Equal([]uint64{64, 20}))

	bools.SetFalse(shared_proto.HeroBoolType(65))
	Ω(bools.serverBools).Should(Equal([]uint64{64, 20}))

	bools.serverBools[0] = math.MaxInt32
	Ω(bools.encodeClient()).Should(Equal(int32(math.MaxInt32)))

	bools.serverBools[0] = math.MaxInt64
	Ω(bools.encodeClient()).Should(Equal(int32(math.MaxInt32)))

	bools.serverBools[0] = math.MaxUint32
	Ω(bools.encodeClient()).Should(Equal(int32(math.MaxInt32)))

	bools.serverBools[0] = math.MaxUint64
	Ω(bools.encodeClient()).Should(Equal(int32(math.MaxInt32)))
}
