package worker

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/rpcpb/game2login"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	service1 "github.com/lightpaw/male7/service"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
)

//func (m *MessageWorker) processRobotLogin(data *service.MsgData) bool {
//	if !service.IndividualServerConfig.GetIsDebug() {
//		logrus.Error("收到RobotLogin消息, 但是不是IsDebug模式")
//		m.Close()
//		return false
//	}
//	m.isRobot = true
//	heroId := int64(0)
//	m.user = service1.NewConnectedUser(heroId, m)
//
//	hero := entity.NewHero(heroId, "机器人", service.ConfigDatas.HeroInitData(), service.TimeService.CurrentTime())
//	hero.SetMale(true)
//
//	// 初始化数据
//	initHeroCreateData(hero, service.TimeService.CurrentTime())
//	hc := service1.NewHeroController(heroId, m, service.HeroDataService.NewHeroLocker(heroId))
//
//	m.user.SetHeroController(hc)
//	m.user.SetLoaded()
//
//	m.Send(login.ROBOT_LOGIN_S2C)
//
//	return true
//}

func getNotHeroLoginMsg() pbutil.Buffer {
	if service.IndividualServerConfig.GetIsDebug() {
		return login.NewS2cLoginMsg(false, false, "", nil, true, service.CountryService.TutorialCountriesProto())
	} else {
		return login.NewS2cLoginMsg(false, false, "", nil, false, service.CountryService.TutorialCountriesProto())
	}
}

const max_hero_id = 1<<56 - 1

func (m *MessageWorker) processLoginMsg(heroId int64, data *service.MsgData) bool {
	if m.user != nil {
		logrus.Errorf("worker received login msg, but already login, received: %d-%d",
			data.ModuleID, data.SequenceID)
		m.Send(login.ERR_LOGIN_FAIL_ALREADY_LOGIN)
		m.Close()
		return false
	}

	id := heroId
	if id == 0 {
		id = m.Id()
	} else {
		// 检查token为机器人
		if !m.isRobot {
			logrus.Errorf("不是机器人，但是发送了机器人登陆消息, received: %d-%d",
				data.ModuleID, data.SequenceID)
			m.Close()
			return false
		}
		//m.session.GetLoginToken().UserSelfID = id
	}

	logrus.Debugf("received login msg, %d", id)

	// 英雄id不能是负数，因为Npc联盟成员统一使用负数作为id，玩家id不能是负数
	// 英雄id最高8位不能使用，为预留id，目前使用到的系统，如出征的部队id，使用英雄id + 部队序号生成
	if id < 0 || id > max_hero_id {
		logrus.Errorf("worker.processLoginModuleMsg invalid user id(negative), %d", id)
		m.Send(login.ERR_LOGIN_FAIL_INVALID_ID)
		m.Close()
		return false
	}

	// 如果是平台登陆，应该附带token，这里获取平台信息
	tencentInfoProto := &shared_proto.TencentInfoProto{}
	if data.Proto != nil {
		if proto, ok := data.Proto.(*login.C2SLoginProto); ok && len(proto.TencentInfo) > 0 {
			if err := tencentInfoProto.Unmarshal(proto.TencentInfo); err != nil {
				logrus.WithError(err).Debug("worker.processLoginModuleMsg 解析TencentInfo失败", id)
				m.Send(login.ERR_LOGIN_FAIL_INVALID_TENCENT_INFO)
				m.Close()
				return false
			}
		}

		tencentInfoProto.ClientIP = m.clientIp
	}

	if m.pf > 0 {
		proto, ok := data.Proto.(*login.C2SLoginProto)
		if !ok {
			logrus.Errorf("不是正常登陆消息，但是又有pf值, received: %d-%d",
				data.ModuleID, data.SequenceID)
			m.Close()
			return false
		}

		if len(proto.Token) <= 0 {
			logrus.Errorf("登陆有pf值，但是没有发送token")
			m.Send(login.ERR_LOGIN_FAIL_INVALID_TOKEN)
			m.Close()
			return false
		}

		// 通过rpc向登陆服请求玩家平台相关数据
		verifySuccess := false
		ctxfunc.NetTimeout3s(func(ctx context.Context) (err error) {
			resp, err := game2login.VerifyLoginToken(service.ClusterService.LoginClient(), ctx, id, proto.Token, m.clientIp, m.pf)
			if err != nil {
				logrus.WithError(err).Error("worker.processLoginModuleMsg verify login token error, %d", id)
				return err
			}

			verifySuccess = resp.Success
			// 这里带回平台相关信息，如果有 TODO
			return nil
		})

		if !verifySuccess {
			logrus.Debug("worker.processLoginModuleMsg verify login token fail, %d", id)
			m.Send(login.ERR_LOGIN_FAIL_INVALID_TOKEN)
			m.Close()
			return false
		}
	}

	hc, err := LoadHeroController(id, m)
	if err != nil {
		logrus.WithError(err).Errorf("worker.processLoginModuleMsg services.LoadHeroController fail, %d", id)
		m.Send(login.ERR_LOGIN_FAIL_SERVER_ERROR)
		m.Close()
		return false
	}

	user := service1.NewConnectedUser(id, m, tencentInfoProto)
	if hc != nil {
		user.SetHeroController(hc)
	}

	if m.doLogin(user) {
		loginSuc := func() (suc bool) {
			hc := user.GetHeroController()
			if hc != nil {
				ctime := service.TimeService.CurrentTime()

				var tlogHero *entity.TlogHeroInfo
				var createRegion, banned bool
				var countryId uint64
				var official shared_proto.CountryOfficialType
				if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

					banned = hero.MiscData().GetBanLoginEndTime().After(service.TimeService.CurrentTime())
					if banned {
						result.Add(login.NewS2cBanLoginMsg(timeutil.DurationMarshal32(hero.MiscData().GetBanLoginEndTime().Sub(service.TimeService.CurrentTime()))))
						logrus.Debug("该账号正被查封")
						result.Add(login.ERR_LOGIN_FAIL_BAN_LOGIN)
						return
					}

					isDebug := service.IndividualServerConfig.GetIsDebug()
					result.Add(login.NewS2cLoginMsg(true, hero.Male(), hero.Head(), hero.Domestic().GetBuildingIds(), isDebug, nil))

					countryId = hero.CountryId()
					official = hero.CountryMisc().ShowOfficialType()

					if hero.BaseRegion() == 0 {
						logrus.WithField("heroId", hero.Id()).
							WithField("level", hero.BaseLevel()).
							WithField("prosperity", hero.Prosperity()).
							WithField("region", hero.BaseRegion()).
							WithField("x", hero.BaseX()).
							WithField("y", hero.BaseY()).
							Debug("英雄登陆，发现没有主城")

						if hero.Prosperity() <= 0 {
							if hero.BaseLevel() > 0 {
								// 城池繁荣度为0，但是存在等级，当成流亡处理
								hero.SetBaseLevel(0)
							}
						} else {
							if hero.BaseLevel() <= 0 {
								// 有繁荣度，没有等级，当成1级城池处理
								hero.SetBaseLevel(1)
							}

							// 只是没有城池地图，随机到别的场景
							createRegion = true
						}
						result.Changed()
					}

					tlogHero = hero.BuildFullTlogHeroInfo(ctime)
					result.Ok()
				}) {
					if !banned {
						logrus.Errorf("登陆时lock英雄失败")
					}
					m.Close()
					return false
				}

				if createRegion {
					// 在所有场景中找一遍，当发现存在主城时候，记录下来，如果只有一个，
					// 则将英雄自己的主城设置成这个，大于1个时，将所有主城删除，重新随机一个 TODO

					if service.RegionModule.InitHeroBase(hc, ctime, countryId, realmface.AddBaseHomeTransfer) {
						logrus.Errorf("登陆时发现英雄没有主城，随机创建主城失败")
						m.Close()
						return false
					} else {
						logrus.Errorf("登陆时发现英雄没有主城，随机创建主城成功")
					}
				}

				realOfficial := service.CountryService.HeroOfficial(hc.LockHeroCountry(), hc.Id())
				// 登录验证官职
				if official != realOfficial {
					logrus.Warnf("hero:%v 修正官职：%v to %v", hc.Id(), official, realOfficial)
					service.CountryService.ForceOfficialDepose(countryId, hc.Id())
					service.CountryService.ForceOfficialAppoint(countryId, hc.Id(), realOfficial)
				}

				if tencentInfo := m.user.TencentInfo(); tencentInfo != nil && tlogHero != nil {
					tlogLogin(tlogHero, tencentInfo)
				}

			} else {
				if service.IndividualServerConfig.GetSkipHeader() {
					user.SetMisc(&server_proto.UserMiscProto{
						IsTutorialComplete: true,
					})
				} else {
					var userMisc *server_proto.UserMiscProto
					err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
						userMisc, err = service.DbService.LoadUserMisc(ctx, id)
						return err
					})

					if err != nil {
						logrus.WithError(err).Errorf("worker received login msg, 加载玩家数据报错: %d-%d",
							data.ModuleID, data.SequenceID)
						m.Send(login.ERR_LOGIN_FAIL_SERVER_ERROR)
						m.Close()
						return false
					}

					if !userMisc.Created {
						if tencentInfo := m.user.TencentInfo(); tencentInfo != nil {
							data := service.TlogService.BuildAccountRegister(id, tencentInfo)
							service.TlogService.WriteLog(data)
							// 注册也算一次登录
							tlogLogin(entity.NewSimpleTlogHeroInfo(id, ""), tencentInfo)
						}

						userMisc.Created = true
					}

					user.SetMisc(userMisc)
				}

				if m.user.Misc().IsTutorialComplete {
					// 发送没有新手教程
					m.Send(getNotHeroLoginMsg())
				} else {
					m.Send(login.NewS2cTutorialProgressMsg(m.user.Misc().TutorialProgress))
				}
			}

			// 登陆日志
			gamelogs.HeroOnlineLog(constants.PID, user.Sid(), id, 0)

			return true
		}()

		if !loginSuc {
			if ok := service.WorldService.RemoveUserIfSame(user); !ok {
				logrus.Error("登录失败，将自己从WorldService移除也失败")
			}
		}

		return loginSuc
	}

	return false
}

func tlogLogin(heroInfo *entity.TlogHeroInfo, tencentInfo *shared_proto.TencentInfoProto) {
	data := service.TlogService.BuildPlayerLogin(heroInfo, tencentInfo,
		heroInfo.GuildId,
		heroInfo.OutCityCount,
		heroInfo.FriendCount,
		heroInfo.TowerMaxFloor,
		heroInfo.JunXianLevel,
		heroInfo.MaxSecretTowerId,
		heroInfo.TopLevelCaptainLevel,
		heroInfo.TopLevelCaptainId,
		heroInfo.TopFightCaptainFightAmount,
		heroInfo.TopFightCaptainId,
		heroInfo.VipLevel(),
		heroInfo.CaptainCount,
		heroInfo.AllFightAmount,
		heroInfo.TopFightTroopFightAmount,
		heroInfo.TroopCaptainIds[0],
		heroInfo.TroopFightAmount[0],
		heroInfo.TroopCaptainIds[1],
		heroInfo.TroopFightAmount[1],
		heroInfo.TroopCaptainIds[2],
		heroInfo.TroopFightAmount[2])

	service.TlogService.WriteLog(data)
}

func tlogLogout(heroInfo *entity.TlogHeroInfo, tencentInfo *shared_proto.TencentInfoProto, logoutType uint64) {

	data := service.TlogService.BuildPlayerLogout(heroInfo, tencentInfo,
		heroInfo.OnlineTime,
		logoutType,
		heroInfo.OnlineTime,
		heroInfo.GuildId,
		heroInfo.OutCityCount,
		heroInfo.FriendCount,
		heroInfo.TowerMaxFloor,
		heroInfo.JunXianLevel,
		heroInfo.MaxSecretTowerId,
		heroInfo.TopLevelCaptainLevel,
		heroInfo.TopLevelCaptainId,
		heroInfo.TopFightCaptainFightAmount,
		heroInfo.TopFightCaptainId,
		heroInfo.VipLevel(),
		heroInfo.CaptainCount,
		heroInfo.AllFightAmount,
		heroInfo.TopFightTroopFightAmount,
		heroInfo.TroopCaptainIds[0],
		heroInfo.TroopFightAmount[0],
		heroInfo.TroopCaptainIds[1],
		heroInfo.TroopFightAmount[1],
		heroInfo.TroopCaptainIds[2],
		heroInfo.TroopFightAmount[2])

	service.TlogService.WriteLog(data)
}

func LoadHeroController(id int64, sender *MessageWorker) (*service1.HeroController, error) {

	exist, err := service.HeroDataService.Exist(id)
	if err != nil {
		logrus.WithError(err).Errorf("worker.LoadHeroController lockSharedHero error, %v", id)
		return nil, err
	}

	if exist {
		return service1.NewHeroController(id, sender, sender.clientIp, sender.clientIp32, sender.pf, service.HeroDataService.NewHeroLocker(id)), nil
	}

	return nil, nil
}

func (m *MessageWorker) doLogin(user *service1.ConnectedUser) bool {
	//// 踢下线
	m.user = user
	if oldUser, ok := service.WorldService.PutConnectedUserIfAbsent(user); !ok {
		logrus.Debug("尝试踢已有的下线")
		oldUser.DisconnectAndWait(misc.ErrDisconectReasonFailKick)
		logrus.Debug("已踢下线")
		// 对方已下线

		// 重新加载一次
		hc, err := LoadHeroController(user.Id(), m)
		if err != nil {
			logrus.WithError(err).WithField("id", user.Id()).Error("踢人成功后, 再加载就出错了")
			m.Close()
			return false
		}
		if hc != nil {
			user.SetHeroController(hc)
		}

		// 再尝试放入一次
		if _, ok := service.WorldService.PutConnectedUserIfAbsent(user); !ok {
			logrus.Error("踢人下线成功之后, 自己再放入worldService还是失败")
			m.Send(login.ERR_LOGIN_FAIL_KICK)
			m.user = nil
			m.Close()
			return false
		}
	}

	// 登录成功
	return true
}
