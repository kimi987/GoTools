package resdata

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

// 掉落奖励

//gogen:config
type PlunderPrize struct {
	_  struct{} `file:"杂项/掉落_奖励.txt"`
	Id uint64

	Plunder                *Plunder                    `default:"nullable"` // 掉落，可为空
	GuildLevelPrizeGroupId uint64                      `validator:"uint"`
	GuildLevelPrizes       map[uint64]*GuildLevelPrize `head:"-"`
	Prize                  *Prize // 奖励
}

func (*PlunderPrize) InitAll(filename string, dataMap map[uint64]*PlunderPrize, configs interface {
	GetGoodsDataArray() []*goods.GoodsData
}) {
	for _, g := range configs.GetGoodsDataArray() {
		if g.EffectType == shared_proto.GoodsEffectType_EFFECT_PARTS {
			if len(g.GoodsEffect.PartsPlunderPrizeId) > 1 {
				for _, id := range g.GoodsEffect.PartsPlunderPrizeId {
					data := dataMap[id]
					check.PanicNotTrue(data != nil, "物品碎片 %d-%s 配置掉落奖励不存在[%d], 请检查物品表和 %s", g.Id, g.Name, id, filename)

					g.GoodsEffect.PartsShowPrize = append(g.GoodsEffect.PartsShowPrize, data.Prize.Encode4Init())
				}
			}
		}
	}
}

func (d *PlunderPrize) Init(filename string, configs interface {
	GetGuildLevelPrizeArray() []*GuildLevelPrize
}) {
	d.GuildLevelPrizes = GuildLevelPrizeGroupMap(d.GuildLevelPrizeGroupId, configs)
}

func (d *PlunderPrize) GetPrize() *Prize {
	if d.Plunder != nil {
		return d.Plunder.Try()
	}
	return d.Prize
}

func (d *PlunderPrize) GetGuildPrize(guildLevel uint64) *Prize {
	p := d.GetPrize()
	gp := d.GuildLevelPrizes[guildLevel]
	if gp == nil {
		return p
	}

	return AppendPrize(p, gp.Prize)
}
