package combat

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"math"
)

const Denominator float64 = 10000

func NewCombat(p *server_proto.CombatRequestServerProto) (*SingleCombat, error) {

	if p == nil {
		return nil, errors.Errorf("请求参数无效，传入的CombatRequestProto为nil")
	}

	if p.Attacker == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Attacker为nil")
	}

	if len(p.Attacker.Troops) == 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Attacker.Troops == 0")
	}

	if p.Defenser == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Defenser为nil")
	}

	if len(p.Defenser.Troops) == 0 && p.Defenser.WallStat == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Defenser.Troops == 0")
	}

	if p.MaxRound <= 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.MaxRound <= 0")
	}

	if len(p.Races) <= 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Races 长度 <= 0")
	}

	if p.GetMapXLen() <= 0 || p.GetMapXLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapXLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapXLen())
	}
	if p.GetMapYLen() <= 0 || p.GetMapYLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapYLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapYLen())
	}

	combat := &SingleCombat{
		misc:  newMisc(p),
		proto: p,
	}

	if err := combat.addPlayer(p.GetAttackerId(), p.Attacker, true); err != nil {
		return nil, err
	}

	if err := combat.addPlayer(p.GetDefenserId(), p.Defenser, false); err != nil {
		return nil, err
	}

	return combat, nil
}

// 单挑
type SingleCombat struct {
	*misc

	proto *server_proto.CombatRequestServerProto

	id int32

	attackerPlayer *CombatPlayer
	defenserPlayer *CombatPlayer
}

func (c *SingleCombat) genId() int32 {
	c.id++
	return c.id
}

func (c *SingleCombat) addPlayer(id int64, p *shared_proto.CombatPlayerProto, isAttacker bool) error {
	player, err := newCombatPlayer(c.genId, id, isAttacker, p, c.misc)
	if err != nil {
		return errors.Wrap(err, "SingleCombat.addPlayer失败")
	}

	if isAttacker {
		c.attackerPlayer = player
	} else {
		c.defenserPlayer = player
	}

	return nil
}

// --- 战斗 ---
func (c *SingleCombat) Calculate() (combatResult *shared_proto.CombatProto, err error) {
	sc, err := newCombat(c.misc, c.attackerPlayer, c.defenserPlayer)
	if err != nil {
		return nil, errors.Wrap(err, "newCombat报错了")
	}

	// 默认队列1战斗
	return sc.Calculate(1), nil
}
