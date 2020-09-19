package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/zhanjiang"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func newHeroZhanJiang() *HeroZhanJiang {
	h := &HeroZhanJiang{
		passGuanQia: make(map[uint64]struct{}),
	}

	return h
}

// 过关斩将
type HeroZhanJiang struct {
	openTimes     uint64              // 当前已经开启的次数
	passGuanQia   map[uint64]struct{} // 通关了的关卡
	curChallenge  *ZhanJiangChallenge // 当前的挑战，可能为空
	lastCaptainId uint64              // 上次通关的设置的武将id
}

func (h *HeroZhanJiang) OpenTimes() uint64 {
	return h.openTimes
}

func (h *HeroZhanJiang) ReduceOpenTimes(toReduce uint64) {
	h.openTimes = u64.Sub(h.openTimes, toReduce)
}

func (h *HeroZhanJiang) IsPass(data *zhanjiang.ZhanJiangGuanQiaData) bool {
	_, exist := h.passGuanQia[data.Id]
	return exist
}

func (h *HeroZhanJiang) Pass(data *zhanjiang.ZhanJiangGuanQiaData) {
	h.passGuanQia[data.Id] = struct{}{}
}

func (h *HeroZhanJiang) StartChallenge(toSet *ZhanJiangChallenge) {
	h.curChallenge = toSet
	h.openTimes++
}

func (h *HeroZhanJiang) CurChallenge() *ZhanJiangChallenge {
	return h.curChallenge
}

func (h *HeroZhanJiang) EndChallenge() {
	h.curChallenge = nil
}

func (h *HeroZhanJiang) LastCaptainId() uint64 {
	return h.lastCaptainId
}

func (h *HeroZhanJiang) SetLastCaptainId(toSet uint64) {
	h.lastCaptainId = toSet
}

func (h *HeroZhanJiang) ResetDaily() {
	h.openTimes = 0
}

func (h *HeroZhanJiang) EncodeClient() *shared_proto.HeroZhanJiangProto {
	proto := &shared_proto.HeroZhanJiangProto{}

	proto.OpenTimes = u64.Int32(h.openTimes)
	if len(h.passGuanQia) > 0 {
		proto.PassGuanQia = make([]int32, 0, len(h.passGuanQia))
		for pass := range h.passGuanQia {
			proto.PassGuanQia = append(proto.PassGuanQia, u64.Int32(pass))
		}
	}

	if h.curChallenge != nil {
		proto.CurChallenge = h.curChallenge.EncodeClient()
	}

	return proto
}

func (h *HeroZhanJiang) unmarshal(proto *server_proto.HeroZhanJiangServerProto, datas *config.ConfigDatas, captainGetter func(id uint64) *Captain) {
	if proto == nil {
		return
	}

	h.openTimes = proto.OpenTimes
	h.lastCaptainId = proto.LastCaptainId

	for _, pass := range proto.PassGuanQia {
		guanQiaData := datas.GetZhanJiangGuanQiaData(pass)
		if guanQiaData == nil {
			logrus.Errorf("玩家上线时，玩家已经通关的关卡找不到了: %d", pass)
			continue
		}

		h.passGuanQia[guanQiaData.Id] = struct{}{}
	}

	if proto.CurChallenge != nil {
		guanQianData := datas.GetZhanJiangGuanQiaData(proto.CurChallenge.GuanQia)
		if guanQianData == nil {
			logrus.Errorf("玩家上线时，当前正在挑战的关卡数据不见了: %d", proto.CurChallenge.GuanQia)
		} else {
			var captainId uint64
			if proto.CurChallenge.CaptainId != 0 {
				if captain := captainGetter(proto.CurChallenge.CaptainId); captain == nil {
					logrus.Errorf("玩家上线时，设置的出战武将不见了：%d", proto.CurChallenge.CaptainId)
				} else {
					captainId = captain.Id()
				}
			}

			h.curChallenge = &ZhanJiangChallenge{
				guanQia:   guanQianData,
				passCount: proto.CurChallenge.PassCount,
				captainId: captainId,
			}
		}
	}
}

func (h *HeroZhanJiang) EncodeServer() *server_proto.HeroZhanJiangServerProto {
	proto := &server_proto.HeroZhanJiangServerProto{}

	proto.OpenTimes = h.openTimes
	if len(h.passGuanQia) > 0 {
		proto.PassGuanQia = make([]uint64, 0, len(h.passGuanQia))
		for pass := range h.passGuanQia {
			proto.PassGuanQia = append(proto.PassGuanQia, pass)
		}
	}

	if h.curChallenge != nil {
		proto.CurChallenge = h.curChallenge.EncodeServer()
	}

	proto.LastCaptainId = h.lastCaptainId

	return proto
}

func NewZhanJiangChallenge(guanQia *zhanjiang.ZhanJiangGuanQiaData, captainId uint64) *ZhanJiangChallenge {
	return &ZhanJiangChallenge{guanQia: guanQia, captainId: captainId}
}

// 当前挑战
type ZhanJiangChallenge struct {
	guanQia   *zhanjiang.ZhanJiangGuanQiaData // 当前开启的关卡
	passCount uint64                          // 通关的数量
	captainId uint64                          // 设置的武将的id
}

func (c *ZhanJiangChallenge) GuanQia() *zhanjiang.ZhanJiangGuanQiaData {
	return c.guanQia
}

func (c *ZhanJiangChallenge) PassCount() uint64 {
	return c.passCount
}

func (c *ZhanJiangChallenge) IncPassCount() {
	c.passCount++
}

func (c *ZhanJiangChallenge) IsAllPass() bool {
	return c.passCount >= uint64(len(c.guanQia.ZhanJiangDatas))
}

func (c *ZhanJiangChallenge) SetCaptainId(toSet uint64) {
	c.captainId = toSet
}

func (c *ZhanJiangChallenge) CaptainId() uint64 {
	return c.captainId
}

func (c *ZhanJiangChallenge) EncodeClient() *shared_proto.ZhanJiangChallengeProto {
	proto := &shared_proto.ZhanJiangChallengeProto{}

	proto.GuanQia = u64.Int32(c.guanQia.Id)
	proto.PassCount = u64.Int32(c.passCount)
	proto.CaptainId = u64.Int32(c.captainId)

	return proto
}

func (c *ZhanJiangChallenge) EncodeServer() *server_proto.ZhanJiangChallengeServerProto {
	proto := &server_proto.ZhanJiangChallengeServerProto{}

	proto.GuanQia = c.guanQia.Id
	proto.PassCount = c.passCount
	proto.CaptainId = c.captainId

	return proto
}
