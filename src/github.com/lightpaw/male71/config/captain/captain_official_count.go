package captain

import (
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type CaptainOfficialCountData struct {
	_ struct{} `file:"武将/武将官职数量.txt"`

	HeroLevel  uint64   `key:"true"`                            // 有数量改变时的君主等级
	OfficialId []uint64 `validator:"int>0,duplicate,notAllNil"` // 官职 ID
	MaxCount   []uint64 `validator:"uint,duplicate,notAllNil"`  // 最大数量

	maxGongXunOfficial *CaptainOfficialData
}

func (d *CaptainOfficialCountData) GetMaxGongXunOfficialData() *CaptainOfficialData {
	return d.maxGongXunOfficial
}

func (*CaptainOfficialCountData) InitAll(array []*CaptainOfficialCountData, configs interface {
	GetCaptainOfficialDataArray() []*CaptainOfficialData
	GetCaptainOfficialData(uint64) *CaptainOfficialData
}) {
	var prev *CaptainOfficialCountData
	for _, v := range array {
		check.PanicNotTrue(len(v.OfficialId) == len(v.MaxCount), "武将官职数量.txt: official_id 和 max_count 必须一一对应")
		check.PanicNotTrue(len(v.OfficialId) == len(configs.GetCaptainOfficialDataArray()), "武将官职数量.txt: official_id 字段必须包含所有的武将官职")
		for _, id := range v.OfficialId {
			data := configs.GetCaptainOfficialData(id)
			check.PanicNotTrue(data != nil, "武将官职数量.txt: official_id 字段的武将官职 id 不存在:%v", id)

			if v.maxGongXunOfficial == nil || v.maxGongXunOfficial.NeedGongxun < data.NeedGongxun {
				v.maxGongXunOfficial = data
			}
		}

		check.PanicNotTrue(v.maxGongXunOfficial != nil, "武将官职数量.txt[Level:%d]: v.maxGongXunOfficial == nil", v.HeroLevel)

		if prev != nil {
			check.PanicNotTrue(v.HeroLevel > prev.HeroLevel, "武将官职数量.txt: 每行的 hero_level 必须大于上一行.%v", v.HeroLevel)
			for i, vcount := range v.MaxCount {
				check.PanicNotTrue(vcount >= prev.MaxCount[i], "武将官职数量.txt: 每行的 max_count 必须不能小于上一行.%v", vcount)
			}
		} else {
			check.PanicNotTrue(v.HeroLevel == 1, "武将官职数量.txt: 第一行的 hero_level 必须 == 1。 hero_level:%v", v.HeroLevel)
		}

		prev = v
	}
}
