package secret_tower

import (
	"bytes"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/service/operate_type"
)

var (
	queueTimeoutWheel = timer.NewTimingWheel(500*time.Millisecond, 8) // 最大timeout是4秒
)

const (
	queueTimeout = 3 * time.Second
)

// 重楼密室
func NewSecretTowerModule(
	dep iface.ServiceDep,
	seasonService iface.SeasonService,
	fightService iface.FightXService) *SecretTowerModule {

	m := &SecretTowerModule{}

	m.dep = dep
	m.timeService = dep.Time()
	m.seasonService = seasonService
	m.configDatas = dep.Datas()
	m.fightService = fightService
	m.heroDataService = dep.HeroData()
	m.worldService = dep.World()
	m.guildService = dep.Guild()
	m.heroSnapshotService = dep.HeroSnapshot()

	m.loopExitNotify = make(chan struct{})
	m.closeNotify = make(chan struct{})

	m.eventChan = make(chan *event, 1024)

	m.manager = newTeamManager(dep.Datas(), dep.World(), dep.HeroData())

	heromodule.RegisterHeroOnlineListener(m)
	heromodule.RegisterHeroOfflineListener(m)
	heromodule.RegisterHeroPveTroopChangeEventHandlers(m.onPveTroopChanged)

	go call.CatchLoopPanic(m.loop, "secret_tower")

	return m
}

//gogen:iface
type SecretTowerModule struct {
	dep                 iface.ServiceDep
	timeService         iface.TimeService
	seasonService       iface.SeasonService
	configDatas         iface.ConfigDatas
	fightService        iface.FightXService
	heroDataService     iface.HeroDataService
	worldService        iface.WorldService
	guildService        iface.GuildService
	tickerService       iface.TickerService
	heroSnapshotService iface.HeroSnapshotService

	loopExitNotify chan struct{}
	closeNotify    chan struct{}
	closeOnce      sync.Once

	manager *team_manager

	eventChan chan *event
}

func (m *SecretTowerModule) miscData() *towerdata.SecretTowerMiscData {
	return m.configDatas.SecretTowerMiscData()
}

type event struct {
	f      func()
	called chan struct{}
}

func (m *SecretTowerModule) Close() {
	m.closeOnce.Do(func() {
		close(m.closeNotify)
	})
	<-m.loopExitNotify
}

func (m *SecretTowerModule) loop() {
	updateTick := time.NewTicker(8 * time.Second)

	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Error("BaiZhanModule.loop recovered from panic")
			metrics.IncPanic()
		}
	}()

	defer close(m.loopExitNotify)
	defer updateTick.Stop()

	for {
		select {
		case e := <-m.eventChan:
			e.f()
			close(e.called) // notify caller
		case <-updateTick.C:
			m.manager.update(m.timeService.CurrentTime())
		case <-m.closeNotify:
			return // quit loop
		}
	}
}

func (m *SecretTowerModule) queueWaitFunc(f func()) (funcCalled bool) {
	return m.queueFunc(true, f)
}

func (m *SecretTowerModule) queueFunc(waited bool, f func()) (funcCalled bool) {
	e := &event{f: f, called: make(chan struct{})}

	select {
	case m.eventChan <- e:
		if waited {
			select {
			case <-m.loopExitNotify:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true
		}

	case <-queueTimeoutWheel.After(queueTimeout):
		return false

	case <-m.closeNotify:
		return false
	}
}

func (m *SecretTowerModule) OnHeroOffline(hc iface.HeroController) {
	m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			return
		}

		destroy, notFound := team.Leave(hc.Id())
		if notFound {
			logrus.Errorf("离开重楼密室队伍，没在队伍中，但是此前判断有队伍，且在该队伍中!")
			return
		}

		if destroy {
			// 队伍都摧毁了，没必要广播了
		} else {
			// 广播XXX离开队伍了，新的队长是谁啦
			team.broadcast(secret_tower.NewS2cOtherLeaveLeaveTeamMsg(hc.IdBytes(), team.leader.IdBytes()))
		}
	})
}

func (m *SecretTowerModule) maxTimes() uint64 {
	return m.miscData().MaxTimes + m.seasonService.Season().SecretTowerTimes
}

func (m *SecretTowerModule) OnHeroOnline(hc iface.HeroController) {
	count, haveNew := m.manager.invites.InviteMeTeamAndHaveNew(hc.Id())

	if count <= 0 {
		return
	}

	hc.Send(secret_tower.NewS2cReceiveInviteMsg(count, haveNew))
}

func (m *SecretTowerModule) onPveTroopChanged(id int64, troopType shared_proto.PveTroopType) {
	if troopType != shared_proto.PveTroopType_DUNGEON {
		return
	}

	// 没有加入队伍
	if !m.manager.isHeroJoinTeam(id) {
		return
	}

	m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(id)
		if team == nil {
			return
		}

		_, mem := team.GetMember(id)
		if mem == nil {
			return
		}

		m.heroDataService.FuncNotError(id, func(hero *entity.Hero) (heroChanged bool) {
			mem.syncCaptains(hero)
			return
		})

		team.changeProtectEndTime(m.timeService.CurrentTime())
		team.invalidCache()
		team.broadcast(secret_tower.NewS2cMemberTroopChangedMsg(mem.EncodeClient(), timeutil.Marshal32(team.ProtectEndTime())))
	})
}

func (m *SecretTowerModule) requestTeamList(secretTowerId uint64, hc iface.HeroController) (processed bool, err msg.ErrMsg) {
	data := m.configDatas.GetSecretTowerData(secretTowerId)
	if data == nil {
		err = secret_tower.ErrRequestTeamListFailUnknownSecretTower
		return
	}

	hasOpen := false
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		hasOpen = hero.SecretTower().HasOpen(secretTowerId)
		return false
	})

	if !hasOpen {
		err = secret_tower.ErrRequestTeamListFailNotOpen
		return
	}

	processed = m.queueWaitFunc(func() {
		hc.Send(m.manager.mustTowerTeams(data).getTeamListCache())
	})

	return
}

func (m *SecretTowerModule) createTeam(secretTowerId uint64, isGuild bool, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	data := m.configDatas.GetSecretTowerData(secretTowerId)
	if data == nil {
		err = secret_tower.ErrCreateTeamFailUnknownTowerId
		return
	}

	processed = m.queueWaitFunc(func() {
		ctime := m.timeService.CurrentTime()

		if m.manager.isHeroJoinTeam(hc.Id()) {
			logrus.Debugf("创建重楼密室队伍，已经有队伍了")
			err = secret_tower.ErrCreateTeamFailHaveTeamNow
			return
		}

		guildId := int64(0)

		var member *secret_tower_team_member

		if m.heroDataService.FuncNotError(hc.Id(), func(hero *entity.Hero) (heroChanged bool) {
			heroSecretTower := hero.SecretTower()
			if !heroSecretTower.HasOpen(secretTowerId) {
				logrus.Debugf("创建重楼密室队伍，密室没有开启")
				err = secret_tower.ErrCreateTeamFailUnopen
				return
			}

			if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
				logrus.Debugf("创建重楼密室队伍，没有次数了")
				err = secret_tower.ErrCreateTeamFailTimesNotEnough
				return
			}

			if isGuild {
				guildId = hero.GuildId()
				if guildId <= 0 {
					logrus.Debugf("创建重楼密室队伍，没有联盟")
					err = secret_tower.ErrCreateTeamFailNoGuild
					return
				}
			}

			pveTroop := hero.PveTroop(shared_proto.PveTroopType_DUNGEON)
			if pveTroop == nil {
				logrus.Errorf("创建重楼密室队伍，队伍没找到")
				err = secret_tower.ErrCreateTeamFailServerError
				return
			}

			failType := hero.CheckCanGenCombatPlayer(shared_proto.PveTroopType_DUNGEON)
			switch failType {
			case entity.SUCCESS:
				break
			case entity.SERVER_ERROR:
				logrus.Errorf("创建重楼密室队伍，队伍没找到")
				err = secret_tower.ErrCreateTeamFailServerError
				return
			case entity.CAPTAIN_COUNT_NOT_ENOUGH:
				logrus.Debugf("创建重楼密室队伍，武将个数未满")
				err = secret_tower.ErrCreateTeamFailCaptainNotFull
				return
			default:
				logrus.Errorf("创建重楼密室队伍，未处理的错误类型 %v", failType)
				err = secret_tower.ErrCreateTeamFailServerError
				return
			}

			member = newSecretTowerTeamMember(hero, m.worldService.Send, shared_proto.TowerTeamMode_CHALLENGE, m.heroSnapshotService)

			return
		}) {
			// 有错误
			if err == nil {
				err = secret_tower.ErrCreateTeamFailServerError
			}
			return
		}

		if err != nil {
			// 里面处理了
			return
		}

		if member == nil {
			logrus.Errorf("创建重楼密室队伍，member为空")
			err = secret_tower.ErrCreateTeamFailServerError
			return
		}

		team := newSecretTowerTeam(m.manager, data, member, guildId, ctime)

		m.manager.addTeam(team)
		m.manager.onHeroJoinTeam(hc.Id(), team)

		// 发送消息，创建队伍成功了
		hc.Send(secret_tower.NewS2cCreateTeamMsg(team.EncodeDetail()))

		ok = true
	})

	return
}

func (m *SecretTowerModule) joinTeam(teamId int64, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	processed = m.queueWaitFunc(func() {
		ctime := m.timeService.CurrentTime()

		oldTeam := m.manager.getHeroJoinTeam(hc.Id())
		if oldTeam != nil {
			logrus.Debug("加入重楼密室队伍，玩家已经在队伍中")
			err = secret_tower.ErrJoinTeamFailHaveTeamNow
			return
		}

		team := m.manager.getTeam(teamId)
		if team == nil {
			logrus.Debugf("加入重楼密室队伍，队伍没找到")
			err = secret_tower.ErrJoinTeamFailTeamNotFound
			return
		}

		if team.IsFull() {
			logrus.Debugf("加入重楼密室队伍，队伍已满")
			err = secret_tower.ErrJoinTeamFailTeamFull
			return
		}

		_, existMember := team.GetMember(hc.Id())
		if existMember != nil {
			logrus.Errorf("加入重楼密室队伍，玩家在前面的检查中没有密室队伍，但是又在这个密室中，我去，什么鬼？")
			err = secret_tower.ErrJoinTeamFailHaveTeamNow
			return
		}

		var member *secret_tower_team_member

		if m.heroDataService.FuncNotError(hc.Id(), func(hero *entity.Hero) (heroChanged bool) {
			heroSecretTower := hero.SecretTower()
			if !heroSecretTower.HasOpen(team.towerData.Id) {
				logrus.Debugf("加入重楼密室队伍，密室没有开启")
				err = secret_tower.ErrJoinTeamFailUnopen
				return
			}

			if team.guildId > 0 && hero.GuildId() != team.guildId {
				logrus.Debugf("加入重楼密室队伍，没有不是目标联盟的")
				err = secret_tower.ErrJoinTeamFailNotTargetGuild
				return
			}

			challengeMode := shared_proto.TowerTeamMode_CHALLENGE
			if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
				if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
					logrus.Debugf("加入重楼密室队伍，没有次数了")
					err = secret_tower.ErrJoinTeamFailTimesNotEnough
					return
				}

				//if team.towerData.Id == heroSecretTower.MaxOpenSecretTowerId() {
				//	logrus.Debugf("加入重楼密室队伍，不可以协助最高自己能够参加的最高层数的密室")
				//	err = secret_tower.ErrJoinTeamFailCanNotHelpMaxTower
				//	return
				//}

				challengeMode = shared_proto.TowerTeamMode_HELP

				// 没有挑战次数了
				//if hero.GuildId() == 0 {
				//	logrus.Debugf("加入重楼密室队伍，没有联盟不可以协助模式加入队伍")
				//	err = secret_tower.ErrJoinTeamFailCanNotHelpNoGuildMemberTeam
				//	return
				//}
				//
				//existChallengeGuildMember := false
				//
				//for _, member := range team.Members() {
				//	if snapshot := member.heroSnapshotGetter(); snapshot != nil && snapshot.GuildId == hero.Guild() {
				//		existChallengeGuildMember = true
				//		break
				//	}
				//}
				//
				//if !existChallengeGuildMember {
				//	logrus.Debugf("加入重楼密室队伍，协助模式下没有盟友在队伍里面")
				//	err = secret_tower.ErrJoinTeamFailCanNotHelpNoGuildMemberTeam
				//	return
				//}
			}

			failType := hero.CheckCanGenCombatPlayer(shared_proto.PveTroopType_DUNGEON)
			switch failType {
			case entity.SUCCESS:
				break
			case entity.SERVER_ERROR:
				logrus.Errorf("加入重楼密室队伍，没找到队伍")
				err = secret_tower.ErrJoinTeamFailServerError
				return
			case entity.CAPTAIN_COUNT_NOT_ENOUGH:
				logrus.Debugf("加入重楼密室队伍，武将个数未满")
				err = secret_tower.ErrJoinTeamFailCaptainNotFull
				return
			default:
				logrus.Errorf("加入重楼密室队伍，未处理的错误类型")
				err = secret_tower.ErrJoinTeamFailServerError
				return
			}

			member = newSecretTowerTeamMember(hero, m.worldService.Send, challengeMode, m.heroSnapshotService)

			return
		}) {
			// 有错误
			if err == nil {
				err = secret_tower.ErrJoinTeamFailServerError
			}
			return
		}

		if err != nil {
			// 里面处理了
			return
		}

		if member == nil {
			logrus.Errorf("加入重楼密室队伍，member 为空")
			err = secret_tower.ErrCreateTeamFailServerError
			return
		}

		isExist, teamFull := team.Add(member, ctime)
		if teamFull {
			logrus.Debugf("加入重楼密室队伍，队伍已满")
			err = secret_tower.ErrJoinTeamFailTeamFull
			return
		}

		if isExist {
			logrus.Debugf("加入重楼密室队伍，已经在队伍中了")
			err = secret_tower.ErrJoinTeamFailHaveTeamNow
			return
		}

		//oldTeam := m.manager.getHeroJoinTeam(hc.Id())
		//if oldTeam != nil {
		//	destroy, notFound := oldTeam.Leave(hc.Id())
		//	if notFound {
		//		logrus.Errorf("离开重楼密室队伍，没在队伍中，但是此前判断有队伍，且在该队伍中!")
		//	} else {
		//		if destroy {
		//			// 队伍都摧毁了，没必要广播了
		//		} else {
		//			// 广播XXX离开队伍了，新的队长是谁啦
		//			oldTeam.broadcast(secret_tower.NewS2cOtherLeaveLeaveTeamMsg(hc.IdBytes(), oldTeam.leader.IdBytes()))
		//		}
		//	}
		//}

		m.manager.onHeroJoinTeam(hc.Id(), team)

		// 发送消息，加入重楼密室队伍成功了
		hc.Send(secret_tower.NewS2cJoinTeamMsg(team.EncodeDetail()))
		// 发送历史聊天记录
		team.SendChatRecord(member)

		// 广播xxx加入了
		team.broadcastIgnore(secret_tower.NewS2cOtherJoinJoinTeamMsg(member.EncodeClient(), timeutil.Marshal32(team.ProtectEndTime())), hc.Id())

		ok = true
	})

	return
}

func (m *SecretTowerModule) autoJoinTeam(dataId uint64, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	processed = m.queueWaitFunc(func() {
		ctime := m.timeService.CurrentTime()

		var data *towerdata.SecretTowerData
		if dataId != 0 {
			data = m.configDatas.GetSecretTowerData(dataId)
			if data == nil {
				logrus.Debugf("一键加入重楼密室队伍，配置ID没找到")
				err = secret_tower.ErrJoinTeamFailNotValidTeam
				return
			}
		}

		oldTeam := m.manager.getHeroJoinTeam(hc.Id())
		if oldTeam != nil {
			logrus.Errorf("一键加入重楼密室队伍，玩家在前面的检查中没有密室队伍，但是又在这个密室中，我去，什么鬼？")
			err = secret_tower.ErrJoinTeamFailHaveTeamNow
			return
		}

		var toJoinTeam *secret_tower_team
		var member *secret_tower_team_member
		if m.heroDataService.FuncNotError(hc.Id(), func(hero *entity.Hero) (heroChanged bool) {
			heroSecretTower := hero.SecretTower()

			if data != nil {
				if !heroSecretTower.HasOpen(data.Id) {
					logrus.Debugf("一键加入重楼密室队伍，密室没有开启")
					err = secret_tower.ErrJoinTeamFailUnopen
					return
				}
			}

			// 根据是否还有次数，确定模式
			teamMode := shared_proto.TowerTeamMode_CHALLENGE
			if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
				// 没有次数了，看下有没有协助次数
				if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
					logrus.Debugf("一键加入重楼密室队伍，没有次数了")
					err = secret_tower.ErrJoinTeamFailTimesNotEnough
					return
				}

				//if data != nil && data.Id == heroSecretTower.MaxOpenSecretTowerId() {
				//	logrus.Debugf("一键加入重楼密室队伍，不可以协助最高自己能够参加的最高层数的密室")
				//	err = secret_tower.ErrJoinTeamFailCanNotHelpMaxTower
				//	return
				//}

				// 没有挑战次数了
				if hero.GuildId() == 0 {
					logrus.Debugf("一键加入重楼密室队伍，没有联盟不可以协助模式加入队伍")
					err = secret_tower.ErrJoinTeamFailTeamNotFound // 没有合适的队伍
					return
				}

				teamMode = shared_proto.TowerTeamMode_HELP
			}

			var toJoinMemberCount uint64
			findTeam := func(team *secret_tower_team) (found bool) {
				if team.towerData != data {
					// 跟限定的对象不一致
					return false
				}

				if team.guildId > 0 && hero.GuildId() != team.guildId {
					return false
				}

				memberCount := uint64(team.MemberCount())
				if memberCount >= team.towerData.MaxAttackerCount {
					// full
					return false
				}

				//if teamMode == shared_proto.TowerTeamMode_HELP {
				//	existChallengeGuildMember := false
				//
				//	for _, member := range team.Members() {
				//		if snapshot := member.heroSnapshotGetter(); snapshot != nil && snapshot.GuildId == hero.Guild() {
				//			existChallengeGuildMember = true
				//			break
				//		}
				//	}
				//
				//	if !existChallengeGuildMember {
				//		// 协助模式必须有盟友存在
				//		return false
				//	}
				//}

				if toJoinTeam == nil || toJoinMemberCount < memberCount {
					// 进行比较，选择最优
					toJoinTeam = team
					toJoinMemberCount = memberCount

					if toJoinMemberCount+1 >= data.MaxAttackerCount {
						// 加我就满了，就他了
						return true
					}
				}

				return false
			}

			if data != nil {
				// 固定列表
				for _, team := range m.manager.mustTowerTeams(data).teams {
					found := findTeam(team)
					if found {
						break
					}
				}
			} else {
				// 邀请列表
				m.manager.invites.WalkInviteMeTeam(hero.Id(), func(teamId int64) (toContinue bool) {
					team := m.manager.getTeam(teamId)
					if team != nil {
						found := findTeam(team)
						if found {
							return false
						}
					}
					return true
				})
			}

			if toJoinTeam == nil {
				logrus.Debugf("一键加入重楼密室队伍，没有找到队伍可以加入")
				err = secret_tower.ErrJoinTeamFailTeamNotFound
				return
			}

			failType := hero.CheckCanGenCombatPlayer(shared_proto.PveTroopType_DUNGEON)
			switch failType {
			case entity.SUCCESS:
				break
			case entity.SERVER_ERROR:
				logrus.Errorf("一键加入重楼密室队伍，没找到队伍")
				err = secret_tower.ErrJoinTeamFailServerError
				return
			case entity.CAPTAIN_COUNT_NOT_ENOUGH:
				logrus.Debugf("一键加入重楼密室队伍，武将个数未满")
				err = secret_tower.ErrJoinTeamFailCaptainNotFull
				return
			default:
				logrus.Errorf("一键加入重楼密室队伍，未处理的错误类型")
				err = secret_tower.ErrJoinTeamFailServerError
				return
			}

			member = newSecretTowerTeamMember(hero, m.worldService.Send, teamMode, m.heroSnapshotService)

			return
		}) {
			// 有错误
			if err == nil {
				err = secret_tower.ErrJoinTeamFailServerError
			}
			return
		}

		if err != nil {
			// 里面处理了
			return
		}

		if toJoinTeam == nil {
			logrus.Errorf("一键加入重楼密室队伍，toJoinTeam 为空")
			err = secret_tower.ErrCreateTeamFailServerError
			return
		}

		if member == nil {
			logrus.Errorf("一键加入重楼密室队伍，member 为空")
			err = secret_tower.ErrCreateTeamFailServerError
			return
		}

		isExist, teamFull := toJoinTeam.Add(member, ctime)
		if teamFull {
			logrus.Debugf("一键加入重楼密室队伍，队伍已满")
			err = secret_tower.ErrJoinTeamFailTeamFull
			return
		}

		if isExist {
			logrus.Debugf("一键加入重楼密室队伍，已经在队伍中了")
			err = secret_tower.ErrJoinTeamFailHaveTeamNow
			return
		}

		m.manager.onHeroJoinTeam(hc.Id(), toJoinTeam)

		// 发送消息，一键加入重楼密室队伍成功了
		hc.Send(secret_tower.NewS2cJoinTeamMsg(toJoinTeam.EncodeDetail()))
		// 发送历史聊天记录
		toJoinTeam.SendChatRecord(member)

		// 广播xxx加入了
		toJoinTeam.broadcastIgnore(secret_tower.NewS2cOtherJoinJoinTeamMsg(member.EncodeClient(), timeutil.Marshal32(toJoinTeam.ProtectEndTime())), hc.Id())

		ok = true
	})

	return
}

func (m *SecretTowerModule) leaveTeam(hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("离开重楼密室队伍，没有队伍")
			err = secret_tower.ErrLeaveTeamFailNoTeam
			return
		}

		destroy, notFound := team.Leave(hc.Id())
		if notFound {
			logrus.Errorf("离开重楼密室队伍，没在队伍中，但是此前判断有队伍，且在该队伍中!")
			err = secret_tower.ErrLeaveTeamFailNoTeam
			return
		}

		hc.Send(secret_tower.LEAVE_TEAM_S2C)

		if destroy {
			// 队伍都摧毁了，没必要广播了
		} else {
			// 广播XXX离开队伍了，新的队长是谁啦
			team.broadcast(secret_tower.NewS2cOtherLeaveLeaveTeamMsg(hc.IdBytes(), team.leader.IdBytes()))
		}

		ok = true
	})
	return
}

func (m *SecretTowerModule) kickMember(id []byte, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	beenKickHeroId, e := idbytes.ToId(id)
	if !e {
		logrus.Debugf("踢出重楼密室队伍成员，解析客户端发送过来的id失败")
		err = secret_tower.ErrKickMemberFailTargetNotFound
		return
	}

	if beenKickHeroId == hc.Id() {
		logrus.Debugf("踢出重楼密室队伍成员，不能够踢出自己")
		err = secret_tower.ErrKickMemberFailCantKickSelf
		return
	}

	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("踢出重楼密室队伍成员，没有队伍")
			err = secret_tower.ErrKickMemberFailNoTeam
			return
		}

		if !team.IsLeader(hc.Id()) {
			logrus.Debugf("踢出重楼密室队伍成员，没有队伍")
			err = secret_tower.ErrKickMemberFailNotLeader
			return
		}

		beenKickMember := team.Kick(beenKickHeroId)
		if beenKickMember == nil {
			logrus.Debugf("踢出重楼密室队伍成员，没有找到目标")
			err = secret_tower.ErrKickMemberFailNotLeader
			return
		}

		// 广播xxx被踢出队伍了
		team.broadcast(secret_tower.NewS2cOtherBeenKickKickMemberMsg(beenKickMember.IdBytes()))

		// 发送给xxx你被踢出队伍了
		beenKickMember.SendMsg(secret_tower.KICK_MEMBER_S2C_YOU_BEEN_KICKED)

		ok = true
	})

	return
}

func (m *SecretTowerModule) moveMember(id []byte, up bool, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	beenMoveHeroId, e := idbytes.ToId(id)
	if !e {
		logrus.Debugf("移动重楼密室队伍成员，解析客户端发送过来的id失败")
		err = secret_tower.ErrMoveMemberFailTargetNotFound
		return
	}

	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("移动重楼密室队伍成员，自己没有队伍")
			err = secret_tower.ErrMoveMemberFailNoTeam
			return
		}

		if !team.IsLeader(hc.Id()) {
			logrus.Debugf("移动重楼密室队伍成员，不是队长")
			err = secret_tower.ErrMoveMemberFailNotLeader
			return
		}

		opSuccess, failAndIsFirst, notFound := team.Move(beenMoveHeroId, up)
		if notFound {
			logrus.Debugf("移动重楼密室队伍成员，目标没找到")
			err = secret_tower.ErrMoveMemberFailNotLeader
			return
		}

		if !opSuccess {
			logrus.Debugf("移动重楼密室队伍成员，操作失败")
			if failAndIsFirst {
				err = secret_tower.ErrMoveMemberFailTargetIsFirst
			} else {
				err = secret_tower.ErrMoveMemberFailTargetIsLast
			}
			return
		}

		// 操作成功了
		team.broadcast(secret_tower.NewS2cBroadcsatMoveMemberMsg(id, up))

		ok = true
	})

	return
}

func (m *SecretTowerModule) updateMemberPos(idArray [][]byte, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {

	ids, ok := idbytes.ToIds(idArray)
	if !ok {
		logrus.Debugf("密室更新位置，解析客户端发送过来的id失败")
		err = secret_tower.ErrUpdateMemberPosFailTargetNotFound
		return
	}

	if check.Int64Duplicate(ids) {
		logrus.Debugf("密室更新位置，发送的id存在重复")
		err = secret_tower.ErrUpdateMemberPosFailTargetDuplicate
		return
	}

	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("密室更新位置，自己没有队伍")
			err = secret_tower.ErrUpdateMemberPosFailNoTeam
			return
		}

		if !team.IsLeader(hc.Id()) {
			logrus.Debugf("密室更新位置，不是队长")
			err = secret_tower.ErrUpdateMemberPosFailNotLeader
			return
		}

		success := team.UpdateMemberPos(ids)
		if !success {
			logrus.WithField("ids", ids).Debugf("密室更新位置，Id列表无效")
			err = secret_tower.ErrUpdateMemberPosFailTargetNotFound
			return
		}

		// 操作成功了
		team.broadcast(secret_tower.NewS2cUpdateMemberPosMsg(idArray))

		ok = true
	})

	return
}

func (m *SecretTowerModule) changeMode(mode shared_proto.TowerTeamMode, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	if mode != shared_proto.TowerTeamMode_CHALLENGE && mode != shared_proto.TowerTeamMode_HELP {
		logrus.Debugf("变更模式，未知模式")
		err = secret_tower.ErrChangeModeFailUnknownMode
		return
	}

	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("变更模式，没有队伍")
			err = secret_tower.ErrChangeModeFailNoTeam
			return
		}

		if team.IsLeader(hc.Id()) {
			logrus.Debugf("变更模式，队长不可以变更")
			err = secret_tower.ErrChangeModeFailIsLeader
			return
		}

		_, member := team.GetMember(hc.Id())
		if member == nil {
			logrus.Debugf("变更模式，玩家竟然没在队伍里面")
			err = secret_tower.ErrChangeModeFailServerError
			return
		}

		if member.mode == mode {
			logrus.Debugf("变更模式，模式没有变化")
			err = secret_tower.ErrChangeModeFailModeNotChange
			return
		}

		if m.heroDataService.FuncNotError(hc.Id(), func(hero *entity.Hero) (heroChanged bool) {
			heroSecretTower := hero.SecretTower()

			if mode == shared_proto.TowerTeamMode_CHALLENGE {
				if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
					logrus.Debugf("变更模式，没有次数了")
					err = secret_tower.ErrChangeModeFailNoTimes
					return
				}
			} else {
				// 协助模式
				//if team.towerData.Id == heroSecretTower.MaxOpenSecretTowerId() {
				//	logrus.Debugf("变更模式，不可以协助最高自己能够参加的最高层数的密室")
				//	err = secret_tower.ErrChangeModeFailCanNotHelpMaxTower
				//	return
				//}

				if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
					logrus.Debugf("变更模式，没有协助次数了")
					err = secret_tower.ErrChangeModeFailNoHelpTimes
					return
				}

				// 没有挑战次数了
				//if hero.GuildId() == 0 {
				//	logrus.Debugf("变更模式，没有联盟不可以变更为协助模式")
				//	err = secret_tower.ErrChangeModeFailCanNotHelpNoGuildMember
				//	return
				//}
				//
				//existChallengeGuildMember := false
				//
				//for _, member := range team.Members() {
				//	if snapshot := member.heroSnapshotGetter(); snapshot != nil && snapshot.GuildId == hero.Guild() {
				//		existChallengeGuildMember = true
				//		break
				//	}
				//}
				//
				//if !existChallengeGuildMember {
				//	logrus.Debugf("变更模式，没有盟友在队伍中，无法变更模式")
				//	err = secret_tower.ErrChangeModeFailCanNotHelpNoGuildMember
				//	return
				//}
			}

			return
		}) {
			if err == nil {
				err = secret_tower.ErrChangeModeFailServerError
			}
			return
		}

		if err != nil {
			// 次数什么的不够，里面已经处理了
			return
		}

		member.ChangeMode(mode)
		team.invalidCache()
		hc.Send(secret_tower.NewS2cChangeModeMsg(int32(mode)))
		team.broadcastIgnore(secret_tower.NewS2cOtherChangedChangeModeMsg(member.IdBytes(), int32(mode)), hc.Id())
		ok = true
	})

	return
}

func (m *SecretTowerModule) invite(inviteId []byte, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	inviteHeroId, e := idbytes.ToId(inviteId)
	if !e {
		logrus.Debugf("邀请他人加入队伍，解析客户端发送过来的id失败")
		hc.Send(secret_tower.NewS2cFailTargetNotFoundInviteMsg(inviteId))
		return
	}

	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("邀请他人加入队伍，你没有队伍，无法邀请")
			err = secret_tower.ErrInviteFailNoTeam
			return
		}

		if team.IsFull() {
			logrus.Debugf("邀请他人加入队伍，队伍已满，无法邀请")
			err = secret_tower.ErrInviteFailTeamFull
			return
		}

		targetTeam := m.manager.getHeroJoinTeam(inviteHeroId)
		if targetTeam == team {
			logrus.Debugf("邀请他人加入队伍，目标在你的队伍中了，无法邀请")
			hc.Send(secret_tower.NewS2cFailTargetInYourTeamInviteMsg(inviteId))
			return
		}

		if !m.worldService.IsOnline(inviteHeroId) {
			logrus.Debugf("邀请他人加入队伍，目标不在线，无法邀请")
			hc.Send(secret_tower.NewS2cFailTargetNotOnlineInviteMsg(inviteId))
			return
		}

		success := false

		m.heroDataService.Func(inviteHeroId, func(hero *entity.Hero, err error) (heroChanged bool) {
			if err != nil {
				if err == lock.ErrEmpty {
					hc.Send(secret_tower.NewS2cFailTargetNotFoundInviteMsg(inviteId))
					return
				}

				logrus.WithError(err).Errorf("查看其它玩家信息，lock hero fail", hc.Id())
				err = secret_tower.ErrInviteFailServerError
				return
			}

			if team.guildId > 0 && hero.GuildId() != team.guildId {
				logrus.Debugf("邀请他人加入队伍，目标不是相同联盟的，无法邀请")
				hc.Send(secret_tower.NewS2cFailTargetNotInMyGuildInviteMsg(inviteId))
				return
			}

			heroSecretTower := hero.SecretTower()
			if !heroSecretTower.HasOpen(team.towerData.Id) {
				logrus.Debugf("邀请他人加入队伍，目标本层没有开启，无法邀请")
				hc.Send(secret_tower.NewS2cFailTargetNotOpenInviteMsg(inviteId))
				return
			}

			if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
				if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
					logrus.Debugf("加入重楼密室队伍，没有次数了")
					hc.Send(secret_tower.NewS2cFailTargetNoTimesInviteMsg(inviteId))
					return
				}
			}

			success = true

			return
		})

		if !success {
			// 返回
			return
		}

		// 即使已经发过了，也要提示邀请成功
		hc.Send(secret_tower.NewS2cInviteMsg(inviteId))

		m.manager.invites.OnInviteJoinTeam(inviteHeroId, team.teamId, m.timeService.CurrentTime())

		ok = true
	})

	return
}

func (m *SecretTowerModule) inviteAll(inviteIds [][]byte, hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	inviteHeroIds, e := idbytes.ToIds(inviteIds)
	if !e {
		logrus.Debugf("一键邀请他人加入队伍，解析客户端发送过来的id失败")
		hc.Send(secret_tower.ERR_INVITE_ALL_FAIL_INVALID_ID)
		return
	}

	if check.Int64Duplicate(inviteHeroIds) {
		logrus.Debugf("一键邀请他人加入队伍，发送的id存在重复id")
		hc.Send(secret_tower.ERR_INVITE_ALL_FAIL_INVALID_ID)
		return
	}

	var toInviteIds []int64
	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugf("邀请他人加入队伍，你没有队伍，无法邀请")
			err = secret_tower.ErrInviteAllFailNoTeam
			return
		}

		if team.IsFull() {
			logrus.Debugf("邀请他人加入队伍，队伍已满，无法邀请")
			err = secret_tower.ErrInviteAllFailTeamFull
			return
		}

		for _, inviteHeroId := range inviteHeroIds {

			targetTeam := m.manager.getHeroJoinTeam(inviteHeroId)
			if targetTeam == team {
				continue
			}

			if !m.worldService.IsOnline(inviteHeroId) {
				continue
			}

			m.heroDataService.FuncNotError(inviteHeroId, func(hero *entity.Hero) (heroChanged bool) {
				if team.guildId > 0 && hero.GuildId() != team.guildId {
					return
				}

				heroSecretTower := hero.SecretTower()
				if !heroSecretTower.HasOpen(team.towerData.Id) {
					return
				}

				if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
					if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
						return
					}
				}

				toInviteIds = append(toInviteIds, inviteHeroId)
				return
			})

		}

		// 即使已经发过了，也要提示邀请成功
		hc.Send(secret_tower.NewS2cInviteAllMsg(inviteIds))

		if len(toInviteIds) > 0 {
			m.manager.invites.OnInviteAllJoinTeam(toInviteIds, team.teamId, m.timeService.CurrentTime())
		}

		ok = true
	})

	return
}

var noInviteListMsg = secret_tower.NewS2cRequestInviteListMsg([][]byte{}).Static()

func (m *SecretTowerModule) requestInviteList(hc iface.HeroController) {
	count, _ := m.manager.invites.InviteMeTeamAndHaveNew(hc.Id())
	if count <= 0 {
		hc.Send(noInviteListMsg)
		return
	}

	inviteList := make([][]byte, 0, count)

	m.manager.invites.WalkInviteMeTeam(hc.Id(), func(teamId int64) (toContinue bool) {
		team := m.manager.getTeam(teamId)
		if team == nil {
			logrus.Errorf("请求邀请列表，存在我有的队伍id，但是队伍找不到")
			return
		}

		inviteList = append(inviteList, team.Encode4Show())
		return true
	})

	hc.Send(secret_tower.NewS2cRequestInviteListMsg(inviteList))

	// 刷新上次请求时间
	if m.manager.invites.RefreshLastGetInviteMeTeamTime(hc.Id(), m.timeService.CurrentTime()) {
		hc.Send(secret_tower.NewS2cReceiveInviteMsg(count, false))
	}
}

func (m *SecretTowerModule) requestTeamDetail(hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugln("请求队伍详情，没有队伍")
			err = secret_tower.ErrRequestTeamDetailFailNoTeam
			return
		}

		hc.Send(team.TeamDetailMsg())

		ok = true
	})
	return
}

func (m *SecretTowerModule) startChallenge(hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	// 武将设置不可以解雇
	// 减少玩家次数
	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugln("开启挑战，但是队伍找不到")
			err = secret_tower.ErrStartChallengeFailNoTeam
			return
		}

		if !team.IsLeader(hc.Id()) {
			logrus.Debugln("开启挑战，但是不是队长")
			err = secret_tower.ErrStartChallengeFailNotTeamLeader
			return
		}

		if team.IsMemberNotEnough() {
			logrus.Debugln("开启挑战，人数不够")
			err = secret_tower.ErrStartChallengeFailTeamMemberNotEnough
			return
		}

		ctime := m.timeService.CurrentTime()
		if ctime.Before(team.ProtectEndTime()) {
			logrus.Debugln("开启挑战，请等待其他人都准备好")
			err = secret_tower.ErrStartChallengeFailWaitProtectEnd
			return
		}

		challengeCount := 0

		// 遍历队员
		for _, member := range team.Members() {
			heroName, guildId, _, guildFlagName := member.HeroNameAndGuildInfo()
			if team.guildId > 0 {
				if guildId == 0 {
					logrus.Debugln("开启挑战，xxx没有联盟")
					hc.Send(secret_tower.NewS2cFailWithMemberNoGuildStartChallengeMsg(member.IdBytes(), heroName, guildFlagName))
					return
				}

				if guildId != team.guildId {
					logrus.Debugln("开启挑战，xxx不是我们联盟的")
					hc.Send(secret_tower.NewS2cFailWithMemberNotMyGuildStartChallengeMsg(member.IdBytes(), heroName, guildFlagName))
					return
				}
			}

			if member.mode != shared_proto.TowerTeamMode_HELP {
				challengeCount++
				continue
			}

			//if guildId <= 0 {
			//	logrus.Debugln("开启挑战，xxx开的是协助模式，但是没有联盟")
			//	hc.Send(secret_tower.NewS2cFailWithMemberIsHelpButNoGuildStartChallengeMsg(member.IdBytes(), heroName, guildFlagName))
			//	return
			//}
			//
			//haveHelpGuildMember := false
			//// 协助模式必须有队员是我的盟友，且必须开着挑战模式
			//for _, guildMember := range team.Members() {
			//	if guildMember == member {
			//		continue
			//	}
			//
			//	if _, guildMemberGuildId, _, _ := guildMember.HeroNameAndGuildInfo(); guildMemberGuildId != guildId {
			//		continue
			//	}
			//
			//	if guildMember.mode != shared_proto.TowerTeamMode_CHALLENGE {
			//		continue
			//	}
			//
			//	haveHelpGuildMember = true
			//	break
			//}
			//
			//if !haveHelpGuildMember {
			//	logrus.Debugln("开启挑战，xxx开的是协助模式，但是没有盟友")
			//	heroName, _, _, guildFlagName := member.HeroNameAndGuildInfo()
			//	hc.Send(secret_tower.NewS2cFailWithMemberIsHelpButNoGuildMemberStartChallengeMsg(member.IdBytes(), heroName, guildFlagName))
			//	return
			//}
		}

		if challengeCount <= 0 {
			logrus.Debugln("开启挑战，没有人是挑战模式，无法开启")
			err = secret_tower.ErrStartChallengeFailNoChallengePeople
			return
		}

		challengeIdBytes := make([][]byte, 0, challengeCount)

		memberCount := team.MemberCount()
		attackerIds := make([]int64, 0, memberCount)
		attackerProtos := make([]*shared_proto.CombatPlayerProto, 0, memberCount)
		for _, member := range team.Members() {
			var errMsg pbutil.Buffer
			if m.heroDataService.FuncNotError(member.Id(), func(hero *entity.Hero) (heroChanged bool) {
				heroSecretTower := hero.SecretTower()
				if member.mode == shared_proto.TowerTeamMode_CHALLENGE {
					// 挑战，看挑战次数够不够
					if heroSecretTower.ChallengeTimes() >= m.maxTimes() {
						logrus.Debugln("开启挑战，xxx挑战次数不够")
						heroName, _, _, guildFlagName := member.HeroNameAndGuildInfo()
						errMsg = secret_tower.NewS2cFailWithMemberTimesNotEnoughStartChallengeMsg(member.IdBytes(), heroName, guildFlagName)
						return
					}
				} else {
					if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
						logrus.Debugln("开启挑战，xxx协助次数不够")
						heroName, _, _, guildFlagName := member.HeroNameAndGuildInfo()
						errMsg = secret_tower.NewS2cFailWithMemberHelpTimesNotEnoughStartChallengeMsg(member.IdBytes(), heroName, guildFlagName)
						return
					}
				}

				combatPlayerProto, failType := hero.GenCombatPlayerProto(true, shared_proto.PveTroopType_DUNGEON, m.guildService.GetSnapshot)
				switch failType {
				case entity.SUCCESS:
					break
				case entity.SERVER_ERROR:
					err = secret_tower.ErrStartChallengeFailServerError
					return
				case entity.CAPTAIN_COUNT_NOT_ENOUGH:
					err = secret_tower.ErrStartChallengeFailNoChallengePeople
					return
				default:
					logrus.Errorf("开启挑战，未处理的错误类型 %v", failType)
					err = secret_tower.ErrStartChallengeFailServerError
					return
				}

				if combatPlayerProto == nil {
					logrus.Error("开启挑战，xxx武将找不到")
					err = secret_tower.ErrStartChallengeFailNoChallengePeople
					return
				}

				attackerIds = append(attackerIds, member.Id())
				attackerProtos = append(attackerProtos, combatPlayerProto)

				if member.mode == shared_proto.TowerTeamMode_CHALLENGE {
					challengeIdBytes = append(challengeIdBytes, member.IdBytes())
				}

				return
			}) {
				err = secret_tower.ErrStartChallengeFailServerError
				return
			}

			if err != nil {
				// 里面处理了
				return
			}

			if errMsg != nil {
				hc.Send(errMsg)
				return
			}
		}

		// 再检查一遍
		if len(challengeIdBytes) <= 0 {
			logrus.Debugln("开启挑战，没有人是挑战模式，无法开启")
			err = secret_tower.ErrStartChallengeFailNoChallengePeople
			return
		}

		tfctx := entity.NewTlogFightContext(operate_type.BattleSecretTower, team.towerData.Id, team.leader.Id(), 0)

		defenserIds, defenserProtos := team.towerData.GenMonsterCombatPlayer()
		if len(attackerIds) <= 0 || len(defenserIds) <= 0 {
			logrus.Error("开启挑战，len(attackerIds) <= 0 || len(defenserIds) <= 0")
			err = secret_tower.ErrStartChallengeFailTeamMemberNotEnough
			return
		}

		defenserCount := len(defenserIds)

		// 左边队伍跟右边队伍，一对一，连续战斗，每打一次换人，直到一方死光（连胜离场也算死）
		//var attackerIndex, defenserIndex int
		//var attackerId, defenserId int64
		//var attackerProto, defenserProto *shared_proto.CombatPlayerProto
		var multiLink []string
		winTimesMap := make(map[int64]int32)

		fightMaxTimes := len(attackerProtos) + len(defenserProtos)
		isAttackerWin := false
		for i := 0; i < fightMaxTimes; i++ {

			attackerId := attackerIds[i]
			attackerProto := attackerProtos[i]

			defenserId := defenserIds[i]
			defenserProto := defenserProtos[i]

			response := m.fightService.SendFightRequest(tfctx, team.towerData.CombatScene, attackerId, defenserId, attackerProto, defenserProto)
			if response == nil {
				logrus.Errorf("开启重楼密室挑战，response==nil")
				err = secret_tower.ErrStartChallengeFailServerError
				return
			}

			if response.ReturnCode != 0 {
				logrus.Errorf("开启重楼密室挑战，战斗计算发生错误，%s", response.ReturnMsg)
				err = secret_tower.ErrStartChallengeFailServerError
				return
			}

			multiLink = append(multiLink, response.Link)
			if response.AttackerWin {
				winTimes := winTimesMap[attackerId]
				winTimes++
				winTimesMap[attackerId] = winTimes

				if i+1 >= len(defenserProtos) {
					// 防守方输了
					isAttackerWin = true
					break
				}

				if team.towerData.MaxAttackerContinuewWinTimes == 0 ||
					uint64(winTimes) < team.towerData.MaxAttackerContinuewWinTimes {
					// 没有连胜离场
					attackerIds = append(attackerIds, attackerId)
					attackerProtos = append(attackerProtos, attackerProto)
				} else {
					if i+1 >= len(attackerProtos) {
						// 连胜离场后，防守方还有，进攻方没有部队了，进攻方输了
						break
					}
				}

			} else {
				winTimesMap[defenserId]++

				defenserIds = append(defenserIds, defenserId)
				defenserProtos = append(defenserProtos, defenserProto)

				if i+1 >= len(attackerProtos) {
					// 进攻方输了
					break
				}
			}
		}

		team.destroy()

		var superPrizeChallengerId []byte
		if team.towerData.SuperPlunder != nil && isAttackerWin && team.towerData.RandomGiveSuperPrize() {
			superPrizeChallengerId = challengeIdBytes[rand.Intn(len(challengeIdBytes))]
		}

		challengeResultProto := &shared_proto.SecretChallengeResultProto{}

		challengeResultProto.AttackLeaderId = hc.IdBytes()
		challengeResultProto.DefenceLeaderId = idbytes.ToBytes(int64(team.towerData.MonsterLeaderId))

		challengeResultProto.AttackCount = int32(memberCount)
		challengeResultProto.DefenceCount = int32(defenserCount)

		challengeResultProto.Win = isAttackerWin
		//challengeResultProto.Link = response.Link
		//challengeResultProto.Share = response.AttackerShare
		challengeResultProto.MultiLink = multiLink
		challengeResultProto.SecretTowerId = u64.Int32(team.towerData.Id)
		challengeResultProto.Members = make([]*shared_proto.SecretMemberResultProto, 0, team.MemberCount())

		//maxContinueWinTimes := int64(0)
		//if isAttackerWin && team.guildId > 0 {
		//	for _, winTimes := range response.WinTimesMap {
		//		if winTimes > maxContinueWinTimes {
		//			maxContinueWinTimes = winTimes
		//		}
		//	}
		//}

		var totalWinTimes int32
		for _, member := range team.Members() {
			totalWinTimes += winTimesMap[member.Id()]
		}
		challengeResultProto.KillMonster = totalWinTimes
		challengeResultProto.LeftMonster = i32.Max(int32(defenserCount)-totalWinTimes, 0)

		// 密室战报
		record := &shared_proto.SecretRecordProto{
			Win:         challengeResultProto.Win,
			KillMonster: challengeResultProto.KillMonster,
			LeftMonster: challengeResultProto.LeftMonster,
			Share:       challengeResultProto.Share,
		}

		// 减少挑战次数
		hctx := heromodule.NewContext(m.dep, operate_type.SecretTowerPrize)
		for _, member := range team.Members() {

			memberResultProto := member.Encode4Result()
			memberResultProto.ContinueKillCount = winTimesMap[member.Id()]

			challengeResultProto.Members = append(challengeResultProto.Members, memberResultProto)

			// 联盟贡献
			var guildContribution uint64

			m.heroDataService.FuncWithSend(member.Id(), func(hero *entity.Hero, result herolock.LockResult) {
				heroSecretTower := hero.SecretTower()

				hero.SecretTower().AddRecord(record, m.miscData().MaxRecord)

				if member.mode == shared_proto.TowerTeamMode_CHALLENGE {
					// 挑战，看挑战次数够不够
					times := heroSecretTower.ChallengeTimes()
					if times >= m.maxTimes() {
						// 次数不够
						logrus.Errorln("开启挑战，xxx挑战次数不够")
						return
					}

					if isAttackerWin {
						heroSecretTower.IncreChallengeTimes()
						heroSecretTower.IncreHistoryChallengeTimes()
						result.Add(secret_tower.NewS2cTimesChangeMsg(u64.Int32(heroSecretTower.ChallengeTimes()), u64.Int32(heroSecretTower.HistoryChallengeTimes())))
					}
				} else {
					if !heroSecretTower.HasEnoughHelpTimes(m.miscData().MaxHelpTimes) {
						logrus.Errorln("开启挑战，xxx协助次数不够")
						return
					}

					if isAttackerWin {
						newTimes := heroSecretTower.IncreHelpTimes()
						result.Add(secret_tower.NewS2cHelpTimesChangeMsg(u64.Int32(newTimes)))

						// 重楼密室协助次数
						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_HelpSecretTower, 0)
						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_HelpSecretTower, team.towerData.Id)
						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HELP_SECRET_TOWER)
					}
				}

				if !isAttackerWin {
					// 失败了，不给奖励
					return
				}

				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_ChallengeSecretTower, 0)
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_ChallengeSecretTower, team.towerData.Id)
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_SECRET_TOWER)
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HISTORY_CHALLENGE_SECRET_TOWER)

				if member.mode == shared_proto.TowerTeamMode_HELP {
					// 协助模式是没奖励的
					guildContribution = team.towerData.GuildHelpContribution
				} else {
					if team.towerData.SuperPlunder != nil && bytes.Equal(superPrizeChallengerId, hero.IdBytes()) {
						superPrize := team.towerData.SuperPlunder.Try()
						// 加超级奖励
						heromodule.AddPveTroopPrize(hctx, hero, result, superPrize, shared_proto.PveTroopType_DUNGEON, ctime)

						challengeResultProto.SuperPrizeId = superPrizeChallengerId
						challengeResultProto.SuperPrize = superPrize.Encode()
					}

					if !heroSecretTower.HasFirstPass(team.towerData.Id) {
						heroSecretTower.GiveFirstPassPrize(team.towerData.Id)
						heromodule.AddPveTroopPrize(hctx, hero, result, team.towerData.FirstPassPrize, shared_proto.PveTroopType_DUNGEON, ctime)

						memberResultProto.FirstPassPrize = true
					}

					prize := team.towerData.Plunder.Try()
					heromodule.AddPveTroopPrize(hctx, hero, result, prize, shared_proto.PveTroopType_DUNGEON, ctime)
					memberResultProto.Prize = prize.Encode()

					//if team.guildId > 0 {
					//	guildContribution = team.towerData.CalcGuih.todayAddCollectGuildContribution ldModeContribution(team.MemberCount(), response.WinTimesMap[member.Id()], maxContinueWinTimes)
					//}
				}

				if guildContribution > 0 {
					//guildContribution = heroSecretTower.AddTodayGuildContributionAmount(m.miscData.MaxGuildContribution, guildContribution)
					heroSecretTower.AddTodayGuildContributionAmount(guildContribution)

					// 加
					memberResultProto.GuildContribution += u64.Int32(guildContribution)

					if guildContribution > 0 {
						// 不管后面的联盟数据里他有没有联盟，都直接在这里加了
						hero.AddGuildContributionCoin(guildContribution)
						result.Add(guild.NewS2cUpdateContributionCoinMsg(u64.Int32(hero.GetGuildContributionCoin())))
					}
				}

				result.Changed()
				result.Ok()
			})

			if guildContribution > 0 {
				// 加
				m.guildService.TimeoutFunc(func(guilds sharedguilddata.Guilds) {
					_, guildId, _, _ := member.HeroNameAndGuildInfo()
					g := guilds.Get(guildId)
					if g == nil {
						return
					}

					member := g.GetMember(member.Id())
					if member == nil {
						return
					}

					member.AddContribution(guildContribution)

					// 清缓存
					m.guildService.ClearSelfGuildMsgCache(g.Id())
				})
			}
		}

		team.broadcast(secret_tower.NewS2cBroadcastStartChallengeMsg(must.Marshal(challengeResultProto)))
		ok = true
	})

	return
}

var changeToNotGuildModeMsg = secret_tower.NewS2cChangeGuildModeBroadcastMsg(0).Static()

func (m *SecretTowerModule) changeGuildMode(hc iface.HeroController) (processed, ok bool, err msg.ErrMsg) {
	processed = m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			logrus.Debugln("变更队伍联盟模式，但是队伍找不到")
			err = secret_tower.ErrChangeGuildModeFailNoTeam
			return
		}

		if !team.IsLeader(hc.Id()) {
			logrus.Debugln("变更队伍联盟模式，但是不是队长")
			err = secret_tower.ErrChangeGuildModeFailNotLeader
			return
		}

		if team.guildId > 0 {
			team.guildId = 0
			// 广播
			team.broadcast(changeToNotGuildModeMsg)
		} else {
			guildId, _ := hc.LockGetGuildId()
			if guildId == 0 {
				logrus.Debugln("变更队伍联盟模式，没有加入联盟")
				err = secret_tower.ErrChangeGuildModeFailNoGuild
				return
			}

			for _, member := range team.Members() {
				_, memGuildId, _, _ := member.HeroNameAndGuildInfo()
				if memGuildId != guildId {
					logrus.Debugln("变更队伍联盟模式，有玩家不是我们联盟的")
					err = secret_tower.ErrChangeGuildModeFailSbNotInMyGuild
					return
				}
			}
			team.guildId = guildId
			// 广播
			team.broadcast(secret_tower.NewS2cChangeGuildModeBroadcastMsg(i64.Int32(guildId)))
		}

		team.invalidCache()

		hc.Send(secret_tower.CHANGE_GUILD_MODE_S2C)

		ok = true
	})

	return
}

func (m *SecretTowerModule) quickQueryTeamBasic(proto *secret_tower.C2SQuickQueryTeamBasicProto, hc iface.HeroController) {
	// 写死，最多请求五支队伍
	maxCount := imath.Min(len(proto.Ids), 5)

	basics := make([][]byte, 0, maxCount)

	notExistTeamIds := []int32{}

	for i := 0; i < maxCount; i++ {
		rTeam := m.manager.getRTeam(int64(proto.Ids[i]))
		if rTeam == nil {
			notExistTeamIds = append(notExistTeamIds, proto.Ids[i])
			continue
		}

		basics = append(basics, rTeam.Encode4Show())
	}

	hc.Send(secret_tower.NewS2cQuickQueryTeamBasicMsg(basics, notExistTeamIds))
}

//gogen:iface
func (m *SecretTowerModule) ProcessTeamTalk(proto *secret_tower.C2STeamTalkProto, hc iface.HeroController) {
	wordsId := u64.FromInt32(proto.WordsId)

	if wordsId > 0 && m.dep.Datas().GetSecretTowerWordsData(wordsId) == nil {
		hc.Send(secret_tower.ERR_TEAM_TALK_FAIL_INVALID_WORDS)
		return
	}

	m.queueWaitFunc(func() {
		team := m.manager.getHeroJoinTeam(hc.Id())
		if team == nil {
			hc.Send(secret_tower.ERR_TEAM_TALK_FAIL_NO_TEAM)
			return
		}

		_, self := team.GetMember(hc.Id())
		if self == nil {
			logrus.Debugf("密室 ProcessTeamTalk，队伍中找不到这个人", hc.Id())
			hc.Send(secret_tower.ERR_TEAM_TALK_FAIL_NO_TEAM)
			return
		}
		hero := self.heroSnapshotGetter()
		team.chatRecord.Add(&shared_proto.SecretTowerChatRecordProto {
			HeroId: idbytes.ToBytes(hero.Id),
			HeroName: hero.Name,
			GuildFlagName: hero.GuildFlagName(),
			WordsId: proto.WordsId,
			Text: proto.Text,
		})

		if wordsId > 0 {
			self.saidWords = wordsId
			self.saidTime = m.timeService.CurrentTime()
			team.invalidCache()
		}

		hc.Send(secret_tower.NewS2cTeamTalkMsg(proto.WordsId, proto.Text))

		for _, mem := range team.Members() {
			if mem.Id() == hc.Id() {
				continue
			}
			m.dep.World().Send(mem.Id(), secret_tower.NewS2cTeamWhoTalkMsg(idbytes.ToBytes(hc.Id()), proto.WordsId, proto.Text))
		}
	})

}
