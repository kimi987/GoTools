package unmarshal

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/util/msg"
	"github.com/pkg/errors"
)

func NewProtoUnmarshaller() protoUnmarshaller {
	var result protoUnmarshaller
	return result
}

type protoUnmarshaller int

// 总的解析消息成为MsgData方法. 所有模块的所有消息都在这里
func (protoUnmarshaller) Unmarshal(data []byte) (interface{}, error) {
	moduleId, msgId, proto, err := msg.ReadMsg(data)
	if err != nil {
		return nil, errors.Wrapf(err, "解析客户端上行消息失败，%d-%d, len: %d", moduleId, msgId, len(data))
	}

	if util.IsCompressMsg(moduleId, msgId) {
		newData, err := util.UncompressMsg(msgId, proto)
		if err != nil {
			return nil, errors.Wrapf(err, "客户端上行压缩消息，解压失败，%d-%d, len: %d", moduleId, msgId, len(data))
		}

		moduleId, msgId, proto, err = msg.ReadMsg(newData)
		if err != nil {
			return nil, errors.Wrapf(err, "解析客户端上行消息失败(解压缩后)，%d-%d, len: %d", moduleId, msgId, len(data))
		}
	}

	//logrus.Debugf("收到客户端上行消息，%d-%d", moduleId, msgId)

	return service.Unmarshal(moduleId, msgId, proto)
}
