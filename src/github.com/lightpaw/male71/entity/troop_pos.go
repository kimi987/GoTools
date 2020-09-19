package entity

import (
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func newTroopPosArray(n uint64) []*TroopPos {
	captains := make([]*TroopPos, n)
	for i := range captains {
		captains[i] = &TroopPos{}
	}

	return captains
}

type TroopPos struct {
	captain *Captain

	xIndex int32
}

func (t *TroopPos) Captain() *Captain {
	return t.captain
}

func (t *TroopPos) XIndex() int32 {
	return t.xIndex
}

func (t *TroopPos) IsEmpty() bool {
	return t.captain == nil
}

func (hero *Hero) EncodeInvaseCaptainInfo(t *TroopPos, fullSoldier bool) *shared_proto.CaptainInfoProto {
	if t.captain != nil {
		return t.captain.EncodeInvaseCaptainInfo(hero, fullSoldier, t.xIndex)
	}
	return nil
}

func (hero *Hero) EncodeDefenseCaptainInfo(t *TroopPos, fullSoldier bool) *shared_proto.CaptainInfoProto {
	if t.captain != nil {
		return t.captain.EncodeDefenseCaptainInfo(hero, fullSoldier, t.xIndex)
	}
	return nil
}

func (hero *Hero) EncodeCopyDefenseCaptainInfo(t *TroopPos, fullSoldier bool) *shared_proto.CaptainInfoProto {
	if t.captain != nil {
		return t.captain.EncodeCopyDefenseCaptainInfo(hero, fullSoldier, t.xIndex)
	}
	return nil
}

func (hero *Hero) EncodeAssistCaptainInfo(t *TroopPos, fullSoldier bool) *shared_proto.CaptainInfoProto {
	if t.captain != nil {
		return t.captain.EncodeAssistCaptainInfo(hero, fullSoldier, t.xIndex)
	}
	return nil
}

func (t *TroopPos) EncodeCaptainInfo(fullSoldier bool) *shared_proto.CaptainInfoProto {
	if t.captain != nil {
		return t.captain.EncodeCaptainInfo(fullSoldier, t.xIndex)
	}
	return nil
}

func doEncodeCaptainPos(tps []*TroopPos) (captainIds, xIndex []int32) {
	if n := len(tps); n > 0 {
		captainIds = make([]int32, n)
		xIndex = make([]int32, n)
		for i, ts := range tps {
			if ts.captain != nil {
				captainIds[i] = u64.Int32(ts.captain.Id())
				xIndex[i] = ts.xIndex
			}
		}
	}

	return
}

func doUnmarshalCaptainPos(tps []*TroopPos, captainIds, xIndex []int32, getCaptain func(uint64) *Captain) {
	n := imath.Minx(len(tps), len(captainIds), len(xIndex))
	for i := 0; i < n; i++ {
		captainId := captainIds[i]
		tps[i].unmarshal(captainId, xIndex[i], getCaptain)

		if captainId > 0 {
			// 一个队伍不允许重复的武将上阵
			start := i + 1
			for i := start; i < n; i++ {
				if captainIds[i] == captainId {
					captainIds[i] = 0
				}
			}
		}
	}
}

func (tp *TroopPos) unmarshal(captainId, xIndex int32, getCaptain func(uint64) *Captain) {
	if captainId > 0 {
		captain := getCaptain(u64.FromInt32(captainId))
		if captain != nil {
			tp.captain = captain
			tp.xIndex = xIndex
		}
	}
}

func doSetTroopPos(tps []*TroopPos, captains []*Captain, xIndex []int32) {
	n := imath.Minx(len(tps), len(captains), len(xIndex))
	for i := 0; i < n; i++ {
		tp := tps[i]
		tp.captain = captains[i]
		tp.xIndex = xIndex[i]
	}
}
