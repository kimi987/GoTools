package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/combine"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

// 玩家开启了的装备合成

func newHeroOpenCombineEquip() *HeroOpenCombineEquip {
	return &HeroOpenCombineEquip{
		openCombineEquipMap: make(map[uint64]struct{}),
	}
}

type HeroOpenCombineEquip struct {
	// 开启了的所有装备合成
	openCombineEquipMap map[uint64]struct{}
}

func (h *HeroOpenCombineEquip) IsAllOpen(combineDatas *combine.EquipCombineDatas) (allOpen bool) {
	return len(h.openCombineEquipMap) >= combineDatas.Len()
}

func (h *HeroOpenCombineEquip) IsOpen(data *combine.EquipCombineData) (isOpen bool) {
	_, ok := h.openCombineEquipMap[data.Id]
	return ok
}

func (h *HeroOpenCombineEquip) Open(data *combine.EquipCombineData) (openSuccess bool) {
	if h.IsOpen(data) {
		return false
	}

	h.openCombineEquipMap[data.Id] = struct{}{}

	return true
}

func (h *HeroOpenCombineEquip) EncodeClient() *shared_proto.HeroOpenCombineEquipProto {
	proto := &shared_proto.HeroOpenCombineEquipProto{}

	proto.OpenEquipCombine = make([]int32, 0, len(h.openCombineEquipMap))
	for id := range h.openCombineEquipMap {
		proto.OpenEquipCombine = append(proto.OpenEquipCombine, u64.Int32(id))
	}

	return proto
}

func (h *HeroOpenCombineEquip) encodeServer() *server_proto.HeroOpenCombineEquipServerProto {
	proto := &server_proto.HeroOpenCombineEquipServerProto{}

	proto.OpenEquipCombine = make([]uint64, 0, len(h.openCombineEquipMap))
	for id := range h.openCombineEquipMap {
		proto.OpenEquipCombine = append(proto.OpenEquipCombine, id)
	}

	return proto
}

func (h *HeroOpenCombineEquip) unmarshal(proto *server_proto.HeroOpenCombineEquipServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	for _, id := range proto.GetOpenEquipCombine() {
		data := datas.GetEquipCombineData(id)
		if data != nil {
			h.Open(data)
		}
	}
}
