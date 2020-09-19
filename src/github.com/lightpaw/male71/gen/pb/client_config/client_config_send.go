package client_config

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
	MODULE_ID = 21

	C2S_CONFIG = 1

	C2S_SET_CLIENT_DATA = 4

	C2S_SET_CLIENT_KEY = 5
)

func NewS2cConfigMsg(data []byte) pbutil.Buffer {
	msg := &S2CConfigProto{
		Data: data,
	}
	return NewS2cConfigProtoMsg(msg)
}

var s2c_config = [...]byte{21, 2} // 2
func NewS2cConfigProtoMsg(object *S2CConfigProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_config[:], "s2c_config")

}

// 路径不存在
var ERR_CONFIG_FAIL_PATH_NOT_EXIST = pbutil.StaticBuffer{3, 21, 3, 1} // 3-1

func NewS2cSetClientKeyMsg(key_type int32, key_value int32) pbutil.Buffer {
	msg := &S2CSetClientKeyProto{
		KeyType:  key_type,
		KeyValue: key_value,
	}
	return NewS2cSetClientKeyProtoMsg(msg)
}

var s2c_set_client_key = [...]byte{21, 6} // 6
func NewS2cSetClientKeyProtoMsg(object *S2CSetClientKeyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_client_key[:], "s2c_set_client_key")

}
