package country

//gogen:config
type FamilyNameData struct {
	_ struct{} `file:"国家/姓氏荐国.txt"`
	_ struct{} `protogen:"true"`

	// id
	Id uint64

	// 姓氏
	Name string

	// 描述
	Desc string

	// 推荐国家
	RecommendCountry []*CountryData `protofield:",config.U64a2I32a(GetCountryDataKeyArray(%s)),int32"`
}
