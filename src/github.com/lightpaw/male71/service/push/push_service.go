package push

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/rpcpb/game2login"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/config/pushdata"
)

func NewPushService(datas iface.ConfigDatas, serverConfig iface.IndividualServerConfig, worldService iface.WorldService,
	dbService iface.DbService, heroSnapshotService iface.HeroSnapshotService,
	clusterService iface.ClusterService) *PushService {
	return &PushService{
		datas:               datas,
		serverConfig:        serverConfig,
		worldService:        worldService,
		dbService:           dbService,
		heroSnapshotService: heroSnapshotService,
		clusterService:      clusterService,
		funcQueue:           event.NewFuncQueue(1024, "push"),
	}
}

//gogen:iface
type PushService struct {
	datas               iface.ConfigDatas
	serverConfig        iface.IndividualServerConfig
	worldService        iface.WorldService
	dbService           iface.DbService
	heroSnapshotService iface.HeroSnapshotService
	clusterService      iface.ClusterService
	funcQueue           *event.FuncQueue
}

func (p *PushService) Push(settingType shared_proto.SettingType, id int64) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	data := p.datas.GetPushData(uint64(settingType))
	if data == nil {
		// 这个类型的推送没有配置，不推送了
		return
	}

	logrus.Debugf("PushService.Push: 给客户端发送push %v, %d", settingType, id)
	p.funcQueue.TryFunc(func() {
		p.push(settingType, data.Title, data.Content, id)
	})
}

func (p *PushService) PushFunc(settingType shared_proto.SettingType, id int64, f pushdata.PushFunc) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	data := p.datas.GetPushData(uint64(settingType))
	if data == nil {
		// 这个类型的推送没有配置，不推送了
		return
	}

	title, content := f(data)

	p.funcQueue.TryFunc(func() {
		p.push(settingType, title, content, id)
	})
}

func (p *PushService) PushTitleContent(settingType shared_proto.SettingType, title, content string, id int64) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	logrus.Debugf("PushService.Push: 给客户端发送push %v, %d", settingType, id)
	p.funcQueue.TryFunc(func() {
		p.push(settingType, title, content, id)
	})
}

func (p *PushService) MultiPush(settingType shared_proto.SettingType, ids []int64, ignore int64) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	data := p.datas.GetPushData(uint64(settingType))
	if data == nil {
		// 这个类型的推送没有配置，不推送了
		return
	}

	logrus.Debugf("PushService.MultiPush: 给客户端发送multi push %v, %v", settingType, ids)
	p.funcQueue.TryFunc(func() {
		p.multiPush(settingType, data.Title, data.Content, ids, ignore)
	})
}

func (p *PushService) MultiPushFunc(settingType shared_proto.SettingType, ids []int64, ignore int64, f pushdata.PushFunc) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	data := p.datas.GetPushData(uint64(settingType))
	if data == nil {
		// 这个类型的推送没有配置，不推送了
		return
	}

	title, content := f(data)

	p.funcQueue.TryFunc(func() {
		p.multiPush(settingType, title, content, ids, ignore)
	})
}

func (p *PushService) MultiPushTitleContent(settingType shared_proto.SettingType, title, content string, ids []int64, ignore int64) {
	if p.serverConfig.GetDisablePush() {
		return
	}

	logrus.Debugf("PushService.MultiPush: 给客户端发送multi push %v, %v", settingType, ids)
	p.funcQueue.TryFunc(func() {
		p.multiPush(settingType, title, content, ids, ignore)
	})
}

func (p *PushService) push(settingType shared_proto.SettingType, title, content string, id int64) {
	p.pushIfSettingsOpen(settingType, title, content, 0, id)
}

func (p *PushService) multiPush(settingType shared_proto.SettingType, title, content string, ids []int64, ignore int64) {
	if len(ids) <= 0 {
		return
	}

	p.pushIfSettingsOpen(settingType, title, content, ignore, ids...)
}

func (p *PushService) pushIfSettingsOpen(settingType shared_proto.SettingType, title, content string, ignore int64, ids ...int64) {

	pushIds := make([]int64, 0, (len(ids)>>2)+1)
	otherIds := make([]int64, 0, (len(ids)>>2)+1)
	for _, id := range ids {
		if id == ignore {
			continue
		}

		if p.worldService.IsDontPush(id) {
			// 在线不处理
			continue
		}

		snapshot := p.heroSnapshotService.GetFromCache(id)
		if snapshot != nil && entity.IsOpen(snapshot.Settings, settingType) {
			pushIds = append(pushIds, id)
			continue
		}

		otherIds = append(otherIds, id)
	}

	if len(otherIds) > 0 {
		var result []int64
		err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			result, err = p.dbService.FindSettingsOpen(ctx, settingType, otherIds)
			return
		})

		if err != nil {
			logrus.WithError(err).Errorf("pushIfSettingsOpen 查找不为开启的玩家数据报错: %v", otherIds)
		} else {
			pushIds = append(pushIds, result...)
		}
	}

	p.doPush(title, content, pushIds...)
}

func (p *PushService) doPush(title, content string, ids ...int64) {
	if len(ids) <= 0 {
		return
	}

	logrus.Debugf("PushService.doPush: 给客户端发送push ids: %v", ids)

	// TODO
	startTime := int64(0)  // 开始推送时间，0表示立即推送（unix时间戳）
	expireTime := int64(0) // 过期时间，超过这个时间不再推送(unix时间戳)
	sid := uint32(p.serverConfig.GetServerID())

	ctxfunc.NetTimeout3s(func(ctx context.Context) (err error) {
		if len(ids) == 1 {
			id := ids[0]
			resp, err := game2login.Push(p.clusterService.LoginClient(), ctx, id, sid, title, content, startTime, expireTime)
			if err != nil {
				if !ctxfunc.IsRpcTimeout(err) {
					logrus.WithError(err).Error("PushService.doPush() single fail, %d", id)
				} else {
					logrus.Debug("PushService.doPush() timeout")
				}

				return nil
			}

			if !resp.Success {
				logrus.WithField("err", resp.ErrMsg).Warn("PushService.doPush() single fail")
			}

			if resp.DontPush {
				// TODO 这个玩家不再推送
			}
		} else {
			resp, err := game2login.PushMulti(p.clusterService.LoginClient(), ctx, ids, sid, title, content, startTime, expireTime)
			if err != nil {
				if !ctxfunc.IsRpcTimeout(err) {
					logrus.WithError(err).Error("PushService.doPush() multi fail")
				} else {
					logrus.Debug("PushService.doPush() timeout")
				}
				return nil
			}

			if !resp.Success {
				logrus.WithField("err", resp.ErrMsg).Warn("PushService.doPush() multi fail")
			}

		}

		return
	})
}

func (p *PushService) GmPush(id int64) {
	p.doPush("服务器测试推送标题", "服务器测试推送正文", id)
}
