package unmarshal

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/util/msg"
	"github.com/pkg/errors"
)

func S2cMsgString(data []byte) (int, int, string, error) {

	// 读取一个varint32的消息长度
	_, data, err := msg.ReadVarint(data)
	if err != nil {
		return 0, 0, "", errors.Wrapf(err, "Print读取消息长度出错")
	}

	// 读取2个varint32的消息号
	moduleID, sequenceID, data, err := msg.ReadMsg(data)
	if err != nil {
		return moduleID, sequenceID, "", errors.Wrapf(err, "Print读取消息出错")
	}

	if util.IsCompressMsg(moduleID, sequenceID) {
		newData, err := util.UncompressMsg(sequenceID, data)
		if err != nil {
			return moduleID, sequenceID, "", errors.Wrapf(err, "Print读取消息出错，解压失败，%d-%d, len: %d", moduleID, sequenceID, len(data))
		}

		moduleID, sequenceID, data, err = msg.ReadMsg(newData)
		if err != nil {
			return moduleID, sequenceID, "", errors.Wrapf(err, "Print读取消息出错(解压缩后)，%d-%d, len: %d", moduleID, sequenceID, len(data))
		}
	}

	// 读取剩余的数据，交给后面switch进行处理
	proto, err := service.PrintObject(moduleID, sequenceID, data)
	if err != nil {
		return moduleID, sequenceID, "", errors.Wrapf(err, "PrintObject，解析出错")
	}

	return moduleID, sequenceID, proto.String(), nil
}
