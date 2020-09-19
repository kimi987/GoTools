package misc

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/rpcpb/game2login"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"golang.org/x/net/context"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/config/settings"
	"github.com/lightpaw/male7/config/charge"
	"github.com/lightpaw/male7/util"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/compress"
)

func NewMiscModule(dep iface.ServiceDep, datas iface.ConfigDatas, time iface.TimeService,
	individualServerConfig *kv.IndividualServerConfig,
	cluster iface.ClusterService, recommendHeroCache iface.LocationHeroCache) *MiscModule {
	m := &MiscModule{
		blockMsg:               misc.NewS2cBlockMsg(datas.BlockData().MinKeyData.GetProtoBytes()).Static(),
		dep:                    dep,
		time:                   time,
		cluster:                cluster,
		recommendHeroCache:     recommendHeroCache,
		individualServerConfig: individualServerConfig,
		datas:                  datas,
	}

	configProto := datas.EncodeClient()
	m.confVersion, m.newConfMsg, m.sameConfMsg = buildNewConfigMsg(must.Marshal(configProto))

	var err error
	m.luaConfVersion, m.newLuaConfMsg, m.sameLuaConfMsg, err = buildNewLuaConfigMsg(individualServerConfig.LuaConfAddr, configProto)
	if err != nil {
		logrus.WithError(err).Panic("初始化LuaConfig配置文件消息失败")
	}

	heromodule.RegisterHeroOnlineListener(m)

	return m
}

//gogen:iface
type MiscModule struct {
	confVersion    string
	newConfMsg     pbutil.Buffer
	sameConfMsg    pbutil.Buffer
	blockMsg       pbutil.Buffer
	luaConfVersion string
	newLuaConfMsg  pbutil.Buffer
	sameLuaConfMsg pbutil.Buffer

	dep     iface.ServiceDep
	time    iface.TimeService
	cluster iface.ClusterService

	recommendHeroCache iface.LocationHeroCache

	datas iface.ConfigDatas

	individualServerConfig *kv.IndividualServerConfig
}

func (m *MiscModule) OnHeroOnline(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heromodule.CheckFuncsOpened(hero, result)

		// 检查外城解锁
		hero.Domestic().OuterCities().SetAutoUnlock()

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_heart_beat
func (m *MiscModule) ProcessHeartBeat(hc iface.HeroController) {
	// 更新最后一次时间消息时间
}

//gogen:iface c2s_ping
func (m *MiscModule) ProcessPing(hc iface.HeroController) {
	hc.Send(misc.PING_S2C)
}

//gogen:iface c2s_config
func (m *MiscModule) ProcessRequestConfig(proto *misc.C2SConfigProto, hc iface.HeroController) {
	m.SendConfig(proto, hc)
}

func (m *MiscModule) SendConfig(proto *misc.C2SConfigProto, sender sender.Sender) {
	if proto.Version == m.confVersion {
		sender.Send(m.sameConfMsg)
	} else {
		sender.Send(m.newConfMsg)
	}
}

//gogen:iface c2s_configlua
func (m *MiscModule) ProcessRequestLuaConfig(proto *misc.C2SConfigluaProto, hc iface.HeroController) {
	m.SendLuaConfig(proto, hc)
}

func (m *MiscModule) SendLuaConfig(proto *misc.C2SConfigluaProto, sender sender.Sender) {
	if proto.Version == m.luaConfVersion {
		sender.Send(m.sameLuaConfMsg)
	} else {
		sender.Send(m.newLuaConfMsg)
	}
}

//gogen:iface
func (m *MiscModule) ProcessClientLog(proto *misc.C2SClientLogProto, hc iface.HeroController) {
	var name string
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		name = hero.Name()
		result.Ok()
		return
	})

	m.PrintClientLog(hc.Id(), name, proto.Level, proto.Text)
}

func (m *MiscModule) PrintClientLog(id int64, name, level, text string) {
	logrus.WithField("level", level).
		WithField("hero_id", id).
		WithField("hero_name", name).
		Errorf(text)
}

//gogen:iface
func (m *MiscModule) ProcessSyncTime(proto *misc.C2SSyncTimeProto, hc iface.HeroController) {
	m.SyncTime(proto.ClientTime, hc)
}

func (m *MiscModule) SyncTime(clientTime int32, hc sender.Sender) {
	hc.Send(misc.NewS2cSyncTimeMsg(clientTime, timeutil.Marshal32(m.time.CurrentTime())))
}

//gogen:iface c2s_block
func (m *MiscModule) ProcessGetBlock(hc iface.HeroController) {
	hc.Send(m.blockMsg)
}

//gogen:iface c2s_client_version
func (m *MiscModule) ProcessClientVersion(proto *misc.C2SClientVersionProto, hc iface.HeroController) {
	m.SendClientVersion(hc)
}

func (m *MiscModule) SendClientVersion(sender sender.Sender) {
	//if msg, ok := m.clientVersionRef.Load().(pbutil.Buffer); ok {
	//	sender.Send(msg)
	//}
}

func (m *MiscModule) setClientVersion(toSet string) {
	//m.clientVersionRef.Store(misc.NewS2cClientVersionMsg(toSet).Static())
}

func (m *MiscModule) GmSetClientVersion(toSet string) {
	m.setClientVersion(toSet)
}

//gogen:iface
func (m *MiscModule) ProcessUpdatePfToken(proto *misc.C2SUpdatePfTokenProto, hc iface.HeroController) {
	hc.Send(misc.UPDATE_PF_TOKEN_S2C)

	if len(proto.Token) <= 0 {
		logrus.Debug("客户端发送更新平台Token，但是token为空")
		return
	}

	ctxfunc.NetTimeout3s(func(ctx context.Context) (err error) {
		resp, err := game2login.VerifyLoginToken(m.cluster.LoginClient(), ctx, hc.Id(), proto.Token, hc.GetClientIp(), hc.GetPf())
		if err != nil {
			logrus.WithError(err).Error("misc.UpdatePfToken() verify login token error, %d", hc.Id())
			return err
		}

		logrus.WithField("success", resp.Success).Debug("更新玩家平台token")
		return nil
	})

}

//gogen:iface c2s_settings
func (m *MiscModule) ProcessSettings(proto *misc.C2SSettingsProto, hc iface.HeroController) {
	settingType := shared_proto.SettingType(proto.SettingType)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroSettings := hero.Settings()
		if !heroSettings.IsValidType(settingType) {
			logrus.Debugf("非法的设置类型: %v", proto.SettingType)
			result.Add(misc.ERR_SETTINGS_FAIL_INVALID_TYPE)
			return
		}

		result.Changed()
		result.Ok()

		heroSettings.Set(settingType, proto.Open)

		result.Add(misc.NewS2cSettingsMsg(settingType, proto.Open))
	})
}

//gogen:iface c2s_settings_to_default
func (m *MiscModule) ProcessSettingsToDefault(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroSettings := hero.Settings()

		heroSettings.ResetDefault(m.datas.SettingMiscData().DefaultSettings)

		result.Add(misc.NewS2cSettingsToDefaultMsg(heroSettings.Encode()))

		result.Changed()
		result.Ok()
	})
}

func buildNewConfigMsg(dataBytes []byte) (confVersion string, diffVersion, sameVersion pbutil.Buffer) {
	confVersion = util.Md5String(dataBytes)
	diffVersion = misc.NewS2cConfigMsg(confVersion, dataBytes).Static()
	sameVersion = misc.NewS2cConfigMarshalMsg(confVersion, &shared_proto.Config{}).Static()
	return
}

//gogen:iface
func (m *MiscModule) ProcessUpdateLocation(proto *misc.C2SUpdateLocationProto, hc iface.HeroController) {
	if proto.Location < 0 {
		hc.Send(misc.ERR_UPDATE_LOCATION_FAIL_LOCATION_ERROR)
		return
	}

	toSet := u64.FromInt32(proto.Location)
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.SetLocation(toSet)
		result.Add(misc.NewS2cUpdateLocationMsg(proto.Location))

		result.Changed()
		result.Ok()
	}) {
		return
	}

	m.recommendHeroCache.UpdateLocation(hc.Id(), toSet)
}

//gogen:iface c2s_background_heart_beat
func (m *MiscModule) ProcessBackgroudHeartBeat(hc iface.HeroController) {
	// 推入后台，30分钟之内都会收到推送
	hc.SetIsInBackgroud(m.time.CurrentTime().Add(30*time.Minute), true)
}

//gogen:iface c2s_background_weakup
func (m *MiscModule) ProcessBackgroudWeakup(hc iface.HeroController) {
	// 从后台唤醒，不再收到推送
	hc.SetIsInBackgroud(time.Time{}, false)
}

////gogen:iface c2s_open_charge_prize_ui
//func (m *MiscModule) ProcessOpenChargePrizeUi(hc iface.HeroController) {
//	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//		result.Add(misc.NewS2cOpenChargePrizeUiMsg(hero.Misc().EncodeCharge()))
//		result.Ok()
//	})
//}

//gogen:iface
func (m *MiscModule) ProcessCollectChargePrize(proto *misc.C2SCollectChargePrizeProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetChargePrizeData(id)
	if data == nil {
		logrus.Debugf("领取充值奖励，无效id")
		hc.Send(misc.ERR_COLLECT_CHARGE_PRIZE_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Misc().ChargeAmount() < data.Amount {
			logrus.Debugf("领取充值奖励，充值不足")
			result.Add(misc.ERR_COLLECT_CHARGE_PRIZE_FAIL_NOT_ENOUGH_CHARGE_AMOUNT)
			return
		}
		if !hero.Misc().TryCollectChargePrize(data.Id) {
			logrus.Debugf("领取充值奖励，已经领取")
			result.Add(misc.ERR_COLLECT_CHARGE_PRIZE_FAIL_COLLECTED)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.MiscCollectChargePrize)
		heromodule.AddPrize(hctx, hero, result, data.Prize, m.time.CurrentTime())
		result.Add(misc.NewS2cCollectChargePrizeMsg(proto.Id, must.Marshal(data.Prize.Encode())))

		result.Ok()
	})
}

//gogen:iface
func (m *MiscModule) ProcessCollectDailyBargain(proto *misc.C2SCollectDailyBargainProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetDailyBargainData(id)
	if data == nil {
		logrus.Debugf("领取每日特惠，无效id")
		hc.Send(misc.ERR_COLLECT_DAILY_BARGAIN_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Misc().GetBoughtBargainTimes(data.Id) >= data.Limit {
			logrus.Debugf("领取每日特惠，无法领取")
			result.Add(misc.ERR_COLLECT_DAILY_BARGAIN_FAIL_CANNOT_COLLECT)
			return
		}
		hero.Misc().IncreaseBoughtBargainTimes(data.Id)

		ctime := m.time.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.MiscCollectDailyBargain)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)
		result.Add(misc.NewS2cCollectDailyBargainMsg(proto.Id, u64.Int32(hero.Misc().GetBoughtBargainTimes(data.Id)), must.Marshal(data.Prize.Encode())))

		// 发邮件通告
		mailData := m.datas.MailHelp().BuyDaillyBargainSuccess
		proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = mailData.Text.New().WithName(data.Name).JsonString()
		m.dep.Mail().SendProtoMail(hero.Id(), proto, ctime)

		result.Ok()
	})
}

//gogen:iface
func (m *MiscModule) ProcessActivateDurationCard(proto *misc.C2SActivateDurationCardProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetDurationCardData(id)
	if data == nil {
		logrus.Debugf("购买（激活）尊享卡，无效id")
		hc.Send(misc.ERR_ACTIVATE_DURATION_CARD_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if data.Duration <= 0 && hero.Misc().IsDurationCardActive(data.Id) { // 永久卡
			logrus.Debugf("购买（激活）尊享卡，永久卡无法重复激活")
			result.Add(misc.ERR_ACTIVATE_DURATION_CARD_FAIL_CANNOT_ACTIVATE)
			return
		}
		//TODO: 等待充值反馈
		charged := true
		if !charged { // 充值未响应
			logrus.Debugf("购买（激活）尊享卡，充值未响应")
			result.Add(misc.ERR_ACTIVATE_DURATION_CARD_FAIL_NO_CHARGE)
			return
		}
		ctime := m.time.CurrentTime()
		// 激活（续费）卡片
		time := hero.Misc().ActivateDurationCard(data, ctime)
		// 给奖励
		hctx := heromodule.NewContext(m.dep, operate_type.MiscActivateDurationCard)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)
		result.Add(misc.NewS2cActivateDurationCardMsg(proto.Id, timeutil.Marshal32(time), must.Marshal(data.Prize.Encode())))

		result.Ok()
	})
}

//gogen:iface
func (m *MiscModule) ProcessCollectDurationCardDailyPrize(proto *misc.C2SCollectDurationCardDailyPrizeProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetDurationCardData(id)
	if data == nil {
		logrus.Debugf("领取尊享卡每日奖励，无效id")
		hc.Send(misc.ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.Misc().IsDurationCardActive(data.Id) {
			logrus.Debugf("领取尊享卡每日奖励，卡未激活")
			result.Add(misc.ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_NOT_ACTIVE)
			return
		}
		ctime := m.time.CurrentTime()
		if hero.Misc().IsDurationCardTimeEnd(data, ctime) {
			logrus.Debugf("领取尊享卡每日奖励，卡片已过期")
			result.Add(misc.ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_OVERDUE)
			return
		}
		if hero.Misc().IsCollectDurationCardDailyPrize(data.Id) {
			logrus.Debugf("领取尊享卡每日奖励，今日已领取")
			result.Add(misc.ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_COLLECTED)
			return
		}
		hero.Misc().SetCollectedDurationCardDailyPrize(data.Id)
		hctx := heromodule.NewContext(m.dep, operate_type.MiscCollectDurationCardDailyPrize)
		heromodule.AddPrize(hctx, hero, result, data.DailyPrize, ctime)
		result.Add(misc.NewS2cCollectDurationCardDailyPrizeMsg(proto.Id, must.Marshal(data.DailyPrize.Encode())))

		result.Ok()
	})
}

////gogen:iface
//func (m *MiscModule) ProcessActivateChargeObj(proto *misc.C2SActivateChargeObjProto, hc iface.HeroController) {
//	id := u64.FromInt32(proto.GetId())
//	data := m.datas.GetChargeObjData(id)
//	if data == nil {
//		logrus.Debugf("购买（激活）充值项，无效id")
//		hc.Send(misc.ERR_ACTIVATE_CHARGE_OBJ_FAIL_INVALID_ID)
//		return
//	}
//	orderNumber := proto.GetOrderNumber()
//	err := m.chargeValidation(orderNumber)
//	if err != nil {
//		switch err.Error() {
//		case "err_order_number":
//			logrus.Debugf("购买（激活）充值项，错误的订单号")
//			hc.Send(misc.ERR_ACTIVATE_CHARGE_OBJ_FAIL_ERR_ORDER_NUMBER)
//		case "invalid_order_number":
//			logrus.Debugf("购买（激活）充值项，无效的订单号")
//			hc.Send(misc.ERR_ACTIVATE_CHARGE_OBJ_FAIL_INVALID_ORDER_NUMBER)
//		default:
//			logrus.Debugf("购买（激活）充值项，验证超时")
//			hc.Send(misc.ERR_ACTIVATE_CHARGE_OBJ_FAIL_VALIDATION_TIMEOUT)
//		}
//		return
//	}
//	ctime := m.time.CurrentTime()
//	hctx := heromodule.NewContext(m.dep, operate_type.MiscActivateChargeObj)
//	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//		firstCharge := hero.Misc().IsFirstChargedObj(id)
//		if firstCharge {
//			hero.Misc().SetFirstChargedObj(id)
//		}
//		// 记录累计充值
//		hero.Misc().AddChargeAmount(data.ChargeAmount)
//		result.Add(misc.NewS2cUpdateChargeAmountMsg(u64.Int32(hero.Misc().ChargeAmount())))
//		// 获取元宝
//		rewardYuanbao := data.GetRewardYuanbao(firstCharge)
//		heromodule.AddYuanbao(hctx,hero, result, rewardYuanbao)
//		// 增加vip经验
//		upgrade := heromodule.AddVipExp(hctx, hero, result, data.VipExp, ctime)
//		if upgrade {
//			result.Changed()
//		}
//		result.Add(misc.NewS2cActivateChargeObjMsg(proto.Id, u64.Int32(hero.Vip().Level()), u64.Int32(hero.Vip().Exp())))
//		result.Ok()
//	})
//}
//
//func (m *MiscModule) chargeValidation(orderNumber string) error {
//	// TODO: 充值订单号验证
//	return nil
//}

// 隐私设置
//gogen:iface
func (m *MiscModule) ProcessSetPrivacySetting(proto *misc.C2SSetPrivacySettingProto, hc iface.HeroController) {
	data := m.datas.GetPrivacySettingData(u64.FromInt32(proto.GetSettingId()))
	if data == nil {
		logrus.Debugf("隐私设置，无效的类型")
		hc.Send(misc.ERR_SET_PRIVACY_SETTING_FAIL_INVALID_TYPE)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.Settings().TrySetPrivacySetting(data.SettingType, proto.GetOpenOrClose()) {
			logrus.Debugf("隐私设置，重复设置")
			result.Add(misc.ERR_SET_PRIVACY_SETTING_FAIL_DUPLICATION)
			return
		}
		result.Add(misc.NewS2cSetPrivacySettingMsg(data.SettingType, proto.GetOpenOrClose()))
		result.Ok()
	})
}

// 恢复默认隐私设置
//gogen:iface c2s_set_default_privacy_settings
func (m *MiscModule) ProcessSetDefaultPrivacySettings(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		changed, isOpen := hero.Settings().SetPrivacySettings(settings.DefaultPrivacySettings)
		if len(changed) <= 0 {
			logrus.Debugf("恢复默认隐私设置，已经恢复默认设置")
			result.Add(misc.ERR_SET_DEFAULT_PRIVACY_SETTINGS_FAIL_HAS_DEFAULT)
			return
		}
		result.Changed()
		for i, t := range changed {
			result.Add(misc.NewS2cSetPrivacySettingMsg(t, isOpen[i]))
		}
		result.Ok()
	})
}

//gogen:iface
func (m *MiscModule) ProcessGetProductInfo(proto *misc.C2SGetProductInfoProto, hc iface.HeroController) {

	data := m.datas.GetProductData(u64.FromInt32(proto.Id))
	if data == nil {
		logrus.Debug("获取购买商品信息，data == nil")
		hc.Send(misc.ERR_GET_PRODUCT_INFO_FAIL_INVALID_ID)
		return
	}

	switch v := data.GetData().(type) {
	case *charge.ChargeObjData:

		ctime := m.time.CurrentTime()
		sid := uint32(m.individualServerConfig.ServerID)
		cpOrderId := util.NewCpOrderId(sid, hc.Id(), data.Id, data.Price, ctime.Unix())
		ext := util.NewCpOrderSign(cpOrderId, m.individualServerConfig.OrderKey)
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			isFirstRecharge := hero.Misc().IsFirstChargedObj(v.Id)
			result.Add(misc.NewS2cGetProductInfoMsg(
				u64.Int32(data.Id), data.ProductId, data.ProductName,
				cpOrderId, u64.Int32(data.Price), u64.Int32(v.GetRewardYuanbao(isFirstRecharge)),
				ext, m.individualServerConfig.P8RechargeDebug))

			result.Ok()
		})

	default:
		logrus.WithField("id", data.Id).Error("获取购买商品信息，该类型没有switch")
		hc.Send(misc.ERR_GET_PRODUCT_INFO_FAIL_CANT_BUY)
	}
}

func buildNewLuaConfigMsg(addr string, proto *shared_proto.Config) (confVersion string, diffVersion, sameVersion pbutil.Buffer, err error) {
	version, luaBytes, err := util.Proto2LuaBytes(addr, proto)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "生成LuaConfig消息出错")
	}

	uncompressBytes, err := compress.GzipUncompress(luaBytes)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "gzip解压LuaConfig失败")
	}
	logrus.Infof("获取客户端配置成功， version: %s length: %d uncompress: %d", version, len(luaBytes), len(uncompressBytes))

	return version, misc.NewS2cConfigluaMsg(version, uncompressBytes).Static(), misc.NewS2cConfigluaMsg(version, nil).Static(), nil
}
