package hebi

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
	MODULE_ID = 41

	C2S_ROOM_LIST = 1

	C2S_HERO_RECORD_LIST = 35

	C2S_CHANGE_CAPTAIN = 3

	C2S_CHECK_IN_ROOM = 28

	C2S_COPY_SELF = 31

	C2S_JOIN_ROOM = 9

	C2S_ROB_POS = 12

	C2S_LEAVE_ROOM = 18

	C2S_ROB = 21

	C2S_VIEW_SHOW_PRIZE = 37
)

func NewS2cRoomListMsg(v int32, list *shared_proto.HebiInfoProto) pbutil.Buffer {
	msg := &S2CRoomListProto{
		V:    v,
		List: list,
	}
	return NewS2cRoomListProtoMsg(msg)
}

var s2c_room_list = [...]byte{41, 2} // 2
func NewS2cRoomListProtoMsg(object *S2CRoomListProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_room_list[:], "s2c_room_list")

}

func NewS2cHeroRecordListMsg(records *shared_proto.HebiHeroRecordProto) pbutil.Buffer {
	msg := &S2CHeroRecordListProto{
		Records: records,
	}
	return NewS2cHeroRecordListProtoMsg(msg)
}

var s2c_hero_record_list = [...]byte{41, 36} // 36
func NewS2cHeroRecordListProtoMsg(object *S2CHeroRecordListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_record_list[:], "s2c_hero_record_list")

}

func NewS2cChangeCaptainMsg(captain_id int32) pbutil.Buffer {
	msg := &S2CChangeCaptainProto{
		CaptainId: captain_id,
	}
	return NewS2cChangeCaptainProtoMsg(msg)
}

var s2c_change_captain = [...]byte{41, 4} // 4
func NewS2cChangeCaptainProtoMsg(object *S2CChangeCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_captain[:], "s2c_change_captain")

}

// 武将不存在
var ERR_CHANGE_CAPTAIN_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 41, 5, 1} // 5-1

func NewS2cChangeRoomCaptainMsg(room_id int32, cpatain *shared_proto.HebiCaptainProto) pbutil.Buffer {
	msg := &S2CChangeRoomCaptainProto{
		RoomId:  room_id,
		Cpatain: cpatain,
	}
	return NewS2cChangeRoomCaptainProtoMsg(msg)
}

var s2c_change_room_captain = [...]byte{41, 34} // 34
func NewS2cChangeRoomCaptainProtoMsg(object *S2CChangeRoomCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_room_captain[:], "s2c_change_room_captain")

}

func NewS2cCheckInRoomMsg(room_id int32, goods_id int32) pbutil.Buffer {
	msg := &S2CCheckInRoomProto{
		RoomId:  room_id,
		GoodsId: goods_id,
	}
	return NewS2cCheckInRoomProtoMsg(msg)
}

var s2c_check_in_room = [...]byte{41, 29} // 29
func NewS2cCheckInRoomProtoMsg(object *S2CCheckInRoomProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_check_in_room[:], "s2c_check_in_room")

}

// 房间号错误
var ERR_CHECK_IN_ROOM_FAIL_INVALID_ROOM_ID = pbutil.StaticBuffer{3, 41, 30, 1} // 30-1

// 不是空房间
var ERR_CHECK_IN_ROOM_FAIL_ROOM_ID_NOT_EMPTY = pbutil.StaticBuffer{3, 41, 30, 2} // 30-2

// 已经在房间里
var ERR_CHECK_IN_ROOM_FAIL_ALREADY_IN_ROOM = pbutil.StaticBuffer{3, 41, 30, 6} // 30-6

// 物品错误或没在背包
var ERR_CHECK_IN_ROOM_FAIL_GOODS_INVALID = pbutil.StaticBuffer{3, 41, 30, 3} // 30-3

// 武将没有设置
var ERR_CHECK_IN_ROOM_FAIL_CAPTAIN_ERR = pbutil.StaticBuffer{3, 41, 30, 4} // 30-4

// 服务器错误
var ERR_CHECK_IN_ROOM_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 30, 5} // 30-5

func NewS2cCopySelfMsg(room_id int32) pbutil.Buffer {
	msg := &S2CCopySelfProto{
		RoomId: room_id,
	}
	return NewS2cCopySelfProtoMsg(msg)
}

var s2c_copy_self = [...]byte{41, 32} // 32
func NewS2cCopySelfProtoMsg(object *S2CCopySelfProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_copy_self[:], "s2c_copy_self")

}

// 房间号错误
var ERR_COPY_SELF_FAIL_INVALID_ROOM_ID = pbutil.StaticBuffer{3, 41, 33, 1} // 33-1

// 当前房间状态不能分身
var ERR_COPY_SELF_FAIL_ROOM_STATE_NOT_WAIT = pbutil.StaticBuffer{3, 41, 33, 4} // 33-4

// 物品不存在或不是分身物品
var ERR_COPY_SELF_FAIL_GOODS_INVALID = pbutil.StaticBuffer{3, 41, 33, 2} // 33-2

// 服务器错误
var ERR_COPY_SELF_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 33, 3} // 33-3

func NewS2cJoinRoomMsg(room_id int32, goods_id int32, prize_id int32) pbutil.Buffer {
	msg := &S2CJoinRoomProto{
		RoomId:  room_id,
		GoodsId: goods_id,
		PrizeId: prize_id,
	}
	return NewS2cJoinRoomProtoMsg(msg)
}

var s2c_join_room = [...]byte{41, 10} // 10
func NewS2cJoinRoomProtoMsg(object *S2CJoinRoomProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_join_room[:], "s2c_join_room")

}

// 房间号错误
var ERR_JOIN_ROOM_FAIL_INVALID_ROOM_ID = pbutil.StaticBuffer{3, 41, 11, 1} // 11-1

// 已经在房间里
var ERR_JOIN_ROOM_FAIL_ALREADY_IN_ROOM = pbutil.StaticBuffer{3, 41, 11, 6} // 11-6

// 玉璧错误，类型或左右错误
var ERR_JOIN_ROOM_FAIL_GOODS_INVALID = pbutil.StaticBuffer{3, 41, 11, 2} // 11-2

// 武将没有设置
var ERR_JOIN_ROOM_FAIL_CAPTAIN_ERR = pbutil.StaticBuffer{3, 41, 11, 5} // 11-5

// 联盟成员才能进和氏璧房间
var ERR_JOIN_ROOM_FAIL_HESHIBI_NOT_SAME_GUILD = pbutil.StaticBuffer{3, 41, 11, 3} // 11-3

// 服务器错误
var ERR_JOIN_ROOM_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 11, 4} // 11-4

func NewS2cSomeoneJoinedRoomMsg(room_id int32) pbutil.Buffer {
	msg := &S2CSomeoneJoinedRoomProto{
		RoomId: room_id,
	}
	return NewS2cSomeoneJoinedRoomProtoMsg(msg)
}

var s2c_someone_joined_room = [...]byte{41, 25} // 25
func NewS2cSomeoneJoinedRoomProtoMsg(object *S2CSomeoneJoinedRoomProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_someone_joined_room[:], "s2c_someone_joined_room")

}

func NewS2cRobPosMsg(room_id int32, succ bool, link string) pbutil.Buffer {
	msg := &S2CRobPosProto{
		RoomId: room_id,
		Succ:   succ,
		Link:   link,
	}
	return NewS2cRobPosProtoMsg(msg)
}

var s2c_rob_pos = [...]byte{41, 13} // 13
func NewS2cRobPosProtoMsg(object *S2CRobPosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rob_pos[:], "s2c_rob_pos")

}

// 已经在房间里了
var ERR_ROB_POS_FAIL_ALREADY_IN_ROOM = pbutil.StaticBuffer{3, 41, 14, 1} // 14-1

// 错误的房间号
var ERR_ROB_POS_FAIL_INVALID_ROOM_ID = pbutil.StaticBuffer{3, 41, 14, 2} // 14-2

// 错误的房间状态
var ERR_ROB_POS_FAIL_INVALID_ROOM_STATE = pbutil.StaticBuffer{3, 41, 14, 6} // 14-6

// 没有玉璧，不能抢
var ERR_ROB_POS_FAIL_NO_YUBI = pbutil.StaticBuffer{3, 41, 14, 3} // 14-3

// 同联盟
var ERR_ROB_POS_FAIL_SAME_GUILD = pbutil.StaticBuffer{3, 41, 14, 9} // 14-9

// 在抢夺 CD 中
var ERR_ROB_POS_FAIL_IN_ROB_CD = pbutil.StaticBuffer{3, 41, 14, 8} // 14-8

// 武将没有设置
var ERR_ROB_POS_FAIL_CAPTAIN_ERR = pbutil.StaticBuffer{3, 41, 14, 5} // 14-5

// 服务器错误
var ERR_ROB_POS_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 14, 4} // 14-4

func NewS2cSomeoneRobbedMyPosMsg(room_id int32) pbutil.Buffer {
	msg := &S2CSomeoneRobbedMyPosProto{
		RoomId: room_id,
	}
	return NewS2cSomeoneRobbedMyPosProtoMsg(msg)
}

var s2c_someone_robbed_my_pos = [...]byte{41, 26} // 26
func NewS2cSomeoneRobbedMyPosProtoMsg(object *S2CSomeoneRobbedMyPosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_someone_robbed_my_pos[:], "s2c_someone_robbed_my_pos")

}

func NewS2cLeaveRoomMsg(room_id int32) pbutil.Buffer {
	msg := &S2CLeaveRoomProto{
		RoomId: room_id,
	}
	return NewS2cLeaveRoomProtoMsg(msg)
}

var s2c_leave_room = [...]byte{41, 19} // 19
func NewS2cLeaveRoomProtoMsg(object *S2CLeaveRoomProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_leave_room[:], "s2c_leave_room")

}

// 没在房间里
var ERR_LEAVE_ROOM_FAIL_NOT_IN_ROOM = pbutil.StaticBuffer{3, 41, 20, 1} // 20-1

// 房间里有别人
var ERR_LEAVE_ROOM_FAIL_INVALID_ROOM_STATE = pbutil.StaticBuffer{3, 41, 20, 5} // 20-5

// 服务器错误
var ERR_LEAVE_ROOM_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 20, 2} // 20-2

func NewS2cRobMsg(room_id int32, first_succ bool, first_link string, second_succ bool, second_link string, prize *shared_proto.PrizeProto) pbutil.Buffer {
	msg := &S2CRobProto{
		RoomId:     room_id,
		FirstSucc:  first_succ,
		FirstLink:  first_link,
		SecondSucc: second_succ,
		SecondLink: second_link,
		Prize:      prize,
	}
	return NewS2cRobProtoMsg(msg)
}

var s2c_rob = [...]byte{41, 22} // 22
func NewS2cRobProtoMsg(object *S2CRobProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rob[:], "s2c_rob")

}

// 在 CD 中
var ERR_ROB_FAIL_IN_CD = pbutil.StaticBuffer{3, 41, 23, 1} // 23-1

// 没有次数
var ERR_ROB_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 41, 23, 7} // 23-7

// 错误的房间号
var ERR_ROB_FAIL_INVALID_ROOM_ID = pbutil.StaticBuffer{3, 41, 23, 2} // 23-2

// 已经在房间里
var ERR_ROB_FAIL_ALREADY_IN_ROOM = pbutil.StaticBuffer{3, 41, 23, 8} // 23-8

// 房间没在合璧中
var ERR_ROB_FAIL_ROOM_NOT_IN_HEBI = pbutil.StaticBuffer{3, 41, 23, 4} // 23-4

// 同联盟不能抢
var ERR_ROB_FAIL_IN_SAME_GUILD = pbutil.StaticBuffer{3, 41, 23, 3} // 23-3

// 武将没有设置
var ERR_ROB_FAIL_CAPTAIN_ERR = pbutil.StaticBuffer{3, 41, 23, 6} // 23-6

// 服务器错误
var ERR_ROB_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 41, 23, 5} // 23-5

func NewS2cSomeoneRobbedMyPrizeMsg(room_id int32) pbutil.Buffer {
	msg := &S2CSomeoneRobbedMyPrizeProto{
		RoomId: room_id,
	}
	return NewS2cSomeoneRobbedMyPrizeProtoMsg(msg)
}

var s2c_someone_robbed_my_prize = [...]byte{41, 27} // 27
func NewS2cSomeoneRobbedMyPrizeProtoMsg(object *S2CSomeoneRobbedMyPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_someone_robbed_my_prize[:], "s2c_someone_robbed_my_prize")

}

func NewS2cCompleteMsg(room_id int32) pbutil.Buffer {
	msg := &S2CCompleteProto{
		RoomId: room_id,
	}
	return NewS2cCompleteProtoMsg(msg)
}

var s2c_complete = [...]byte{41, 24} // 24
func NewS2cCompleteProtoMsg(object *S2CCompleteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_complete[:], "s2c_complete")

}

func NewS2cViewShowPrizeMsg(prize *shared_proto.PrizeProto) pbutil.Buffer {
	msg := &S2CViewShowPrizeProto{
		Prize: prize,
	}
	return NewS2cViewShowPrizeProtoMsg(msg)
}

var s2c_view_show_prize = [...]byte{41, 38} // 38
func NewS2cViewShowPrizeProtoMsg(object *S2CViewShowPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_show_prize[:], "s2c_view_show_prize")

}

// 等级错误
var ERR_VIEW_SHOW_PRIZE_FAIL_INVALID_HERO_LEVEL = pbutil.StaticBuffer{3, 41, 39, 1} // 39-1

// 物品错误
var ERR_VIEW_SHOW_PRIZE_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 41, 39, 2} // 39-2

// 没找到奖励
var ERR_VIEW_SHOW_PRIZE_FAIL_NO_PRIZE = pbutil.StaticBuffer{3, 41, 39, 3} // 39-3
