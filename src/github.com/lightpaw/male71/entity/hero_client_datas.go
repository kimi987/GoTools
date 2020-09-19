package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func newClientDatas(maps *hero_maps) *HeroClientDatas {
	return &HeroClientDatas{
		int32Array: make([]int32, 20),
		clientKeys: maps.getOrCreateKeys(server_proto.HeroMapCategory_client_keys, false),
	}
}

const (
	Int32BitCount = 32
	KeyCount      = 256
)

type HeroClientDatas struct {
	int32Array []int32

	clientKeys *herokeys
}

func (h *HeroClientDatas) SetBool(index int, value bool) {
	if index < 0 || index >= len(h.int32Array)*Int32BitCount {
		logrus.WithField("len", len(h.int32Array)).WithField("index", index).WithField("value", value).Debugln("HeroClientDatas.SetInt32 越界了")
		return
	}

	arrayIndex, realIndex := toIndex(index)

	shiftValue := uint32(1) << realIndex

	if value {
		h.int32Array[arrayIndex] |= int32(shiftValue)
	} else {
		h.int32Array[arrayIndex] &= int32(^shiftValue)
	}
}

func toIndex(index int) (arrayIndex uint32, realIndex uint32) {
	arrayIndex = uint32(index / Int32BitCount)
	realIndex = uint32(index % Int32BitCount)
	return
}

func (h *HeroClientDatas) Bool(index int) bool {
	if index < 0 || index >= len(h.int32Array)*Int32BitCount {
		logrus.WithField("len", len(h.int32Array)).WithField("index", index).Debugln("HeroClientDatas.SetInt32 越界了")
		return false
	}

	arrayIndex, realIndex := toIndex(index)

	v := h.int32Array[arrayIndex] >> realIndex

	return v&1 == 1
}

func (h *HeroClientDatas) IsExistClientKey(ckType, key uint64) bool {
	return h.clientKeys.Exist(combineClientKeyType(ckType, key))
}

func (h *HeroClientDatas) SetClientKey(ckType, key uint64) {
	if h.clientKeys.Size() < KeyCount {
		h.clientKeys.Add(combineClientKeyType(ckType, key))
	}
}

func combineClientKeyType(ckType, key uint64) uint64 {
	return key<<4 | ckType
}

func (h *HeroClientDatas) unmarshal(proto *server_proto.HeroClientDatasServerProto) {
	if proto == nil {
		return
	}

	copy(h.int32Array, proto.IntValue)
}

func (h *HeroClientDatas) encodeClient() *shared_proto.HeroClientDatasProto {
	proto := &shared_proto.HeroClientDatasProto{}

	for index := Int32BitCount*len(h.int32Array) - 1; index >= 0; index-- {
		if h.Bool(index) {
			proto.IntValue = append(proto.IntValue, int32(index))
		}
	}

	for k := range h.clientKeys.internalMap {
		proto.ClientKeys = append(proto.ClientKeys, u64.Int32(k))
	}

	return proto
}

func (h *HeroClientDatas) encodeServer() *server_proto.HeroClientDatasServerProto {
	proto := &server_proto.HeroClientDatasServerProto{}

	proto.IntValue = h.int32Array

	return proto
}
