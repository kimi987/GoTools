package entity

import (
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
)

const year = time.Hour * 24 * 365

func NewBuffInfo(start time.Time, data *data.BuffEffectData, operId int64) *BuffInfo {
	b := &BuffInfo{}
	b.EffectData = data
	b.OperHeroId = operId
	b.StartTime = start

	if data.NoDuration {
		b.EndTime = start.Add(year)
		b.Permanent = true
	} else {
		b.EndTime = start.Add(data.KeepDuration)
	}

	return b
}

type BuffInfo struct {
	StartTime  time.Time
	EndTime    time.Time
	EffectData *data.BuffEffectData
	OperHeroId int64
	Permanent  bool // 永久 buff
}

// 身上带的所有 buff
type HeroBuff struct {
	buffs           map[uint64]*BuffInfo // 所有的 buff
	captainBuffStat *data.SpriteStat     // 武将 buff 属性效果
	pvpBuffStat     *data.SpriteStat     // pvp buff 属性效果
}

func NewHeroBuff() *HeroBuff {
	h := &HeroBuff{}
	h.buffs = make(map[uint64]*BuffInfo)
	return h
}

func (h *HeroBuff) getCaptainBuffStat() *data.SpriteStat {
	return h.captainBuffStat
}

func (h *HeroBuff) getPvpBuffStat() *data.SpriteStat {
	return h.pvpBuffStat
}

func (h *HeroBuff) Update() {
	var capStatBuffs []*data.SpriteStat
	var pvpStatBuffs []*data.SpriteStat
	for _, b := range h.buffs {
		if b.EffectData.EffectType == shared_proto.BuffEffectType_Buff_ET_sprite_stat {
			if b.EffectData.PvpBuff {
				pvpStatBuffs = append(pvpStatBuffs, b.EffectData.StatBuff)
			} else {
				capStatBuffs = append(capStatBuffs, b.EffectData.StatBuff)
			}
		}
	}
	h.captainBuffStat = data.AppendSpriteStat(capStatBuffs...)
	h.pvpBuffStat = data.AppendSpriteStat(pvpStatBuffs...)
}

func (h *HeroBuff) Add(newBuffData *data.BuffEffectData, operHeroId int64, ctime time.Time) (succ bool, newBuff, oldBuff *BuffInfo) {
	if old, ok := h.buffs[newBuffData.Group]; ok {
		// 调用方根据自己的逻辑做验证
		//if old.EffectData.Level > newBuffData.Level {
		//	return
		//} else if old.EffectData.Level == newBuffData.Level {
		//	if old.EndTime.After(ctime.Add(newBuffData.KeepDuration)) {
		//		return
		//	}
		//}

		oldBuff = old
		delete(h.buffs, newBuffData.Group)
	}

	newBuff = NewBuffInfo(ctime, newBuffData, operHeroId)
	h.buffs[newBuffData.Group] = newBuff
	succ = true
	return
}

func (h *HeroBuff) Del(g uint64, ctime time.Time) (oldBuff *BuffInfo) {
	oldBuff = h.buffs[g]
	delete(h.buffs, g)
	return
}

func (h *HeroBuff) Walk(f func(buff *BuffInfo)) {
	for _, buff := range h.buffs {
		if buff == nil {
			continue
		}

		f(buff)
	}
}

func (h *HeroBuff) BuffCount() int {
	return len(h.buffs)
}

func (h *HeroBuff) Buff(g uint64) (buff *BuffInfo) {
	return h.buffs[g]
}

func (h *HeroBuff) Buffs(t shared_proto.BuffEffectType) (buffs []*BuffInfo) {
	for _, b := range h.buffs {
		if b == nil {
			continue
		}
		if b.EffectData.EffectType == t {
			buffs = append(buffs, b)
		}
	}
	return
}

func (h *HeroBuff) encodeServer() *server_proto.HeroBuffServerProto {
	p := &server_proto.HeroBuffServerProto{}
	p.Buff = make(map[int32]*server_proto.BuffInfoServerProto, len(h.buffs))
	for k, v := range h.buffs {
		if v == nil {
			continue
		}
		p.Buff[int32(k)] = v.EncodeServer()
	}

	return p
}

func (h *HeroBuff) encode() *shared_proto.HeroBuffProto {
	p := &shared_proto.HeroBuffProto{}
	for _, v := range h.buffs {
		if v == nil {
			continue
		}
		p.Buff = append(p.Buff, v.Encode())
	}

	return p
}

// 武将unmarshal之前先解码数据
func (h *HeroBuff) unmarshal(p *server_proto.HeroBuffServerProto, datas *config.BuffEffectDataConfig, ctime time.Time) {
	if p == nil {
		return
	}

	for k, v := range p.Buff {
		if v == nil {
			continue
		}
		if _, ok := shared_proto.BuffEffectType_name[k]; k <= 0 || !ok {
			logrus.Errorf("heroBuff.unmarshal, proto 中的 buff 类型不存在。k:%v", k)
			continue
		}

		if data := datas.Get(v.BuffEffectId); data != nil {
			b := NewBuffInfo(timeutil.Unix64(v.StartTime), data, v.OperId)
			b.EndTime = timeutil.Unix64(v.EndTime)

			h.buffs[u64.FromInt32(k)] = b
		}

	}
	h.Update()
}

func (b *BuffInfo) EncodeServer() *server_proto.BuffInfoServerProto {
	p := &server_proto.BuffInfoServerProto{}
	p.StartTime = timeutil.Marshal64(b.StartTime)
	p.EndTime = timeutil.Marshal64(b.EndTime)
	p.BuffEffectId = b.EffectData.Id
	p.OperId = b.OperHeroId

	return p
}

func (b *BuffInfo) Encode() *shared_proto.BuffInfoProto {
	p := &shared_proto.BuffInfoProto{}
	p.StartTime = timeutil.Marshal32(b.StartTime)
	p.EndTime = timeutil.Marshal32(b.EndTime)
	p.BuffEffectId = u64.Int32(b.EffectData.Id)

	return p
}
