package module

import "github.com/lightpaw/male7/gen/iface"

func NewModules(
	activity iface.ActivityModule,
	baizhan iface.BaiZhanModule,
	chat iface.ChatModule,
	clientconfig iface.ClientConfigModule,
	country iface.CountryModule,
	depot iface.DepotModule,
	dianquan iface.DianquanModule,
	domestic iface.DomesticModule,
	dungeon iface.DungeonModule,
	equipment iface.EquipmentModule,
	farm iface.FarmModule,
	fishing iface.FishingModule,
	garden iface.GardenModule,
	gem iface.GemModule,
	guild iface.GuildModule,
	hebi iface.HebiModule,
	mail iface.MailModule,
	military iface.MilitaryModule,
	mingc iface.MingcModule,
	mingcwar iface.MingcWarModule,
	misc iface.MiscModule,
	promotion iface.PromotionModule,
	question iface.QuestionModule,
	randomevent iface.RandomEventModule,
	rank iface.RankModule,
	redpacket iface.RedPacketModule,
	region iface.RegionModule,
	relation iface.RelationModule,
	secrettower iface.SecretTowerModule,
	shop iface.ShopModule,
	strategy iface.StrategyModule,
	stress iface.StressModule,
	survey iface.SurveyModule,
	tag iface.TagModule,
	task iface.TaskModule,
	teach iface.TeachModule,
	tower iface.TowerModule,
	vip iface.VipModule,
	xiongnu iface.XiongNuModule,
	xuanyuan iface.XuanyuanModule,
	zhanjiang iface.ZhanJiangModule,
	zhengwu iface.ZhengWuModule,

) *Modules {
	m := &Modules{}
	m.activity = activity
	m.baizhan = baizhan
	m.chat = chat
	m.clientconfig = clientconfig
	m.country = country
	m.depot = depot
	m.dianquan = dianquan
	m.domestic = domestic
	m.dungeon = dungeon
	m.equipment = equipment
	m.farm = farm
	m.fishing = fishing
	m.garden = garden
	m.gem = gem
	m.guild = guild
	m.hebi = hebi
	m.mail = mail
	m.military = military
	m.mingc = mingc
	m.mingcwar = mingcwar
	m.misc = misc
	m.promotion = promotion
	m.question = question
	m.randomevent = randomevent
	m.rank = rank
	m.redpacket = redpacket
	m.region = region
	m.relation = relation
	m.secrettower = secrettower
	m.shop = shop
	m.strategy = strategy
	m.stress = stress
	m.survey = survey
	m.tag = tag
	m.task = task
	m.teach = teach
	m.tower = tower
	m.vip = vip
	m.xiongnu = xiongnu
	m.xuanyuan = xuanyuan
	m.zhanjiang = zhanjiang
	m.zhengwu = zhengwu

	return m
}

//gogen:iface
type Modules struct {
	activity iface.ActivityModule

	baizhan iface.BaiZhanModule

	chat iface.ChatModule

	clientconfig iface.ClientConfigModule

	country iface.CountryModule

	depot iface.DepotModule

	dianquan iface.DianquanModule

	domestic iface.DomesticModule

	dungeon iface.DungeonModule

	equipment iface.EquipmentModule

	farm iface.FarmModule

	fishing iface.FishingModule

	garden iface.GardenModule

	gem iface.GemModule

	guild iface.GuildModule

	hebi iface.HebiModule

	mail iface.MailModule

	military iface.MilitaryModule

	mingc iface.MingcModule

	mingcwar iface.MingcWarModule

	misc iface.MiscModule

	promotion iface.PromotionModule

	question iface.QuestionModule

	randomevent iface.RandomEventModule

	rank iface.RankModule

	redpacket iface.RedPacketModule

	region iface.RegionModule

	relation iface.RelationModule

	secrettower iface.SecretTowerModule

	shop iface.ShopModule

	strategy iface.StrategyModule

	stress iface.StressModule

	survey iface.SurveyModule

	tag iface.TagModule

	task iface.TaskModule

	teach iface.TeachModule

	tower iface.TowerModule

	vip iface.VipModule

	xiongnu iface.XiongNuModule

	xuanyuan iface.XuanyuanModule

	zhanjiang iface.ZhanJiangModule

	zhengwu iface.ZhengWuModule
}

func (m *Modules) ActivityModule() iface.ActivityModule {
	return m.activity
}

func (m *Modules) BaiZhanModule() iface.BaiZhanModule {
	return m.baizhan
}

func (m *Modules) ChatModule() iface.ChatModule {
	return m.chat
}

func (m *Modules) ClientConfigModule() iface.ClientConfigModule {
	return m.clientconfig
}

func (m *Modules) CountryModule() iface.CountryModule {
	return m.country
}

func (m *Modules) DepotModule() iface.DepotModule {
	return m.depot
}

func (m *Modules) DianquanModule() iface.DianquanModule {
	return m.dianquan
}

func (m *Modules) DomesticModule() iface.DomesticModule {
	return m.domestic
}

func (m *Modules) DungeonModule() iface.DungeonModule {
	return m.dungeon
}

func (m *Modules) EquipmentModule() iface.EquipmentModule {
	return m.equipment
}

func (m *Modules) FarmModule() iface.FarmModule {
	return m.farm
}

func (m *Modules) FishingModule() iface.FishingModule {
	return m.fishing
}

func (m *Modules) GardenModule() iface.GardenModule {
	return m.garden
}

func (m *Modules) GemModule() iface.GemModule {
	return m.gem
}

func (m *Modules) GuildModule() iface.GuildModule {
	return m.guild
}

func (m *Modules) HebiModule() iface.HebiModule {
	return m.hebi
}

func (m *Modules) MailModule() iface.MailModule {
	return m.mail
}

func (m *Modules) MilitaryModule() iface.MilitaryModule {
	return m.military
}

func (m *Modules) MingcModule() iface.MingcModule {
	return m.mingc
}

func (m *Modules) MingcWarModule() iface.MingcWarModule {
	return m.mingcwar
}

func (m *Modules) MiscModule() iface.MiscModule {
	return m.misc
}

func (m *Modules) PromotionModule() iface.PromotionModule {
	return m.promotion
}

func (m *Modules) QuestionModule() iface.QuestionModule {
	return m.question
}

func (m *Modules) RandomEventModule() iface.RandomEventModule {
	return m.randomevent
}

func (m *Modules) RankModule() iface.RankModule {
	return m.rank
}

func (m *Modules) RedPacketModule() iface.RedPacketModule {
	return m.redpacket
}

func (m *Modules) RegionModule() iface.RegionModule {
	return m.region
}

func (m *Modules) RelationModule() iface.RelationModule {
	return m.relation
}

func (m *Modules) SecretTowerModule() iface.SecretTowerModule {
	return m.secrettower
}

func (m *Modules) ShopModule() iface.ShopModule {
	return m.shop
}

func (m *Modules) StrategyModule() iface.StrategyModule {
	return m.strategy
}

func (m *Modules) StressModule() iface.StressModule {
	return m.stress
}

func (m *Modules) SurveyModule() iface.SurveyModule {
	return m.survey
}

func (m *Modules) TagModule() iface.TagModule {
	return m.tag
}

func (m *Modules) TaskModule() iface.TaskModule {
	return m.task
}

func (m *Modules) TeachModule() iface.TeachModule {
	return m.teach
}

func (m *Modules) TowerModule() iface.TowerModule {
	return m.tower
}

func (m *Modules) VipModule() iface.VipModule {
	return m.vip
}

func (m *Modules) XiongNuModule() iface.XiongNuModule {
	return m.xiongnu
}

func (m *Modules) XuanyuanModule() iface.XuanyuanModule {
	return m.xuanyuan
}

func (m *Modules) ZhanJiangModule() iface.ZhanJiangModule {
	return m.zhanjiang
}

func (m *Modules) ZhengWuModule() iface.ZhengWuModule {
	return m.zhengwu
}
