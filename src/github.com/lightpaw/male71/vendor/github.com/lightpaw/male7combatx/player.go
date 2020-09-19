package combatx

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
)

func newCombatPlayer(genIndex func() int32, id int64, isAttacker bool, proto *shared_proto.CombatPlayerProto, config *Config, misc *misc) (player *CombatPlayer, err error) {
	troopsArray, err := newTroopsArray(genIndex, isAttacker, proto.GetTroops(), config, misc)
	if err != nil {
		return nil, errors.Wrap(err, "创建部队报错")
	}

	player = &CombatPlayer{
		id:           id,
		proto:        proto,
		isAttacker:   isAttacker,
		troops:       troopsArray,
		wall:         newWall(isAttacker, proto, config.WallBeenHurtLostMaxPercent, config.InitWallX),
		troopRaceMap: make(map[shared_proto.Race][]*Troops),
	}

	return
}

// 君主战斗
type CombatPlayer struct {
	id         int64
	proto      *shared_proto.CombatPlayerProto
	isAttacker bool
	troops     []*Troops
	wall       *Wall // 城墙

	troopRaceMap map[shared_proto.Race][]*Troops
}

func (p *CombatPlayer) TroopInitData() []*shared_proto.CombatTroopsInitProto {
	ps := make([]*shared_proto.CombatTroopsInitProto, 0, len(p.troops))
	for _, v := range p.troops {
		ps = append(ps, v.initProto)
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

func (p *CombatPlayer) isAllDead() bool {
	return p.isAllTroopDead() && !p.HasAliveWall()
}

func (p *CombatPlayer) isAllTroopDead() bool {
	for _, troops := range p.troops {
		if troops.isAlive() {
			return false
		}
	}

	return true
}
