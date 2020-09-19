package data

//gogen:config
type UnlockCondition struct {
	_ struct{} `proto:"shared_proto.UnlockConditionProto"`
	_ struct{} `int:"uint"`

	// 需要君主等级
	RequiredHeroLevel uint64 `default:"0"`

	// 需要主城等级
	RequiredBaseLevel uint64 `default:"0"`

	// 需要联盟等级
	RequiredGuildLevel uint64 `default:"0"`

	// 需要vip等级
	RequiredVipLevel uint64 `default:"0"`

	// true表示空的条件
	isEmptyCondition bool
}

func (c *UnlockCondition) IsEmptyCondition() bool {
	return c.isEmptyCondition
}

func (c *UnlockCondition) Init(filename string) {

	c.isEmptyCondition = c.RequiredHeroLevel <= 0 &&
		c.RequiredBaseLevel <= 0 &&
		c.RequiredGuildLevel <= 0
}
