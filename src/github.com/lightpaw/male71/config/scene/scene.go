package scene

// 战斗场景
//gogen:config
type CombatScene struct {
	_ struct{} `file:"地图/战斗场景.txt"`

	Id         string `validator:"string>0"` // 场景id
	Name       string `validator:"string>0"` // 名字
	MapRes     string `validator:"string>0"` // 地图资源
	WallMapRes string `validator:"string>0"` // 带城墙的地图资源
}
