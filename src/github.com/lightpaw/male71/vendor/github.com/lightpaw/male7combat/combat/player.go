package combat

import (
	"fmt"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
)

func newCombatPlayer(genIndex func() int32, id int64, isAttacker bool, proto *shared_proto.CombatPlayerProto, misc *misc) (player *CombatPlayer, err error) {
	troopsArray, err := newTroopsArray(genIndex, isAttacker, proto.GetTroops(), misc)
	if err != nil {
		return nil, errors.Wrap(err, "创建部队报错")
	}

	player = &CombatPlayer{
		id:         id,
		proto:      proto,
		isAttacker: isAttacker,
		troops:     troopsArray,
		wall:       newWall(isAttacker, proto, misc.maxWallBeenHurtPercent),
	}

	return
}

// 战斗对象
type CombatPlayer struct {
	id             int64
	proto          *shared_proto.CombatPlayerProto
	isAttacker     bool
	troops         []*Troops
	curStep        int   // 最大步数
	winTimes       int   // 连胜次数
	lastFightQueue int   // 最近一次战斗在哪条路上面进行的
	wall           *Wall // 城墙
}

func (p *CombatPlayer) String() string {
	return fmt.Sprintf("CombatPlayer{id: %d, isAttacker: %v, curStep: %d, winTimes: %d, lastFightQueue: %d}", p.id, p.isAttacker, p.curStep, p.winTimes, p.lastFightQueue)
}

func (p *CombatPlayer) Id() int64 {
	return p.id
}

func (p *CombatPlayer) IdBytes() []byte {
	return p.proto.Hero.Id
}

func (p *CombatPlayer) TroopPos() []*shared_proto.CombatTroopsPosProto {
	ps := make([]*shared_proto.CombatTroopsPosProto, 0, len(p.troops))
	for _, v := range p.troops {
		ps = append(ps, v.pos)
	}
	return ps
}

func (p *CombatPlayer) Wall() *Wall {
	return p.wall
}

// 是否有活着的城墙
func (p *CombatPlayer) HasAliveWall() bool {
	return p.wall != nil && p.wall.isAlive()
}

func (p *CombatPlayer) AllDead() bool {
	for _, troops := range p.troops {
		if troops.isAlive() {
			return false
		}
	}

	return !p.HasAliveWall()
}

func (p *CombatPlayer) IsContinueWinLeave(continueWinMaxTimes int) bool {
	if continueWinMaxTimes <= 0 {
		// 没有连胜限制
		return false
	}
	return p.winTimes >= continueWinMaxTimes
}

func (p *CombatPlayer) IncreWinTimes() int {
	p.winTimes++
	return p.winTimes
}

func (p *CombatPlayer) WinTimes() int {
	return p.winTimes
}

func (p *CombatPlayer) CurStep() int {
	return p.curStep
}

func (p *CombatPlayer) SetCurStep(toSet int) {
	p.curStep = toSet
}

func (p *CombatPlayer) EncodeAliveSoldiers() *server_proto.AliveSoldierProto {
	proto := &server_proto.AliveSoldierProto{}

	proto.Id = p.id
	proto.AliveSoldier = make(map[int32]int32, len(p.troops))
	for _, v := range p.troops {
		proto.AliveSoldier[v.proto.Captain.GetId()] = int32(v.soldier)
	}

	return proto
}
