package relation

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/relation"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"github.com/lightpaw/male7/entity/npcid"
)

func NewRelationModule(dep iface.ServiceDep, chat iface.ChatService, tssClient iface.TssClient, recommendHeroCache *LocationHeroCache) *RelationModule {

	m := &RelationModule{
		dep:                 dep,
		heroSnapshotService: dep.HeroSnapshot(),
		chat:                chat,
		tssClient:           tssClient,
		recommendHeroCache:  recommendHeroCache,
	}
	d, err := time.ParseDuration("1m")
	if err != nil {
		logrus.WithError(err).Panicf("刷新推荐好友列表间隔错误")
	}
	m.nextRefreshMsgDuration = d

	return m
}

//gogen:iface
type RelationModule struct {
	dep                 iface.ServiceDep
	heroSnapshotService iface.HeroSnapshotService
	chat                iface.ChatService
	tssClient           iface.TssClient

	recommendHeroCache *LocationHeroCache

	nextRefreshMsgDuration time.Duration
}

//gogen:iface
func (m *RelationModule) ProcessAddRelation(proto *relation.C2SAddRelationProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok || targetId <= 0 {
		logrus.Debug("添加好友黑名单，目标id无效")
		hc.Send(relation.ERR_ADD_RELATION_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("添加好友黑名单，自己的id")
		hc.Send(relation.ERR_ADD_RELATION_FAIL_SELF_ID)
		return
	}

	targetSnapshot := m.heroSnapshotService.Get(targetId)
	if targetSnapshot == nil {
		logrus.Warn("添加好友黑名单，目标id玩家不存在")
		hc.Send(relation.ERR_ADD_RELATION_FAIL_INVALID_ID)
		return
	}

	ctime := m.dep.Time().CurrentTime()

	var addFriendSucc bool
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		var snsType uint64
		rt := hero.Relation().GetRelation(targetId)
		if proto.Friend {
			if hero.Relation().RelationCount(entity.Friend) >= m.dep.Datas().MiscGenConfig().FriendMaxCount {
				result.Add(relation.ERR_ADD_RELATION_FAIL_LIMIT)
				return
			}

			if rt == entity.Friend {
				logrus.Debug("添加好友黑名单，目标已经是好友")
				result.Add(relation.ERR_ADD_RELATION_FAIL_FRIEND)
				return
			}

			hero.Relation().AddFriend(targetId, ctime)
			addFriendSucc = true

			snsType = operate_type.SNSAddFriend
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FRIEND_AMOUNT)
		} else {
			if hero.Relation().RelationCount(entity.Black) >= m.dep.Datas().MiscGenConfig().BlackMaxCount {
				result.Add(relation.ERR_ADD_RELATION_FAIL_LIMIT)
				return
			}

			if rt == entity.Black {
				logrus.Debug("添加好友黑名单，目标已经是黑名单")
				result.Add(relation.ERR_ADD_RELATION_FAIL_BLACK)
				return
			}

			hero.Relation().AddBlack(targetId, ctime)

			snsType = operate_type.SNSAddBlack
		}

		result.Add(relation.NewS2cAddRelationMarshalMsg(proto.Friend, proto.Id, targetSnapshot.EncodeClient(), timeutil.Marshal32(ctime)))

		m.dep.Tlog().TlogSnsFlow(hero, snsType, u64.FromInt64(targetId))
		result.Changed()
		result.Ok()
	}) {
		return
	}

	if addFriendSucc {
		targetName := heromodule.GetFlagHeroName(m.dep.Datas(), targetSnapshot)
		text := m.dep.Datas().TextHelp().SysChatFriendAdded.New().WithHeroName(targetName).JsonString()
		m.chat.SysChat(targetId, hc.Id(), shared_proto.ChatType_ChatPrivate, text, shared_proto.ChatMsgType_ChatMsgFriendAdded, false, true, true, false)
	}

}

//gogen:iface
func (m *RelationModule) ProcessRemoveRelation(proto *relation.C2SRemoveRelationProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Id)
	if !ok || targetId <= 0 {
		logrus.Debug("移除好友黑名单，目标id无效")
		hc.Send(relation.ERR_REMOVE_RELATION_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("移除好友黑名单，自己的id")
		hc.Send(relation.ERR_REMOVE_RELATION_FAIL_SELF_ID)
		return
	}

	var relationType entity.RelationType
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		relationType = hero.Relation().RemoveRelation(targetId)
		result.Add(relation.NewS2cRemoveRelationMsg(proto.Id, int32(relationType)))

		var snsType uint64
		if relationType == entity.Friend {
			snsType = operate_type.SNSDelFriend
		} else if relationType == entity.Black {
			snsType = operate_type.SNSDelBlack
		}
		m.dep.Tlog().TlogSnsFlow(hero, snsType, u64.FromInt64(targetId))

		result.Changed()
		result.Ok()
	}) {
		return
	}
}

//gogen:iface
func (m *RelationModule) ProcessRemoveEnemy(proto *relation.C2SRemoveEnemyProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Id)
	if !ok || targetId <= 0 {
		logrus.Debug("移除仇人，目标id无效")
		hc.Send(relation.ERR_REMOVE_ENEMY_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("移除仇人，自己的id")
		hc.Send(relation.ERR_REMOVE_ENEMY_FAIL_SELF_ID)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Relation().RemoveEnemy(targetId)
		result.Add(relation.NewS2cRemoveEnemyMsg(proto.Id))

		result.Changed()
		result.Ok()

		m.dep.Tlog().TlogSnsFlow(hero, operate_type.SNSDelEnemy, u64.FromInt64(targetId))
	}) {
		return
	}

}

// 作废，用 new_list_relation
//gogen:iface
func (m *RelationModule) ProcessListRelation(proto *relation.C2SListRelationProto, hc iface.HeroController) {

	var targetIds []int64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		targetIds = hero.Relation().RelationAndEnemyIds()
		return false
	})

	toSendProto := &relation.S2CListRelationProto{}
	toSendProto.Version = proto.Version

	for _, t := range targetIds {
		target := m.heroSnapshotService.Get(t)
		if target != nil {
			if proto.Version < target.SnapshotCreateTime() {
				toSendProto.Version = i32.Max(toSendProto.Version, target.SnapshotCreateTime())

				toSendProto.Proto = append(toSendProto.Proto, target.EncodeBasic4ClientBytes())
			}
		}
	}

	hc.Send(relation.NewS2cListRelationProtoMsg(toSendProto))
}

//gogen:iface
func (m *RelationModule) ProcessNewListRelation(proto *relation.C2SNewListRelationProto, hc iface.HeroController) {

	var targetIds []int64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		targetIds = hero.Relation().RelationAndEnemyIds()
		return false
	})

	toSendProto := &relation.S2CNewListRelationProto{}
	toSendProto.Version = proto.Version

	for _, t := range targetIds {
		target := m.heroSnapshotService.Get(t)
		if target != nil {
			if proto.Version < target.SnapshotCreateTime() {
				toSendProto.Version = i32.Max(toSendProto.Version, target.SnapshotCreateTime())

				toSendProto.Proto = append(toSendProto.Proto, target.EncodeClientBytes())
			}
		}
	}

	hc.Send(relation.NewS2cNewListRelationProtoMsg(toSendProto))
}

//gogen:iface
func (m *RelationModule) ProcessRecommendHeroList(proto *relation.C2SRecommendHeroListProto, hc iface.HeroController) {
	ctime := m.dep.Time().CurrentTime()
	// check cd
	if ctime.Before(hc.NextRefreshRecommendHeroTime()) {
		hc.Send(relation.ERR_RECOMMEND_HERO_LIST_FAIL_IN_CD)
		return
	}
	hc.UpdateNextRefreshRecommendHeroTime(ctime.Add(m.dep.Datas().MiscConfig().RefreshRecommendHeroDuration))

	var heroFriends []int64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		heroFriends = hero.Relation().FriendIds()
		return false
	})

	loc := u64.FromInt32(proto.Loc)
	if !proto.NeedLoc {
		loc = 0
	}

	// 尝试加载数据
	m.tryLoadRecommendHeroCache(proto.NeedLoc, loc)

	countPerPage := m.dep.Datas().MiscConfig().RefreshRecommendHeroPageSize

	startIndex := int(proto.Page)
	if startIndex <= 0 {
		startIndex = rand.Intn(u64.Int(countPerPage*10) + 1)
	}

	toSendProto := &relation.S2CRecommendHeroListProto{}

	m.recommendHeroCache.rangeHeros(loc, startIndex, func(i int, o *HeroObject) bool {

		if o.heroId == hc.Id() {
			return true
		}

		if i64.Contains(heroFriends, o.heroId) {
			// 已经是好友，跳过
			return true
		}

		heroProto := o.heroProto
		if target := m.heroSnapshotService.GetFromCache(o.heroId); target != nil {
			heroProto = target.EncodeClient()
		}

		if loc == 0 || u64.FromInt32(heroProto.Basic.Location) == loc {
			toSendProto.Page = int32(startIndex + i)
			toSendProto.Heros = append(toSendProto.Heros, heroProto)

			if uint64(len(toSendProto.Heros)) >= countPerPage {
				return false
			}
		}

		return true
	})

	hc.Send(relation.NewS2cRecommendHeroListProtoMsg(toSendProto))
}

func (m *RelationModule) tryLoadRecommendHeroCache(needLoc bool, loc uint64) {
	if m.recommendHeroCache.isLoaded(loc) {
		return
	}

	minLvl := m.dep.Datas().MiscConfig().RefreshRecommendHeroMinLevel
	pageSize := u64.Int(m.dep.Datas().MiscConfig().RefreshRecommendHeroPageSize)
	pageCount := u64.Int(m.dep.Datas().MiscConfig().RefreshRecommendHeroPageCount)
	allSize := pageSize * pageCount

	var heros []*entity.Hero
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		heros, err = m.dep.Db().LoadRecommendHeros(ctx, needLoc, loc, minLvl, 0, u64.FromInt(allSize), 0)
		return
	})
	if err != nil {
		logrus.WithError(err).Debugf("查询推荐好友列表错误")
		return
	}

	var heroProtos []*shared_proto.HeroBasicSnapshotProto
	for _, h := range heros {
		if h == nil {
			continue
		}
		hss := m.dep.HeroSnapshot().GetFromCache(h.Id())
		if hss == nil {
			hss = m.dep.HeroSnapshot().NewSnapshot(h)
		}

		heroProtos = append(heroProtos, hss.EncodeClient())
	}

	m.recommendHeroCache.updateLocationHeros(loc, heroProtos)
}

//gogen:iface
func (m *RelationModule) ProcessSearchHeros(proto *relation.C2SSearchHerosProto, hc iface.HeroController) {
	name := strings.TrimSpace(proto.Name)
	if len(name) <= 0 {
		logrus.Debugf("模糊搜索玩家，关键字错误 name: %v", name)
		hc.Send(relation.ERR_SEARCH_HEROS_FAIL_INVALID_ARG)
		return
	}

	if proto.Page < 0 {
		logrus.Debugf("模糊搜索玩家，页数错误 page: %v", proto.Page)
		hc.Send(relation.ERR_SEARCH_HEROS_FAIL_INVALID_ARG)
		return
	}

	ctime := m.dep.Time().CurrentTime()
	if ctime.Before(hc.NextSearchHeroTime()) {
		logrus.Debugf("模糊搜索玩家，请求太频繁")
		hc.Send(relation.ERR_SEARCH_HEROS_FAIL_IN_CD)
		return
	}
	hc.UpdateNextSearchHeroTime(ctime.Add(m.dep.Datas().MiscConfig().SearchHeroDuration))

	if !m.tssClient.TryCheckName("模糊搜索玩家", hc, name, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, relation.ERR_SEARCH_HEROS_FAIL_SERVER_ERROR) {
		return
	}

	page := u64.FromInt32(proto.Page)
	size := m.dep.Datas().MiscConfig().RefreshRecommendHeroPageSize
	startIndex := page * size

	var heros []*entity.Hero
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		heros, err = m.dep.Db().LoadHerosByName(ctx, name, startIndex, size)
		return
	}); err != nil {
		logrus.WithError(err).Debugf("模糊搜索玩家，页数错误 page: %v", proto.Page)
		hc.Send(relation.ERR_SEARCH_HEROS_FAIL_SERVER_ERROR)
		return
	}

	msgProto := &relation.S2CSearchHerosProto{}
	// 先根据 id 搜索
	if heroId, err := strconv.ParseInt(name, 10, 64); err == nil && !npcid.IsNpcId(heroId) {
		heroSnapshot := m.heroSnapshotService.Get(heroId)
		if heroSnapshot != nil {
			msgProto.Heros = append(msgProto.Heros, heroSnapshot.EncodeClient())
		}
	}

	for _, hero := range heros {
		if hero == nil {
			logrus.Debugf("模糊搜索玩家，hero == nil")
			continue
		}

		heroSnapshot := m.heroSnapshotService.GetFromCache(hero.Id())
		if heroSnapshot == nil {
			heroSnapshot = m.heroSnapshotService.NewSnapshot(hero)
		}
		msgProto.Heros = append(msgProto.Heros, heroSnapshot.EncodeClient())
	}

	hc.Send(relation.NewS2cSearchHerosProtoMsg(msgProto))
}

//gogen:iface
func (m *RelationModule) ProcessSearchHeroById(proto *relation.C2SSearchHeroByIdProto, hc iface.HeroController) {
	heroId, ok := idbytes.ToId(proto.HeroId)
	if !ok || heroId <= 0 {
		hc.Send(relation.ERR_SEARCH_HERO_BY_ID_FAIL_INVALID_ARG)
		return
	}

	ctime := m.dep.Time().CurrentTime()
	if ctime.Before(hc.NextSearchHeroTime()) {
		hc.Send(relation.ERR_SEARCH_HERO_BY_ID_FAIL_IN_CD)
		return
	}
	hc.UpdateNextSearchHeroTime(ctime.Add(m.dep.Datas().MiscConfig().SearchHeroDuration))

	var hero *entity.Hero
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		hero, err = m.dep.Db().LoadHero(ctx, heroId)
		return
	})
	if err != nil {
		logrus.WithError(err).Debugf("搜索玩家，db错误")
		hc.Send(relation.ERR_SEARCH_HERO_BY_ID_FAIL_SERVER_ERROR)
		return
	}
	if hero == nil {
		hc.Send(relation.ERR_SEARCH_HERO_BY_ID_FAIL_NO_HERO)
		return
	}

	heroSnapshot := m.heroSnapshotService.GetFromCache(hero.Id())
	if heroSnapshot == nil {
		heroSnapshot = m.heroSnapshotService.NewSnapshot(hero)
	}

	hc.Send(relation.NewS2cSearchHeroByIdMsg(heroSnapshot.EncodeClient()))
}

////gogen:iface
//func (m *RelationModule) ProcessSearchRelation(proto *relation.C2SSearchRelationProto, hc iface.HeroController)  {
//
//	ctime := m.dep.Time().CurrentTime()
//	if ctime.Before(hc.NextSearchHeroTime()) {
//		hc.Send(relation.ERR_SEARCH_RELATION_FAIL_IN_CD)
//		return
//	}
//	hc.UpdateNextSearchHeroTime(ctime.Add(m.dep.Datas().MiscConfig().SearchHeroDuration))
//
//	name := strings.TrimSpace(proto.Name)
//	if len(name) <= 0 {
//		hc.Send(relation.ERR_SEARCH_RELATION_FAIL_EMPTY_NAME)
//		return
//	}
//	var ids []int64
//	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
//		ids = hero.Relation().RelationIds()
//		return false
//	})
//	sendProto :=  &relation.S2CSearchRelationProto{}
//	for _, id := range ids {
//		if heroSnapshot := m.heroSnapshotService.Get(id); heroSnapshot != nil && strings.Contains(heroSnapshot.Name, name){
//			sendProto.Heros = append(sendProto.Heros, heroSnapshot.EncodeClient())
//		}
//	}
//	if len(sendProto.Heros) <= 0 {
//		hc.Send(relation.ERR_SEARCH_RELATION_FAIL_EMPTY_LIST)
//		return
//	}
//	hc.Send(relation.NewS2cSearchRelationProtoMsg(sendProto))
//}

//gogen:iface
func (m *RelationModule) ProcessSetImportantFriend(proto *relation.C2SSetImportantFriendProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Id)
	if !ok || targetId <= 0 {
		logrus.Debug("设置星标好友，目标id无效")
		hc.Send(relation.ERR_SET_IMPORTANT_FRIEND_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("设置星标好友，自己的id")
		hc.Send(relation.ERR_SET_IMPORTANT_FRIEND_FAIL_SELF_ID)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		relationType := hero.Relation().GetRelation(targetId)
		if relationType != entity.Friend {
			logrus.Debug("设置星标好友，不是自己好友")
			result.Add(relation.ERR_SET_IMPORTANT_FRIEND_FAIL_NOT_FRIEND)
			return
		}
		ctime := m.dep.Time().CurrentTime()
		if !hero.Relation().SetImportantFriend(targetId, ctime) {
			logrus.Debug("设置星标好友，已经是星标好友")
			result.Add(relation.ERR_SET_IMPORTANT_FRIEND_FAIL_HAS_SET)
			return
		}

		result.Add(relation.NewS2cSetImportantFriendMsg(proto.Id, timeutil.Marshal32(ctime)))

		m.dep.Tlog().TlogSnsFlow(hero, operate_type.SNSSetImportant, u64.FromInt64(targetId))

		result.Changed()
		result.Ok()
	}) {
		return
	}
}

//gogen:iface
func (m *RelationModule) ProcessCancelImportantFriend(proto *relation.C2SCancelImportantFriendProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Id)
	if !ok || targetId <= 0 {
		logrus.Debug("取消星标好友，目标id无效")
		hc.Send(relation.ERR_CANCEL_IMPORTANT_FRIEND_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("取消星标好友，自己的id")
		hc.Send(relation.ERR_CANCEL_IMPORTANT_FRIEND_FAIL_SELF_ID)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		relationType := hero.Relation().GetRelation(targetId)
		if relationType != entity.Friend {
			logrus.Debug("取消星标好友，不是自己好友")
			result.Add(relation.ERR_CANCEL_IMPORTANT_FRIEND_FAIL_NOT_FRIEND)
			return
		}
		if !hero.Relation().CancelImportantFriend(targetId) {
			logrus.Debug("取消星标好友，无法取消非星标好友")
			result.Add(relation.ERR_CANCEL_IMPORTANT_FRIEND_FAIL_HAS_CANCEL)
			return
		}

		result.Add(relation.NewS2cCancelImportantFriendMsg(proto.Id))

		m.dep.Tlog().TlogSnsFlow(hero, operate_type.SNSCancelImportant, u64.FromInt64(targetId))

		result.Changed()
		result.Ok()
	}) {
		return
	}
}
