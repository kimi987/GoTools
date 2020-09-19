package resdata

import "github.com/lightpaw/male7/config/data"

//gogen:config
type ConditionPlunder struct {
	_ struct{} `file:"杂项/条件掉落.txt"`

	Id uint64

	DefPlunder *Plunder

	CondItem []*ConditionPlunderItem
}

func (c *ConditionPlunder) GetPrize(heroLevel uint64) *Prize {

	// 从左到右，找到第一个符合条件的，如果找不到，使用DefPlunder
	for _, item := range c.CondItem {
		if item.HeroLevel.Compare(heroLevel) {
			return item.Plunder.Try()
		}
	}
	return c.DefPlunder.Try()
}

//gogen:config
type ConditionPlunderItem struct {
	_ struct{} `file:"杂项/条件掉落项.txt"`

	Id uint64

	HeroLevel *data.CompareCondition

	Plunder *Plunder `protofield:"-"`
}
