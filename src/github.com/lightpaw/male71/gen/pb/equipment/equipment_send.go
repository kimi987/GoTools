package equipment

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
	MODULE_ID = 12

	C2S_VIEW_CHAT_EQUIP = 40

	C2S_WEAR_EQUIPMENT = 1

	C2S_UPGRADE_EQUIPMENT = 4

	C2S_UPGRADE_EQUIPMENT_ALL = 19

	C2S_REFINED_EQUIPMENT = 7

	C2S_SMELT_EQUIPMENT = 10

	C2S_REBUILD_EQUIPMENT = 13

	C2S_ONE_KEY_TAKE_OFF = 43
)

func NewS2cViewChatEquipMsg(data []byte) pbutil.Buffer {
	msg := &S2CViewChatEquipProto{
		Data: data,
	}
	return NewS2cViewChatEquipProtoMsg(msg)
}

func NewS2cViewChatEquipMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CViewChatEquipProto{
		Data: safeMarshal(data),
	}
	return NewS2cViewChatEquipProtoMsg(msg)
}

var s2c_view_chat_equip = [...]byte{12, 41} // 41
func NewS2cViewChatEquipProtoMsg(object *S2CViewChatEquipProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_chat_equip[:], "s2c_view_chat_equip")

}

// // 无效的链接
var ERR_VIEW_CHAT_EQUIP_FAIL_INVALID = pbutil.StaticBuffer{3, 12, 42, 1} // 42-1

func NewS2cAddEquipmentMsg(data []byte) pbutil.Buffer {
	msg := &S2CAddEquipmentProto{
		Data: data,
	}
	return NewS2cAddEquipmentProtoMsg(msg)
}

func NewS2cAddEquipmentMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CAddEquipmentProto{
		Data: safeMarshal(data),
	}
	return NewS2cAddEquipmentProtoMsg(msg)
}

var s2c_add_equipment = [...]byte{12, 18} // 18
func NewS2cAddEquipmentProtoMsg(object *S2CAddEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_equipment[:], "s2c_add_equipment")

}

func NewS2cAddEquipmentWithExpireTimeMsg(data []byte, expire_time int32) pbutil.Buffer {
	msg := &S2CAddEquipmentWithExpireTimeProto{
		Data:       data,
		ExpireTime: expire_time,
	}
	return NewS2cAddEquipmentWithExpireTimeProtoMsg(msg)
}

func NewS2cAddEquipmentWithExpireTimeMarshalMsg(data marshaler, expire_time int32) pbutil.Buffer {
	msg := &S2CAddEquipmentWithExpireTimeProto{
		Data:       safeMarshal(data),
		ExpireTime: expire_time,
	}
	return NewS2cAddEquipmentWithExpireTimeProtoMsg(msg)
}

var s2c_add_equipment_with_expire_time = [...]byte{12, 34} // 34
func NewS2cAddEquipmentWithExpireTimeProtoMsg(object *S2CAddEquipmentWithExpireTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_equipment_with_expire_time[:], "s2c_add_equipment_with_expire_time")

}

func NewS2cWearEquipmentMsg(captain_id int32, up_id int32, down_id int32, taoz int32) pbutil.Buffer {
	msg := &S2CWearEquipmentProto{
		CaptainId: captain_id,
		UpId:      up_id,
		DownId:    down_id,
		Taoz:      taoz,
	}
	return NewS2cWearEquipmentProtoMsg(msg)
}

var s2c_wear_equipment = [...]byte{12, 2} // 2
func NewS2cWearEquipmentProtoMsg(object *S2CWearEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_wear_equipment[:], "s2c_wear_equipment")

}

// 无效的武将id
var ERR_WEAR_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 12, 3, 1} // 3-1

// 无效的装备id
var ERR_WEAR_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID = pbutil.StaticBuffer{3, 12, 3, 2} // 3-2

// 武将出征中，不能卸下装备
var ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 12, 3, 4} // 3-4

// 装备背包已满
var ERR_WEAR_EQUIPMENT_FAIL_DEPOT_EQUIP_FULL = pbutil.StaticBuffer{3, 12, 3, 5} // 3-5

// 服务器忙，请稍后再试
var ERR_WEAR_EQUIPMENT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 12, 3, 3} // 3-3

// 武将出征中，替换装备必须继承
var ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE_MUST_INHERIT = pbutil.StaticBuffer{3, 12, 3, 6} // 3-6

// 武将出征中，必须替换更高品质的装备
var ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE_QUALITY_ERR = pbutil.StaticBuffer{3, 12, 3, 7} // 3-7

func NewS2cUpgradeEquipmentMsg(captain_id int32, equipment_id int32, level int32) pbutil.Buffer {
	msg := &S2CUpgradeEquipmentProto{
		CaptainId:   captain_id,
		EquipmentId: equipment_id,
		Level:       level,
	}
	return NewS2cUpgradeEquipmentProtoMsg(msg)
}

var s2c_upgrade_equipment = [...]byte{12, 5} // 5
func NewS2cUpgradeEquipmentProtoMsg(object *S2CUpgradeEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_equipment[:], "s2c_upgrade_equipment")

}

// 无效的武将id
var ERR_UPGRADE_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 12, 6, 1} // 6-1

// 无效的装备id
var ERR_UPGRADE_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID = pbutil.StaticBuffer{3, 12, 6, 2} // 6-2

// 已达当前最大等级
var ERR_UPGRADE_EQUIPMENT_FAIL_LEVEL_LIMIT = pbutil.StaticBuffer{3, 12, 6, 3} // 6-3

// 消耗不足
var ERR_UPGRADE_EQUIPMENT_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 12, 6, 4} // 6-4

// 武将出征中，不能操作
var ERR_UPGRADE_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 12, 6, 6} // 6-6

// 服务器忙，请稍后再试
var ERR_UPGRADE_EQUIPMENT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 12, 6, 5} // 6-5

func NewS2cUpgradeEquipmentAllMsg(captain_id int32, level []int32) pbutil.Buffer {
	msg := &S2CUpgradeEquipmentAllProto{
		CaptainId: captain_id,
		Level:     level,
	}
	return NewS2cUpgradeEquipmentAllProtoMsg(msg)
}

var s2c_upgrade_equipment_all = [...]byte{12, 20} // 20
func NewS2cUpgradeEquipmentAllProtoMsg(object *S2CUpgradeEquipmentAllProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_equipment_all[:], "s2c_upgrade_equipment_all")

}

// 无效的武将id
var ERR_UPGRADE_EQUIPMENT_ALL_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 12, 21, 1} // 21-1

// 没有可以升级的装备
var ERR_UPGRADE_EQUIPMENT_ALL_FAIL_CANT_UPGRADE = pbutil.StaticBuffer{3, 12, 21, 2} // 21-2

// 消耗不足
var ERR_UPGRADE_EQUIPMENT_ALL_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 12, 21, 4} // 21-4

// 武将出征中，不能操作
var ERR_UPGRADE_EQUIPMENT_ALL_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 12, 21, 5} // 21-5

// 服务器忙，请稍后再试
var ERR_UPGRADE_EQUIPMENT_ALL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 12, 21, 6} // 21-6

func NewS2cRefinedEquipmentMsg(captain_id int32, equipment_id int32, level int32, taoz int32) pbutil.Buffer {
	msg := &S2CRefinedEquipmentProto{
		CaptainId:   captain_id,
		EquipmentId: equipment_id,
		Level:       level,
		Taoz:        taoz,
	}
	return NewS2cRefinedEquipmentProtoMsg(msg)
}

var s2c_refined_equipment = [...]byte{12, 8} // 8
func NewS2cRefinedEquipmentProtoMsg(object *S2CRefinedEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_refined_equipment[:], "s2c_refined_equipment")

}

// 无效的武将id
var ERR_REFINED_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 12, 9, 1} // 9-1

// 无效的装备id
var ERR_REFINED_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID = pbutil.StaticBuffer{3, 12, 9, 2} // 9-2

// 已达当前最大等级
var ERR_REFINED_EQUIPMENT_FAIL_LEVEL_LIMIT = pbutil.StaticBuffer{3, 12, 9, 3} // 9-3

// 君主等级不足
var ERR_REFINED_EQUIPMENT_FAIL_HERO_LEVEL_LIMIT = pbutil.StaticBuffer{3, 12, 9, 7} // 9-7

// 消耗不足
var ERR_REFINED_EQUIPMENT_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 12, 9, 4} // 9-4

// 武将出征中，不能操作
var ERR_REFINED_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 12, 9, 6} // 9-6

// 服务器忙，请稍后再试
var ERR_REFINED_EQUIPMENT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 12, 9, 5} // 9-5

func NewS2cUpdateEquipmentMsg(data []byte) pbutil.Buffer {
	msg := &S2CUpdateEquipmentProto{
		Data: data,
	}
	return NewS2cUpdateEquipmentProtoMsg(msg)
}

func NewS2cUpdateEquipmentMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CUpdateEquipmentProto{
		Data: safeMarshal(data),
	}
	return NewS2cUpdateEquipmentProtoMsg(msg)
}

var s2c_update_equipment = [...]byte{12, 25} // 25
func NewS2cUpdateEquipmentProtoMsg(object *S2CUpdateEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_equipment[:], "s2c_update_equipment")

}

func NewS2cUpdateMultiEquipmentMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CUpdateMultiEquipmentProto{
		Data: data,
	}
	return NewS2cUpdateMultiEquipmentProtoMsg(msg)
}

func NewS2cUpdateMultiEquipmentMarshalMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CUpdateMultiEquipmentProto{
		Data: data,
	}
	return NewS2cUpdateMultiEquipmentProtoMsg(msg)
}

var s2c_update_multi_equipment = [...]byte{12, 26} // 26
func NewS2cUpdateMultiEquipmentProtoMsg(object *S2CUpdateMultiEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_equipment[:], "s2c_update_multi_equipment")

}

func NewS2cSmeltEquipmentMsg(equipment_id []int32) pbutil.Buffer {
	msg := &S2CSmeltEquipmentProto{
		EquipmentId: equipment_id,
	}
	return NewS2cSmeltEquipmentProtoMsg(msg)
}

var s2c_smelt_equipment = [...]byte{12, 11} // 11
func NewS2cSmeltEquipmentProtoMsg(object *S2CSmeltEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_smelt_equipment[:], "s2c_smelt_equipment")

}

// 无效的装备id
var ERR_SMELT_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID = pbutil.StaticBuffer{3, 12, 12, 1} // 12-1

func NewS2cRebuildEquipmentMsg(equipment_id []int32) pbutil.Buffer {
	msg := &S2CRebuildEquipmentProto{
		EquipmentId: equipment_id,
	}
	return NewS2cRebuildEquipmentProtoMsg(msg)
}

var s2c_rebuild_equipment = [...]byte{12, 14} // 14
func NewS2cRebuildEquipmentProtoMsg(object *S2CRebuildEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rebuild_equipment[:], "s2c_rebuild_equipment")

}

// 无效的装备id
var ERR_REBUILD_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID = pbutil.StaticBuffer{3, 12, 15, 1} // 15-1

// 装备没有升级和强化
var ERR_REBUILD_EQUIPMENT_FAIL_CANT_REBUILD = pbutil.StaticBuffer{3, 12, 15, 2} // 15-2

func NewS2cOpenEquipCombineMsg(id int32) pbutil.Buffer {
	msg := &S2COpenEquipCombineProto{
		Id: id,
	}
	return NewS2cOpenEquipCombineProtoMsg(msg)
}

var s2c_open_equip_combine = [...]byte{12, 33} // 33
func NewS2cOpenEquipCombineProtoMsg(object *S2COpenEquipCombineProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open_equip_combine[:], "s2c_open_equip_combine")

}

func NewS2cRebuildUpgradeEquipmentMsg(equipment_id int32, captain_id int32) pbutil.Buffer {
	msg := &S2CRebuildUpgradeEquipmentProto{
		EquipmentId: equipment_id,
		CaptainId:   captain_id,
	}
	return NewS2cRebuildUpgradeEquipmentProtoMsg(msg)
}

var s2c_rebuild_upgrade_equipment = [...]byte{12, 35} // 35
func NewS2cRebuildUpgradeEquipmentProtoMsg(object *S2CRebuildUpgradeEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rebuild_upgrade_equipment[:], "s2c_rebuild_upgrade_equipment")

}

func NewS2cRebuildRefineEquipmentMsg(equipment_id int32, captain_id int32) pbutil.Buffer {
	msg := &S2CRebuildRefineEquipmentProto{
		EquipmentId: equipment_id,
		CaptainId:   captain_id,
	}
	return NewS2cRebuildRefineEquipmentProtoMsg(msg)
}

var s2c_rebuild_refine_equipment = [...]byte{12, 36} // 36
func NewS2cRebuildRefineEquipmentProtoMsg(object *S2CRebuildRefineEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rebuild_refine_equipment[:], "s2c_rebuild_refine_equipment")

}

func NewS2cOneKeyTakeOffMsg(captain_id int32) pbutil.Buffer {
	msg := &S2COneKeyTakeOffProto{
		CaptainId: captain_id,
	}
	return NewS2cOneKeyTakeOffProtoMsg(msg)
}

var s2c_one_key_take_off = [...]byte{12, 44} // 44
func NewS2cOneKeyTakeOffProtoMsg(object *S2COneKeyTakeOffProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_one_key_take_off[:], "s2c_one_key_take_off")

}

// 没有该武将
var ERR_ONE_KEY_TAKE_OFF_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{3, 12, 45, 1} // 45-1

// 武将出征中
var ERR_ONE_KEY_TAKE_OFF_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 12, 45, 2} // 45-2

// 没有任何装备
var ERR_ONE_KEY_TAKE_OFF_FAIL_NO_EQUIPMENT = pbutil.StaticBuffer{3, 12, 45, 3} // 45-3

// 装备背包空间不足
var ERR_ONE_KEY_TAKE_OFF_FAIL_DEPOT_SPACE_NOT_ENOUGH = pbutil.StaticBuffer{3, 12, 45, 4} // 45-4
