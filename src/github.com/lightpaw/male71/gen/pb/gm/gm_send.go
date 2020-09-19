package gm

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
	MODULE_ID = 3

	C2S_LIST_CMD = 5

	C2S_GM = 1

	C2S_INVASE_TARGET_ID = 7
)

func NewS2cListCmdMsg(datas [][]byte) pbutil.Buffer {
	msg := &S2CListCmdProto{
		Datas: datas,
	}
	return NewS2cListCmdProtoMsg(msg)
}

func NewS2cListCmdMarshalMsg(datas [][]byte) pbutil.Buffer {
	msg := &S2CListCmdProto{
		Datas: datas,
	}
	return NewS2cListCmdProtoMsg(msg)
}

var s2c_list_cmd = [...]byte{3, 6} // 6
func NewS2cListCmdProtoMsg(object *S2CListCmdProto) pbutil.Buffer {

	return util.NewGzipCompressMsg(object, s2c_list_cmd[:], "s2c_list_cmd")

}

func NewS2cGmMsg(result string) pbutil.Buffer {
	msg := &S2CGmProto{
		Result: result,
	}
	return NewS2cGmProtoMsg(msg)
}

var s2c_gm = [...]byte{3, 2} // 2
func NewS2cGmProtoMsg(object *S2CGmProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_gm[:], "s2c_gm")

}

func NewS2cInvaseTargetIdMsg(target_id []byte, target_x int32, target_y int32) pbutil.Buffer {
	msg := &S2CInvaseTargetIdProto{
		TargetId: target_id,
		TargetX:  target_x,
		TargetY:  target_y,
	}
	return NewS2cInvaseTargetIdProtoMsg(msg)
}

var s2c_invase_target_id = [...]byte{3, 8} // 8
func NewS2cInvaseTargetIdProtoMsg(object *S2CInvaseTargetIdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_invase_target_id[:], "s2c_invase_target_id")

}
