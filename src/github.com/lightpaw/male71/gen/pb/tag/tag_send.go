package tag

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
	MODULE_ID = 29

	C2S_ADD_OR_UPDATE_TAG = 1

	C2S_DELETE_TAG = 5
)

func NewS2cAddOrUpdateTagMsg(id []byte, record []byte, tag []byte) pbutil.Buffer {
	msg := &S2CAddOrUpdateTagProto{
		Id:     id,
		Record: record,
		Tag:    tag,
	}
	return NewS2cAddOrUpdateTagProtoMsg(msg)
}

var s2c_add_or_update_tag = [...]byte{29, 2} // 2
func NewS2cAddOrUpdateTagProtoMsg(object *S2CAddOrUpdateTagProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_or_update_tag[:], "s2c_add_or_update_tag")

}

// 你被加入黑名单了
var ERR_ADD_OR_UPDATE_TAG_FAIL_BLACK = pbutil.StaticBuffer{3, 29, 3, 1} // 3-1

// 目标没找到
var ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 29, 3, 2} // 3-2

// 标签太长
var ERR_ADD_OR_UPDATE_TAG_FAIL_CONTENT_TOO_LONG = pbutil.StaticBuffer{3, 29, 3, 3} // 3-3

// 标签太端
var ERR_ADD_OR_UPDATE_TAG_FAIL_CONTENT_TOO_SHORT = pbutil.StaticBuffer{3, 29, 3, 4} // 3-4

// 目标标签已满，无法添加
var ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_TAG_FULL = pbutil.StaticBuffer{3, 29, 3, 5} // 3-5

// 目标不可以是自己
var ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_CANT_ME = pbutil.StaticBuffer{3, 29, 3, 6} // 3-6

// 服务器繁忙，请稍后再试
var ERR_ADD_OR_UPDATE_TAG_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 29, 3, 7} // 3-7

func NewS2cOtherTagMeMsg(record []byte, tag []byte) pbutil.Buffer {
	msg := &S2COtherTagMeProto{
		Record: record,
		Tag:    tag,
	}
	return NewS2cOtherTagMeProtoMsg(msg)
}

var s2c_other_tag_me = [...]byte{29, 4} // 4
func NewS2cOtherTagMeProtoMsg(object *S2COtherTagMeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_tag_me[:], "s2c_other_tag_me")

}

func NewS2cDeleteTagMsg(tags []string) pbutil.Buffer {
	msg := &S2CDeleteTagProto{
		Tags: tags,
	}
	return NewS2cDeleteTagProtoMsg(msg)
}

var s2c_delete_tag = [...]byte{29, 6} // 6
func NewS2cDeleteTagProtoMsg(object *S2CDeleteTagProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_delete_tag[:], "s2c_delete_tag")

}
