package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/pvetroop"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/data"
)

func newPveTroop(data *pvetroop.PveTroopData) *PveTroop {
	return &PveTroop{
		data:     data,
		captains: newTroopPosArray(data.Capacity),
	}
}

// pve队伍
type PveTroop struct {
	data     *pvetroop.PveTroopData
	captains []*TroopPos
}

func (troop *PveTroop) TroopData() *pvetroop.PveTroopData {
	return troop.data
}

func (troop *PveTroop) TotalFullFightAmount() (total uint64) {
	tfa := data.NewTroopFightAmount()
	for _, pos := range troop.captains {
		c := pos.captain
		if c != nil {
			tfa.Add(c.FullSoldierFightAmount())
		}
	}
	return tfa.ToU64()
}

func (troop *PveTroop) CaptainRaces() (races []shared_proto.Race) {
	races = make([]shared_proto.Race, len(troop.captains))

	for idx, pos := range troop.captains {
		c := pos.captain
		if c != nil {
			races[idx] = c.Race().Race
		}
	}

	return
}

func (troop *PveTroop) CaptainCount() (count uint64) {
	for _, c := range troop.captains {
		if c.captain != nil {
			count++
		}
	}

	return
}

func (troop *PveTroop) Captains() []*TroopPos {
	return troop.captains
}

func (troop *PveTroop) SetCaptain(captains []*Captain, xIndex []int32) {
	doSetTroopPos(troop.captains, captains, xIndex)
}

func (troop *PveTroop) AddCaptain(captain *Captain) bool {

	var emptyPos *TroopPos
	for _, c := range troop.captains {
		if c.captain == nil {
			if emptyPos == nil {
				emptyPos = c
			}
		} else if captain == c.captain {
			// 已经在里面了
			logrus.Errorf("竟然存在 AddCaptain 的时候，已经有武将在这支队伍里面了")
			return false
		}
	}

	if emptyPos == nil {
		return false
	}

	emptyPos.captain = captain

	return true
}

func (troop *PveTroop) Encode() *shared_proto.HeroPveTroopProto {
	proto := &shared_proto.HeroPveTroopProto{}

	proto.Type = troop.data.PveTroopType
	proto.Captains, proto.XIndex = doEncodeCaptainPos(troop.captains)
	return proto
}

func (troop *PveTroop) unmarshal(getCaptainFunc func(cid uint64) *Captain, proto *shared_proto.HeroPveTroopProto) {
	doUnmarshalCaptainPos(troop.captains, proto.Captains, proto.XIndex, getCaptainFunc)
}
