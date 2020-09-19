package mingc_war

import (
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/pb/server_proto"
)

func applyAttack(w *MingcWar, gid int64, mcData *mingcdata.MingcBaseData, cost uint64, ctime time.Time) (errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_APPLY_ATK {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_INVALID_TIME
		return
	}

	if w.bidValue(gid) <= 0 && cost < mcData.AtkMinHufu {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_HUFU_LIMIT
		return
	}

	stage := w.currStage().(*ApplyAtkObj)
	if mc, ok := stage.applicants[gid]; ok {
		if mc != mcData.Id {
			errMsg = mingc_war.ERR_APPLY_ATK_FAIL_APPLIED
			return
		}
	}

	if mc := w.mingcs[mcData.Id]; mc != nil && mc.defId == gid {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_IS_HOST
		return
	}

	stage.applicants[gid] = mcData.Id
	cost += w.mingcs[mcData.Id].bid.GetBid(gid)
	_, succ := w.mingcs[mcData.Id].bid.Bidding(gid, cost, ctime)
	if !succ {
		logrus.Debugf("名城申请，Bidding 失败")
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_SERVER_ERR
		return
	}

	return
}

func applyAssist(w *MingcWar, gid int64, isAtk bool, mcData *mingcdata.MingcBaseData, d *mingcdata.MingcMiscData) (recvGId int64, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_APPLY_AST {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_INVALID_TIME
		return
	}

	stage := w.currStage().(*ApplyAstObj)
	if len(stage.AstGuilds[gid])+len(stage.ApplyAstGuilds[gid]) >= u64.Int(d.ApplyAstLimit) {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_APPLIED_LIMIT
		return
	}

	mc := w.mingcs[mcData.Id]
	if mc == nil {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_INVALID_MCID
		return
	}

	if mc.atkId <= 0 {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_MCID_CANNOT_AST
		return
	}

	if isAtk && len(mc.astAtkList) >= int(mcData.AstMaxGuild) {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_MC_AST_LIMIT
		return
	}

	if !isAtk && len(mc.astDefList) >= int(mcData.AstMaxGuild) {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_MC_AST_LIMIT
		return
	}

	if !isAtk && mc.defId <= 0 {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_MCID_CANNOT_AST
		return
	}

	if isApplyAtk, ok := mc.applyAsts[gid]; ok && isAtk == isApplyAtk {
		errMsg = mingc_war.ERR_APPLY_AST_FAIL_ALREADY_AST
		return
	}

	mc.applyAsts[gid] = isAtk
	if isAtk {
		recvGId = mc.atkId
	} else {
		recvGId = mc.defId
	}

	stage.addApplyAstGuilds(gid, mc.id)

	return
}

func cancelApplyAst(w *MingcWar, gid int64, mcData *mingcdata.MingcBaseData) (receiveId int64, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_APPLY_AST {
		errMsg = mingc_war.ERR_CANCEL_APPLY_AST_FAIL_INVALID_TIME
		return
	}

	mc := w.mingcs[mcData.Id]
	if mc == nil {
		errMsg = mingc_war.ERR_CANCEL_APPLY_AST_FAIL_INVALID_MCID
		return
	}

	if isAtk, ok := mc.applyAsts[gid]; !ok {
		errMsg = mingc_war.ERR_CANCEL_APPLY_AST_FAIL_NOT_APPLY
		return
	} else {
		if isAtk {
			receiveId = mc.atkId
		} else {
			receiveId = mc.defId
		}
	}

	cancelApplyAst0(w, mc, gid)
	return
}

func cleanApplyAstList(w *MingcWar, mc *MingcObj, astAtk bool) {
	for gid, isAtk := range mc.applyAsts {
		if isAtk != astAtk {
			continue
		}
		cancelApplyAst0(w, mc, gid)
	}
}

func cancelApplyAst0(w *MingcWar, mc *MingcObj, gid int64) {
	delete(mc.applyAsts, gid)
	stage := w.currStage().(*ApplyAstObj)
	stage.delApplyAstGuilds(gid, mc.id)
}

func replyApplyAst(w *MingcWar, operId, gid int64, mcData *mingcdata.MingcBaseData, agree bool, d *mingcdata.MingcMiscData) (errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_APPLY_AST {
		errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_INVALID_TIME
		return
	}

	stage := w.currStage().(*ApplyAstObj)
	mc := w.mingcs[mcData.Id]

	if mc == nil {
		errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_INVALID_MCID
		return
	}

	if mc.defId != operId && mc.atkId != operId {
		errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_REPLY_PERMISSION_DENIED
		return
	}

	if agree {
		if mc.atkId == operId && len(mc.astAtkList) >= int(mcData.AstMaxGuild) {
			errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_AST_LIMIT
			return
		} else if mc.defId == operId && len(mc.astDefList) >= int(mcData.AstMaxGuild) {
			errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_AST_LIMIT
			return
		}
	}

	if mc.isAtk(gid) || mc.isDef(gid) {
		delete(mc.applyAsts, gid)
		errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_TARGET_ALREADY_AST
		return
	}

	if !mc.replyApplyAst(operId, gid, mc.id, agree) {
		errMsg = mingc_war.ERR_REPLY_APPLY_AST_FAIL_NOT_APPLY
		return
	}

	// 援助申请通过
	if agree {
		stage.addAstGuilds(gid, mc.id)

		// 援助满，取消其他报名
		if mc.atkId == operId && len(mc.astAtkList) >= int(mcData.AstMaxGuild) {
			cleanApplyAstList(w, mc, true)
		} else if mc.defId == operId && len(mc.astDefList) >= int(mcData.AstMaxGuild) {
			cleanApplyAstList(w, mc, false)
		}
	}

	return
}

// ********** watch **********

func watch(w *MingcWar, mcId uint64, heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_WATCH_FAIL_INVALID_TIME
		return
	}

	mc := w.mingcs[mcId]
	if mc == nil {
		errMsg = mingc_war.ERR_WATCH_FAIL_INVALID_MCID
		return
	}
	if mc.scene == nil {
		errMsg = mingc_war.ERR_WATCH_FAIL_INVALID_TIME
		return
	}
	scene := mc.scene
	if scene.ended {
		logrus.Debugf("名城战观战，mcId:%v 已结束", mcId)
		errMsg = mingc_war.ERR_WATCH_FAIL_INVALID_TIME
		return
	}

	scene.watch(heroId)
	succMsg = mingc_war.WATCH_S2C

	return
}

func quitWatch(w *MingcWar, mcId uint64, heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_QUIT_WATCH_FAIL_INVALID_TIME
		return
	}

	mc := w.mingcs[mcId]
	if mc == nil {
		errMsg = mingc_war.ERR_QUIT_WATCH_FAIL_INVALID_MCID
		return
	}
	if mc.scene == nil {
		errMsg = mingc_war.ERR_QUIT_WATCH_FAIL_INVALID_TIME
		return
	}
	scene := mc.scene
	if scene.ended {
		logrus.Debugf("名城战退出观战，mcId:%v 已结束", mcId)
		errMsg = mingc_war.ERR_QUIT_WATCH_FAIL_INVALID_TIME
		return
	}

	scene.quitWatch(heroId)
	succMsg = mingc_war.QUIT_WATCH_S2C

	return
}

// ********** fight **********

func joinFight(w *MingcWar, mcId uint64, gid, heroId int64, hero *shared_proto.HeroBasicProto, captains []*shared_proto.CaptainInfoProto, burnPos cb.Cube, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_INVALID_TIME
		return
	}

	if _, ok := w.joinedHeros[heroId]; ok {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_ALREADY_IN_WAR
		return
	}

	mc := w.mingcs[mcId]
	if mc == nil {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_INVALID_MCID
		return
	}
	if mc.scene == nil {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_MC_WAR_END
		return
	}
	scene := mc.scene
	if scene.ended {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_MC_WAR_END
		return
	}

	isAtk := mc.isAtk(gid)
	if !isAtk && !mc.isDef(gid) {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_NOT_APPLY
		return
	}

	fightRunStartTime := scene.startTime.Add(mc.dep.Datas().MingcMiscData().FightPrepareDuration)
	var endTime time.Time
	if ctime.After(fightRunStartTime) {
		endTime = ctime.Add(mc.dep.Datas().MingcMiscData().JoinFightDuration)
	} else {
		endTime = fightRunStartTime
	}

	rankObj := scene.troopsRank.getRankObj(heroId)
	if rankObj == nil {
		rankObj = newMcWarTroopRankObject(heroId, hero, isAtk)
		scene.troopsRank.putRankObj(rankObj)
		mc.dep.HeroData().FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_McWarJoin)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_MCWAR_JOIN)
			result.Ok()
		})
	}

	troop := newMcWarTroop(rankObj, gid, captains, burnPos, ctime, endTime, w.datas)
	w.joinedHeros[heroId] = mcId
	scene.joinFight(troop)

	succMsg = mingc_war.NewS2cJoinFightMsg(u64.Int32(mc.id), isAtk)

	return
}

func quitFight(w *MingcWar, heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_QUIT_FIGHT_FAIL_INVALID_TIME
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_QUIT_FIGHT_FAIL_NOT_JOIN_FIGHT
		return
	} else {
		mc = w.mingcs[mcId]
		if mc == nil {
			logrus.Errorf("quitFight，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_QUIT_FIGHT_FAIL_SERVER_ERR
			return
		}
	}

	delete(w.joinedHeros, heroId)
	troop := mc.scene.quitFight(heroId)
	mc.saveTroopRecord(w.id, troop)

	succMsg = mingc_war.NewS2cQuitFightMsg(u64.Int32(mc.id))
	return
}

func sceneMove(w *MingcWar, heroId int64, dest cb.Cube, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_INVALID_TIME
		return
	}

	miscData := w.datas.MingcMiscData()
	if ctime.Before(w.currStage().stageStartTime().Add(miscData.FightPrepareDuration)) {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_IN_PREPARE_DURATION
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneMove，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_MC_WAR_END
		return
	}

	if t, _ := mc.scene.troop(heroId); t == nil {
		logrus.Debugf("名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_NOT_IN_SCENE
		return
	} else if t.action.getState() == shared_proto.MingcWarTroopState_MC_TP_WAIT && ctime.Before(t.action.getEndTime()) {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_IN_JOIN_DURATION
		return
	} else if t.action.getState() == shared_proto.MingcWarTroopState_MC_TP_MOVING {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_IS_MOVING
		return
	}else if t.action.getState() == shared_proto.MingcWarTroopState_MC_TP_RELIVE {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_IS_RELIVE
		return
	}

	pos, ok := mc.scene.pos(heroId)
	if !ok || pos == dest {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_ALREADY_ON_DEST_POS
		return
	}

	if mcSceneData := w.datas.GetMingcWarSceneData(mc.id); mcSceneData == nil {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_DEST_CANNOT_ARRIVE
		return
	} else if !mcSceneData.CanArrive(pos, dest) {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_DEST_CANNOT_ARRIVE
		return
	}

	if !mc.scene.canArrive(heroId, dest) {
		logrus.Debugf("不能在正摧毁一个敌方据点时，移动到另一个没被摧毁的敌方据点。heroId:%v", heroId)
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_NOT_DESTROY
		return
	}

	endTime, succ := mc.scene.move(heroId, dest, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_NOT_STATION
		return
	}

	x, y := dest.XYI32()
	succMsg = mingc_war.NewS2cSceneMoveMsg(x, y, timeutil.Marshal32(endTime))

	return
}

func sceneBack(w *MingcWar, heroId int64, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_INVALID_TIME
		return
	}

	miscData := w.datas.MingcMiscData()
	if ctime.Before(w.currStage().stageStartTime().Add(miscData.FightPrepareDuration)) {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_IN_PREPARE_DURATION
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneBack，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_BACK_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_MC_WAR_END
		return
	}

	if t, _ := mc.scene.troop(heroId); t == nil {
		logrus.Debugf("sceneBack 名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_NOT_IN_SCENE
		return
	}

	endTime, succ := mc.scene.back(heroId, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_NOT_MOVE
		return
	}

	pos, ok := mc.scene.pos(heroId)
	if !ok {
		logrus.Warnf("sceneBack，找不到队伍起始点。mcid:%v heroId:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_SERVER_ERR
		return
	}

	x, y := pos.XYI32()
	succMsg = mingc_war.NewS2cSceneBackMsg(x, y, timeutil.Marshal32(endTime))

	return
}

func sceneSpeedUp(w *MingcWar, heroId int64, speedUpRate float64, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_INVALID_TIME
		return
	}

	miscData := w.datas.MingcMiscData()
	if ctime.Before(w.currStage().stageStartTime().Add(miscData.FightPrepareDuration)) {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_IN_PREPARE_DURATION
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneBack，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_SERVER_ERROR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_MC_WAR_END
		return
	}

	if t, _ := mc.scene.troop(heroId); t == nil {
		logrus.Debugf("sceneBack 名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_NOT_IN_SCENE
		return
	} else {
		if ctime.Before(t.joinTime) {
			errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_IN_JOIN_DURATION
			return
		}
	}

	endTime, succ := mc.scene.speedUp(heroId, speedUpRate, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_NO_MOVING
		return
	}

	succMsg = mingc_war.NewS2cSceneSpeedUpMsg(idbytes.ToBytes(heroId), timeutil.Marshal32(endTime))

	return
}

func sceneRelive(w *MingcWar, heroId int64, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_INVALID_TIME
		return
	}

	miscData := w.datas.MingcMiscData()
	if ctime.Before(w.currStage().stageStartTime().Add(miscData.FightPrepareDuration)) {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_IN_PREPARE_DURATION
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneRelive，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_MC_WAR_END
		return
	}

	var t *McWarTroop
	if t, _ = mc.scene.troop(heroId); t == nil {
		logrus.Debugf("sceneRelive,名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_SCENE
		return
	} else {
		if ctime.Before(t.joinTime.Add(miscData.JoinFightDuration)) {
			errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_IN_JOIN_DURATION
			return
		}
		if t.isFullSolider() {
			errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_FULL_SOLIDER
			return
		}
	}

	endTime, succ := mc.scene.relive(heroId, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_NOT_STATION
		return
	}

	succMsg = mingc_war.NewS2cSceneTroopReliveMsg(timeutil.Marshal32(endTime))

	return
}

func sceneChangeMode(w *MingcWar, heroId int64, newMode shared_proto.MingcWarModeType) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_INVALID_TIME
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("changeMode，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_MC_WAR_END
		return
	}

	var t *McWarTroop
	if t, _ = mc.scene.troop(heroId); t == nil {
		logrus.Debugf("changeMode,名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_SCENE
		return
	}

	if t.mode == newMode {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_SAME_MODE
		return
	}

	succ := mc.scene.changeMode(heroId, newMode)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_NOT_STATION
		return
	}

	succMsg = mingc_war.NewS2cSceneChangeModeMsg(int32(t.mode))

	return
}

func sceneTouShiBuildingTurnTo(w *MingcWar, heroId int64, pos cb.Cube, left bool, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_TIME
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneTouShiBuildingTurnTo，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_MC_WAR_END
		return
	}

	if b, ok := mc.scene.buildings[pos]; !ok || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_POS
		return
	} else if t, ok := b.troops[heroId]; !ok || t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_HOST
		return
	} else if len(b.touShiTargets) <= 0 {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NO_TARGET
		return
	} else if ctime.Before(b.touShiPrepareEndTime) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_PREPARE_CD
		return
	} else if ctime.Before(b.touShiTurnEndTime) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_TURN_CD
		return
	}

	succ, newTargetIndex, turnEndTime := mc.scene.touShiBuildingTurnTo(heroId, pos, left, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_SERVER_ERR
		return
	}

	x, y := pos.XYI32()
	succMsg = mingc_war.NewS2cSceneTouShiBuildingTurnToMsg(x, y, left, int32(newTargetIndex), timeutil.Marshal32(turnEndTime))

	return
}

func sceneTouShiBuildingFire(w *MingcWar, heroId int64, pos cb.Cube, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_TIME
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneTouShiBuildingFire，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_MC_WAR_END
		return
	}

	b, ok := mc.scene.buildings[pos]
	if !ok || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_POS
		return
	} else if t, ok := b.troops[heroId]; !ok || t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_HOST
		return
	} else if len(b.touShiTargets) <= 0 || b.touShiTargetIndex >= len(b.touShiTargets) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NO_TARGET
		return
	} else if _, ok := mc.scene.buildings[b.touShiTargets[b.touShiTargetIndex]]; !ok {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NO_TARGET
		return
	} else if ctime.Before(b.touShiPrepareEndTime) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_PREPARE_CD
		return
	} else if ctime.Before(b.touShiTurnEndTime) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_TURN_CD
		return
	}

	succ, bombExplodeTime := mc.scene.touShiBuildingFire(heroId, pos, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_SERVER_ERR
		return
	}

	x, y := pos.XYI32()
	succMsg = mingc_war.NewS2cSceneTouShiBuildingFireMsg(x, y, int32(b.touShiTargetIndex), timeutil.Marshal32(b.touShiPrepareEndTime), timeutil.Marshal32(bombExplodeTime))

	return
}

func sceneDrum(w *MingcWar, heroId int64, ctime time.Time) (succMsg, errMsg pbutil.Buffer) {
	if w.state != shared_proto.MingcWarState_MC_T_FIGHT {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_INVALID_TIME
		return
	}

	var mc *MingcObj
	if mcId, ok := w.joinedHeros[heroId]; !ok {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_NOT_IN_SCENE
		return
	} else {
		if mc = w.mingcs[mcId]; mc == nil {
			logrus.Errorf("sceneDrum，找不到名城.id:%v", mcId)
			errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_SERVER_ERR
			return
		}
	}

	if mc.scene == nil || mc.scene.ended {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_INVALID_TIME
		return
	}

	if ctime.After(mc.scene.durmStopTime) {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_INVALID_TIME
		return
	}

	var t *McWarTroop
	if t, _ = mc.scene.troop(heroId); t == nil {
		logrus.Debugf("sceneDrum,名城战场景中找不到部队。mcid:%v heroid:%v", mc.id, heroId)
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_NOT_IN_SCENE
		return
	}

	if ctime.Before(t.nextDrumTime) {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_IN_CD
		return
	}

	succ, toAdd, desc, nextTime := mc.scene.drum(heroId, ctime)
	if !succ {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_SERVER_ERR
		return
	}
	succMsg = mingc_war.NewS2cSceneDrumMsg(desc, timeutil.Marshal32(nextTime), toAdd)

	return
}

func getGuildBasicProto(id int64, dep iface.ServiceDep) (p *shared_proto.GuildBasicProto) {
	p = dep.GuildSnapshot().GetGuildBasicProto(id)
	if p != nil {
		return
	}

	guildDataId := mingcdata.RecoverMcWarGuildDataId(id)
	if d := dep.Datas().GetMingcWarNpcGuildData(guildDataId); d != nil {
		p = d.GuildBasicProto()
	}

	return
}

func gmTransMingcWarTime(w *MingcWar, toTrans time.Duration) {
	w.stages[shared_proto.MingcWarState_MC_T_NOT_START].(*Stage).startTime = w.stages[shared_proto.MingcWarState_MC_T_NOT_START].(*Stage).startTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_NOT_START].(*Stage).endTime = w.stages[shared_proto.MingcWarState_MC_T_NOT_START].(*Stage).endTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK].(*ApplyAtkObj).startTime = w.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK].(*ApplyAtkObj).startTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK].(*ApplyAtkObj).endTime = w.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK].(*ApplyAtkObj).endTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_AST].(*ApplyAstObj).startTime = w.stages[shared_proto.MingcWarState_MC_T_APPLY_AST].(*ApplyAstObj).startTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_AST].(*ApplyAstObj).endTime = w.stages[shared_proto.MingcWarState_MC_T_APPLY_AST].(*ApplyAstObj).endTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT].(*FightObj).startTime = w.stages[shared_proto.MingcWarState_MC_T_FIGHT].(*FightObj).startTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT].(*FightObj).endTime = w.stages[shared_proto.MingcWarState_MC_T_FIGHT].(*FightObj).endTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT_END].(*Stage).startTime = w.stages[shared_proto.MingcWarState_MC_T_FIGHT_END].(*Stage).startTime.Add(toTrans)
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT_END].(*Stage).endTime = w.stages[shared_proto.MingcWarState_MC_T_FIGHT_END].(*Stage).endTime.Add(toTrans)

	w.startTime = w.startTime.Add(toTrans)
	w.endTime = w.endTime.Add(toTrans)
}
