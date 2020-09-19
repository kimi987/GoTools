package gem

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
	MODULE_ID = 19

	C2S_USE_GEM = 3

	C2S_INLAY_GEM = 21

	C2S_COMBINE_GEM = 6

	C2S_ONE_KEY_USE_GEM = 9

	C2S_ONE_KEY_COMBINE_GEM = 11

	C2S_REQUEST_ONE_KEY_COMBINE_COST = 15

	C2S_ONE_KEY_COMBINE_DEPOT_GEM = 18
)

func NewS2cUseGemMsg(captain_id int32, slot_idx int32, up_id int32) pbutil.Buffer {
	msg := &S2CUseGemProto{
		CaptainId: captain_id,
		SlotIdx:   slot_idx,
		UpId:      up_id,
	}
	return NewS2cUseGemProtoMsg(msg)
}

var s2c_use_gem = [...]byte{19, 4} // 4
func NewS2cUseGemProtoMsg(object *S2CUseGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_gem[:], "s2c_use_gem")

}

// 无效的武将id
var ERR_USE_GEM_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 19, 5, 1} // 5-1

// 无效的宝石id
var ERR_USE_GEM_FAIL_INVALID_GEM_ID = pbutil.StaticBuffer{3, 19, 5, 2} // 5-2

// 无效的宝石槽位
var ERR_USE_GEM_FAIL_INVALID_SLOT_IDX = pbutil.StaticBuffer{3, 19, 5, 3} // 5-3

// 在该部位已经有该类型的宝石了
var ERR_USE_GEM_FAIL_HAVE_SAME_TYPE_GEM_IN_THIS_TYPE = pbutil.StaticBuffer{3, 19, 5, 4} // 5-4

// 武将出征中，不能卸下宝石
var ERR_USE_GEM_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 19, 5, 5} // 5-5

// 成长值不够
var ERR_USE_GEM_FAIL_ABILITY_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 5, 6} // 5-6

// 宝石不够
var ERR_USE_GEM_FAIL_GEM_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 5, 9} // 5-9

// 上面镶嵌了相同的宝石
var ERR_USE_GEM_FAIL_HAS_SAME_GEM_IN_THE_SLOT = pbutil.StaticBuffer{3, 19, 5, 10} // 5-10

// 服务器忙，请稍后再试
var ERR_USE_GEM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 19, 5, 8} // 5-8

// 武将出征中，只能更换更高品质的宝石
var ERR_USE_GEM_FAIL_CAPTAIN_OUTSIDE_QUALITY_ERR = pbutil.StaticBuffer{3, 19, 5, 11} // 5-11

func NewS2cInlayGemMsg(captain_id int32, slot_idx int32, gem_id int32) pbutil.Buffer {
	msg := &S2CInlayGemProto{
		CaptainId: captain_id,
		SlotIdx:   slot_idx,
		GemId:     gem_id,
	}
	return NewS2cInlayGemProtoMsg(msg)
}

var s2c_inlay_gem = [...]byte{19, 22} // 22
func NewS2cInlayGemProtoMsg(object *S2CInlayGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_inlay_gem[:], "s2c_inlay_gem")

}

// 宝石不存在
var ERR_INLAY_GEM_FAIL_NO_GEM = pbutil.StaticBuffer{3, 19, 23, 1} // 23-1

// 武将不存在
var ERR_INLAY_GEM_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{3, 19, 23, 2} // 23-2

// 该槽位无法镶嵌
var ERR_INLAY_GEM_FAIL_WRONG_SLOT_IDX = pbutil.StaticBuffer{3, 19, 23, 3} // 23-3

// 无法镶嵌这类宝石
var ERR_INLAY_GEM_FAIL_WRONG_GEM = pbutil.StaticBuffer{3, 19, 23, 4} // 23-4

// 已镶嵌相同宝石
var ERR_INLAY_GEM_FAIL_HAS_SAME_GEM = pbutil.StaticBuffer{3, 19, 23, 5} // 23-5

func NewS2cCombineGemMsg(captain_id int32, slot_idx int32, gem_id int32) pbutil.Buffer {
	msg := &S2CCombineGemProto{
		CaptainId: captain_id,
		SlotIdx:   slot_idx,
		GemId:     gem_id,
	}
	return NewS2cCombineGemProtoMsg(msg)
}

var s2c_combine_gem = [...]byte{19, 7} // 7
func NewS2cCombineGemProtoMsg(object *S2CCombineGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_combine_gem[:], "s2c_combine_gem")

}

// 无效的武将id
var ERR_COMBINE_GEM_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 19, 8, 1} // 8-1

// 无效的宝石id
var ERR_COMBINE_GEM_FAIL_INVALID_GEM_ID = pbutil.StaticBuffer{3, 19, 8, 2} // 8-2

// 无效的宝石槽位
var ERR_COMBINE_GEM_FAIL_INVALID_SLOT_IDX = pbutil.StaticBuffer{3, 19, 8, 3} // 8-3

// 宝石满级了
var ERR_COMBINE_GEM_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 19, 8, 4} // 8-4

// 宝石槽位没有开启
var ERR_COMBINE_GEM_FAIL_SLOT_NOT_OPEN = pbutil.StaticBuffer{3, 19, 8, 7} // 8-7

// 武将出征中，不能够合成宝石
var ERR_COMBINE_GEM_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 19, 8, 9} // 8-9

// 宝石数量不够
var ERR_COMBINE_GEM_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 8, 11} // 8-11

func NewS2cOneKeyUseGemMsg(captain_id int32, down_all bool, gem_id []int32, equip_type int32) pbutil.Buffer {
	msg := &S2COneKeyUseGemProto{
		CaptainId: captain_id,
		DownAll:   down_all,
		GemId:     gem_id,
		EquipType: equip_type,
	}
	return NewS2cOneKeyUseGemProtoMsg(msg)
}

var s2c_one_key_use_gem = [...]byte{19, 10} // 10
func NewS2cOneKeyUseGemProtoMsg(object *S2COneKeyUseGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_one_key_use_gem[:], "s2c_one_key_use_gem")

}

// 无效的武将id
var ERR_ONE_KEY_USE_GEM_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 19, 14, 1} // 14-1

// 没有可以升级的宝石
var ERR_ONE_KEY_USE_GEM_FAIL_NO_CAN_UPGRADE_GEM = pbutil.StaticBuffer{3, 19, 14, 2} // 14-2

// 武将出征中，不能卸下宝石
var ERR_ONE_KEY_USE_GEM_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 19, 14, 4} // 14-4

// 武将没有宝石可以卸下
var ERR_ONE_KEY_USE_GEM_FAIL_NO_GEM_TO_DOWN_ALL = pbutil.StaticBuffer{3, 19, 14, 5} // 14-5

// 服务器忙，请稍后再试
var ERR_ONE_KEY_USE_GEM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 19, 14, 3} // 14-3

func NewS2cOneKeyCombineGemMsg(captain_id int32, slot_idx int32, gem_id int32) pbutil.Buffer {
	msg := &S2COneKeyCombineGemProto{
		CaptainId: captain_id,
		SlotIdx:   slot_idx,
		GemId:     gem_id,
	}
	return NewS2cOneKeyCombineGemProtoMsg(msg)
}

var s2c_one_key_combine_gem = [...]byte{19, 12} // 12
func NewS2cOneKeyCombineGemProtoMsg(object *S2COneKeyCombineGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_one_key_combine_gem[:], "s2c_one_key_combine_gem")

}

// 无效的武将id
var ERR_ONE_KEY_COMBINE_GEM_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 19, 13, 1} // 13-1

// 无效的宝石槽位
var ERR_ONE_KEY_COMBINE_GEM_FAIL_INVALID_SLOT_IDX = pbutil.StaticBuffer{3, 19, 13, 2} // 13-2

// 宝石满级了
var ERR_ONE_KEY_COMBINE_GEM_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 19, 13, 3} // 13-3

// 宝石槽位没有开启
var ERR_ONE_KEY_COMBINE_GEM_FAIL_SLOT_NOT_OPEN = pbutil.StaticBuffer{3, 19, 13, 4} // 13-4

// 武将出征中，不能够合成宝石
var ERR_ONE_KEY_COMBINE_GEM_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 19, 13, 5} // 13-5

// 宝石数量不够
var ERR_ONE_KEY_COMBINE_GEM_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 13, 6} // 13-6

// 宝石未开放购买
var ERR_ONE_KEY_COMBINE_GEM_FAIL_CANT_BUY = pbutil.StaticBuffer{3, 19, 13, 8} // 13-8

// 购买宝石，消耗不足
var ERR_ONE_KEY_COMBINE_GEM_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 13, 9} // 13-9

// 服务器忙，请稍后再试
var ERR_ONE_KEY_COMBINE_GEM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 19, 13, 7} // 13-7

func NewS2cRequestOneKeyCombineCostMsg(captain_id int32, slot_idx int32, can_combine bool, gem_id []int32, gem_count []int32, buy_count int32, buy_yuanbao int32) pbutil.Buffer {
	msg := &S2CRequestOneKeyCombineCostProto{
		CaptainId:  captain_id,
		SlotIdx:    slot_idx,
		CanCombine: can_combine,
		GemId:      gem_id,
		GemCount:   gem_count,
		BuyCount:   buy_count,
		BuyYuanbao: buy_yuanbao,
	}
	return NewS2cRequestOneKeyCombineCostProtoMsg(msg)
}

var s2c_request_one_key_combine_cost = [...]byte{19, 16} // 16
func NewS2cRequestOneKeyCombineCostProtoMsg(object *S2CRequestOneKeyCombineCostProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_one_key_combine_cost[:], "s2c_request_one_key_combine_cost")

}

// 无效的武将id
var ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 19, 17, 1} // 17-1

// 无效的宝石槽位
var ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_SLOT_IDX = pbutil.StaticBuffer{3, 19, 17, 2} // 17-2

// 宝石满级了
var ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 19, 17, 3} // 17-3

// 宝石槽位没有开启
var ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_SLOT_NOT_OPEN = pbutil.StaticBuffer{3, 19, 17, 4} // 17-4

// 武将出征中，不能够合成宝石
var ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 19, 17, 5} // 17-5

func NewS2cOneKeyCombineDepotGemMsg(new_gem_id int32, new_gem_count int32) pbutil.Buffer {
	msg := &S2COneKeyCombineDepotGemProto{
		NewGemId:    new_gem_id,
		NewGemCount: new_gem_count,
	}
	return NewS2cOneKeyCombineDepotGemProtoMsg(msg)
}

var s2c_one_key_combine_depot_gem = [...]byte{19, 19} // 19
func NewS2cOneKeyCombineDepotGemProtoMsg(object *S2COneKeyCombineDepotGemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_one_key_combine_depot_gem[:], "s2c_one_key_combine_depot_gem")

}

// 无效的宝石id
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_INVALID_GEM_ID = pbutil.StaticBuffer{3, 19, 20, 1} // 20-1

// 宝石满级了
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 19, 20, 2} // 20-2

// 宝石数量不够
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 20, 3} // 20-3

// 宝石未开放购买
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_CANT_BUY = pbutil.StaticBuffer{3, 19, 20, 5} // 20-5

// 购买宝石，消耗不足
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 19, 20, 6} // 20-6

// 服务器忙，请稍后再试
var ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 19, 20, 4} // 20-4
