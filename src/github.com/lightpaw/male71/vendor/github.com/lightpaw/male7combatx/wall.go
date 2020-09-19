package combatx

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/imath"
)

func newWall(isAttacker bool, proto *shared_proto.CombatPlayerProto, maxWallBeenHurtPercent float64, x int) *Wall {
	if isAttacker {
		return nil
	}

	if proto.WallStat == nil {
		return nil
	}

	// 防守方，且有城墙初始化血量
	wall := &Wall{}

	wall.maxLife = int(i32.Max(proto.WallStat.Strength, 1))
	wall.life = wall.maxLife
	wall.attack = float64(proto.WallStat.Attack)
	wall.defense = float64(proto.WallStat.Defense)
	wall.soldierCapcity = float64(proto.WallStat.SoldierCapcity)
	wall.fixDamage = int(proto.WallFixDamage)
	wall.maxBeenHurt = int(float64(wall.maxLife) * maxWallBeenHurtPercent)
	wall.x = x

	return wall
}

// 城墙
type Wall struct {
	life           int     // 血量
	maxLife        int     // 最大血量
	attack         float64 // 攻击
	defense        float64 // 防御
	soldierCapcity float64
	fixDamage      int // 固定死兵数
	maxBeenHurt    int // 每次扣血最多扣多少

	x int // 城墙的X坐标，y坐标全屏都是

	nextReleaseFrame int
	attackTimes      int

	killSoldier int

	delayDamage          []*DelayDamage
	nextDelayDamageFrame int
}

func (w *Wall) addKillSoldier(toAdd int) {
	w.killSoldier += toAdd
}

func (w *Wall) isAlive() bool {
	return w.life > 0
}

func (w *Wall) Life() int {
	return w.life
}

func (w *Wall) GetMaxBeenHurt() int {
	return w.maxBeenHurt
}

func (w *Wall) ReduceLife(toReduce int) {
	w.life = imath.Max(w.life-toReduce, 0)
}

func (target *Wall) addDelayDamage(caster spellcaster, damage, effectFrame int, spellId int32) {
	proto := &DelayDamage{
		caster:      caster,
		damage:      damage,
		effectFrame: effectFrame,
		spellId:     spellId,
	}

	target.nextDelayDamageFrame = IMin(target.nextDelayDamageFrame, effectFrame)

	for i, v := range target.delayDamage {
		if v == nil {
			target.delayDamage[i] = proto
			return
		}
	}

	target.delayDamage = append(target.delayDamage, proto)
}
