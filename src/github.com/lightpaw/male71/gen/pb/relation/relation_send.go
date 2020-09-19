package relation

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
	MODULE_ID = 35

	C2S_ADD_RELATION = 1

	C2S_REMOVE_ENEMY = 10

	C2S_REMOVE_RELATION = 4

	C2S_LIST_RELATION = 7

	C2S_NEW_LIST_RELATION = 28

	C2S_RECOMMEND_HERO_LIST = 16

	C2S_SEARCH_HEROS = 22

	C2S_SEARCH_HERO_BY_ID = 25

	C2S_SET_IMPORTANT_FRIEND = 33

	C2S_CANCEL_IMPORTANT_FRIEND = 36
)

func NewS2cAddRelationMsg(friend bool, id []byte, proto []byte, create_time int32) pbutil.Buffer {
	msg := &S2CAddRelationProto{
		Friend:     friend,
		Id:         id,
		Proto:      proto,
		CreateTime: create_time,
	}
	return NewS2cAddRelationProtoMsg(msg)
}

func NewS2cAddRelationMarshalMsg(friend bool, id []byte, proto marshaler, create_time int32) pbutil.Buffer {
	msg := &S2CAddRelationProto{
		Friend:     friend,
		Id:         id,
		Proto:      safeMarshal(proto),
		CreateTime: create_time,
	}
	return NewS2cAddRelationProtoMsg(msg)
}

var s2c_add_relation = [...]byte{35, 2} // 2
func NewS2cAddRelationProtoMsg(object *S2CAddRelationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_relation[:], "s2c_add_relation")

}

// 无效的玩家id
var ERR_ADD_RELATION_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 35, 3, 1} // 3-1

// 不能是自己的id
var ERR_ADD_RELATION_FAIL_SELF_ID = pbutil.StaticBuffer{3, 35, 3, 2} // 3-2

// 玩家已经是你的好友
var ERR_ADD_RELATION_FAIL_FRIEND = pbutil.StaticBuffer{3, 35, 3, 3} // 3-3

// 玩家已经是你的黑名单
var ERR_ADD_RELATION_FAIL_BLACK = pbutil.StaticBuffer{3, 35, 3, 4} // 3-4

// 好友个数已达上限
var ERR_ADD_RELATION_FAIL_LIMIT = pbutil.StaticBuffer{3, 35, 3, 5} // 3-5

func NewS2cAddEnemyMsg(id []byte) pbutil.Buffer {
	msg := &S2CAddEnemyProto{
		Id: id,
	}
	return NewS2cAddEnemyProtoMsg(msg)
}

var s2c_add_enemy = [...]byte{35, 9} // 9
func NewS2cAddEnemyProtoMsg(object *S2CAddEnemyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_enemy[:], "s2c_add_enemy")

}

func NewS2cRemoveEnemyMsg(id []byte) pbutil.Buffer {
	msg := &S2CRemoveEnemyProto{
		Id: id,
	}
	return NewS2cRemoveEnemyProtoMsg(msg)
}

var s2c_remove_enemy = [...]byte{35, 11} // 11
func NewS2cRemoveEnemyProtoMsg(object *S2CRemoveEnemyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_enemy[:], "s2c_remove_enemy")

}

// 无效的玩家id
var ERR_REMOVE_ENEMY_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 35, 12, 1} // 12-1

// 不能是自己的id
var ERR_REMOVE_ENEMY_FAIL_SELF_ID = pbutil.StaticBuffer{3, 35, 12, 2} // 12-2

func NewS2cRemoveRelationMsg(id []byte, rt int32) pbutil.Buffer {
	msg := &S2CRemoveRelationProto{
		Id: id,
		Rt: rt,
	}
	return NewS2cRemoveRelationProtoMsg(msg)
}

var s2c_remove_relation = [...]byte{35, 5} // 5
func NewS2cRemoveRelationProtoMsg(object *S2CRemoveRelationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_relation[:], "s2c_remove_relation")

}

// 无效的玩家id
var ERR_REMOVE_RELATION_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 35, 6, 1} // 6-1

// 不能是自己的id
var ERR_REMOVE_RELATION_FAIL_SELF_ID = pbutil.StaticBuffer{3, 35, 6, 2} // 6-2

func NewS2cListRelationMsg(version int32, proto [][]byte) pbutil.Buffer {
	msg := &S2CListRelationProto{
		Version: version,
		Proto:   proto,
	}
	return NewS2cListRelationProtoMsg(msg)
}

var s2c_list_relation = [...]byte{35, 8} // 8
func NewS2cListRelationProtoMsg(object *S2CListRelationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_relation[:], "s2c_list_relation")

}

func NewS2cNewListRelationMsg(version int32, proto [][]byte) pbutil.Buffer {
	msg := &S2CNewListRelationProto{
		Version: version,
		Proto:   proto,
	}
	return NewS2cNewListRelationProtoMsg(msg)
}

var s2c_new_list_relation = [...]byte{35, 29} // 29
func NewS2cNewListRelationProtoMsg(object *S2CNewListRelationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_new_list_relation[:], "s2c_new_list_relation")

}

func NewS2cRecommendHeroListMsg(page int32, heros []*shared_proto.HeroBasicSnapshotProto) pbutil.Buffer {
	msg := &S2CRecommendHeroListProto{
		Page:  page,
		Heros: heros,
	}
	return NewS2cRecommendHeroListProtoMsg(msg)
}

var s2c_recommend_hero_list = [...]byte{35, 17} // 17
func NewS2cRecommendHeroListProtoMsg(object *S2CRecommendHeroListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recommend_hero_list[:], "s2c_recommend_hero_list")

}

// 在刷新 CD 中
var ERR_RECOMMEND_HERO_LIST_FAIL_IN_CD = pbutil.StaticBuffer{3, 35, 18, 1} // 18-1

// 服务器错误
var ERR_RECOMMEND_HERO_LIST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 35, 18, 2} // 18-2

func NewS2cSearchHerosMsg(heros []*shared_proto.HeroBasicSnapshotProto) pbutil.Buffer {
	msg := &S2CSearchHerosProto{
		Heros: heros,
	}
	return NewS2cSearchHerosProtoMsg(msg)
}

var s2c_search_heros = [...]byte{35, 23} // 23
func NewS2cSearchHerosProtoMsg(object *S2CSearchHerosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_heros[:], "s2c_search_heros")

}

// 参数错误
var ERR_SEARCH_HEROS_FAIL_INVALID_ARG = pbutil.StaticBuffer{3, 35, 24, 1} // 24-1

// 服务器忙，请稍后再试
var ERR_SEARCH_HEROS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 35, 24, 2} // 24-2

// 在刷新 CD 中
var ERR_SEARCH_HEROS_FAIL_IN_CD = pbutil.StaticBuffer{3, 35, 24, 3} // 24-3

func NewS2cSearchHeroByIdMsg(hero *shared_proto.HeroBasicSnapshotProto) pbutil.Buffer {
	msg := &S2CSearchHeroByIdProto{
		Hero: hero,
	}
	return NewS2cSearchHeroByIdProtoMsg(msg)
}

var s2c_search_hero_by_id = [...]byte{35, 26} // 26
func NewS2cSearchHeroByIdProtoMsg(object *S2CSearchHeroByIdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_hero_by_id[:], "s2c_search_hero_by_id")

}

// 没有这个人
var ERR_SEARCH_HERO_BY_ID_FAIL_NO_HERO = pbutil.StaticBuffer{3, 35, 27, 4} // 27-4

// 参数错误
var ERR_SEARCH_HERO_BY_ID_FAIL_INVALID_ARG = pbutil.StaticBuffer{3, 35, 27, 1} // 27-1

// 服务器忙，请稍后再试
var ERR_SEARCH_HERO_BY_ID_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 35, 27, 2} // 27-2

// 在刷新 CD 中
var ERR_SEARCH_HERO_BY_ID_FAIL_IN_CD = pbutil.StaticBuffer{3, 35, 27, 3} // 27-3

func NewS2cSetImportantFriendMsg(id []byte, set_time int32) pbutil.Buffer {
	msg := &S2CSetImportantFriendProto{
		Id:      id,
		SetTime: set_time,
	}
	return NewS2cSetImportantFriendProtoMsg(msg)
}

var s2c_set_important_friend = [...]byte{35, 34} // 34
func NewS2cSetImportantFriendProtoMsg(object *S2CSetImportantFriendProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_important_friend[:], "s2c_set_important_friend")

}

// 不是自己的好友
var ERR_SET_IMPORTANT_FRIEND_FAIL_NOT_FRIEND = pbutil.StaticBuffer{3, 35, 35, 1} // 35-1

// 已经是星标好友
var ERR_SET_IMPORTANT_FRIEND_FAIL_HAS_SET = pbutil.StaticBuffer{3, 35, 35, 2} // 35-2

// 无效的玩家id
var ERR_SET_IMPORTANT_FRIEND_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 35, 35, 3} // 35-3

// 不能是自己的id
var ERR_SET_IMPORTANT_FRIEND_FAIL_SELF_ID = pbutil.StaticBuffer{3, 35, 35, 4} // 35-4

func NewS2cCancelImportantFriendMsg(id []byte) pbutil.Buffer {
	msg := &S2CCancelImportantFriendProto{
		Id: id,
	}
	return NewS2cCancelImportantFriendProtoMsg(msg)
}

var s2c_cancel_important_friend = [...]byte{35, 37} // 37
func NewS2cCancelImportantFriendProtoMsg(object *S2CCancelImportantFriendProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_cancel_important_friend[:], "s2c_cancel_important_friend")

}

// 不是自己的好友
var ERR_CANCEL_IMPORTANT_FRIEND_FAIL_NOT_FRIEND = pbutil.StaticBuffer{3, 35, 38, 1} // 38-1

// 不是星标好友
var ERR_CANCEL_IMPORTANT_FRIEND_FAIL_HAS_CANCEL = pbutil.StaticBuffer{3, 35, 38, 2} // 38-2

// 无效的玩家id
var ERR_CANCEL_IMPORTANT_FRIEND_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 35, 38, 3} // 38-3

// 不能是自己的id
var ERR_CANCEL_IMPORTANT_FRIEND_FAIL_SELF_ID = pbutil.StaticBuffer{3, 35, 38, 4} // 38-4
