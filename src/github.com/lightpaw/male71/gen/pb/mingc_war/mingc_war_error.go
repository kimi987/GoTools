package mingc_war

import (
	"github.com/lightpaw/pbutil"
)

// view_mc_war_self_guild
var (
	ErrViewMcWarSelfGuildFailNoGuild   = newMsgError("view_mc_war_self_guild 没有联盟", ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_NO_GUILD)    // 33-3
	ErrViewMcWarSelfGuildFailServerErr = newMsgError("view_mc_war_self_guild 服务器错误", ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_SERVER_ERR) // 33-2
)

// apply_atk
var (
	ErrApplyAtkFailInvalidTime            = newMsgError("apply_atk 当前不是名城战报名时间", ERR_APPLY_ATK_FAIL_INVALID_TIME)              // 18-1
	ErrApplyAtkFailInvalidMcid            = newMsgError("apply_atk 无效的名城", ERR_APPLY_ATK_FAIL_INVALID_MCID)                    // 18-2
	ErrApplyAtkFailApplied                = newMsgError("apply_atk 已经报名攻打别的名城", ERR_APPLY_ATK_FAIL_APPLIED)                    // 18-3
	ErrApplyAtkFailIsHost                 = newMsgError("apply_atk 城主不能攻打自己的名城", ERR_APPLY_ATK_FAIL_IS_HOST)                   // 18-4
	ErrApplyAtkFailHufuNotEnough          = newMsgError("apply_atk 虎符不足", ERR_APPLY_ATK_FAIL_HUFU_NOT_ENOUGH)                  // 18-5
	ErrApplyAtkFailHufuLimit              = newMsgError("apply_atk 虎符未达最低要求", ERR_APPLY_ATK_FAIL_HUFU_LIMIT)                   // 18-6
	ErrApplyAtkFailNotLeader              = newMsgError("apply_atk 盟主才能申请", ERR_APPLY_ATK_FAIL_NOT_LEADER)                     // 18-7
	ErrApplyAtkFailGuildLevelLimit        = newMsgError("apply_atk 联盟等级不够", ERR_APPLY_ATK_FAIL_GUILD_LEVEL_LIMIT)              // 18-8
	ErrApplyAtkFailDuChengNotOpen         = newMsgError("apply_atk 本国都城没开启", ERR_APPLY_ATK_FAIL_DU_CHENG_NOT_OPEN)             // 18-10
	ErrApplyAtkFailServerErr              = newMsgError("apply_atk 服务器错误", ERR_APPLY_ATK_FAIL_SERVER_ERR)                      // 18-9
	ErrApplyAtkFailOtherDuChengNotOpen    = newMsgError("apply_atk 他国都城没开启", ERR_APPLY_ATK_FAIL_OTHER_DU_CHENG_NOT_OPEN)       // 18-11
	ErrApplyAtkFailOtherCountryHasOtherMc = newMsgError("apply_atk 他国还在占领其他名城", ERR_APPLY_ATK_FAIL_OTHER_COUNTRY_HAS_OTHER_MC) // 18-12
	ErrApplyAtkFailNotHoldCapital         = newMsgError("apply_atk 没有占领本国都城", ERR_APPLY_ATK_FAIL_NOT_HOLD_CAPITAL)             // 18-14
)

// apply_ast
var (
	ErrApplyAstFailInvalidTime   = newMsgError("apply_ast 当前不是名城战报名时间", ERR_APPLY_AST_FAIL_INVALID_TIME) // 23-1
	ErrApplyAstFailInvalidMcid   = newMsgError("apply_ast 无效的名城", ERR_APPLY_AST_FAIL_INVALID_MCID)       // 23-2
	ErrApplyAstFailIsHost        = newMsgError("apply_ast 城主不能攻打自己的名城", ERR_APPLY_AST_FAIL_IS_HOST)      // 23-3
	ErrApplyAstFailNotLeader     = newMsgError("apply_ast 盟主才能申请", ERR_APPLY_AST_FAIL_NOT_LEADER)        // 23-4
	ErrApplyAstFailMcidCannotAst = newMsgError("apply_ast 这座城不能协助", ERR_APPLY_AST_FAIL_MCID_CANNOT_AST)  // 23-5
	ErrApplyAstFailMcAstLimit    = newMsgError("apply_ast 这座城协助已达上限", ERR_APPLY_AST_FAIL_MC_AST_LIMIT)   // 23-6
	ErrApplyAstFailAppliedLimit  = newMsgError("apply_ast 达到申请上限", ERR_APPLY_AST_FAIL_APPLIED_LIMIT)     // 23-7
	ErrApplyAstFailAlreadyAst    = newMsgError("apply_ast 已经申请", ERR_APPLY_AST_FAIL_ALREADY_AST)         // 23-9
	ErrApplyAstFailServerErr     = newMsgError("apply_ast 服务器错误", ERR_APPLY_AST_FAIL_SERVER_ERR)         // 23-8
)

// cancel_apply_ast
var (
	ErrCancelApplyAstFailInvalidTime = newMsgError("cancel_apply_ast 当前不是名城战报名时间", ERR_CANCEL_APPLY_AST_FAIL_INVALID_TIME) // 82-1
	ErrCancelApplyAstFailInvalidMcid = newMsgError("cancel_apply_ast 无效的名城", ERR_CANCEL_APPLY_AST_FAIL_INVALID_MCID)       // 82-2
	ErrCancelApplyAstFailNotApply    = newMsgError("cancel_apply_ast 没有申请", ERR_CANCEL_APPLY_AST_FAIL_NOT_APPLY)           // 82-3
	ErrCancelApplyAstFailNotLeader   = newMsgError("cancel_apply_ast 盟主才能操作", ERR_CANCEL_APPLY_AST_FAIL_NOT_LEADER)        // 82-5
	ErrCancelApplyAstFailServerErr   = newMsgError("cancel_apply_ast 服务器错误", ERR_CANCEL_APPLY_AST_FAIL_SERVER_ERR)         // 82-4
)

// reply_apply_ast
var (
	ErrReplyApplyAstFailInvalidTime           = newMsgError("reply_apply_ast 当前不是名城战报名时间", ERR_REPLY_APPLY_AST_FAIL_INVALID_TIME)       // 27-1
	ErrReplyApplyAstFailInvalidMcid           = newMsgError("reply_apply_ast 无效的名城", ERR_REPLY_APPLY_AST_FAIL_INVALID_MCID)             // 27-2
	ErrReplyApplyAstFailNotLeader             = newMsgError("reply_apply_ast 盟主才能审批", ERR_REPLY_APPLY_AST_FAIL_NOT_LEADER)              // 27-8
	ErrReplyApplyAstFailReplyPermissionDenied = newMsgError("reply_apply_ast 无权审批", ERR_REPLY_APPLY_AST_FAIL_REPLY_PERMISSION_DENIED)   // 27-6
	ErrReplyApplyAstFailNotApply              = newMsgError("reply_apply_ast 联盟id 错误", ERR_REPLY_APPLY_AST_FAIL_NOT_APPLY)              // 27-7
	ErrReplyApplyAstFailAstLimit              = newMsgError("reply_apply_ast 协助的联盟数已经到最大了", ERR_REPLY_APPLY_AST_FAIL_AST_LIMIT)         // 27-3
	ErrReplyApplyAstFailTargetAstLimit        = newMsgError("reply_apply_ast 对方协助的名城数已经最大了", ERR_REPLY_APPLY_AST_FAIL_TARGET_AST_LIMIT) // 27-4
	ErrReplyApplyAstFailTargetAlreadyAst      = newMsgError("reply_apply_ast 对方已经协助这座名城了", ERR_REPLY_APPLY_AST_FAIL_TARGET_ALREADY_AST) // 27-9
	ErrReplyApplyAstFailServerErr             = newMsgError("reply_apply_ast 服务器错误", ERR_REPLY_APPLY_AST_FAIL_SERVER_ERR)               // 27-5
)

// view_mingc_war_mc
var (
	ErrViewMingcWarMcFailInvalidMcid = newMsgError("view_mingc_war_mc 无效的名城", ERR_VIEW_MINGC_WAR_MC_FAIL_INVALID_MCID) // 77-1
	ErrViewMingcWarMcFailServerErr   = newMsgError("view_mingc_war_mc 服务器错误", ERR_VIEW_MINGC_WAR_MC_FAIL_SERVER_ERR)   // 77-2
)

// join_fight
var (
	ErrJoinFightFailInvalidTime      = newMsgError("join_fight 当前不是名城战战斗时间", ERR_JOIN_FIGHT_FAIL_INVALID_TIME)       // 37-1
	ErrJoinFightFailInvalidMcid      = newMsgError("join_fight 无效的名城", ERR_JOIN_FIGHT_FAIL_INVALID_MCID)             // 37-2
	ErrJoinFightFailMcWarEnd         = newMsgError("join_fight 这座城的名城战已经结束了", ERR_JOIN_FIGHT_FAIL_MC_WAR_END)        // 37-7
	ErrJoinFightFailNotApply         = newMsgError("join_fight 没有申请攻打/协助名城或不是占领盟", ERR_JOIN_FIGHT_FAIL_NOT_APPLY)    // 37-3
	ErrJoinFightFailAlreadyInWar     = newMsgError("join_fight 已经在其他名城参战了", ERR_JOIN_FIGHT_FAIL_ALREADY_IN_WAR)      // 37-4
	ErrJoinFightFailInvalidCaptainId = newMsgError("join_fight 武将不存在", ERR_JOIN_FIGHT_FAIL_INVALID_CAPTAIN_ID)       // 37-6
	ErrJoinFightFailJoinGuildTooLate = newMsgError("join_fight 战斗开始后加入不让进", ERR_JOIN_FIGHT_FAIL_JOIN_GUILD_TOO_LATE) // 37-8
	ErrJoinFightFailHeroLevelLimit   = newMsgError("join_fight 君主等级不够", ERR_JOIN_FIGHT_FAIL_HERO_LEVEL_LIMIT)        // 37-9
	ErrJoinFightFailServerErr        = newMsgError("join_fight 服务器错误", ERR_JOIN_FIGHT_FAIL_SERVER_ERR)               // 37-5
)

// quit_fight
var (
	ErrQuitFightFailInvalidTime  = newMsgError("quit_fight 当前不是名城战战斗时间", ERR_QUIT_FIGHT_FAIL_INVALID_TIME) // 40-1
	ErrQuitFightFailNotJoinFight = newMsgError("quit_fight 没有参战", ERR_QUIT_FIGHT_FAIL_NOT_JOIN_FIGHT)      // 40-2
	ErrQuitFightFailServerErr    = newMsgError("quit_fight 服务器错误", ERR_QUIT_FIGHT_FAIL_SERVER_ERR)         // 40-3
)

// scene_move
var (
	ErrSceneMoveFailInvalidTime       = newMsgError("scene_move 当前不是名城战战斗时间", ERR_SCENE_MOVE_FAIL_INVALID_TIME)     // 51-1
	ErrSceneMoveFailInvalidMcid       = newMsgError("scene_move 无效的名城", ERR_SCENE_MOVE_FAIL_INVALID_MCID)           // 51-2
	ErrSceneMoveFailNotInScene        = newMsgError("scene_move 没有参战", ERR_SCENE_MOVE_FAIL_NOT_IN_SCENE)            // 51-6
	ErrSceneMoveFailMcWarEnd          = newMsgError("scene_move 这座城的名城战已经结束了", ERR_SCENE_MOVE_FAIL_MC_WAR_END)      // 51-11
	ErrSceneMoveFailInPrepareDuration = newMsgError("scene_move 在准备时间", ERR_SCENE_MOVE_FAIL_IN_PREPARE_DURATION)    // 51-9
	ErrSceneMoveFailInJoinDuration    = newMsgError("scene_move 在入场时间", ERR_SCENE_MOVE_FAIL_IN_JOIN_DURATION)       // 51-10
	ErrSceneMoveFailNotStation        = newMsgError("scene_move 不在驻扎状态", ERR_SCENE_MOVE_FAIL_NOT_STATION)           // 51-8
	ErrSceneMoveFailAlreadyOnDestPos  = newMsgError("scene_move 已经在目的地了", ERR_SCENE_MOVE_FAIL_ALREADY_ON_DEST_POS)  // 51-4
	ErrSceneMoveFailServerErr         = newMsgError("scene_move 服务器错误", ERR_SCENE_MOVE_FAIL_SERVER_ERR)             // 51-5
	ErrSceneMoveFailDestCannotArrive  = newMsgError("scene_move 目的地不能直接到达", ERR_SCENE_MOVE_FAIL_DEST_CANNOT_ARRIVE) // 51-12
	ErrSceneMoveFailNotDestroy        = newMsgError("scene_move 当前占领据点没有被摧毁", ERR_SCENE_MOVE_FAIL_NOT_DESTROY)      // 51-13
	ErrSceneMoveFailIsMoving          = newMsgError("scene_move 在行军中", ERR_SCENE_MOVE_FAIL_IS_MOVING)               // 51-14
	ErrSceneMoveFailIsRelive          = newMsgError("scene_move 在补兵中", ERR_SCENE_MOVE_FAIL_IS_RELIVE)               // 51-15
)

// scene_back
var (
	ErrSceneBackFailInvalidTime       = newMsgError("scene_back 当前不是名城战战斗时间", ERR_SCENE_BACK_FAIL_INVALID_TIME)  // 87-1
	ErrSceneBackFailInvalidMcid       = newMsgError("scene_back 无效的名城", ERR_SCENE_BACK_FAIL_INVALID_MCID)        // 87-2
	ErrSceneBackFailNotInScene        = newMsgError("scene_back 没有参战", ERR_SCENE_BACK_FAIL_NOT_IN_SCENE)         // 87-3
	ErrSceneBackFailMcWarEnd          = newMsgError("scene_back 这座城的名城战已经结束了", ERR_SCENE_BACK_FAIL_MC_WAR_END)   // 87-4
	ErrSceneBackFailInPrepareDuration = newMsgError("scene_back 在准备时间", ERR_SCENE_BACK_FAIL_IN_PREPARE_DURATION) // 87-5
	ErrSceneBackFailInJoinDuration    = newMsgError("scene_back 在入场时间", ERR_SCENE_BACK_FAIL_IN_JOIN_DURATION)    // 87-6
	ErrSceneBackFailNotMove           = newMsgError("scene_back 不是移动状态", ERR_SCENE_BACK_FAIL_NOT_MOVE)           // 87-7
	ErrSceneBackFailServerErr         = newMsgError("scene_back 服务器错误", ERR_SCENE_BACK_FAIL_SERVER_ERR)          // 87-8
)

// scene_speed_up
var (
	ErrSceneSpeedUpFailInvalidGoods      = newMsgError("scene_speed_up 发送的不是行军加速道具", ERR_SCENE_SPEED_UP_FAIL_INVALID_GOODS) // 90-1
	ErrSceneSpeedUpFailGoodsNotEnough    = newMsgError("scene_speed_up 物品个数不足", ERR_SCENE_SPEED_UP_FAIL_GOODS_NOT_ENOUGH)   // 90-2
	ErrSceneSpeedUpFailCostNotSupport    = newMsgError("scene_speed_up 不支持点券购买", ERR_SCENE_SPEED_UP_FAIL_COST_NOT_SUPPORT)  // 90-3
	ErrSceneSpeedUpFailCostNotEnough     = newMsgError("scene_speed_up 点券购买，点券不足", ERR_SCENE_SPEED_UP_FAIL_COST_NOT_ENOUGH) // 90-4
	ErrSceneSpeedUpFailInvalidTime       = newMsgError("scene_speed_up 当前不是名城战战斗时间", ERR_SCENE_SPEED_UP_FAIL_INVALID_TIME)  // 90-5
	ErrSceneSpeedUpFailInvalidMcid       = newMsgError("scene_speed_up 无效的名城", ERR_SCENE_SPEED_UP_FAIL_INVALID_MCID)        // 90-6
	ErrSceneSpeedUpFailNotInScene        = newMsgError("scene_speed_up 没有参战", ERR_SCENE_SPEED_UP_FAIL_NOT_IN_SCENE)         // 90-7
	ErrSceneSpeedUpFailMcWarEnd          = newMsgError("scene_speed_up 这座城的名城战已经结束了", ERR_SCENE_SPEED_UP_FAIL_MC_WAR_END)   // 90-8
	ErrSceneSpeedUpFailInPrepareDuration = newMsgError("scene_speed_up 在准备时间", ERR_SCENE_SPEED_UP_FAIL_IN_PREPARE_DURATION) // 90-9
	ErrSceneSpeedUpFailInJoinDuration    = newMsgError("scene_speed_up 在入场时间", ERR_SCENE_SPEED_UP_FAIL_IN_JOIN_DURATION)    // 90-10
	ErrSceneSpeedUpFailNoMoving          = newMsgError("scene_speed_up 部队不是行军中，不能加速", ERR_SCENE_SPEED_UP_FAIL_NO_MOVING)    // 90-11
	ErrSceneSpeedUpFailServerError       = newMsgError("scene_speed_up 服务器忙，请稍后再试", ERR_SCENE_SPEED_UP_FAIL_SERVER_ERROR)   // 90-12
)

// scene_troop_relive
var (
	ErrSceneTroopReliveFailInvalidTime       = newMsgError("scene_troop_relive 当前不是名城战战斗时间", ERR_SCENE_TROOP_RELIVE_FAIL_INVALID_TIME)  // 73-1
	ErrSceneTroopReliveFailInvalidMcid       = newMsgError("scene_troop_relive 无效的名城", ERR_SCENE_TROOP_RELIVE_FAIL_INVALID_MCID)        // 73-2
	ErrSceneTroopReliveFailNotInScene        = newMsgError("scene_troop_relive 没有参战", ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_SCENE)         // 73-3
	ErrSceneTroopReliveFailMcWarEnd          = newMsgError("scene_troop_relive 这座城的名城战已经结束了", ERR_SCENE_TROOP_RELIVE_FAIL_MC_WAR_END)   // 73-4
	ErrSceneTroopReliveFailInPrepareDuration = newMsgError("scene_troop_relive 在准备时间", ERR_SCENE_TROOP_RELIVE_FAIL_IN_PREPARE_DURATION) // 73-5
	ErrSceneTroopReliveFailInJoinDuration    = newMsgError("scene_troop_relive 在入场时间", ERR_SCENE_TROOP_RELIVE_FAIL_IN_JOIN_DURATION)    // 73-6
	ErrSceneTroopReliveFailNotStation        = newMsgError("scene_troop_relive 不是驻扎状态", ERR_SCENE_TROOP_RELIVE_FAIL_NOT_STATION)        // 73-7
	ErrSceneTroopReliveFailNotInRelivePos    = newMsgError("scene_troop_relive 不在复活点", ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_RELIVE_POS)   // 73-10
	ErrSceneTroopReliveFailFullSolider       = newMsgError("scene_troop_relive 兵是满的", ERR_SCENE_TROOP_RELIVE_FAIL_FULL_SOLIDER)         // 73-8
	ErrSceneTroopReliveFailServerErr         = newMsgError("scene_troop_relive 服务器错误", ERR_SCENE_TROOP_RELIVE_FAIL_SERVER_ERR)          // 73-9
)

// view_mc_war_scene
var (
	ErrViewMcWarSceneFailInvalidTime = newMsgError("view_mc_war_scene 当前不是名城战战斗时间", ERR_VIEW_MC_WAR_SCENE_FAIL_INVALID_TIME) // 48-1
	ErrViewMcWarSceneFailInvalidMcid = newMsgError("view_mc_war_scene 无效的名城", ERR_VIEW_MC_WAR_SCENE_FAIL_INVALID_MCID)       // 48-2
	ErrViewMcWarSceneFailServerErr   = newMsgError("view_mc_war_scene 服务器错误", ERR_VIEW_MC_WAR_SCENE_FAIL_SERVER_ERR)         // 48-3
)

// watch
var (
	ErrWatchFailInvalidTime = newMsgError("watch 当前不是名城战战斗时间", ERR_WATCH_FAIL_INVALID_TIME) // 141-1
	ErrWatchFailInvalidMcid = newMsgError("watch 无效的名城", ERR_WATCH_FAIL_INVALID_MCID)       // 141-2
	ErrWatchFailServerErr   = newMsgError("watch 服务器错误", ERR_WATCH_FAIL_SERVER_ERR)         // 141-3
)

// quit_watch
var (
	ErrQuitWatchFailInvalidTime = newMsgError("quit_watch 当前不是名城战战斗时间", ERR_QUIT_WATCH_FAIL_INVALID_TIME) // 138-2
	ErrQuitWatchFailInvalidMcid = newMsgError("quit_watch 无效的名城", ERR_QUIT_WATCH_FAIL_INVALID_MCID)       // 138-3
	ErrQuitWatchFailServerErr   = newMsgError("quit_watch 服务器错误", ERR_QUIT_WATCH_FAIL_SERVER_ERR)         // 138-4
)

// view_mc_war_record
var (
	ErrViewMcWarRecordFailNoRecord = newMsgError("view_mc_war_record 找不到记录", ERR_VIEW_MC_WAR_RECORD_FAIL_NO_RECORD) // 93-1
)

// view_mc_war_troop_record
var (
	ErrViewMcWarTroopRecordFailNoRecord = newMsgError("view_mc_war_troop_record 找不到记录", ERR_VIEW_MC_WAR_TROOP_RECORD_FAIL_NO_RECORD) // 96-1
)

// view_scene_troop_record
var (
	ErrViewSceneTroopRecordFailInvalidTime = newMsgError("view_scene_troop_record 当前不是名城战战斗时间", ERR_VIEW_SCENE_TROOP_RECORD_FAIL_INVALID_TIME) // 101-1
	ErrViewSceneTroopRecordFailNotInScene  = newMsgError("view_scene_troop_record 没有参战", ERR_VIEW_SCENE_TROOP_RECORD_FAIL_NOT_IN_SCENE)        // 101-3
	ErrViewSceneTroopRecordFailMcWarEnd    = newMsgError("view_scene_troop_record 这座城的名城战已经结束了", ERR_VIEW_SCENE_TROOP_RECORD_FAIL_MC_WAR_END)  // 101-4
	ErrViewSceneTroopRecordFailServerErr   = newMsgError("view_scene_troop_record 服务器错误", ERR_VIEW_SCENE_TROOP_RECORD_FAIL_SERVER_ERR)         // 101-5
)

// apply_refresh_rank
var (
	ErrApplyRefreshRankFailServerErr   = newMsgError("apply_refresh_rank 服务器错误", ERR_APPLY_REFRESH_RANK_FAIL_SERVER_ERR)         // 109-1
	ErrApplyRefreshRankFailInvalidTime = newMsgError("apply_refresh_rank 当前不是名城战战斗时间", ERR_APPLY_REFRESH_RANK_FAIL_INVALID_TIME) // 109-2
	ErrApplyRefreshRankFailNotInScene  = newMsgError("apply_refresh_rank 没有参战", ERR_APPLY_REFRESH_RANK_FAIL_NOT_IN_SCENE)        // 109-3
)

// view_my_guild_member_rank
var (
	ErrViewMyGuildMemberRankFailNoRecord = newMsgError("view_my_guild_member_rank 找不到记录", ERR_VIEW_MY_GUILD_MEMBER_RANK_FAIL_NO_RECORD) // 113-1
)

// scene_change_mode
var (
	ErrSceneChangeModeFailInvalidMode    = newMsgError("scene_change_mode mode 不存在", ERR_SCENE_CHANGE_MODE_FAIL_INVALID_MODE)    // 117-9
	ErrSceneChangeModeFailInvalidTime    = newMsgError("scene_change_mode 当前不是名城战战斗时间", ERR_SCENE_CHANGE_MODE_FAIL_INVALID_TIME) // 117-5
	ErrSceneChangeModeFailNotInScene     = newMsgError("scene_change_mode 没有参战", ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_SCENE)        // 117-7
	ErrSceneChangeModeFailMcWarEnd       = newMsgError("scene_change_mode 这座城的名城战已经结束了", ERR_SCENE_CHANGE_MODE_FAIL_MC_WAR_END)  // 117-8
	ErrSceneChangeModeFailSameMode       = newMsgError("scene_change_mode 已经是新形态了", ERR_SCENE_CHANGE_MODE_FAIL_SAME_MODE)        // 117-1
	ErrSceneChangeModeFailNotStation     = newMsgError("scene_change_mode 不是驻扎状态", ERR_SCENE_CHANGE_MODE_FAIL_NOT_STATION)       // 117-2
	ErrSceneChangeModeFailNotInRelivePos = newMsgError("scene_change_mode 不在复活点", ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_RELIVE_POS)  // 117-3
	ErrSceneChangeModeFailServerErr      = newMsgError("scene_change_mode 服务器错误", ERR_SCENE_CHANGE_MODE_FAIL_SERVER_ERR)         // 117-4
)

// scene_tou_shi_building_turn_to
var (
	ErrSceneTouShiBuildingTurnToFailInvalidTime = newMsgError("scene_tou_shi_building_turn_to 当前不是名城战战斗时间", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_TIME) // 121-6
	ErrSceneTouShiBuildingTurnToFailNotInScene  = newMsgError("scene_tou_shi_building_turn_to 没有参战", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_IN_SCENE)        // 121-7
	ErrSceneTouShiBuildingTurnToFailMcWarEnd    = newMsgError("scene_tou_shi_building_turn_to 这座城的名城战已经结束了", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_MC_WAR_END)  // 121-8
	ErrSceneTouShiBuildingTurnToFailInvalidPos  = newMsgError("scene_tou_shi_building_turn_to 坐标错误或不是投石机", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_POS)   // 121-1
	ErrSceneTouShiBuildingTurnToFailNotHost     = newMsgError("scene_tou_shi_building_turn_to 不是占领者", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_HOST)           // 121-2
	ErrSceneTouShiBuildingTurnToFailNoTarget    = newMsgError("scene_tou_shi_building_turn_to 没有目标", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NO_TARGET)           // 121-3
	ErrSceneTouShiBuildingTurnToFailInTurnCd    = newMsgError("scene_tou_shi_building_turn_to 正在转向中", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_TURN_CD)         // 121-4
	ErrSceneTouShiBuildingTurnToFailInPrepareCd = newMsgError("scene_tou_shi_building_turn_to 正在装填中", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_PREPARE_CD)      // 121-5
	ErrSceneTouShiBuildingTurnToFailServerErr   = newMsgError("scene_tou_shi_building_turn_to 服务器错误", ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_SERVER_ERR)         // 121-9
)

// scene_tou_shi_building_fire
var (
	ErrSceneTouShiBuildingFireFailInvalidTime = newMsgError("scene_tou_shi_building_fire 当前不是名城战战斗时间", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_TIME) // 125-6
	ErrSceneTouShiBuildingFireFailNotInScene  = newMsgError("scene_tou_shi_building_fire 没有参战", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_IN_SCENE)        // 125-7
	ErrSceneTouShiBuildingFireFailMcWarEnd    = newMsgError("scene_tou_shi_building_fire 这座城的名城战已经结束了", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_MC_WAR_END)  // 125-8
	ErrSceneTouShiBuildingFireFailInvalidPos  = newMsgError("scene_tou_shi_building_fire 坐标错误或不是投石机", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_POS)   // 125-1
	ErrSceneTouShiBuildingFireFailNotHost     = newMsgError("scene_tou_shi_building_fire 不是占领者", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_HOST)           // 125-2
	ErrSceneTouShiBuildingFireFailNoTarget    = newMsgError("scene_tou_shi_building_fire 没有目标", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NO_TARGET)           // 125-3
	ErrSceneTouShiBuildingFireFailInTurnCd    = newMsgError("scene_tou_shi_building_fire 正在转向中", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_TURN_CD)         // 125-4
	ErrSceneTouShiBuildingFireFailInPrepareCd = newMsgError("scene_tou_shi_building_fire 正在装填中", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_PREPARE_CD)      // 125-5
	ErrSceneTouShiBuildingFireFailServerErr   = newMsgError("scene_tou_shi_building_fire 服务器错误", ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_SERVER_ERR)         // 125-9
)

// scene_drum
var (
	ErrSceneDrumFailNotInScene        = newMsgError("scene_drum 没有参战", ERR_SCENE_DRUM_FAIL_NOT_IN_SCENE)             // 130-1
	ErrSceneDrumFailInvalidTime       = newMsgError("scene_drum 击鼓时间已经结束", ERR_SCENE_DRUM_FAIL_INVALID_TIME)         // 130-2
	ErrSceneDrumFailInCd              = newMsgError("scene_drum 还在击鼓 CD 中", ERR_SCENE_DRUM_FAIL_IN_CD)               // 130-3
	ErrSceneDrumFailBaiZhanLevelLimit = newMsgError("scene_drum 百战千军等级不够", ERR_SCENE_DRUM_FAIL_BAI_ZHAN_LEVEL_LIMIT) // 130-5
	ErrSceneDrumFailServerErr         = newMsgError("scene_drum 服务器错误", ERR_SCENE_DRUM_FAIL_SERVER_ERR)              // 130-4
)

func newMsgError(msg string, buffer pbutil.StaticBuffer) *error_msg {
	return &error_msg{
		msg:  msg,
		buff: buffer,
	}
}

type error_msg struct {
	msg  string
	buff pbutil.Buffer
}

func (f *error_msg) Error() string         { return f.msg }
func (f *error_msg) ErrMsg() pbutil.Buffer { return f.buff }
