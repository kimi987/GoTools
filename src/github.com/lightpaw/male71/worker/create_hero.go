package worker

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/module/realm/realmface"
	service1 "github.com/lightpaw/male7/service"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math/rand"
	"strings"
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/util/timeutil"
)

func (m *MessageWorker) processCreateHero(data *service.MsgData) bool {
	if m.user == nil {
		logrus.Errorf("worker received create hero msg, but not login, received: %d-%d",
			data.ModuleID, data.SequenceID)
		m.Send(login.ERR_CREATE_HERO_FAIL_NO_LOGIN)
		m.Close()
		return false
	}

	if m.user.GetHeroController() != nil {
		// 这里是怎么进来的...
		logrus.Errorf("worker received create hero msg, hero exist, received: %d-%d",
			data.ModuleID, data.SequenceID)
		m.Send(login.ERR_CREATE_HERO_FAIL_CREATED)
		m.Close()
		return false
	}

	proto, ok := data.Proto.(*login.C2SCreateHeroProto)
	if !ok {
		logrus.Errorf("worker.processCreateHero proto marshal fail")
		m.Send(login.ERR_CREATE_HERO_FAIL_INVALID_PROTO)
		m.Close()
		return false
	}

	// 角色创建流程
	heroId := m.user.Id()
	logrus.WithField("hero_id", heroId).WithField("name", proto.Name).Debug("新建英雄")

	countryId := u64.FromInt32(proto.Country)
	if service.ConfigDatas.GetCountryData(countryId) == nil {
		if array := service.ConfigDatas.GetCountryDataArray(); len(array) > 0 {
			countryId = array[rand.Intn(len(array))].Id
			logrus.Debugf("创建角色，无效的国家(%d)，随机一个给玩家（%d）", proto.Country, countryId)
		} else {
			logrus.Errorf("创建角色，无效的国家(%d)，无法随机", proto.Country)
			m.user.Send(login.ERR_CREATE_HERO_FAIL_COUNTRY_ERR)
			return false
		}
	}

	ctime := service.TimeService.CurrentTime()
	heroName, errMsg := onCreateHero(heroId, countryId, proto, ctime)
	if errMsg != nil {
		m.user.Send(errMsg)
		return false
	}

	hc := service1.NewHeroController(heroId, m, m.clientIp, m.clientIp32, m.pf, service.HeroDataService.NewHeroLocker(heroId))
	m.user.SetHeroController(hc)

	hc.Send(login.CREATE_HERO_S2C)

	if service.RegionModule.InitHeroBase(hc, service.TimeService.CurrentTime(), countryId, realmface.AddBaseHomeNewHero) {
		logrus.Errorf("新建英雄，创建主城失败")
	}

	if tencentInfo := m.user.TencentInfo(); tencentInfo != nil {
		tlogHero := entity.NewSimpleTlogHeroInfo(hc.Id(), heroName)
		data := service.TlogService.BuildPlayerRegister(tlogHero, tencentInfo)
		service.TlogService.WriteLog(data)
	}

	sinceDuration := ctime.Sub(service.IndividualServerConfig.GetServerStartTime())
	sinceDay := timeutil.DivideTimes(sinceDuration, timeutil.Day)

	gamelogs.CreateHeroLog(constants.PID, m.user.Sid(), heroId, ctime.Unix(), sinceDay)
	gamelogs.UpgradeHeroLevelLog(constants.PID, m.user.Sid(), heroId, 1)

	return true
}

func onCreateHero(heroId int64, countryId uint64, proto *login.C2SCreateHeroProto, ctime time.Time) (string, pbutil.Buffer) {

	name := proto.Name

	if len(name) > 0 {
		name = util.ReplaceInvalidChar(name)
		name = truncateName(name)

		// 敏感词查询
		resp, err := service.TssClient.CheckName(name, false)
		if err != nil {
			logrus.WithField("name", name).WithError(err).Error("创建英雄，敏感词查询失败")
			goto DEFAULT_NAME
		}

		if resp.Ret != 0 {
			logrus.WithField("ret", resp.Ret).WithField("ret_msg", resp.RetMsg).WithField("name", name).Error("创建英雄，敏感词查询失败")
			goto DEFAULT_NAME
		}

		if resp.MsgResultFlag != 0 {
			logrus.WithField("ret", resp.MsgResultFlag).Error("创建英雄，名字包含敏感词")
			goto DEFAULT_NAME
		}

		suc, errMsg := tryCreateHero(heroId, name, countryId, proto, ctime)
		if errMsg != nil {
			goto DEFAULT_NAME
		} else if suc {
			return name, nil // 创建成功
		}

		// 已经存在了，随机处理
		for genRuneCount := 1; genRuneCount <= 3; genRuneCount++ {
			newName := genName(name, genRuneCount)

			suc, errMsg := tryCreateHero(heroId, newName, countryId, proto, ctime)
			if errMsg != nil {
				continue
			} else if suc {
				name = newName
				return newName, nil // 创建成功
			}
		}
	}

DEFAULT_NAME:

// 名字最长20个字符
	name = idbytes.PlayerName(heroId)
	//if c := util.GetCharLen(name); c <= 0 {
	//	logrus.Debugf("worker.processCreateHero 名字长度不符合条件, name: %v, len: %v", name, c)
	//	return login.ERR_CREATE_HERO_FAIL_INVALID_NAME
	//}

	suc, errMsg := tryCreateHero(heroId, name, countryId, proto, ctime)
	if errMsg != nil {
		return "", errMsg
	} else if suc {
		return name, nil // 创建成功
	}

	// 到这里还是没创建出来，我去
	return "", login.ERR_CREATE_HERO_FAIL_SERVER_ERROR
}

func tryCreateHero(heroId int64, name string, countryId uint64, proto *login.C2SCreateHeroProto, ctime time.Time) (suc bool, errMsg pbutil.Buffer) {
	// 检查名字是否存在

	exist := false
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		exist, err = service.DbService.HeroNameExist(ctx, name)
		return
	})
	if err != nil {
		// 已经存在了，创建
		logrus.WithError(err).Error("创建角色失败")
		errMsg = login.ERR_CREATE_HERO_FAIL_SERVER_ERROR
		return
	}

	if !exist {
		err := createHero(heroId, name, proto.Male, proto.HeadUrl, countryId, ctime)
		if err == nil {
			// 创建成功
			suc = true
			service.GameExporter.GetRegisterCounter().Inc()
			return
		}

		if err == lock.ErrCreateExist {
			// 英雄已经创建了
			logrus.WithError(err).Errorf("创建英雄的时候，发现英雄已经创建了??")
			errMsg = login.ERR_CREATE_HERO_FAIL_SERVER_ERROR
			return
		}

		//重复的名字
	}

	return
}

var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genName(name string, genRuneCount int) string {
	var appendStr string
	for i := 0; i < genRuneCount; i++ {
		appendStr += string(runes[rand.Intn(len(runes))]) // 随机一个
	}

	newName := name + appendStr
	if uint64(util.GetCharLen(newName)) <= service.ConfigDatas.MiscConfig().MaxNameCharLen {
		return newName
	}

	runeArray := []rune(name)

	// 干掉最后一个字符
	for i := 0; i < genRuneCount; i++ {
		if len(runeArray) <= 0 {
			break
		}

		runeArray = runeArray[:len(runeArray)-1]
		result := string(runeArray) + appendStr
		if uint64(util.GetCharLen(result)) <= service.ConfigDatas.MiscConfig().MaxNameCharLen {
			return result
		}
	}

	logrus.Errorf("随机不到")
	return name
}

// 裁减名字
func truncateName(name string) string {
	charLen := uint64(util.GetCharLen(name))

	if charLen < service.ConfigDatas.MiscConfig().MinNameCharLen {
		return genName(name, int(service.ConfigDatas.MiscConfig().MinNameCharLen-charLen))
	} else if charLen > service.ConfigDatas.MiscConfig().MaxNameCharLen {
		return util.TruncateCharLen(name, int(service.ConfigDatas.MiscConfig().MaxNameCharLen))
	}

	return name
}

func createHero(heroId int64, name string, male bool, headUrl string, countryId uint64, ctime time.Time) error {

	hero := entity.NewHero(heroId, name, service.ConfigDatas.HeroInitData(), ctime)
	hero.SetMale(male)
	hero.SetCountryId(countryId)

	if len(headUrl) > 0 && strings.HasPrefix(headUrl, "http") {
		hero.SetHeadUrl(headUrl)
	}

	// 初始化数据
	initHeroCreateData(hero, ctime)

	return service.HeroDataService.Create(hero)
}

func initHeroCreateData(hero *entity.Hero, ctime time.Time) {
	data := service.ConfigDatas.HeroCreateData()

	heroMilitary := hero.Military()

	heroMilitary.SetNewSoldierCount(data.NewSoldier, ctime)
	hero.SetDailyResetTime(ctime)
	hero.SetDailyZeroResetTime(ctime)
	hero.SetSeasonResetTime(ctime)

	heroMilitary.AddFreeSoldier(data.NewSoldier, ctime)

	hero.GetSafeResource().AddGold(data.Gold)
	hero.GetSafeResource().AddFood(data.Food)
	hero.GetSafeResource().AddWood(data.Wood)
	hero.GetSafeResource().AddStone(data.Stone)

	// 加初始仇恨
	for _, t := range service.ConfigDatas.GetRegionMultiLevelNpcTypeDataArray() {
		if t.InitHate > 0 {
			hero.GetOrCreateNpcTypeInfo(t.Type).SetHate(t.InitHate)
		}
	}

	// 初始入侵野怪
	for _, v := range service.ConfigDatas.GetHomeNpcBaseDataArray() {
		if v.HomeBaseLevel == 1 || v.BaYeStage == 1 {
			hero.CreateHomeNpcBase(v)
		}
	}

	// 加武将
	configDatas := service.ConfigDatas

	for _, data := range configDatas.HeroCreateData().Captain {

		captain := hero.NewCaptain(data, ctime)
		captain.SetTrainAccExp(data.GetInitTrainExp())
		captain.CalculateProperties()
		captain.AddSoldier(u64.Sub(captain.SoldierCapcity(), captain.Soldier()))

		hero.Military().AddCaptain(captain)

		hero.WalkPveTroop(func(troop *entity.PveTroop) (endWalk bool) {
			troop.AddCaptain(captain)
			return
		})

		troop, index := hero.GetRecruitCaptainTroop()
		if troop != nil {
			troop.SetCaptainIfAbsent(index, captain, 0)
		}

	}

	// 修炼馆开始时间
	hero.Military().SetGlobalTrainStartTime(ctime)

	for _, t := range hero.Troops() {
		t.UpdateFightAmountIfChanged()
	}

	// 轩辕会武初始积分
	hero.Xuanyuan().AddScore(service.ConfigDatas.XuanyuanMiscData().InitScore)

	if service.IndividualServerConfig.GetSkipHeader() {
		for k := range shared_proto.HeroBoolType_name {
			bt := shared_proto.HeroBoolType(k)
			switch bt {
			case shared_proto.HeroBoolType_BOOL_XUAN_YUAN:
			default:
				hero.Bools().SetTrue(bt)
			}
		}
	}

	// 功能开启
	for _, data := range hero.Function().FunctionOpenDataArray {
		if heromodule.GetIsFuncOpened(hero, data) {
			hero.Function().OpenFunction(data.FunctionType)
		}
	}

	// 任务
	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumLoginDays)

	taskTypes := []shared_proto.TaskTargetType{
		shared_proto.TaskTargetType_TASK_TARGET_ACCUM_LOGIN_DAY,
	}

	if len(taskTypes) > 0 {
		hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
			for _, t := range taskTypes {
				task.Progress().UpdateTaskTypeProgress(t, hero)
			}
			return false
		})
	}

	// 初始不自动补兵
	hero.MiscData().SetDefenserDontAutoFullSoldier(true)

	// 初始化随机事件
	heromodule.ResetRandomEvent(hero, ctime, configDatas)
}
