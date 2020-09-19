package captain

import (
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/spell"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

// 无名之辈（杂兵）
//gogen:config
type NamelessCaptainData struct {
	_ struct{} `file:"武将/无名武将.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	Name string

	Icon *icon.Icon `protofield:"IconId,%s.Id,string"` // 图标

	// 脊柱动画
	Spine string `default:" "`

	Male bool

	Race shared_proto.Race `type:"enum"`

	BaseSpell *spell.SpellFacadeData `protofield:",config.U64ToI32(%s.Id),int32"`

	Spell [] *spell.SpellFacadeData `protofield:",config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(%s)),int32"`

	// 技能战斗力系数
	spellFightAmountCoef []uint64

	InitRage uint64 `validator:"uint" default:"0"` // 初始怒气
}

func (c *NamelessCaptainData) Init(filename string) {
	spellFightAmountCoef := make([]uint64, len(c.Spell))
	for i, v := range c.Spell {

		// 技能战斗力系数
		spellFightAmountCoef[i] = v.FightAmountCoef
		if i > 0 {
			spellFightAmountCoef[i] += spellFightAmountCoef[i-1]
		}

		check.PanicNotTrue(v.BuildingEffect == nil, "%s配置的%v %s的技能中包含内政技能", filename)
	}

	c.spellFightAmountCoef = spellFightAmountCoef
}

func (c *NamelessCaptainData) GetSpellFightAmountCoef(unlockSpellCount uint64) uint64 {
	if unlockSpellCount > 0 && len(c.spellFightAmountCoef) > 0 {
		index := u64.Int(unlockSpellCount - 1)
		if index < len(c.spellFightAmountCoef) {
			return c.spellFightAmountCoef[index]
		}
		return c.spellFightAmountCoef[len(c.spellFightAmountCoef)-1]
	}

	return 0
}

//gogen:config
type CaptainRarityData struct {
	_ struct{} `file:"武将/稀有度.txt"`
	_ struct{} `protogen:"true"`

	Id   uint64 `validator:"int>0"` // 稀有度id
	Name string                     // 稀有度名称

	Color shared_proto.Quality `default:"3" protofield:"-"` // 稀有度品质颜色

	// 武将属性品质系数
	Coef float64
}

func GetCaptainData(obj *resdata.ResCaptainData) *CaptainData {
	return obj.GetObject().(*CaptainData)
}

//gogen:config
type CaptainData struct {
	_ struct{} `file:"武将/武将.txt"`
	_ struct{} `protogen:"true"`

	// 武将id
	Id uint64

	// 武将稀有度 S,SR,SSR...
	Rarity *CaptainRarityData `protofield:",config.U64ToI32(%s.Id),int32"`

	// 名字
	Name string

	// 图标
	Icon *icon.Icon `protofield:"IconId,%s.Id,string"`

	// 脊柱动画
	Spine string `default:" "`

	// 描述
	Desc string

	// 默认职业
	Race *race.RaceData `protofield:"Race,%s.Race,Race"`

	// 在已有该武将的情况下，给的奖励
	PrizeIfHas *resdata.Prize
	ResObject  *resdata.ResCaptainData `head:"id" protofield:"-"`

	ObtainWays    []uint64 `validator:"int" default:"nullable"` // 获得的途径
	FishingObtain bool                                          // true表示钓鱼可获得

	Sound string `default:" "` // 语音

	// 普攻
	BaseSpell *spell.SpellFacadeData `protofield:",config.U64ToI32(%s.Id),int32"`

	// 星数
	Star []*CaptainStarData `head:"-"`

	InitRage uint64 `validator:"uint" default:"0"` // 初始怒气

	GiftData *promdata.EventLimitGiftData `head:"-" protofield:"-"` // 事件礼包

	firstStar *CaptainStarData

	initTrainexp uint64
}

func (c *CaptainData) SetInitTrainExp(toSet uint64) {
	c.initTrainexp = toSet
}

func (c *CaptainData) GetInitTrainExp() uint64 {
	return c.initTrainexp
}

func (c *CaptainData) Init(filename string, configs interface {
	GetCaptainStarDataArray() []*CaptainStarData
	EventLimitGiftConfig() *promdata.EventLimitGiftConfig
}) {

	var stars []*CaptainStarData
	for _, star := range configs.GetCaptainStarDataArray() {
		if star.CaptainId == c.Id {
			stars = append(stars, star)
		}
	}

	for i := 1; i < len(stars); i++ {
		prev := stars[i-1]
		cur := stars[i]
		prev.nextStar = cur

		cur.AddedStat = data.DiffStat(cur.SpriteStat, prev.SpriteStat)

		check.PanicNotTrue(prev.Star == uint64(i), "武将 %d 配置的武将星数错误，必须从1星开始连续配置")
		check.PanicNotTrue(prev.Star+1 == stars[i].Star, "武将 %d 配置的武将星数错误，必须从1星开始连续配置")
	}
	check.PanicNotTrue(len(stars) > 0, "武将 %d 没有配置对应的星数")

	c.Star = stars
	c.GiftData = configs.EventLimitGiftConfig().GetCaptainGift(c.Rarity.Id)
	c.firstStar = stars[0]

	c.ResObject.InitObject(c)
}

func (c *CaptainData) GetStar(star uint64) *CaptainStarData {
	idx := u64.Sub(star, 1)
	if idx < uint64(len(c.Star)) {
		return c.Star[idx]
	}
	return c.Star[len(c.Star)-1]
}

func (c *CaptainData) GetFirstStar() *CaptainStarData {
	return c.firstStar
}

func CalculateCaptainStarId(captainId, star uint64) uint64 {
	return captainId*100 + star
}

//gogen:config
type CaptainStarData struct {
	_ struct{} `file:"武将/武将星数.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id        uint64 `head:"-,CalculateCaptainStarId(%s.CaptainId%c %s.Star)" protofield:"-"`
	CaptainId uint64 `protofield:"-"`
	Star      uint64

	// 武将星级系数
	Coef float64

	// 星级属性
	SpriteStat *data.SpriteStat `protofield:"-"`

	// 1升2取2星的属性数据
	AddedStat *data.SpriteStat `head:"-"`

	// 1星的当成激活的消耗，1升2取2星的消耗数据
	Cost *resdata.Cost

	// 武将技能列表
	Spell []*spell.SpellFacadeData `protofield:",config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(%s)),int32"`

	// 技能战斗力系数
	spellFightAmountCoef []uint64

	// true 表示有内政技能
	hasBuildingEffectSpell bool

	buildingEffectSpell [][]*sub.BuildingEffectData

	nextStar *CaptainStarData
}

func (c *CaptainStarData) Init(filename string) {

	spellFightAmountCoef := make([]uint64, len(c.Spell))
	bes := make([][]*sub.BuildingEffectData, len(c.Spell))
	for i, v := range c.Spell {

		// 技能战斗力系数
		spellFightAmountCoef[i] = v.FightAmountCoef
		if i > 0 {
			spellFightAmountCoef[i] += spellFightAmountCoef[i-1]
		}

		// 内政技能
		var prev []*sub.BuildingEffectData
		if i > 0 {
			prev = bes[i-1]
		}

		newArray := prev
		if v.BuildingEffect != nil {
			c.hasBuildingEffectSpell = true

			newArray = make([]*sub.BuildingEffectData, len(prev)+1)
			copy(newArray, prev)
			newArray[len(prev)] = v.BuildingEffect
		}

		bes[i] = newArray
	}

	c.spellFightAmountCoef = spellFightAmountCoef
	c.buildingEffectSpell = bes
}

func (c *CaptainStarData) HasBuildingEffectSpell() bool {
	return c.hasBuildingEffectSpell
}

func (c *CaptainStarData) GetBuildingEffectSpell(unlockSpellCount uint64) []*sub.BuildingEffectData {
	if c.hasBuildingEffectSpell && unlockSpellCount > 0 {
		index := u64.Int(unlockSpellCount - 1)
		if index < len(c.buildingEffectSpell) {
			return c.buildingEffectSpell[index]
		}
		return c.buildingEffectSpell[len(c.buildingEffectSpell)-1]
	}

	return nil
}

func (c *CaptainStarData) GetSpellFightAmountCoef(unlockSpellCount uint64) uint64 {
	if unlockSpellCount > 0 && len(c.spellFightAmountCoef) > 0 {
		index := u64.Int(unlockSpellCount - 1)
		if index < len(c.spellFightAmountCoef) {
			return c.spellFightAmountCoef[index]
		}
		return c.spellFightAmountCoef[len(c.spellFightAmountCoef)-1]
	}

	return 0
}

func (c *CaptainStarData) GetNextStar() *CaptainStarData {
	return c.nextStar
}

//gogen:config
type CaptainFriendshipData struct {
	_ struct{} `file:"武将/武将羁绊.txt"`
	_ struct{} `protogen:"true"`

	Id   uint64 // 羁绊id
	Name string // 羁绊名
	Desc string // 描述
	Tips string // 叹号弹窗提示

	// 羁绊武将id
	Captains []*CaptainData `protofield:",config.U64a2I32a(GetCaptainDataKeyArray(%s)),int32"`

	// 行军加速
	MoveSpeedRate float64 `validator:"float64" protofield:"-"`

	// 给所有武将加属性
	AllStat *data.SpriteStat `default:"nullable" protofield:"-"`

	// 给职业加属性
	Race     []shared_proto.Race `validator:"string" type:"enum" protofield:"-"`
	RaceStat []*data.SpriteStat  `validator:"uint,duplicate" default:"nullable" protofield:"-"`

	raceStatMap map[shared_proto.Race]*data.SpriteStat
	raceStat    []*data.RaceStat

	// 词条描述
	EffectDesc   []string
	EffectAmount []uint64 `validator:"uint,duplicate"`
}

func (d *CaptainFriendshipData) Init(filename string) {
	check.PanicNotTrue(len(d.Race) == len(d.RaceStat), "%s %v-%v 配置的职业个数跟职业属性个数不一致", filename, d.Id, d.Name)

	raceStatMap := make(map[shared_proto.Race]*data.SpriteStat)
	for i, r := range d.Race {
		raceStatMap[r] = data.AppendSpriteStat(d.RaceStat[i], d.AllStat)
	}

	if d.AllStat != nil {
		for _, r := range race.Array {
			if raceStatMap[r] == nil {
				raceStatMap[r] = d.AllStat
			}
		}
	}

	var raceStat []*data.RaceStat
	for k, v := range raceStatMap {
		raceStat = append(raceStat, data.NewRaceStat(k, v))
	}

	d.raceStatMap = raceStatMap
	d.raceStat = raceStat
}

func (d *CaptainFriendshipData) GetStat(race shared_proto.Race) *data.SpriteStat {
	return d.raceStatMap[race]
}

func (d *CaptainFriendshipData) GetRaceStat() []*data.RaceStat {
	return d.raceStat
}

func (d *CaptainFriendshipData) IsValidRace(race shared_proto.Race) bool {
	if d.AllStat != nil {
		return true
	}

	for _, r := range d.Race {
		if r == race {
			return true
		}
	}

	return false
}
