package chat

import (
	"context"
	"github.com/eapache/queue"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/sortkeys"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sync"
)

func NewChatService(db iface.DbService, push iface.PushService, time iface.TimeService, guild iface.GuildSnapshotService,
	heroSnapshot iface.HeroSnapshotService, world iface.WorldService, datas iface.ConfigDatas, redPacket iface.RedPacketService) *ChatService {
	m := &ChatService{}
	m.db = db
	m.push = push
	m.guild = guild
	m.time = time
	m.world = world
	m.heroSnapshot = heroSnapshot
	m.datas = datas
	m.redPacket = redPacket

	count := u64.Int(u64.Max(datas.MiscConfig().ChatBatchCount, 1))
	m.systemChatRecord = NewChatRecord(count, datas.MiscGenConfig().FirstHistoryChatSend)
	m.worldChatCache = &ChatCache{
		maxCount: count,
		queue:    queue.New(),
	}

	// 从DB中load初始数据上来
	go func() {
		_ = ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			if msgs, err := db.ListHeroChatMsg(ctx, must.Marshal(WorldChatRoomId), 0); err != nil {
				logrus.WithError(err).Error("初始化世界聊天失败")
			} else {
				m.worldChatCache.initMsg(msgs)
			}
			return
		})
	}()

	heromodule.RegisterHeroOnlineListener(m)

	return m
}

//gogen:iface
type ChatService struct {
	db           iface.DbService
	push         iface.PushService
	guild        iface.GuildSnapshotService
	time         iface.TimeService
	world        iface.WorldService
	heroSnapshot iface.HeroSnapshotService
	datas        iface.ConfigDatas
	redPacket    iface.RedPacketService

	worldChatCache   *ChatCache
	systemChatRecord *ChatRecord // systemCache
}

func (m *ChatService) BroadcastSystemChat(text string) {
	proto := m.buildSysChatProto(0, 0, shared_proto.ChatType_ChatSystem, text, shared_proto.ChatMsgType_ChatMsgText, true, true, false)
	m.systemChatRecord.AddChat(proto)
	m.world.Broadcast(chat.NewS2cOtherSendChatMarshalMsg(proto))
}

func (m *ChatService) GetSystemChatRecord(minChatId int64) pbutil.Buffer {
	if minChatId == 0 {
		return m.systemChatRecord.GetFirstChatRecord()
	}
	return m.systemChatRecord.GetChatRecored(minChatId)
}

// 系统自动聊天，有DB操作
func (m *ChatService) SysChat(senderId, targetId int64, chatType shared_proto.ChatType, text string, msgType shared_proto.ChatMsgType, showSysSender, save, isSys, isJson bool) {
	m.SysChatSendFunc(senderId, targetId, chatType, text, msgType, showSysSender, save, isSys, isJson, func(proto *shared_proto.ChatMsgProto) {
		m.defaultSendFunc(chatType, targetId, proto)
	})
}

func (m *ChatService) defaultSendFunc(chatType shared_proto.ChatType, targetId int64, proto *shared_proto.ChatMsgProto) {
	switch chatType {
	case shared_proto.ChatType_ChatWorld:
		m.worldChatCache.addMsg(proto)
		m.world.Broadcast(chat.NewS2cOtherSendChatMarshalMsg(proto))

	case shared_proto.ChatType_ChatGuild:
		if g := m.guild.GetSnapshot(targetId); g != nil {
			m.world.MultiSend(g.UserMemberIds, chat.NewS2cOtherSendChatMarshalMsg(proto))
		}
	case shared_proto.ChatType_ChatPrivate:
		m.world.SendFunc(targetId, func() pbutil.Buffer {
			return chat.NewS2cOtherSendChatMarshalMsg(proto)
		})
	}
}

var emptyFunc = func(proto *shared_proto.ChatMsgProto) {}

func (m *ChatService) SysChatProtoFunc(senderId, targetId int64, chatType shared_proto.ChatType, text string, msgType shared_proto.ChatMsgType, showSysSender, save, isSys, isJson bool, f constants.ChatFunc) (chatId int64) {
	return m.SysChatFunc(senderId, targetId, chatType, text, msgType, showSysSender, save, isSys, isJson, f, func(proto *shared_proto.ChatMsgProto) {
		m.defaultSendFunc(chatType, targetId, proto)
	})
}

func (m *ChatService) SysChatSendFunc(senderId, targetId int64, chatType shared_proto.ChatType, text string, msgType shared_proto.ChatMsgType, showSysSender, save, isSys, isJson bool, f constants.ChatFunc) (chatId int64) {
	return m.SysChatFunc(senderId, targetId, chatType, text, msgType, showSysSender, save, isSys, isJson, emptyFunc, f)
}

func (m *ChatService) SysChatFunc(senderId, targetId int64, chatType shared_proto.ChatType, text string, msgType shared_proto.ChatMsgType, showSysSender, save, isSys, isJson bool, protoFunc, sendFunc constants.ChatFunc) (chatId int64) {
	proto := m.buildSysChatProto(senderId, targetId, chatType, text, msgType, showSysSender, isSys, isJson)

	protoFunc(proto)

	roomId := NewChatRoomId(senderId, chatType, targetId)
	roomIdBytes := must.Marshal(roomId)

	var targetSender *shared_proto.ChatSenderProto
	var targetSenderBytes []byte

	if chatType == shared_proto.ChatType_ChatPrivate {
		if targetId == 0 || npcid.IsNpcId(targetId) || targetId == senderId {
			logrus.Debug("发送系统聊天，无效的目标")
			return
		}

		target := m.heroSnapshot.Get(targetId)
		if target == nil {
			logrus.Debug("发送系统聊天，目标不存在")
			return
		}

		targetSender = newChatSender(targetId, target)
		targetSenderBytes = must.Marshal(targetSender)
	} else if chatType == shared_proto.ChatType_ChatSystem {
		if targetId == 0 || npcid.IsNpcId(targetId) {
			logrus.Debug("发送系统聊天，无效的目标")
			return
		}
		if msgType == shared_proto.ChatMsgType_ChatMsgFriendAdded {
			if targetId == senderId {
				logrus.Debug("发送系统聊天，无效的目标")
				return
			}

		}
	}

	// 保存DB
	if save {
		chatId = m.SaveDB(senderId, targetId, chatType, roomIdBytes, targetSenderBytes, proto)
	}

	sendFunc(proto)

	return
}

func (m *ChatService) buildSysChatProto(senderId, targetId int64, chatType shared_proto.ChatType, text string, msgType shared_proto.ChatMsgType, showSysSender, isSys, isJson bool) (proto *shared_proto.ChatMsgProto) {
	ctime := m.time.CurrentTime()
	proto = &shared_proto.ChatMsgProto{Sys: isSys}
	proto.ChatType = chatType
	proto.ChatTarget = idbytes.ToBytes(targetId)
	proto.MsgType = int32(msgType)
	if isJson {
		proto.Json = text
	} else {
		proto.Text = text
	}
	proto.SendTime = timeutil.Marshal32(ctime)
	if showSysSender {
		proto.Sender = m.getSysChatSender()
	} else {
		proto.Sender = m.GetChatSender(senderId)
	}
	return
}

func (m *ChatService) SaveDB(senderId, targetId int64, chatType shared_proto.ChatType, roomIdBytes, targetSender []byte, proto *shared_proto.ChatMsgProto) (chatId int64) {

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		if id, err0 := m.db.AddChatMsg(ctx, senderId, roomIdBytes, proto); err0 != nil {
			if err == nil {
				err = err0
			}
		} else {
			chatId = id
		}

		if chatType == shared_proto.ChatType_ChatPrivate {
			if senderId > 0 {
				if err0 := m.db.UpdateChatWindow(ctx, senderId, roomIdBytes, targetSender, false, proto.SendTime, true); err0 != nil {
					if err == nil {
						err = err0
					}
				}
			}

			if targetId > 0 {
				if err0 := m.db.UpdateChatWindow(ctx, targetId, roomIdBytes, must.Marshal(proto.Sender), true, proto.SendTime, true); err0 != nil {
					if err == nil {
						err = err0
					}
				}
			}
		}

		return
	}); err != nil {
		logrus.WithError(err).Error("发送聊天，保存聊天失败")
	}

	return
}

func (m *ChatService) UpdateDBRedPacket(chatId int64, allGrabbed bool) (succ bool) {
	if chatId <= 0 {
		return
	}

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		var msg *shared_proto.ChatMsgProto
		if msg, err = m.db.LoadChatMsg(ctx, chatId); err != nil {
			return
		}
		if msg == nil {
			logrus.Debugf("更新聊天红包状态失败, 红包聊天id:%v 不存在", chatId)
			return
		}

		obj := entity.ChatJsonObjMarshal(msg.Json)
		if obj == nil {
			return
		}

		obj.AllGrabbed = allGrabbed
		if msg.Json, err = obj.BuildJson(); err != nil {
			return
		}

		succ = m.db.UpdateChatMsg(ctx, chatId, msg)
		return
	}); err != nil {
		logrus.WithError(err).Error("更新聊天红包状态失败")
	}

	return
}

func (m *ChatService) getSysChatSender() *shared_proto.ChatSenderProto {
	return &shared_proto.ChatSenderProto{
		Id:                     idbytes.ToBytes(0),
		Name:                   m.datas.TextHelp().SysChatSenderName.New().JsonString(),
		Head:                   "", // 默认值
		GuildFlag:              "", // 默认值
		WhiteFlagGuildFlagName: "", // 默认值
		Level: 1,
	}
}

func (m *ChatService) GetChatSender(heroId int64) *shared_proto.ChatSenderProto {
	return newChatSender(heroId, m.heroSnapshot.Get(heroId))
}

func newChatSender(heroId int64, hero *snapshotdata.HeroSnapshot) *shared_proto.ChatSenderProto {
	if hero != nil {
		return &shared_proto.ChatSenderProto{
			Id:                     hero.IdBytes,
			Name:                   hero.Name,
			Head:                   hero.Head,
			GuildFlag:              hero.GuildFlagName(),
			WhiteFlagGuildFlagName: "", // TODO 插白旗
			Level: u64.Int32(hero.Level),

			Basic: hero.EncodeBasic4Client(),
			Title: u64.Int32(hero.Title),
		}
	} else {
		return &shared_proto.ChatSenderProto{
			Id:                     idbytes.ToBytes(heroId),
			Name:                   idbytes.PlayerName(heroId),
			Head:                   "", // 默认值
			GuildFlag:              "", // 默认值
			WhiteFlagGuildFlagName: "", // 默认值
			Level: 1,

			Basic: snapshotdata.NewIdBasicProto(heroId),
			Title: 0, // 默认值
		}
	}
}

var WorldChatRoomId = &shared_proto.ChatRoomId{T: shared_proto.ChatType_ChatWorld, MemberIds: [][]byte{{0}}}

func NewChatRoomId(heroId int64, t shared_proto.ChatType, targetId int64) *shared_proto.ChatRoomId {
	if t == shared_proto.ChatType_ChatWorld {
		return WorldChatRoomId
	}

	proto := &shared_proto.ChatRoomId{}
	proto.T = t

	switch t {
	case shared_proto.ChatType_ChatGuild:
		proto.MemberIds = append(proto.MemberIds, i64.ToBytes(targetId))
	case shared_proto.ChatType_ChatPrivate:
		if heroId < targetId {
			proto.MemberIds = append(proto.MemberIds, i64.ToBytes(heroId), i64.ToBytes(targetId))
		} else {
			proto.MemberIds = append(proto.MemberIds, i64.ToBytes(targetId), i64.ToBytes(heroId))
		}
	}

	return proto
}

type ChatCache struct {
	sync.RWMutex

	maxCount int

	//cache []*shared_proto.ChatMsgProto
	queue *queue.Queue

	cacheMsg    pbutil.Buffer
	cacheTopMsg pbutil.Buffer
}

func (c *ChatService) OnHeroOnline(hc iface.HeroController) {
	hc.Send(c.worldChatCache.getCacheTopMsg())

	f := func() bool {
		if c.db.CallingTimes() > constants.DBBusyCallingTimes {
			return false
		}

		// 如果英雄离线期间有私聊信息，发消息通知
		var unreadChatCount uint64
		if err := ctxfunc.NetTimeout1s(func(ctx context.Context) (err error) {
			unreadChatCount, err = c.db.LoadUnreadChatCount(ctx, hc.Id())
			return
		}); err != nil {
			logrus.WithError(err).Error("玩家上线，获取未读私聊信息错误")
		}

		if unreadChatCount > 0 {
			// 有私聊消息，发消息通知
			hc.Send(chat.OFFLINE_CHAT_S2C)
		}

		return true
	}

	if !f() {
		hc.AddTickFunc(f)
	}

}

func (c *ChatCache) initMsg(toAdds []*shared_proto.ChatMsgProto) {
	if len(toAdds) <= 0 {
		return
	}

	c.Lock()
	defer c.Unlock()

	if c.queue.Length() >= c.maxCount {
		// 缓存中已经达到最大个数，什么都不干了
		return
	}

	// 2个列表放入map中去重，然后变回list，然后排序，然后取前X个，一个个加回去
	msgMap := make(map[int64]*shared_proto.ChatMsgProto)
	for i := 0; i < c.queue.Length(); i++ {
		if cc := c.queue.Remove().(*shared_proto.ChatMsgProto); cc != nil {
			if id, ok := i64.FromBytes(cc.ChatId); ok {
				msgMap[id] = cc
			}
		}
	}

	for _, cc := range toAdds {
		if cc != nil {
			if id, ok := i64.FromBytes(cc.ChatId); ok {
				msgMap[id] = cc
			}
		}
	}

	// 排序
	var keys []int64
	for k := range msgMap {
		keys = append(keys, k)
	}
	sortkeys.Int64s(keys)

	// 取前X个，一个个加回去
	startIndex := 0
	if len(keys) > c.maxCount {
		startIndex = c.maxCount - len(keys)
	}
	for i := startIndex; i < len(keys); i++ {
		c.queue.Add(msgMap[keys[i]])
	}

	c.cacheMsg = nil
	c.cacheTopMsg = nil
}

func (c *ChatService) AddMsg(toAdd *shared_proto.ChatMsgProto) {
	c.worldChatCache.addMsg(toAdd)
}

func (c *ChatCache) addMsg(toAdd *shared_proto.ChatMsgProto) {
	c.Lock()
	defer c.Unlock()

	if c.queue.Length() >= c.maxCount {
		// 移除一个
		c.queue.Remove()
	}

	// 添加到尾巴，最新的消息放在最后的位置
	c.queue.Add(toAdd)

	c.cacheMsg = nil
	c.cacheTopMsg = nil
}

func (c *ChatService) GetCacheMsg() pbutil.Buffer {
	return c.worldChatCache.getCacheMsg()
}

func (c *ChatCache) getCacheMsg() pbutil.Buffer {
	c.RLock()
	defer c.RUnlock()

	if c.cacheMsg != nil {
		return c.cacheMsg
	}

	n := c.queue.Length()
	c.cacheMsg = c.buildCacheMsg(n)

	return c.cacheMsg
}

func (c *ChatCache) getCacheTopMsg() pbutil.Buffer {
	c.RLock()
	defer c.RUnlock()

	if c.cacheTopMsg != nil {
		return c.cacheTopMsg
	}

	n := imath.Min(c.queue.Length(), 2)
	c.cacheTopMsg = c.buildCacheMsg(n)

	return c.cacheTopMsg
}

func (c *ChatCache) buildCacheMsg(n int) pbutil.Buffer {
	data := make([][]byte, 0, n)
	msgs := make([]*shared_proto.ChatMsgProto, 0, n)
	// 最新的消息放在最后的位置，倒序取出来
	for j := 1; j <= n; j++ {
		v := c.queue.Get(-j).(*shared_proto.ChatMsgProto)
		data = append(data, must.Marshal(v))
		msgs = append(msgs, v)
	}

	return chat.NewS2cListHistoryChatMarshalMsg(data).Static()
}
