package regdata

//gogen:config
type RegionAreaData struct {
	_ struct{} `file:"地图/地区区域带.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	Name string

	// 区域范围
	Area *AreaData

	// 联盟工坊额外奖励
	WorkshopPrizeCoef float64 `validator:"float64"`
}
