package gm

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"io/ioutil"
	"math/rand"
	data2 "github.com/lightpaw/male7/config/data"
)

func (m *GmModule) newMailGmGroup() *gm_group {
	return &gm_group{
		tab: "邮件",
		handler: []*gm_handler{
			newIntHandler("给自己发系统(联盟)邮件", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendMail(hc, false)
				}
			}),
			newIntHandler("给自己发系统(联盟)邮件(奖励)", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendMail(hc, true)
				}
			}),
			newIntHandler("给自己发联盟邮件（奖励）", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendMail(hc, true)
				}
			}),
			newIntHandler("给自己发文字战报", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendReportTextMail(hc)
				}
			}),
			newIntHandler("给自己发打架战报", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendFightReportMail(hc)
				}
			}),
			newIntHandler("给自己发集结打架战报", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendAssemblyFightReportMail(hc)
				}
			}),
			newIntHandler("给自己发集结结算战报", "1", func(amount int64, hc iface.HeroController) {
				n := i64.Max(amount, 1)
				for i := int64(0); i < n; i++ {
					m.sendAssemblyDoneMail(hc)
				}
			}),
		},
	}
}

func (m *GmModule) sendMail(hc iface.HeroController, hasPrize bool) {
	title := "GM邮件: 无奖励邮件"
	var prize *shared_proto.PrizeProto
	if hasPrize {
		switch rand.Intn(3) {
		case 1:
			title = "GM邮件: 资源奖励邮件"
			prize = resdata.NewPrizeBuilder().AddSafeResource(100, 100, 100, 100).Build().Encode()
		case 2:
			title = "GM邮件: 物品奖励邮件"
			prize = resdata.NewPrizeBuilder().AddGoods(m.datas.GoodsData().Array[rand.Intn(len(m.datas.GoodsData().Array))], 1).Build().Encode()
		default:
			title = "GM邮件: 装备奖励邮件"
			prize = resdata.NewPrizeBuilder().AddEquipment(m.datas.EquipmentData().Array[rand.Intn(len(m.datas.EquipmentData().Array))], 1).Build().Encode()
		}
	}

	// 联盟邮件
	if guildId, _ := hc.LockGetGuildId(); guildId != 0 {
		if g := m.sharedGuildService.GetSnapshot(guildId); g != nil {
			proto := &shared_proto.MailProto{}
			proto.Icon = rand.Int31n(5)
			proto.Title = title
			proto.Text = "测试邮件正文"
			proto.Keep = false
			proto.Report = nil
			proto.Prize = prize

			proto.GuildName = g.Name
			proto.GuildFlagName = g.FlagName

			m.modules.MailModule().SendProtoMail(hc.Id(), proto, m.time.CurrentTime())
			return
		}
	}

	m.modules.MailModule().SendMail(hc.Id(), uint64(rand.Int31n(5)), title, "测试邮件正文", false, nil, prize, m.time.CurrentTime())
}

func (m *GmModule) sendReportTextMail(hc iface.HeroController) {

	var mailProto *shared_proto.MailProto
	data := m.datas.MailData().Array[rand.Intn(len(m.datas.MailData().Array))]
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		mailProto = data.NewTextMail(shared_proto.MailType_MailReport)
		mailProto.Text = data.NewTextFields().
			WithFields("attacker", hero.Name()).
			WithFields("defenser", hero.Name()).
			WithFields("assister", hero.Name()).JsonString()
		mailProto.Text = data.NewTextFields().
			WithFields("attacker", hero.Name()).
			WithFields("defenser", hero.Name()).
			WithFields("assister", hero.Name()).JsonString()
	})

	m.modules.MailModule().SendReportMail(hc.Id(), mailProto, m.time.CurrentTime())
}

func (m *GmModule) sendFightReportMail(hc iface.HeroController) {

	datas := m.datas

	var mailProto *shared_proto.MailProto
	data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroProto := &shared_proto.ReportHeroProto{
			Id:         hero.IdBytes(),
			Name:       hero.Name(),
			Level:      int32(hero.Level()),
			Head:       hero.Head(),
			BaseRegion: i64.Int32(hero.BaseRegion()),
			BaseX:      imath.Int32(hero.BaseX()),
			BaseY:      imath.Int32(hero.BaseY()),
			IsTent:     false,
		}

		for _, captain := range hero.Military().Captains() {
			heroProto.TotalSoldier += u64.Int32(captain.SoldierCapcity())

			cp := &shared_proto.ReportCaptainProto{
				Index:         imath.Int32(len(heroProto.Captains) + 1),
				Captain:       captain.EncodeCaptainInfo(false, 0),
				CombatSoldier: rand.Int31n(u64.Int32(captain.SoldierCapcity())),
			}

			tfa := data2.NewTroopFightAmount()
			if cp.CombatSoldier > 0 {
				cp.AliveSoldier = rand.Int31n(cp.CombatSoldier)
				heroProto.AliveSoldier += cp.AliveSoldier

				fm := captain.FightAmount()
				cp.FightAmount = u64.Int32(fm)
				tfa.Add(fm)
			}
			heroProto.TotalFightAmount = tfa.ToI32()

			// 武将详情
			heroProto.Captains = append(heroProto.Captains, cp)
			if len(heroProto.Captains) >= 5 {
				break
			}
		}

		proto := &shared_proto.FightReportProto{}
		proto.AttackerWin = rand.Int31n(2) == 0
		proto.ReplayUrl = m.randomFightUrl()

		proto.Attacker = heroProto
		proto.Defenser = heroProto

		proto.ShowPrize = &shared_proto.PrizeProto{}
		proto.ShowPrize.Gold = rand.Int31()
		proto.ShowPrize.Food = rand.Int31()
		proto.ShowPrize.Wood = rand.Int31()
		proto.ShowPrize.Stone = rand.Int31()
		proto.ShowPrize.JadeOre = rand.Int31()
		resdata.SetPrizeProtoIsNotEmpty(proto.ShowPrize)

		proto.FightType = rand.Int31n(3)
		proto.FightX, proto.FightY = heroProto.BaseX, heroProto.BaseY
		proto.FightTargetId = proto.Attacker.Id
		proto.FightTargetName = proto.Attacker.Name
		proto.FightTargetFlagName = proto.Attacker.GuildFlagName

		proto.AttackerDesc, proto.DefenserDesc = datas.TextHelp().MailReportWinnerDesc.Text.KeysOnlyJson(), datas.TextHelp().MailReportLoserDesc.Text.KeysOnlyJson()
		//switch rand.Intn(4) {
		//case 0:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionInvadeActDesc.OneText, m.datas.MailConfig().RegionInvadeDefDesc.OneText
		//case 1:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionAssistActDesc.OneText, m.datas.MailConfig().RegionAssistDefDesc.OneText
		//case 2:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionExpelActDesc.OneText, m.datas.MailConfig().RegionExpelDefDesc.OneText
		//default:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionBackActDesc.OneText, m.datas.MailConfig().RegionBackDefDesc.OneText
		//}

		mailProto = data.NewTextMail(shared_proto.MailType_MailReport)
		mailProto.Text = data.NewTextFields().
			WithFields("attacker", proto.Attacker.Name).
			WithFields("defenser", proto.Attacker.Name).
			WithFields("assister", proto.Attacker.Name).JsonString()
		mailProto.Report = proto

		//data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
		//mailProto = data.NewArgsTextMail(hero.Name(), hero.Name(), hero.Name())
		//mailProto.Title = "GM邮件: " + mailProto.Title
		//mailProto.Report = proto

		return

	})

	m.modules.MailModule().SendReportMail(hc.Id(), mailProto, m.time.CurrentTime())
}

func (m *GmModule) randomFightUrl() string {
	files, _ := ioutil.ReadDir("temp/")
	if len(files) > 0 {
		return "{{local}}/" + files[rand.Intn(len(files))].Name()
	}
	return ""
}

func (m *GmModule) sendAssemblyFightReportMail(hc iface.HeroController) {

	datas := m.datas

	var mailProto *shared_proto.MailProto
	data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroProto := &shared_proto.ReportHeroProto{
			Id:         hero.IdBytes(),
			Name:       hero.Name(),
			Level:      int32(hero.Level()),
			Head:       hero.Head(),
			BaseRegion: i64.Int32(hero.BaseRegion()),
			BaseX:      imath.Int32(hero.BaseX()),
			BaseY:      imath.Int32(hero.BaseY()),
			IsTent:     false,
		}

		for _, captain := range hero.Military().Captains() {
			heroProto.TotalSoldier += u64.Int32(captain.SoldierCapcity())

			cp := &shared_proto.ReportCaptainProto{
				Index:         imath.Int32(len(heroProto.Captains) + 1),
				Captain:       captain.EncodeCaptainInfo(false, 0),
				CombatSoldier: rand.Int31n(u64.Int32(captain.SoldierCapcity())),
			}

			tfa := data2.NewTroopFightAmount()
			if cp.CombatSoldier > 0 {
				cp.AliveSoldier = rand.Int31n(cp.CombatSoldier)
				heroProto.AliveSoldier += cp.AliveSoldier

				fm := captain.FightAmount()
				cp.FightAmount = u64.Int32(fm)
				tfa.Add(fm)
			}
			heroProto.TotalFightAmount = tfa.ToI32()

			// 武将详情
			heroProto.Captains = append(heroProto.Captains, cp)
			if len(heroProto.Captains) >= 5 {
				break
			}
		}

		proto := &shared_proto.FightReportProto{}
		proto.AttackerWin = rand.Int31n(2) == 0
		proto.ReplayUrl = m.randomFightUrl()

		proto.Attacker = heroProto
		proto.Defenser = heroProto

		proto.ShowPrize = &shared_proto.PrizeProto{}
		proto.ShowPrize.Gold = rand.Int31()
		proto.ShowPrize.Food = rand.Int31()
		proto.ShowPrize.Wood = rand.Int31()
		proto.ShowPrize.Stone = rand.Int31()
		proto.ShowPrize.JadeOre = rand.Int31()
		resdata.SetPrizeProtoIsNotEmpty(proto.ShowPrize)

		proto.FightType = rand.Int31n(3)
		proto.FightX, proto.FightY = heroProto.BaseX, heroProto.BaseY
		proto.FightTargetId = proto.Attacker.Id
		proto.FightTargetName = proto.Attacker.Name
		proto.FightTargetFlagName = proto.Attacker.GuildFlagName

		proto.AttackerDesc, proto.DefenserDesc = datas.TextHelp().MailReportWinnerDesc.Text.KeysOnlyJson(), datas.TextHelp().MailReportLoserDesc.Text.KeysOnlyJson()
		//switch rand.Intn(4) {
		//case 0:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionInvadeActDesc.OneText, m.datas.MailConfig().RegionInvadeDefDesc.OneText
		//case 1:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionAssistActDesc.OneText, m.datas.MailConfig().RegionAssistDefDesc.OneText
		//case 2:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionExpelActDesc.OneText, m.datas.MailConfig().RegionExpelDefDesc.OneText
		//default:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionBackActDesc.OneText, m.datas.MailConfig().RegionBackDefDesc.OneText
		//}

		proto.AttackerTroopCount = 3
		proto.AttackerTroopTotalCount = 5

		proto.DefenserTroopCount = 3
		proto.DefenserTroopTotalCount = 3

		proto.AttackerTroopWinTimes = 5
		proto.DefenserTroopWinTimes = 3

		for i := 0; i < 5; i++ {
			fight := &shared_proto.AssemblyFightProto{}
			fight.AttackerId = hero.IdBytes()
			fight.DefenserId = hero.IdBytes()

			fight.AttackerFightAmount = rand.Int31n(100000)
			fight.DefenserFightAmount = rand.Int31n(100000)

			fight.AttackerWin = rand.Int31n(2) == 0

			fight.WinTimes = int32(i)

			fight.Share = &shared_proto.CombatShareProto{
				Link: m.randomFightUrl(),
				Type: shared_proto.CombatType_SINGLE_X,
			}

			proto.Fight = append(proto.Fight, fight)
		}

		mailProto = data.NewTextMail(shared_proto.MailType_MailAssemblyReport)
		mailProto.Text = data.NewTextFields().
			WithFields("attacker", proto.Attacker.Name).
			WithFields("defenser", proto.Attacker.Name).
			WithFields("assister", proto.Attacker.Name).JsonString()
		mailProto.Report = proto

		//data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
		//mailProto = data.NewArgsTextMail(hero.Name(), hero.Name(), hero.Name())
		//mailProto.Title = "GM邮件: " + mailProto.Title
		//mailProto.Report = proto

		return

	})

	m.modules.MailModule().SendReportMail(hc.Id(), mailProto, m.time.CurrentTime())
}

func (m *GmModule) sendAssemblyDoneMail(hc iface.HeroController) {

	datas := m.datas

	var mailProto *shared_proto.MailProto
	data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroProto := &shared_proto.ReportHeroProto{
			Id:         hero.IdBytes(),
			Name:       hero.Name(),
			Level:      int32(hero.Level()),
			Head:       hero.Head(),
			BaseRegion: i64.Int32(hero.BaseRegion()),
			BaseX:      imath.Int32(hero.BaseX()),
			BaseY:      imath.Int32(hero.BaseY()),
			IsTent:     false,
		}

		for _, captain := range hero.Military().Captains() {
			heroProto.TotalSoldier += u64.Int32(captain.SoldierCapcity())

			cp := &shared_proto.ReportCaptainProto{
				Index:         imath.Int32(len(heroProto.Captains) + 1),
				Captain:       captain.EncodeCaptainInfo(false, 0),
				CombatSoldier: rand.Int31n(u64.Int32(captain.SoldierCapcity())),
			}

			tfa := data2.NewTroopFightAmount()
			if cp.CombatSoldier > 0 {
				cp.AliveSoldier = rand.Int31n(cp.CombatSoldier)
				heroProto.AliveSoldier += cp.AliveSoldier

				fm := captain.FightAmount()
				cp.FightAmount = u64.Int32(fm)
				tfa.Add(fm)
			}
			heroProto.TotalFightAmount = tfa.ToI32()

			// 武将详情
			heroProto.Captains = append(heroProto.Captains, cp)
			if len(heroProto.Captains) >= 5 {
				break
			}
		}

		proto := &shared_proto.FightReportProto{}
		proto.AttackerWin = rand.Int31n(2) == 0
		proto.ReplayUrl = m.randomFightUrl()

		proto.Attacker = heroProto
		proto.Defenser = heroProto

		proto.ShowPrize = &shared_proto.PrizeProto{}
		proto.ShowPrize.Gold = rand.Int31()
		proto.ShowPrize.Food = rand.Int31()
		proto.ShowPrize.Wood = rand.Int31()
		proto.ShowPrize.Stone = rand.Int31()
		proto.ShowPrize.JadeOre = rand.Int31()
		resdata.SetPrizeProtoIsNotEmpty(proto.ShowPrize)

		proto.FightType = rand.Int31n(3)
		proto.FightX, proto.FightY = heroProto.BaseX, heroProto.BaseY
		proto.FightTargetId = proto.Attacker.Id
		proto.FightTargetName = proto.Attacker.Name
		proto.FightTargetFlagName = proto.Attacker.GuildFlagName

		proto.AttackerDesc, proto.DefenserDesc = datas.TextHelp().MailReportWinnerDesc.Text.KeysOnlyJson(), datas.TextHelp().MailReportLoserDesc.Text.KeysOnlyJson()
		//switch rand.Intn(4) {
		//case 0:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionInvadeActDesc.OneText, m.datas.MailConfig().RegionInvadeDefDesc.OneText
		//case 1:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionAssistActDesc.OneText, m.datas.MailConfig().RegionAssistDefDesc.OneText
		//case 2:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionExpelActDesc.OneText, m.datas.MailConfig().RegionExpelDefDesc.OneText
		//default:
		//	proto.AttackerDesc, proto.DefenserDesc = datas.MailConfig().RegionBackActDesc.OneText, m.datas.MailConfig().RegionBackDefDesc.OneText
		//}

		proto.AttackerTroopCount = 3
		proto.AttackerTroopTotalCount = 5

		for i := 0; i < 5; i++ {
			fight := &shared_proto.AssemblyFightProto{}
			fight.Attacker = proto.Attacker
			fight.ShowPrize = &shared_proto.PrizeProto{}
			fight.ShowPrize.Gold = rand.Int31()
			fight.ShowPrize.Stone = rand.Int31()
			resdata.SetPrizeProtoIsNotEmpty(fight.ShowPrize)

			proto.Fight = append(proto.Fight, fight)
		}

		mailProto = data.NewTextMail(shared_proto.MailType_MailAssemblyRobFinished)
		mailProto.Text = data.NewTextFields().
			WithFields("attacker", proto.Attacker.Name).
			WithFields("defenser", proto.Attacker.Name).
			WithFields("assister", proto.Attacker.Name).JsonString()
		mailProto.Report = proto

		//data := datas.MailData().Array[rand.Intn(len(datas.MailData().Array))]
		//mailProto = data.NewArgsTextMail(hero.Name(), hero.Name(), hero.Name())
		//mailProto.Title = "GM邮件: " + mailProto.Title
		//mailProto.Report = proto

		return

	})

	m.modules.MailModule().SendReportMail(hc.Id(), mailProto, m.time.CurrentTime())
}
