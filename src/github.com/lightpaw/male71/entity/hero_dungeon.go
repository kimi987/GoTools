package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/recovtimes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func newHeroDungeon(initData *heroinit.HeroInitData, ctime time.Time) *HeroDungeon {
	h := &HeroDungeon{}
	h.collectedChapterStarPrizes = make(map[uint64]map[int]struct{})
	h.dungeonStars = make(map[uint64][]bool)
	h.dungeonPassLimits = make(map[uint64]uint64)
	h.passDungeons = make(map[uint64]struct{})
	h.collectedPassDungeonPrizes = make(map[uint64]struct{})
	h.collectedChapterPrizes = make(map[uint64]struct{})
	h.challengeTimes = NewRecoverableTimes(ctime, initData.DungeonChallengeRecoverDuration, initData.DungeonChallengeMaxTimes)
	h.challengeTimes.SetTimes(initData.DungeonChallengeDefaultTimes, ctime)
	h.vipAddedPassLimit = make(map[uint64]uint64)
	h.vipAddedPassBoughtTimes = make(map[uint64]uint64)

	return h
}

// 玩家副本
type HeroDungeon struct {
	// 记录章节已经领取的星数宝箱<chapterId, 宝箱下标>
	collectedChapterStarPrizes map[uint64]map[int]struct{}

	// 带有星数的副本的完成星数以及点亮位置
	dungeonStars map[uint64][]bool

	// 记录带有每日通过次数限制的副本次数
	dungeonPassLimits map[uint64]uint64

	// 通关了的副本
	passDungeons map[uint64]struct{}

	// 领取了的副本奖励
	collectedPassDungeonPrizes map[uint64]struct{}

	// 领取了的章节奖励
	collectedChapterPrizes map[uint64]struct{}

	// 挑战次数
	challengeTimes *recovtimes.RecoverTimes

	// 每天vip特权购买的挑战次数
	vipAddedPassLimit map[uint64]uint64

	// 每天vip特权购买了几次
	vipAddedPassBoughtTimes map[uint64]uint64

}

func (h *HeroDungeon) AddVipAddedPassLimit(dungeonId, toAdd uint64) uint64 {
	h.vipAddedPassLimit[dungeonId] += toAdd
	h.vipAddedPassBoughtTimes[dungeonId]++
	return h.vipAddedPassLimit[dungeonId]
}

func (h *HeroDungeon) VipAddedPassBoughtTimes(dungeonId uint64) uint64 {
	return h.vipAddedPassBoughtTimes[dungeonId]
}

func (h *HeroDungeon) AddLimitDungeonPassTimes(id uint64, times uint64) {
	h.dungeonPassLimits[id] += times
}

func (h *HeroDungeon) GetLimitDungeonPassTimes(id uint64) uint64 {
	return h.dungeonPassLimits[id]
}

func (h *HeroDungeon) IsCanAutoPass(data *dungeon.DungeonData) bool {
	if data.Star > 0 {
		return h.GetDungeonStar(data.Id) >= data.Star
	}
	//return h.IsPass(data) // 原来的module代码已经有判定
	return true
}

func (h *HeroDungeon) IsPassLimit(data *dungeon.DungeonData, challengeTimes uint64) bool {
	if data.PassLimit > 0 {
		if times, ok := h.dungeonPassLimits[data.Id]; ok && times+challengeTimes > data.PassLimit+h.vipAddedPassLimit[data.Id] {
			return true
		}
	}
	return false
}

// 是否已经解锁全部的前置关卡
func (h *HeroDungeon) IsUnlockPreDungeon(data *dungeon.DungeonData) bool {
	if len(data.UnlockPassDungeon) <= 0 {
		return true
	}
	for _, dungeonData := range data.UnlockPassDungeon {
		if !h.IsPass(dungeonData) {
			return false
		}
	}
	return true
}

func (h *HeroDungeon) IsPass(data *dungeon.DungeonData) bool {
	_, ok := h.passDungeons[data.Id]
	return ok
}

func (h *HeroDungeon) Pass(data *dungeon.DungeonData, enabledStars []bool) (starRefreshed bool) {
	h.passDungeons[data.Id] = struct{}{}
	// 带有次数限制的副本需要记录当日已经消耗的次数
	if data.PassLimit > 0 {
		h.dungeonPassLimits[data.Id]++
	}
	// 带星副本才需要记录
	if data.Star > 0 {
		passStar := util.TrueCount(enabledStars)
		if passStar <= 0 {
			return
		}
		oldStar := h.GetDungeonStar(data.Id)
		if passStar < oldStar {
			return
		}
		h.dungeonStars[data.Id] = enabledStars
		starRefreshed = true
	}
	return
}

func (h *HeroDungeon) GetDungeonStar(id uint64) uint64 {
	enabledStars, ok := h.dungeonStars[id]
	if !ok {
		return 0
	}
	return util.TrueCount(enabledStars)
}

func (h *HeroDungeon) GetChapterStar(dungeunIds []uint64) (star uint64) {
	for _, id := range dungeunIds {
		star += h.GetDungeonStar(id)
	}
	return
}

func (h *HeroDungeon) TryCollectChapterStarPrize(id uint64, prize int) bool {
	if prizes, ok := h.collectedChapterStarPrizes[id]; ok {
		if _, ok := prizes[prize]; ok {
			return false
		} else {
			prizes[prize] = struct{}{}
		}
	} else {
		prizes = make(map[int]struct{})
		prizes[prize] = struct{}{}
		h.collectedChapterStarPrizes[id] = prizes
	}
	return true
}

// GM命令重置通关的副本
func (h *HeroDungeon) GMReset() {
	h.collectedChapterStarPrizes = make(map[uint64]map[int]struct{})
	h.dungeonStars = make(map[uint64][]bool)
	h.dungeonPassLimits = make(map[uint64]uint64)
	h.passDungeons = make(map[uint64]struct{})
	h.collectedChapterPrizes = make(map[uint64]struct{})
	h.collectedPassDungeonPrizes = make(map[uint64]struct{})
}

func (h *HeroDungeon) IsCollectChapterPrize(data *dungeon.DungeonChapterData) bool {
	_, ok := h.collectedChapterPrizes[data.Id]
	return ok
}

func (h *HeroDungeon) CollectChapterPrize(data *dungeon.DungeonChapterData) {
	h.collectedChapterPrizes[data.Id] = struct{}{}
}

func (h *HeroDungeon) IsCollectPassDungeonPrize(data *dungeon.DungeonData) bool {
	_, ok := h.collectedPassDungeonPrizes[data.Id]
	return ok
}

func (h *HeroDungeon) CollectPassDungeonPrize(data *dungeon.DungeonData) {
	h.collectedPassDungeonPrizes[data.Id] = struct{}{}
}

func (h *HeroDungeon) ChallengeTimes() *recovtimes.RecoverTimes {
	return h.challengeTimes
}

func (h *HeroDungeon) ResetDaily() {
	h.dungeonPassLimits = make(map[uint64]uint64)
	h.vipAddedPassLimit = make(map[uint64]uint64)
	h.vipAddedPassBoughtTimes = make(map[uint64]uint64)
}

func (h *HeroDungeon) EncodeClient(datas config.Configs) *shared_proto.HeroDungeonProto {
	proto := &shared_proto.HeroDungeonProto{}

	if len(h.collectedChapterStarPrizes) > 0 {
		proto.ChapterStarPrizes = make([]*shared_proto.CollectedChapterStarPrizes, 0, len(h.collectedChapterStarPrizes))
		for chapterId, prizeIndexs := range h.collectedChapterStarPrizes {
			p := &shared_proto.CollectedChapterStarPrizes{}
			p.Chapter = u64.Int32(chapterId)
			for index := range prizeIndexs {
				p.PrizeIndexs = append(p.PrizeIndexs, int32(index))
			}
			proto.ChapterStarPrizes = append(proto.ChapterStarPrizes, p)
		}
	}

	if len(h.dungeonStars) > 0 {
		proto.DungeonStars = make([]*shared_proto.DungeonStar, 0, len(h.dungeonStars))
		for dungeonId, stars := range h.dungeonStars {
			proto.DungeonStars = append(proto.DungeonStars, &shared_proto.DungeonStar{
				Dungeon:     u64.Int32(dungeonId),
				EnabledStar: stars,
			})
		}
	}

	if len(h.dungeonPassLimits) > 0 {
		proto.DungeonLimits = make([]*shared_proto.DungeonPassLimit, 0, len(h.dungeonPassLimits))
		for dungeonId, times := range h.dungeonPassLimits {
			proto.DungeonLimits = append(proto.DungeonLimits, &shared_proto.DungeonPassLimit{
				Dungeon: u64.Int32(dungeonId),
				Times:   u64.Int32(times),
			})
		}
	}

	if len(h.passDungeons) > 0 {
		proto.PassDungeons = make([]int32, 0, len(h.passDungeons))
		for id := range h.passDungeons {
			proto.PassDungeons = append(proto.PassDungeons, u64.Int32(id))
		}
	}

	if len(h.collectedPassDungeonPrizes) > 0 {
		proto.CollectedPassDungeonPrizes = make([]int32, 0, len(h.collectedPassDungeonPrizes))
		for id := range h.collectedPassDungeonPrizes {
			proto.CollectedPassDungeonPrizes = append(proto.CollectedPassDungeonPrizes, u64.Int32(id))
		}
	}

	if len(h.collectedChapterPrizes) > 0 {
		proto.CollectedChapterPrizes = make([]int32, 0, len(h.collectedChapterPrizes))
		for id := range h.collectedChapterPrizes {
			proto.CollectedChapterPrizes = append(proto.CollectedChapterPrizes, u64.Int32(id))
		}
	}

	for _, data := range datas.GetDungeonChapterDataArray() {
		star := h.GetChapterStar(data.GetStarDungeonIds())
		if star <= 0 {
			continue
		}

		proto.ChapterStars = append(proto.ChapterStars, &shared_proto.DungeonChapterStar{
			Chapter: u64.Int32(data.Id),
			Star:    u64.Int32(star),
		})
	}

	proto.AutoRecoverStartTime = h.challengeTimes.StartTimeUnix32()

	for k, v := range h.vipAddedPassBoughtTimes {
		vp := &shared_proto.HeroVipBoughtTimes{DungeonId:u64.Int32(k), BoughtTimes: u64.Int32(v)}
		proto.VipBoughtTimes = append(proto.VipBoughtTimes, vp)
	}

	return proto
}

func (h *HeroDungeon) unmarshal(proto *server_proto.HeroDungeonServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	for _, v := range proto.GetChapterStarPrizes() {
		data := datas.GetDungeonChapterData(u64.FromInt32(v.Chapter))
		if data == nil {
			logrus.WithField("id", v.Chapter).Debugln("HeroDungeon 中章节没找到!")
			continue
		}
		indexs := make(map[int]struct{})
		for _, index := range v.PrizeIndexs {
			indexs[int(index)] = struct{}{}
		}
		h.collectedChapterStarPrizes[data.Id] = indexs
	}

	for _, v := range proto.GetDungeonStars() {
		data := datas.GetDungeonData(u64.FromInt32(v.Dungeon))
		if data == nil {
			logrus.WithField("id", v.Dungeon).Debugln("HeroDungeon 中副本没找到!")
			continue
		}
		h.dungeonStars[data.Id] = v.EnabledStar
	}

	for _, v := range proto.GetDungeonLimits() {
		data := datas.GetDungeonData(u64.FromInt32(v.Dungeon))
		if data == nil {
			logrus.WithField("id", v.Dungeon).Debugln("HeroDungeon 中副本没找到!")
			continue
		}
		h.dungeonPassLimits[data.Id] = u64.FromInt32(v.Times)
	}

	for _, id := range proto.GetPassDungeons() {
		data := datas.GetDungeonData(id)
		if data == nil {
			logrus.WithField("id", id).Debugln("HeroDungeon 中通关的副本没找到!")
			continue
		}

		h.passDungeons[data.Id] = struct{}{}
	}

	for _, id := range proto.GetCollectedPassDungeonPrizes() {
		h.collectedPassDungeonPrizes[id] = struct{}{}
	}

	for _, id := range proto.GetCollectedChapterPrizes() {
		data := datas.GetDungeonChapterData(id)
		if data == nil {
			logrus.WithField("id", id).Debugln("HeroDungeon 中领取的通关难度奖励没找到!")
			continue
		}

		h.collectedChapterPrizes[data.Id] = struct{}{}
	}

	h.challengeTimes.SetStartTime(timeutil.Unix64(proto.GetAutoRecoverStartTime()))

	if proto.VipAddPassLimit != nil {
		h.vipAddedPassLimit = proto.VipAddPassLimit
	}

	if proto.VipAddPassBoughtTimes != nil {
		h.vipAddedPassBoughtTimes = proto.VipAddPassBoughtTimes
	}
}

func (h *HeroDungeon) EncodeServer() *server_proto.HeroDungeonServerProto {
	proto := &server_proto.HeroDungeonServerProto{}

	if len(h.collectedChapterStarPrizes) > 0 {
		proto.ChapterStarPrizes = make([]*shared_proto.CollectedChapterStarPrizes, 0, len(h.collectedChapterStarPrizes))
		for chapterId, prizeIndexs := range h.collectedChapterStarPrizes {
			p := &shared_proto.CollectedChapterStarPrizes{}
			p.Chapter = u64.Int32(chapterId)
			for index := range prizeIndexs {
				p.PrizeIndexs = append(p.PrizeIndexs, int32(index))
			}
			proto.ChapterStarPrizes = append(proto.ChapterStarPrizes, p)
		}
	}

	if len(h.dungeonStars) > 0 {
		proto.DungeonStars = make([]*shared_proto.DungeonStar, 0, len(h.dungeonStars))
		for dungeonId, stars := range h.dungeonStars {
			proto.DungeonStars = append(proto.DungeonStars, &shared_proto.DungeonStar{
				Dungeon:     u64.Int32(dungeonId),
				EnabledStar: stars,
			})
		}
	}

	if len(h.dungeonPassLimits) > 0 {
		proto.DungeonLimits = make([]*shared_proto.DungeonPassLimit, 0, len(h.dungeonPassLimits))
		for dungeonId, times := range h.dungeonPassLimits {
			proto.DungeonLimits = append(proto.DungeonLimits, &shared_proto.DungeonPassLimit{
				Dungeon: u64.Int32(dungeonId),
				Times:   u64.Int32(times),
			})
		}
	}

	if len(h.passDungeons) > 0 {
		proto.PassDungeons = make([]uint64, 0, len(h.passDungeons))
		for id := range h.passDungeons {
			proto.PassDungeons = append(proto.PassDungeons, id)
		}
	}

	if len(h.collectedPassDungeonPrizes) > 0 {
		proto.CollectedPassDungeonPrizes = make([]uint64, 0, len(h.collectedPassDungeonPrizes))
		for id := range h.collectedPassDungeonPrizes {
			proto.CollectedPassDungeonPrizes = append(proto.CollectedPassDungeonPrizes, id)
		}
	}

	if len(h.collectedChapterPrizes) > 0 {
		proto.CollectedChapterPrizes = make([]uint64, 0, len(h.collectedChapterPrizes))
		for id := range h.collectedChapterPrizes {
			proto.CollectedChapterPrizes = append(proto.CollectedChapterPrizes, id)
		}
	}

	proto.AutoRecoverStartTime = h.challengeTimes.StartTimeUnix64()
	proto.VipAddPassLimit = h.vipAddedPassLimit
	proto.VipAddPassBoughtTimes = h.vipAddedPassBoughtTimes

	return proto
}
