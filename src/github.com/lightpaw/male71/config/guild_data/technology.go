package guild_data

import (
	"time"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

func GetTechnologyDataId(group, level uint64) uint64 {
	return group*10000 + level
}

func GetTechnologyGroupById(id uint64) uint64 {
	return id / 10000
}

func GetTechnologyLevelById(id uint64) uint64 {
	return id % 10000
}

//联盟科技
//gogen:config
type GuildTechnologyData struct {
	_ struct{} `file:"联盟/联盟科技.txt"`
	_ struct{} `proto:"shared_proto.GuildTechnologyDataProto"`
	_ struct{} `protoconfig:"guild_technology"`

	Id    uint64   `head:"-,GetTechnologyDataId(%s.Group%c %s.Level)"`
	Name  string
	Desc  []string `validator:"string,notAllNil,duplicate"`
	Icon  string
	Group uint64 // 科技分组
	Level uint64 // 科技等级

	RequireGuildLevel uint64 // 升级所需联盟等级

	// 升级所需建设值，1级放着0升1的升级所需建设值
	UpgradeBuilding uint64

	// 升级所需时长，1级放着0升1的升级所需时长
	UpgradeDuration time.Duration

	Cdrs []*GuildLevelCdrData `head:"-"`

	// 盟友协助减cd次数
	HelpCdr time.Duration `default:"1m"`

	// 升级效果
	Effect *sub.BuildingEffectData `default:"nullable"` // 加建筑属性

	BigBox *GuildBigBoxData `default:"nullable" protofield:",config.U64ToI32(%s.Id)"` // 宝箱科技

	prevLevel *GuildTechnologyData
	nextLevel *GuildTechnologyData
}

func (d *GuildTechnologyData) Init(filename string, configs interface {
	GetGuildLevelCdrDataArray() []*GuildLevelCdrData
}) {
	if d.Effect != nil {
		check.PanicNotTrue(d.BigBox == nil, "联盟科技配置%v 同一类型[Group:%v]等级[level:%v]的配置效果，不能同时有宝箱和建筑效果", filename, d.Group, d.Level)
	} else {
		check.PanicNotTrue(d.BigBox != nil, "联盟科技配置%v 同一类型[Group:%v]等级[level:%v]的配置效果，必须有宝箱，或者建筑效果（2选1）", filename, d.Group, d.Level)

		check.PanicNotTrue(d.BigBox.TechLevel == 0, "联盟科技配置%v 同一类型[Group:%v]等级[level:%v]的配置效果有宝箱，同一个宝箱不能配置在多个科技等级", filename, d.Group, d.Level)
		d.BigBox.TechLevel = d.Level
	}

	check.PanicNotTrue(d.UpgradeDuration >= 0, "联盟科技配置%v 同一类型[Group:%v]等级[level:%v]的升级所需时间必须>0", filename, d.Group, d.Level)

	// 缓存好加速数据
	var cdr []*GuildLevelCdrData
	for _, c := range configs.GetGuildLevelCdrDataArray() {
		if c.Group == d.Group && c.Level == d.Level {
			cdr = append(cdr, c)

			check.PanicNotTrue(int(c.Times) == len(cdr), "联盟科技配置类型[Group:%v]等级[level:%v]的配置效果，找不到加速%v 次的数据，加速次数必须从1开始连续配置", d.Group, d.Level, len(cdr))
		}
	}
	d.Cdrs = cdr
}

func (*GuildTechnologyData) InitAll(filename string, array []*GuildTechnologyData) {

	groupMap := make(map[uint64]map[uint64]*GuildTechnologyData)
	for _, v := range array {
		levelMap := groupMap[v.Group]
		if levelMap == nil {
			levelMap = make(map[uint64]*GuildTechnologyData)
			groupMap[v.Group] = levelMap
		}

		check.PanicNotTrue(levelMap[v.Level] == nil, "联盟科技配置%v 同一类型[Type:%v]-[Group:%v]存在重复的等级[level:%v] id: %v", filename, v.Name, v.Group, v.Level, v.Id)
		levelMap[v.Level] = v
	}

	for k, levelMap := range groupMap {
		check.PanicNotTrue(len(levelMap) < 10000, "联盟科技配置%v 同一类型[Group:%v]最高等级是9999，当前配置已超出最大等级上限", filename, k)

		for i := 0; i < len(levelMap); i++ {
			level := uint64(i + 1)
			data := levelMap[level]
			check.PanicNotTrue(data != nil, "联盟科技配置%v 同一类型[Group:%v]等级[level:%v]的数据没找到，必须从1级开始连续配置", filename, k, level)

			if i > 0 {
				prevLevel := levelMap[uint64(i)]
				prevLevel.nextLevel = data
				data.prevLevel = prevLevel

				if data.Effect != nil {
					check.PanicNotTrue(prevLevel.Effect != nil, "联盟科技配置%v 同一类型[Group:%v]不同等级[level1:%v level2:%v]的配置效果必须类型一致，不能一个是宝箱，一个是建筑效果", filename, k, level-1, level)
				}
			}
		}
	}

}

func (d *GuildTechnologyData) GetPrevLevel() *GuildTechnologyData {
	return d.prevLevel
}

func (d *GuildTechnologyData) GetNextLevel() *GuildTechnologyData {
	return d.nextLevel
}

func (d *GuildTechnologyData) GetCdr(times uint64) *GuildLevelCdrData {
	if times > 0 && times <= uint64(len(d.Cdrs)) {
		return d.Cdrs[times-1]
	}

	return nil
}

type TechnologySlice []*GuildTechnologyData

func (p TechnologySlice) Len() int           { return len(p) }
func (p TechnologySlice) Less(i, j int) bool { return p[i].Id < p[j].Id }
func (p TechnologySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func GetDiffTechnology(oldArray, newArray []*GuildTechnologyData) (diff []*GuildTechnologyData) {

out0:
	for _, o := range oldArray {
		for _, n := range newArray {
			if o == n {
				continue out0
			}
		}
		diff = append(diff, o)
	}

out1:
	for _, n := range newArray {
		for _, o := range oldArray {
			if o == n {
				continue out1
			}
		}
		diff = append(diff, n)
	}

	return
}

func GetTechnologyEffects(array []*GuildTechnologyData) ([]*sub.BuildingEffectData) {
	effects := make([]*sub.BuildingEffectData, 0, len(array))
	for _, v := range array {
		if v.Effect != nil {
			effects = append(effects, v.Effect)
		}
	}
	return effects
}
