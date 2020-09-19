package entity

import (
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func newHeroCaptainFriendship() *HeroCaptainFriendship {
	c := &HeroCaptainFriendship{}
	c.friendship = make(map[uint64]*captain.CaptainFriendshipData)
	c.raceStat = make(map[shared_proto.Race]*data.SpriteStat)

	return c
}

type HeroCaptainFriendship struct {
	// 激活的武将羁绊
	friendship map[uint64]*captain.CaptainFriendshipData

	// 羁绊效果
	moveSpeedRate float64 // 移动速度

	raceStat map[shared_proto.Race]*data.SpriteStat
}

func (c *HeroCaptainFriendship) getRaceStat(race shared_proto.Race) *data.SpriteStat {
	return c.raceStat[race]
}

func (c *HeroCaptainFriendship) tryAdd(toAdd *captain.CaptainFriendshipData) bool {
	if d := c.friendship[toAdd.Id]; d != nil {
		return false
	}

	c.friendship[toAdd.Id] = toAdd

	if toAdd.MoveSpeedRate > 0 {
		c.moveSpeedRate += toAdd.MoveSpeedRate
	}

	for _, v := range toAdd.GetRaceStat() {
		oldStat := c.raceStat[v.Race]
		c.raceStat[v.Race] = data.AppendSpriteStat(oldStat, v.SpriteStat)
	}

	return true
}

func (c *HeroCaptainFriendship) calculateProsperity() {

	var moveSpeedRate float64
	raceStatBuilder := make(map[shared_proto.Race]*data.SpriteStatBuilder)
	for _, f := range c.friendship {
		if f.MoveSpeedRate > 0 {
			moveSpeedRate += f.MoveSpeedRate
		}

		for _, v := range f.GetRaceStat() {
			builder := raceStatBuilder[v.Race]
			if builder == nil {
				builder = data.NewSpriteStatBuilder()
				raceStatBuilder[v.Race] = builder
			}
			builder.Add(v.SpriteStat)
		}
	}

	c.moveSpeedRate = moveSpeedRate

	for k, v := range raceStatBuilder {
		if v != nil {
			c.raceStat[k] = v.Build()
		}
	}
}
