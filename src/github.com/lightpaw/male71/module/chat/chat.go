package chat

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"github.com/eapache/queue"
	"golang.org/x/net/context"
	"sync"
	"github.com/lightpaw/male7/service/tss"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/golang/snappy"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/pb/rpcpb/game2tss"
	"strconv"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	chat2 "github.com/lightpaw/male7/service/chat"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/service/operate_type"
	"time"
)

func NewChatModule(dep iface.ServiceDep, db iface.DbService, pushService iface.PushService,
	tssClient iface.TssClient, serverConfig iface.IndividualServerConfig, chatService iface.ChatService,
	baizhanService iface.BaiZhanService, mcWar iface.MingcWarService) *ChatModule {
	m := &ChatModule{}
	m.dep = dep
	m.db = db
	m.time = dep.Time()
	m.datas = dep.Datas()
	m.world = dep.World()
	m.guildService = dep.Guild()
	m.heroSnapshotService = dep.HeroSnapshot()
	m.pushService = pushService
	m.tssClient = tssClient
	m.serverConfig = serverConfig
	m.broadcast = dep.Broadcast()
	m.chatService = chatService
	m.baizhanService = baizhanService
	m.mcWar = mcWar

	// 注册tss回调
	tssClient.RegisterCallback(tss.Chat, m.tssChatCallback)

	return m
}

//gogen:iface
type ChatModule struct {
	dep                 iface.ServiceDep
	db                  iface.DbService
	time                iface.TimeService
	datas               iface.ConfigDatas
	world               iface.WorldService
	guildService        iface.GuildService
	heroSnapshotService iface.HeroSnapshotService
	pushService         iface.PushService
	tssClient           iface.TssClient
	serverConfig        iface.IndividualServerConfig
	broadcast           iface.BroadcastService
	chatService         iface.ChatService
	baizhanService      iface.BaiZhanService
	mcWar               iface.MingcWarService
}

type chatCache struct {
	sync.RWMutex

	maxCount int

	//cache []*shared_proto.ChatMsgProto
	queue *queue.Queue

	cacheMsg    pbutil.Buffer
	cacheTopMsg pbutil.Buffer
}

//gogen:iface
func (m *ChatModule) ProcessWorldChat(proto *chat.C2SWorldChatProto, hc iface.HeroController) {

	hc.Send(chat.WORLD_CHAT_S2C)

	//if len(proto.Text) <= 0 {
	//	return
	//}
	//
	//var toSend pbutil.Buffer
	//if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	g := m.guildService.GetSnapshot(hero.GuildId())
	//	var guildFlagName string
	//	if g != nil {
	//		guildFlagName = g.FlagName
	//	}
	//
	//	var whiteFlagGuildFlagName string
	//	if g := m.guildService.GetSnapshot(hero.GetWhiteFlagGuildId()); g != nil {
	//		whiteFlagGuildFlagName = g.FlagName
	//	}
	//
	//	toSend = chat.NewS2cWorldOtherChatMsg(hero.IdBytes(), hero.Name(), hero.Head(), guildFlagName, proto.Text, whiteFlagGuildFlagName)
	//
	//	result.Ok()
	//
	//	heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHAT_TIMES)
	//
	//}) {
	//	return
	//}
	//
	//m.world.BroadcastIgnore(toSend, hc.Id())
}

//gogen:iface
func (m *ChatModule) ProcessGuildChat(proto *chat.C2SGuildChatProto, hc iface.HeroController) {
	hc.Send(chat.GUILD_CHAT_S2C)

	//if len(proto.Text) <= 0 {
	//	hc.Send(chat.GUILD_CHAT_S2C)
	//	return
	//}
	//
	//var guildId int64
	//var toSend pbutil.Buffer
	//if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	guildId = hero.GuildId()
	//	if guildId == 0 {
	//		result.Add(chat.ERR_GUILD_CHAT_FAIL_NOT_GUILD)
	//		return
	//	}
	//
	//	result.Add(chat.GUILD_CHAT_S2C)
	//
	//	var whiteFlagGuildFlagName string
	//	if g := m.guildService.GetSnapshot(hero.GetWhiteFlagGuildId()); g != nil {
	//		whiteFlagGuildFlagName = g.FlagName
	//	}
	//
	//	toSend = chat.NewS2cGuildOtherChatMsg(hero.IdBytes(), hero.Name(), hero.Head(), proto.Text, whiteFlagGuildFlagName)
	//
	//	result.Ok()
	//
	//	heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHAT_TIMES)
	//}) {
	//	return
	//}
	//
	//snapshot := m.guildService.GetSnapshot(guildId)
	//if snapshot == nil {
	//	hc.Send(chat.ERR_GUILD_CHAT_FAIL_NOT_GUILD)
	//	return
	//}
	//
	//m.world.MultiSend(i64.RemoveIfPresent(snapshot.UserMemberIds, hc.Id()), toSend)
}

//gogen:iface c2s_self_chat_window
func (m *ChatModule) ProcessSelfChatWindow(hc iface.HeroController) {
	var unreadCount []uint64
	var lastMsg [][]byte
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		unreadCount, lastMsg, err = m.db.ListHeroChatWindow(ctx, hc.Id())
		return
	}); err != nil {
		logrus.WithError(err).Error("获取聊天窗口列表")
		hc.Send(chat.NewS2cSelfChatWindowMsg(nil, nil))
		return
	}

	hc.Send(chat.NewS2cSelfChatWindowMsg(lastMsg, u64.Int32Array(unreadCount)))
}

//gogen:iface
func (m *ChatModule) ProcessCreateSelfChatWindow(proto *chat.C2SCreateSelfChatWindowProto, hc iface.HeroController) {

	hc.Send(chat.NewS2cCreateSelfChatWindowMsg(proto.Target, proto.SetUp))

	targetId, ok := idbytes.ToId(proto.Target)
	if !ok || targetId == 0 {
		logrus.Debug("添加聊天窗口，无效的id")
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("添加聊天窗口，目标是Npc")
		return
	}

	target := m.dep.HeroSnapshot().Get(targetId)
	if target == nil {
		logrus.Debug("添加聊天窗口，目标是不存在")
		return
	}

	senderId := hc.Id()

	targetSender := m.chatService.GetChatSender(targetId)

	roomId := chat2.NewChatRoomId(hc.Id(), shared_proto.ChatType_ChatPrivate, targetId)

	ctime := m.time.CurrentTime()

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		return m.db.UpdateChatWindow(ctx, senderId, must.Marshal(roomId), must.Marshal(targetSender), false, timeutil.Marshal32(ctime), proto.SetUp)
	}); err != nil {
		logrus.WithError(err).Error("发送聊天，保存聊天失败")
	}
}

//gogen:iface
func (m *ChatModule) ProcessRemoveChatWindow(proto *chat.C2SRemoveChatWindowProto, hc iface.HeroController) {

	defer hc.Send(chat.NewS2cRemoveChatWindowMsg(proto.ChatType, proto.ChatTarget))

	t := shared_proto.ChatType(proto.ChatType)
	if t != shared_proto.ChatType_ChatPrivate {
		logrus.Debug("移除聊天，无效的类型")
		return
	}

	targetId, ok := i64.FromBytes(proto.ChatTarget)
	if !ok || targetId == 0 || npcid.IsNpcId(targetId) || targetId == hc.Id() {
		logrus.Debug("移除聊天，无效的目标")
		return
	}

	roomId := chat2.NewChatRoomId(hc.Id(), t, targetId)
	roomIdBytes := must.Marshal(roomId)

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.DeleteChatWindow(ctx, hc.Id(), roomIdBytes)
		return
	}); err != nil {
		logrus.WithError(err).Error("移除聊天失败")
		return
	}
}

//gogen:iface
func (m *ChatModule) ProcessReadChatMsg(proto *chat.C2SReadChatMsgProto, hc iface.HeroController) {
	defer hc.Send(chat.NewS2cReadChatMsgMsg(proto.ChatType, proto.ChatTarget))

	t := shared_proto.ChatType(proto.ChatType)
	if t != shared_proto.ChatType_ChatPrivate {
		logrus.Debug("已读聊天，无效的类型")
		return
	}

	targetId, ok := i64.FromBytes(proto.ChatTarget)
	if !ok || targetId == 0 || npcid.IsNpcId(targetId) || targetId == hc.Id() {
		logrus.Debug("已读聊天，无效的目标")
		return
	}

	roomId := chat2.NewChatRoomId(hc.Id(), t, targetId)
	roomIdBytes := must.Marshal(roomId)

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.ReadChat(ctx, hc.Id(), roomIdBytes)
		return
	}); err != nil {
		logrus.WithError(err).Error("已读聊天失败")
		return
	}
}

var emptyHistoryChat = chat.NewS2cListHistoryChatMarshalMsg(nil)

//gogen:iface
func (m *ChatModule) ProcessListHistoryChat(proto *chat.C2SListHistoryChatProto, hc iface.HeroController) {

	minChatId, ok := i64.FromBytes(proto.MinChatId)
	if !ok {
		logrus.Debug("获取聊天记录，无效的MinChatId")
		hc.Send(emptyHistoryChat)
		return
	}

	if t := shared_proto.ChatType(proto.ChatType); t != shared_proto.ChatType_ChatMcWar {
		var targetId int64
		switch t {
		case shared_proto.ChatType_ChatWorld:
			// 世界聊天，第一次取数据，从缓存中获取
			if minChatId == 0 {
				hc.Send(m.chatService.GetCacheMsg())
				return
			}
		case shared_proto.ChatType_ChatGuild:
			if guildId, ok := hc.LockGetGuildId(); !ok {
				logrus.Debug("获取聊天记录，获取玩家联盟id失败")
				hc.Send(emptyHistoryChat)
				return
			} else {
				if guildId == 0 {
					logrus.Debug("获取聊天记录，玩家没有联盟")
					hc.Send(emptyHistoryChat)
					return
				}

				targetId = guildId
			}

		case shared_proto.ChatType_ChatPrivate:
			targetId, ok = i64.FromBytes(proto.ChatTarget)
			if !ok {
				logrus.Debug("获取聊天记录，解析ChatTarget失败")
				hc.Send(emptyHistoryChat)
				return
			}
			if targetId == 0 || npcid.IsNpcId(targetId) || targetId == hc.Id() {
				logrus.Debug("获取聊天记录，无效的目标")
				hc.Send(emptyHistoryChat)
				return
			}
		case shared_proto.ChatType_ChatSystem:
			sndMsg := m.chatService.GetSystemChatRecord(minChatId)
			if sndMsg != nil {
				hc.Send(sndMsg)
			} else {
				hc.Send(emptyHistoryChat)
			}
			return
		default:
			logrus.Debug("获取聊天记录，无效的类型")
			hc.Send(emptyHistoryChat)
			return
		}

		roomId := chat2.NewChatRoomId(hc.Id(), t, targetId)
		roomIdBytes := must.Marshal(roomId)

		var msgProto []*shared_proto.ChatMsgProto
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			msgProto, err = m.db.ListHeroChatMsg(ctx, roomIdBytes, uint64(minChatId))
			return
		}); err != nil {
			logrus.WithError(err).Error("获取聊天记录失败")
			hc.Send(emptyHistoryChat)
			return
		}

		if len(msgProto) <= 0 {
			hc.Send(emptyHistoryChat)
			return
		}

		data := make([][]byte, 0, len(msgProto))
		msgs := make([]*shared_proto.ChatMsgProto, 0, len(msgProto))
		for _, v := range msgProto {
			//m.updateProtoHeroInfo(v)
			data = append(data, must.Marshal(v))
			msgs= append(msgs, v)
		}
		hc.Send(chat.NewS2cListHistoryChatMarshalMsg(data))

	} else {
		sndMsg := m.mcWar.CatchHistoryRecord(hc.Id(), minChatId)
		if sndMsg != nil {
			hc.Send(sndMsg)
		} else {
			hc.Send(emptyHistoryChat)
		}
	}
}

func (m *ChatModule) updateProtoHeroInfo(proto *shared_proto.ChatMsgProto) {
	if proto == nil {
		return
	}

	senderProto := proto.Sender
	if senderProto == nil {
		return
	}

	senderId, ok := idbytes.ToId(senderProto.Id)
	if !ok || senderId <= 0 {
		return
	}

	// 暂时先从缓存中取
	heroSnapshot := m.heroSnapshotService.GetFromCache(senderId)
	if heroSnapshot == nil {
		return
	}
	senderProto.Name = heroSnapshot.Name
	senderProto.Level = u64.Int32(heroSnapshot.Level)
	senderProto.GuildFlag = heroSnapshot.GuildFlagName()
	senderProto.Head = heroSnapshot.Head
}

//gogen:iface
func (m *ChatModule) ProcessSendChat(proto0 *chat.C2SSendChatProto, hc iface.HeroController) {

	proto := &shared_proto.ChatMsgProto{}
	if err := proto.Unmarshal(proto0.ChatMsg); err != nil {
		logrus.Debug("发送聊天，解析聊天Proto失败")
		hc.Send(chat.ERR_SEND_CHAT_FAIL_INVALID_TARGET)
		return
	}

	if c := uint64(util.GetCharLen(proto.Text)); c > m.datas.MiscConfig().ChatTextLength {
		if proto.MsgType == int32(shared_proto.ChatMsgType_ChatMsgText) ||
			proto.MsgType == int32(shared_proto.ChatMsgType_ChatMsgGuildAllMembers) {
			logrus.Debug("发送聊天，文字长度超出限制")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_TEXT_TOO_LONG)
			return
		}

		// 分享回放，这个不做处理
	}

	if c := uint64(util.GetCharLen(proto.Json)); c > m.datas.MiscConfig().ChatJsonLength {
		logrus.Debug("发送聊天，Json长度超出限制")
		hc.Send(chat.ERR_SEND_CHAT_FAIL_TEXT_TOO_LONG)
		return
	}

	ctime := m.time.CurrentTime()

	t := shared_proto.ChatType(proto.ChatType)

	var guildId int64
	var heroLevel uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if ctime.Before(hero.MiscData().GetBanChatEndTime()) {
			logrus.Debug("被系统禁言")
			result.Add(chat.ERR_SEND_CHAT_FAIL_BAN_CHAT)
			return
		}

		guildId = hero.GuildId()
		heroLevel = hero.Level()

		switch t {
		case shared_proto.ChatType_ChatWorld:
			if heroLevel < m.datas.MiscConfig().WorldChatLevel {
				logrus.Debug("发送聊天，世界聊天等级不足")
				result.Add(chat.ERR_SEND_CHAT_FAIL_WORLD_CHAT_LEVEL)
				return
			}

			if ctime.Before(hero.LastWorldChatTime().Add(m.datas.MiscConfig().WorldChatDuration)) {
				logrus.Debug("发送world聊天，CD未到")
				result.Add(chat.ERR_SEND_CHAT_FAIL_WORLD_CHAT_TOO_FAST)
				return
			}

			if proto.Laba {
				if !heromodule.HasEnoughGoodsOrBuy(hero, m.datas.MiscConfig().BroadcastGoods, 1, proto.AutoBuyLaba) {
					logrus.Debug("发送聊天，喇叭消耗不足")
					result.Add(chat.ERR_SEND_CHAT_FAIL_BROADCAST_GOODS_NOT_ENOUGH)
					return
				}
			}
		default:
			if ctime.Before(hero.LastChatTime().Add(m.datas.MiscConfig().ChatDuration)) {
				logrus.Debug("发送聊天，CD未到")
				result.Add(chat.ERR_SEND_CHAT_FAIL_CHAT_TOO_FAST)
				return
			}
		}

		// 和策划确认过以后特殊消息可能会改变频道发送，所以放在这里单独判定，不嵌套在某个频道中了
		var cd time.Duration
		cdType := server_proto.OperationCDType_InvalidCDType
		switch shared_proto.ChatMsgType(proto.MsgType) {
		case shared_proto.ChatMsgType_ChatMsgSecretTower:
			cd = m.datas.MiscGenConfig().SecretTowerCd
			cdType = server_proto.OperationCDType_SecretTowerShout
		case shared_proto.ChatMsgType_ChatMsgFight:
			cd = timeutil.MaxDuration(m.datas.MiscGenConfig().BaizhanCd, m.datas.MiscGenConfig().XuanyuanCd)
			cdType = server_proto.OperationCDType_FightShare
		case shared_proto.ChatMsgType_ChatMsgHebi:
			cd = m.datas.MiscGenConfig().HebiCd
			cdType = server_proto.OperationCDType_HebiShout
		case shared_proto.ChatMsgType_ChatMsgXiongnu:
			cd = m.datas.MiscGenConfig().XiongnuCd
			cdType = server_proto.OperationCDType_XiongnuRemind
		case shared_proto.ChatMsgType_ChatMsgReport:
			cd = m.datas.MiscGenConfig().MailCd
			cdType = server_proto.OperationCDType_MailShare
		}
		if cdType != server_proto.OperationCDType_InvalidCDType {
			if ctime.Before(hero.Misc().GetNextOperationTime(cdType)) {
				logrus.Debug("发送聊天，发送太频繁")
				result.Add(chat.ERR_SEND_CHAT_FAIL_TOO_FAST)
				return
			}
			hero.Misc().SetNextOperationTime(cdType, ctime.Add(cd))
		}

		result.Changed()
		result.Ok()
	}) {
		return
	}

	var targetId int64
	var targetSender *shared_proto.ChatSenderProto
	var targetSenderBytes []byte
	switch t {
	case shared_proto.ChatType_ChatWorld:
	case shared_proto.ChatType_ChatGuild:
		if guildId == 0 {
			logrus.Debug("发送聊天，玩家没有联盟")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_NOT_IN_GUILD)
			return
		}

		if proto.MsgType == int32(shared_proto.ChatMsgType_ChatMsgGuildAskForHelp) {
			if ctime.Sub(hc.LastClickTime()) < 5*time.Second {
				logrus.Debug("发送聊天，求援按钮点击太频繁")
				hc.Send(chat.ERR_SEND_CHAT_FAIL_TOO_FAST)
				return
			}
			hc.SetLastClickTime(ctime)
		}

		targetId = guildId
		proto.ChatTarget = i64.ToBytes(targetId)

	case shared_proto.ChatType_ChatGuildAllMembers:
		// 联盟全体消息
		if guildId <= 0 {
			logrus.Debugf("发送聊天，联盟ID错误 id:%v", guildId)
			hc.Send(chat.ERR_SEND_CHAT_FAIL_NOT_IN_GUILD)
			return
		}

		var errMsg pbutil.Buffer
		m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				logrus.Debugf("发送联盟全体消息，没找到联盟 id:%v。", guildId)
				errMsg = chat.ERR_SEND_CHAT_FAIL_NOT_IN_GUILD
				return
			}

			member := g.GetMember(hc.Id())
			if member == nil {
				logrus.Debugf("发送联盟全体消息，没找到联盟成员 gid:%v。hero id:%v", guildId, hc.Id())
				errMsg = chat.ERR_SEND_CHAT_FAIL_NOT_IN_GUILD
				return
			}

			if !member.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
				return permission.SendToAllMembers
			}) {
				logrus.Debugf("发送联盟全体消息，没有权限")
				errMsg = chat.ERR_SEND_CHAT_FAIL_GUILD_PERM_DENY
				return
			}
		})

		if errMsg != nil {
			hc.Send(errMsg)
			return
		}

		targetId = guildId
		proto.ChatTarget = i64.ToBytes(targetId)

	case shared_proto.ChatType_ChatPrivate:
		if heroLevel < m.datas.MiscConfig().ChatPrivateMinLevel {
			hc.Send(chat.ERR_SEND_CHAT_FAIL_PRIVATE_CHAT_LEVEL)
			return
		}
		// 私聊
		var ok bool
		targetId, ok = i64.FromBytes(proto.ChatTarget)
		if !ok {
			logrus.Debug("发送聊天，解析ChatTarget失败")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_INVALID_TARGET)
			return
		}

		if targetId == 0 || npcid.IsNpcId(targetId) || targetId == hc.Id() {
			logrus.Debug("发送聊天，无效的目标")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_INVALID_TARGET)
			return
		}

		target := m.heroSnapshotService.Get(targetId)
		if target == nil {
			logrus.Debug("发送聊天，目标不存在")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_INVALID_TARGET)
			return
		}

		targetSender = m.chatService.GetChatSender(targetId)
		targetSenderBytes = must.Marshal(targetSender)

	case shared_proto.ChatType_ChatMcWar:
		if state, _, _ := m.mcWar.CurrMcWarStage(); state != int32(shared_proto.MingcWarState_MC_T_FIGHT) {
			hc.Send(chat.ERR_SEND_CHAT_FAIL_MC_WAR_NOT_IN_FIGHT_STAGE)
			return
		}
		if _, ok := m.mcWar.JoiningFightMingc(hc.Id()); !ok {
			hc.Send(chat.ERR_SEND_CHAT_FAIL_MC_WAR_NOT_JOIN_FIGHT)
			return
		}

	default:
		logrus.Debug("发送聊天，无效的类型")
		hc.Send(chat.ERR_SEND_CHAT_FAIL_INVALID_TARGET)
		return
	}

	// 设置数据
	ctime = m.time.CurrentTime()
	proto.SendTime = timeutil.Marshal32(ctime)
	proto.Sender = m.chatService.GetChatSender(hc.Id())

	roomId := chat2.NewChatRoomId(hc.Id(), t, targetId)
	roomIdBytes := must.Marshal(roomId)

	//if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
	//	if err0 := m.db.AddChatMsg(ctx, hc.Id(), roomIdBytes, proto); err0 != nil {
	//		if err == nil {
	//			err = err0
	//		}
	//	}
	//
	//	if t == shared_proto.ChatType_ChatPrivate {
	//		if err0 := m.db.UpdateChatWindow(ctx, hc.Id(), roomIdBytes, targetSenderBytes, false, proto.SendTime); err0 != nil {
	//			if err == nil {
	//				err = err0
	//			}
	//		}
	//
	//		if err0 := m.db.UpdateChatWindow(ctx, targetId, roomIdBytes, must.Marshal(proto.Sender), true, proto.SendTime); err0 != nil {
	//			if err == nil {
	//				err = err0
	//			}
	//		}
	//	}
	//
	//	return
	//}); err != nil {
	//	logrus.WithError(err).Error("发送聊天，保存聊天失败")
	//}
	//
	//// 发送成功，给自己
	//hc.Send(chat.NewS2cSendChatMsg(proto.ChatId, targetSenderBytes))
	//
	// 发给其他人
	//switch t {
	//case shared_proto.ChatType_ChatWorld:
	//	m.worldChatCache.addMsg(proto)
	//	m.world.BroadcastIgnore(chat.NewS2cOtherSendChatMarshalMsg(proto), hc.Id())
	//case shared_proto.ChatType_ChatGuild:
	//	if g := m.guildService.GetSnapshot(targetId); g != nil {
	//		m.world.MultiSendIgnore(g.UserMemberIds, chat.NewS2cOtherSendChatMarshalMsg(proto), hc.Id())
	//	}
	//case shared_proto.ChatType_ChatPrivate:
	//	m.world.SendFunc(targetId, func() pbutil.Buffer {
	//		return chat.NewS2cOtherSendChatMarshalMsg(proto)
	//	})
	//}

	if m.tssClient.IsEnable() && len(proto.Text) > 0 {
		// 发送到tss验证敏感词
		callbackData := &server_proto.TssChatCallbackProto{
			SenderId:     hc.Id(),
			TargetId:     targetId,
			RoomIdBytes:  roomIdBytes,
			TargetSender: targetSenderBytes,
			Proto:        proto,
		}
		callbackDataBytes := snappy.Encode(nil, must.Marshal(callbackData))

		// TODO
		toSend := &game2tss.C2SUicJudgeUserInputChatV2Proto{
			Openid:       strconv.FormatInt(hc.Id(), 10),      // OpenId TODO
			Platid:       0,                                   // 0-IOS, 1-Android TODO
			WorldId:      int64(m.serverConfig.GetServerID()), // sid
			MsgCategory:  int32(tss.Chat),                     // 消息内容类别 1： 邮件；  2：聊天
			ChannelId:    int64(t),                            // 聊天频道
			ClientIp:     int32(hc.GetClientIp32()),           // 客户端ip /* 127.0.0.1 */ = 0x100007f
			RoldId:       hc.Id(),
			RoleName:     proto.Sender.Name,
			RoleLevel:    proto.Sender.Level,
			Msg:          proto.Text,
			CallbackData: callbackDataBytes,
		}

		resp, err := m.tssClient.JudgeChat(toSend)
		if err != nil {
			logrus.WithError(err).Error("发送聊天，tss敏感词校验出错")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_SERVER_ERROR)
			return
		}

		if resp.Ret != 0 {
			logrus.WithField("ret", resp.Ret).WithField("ret_msg", resp.RetMsg).Error("发送聊天，tss敏感词校验失败")
			hc.Send(chat.ERR_SEND_CHAT_FAIL_SERVER_ERROR)
			return
		}

		// tss接收了，这里不处理，等待回调

	} else {
		// 不验证，直接发出来
		m.doSendMsg(hc.Id(), targetId, roomIdBytes, targetSenderBytes, proto, "")
	}
}

func (m *ChatModule) doSendMsg(senderId, targetId int64, roomIdBytes, targetSender []byte, proto *shared_proto.ChatMsgProto, replaceText string) {

	t := proto.ChatType

	isLabaBroadcast := t == shared_proto.ChatType_ChatWorld && proto.Laba

	hctx := heromodule.NewContext(m.dep, operate_type.ChatLaba)

	var guildId int64
	var heroName string
	if m.dep.HeroData().FuncWithSend(senderId, func(hero *entity.Hero, result herolock.LockResult) {

		if isLabaBroadcast {
			if !heromodule.TryReduceOrBuyGoods(hctx, hero, result, m.datas.MiscConfig().BroadcastGoods, 1, proto.AutoBuyLaba) {
				logrus.Debug("发送聊天，喇叭消耗不足")
				result.Add(chat.ERR_SEND_CHAT_FAIL_BROADCAST_GOODS_NOT_ENOUGH)
				return
			}
		}

		guildId = hero.GuildId()
		heroName = hero.Name()

		// 加聊天cd
		switch t {
		case shared_proto.ChatType_ChatWorld:
			hero.SetLastWorldChatTime(timeutil.Unix32(proto.SendTime))
		default:
			hero.SetLastChatTime(timeutil.Unix32(proto.SendTime))
		}

		result.Ok()
	}) {
		return
	}

	if isLabaBroadcast {
		// 喇叭广播
		var guildFlag string
		g := m.dep.GuildSnapshot().GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
		}

		msg := misc.NewS2cHeroBroadcastMsg(proto.Text, heroName, guildFlag).Static()
		m.world.Broadcast(msg)
	}

	pushFunc := func(d *pushdata.PushData) (title, content string) {
		return d.Title, "[" + heroName + "]: " + proto.Text
	}

	if t == shared_proto.ChatType_ChatGuildAllMembers {
		// 特殊处理联盟全体消息，转成私聊类型。也会发自己

		g := m.guildService.GetSnapshot(targetId)
		if g == nil {
			logrus.Debugf("发送联盟全体消息，没找到联盟 id:%v。", targetId)
			m.world.Send(senderId, chat.ERR_SEND_CHAT_FAIL_NOT_IN_GUILD)
			return
		}

		// 自己的成功响应消息
		m.world.Send(senderId, chat.NewS2cSendChatMsg(proto.ChatId, targetSender, replaceText))

		// 给联盟成员发私聊
		proto.ChatType = shared_proto.ChatType_ChatPrivate
		for _, memId := range g.UserMemberIds {
			if senderId == memId {
				// 不发给自己
				continue
			}

			member := m.heroSnapshotService.Get(memId)
			if member == nil {
				logrus.Debug("发送联盟全体消息，联盟成员不存在")
				continue
			}

			proto.ChatTarget = idbytes.ToBytes(memId)

			targetSender := m.chatService.GetChatSender(memId)
			targetSenderBytes := must.Marshal(targetSender)

			roomId := chat2.NewChatRoomId(senderId, proto.ChatType, memId)
			roomIdBytes := must.Marshal(roomId)

			m.chatService.SaveDB(senderId, memId, proto.ChatType, roomIdBytes, targetSenderBytes, proto)

			m.world.SendFunc(memId, func() pbutil.Buffer {
				return chat.NewS2cOtherSendChatMarshalMsg(proto)
			})

			if memId != senderId {
				if len(proto.Text) > 0 && proto.MsgType == int32(shared_proto.ChatMsgType_ChatMsgText) {
					m.pushService.PushFunc(shared_proto.SettingType_ST_CHAT_PRIVATE, targetId, pushFunc)
				}
			}
		}

		return
	}

	// 发给其他人
	switch t {
	case shared_proto.ChatType_ChatWorld:
		m.chatService.SaveDB(senderId, targetId, t, roomIdBytes, targetSender, proto) // 保存DB
		m.chatService.AddMsg(proto)
		m.world.BroadcastIgnore(chat.NewS2cOtherSendChatMarshalMsg(proto), senderId)

	case shared_proto.ChatType_ChatGuild:
		m.chatService.SaveDB(senderId, targetId, t, roomIdBytes, targetSender, proto) // 保存DB
		if g := m.guildService.GetSnapshot(targetId); g != nil {
			m.world.MultiSendIgnore(g.UserMemberIds, chat.NewS2cOtherSendChatMarshalMsg(proto), senderId)
		}
	case shared_proto.ChatType_ChatPrivate:
		m.chatService.SaveDB(senderId, targetId, t, roomIdBytes, targetSender, proto) // 保存DB
		m.world.SendFunc(targetId, func() pbutil.Buffer {
			return chat.NewS2cOtherSendChatMarshalMsg(proto)
		})

		if len(proto.Text) > 0 && proto.MsgType == int32(shared_proto.ChatMsgType_ChatMsgText) {
			m.pushService.PushFunc(shared_proto.SettingType_ST_CHAT_PRIVATE, targetId, pushFunc)
		}
	case shared_proto.ChatType_ChatMcWar:
		if errMsg := m.mcWar.SendChat(senderId, proto); errMsg != nil {
			m.world.Send(senderId, errMsg)
			return
		}
	}
	// 发送成功，最后才发给自己（名城战生成ChatID机制决定）
	m.world.Send(senderId, chat.NewS2cSendChatMsg(proto.ChatId, targetSender, replaceText))

	m.dep.Tlog().TlogChatFlowById(senderId, uint64(proto.ChatType), u64.FromInt64(targetId), 0)
}

func (m *ChatModule) tssChatCallback(heroId int64, msgResultFlag int32, replaceMsg string, callbackData []byte) {

	// 0:合法 1:不合法，不能显示 2:合法，但包含敏感词
	if msgResultFlag == 1 {
		logrus.Debug("tss聊天回调，内容不合法，不能显示")
		m.world.Send(heroId, chat.ERR_SEND_CHAT_FAIL_SENSITIVE_WORDS)
		return
	}

	if len(callbackData) <= 0 {
		logrus.Debug("tss聊天回调，protoBytes.len <= 0")
		return
	}

	protoBytes, err := snappy.Decode(nil, callbackData)
	if err != nil {
		logrus.WithError(err).Error("tss聊天回调，snappy.Decode() 失败")
		return
	}

	proto := &server_proto.TssChatCallbackProto{}
	if err := proto.Unmarshal(protoBytes); err != nil {
		logrus.WithError(err).Error("tss聊天回调，TssChatCallbackProto.Unmarshal() 失败")
		return
	}

	replaceText := ""
	if msgResultFlag == 2 {
		proto.Proto.Text = replaceMsg
		replaceText = replaceMsg
	}

	m.doSendMsg(proto.SenderId, proto.TargetId, proto.RoomIdBytes, proto.TargetSender, proto.Proto, replaceText)
}

//gogen:iface
func (m *ChatModule) ProcessGetHeroChatInfo(proto *chat.C2SGetHeroChatInfoProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debug("聊天获取玩家信息，无效的id")
		hc.Send(chat.NewS2cGetHeroChatInfoMsg(proto.Id, 0, 0))
		return
	}

	hero := m.heroSnapshotService.Get(targetId)
	if hero == nil {
		logrus.Debug("聊天获取玩家信息，玩家id不存在")
		hc.Send(chat.NewS2cGetHeroChatInfoMsg(proto.Id, 0, 0))
		return
	}

	hc.Send(chat.NewS2cGetHeroChatInfoMsg(proto.Id, u64.Int32(hero.TowerMaxFloor), u64.Int32(hero.BaiZhanJunXianLevel)))
}
