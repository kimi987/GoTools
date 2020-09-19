package combat

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/imath"
)

func newWall(isAttacker bool, proto *shared_proto.CombatPlayerProto, maxWallBeenHurtPercent float64) *Wall {
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

	killSoldier int
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
