package httphandler

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/config/goods"
)

func newGMResponse(errMsg string, resp interface{}) *GMResponse {
	return &GMResponse {
		RespCode: len(errMsg),
		ErrMsg:   errMsg,
		Resp:     resp,
	}
}

type GMResponse struct {
	RespCode int          `json:"resp_code"` // 错误码,0表示成功（OK）
	ErrMsg   string       `json:"err_msg,omitempty"`  // 错误信息，错误码不为0时才有内容
	Resp     interface{} `json:"resp,omitempty"` // 回复内容
}

// 英雄快照里面有些数据拿不到，所以直接从hero里面拿，但是
// 这必须放在英雄锁里面才能new，不允许放外面，
func newGMHeroInfoResponse(hero *entity.Hero) *GMHeroInfoResponse {
	resp := &GMHeroInfoResponse {
		Id:             hero.Id(),
		Name:           hero.Name(),
		Level:          hero.Level(),
		VipLevel:       hero.VipLevel(),
		Gold:           hero.GetGold(),
		Stone:          hero.GetStone(),
		Dianquan:       hero.GetDianquan(),
		Yinliang:       hero.GetYinliang(),
		CreateTime:     timeutil.Marshal64(hero.CreateTime()),
		LastLoginTime:  timeutil.Marshal64(hero.LastOnlineTime()),
		LastChargeTime: hero.Misc().LastChargeTime(),
		SystemYuanbao:  0,
		ChargeYuanbao:  hero.GetYuanbao(),
	}
	if oldNames := hero.OldName(); len(oldNames) > 0 {
		resp.OldName = oldNames[0]
	}
	return resp
}

// 角色查询回复
type GMHeroInfoResponse struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	OldName        string `json:"old_name,omitempty"`
	Level          uint64 `json:"level,omitempty"`
	VipLevel       uint64 `json:"vip_level,omitempty"`
	GuildName      string `json:"guild_name,omitempty"`
	GuildClassName string `json:"guild_class_name,omitempty"`
	Gold           uint64 `json:"gold,omitempty"`
	Stone          uint64 `json:"stone,omitempty"`
	Dianquan       uint64 `json:"dianquan,omitempty"`
	Yinliang       uint64 `json:"yinliang,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	LastLoginTime  int64  `json:"last_login_time,omitempty"`
	LastChargeTime int64  `json:"last_charge_time,omitempty"`
	SystemYuanbao  uint64 `json:"system_yuanbao,omitempty"`
	ChargeYuanbao  uint64 `json:"charge_yuanbao,omitempty"`
}

func (resp *GMHeroInfoResponse) setGuildInfo(guild *sharedguilddata.Guild) {
	resp.GuildName = guild.Name()
	if member := guild.GetMember(resp.Id); member != nil {
		resp.GuildClassName = member.ClassLevelData().Name
	}
}

// 角色仓库（物品）查询回复
type GMHeroDepotGoodsResponse struct {
	Goods     []*GMGoodsInfo `json:"goods,omitempty"`
}

// 角色仓库（装备）查询回复
type GMHeroDepotEquipResponse struct {
	Equips    []*GMEquipInfo `json:"equips,omitempty"`
}

// 武将查询回复
type GMCaptainInfoResponse struct {
	Captains []*GMCaptianInfo `json:"captains,omitempty"`
}

func newGMCaptianInfo(captain *entity.Captain) *GMCaptianInfo {
	return &GMCaptianInfo {
		Id:          captain.Id(),
		Name:        captain.Name(),
		Troop:       captain.GetTroopSequence(),
		Rarity:      captain.Data().Rarity.Name,
		Level:       captain.Level(),
		Ability:     captain.Ability(),
		Star:        captain.Star(),
		FightAmount: captain.FullSoldierFightAmount(),
		EquipCount:  captain.GetEquipmentCount(),
		GemCount:    captain.GetGemCount(),
	}
}

// 武将基础信息
type GMCaptianInfo struct {
	Id          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Troop       uint64 `json:"troop,omitempty"`
	Rarity      string `json:"rarity,omitempty"`
	Level       uint64 `json:"level,omitempty"`
	Ability     uint64 `json:"ability,omitempty"`
	Star        uint64 `json:"star,omitempty"`
	FightAmount uint64 `json:"fight_amount,omitempty"` // 战斗力
	EquipCount  uint64 `json:"equip_count,omitempty"` // 装备数
	GemCount    uint64 `json:"gem_count,omitempty"` // 宝石数
}

// 武将装备查询回复
type GMCaptainEquipInfoResponse struct {
	CaptainId uint64         `json:"captain_id,omitempty"`
	Equips    []*GMEquipInfo `json:"equips,omitempty"`
}

func newGMEquipInfo(equip *entity.Equipment) *GMEquipInfo {
	return &GMEquipInfo {
		Id:      equip.Data().Id,
		Name:    equip.Data().Name,
		Quality: u64.FromInt32(int32(equip.Data().Quality.GoodsQuality.Quality)),
		Level:   equip.Level(),
		Star:    equip.RefinedStar(),
	}
}

// 装备基础信息
type GMEquipInfo struct {
	Id      uint64 `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Quality uint64 `json:"quality,omitempty"`
	Level   uint64 `json:"level,omitempty"`
	Star    uint64 `json:"star,omitempty"`
}

// 武将宝石查询回复
type GMCaptainGemInfoResponse struct {
	CaptainId uint64       `json:"captain_id,omitempty"`
	Gems      []*GMGemInfo `json:"gems,omitempty"`
}

func newGMGemInfo(gem *goods.GemData) *GMGemInfo {
	return &GMGemInfo {
		Id:    gem.Id,
		Name:  gem.Name,
		Level: gem.Level,
	}
}

// 宝石基础信息
type GMGemInfo struct {
	Id      uint64 `json:"id,omitempty"`    // 宝石id
	Name    string `json:"name,omitempty"`  // 宝石名称
	Level   uint64 `json:"level,omitempty"` // 宝石等级
}

func newGMGoodsInfo(id, count uint64, config interface {
	GetGemData(uint64) *goods.GemData
	GetGoodsData(uint64) *goods.GoodsData
}) *GMGoodsInfo {
	p := &GMGoodsInfo {
		Id:    id,
		Count: count,
	}
	if goodsData := config.GetGoodsData(id); goodsData != nil {
		p.Name = goodsData.Name
	} else if gemData := config.GetGemData(id); gemData != nil {
		p.Name = gemData.Name
	}
	return p
}

// 物品基础信息
type GMGoodsInfo struct {
	Id      uint64 `json:"id,omitempty"`    // 物品id
	Name    string `json:"name,omitempty"`  // 物品名称
	Count   uint64 `json:"count,omitempty"` // 物品数量
}

// 发送礼品邮件
type GMSendPrizeMailRequest struct {
	Title   string                   `json:"title,omitempty"`
	Content string                   `json:"content,omitempty"`
	Prize   *shared_proto.PrizeProto `json:"prize,omitempty"`
}

func checkGMSendPrizeMailRequest(req *GMSendPrizeMailRequest, datas iface.ConfigDatas) error {
	if len(req.Title) <= 0 {
		return errors.Errorf("no title")
	}
	if len(req.Content) <= 0 {
		return errors.Errorf("no content")
	}
	prize := req.Prize
	if prize.Gold < 0 || prize.SafeGold < 0 {
		return errors.Errorf("invalid gold")
	}
	if prize.Food < 0 || prize.SafeFood < 0 {
		return errors.Errorf("invalid food")
	}
	if prize.Wood < 0 || prize.SafeWood < 0 {
		return errors.Errorf("invalid wood")
	}
	if prize.Stone < 0 || prize.SafeStone < 0 {
		return errors.Errorf("invalid stone")
	}
	if prize.HeroExp < 0 {
		return errors.Errorf("invalid heroExp")
	}
	if prize.CaptainExp < 0 {
		return errors.Errorf("invalid captainExp")
	}
	if prize.Prosperity < 0 {
		return errors.Errorf("invalid prosperity")
	}
	if prize.Yuanbao < 0 {
		return errors.Errorf("invalid yuanbao")
	}
	if prize.Dianquan < 0 {
		return errors.Errorf("invalid dianquan")
	}
	if prize.GuildContributionCoin < 0 {
		return errors.Errorf("invalid guildContributionCoin")
	}
	if prize.Jade < 0 {
		return errors.Errorf("invalid jade")
	}
	if prize.JadeOre < 0 {
		return errors.Errorf("invalid jadeOre")
	}
	if prize.Sp < 0 {
		return errors.Errorf("invalid sp")
	}

	if len(prize.GoodsId) != len(prize.GoodsCount) {
		return errors.Errorf("wrong goods array length")
	}
	for _, count := range prize.GoodsCount {
		if count <= 0 {
			return errors.Errorf("invalid goodsCount")
		}
	}
	for _, id := range prize.GoodsId {
		if id <= 0 || datas.GetGoodsData(u64.FromInt32(id)) == nil {
			return errors.Errorf("invalid goodsId")
		}
	}

	if len(prize.EquipmentId) != len(prize.EquipmentCount) {
		return errors.Errorf("wrong equip array length")
	}
	for _, count := range prize.EquipmentCount {
		if count <= 0 {
			return errors.Errorf("invalid equipCount")
		}
	}
	for _, id := range prize.EquipmentId {
		if id <= 0 || datas.GetEquipmentData(u64.FromInt32(id)) == nil {
			return errors.Errorf("invalid equipId")
		}
	}

	if len(prize.GemId) != len(prize.GemCount) {
		return errors.Errorf("wrong gem array length")
	}
	for _, count := range prize.GemCount {
		if count <= 0 {
			return errors.Errorf("invalid gemCount")
		}
	}
	for _, id := range prize.GemId {
		if id <= 0 || datas.GetGemData(u64.FromInt32(id)) == nil {
			return errors.Errorf("invalid gemId")
		}
	}

	if len(prize.BaowuId) != len(prize.BaowuCount) {
		return errors.Errorf("wrong baowu array length")
	}
	for _, count := range prize.BaowuCount {
		if count <= 0 {
			return errors.Errorf("invalid baowuCount")
		}
	}
	for _, id := range prize.BaowuId {
		if id <= 0 || datas.GetBaowuData(u64.FromInt32(id)) == nil {
			return errors.Errorf("invalid baowuId")
		}
	}

	if len(prize.CaptainId) != len(prize.CaptainCount) {
		return errors.Errorf("wrong captain array length")
	}
	for _, count := range prize.CaptainCount {
		if count <= 0 {
			return errors.Errorf("invalid captainCount")
		}
	}
	for _, id := range prize.CaptainId {
		if id <= 0 || datas.GetResCaptainData(u64.FromInt32(id)) == nil {
			return errors.Errorf("invalid captainId")
		}
	}

	resdata.SetPrizeProtoIsNotEmpty(prize)

	return nil
}

func (req *GMSendPrizeMailRequest) encodeMail(datas iface.ConfigDatas) *shared_proto.MailProto {
	proto := &shared_proto.MailProto{}
	proto.Title = req.Title
	proto.Text = req.Content
	proto.MailType = shared_proto.MailType_MailNormal
	proto.Prize = req.Prize

	return proto
}
