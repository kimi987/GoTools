package location

import "github.com/lightpaw/male7/config/country"

//gogen:config
type LocationData struct {
	_ struct{} `file:"杂项/省市.txt"`
	_ struct{} `protogen:"true"`

	// id
	Id uint64

	// 名称
	Name string

	// 地域荐国
	RecommendCountry []*country.CountryData `protofield:",config.U64a2I32a(country.GetCountryDataKeyArray(%s)),int32"`
}
