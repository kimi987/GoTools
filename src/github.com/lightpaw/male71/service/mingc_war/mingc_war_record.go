package mingc_war

import (
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/gen/iface"
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/pb/server_proto"
)

type McWarFightRecord struct {
	guilds      map[int64]*McWarFightGuildRecord
	atkYinliang uint64
	defYinliang uint64
}

type McWarFightGuildRecord struct {
	gid           int64
	joinedCount   uint64
	killedAmount  uint64
	woundedAmount uint64
	destroyed     uint64
	joinedTroops  map[int64]bool
}

func NewMcWarFightGuildRecord(gid int64) *McWarFightGuildRecord {
	c := &McWarFightGuildRecord{}
	c.gid = gid
	c.joinedTroops = make(map[int64]bool)
	return c
}

func (r *McWarFightGuildRecord) addTroop(tid int64) {
	r.joinedTroops[tid] = false
}

func (r *McWarFightRecord) unmarshal(p *server_proto.McWarFightRecordServerProto) {
	r.guilds = make(map[int64]*McWarFightGuildRecord)
	for _, gp := range p.GuildRecord {
		g := NewMcWarFightGuildRecord(gp.Gid)
		g.unmarshal(gp)
		r.guilds[gp.Gid] = g
	}
}

func (r *McWarFightRecord) encodeServer() *server_proto.McWarFightRecordServerProto {
	p := &server_proto.McWarFightRecordServerProto{}
	for _, g := range r.guilds {
		p.GuildRecord = append(p.GuildRecord, g.encodeServer())
	}

	return p
}

func (r *McWarFightRecord) encode(mc *MingcObj, timeout bool, ctime time.Time,
	guildGetter func(id int64, dep iface.ServiceDep) *shared_proto.GuildBasicProto,
	dep iface.ServiceDep) (*shared_proto.McWarFightRecordProto, *shared_proto.McWarTroopsInfoProto) {
	p := &shared_proto.McWarFightRecordProto{}
	atk := r.guilds[mc.atkId]
	p.McId = u64.Int32(mc.id)
	if atk != nil {
		p.Atk = atk.encode(guildGetter, dep)
	}

	def := r.guilds[mc.defId]
	if def != nil {
		p.Def = def.encode(guildGetter, dep)
	}
	for _, gid := range mc.astAtkList {
		g := r.guilds[gid]
		if g != nil {
			p.AstAtk = append(p.AstAtk, g.encode(guildGetter, dep))
		}
	}
	for _, gid := range mc.astDefList {
		g := r.guilds[gid]
		if g != nil {
			p.AstDef = append(p.AstDef, g.encode(guildGetter, dep))
		}
	}
	p.AtkYinliang = u64.Int32(r.atkYinliang)
	p.DefYinliang = u64.Int32(r.defYinliang)
	p.AtkWin = mc.scene.atkWin
	p.Timeout = timeout

	fightRunStartTime := mc.scene.startTime.Add(mc.dep.Datas().MingcMiscData().FightPrepareDuration)
	p.FightDuration = timeutil.DurationMarshal32(ctime.Sub(fightRunStartTime))

	troops := mc.scene.troopsRank.troopRankMap
	t := &shared_proto.McWarTroopsInfoProto{
		TroopsInfo: make([]*shared_proto.McWarTroopInfoProto, 0, len(troops)),
	}
	for _, v := range troops {
		t.TroopsInfo = append(t.TroopsInfo, v.encode4Info())
	}
	return p, t
}

func (r *McWarFightGuildRecord) encode(guildGetter func(id int64, dep iface.ServiceDep) *shared_proto.GuildBasicProto, dep iface.ServiceDep) *shared_proto.McWarFightGuildRecordProto {
	p := &shared_proto.McWarFightGuildRecordProto{}
	p.Guild = guildGetter(r.gid, dep)
	p.JoinedCount = int32(len(r.joinedTroops))
	p.KilledAmount = u64.Int32(r.killedAmount)
	p.WoundedAmount = u64.Int32(r.woundedAmount)
	p.Destroyed = u64.Int32(r.destroyed)
	return p
}

func (r *McWarFightGuildRecord) encodeServer() *server_proto.McWarFightGuildRecordServerProto {
	p := &server_proto.McWarFightGuildRecordServerProto{}
	p.Gid = r.gid
	p.KilledAmount = r.killedAmount
	p.WoundedAmount = r.woundedAmount
	p.Destroyed = r.destroyed
	p.JoinedTroops = r.joinedTroops
	return p
}

func (r *McWarFightGuildRecord) unmarshal(p *server_proto.McWarFightGuildRecordServerProto) {
	r.gid = p.Gid
	r.killedAmount = p.KilledAmount
	r.woundedAmount = p.WoundedAmount
	r.destroyed = p.Destroyed
	r.joinedTroops = p.JoinedTroops
}
