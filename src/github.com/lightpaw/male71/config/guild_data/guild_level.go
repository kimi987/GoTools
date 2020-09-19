package guild_data

import (
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

//gogen:config
type GuildLevelData struct {
	_     struct{} `file:"联盟/联盟等级.txt"`
	_     struct{} `proto:"shared_proto.GuildLevelProto"`
	_     struct{} `protoconfig:"GuildLevel"`
	Level uint64   `validator:"int>0"`

	MemberCount uint64

	ClassMemberCount []uint64 `validator:"uint,duplicate"`

	// 升级所需建设值，1级放着1升2的升级所需建设值
	UpgradeBuilding uint64

	// 升级所需时长，1级放着1升2的升级所需时长
	UpgradeDuration time.Duration

	Cdrs []*GuildLevelCdrData `head:"-"`

	nextLevel *GuildLevelData
}

func (d *GuildLevelData) NextLevel() *GuildLevelData {
	return d.nextLevel
}

func (d *GuildLevelData) Init(filename string, dataMap map[uint64]*GuildLevelData, configs interface {
	GetGuildClassLevelDataArray() []*GuildClassLevelData
	GetGuildLevelCdrDataArray() []*GuildLevelCdrData
}) {

	if d.Level > 1 {
		prevLevel := dataMap[d.Level-1]
		check.PanicNotTrue(prevLevel != nil, "%s 没有找到%v 级的帮派等级数据，等级必须从1开始连续配置", filename, d.Level-1)

		check.PanicNotTrue(prevLevel.MemberCount <= d.MemberCount, "%s 配置的%v 级的帮派成员数比上一级的少", filename, d.Level)

		prevLevel.nextLevel = d
	}

	check.PanicNotTrue(len(d.ClassMemberCount) == len(configs.GetGuildClassLevelDataArray()), "%s 配置的%v 级的帮派阶级成员个数必须更阶级数一致，阶级个数参考帮派阶级表", filename, d.Level)
	check.PanicNotTrue(len(d.ClassMemberCount) >= 2, "%s 配置的%v 级帮派等级数据，阶级个数必须>=2，至少要有帮主和帮众", filename, d.Level)

	check.PanicNotTrue(d.ClassMemberCount[0] == 0, "%s 配置的%v 级帮派等级数据，帮众人数上限必须=0", filename, d.Level)
	check.PanicNotTrue(d.ClassMemberCount[len(d.ClassMemberCount)-1] == 1, "%s 配置的%v 级帮派等级数据，帮主人数上限必须=1", filename, d.Level)

	// 缓存好加速数据
	var cdr []*GuildLevelCdrData
	for _, c := range configs.GetGuildLevelCdrDataArray() {
		if c.Group == 0 && c.Level == d.Level {
			cdr = append(cdr, c)

			check.PanicNotTrue(int(c.Times) == len(cdr), "%v 级的帮派升级加速配置无效，找不到加速%v 次的数据，加速次数必须从1开始连续配置", d.Level, len(cdr))
		}
	}
	check.PanicNotTrue(len(cdr) > 0, "%v 级的帮派升级加速没有配置!%d", d.Level, len(cdr))
	d.Cdrs = cdr
}

func (d *GuildLevelData) GetCdr(times uint64) *GuildLevelCdrData {
	if times > 0 && times <= uint64(len(d.Cdrs)) {
		return d.Cdrs[times-1]
	}

	return nil
}

func (d *GuildLevelData) GetClassMemberCount(level uint64) uint64 {
	index := u64.Sub(level, 1)
	if index >= 0 && index < uint64(len(d.ClassMemberCount)) {
		return d.ClassMemberCount[index]
	}

	return d.ClassMemberCount[len(d.ClassMemberCount)-1]
}

func GuildLevelCdrId(group, level, times uint64) uint64 {
	return group * 100000 + level*100 + times
}

//gogen:config
type GuildLevelCdrData struct {
	_ struct{} `file:"联盟/联盟升级加速.txt"`
	_ struct{} `proto:"shared_proto.GuildLevelCdrProto"`

	Id uint64 `head:"-,GuildLevelCdrId(%s.Group%c %s.Level%c %s.Times)" protofield:"-"`

	Group uint64 `default:"0" validator:"uint" protofield:"-"`

	// 帮派等级
	Level uint64 `protofield:"-"`

	// 加速次数
	Times uint64

	// 本次加速消耗的建设值
	Cost uint64

	// 减多少CD
	CDR time.Duration `protofield:"Cdr"`
}
