package depot

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
	_              = shared_proto.ErrIntOverflowConfig
)

type marshaler util.Marshaler

const (
	MODULE_ID = 11

	C2S_USE_GOODS = 2

	C2S_USE_CDR_GOODS = 6

	C2S_GOODS_COMBINE = 9

	C2S_GOODS_PARTS_COMBINE = 18

	C2S_UNLOCK_BAOWU = 23

	C2S_COLLECT_BAOWU = 26

	C2S_LIST_BAOWU_LOG = 30

	C2S_DECOMPOSE_BAOWU = 35
)

func NewS2cUpdateGoodsMsg(id int32, count int32) pbutil.Buffer {
	msg := &S2CUpdateGoodsProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateGoodsProtoMsg(msg)
}

var s2c_update_goods = [...]byte{11, 1} // 1
func NewS2cUpdateGoodsProtoMsg(object *S2CUpdateGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_goods[:], "s2c_update_goods")

}

func NewS2cUpdateMultiGoodsMsg(id []int32, count []int32) pbutil.Buffer {
	msg := &S2CUpdateMultiGoodsProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateMultiGoodsProtoMsg(msg)
}

var s2c_update_multi_goods = [...]byte{11, 5} // 5
func NewS2cUpdateMultiGoodsProtoMsg(object *S2CUpdateMultiGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_goods[:], "s2c_update_multi_goods")

}

func NewS2cUseGoodsMsg(id int32, count int32) pbutil.Buffer {
	msg := &S2CUseGoodsProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUseGoodsProtoMsg(msg)
}

var s2c_use_goods = [...]byte{11, 3} // 3
func NewS2cUseGoodsProtoMsg(object *S2CUseGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_goods[:], "s2c_use_goods")

}

// 无效的物品id
var ERR_USE_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 11, 4, 1} // 4-1

// 无效的物品个数
var ERR_USE_GOODS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 11, 4, 2} // 4-2

// 物品个数不足
var ERR_USE_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{3, 11, 4, 3} // 4-3

// 发送的物品不能直接使用
var ERR_USE_GOODS_FAIL_CANT_USE = pbutil.StaticBuffer{3, 11, 4, 6} // 4-6

// 服务器忙，请稍后再试
var ERR_USE_GOODS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 11, 4, 4} // 4-4

// 资源已满，无法使用道具增加资源
var ERR_USE_GOODS_FAIL_RESOURCE_FULL = pbutil.StaticBuffer{3, 11, 4, 5} // 4-5

// 您未加入联盟不可使用
var ERR_USE_GOODS_FAIL_NEED_GUILD = pbutil.StaticBuffer{3, 11, 4, 7} // 4-7

// 千重楼今日还未扫荡过不可重置
var ERR_USE_GOODS_FAIL_FLOOR_LESS_THAN_MAX = pbutil.StaticBuffer{3, 11, 4, 8} // 4-8

// 过关斩将今日开启次数尚未用完
var ERR_USE_GOODS_FAIL_OPEN_TIMES_NOT_ZERO = pbutil.StaticBuffer{3, 11, 4, 9} // 4-9

func NewS2cUseCdrGoodsMsg(id int32, count int32, cdr_type int32, index int32, cooldown int32) pbutil.Buffer {
	msg := &S2CUseCdrGoodsProto{
		Id:       id,
		Count:    count,
		CdrType:  cdr_type,
		Index:    index,
		Cooldown: cooldown,
	}
	return NewS2cUseCdrGoodsProtoMsg(msg)
}

var s2c_use_cdr_goods = [...]byte{11, 7} // 7
func NewS2cUseCdrGoodsProtoMsg(object *S2CUseCdrGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_cdr_goods[:], "s2c_use_cdr_goods")

}

// 无效的物品id
var ERR_USE_CDR_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 11, 8, 1} // 8-1

// 无效的物品个数
var ERR_USE_CDR_GOODS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 11, 8, 2} // 8-2

// 物品个数不足
var ERR_USE_CDR_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{3, 11, 8, 3} // 8-3

// 发送的物品id不是减cd道具
var ERR_USE_CDR_GOODS_FAIL_CANT_USE = pbutil.StaticBuffer{3, 11, 8, 4} // 8-4

// 物品不是这个类型的减cd物品
var ERR_USE_CDR_GOODS_FAIL_INVALID_CDR_TYPE = pbutil.StaticBuffer{3, 11, 8, 6} // 8-6

// 无效的索引
var ERR_USE_CDR_GOODS_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 11, 8, 8} // 8-8

// 当前不处于cd状态
var ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD = pbutil.StaticBuffer{3, 11, 8, 7} // 8-7

// 服务器忙，请稍后再试
var ERR_USE_CDR_GOODS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 11, 8, 5} // 8-5

var GOODS_COMBINE_S2C = pbutil.StaticBuffer{2, 11, 10} // 10

// 物品不够
var ERR_GOODS_COMBINE_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{3, 11, 11, 1} // 11-1

// 合成没找到
var ERR_GOODS_COMBINE_FAIL_COMBINE_NOT_FOUND = pbutil.StaticBuffer{3, 11, 11, 2} // 11-2

// 合成数量非法
var ERR_GOODS_COMBINE_FAIL_COMBINE_COUNT_INVALID = pbutil.StaticBuffer{3, 11, 11, 3} // 11-3

func NewS2cGoodsPartsCombineMsg(id int32, count int32, select_index int32, prize []byte) pbutil.Buffer {
	msg := &S2CGoodsPartsCombineProto{
		Id:          id,
		Count:       count,
		SelectIndex: select_index,
		Prize:       prize,
	}
	return NewS2cGoodsPartsCombineProtoMsg(msg)
}

func NewS2cGoodsPartsCombineMarshalMsg(id int32, count int32, select_index int32, prize marshaler) pbutil.Buffer {
	msg := &S2CGoodsPartsCombineProto{
		Id:          id,
		Count:       count,
		SelectIndex: select_index,
		Prize:       safeMarshal(prize),
	}
	return NewS2cGoodsPartsCombineProtoMsg(msg)
}

var s2c_goods_parts_combine = [...]byte{11, 19} // 19
func NewS2cGoodsPartsCombineProtoMsg(object *S2CGoodsPartsCombineProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_goods_parts_combine[:], "s2c_goods_parts_combine")

}

// 物品碎片id不存在
var ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_FOUND = pbutil.StaticBuffer{3, 11, 20, 1} // 20-1

// 物品碎片不够
var ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{3, 11, 20, 2} // 20-2

// 无效的合成索引
var ERR_GOODS_PARTS_COMBINE_FAIL_INVALID_SELECT_INDEX = pbutil.StaticBuffer{3, 11, 20, 3} // 20-3

func NewS2cGoodsExpiredMsg(id []int32) pbutil.Buffer {
	msg := &S2CGoodsExpiredProto{
		Id: id,
	}
	return NewS2cGoodsExpiredProtoMsg(msg)
}

var s2c_goods_expired = [...]byte{11, 16} // 16
func NewS2cGoodsExpiredProtoMsg(object *S2CGoodsExpiredProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_goods_expired[:], "s2c_goods_expired")

}

func NewS2cGoodsExpireTimeRemoveMsg(id []int32) pbutil.Buffer {
	msg := &S2CGoodsExpireTimeRemoveProto{
		Id: id,
	}
	return NewS2cGoodsExpireTimeRemoveProtoMsg(msg)
}

var s2c_goods_expire_time_remove = [...]byte{11, 17} // 17
func NewS2cGoodsExpireTimeRemoveProtoMsg(object *S2CGoodsExpireTimeRemoveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_goods_expire_time_remove[:], "s2c_goods_expire_time_remove")

}

func NewS2cUpdateBaowuMsg(id int32, count int32) pbutil.Buffer {
	msg := &S2CUpdateBaowuProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateBaowuProtoMsg(msg)
}

var s2c_update_baowu = [...]byte{11, 21} // 21
func NewS2cUpdateBaowuProtoMsg(object *S2CUpdateBaowuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_baowu[:], "s2c_update_baowu")

}

func NewS2cUpdateMultiBaowuMsg(id []int32, count []int32) pbutil.Buffer {
	msg := &S2CUpdateMultiBaowuProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateMultiBaowuProtoMsg(msg)
}

var s2c_update_multi_baowu = [...]byte{11, 22} // 22
func NewS2cUpdateMultiBaowuProtoMsg(object *S2CUpdateMultiBaowuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_baowu[:], "s2c_update_multi_baowu")

}

func NewS2cUnlockBaowuMsg(id int32, end_time int32) pbutil.Buffer {
	msg := &S2CUnlockBaowuProto{
		Id:      id,
		EndTime: end_time,
	}
	return NewS2cUnlockBaowuProtoMsg(msg)
}

var s2c_unlock_baowu = [...]byte{11, 24} // 24
func NewS2cUnlockBaowuProtoMsg(object *S2CUnlockBaowuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_baowu[:], "s2c_unlock_baowu")

}

// 你没有这个宝物，无法解锁
var ERR_UNLOCK_BAOWU_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 11, 25, 1} // 25-1

// 正在解锁其他的宝物，请稍后再试
var ERR_UNLOCK_BAOWU_FAIL_UNLOCKING = pbutil.StaticBuffer{3, 11, 25, 2} // 25-2

func NewS2cCollectBaowuMsg(prize *shared_proto.PrizeProto) pbutil.Buffer {
	msg := &S2CCollectBaowuProto{
		Prize: prize,
	}
	return NewS2cCollectBaowuProtoMsg(msg)
}

var s2c_collect_baowu = [...]byte{11, 27} // 27
func NewS2cCollectBaowuProtoMsg(object *S2CCollectBaowuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_baowu[:], "s2c_collect_baowu")

}

// 没有解锁中的宝箱，不能领取
var ERR_COLLECT_BAOWU_FAIL_NO_UNLOCKING = pbutil.StaticBuffer{3, 11, 28, 2} // 28-2

// 宝物开启中，不能领取
var ERR_COLLECT_BAOWU_FAIL_UNLOCKING = pbutil.StaticBuffer{3, 11, 28, 1} // 28-1

// 提前开启，消耗不足
var ERR_COLLECT_BAOWU_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 11, 28, 3} // 28-3

// 提前开启次数已达上限
var ERR_COLLECT_BAOWU_FAIL_MIAO_TIMES_LIMIT = pbutil.StaticBuffer{3, 11, 28, 4} // 28-4

func NewS2cAddBaowuLogMsg(data []byte) pbutil.Buffer {
	msg := &S2CAddBaowuLogProto{
		Data: data,
	}
	return NewS2cAddBaowuLogProtoMsg(msg)
}

func NewS2cAddBaowuLogMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CAddBaowuLogProto{
		Data: safeMarshal(data),
	}
	return NewS2cAddBaowuLogProtoMsg(msg)
}

var s2c_add_baowu_log = [...]byte{11, 29} // 29
func NewS2cAddBaowuLogProtoMsg(object *S2CAddBaowuLogProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_baowu_log[:], "s2c_add_baowu_log")

}

func NewS2cListBaowuLogMsg(datas [][]byte) pbutil.Buffer {
	msg := &S2CListBaowuLogProto{
		Datas: datas,
	}
	return NewS2cListBaowuLogProtoMsg(msg)
}

func NewS2cListBaowuLogMarshalMsg(datas [][]byte) pbutil.Buffer {
	msg := &S2CListBaowuLogProto{
		Datas: datas,
	}
	return NewS2cListBaowuLogProtoMsg(msg)
}

var s2c_list_baowu_log = [...]byte{11, 31} // 31
func NewS2cListBaowuLogProtoMsg(object *S2CListBaowuLogProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_baowu_log[:], "s2c_list_baowu_log")

}

func NewS2cDecomposeBaowuMsg(baowu_id int32, new_count int32) pbutil.Buffer {
	msg := &S2CDecomposeBaowuProto{
		BaowuId:  baowu_id,
		NewCount: new_count,
	}
	return NewS2cDecomposeBaowuProtoMsg(msg)
}

var s2c_decompose_baowu = [...]byte{11, 36} // 36
func NewS2cDecomposeBaowuProtoMsg(object *S2CDecomposeBaowuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_decompose_baowu[:], "s2c_decompose_baowu")

}

// 宝物id不存在
var ERR_DECOMPOSE_BAOWU_FAIL_INVALID_BAOWU = pbutil.StaticBuffer{3, 11, 37, 1} // 37-1

// 宝物数量错误
var ERR_DECOMPOSE_BAOWU_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 11, 37, 2} // 37-2

// 宝物正在解锁
var ERR_DECOMPOSE_BAOWU_FAIL_BAOWU_IS_UNLOCKING = pbutil.StaticBuffer{3, 11, 37, 3} // 37-3
