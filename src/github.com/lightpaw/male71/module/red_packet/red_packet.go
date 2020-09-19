package red_packet

import (
	"github.com/lightpaw/logrus"
	red_packet_conf "github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/red_packet"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/u64"
	"strings"
	"unicode/utf8"
)

func NewRedPacketModule(dep iface.ServiceDep, srv iface.RedPacketService) *RedPacketModule {
	m := &RedPacketModule{}
	m.dep = dep
	m.srv = srv

	return m
}

//gogen:iface
type RedPacketModule struct {
	dep iface.ServiceDep
	srv iface.RedPacketService
}

//gogen:iface
func (m *RedPacketModule) ProcessCreate(proto *red_packet.C2SCreateProto, hc iface.HeroController) {
	dataId := u64.FromInt32(proto.DataId)
	data := m.dep.Datas().GetRedPacketData(dataId)
	if data == nil {
		hc.Send(red_packet.ERR_CREATE_FAIL_INVALID_DATA_ID)
		return
	}

	count := u64.FromInt32(proto.Count)
	if count <= 0 || !entity.CheckCanCreateRedPacket(data.Money, data.MinPartMoney, count) {
		hc.Send(red_packet.ERR_CREATE_FAIL_COUNT_ERR)
		return
	}

	if count < data.MinCount || count > data.MaxCount {
		logrus.Debugf("发红包，数量不符合配置的范围。count:%v data:%v", count, dataId)
		hc.Send(red_packet.ERR_CREATE_FAIL_COUNT_ERR)
		return
	}

	var gid int64

	switch proto.ChatType {
	case shared_proto.ChatType_ChatGuild:
		if id, ok := hc.LockGetGuildId(); !ok || id <= 0 {
			hc.Send(red_packet.ERR_CREATE_FAIL_NO_GUILD)
			return
		} else {
			gid = id
		}

		if m.dep.GuildSnapshot().GetSnapshot(gid).MemberCount < m.dep.Datas().MiscConfig().RedPacketGuildMemberMinCount {
			hc.Send(red_packet.ERR_CREATE_FAIL_GUILD_LIMIT)
			return
		}
	case shared_proto.ChatType_ChatWorld:
	default:
		hc.Send(red_packet.ERR_CREATE_FAIL_INVALID_CHAT_TYPE)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.HeroRedPacket().Reduce(dataId) {
			result.Add(red_packet.ERR_CREATE_FAIL_NOT_BOUGHT)
			return
		}
		result.Changed()
		result.Ok()
	}) {
		return
	}

	if utf8.RuneCountInString(proto.Text) > 140 {
		hc.Send(red_packet.ERR_CREATE_FAIL_TEXT_TOO_LONG)
		return
	}

	text := strings.TrimSpace(proto.Text)
	if text == "" {
		text = data.DefaultText
	}

	id, jsonStr, errMsg := m.srv.Create(hc.Id(), data, count, text, proto.ChatType)
	if errMsg != nil {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			hero.HeroRedPacket().Add(dataId, 1)
			result.Changed()
			result.Ok()
		})

		hc.Send(errMsg.ErrMsg())
		return
	}

	hc.Send(red_packet.NewS2cCreateMsg(proto.DataId))

	// 推聊天
	chatId := m.dep.Chat().SysChatProtoFunc(hc.Id(), gid, proto.ChatType, jsonStr, shared_proto.ChatMsgType_ChatMsgRedPacket, false, true, false, true, func(proto *shared_proto.ChatMsgProto) {
		proto.RedPacketId = idbytes.ToBytes(id)
	})

	m.srv.SetRedPacketChatId(id, chatId)
}

//gogen:iface
func (m *RedPacketModule) ProcessGrab(proto *red_packet.C2SGrabProto, hc iface.HeroController) {
	id, ok := idbytes.ToId(proto.Id)
	if !ok {
		hc.Send(red_packet.ERR_GRAB_FAIL_INVALID_ID)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok {
		hc.Send(red_packet.ERR_GRAB_FAIL_NOT_SAME_GUILD)
		return
	}

	grabMoney, allGrabbed, packet, errMsg := m.srv.Grab(id, hc.Id(), gid)
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	if packet == nil {
		hc.Send(red_packet.ERR_GRAB_FAIL_SERVER_ERR)
		return
	}

	data := m.dep.Datas().GetRedPacketData(u64.FromInt32(packet.DataId))
	if data == nil {
		logrus.Errorf("抢红包，dataId 不存在:%v", data)
		hc.Send(red_packet.ERR_GRAB_FAIL_SERVER_ERR)
		return
	}

	hc.Send(red_packet.NewS2cGrabMsg(packet, u64.Int32(grabMoney)))

	hctx := heromodule.NewContext(m.dep, operate_type.RedPacketAllGrabbedPrize)

	if grabMoney <= 0 {
		return
	}

	// 加钱
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		switch data.AmountType {
		case shared_proto.AmountType_PT_YUANBAO:
			heromodule.AddYuanbao(hctx, hero, result, grabMoney)
		case shared_proto.AmountType_PT_DIANQUAN:
			heromodule.AddDianquan(hctx, hero, result, grabMoney)
		}
		result.Ok()
	})

	if allGrabbed {
		m.onAllGrabbed(hctx, id, data, packet)
	}
}

func (m *RedPacketModule) onAllGrabbed(hctx *heromodule.HeroContext, redPacketId int64, data *red_packet_conf.RedPacketData, redPacket *shared_proto.RedPacketProto) {
	noticeMsg := red_packet.NewS2cAllGrabbedNoticeMsg(idbytes.ToBytes(redPacketId)).Static()
	if redPacket.ChatType == shared_proto.ChatType_ChatGuild {
		gid := int64(redPacket.CreateHero.GuildId)
		m.dep.Guild().Broadcast(gid, noticeMsg)
	} else if redPacket.ChatType == shared_proto.ChatType_ChatWorld {
		m.dep.World().Broadcast(noticeMsg)
	}

	if createHeroId, ok := idbytes.ToId(redPacket.CreateHero.Id); ok {
		m.dep.HeroData().FuncWithSend(createHeroId, func(hero *entity.Hero, result herolock.LockResult) {
			ctime := m.dep.Time().CurrentTime()
			heromodule.AddPrize(hctx, hero, result, data.AllGarbbedPrize, ctime)
		})
	}

	// 更新红包状态
	m.dep.Chat().UpdateDBRedPacket(m.srv.RedPacketChatId(redPacketId), true)
}

//gogen:iface
func (m *RedPacketModule) ProcessBuy(proto *red_packet.C2SBuyProto, hc iface.HeroController) {
	dataId := u64.FromInt32(proto.DataId)
	data := m.dep.Datas().GetRedPacketData(dataId)
	if data == nil {
		hc.Send(red_packet.ERR_BUY_FAIL_INVALID_DATA_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Product().GetYuanbaoGiftLimit() < data.Money {
			result.Add(red_packet.ERR_BUY_FAIL_RECHANGE_YUANBAO_LIMIT)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.RedPacketBuy)
		if !heromodule.TryReduceCost(hctx, hero, result, data.Cost) {
			result.Add(red_packet.ERR_BUY_FAIL_COST_NOT_ENOUGH)
			return
		}

		heromodule.ReduceYuanbaoGiftLimit(hctx, hero, result, data.Money)
		hero.HeroRedPacket().Add(dataId, 1)
		result.Add(red_packet.NewS2cBuyMsg(u64.Int32(dataId), u64.Int32(hero.HeroRedPacket().Count(dataId))))
		result.Changed()
		result.Ok()
	})
}
