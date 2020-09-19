package idbytes

import (
	"encoding/hex"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity/heroid"
	"strconv"
)

func NewIdHolder(id int64) (result IdHolder) {
	return IdHolder{
		id:      id,
		idBytes: ToBytes(id),
	}
}

type IdHolder struct {
	id      int64
	idBytes []byte
}

func (h IdHolder) Id() int64 {
	return h.id
}

func (h IdHolder) IdBytes() []byte {
	return h.idBytes
}

func PlayerName(id int64) string {
	sid := heroid.GetSid(id)
	accountId := heroid.GetAccountId(id)
	return constants.PlayerNamePrefix + hex.EncodeToString(ToBytes(accountId)) + "s" + strconv.FormatUint(uint64(sid), 10)
}

func HeroBasicProto(id int64) *shared_proto.HeroBasicProto {
	return &shared_proto.HeroBasicProto{
		Id:    ToBytes(id),
		Name:  PlayerName(id),
		Level: 1,
	}
}
