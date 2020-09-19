package entity

import (
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timeutil"
)

func NewHeroHebi() *HeroHebi {
	h := &HeroHebi{}
	h.RobCdEndTime = time.Time{}
	return h
}

type HeroHebi struct {

	CaptainId uint64
	DailyRobCount uint64
	RobCdEndTime time.Time
	RobPosCdEndTime time.Time

	CurrentGoodsId uint64
}

func (m *HeroHebi) ResetDaily() {
	m.DailyRobCount = 0
}


func (h *HeroHebi) encode() *shared_proto.HeroHebiProto {
	p := &shared_proto.HeroHebiProto{}
	p.CaptainId = u64.Int32(h.CaptainId)
	p.RobCount = u64.Int32(h.DailyRobCount)
	p.RobCdEndTime = timeutil.Marshal32(h.RobCdEndTime)
	p.RobPosCdEndTime = timeutil.Marshal32(h.RobPosCdEndTime)
	p.CurrentGoodsId = u64.Int32(h.CurrentGoodsId)
	return p
}

func (h *HeroHebi) unmarshal(p *shared_proto.HeroHebiProto) {
	if p == nil {
		return
	}

	h.CaptainId = u64.FromInt32(p.CaptainId)
	h.DailyRobCount = u64.FromInt32(p.RobCount)
	h.RobCdEndTime = timeutil.Unix32(p.RobCdEndTime)
	h.RobPosCdEndTime = timeutil.Unix32(p.RobPosCdEndTime)
	h.CurrentGoodsId = u64.FromInt32(p.CurrentGoodsId)
}


