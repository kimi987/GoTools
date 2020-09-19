package guild

import (
	"github.com/lightpaw/pbutil"
)

// list_guild
var (
	ErrListGuildFailServerError = newMsgError("list_guild 服务器忙，请稍后再试", ERR_LIST_GUILD_FAIL_SERVER_ERROR) // 3-1
)

// search_guild
var (
	ErrSearchGuildFailInvalidName = newMsgError("search_guild 无效的搜索名字", ERR_SEARCH_GUILD_FAIL_INVALID_NAME)    // 6-2
	ErrSearchGuildFailInvalidNum  = newMsgError("search_guild 无效的页数", ERR_SEARCH_GUILD_FAIL_INVALID_NUM)       // 6-3
	ErrSearchGuildFailServerError = newMsgError("search_guild 服务器忙，请稍后再试", ERR_SEARCH_GUILD_FAIL_SERVER_ERROR) // 6-1
)

// create_guild
var (
	ErrCreateGuildFailInTheGuild         = newMsgError("create_guild 已经在联盟中，不能创建联盟", ERR_CREATE_GUILD_FAIL_IN_THE_GUILD)    // 9-1
	ErrCreateGuildFailInvalidNameLen     = newMsgError("create_guild 无效的名字长度", ERR_CREATE_GUILD_FAIL_INVALID_NAME_LEN)      // 9-5
	ErrCreateGuildFailNameDuplicate      = newMsgError("create_guild 联盟名字已经存在", ERR_CREATE_GUILD_FAIL_NAME_DUPLICATE)       // 9-2
	ErrCreateGuildFailFlagNameDuplicate  = newMsgError("create_guild 联盟旗号已经存在", ERR_CREATE_GUILD_FAIL_FLAG_NAME_DUPLICATE)  // 9-3
	ErrCreateGuildFailInvalidFlagNameLen = newMsgError("create_guild 无效的旗号长度", ERR_CREATE_GUILD_FAIL_INVALID_FLAG_NAME_LEN) // 9-6
	ErrCreateGuildFailCostNotEnough      = newMsgError("create_guild 创建联盟消耗不足", ERR_CREATE_GUILD_FAIL_COST_NOT_ENOUGH)      // 9-7
	ErrCreateGuildFailHeroNoCountry      = newMsgError("create_guild 玩家没有国家", ERR_CREATE_GUILD_FAIL_HERO_NO_COUNTRY)        // 9-9
	ErrCreateGuildFailSensitiveWords     = newMsgError("create_guild 输入包含敏感词", ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS)       // 9-8
	ErrCreateGuildFailServerError        = newMsgError("create_guild 服务器忙，请稍后再试", ERR_CREATE_GUILD_FAIL_SERVER_ERROR)       // 9-4
)

// self_guild
var (
	ErrSelfGuildFailNotInGuild  = newMsgError("self_guild 你没有联盟", ERR_SELF_GUILD_FAIL_NOT_IN_GUILD)      // 12-2
	ErrSelfGuildFailServerError = newMsgError("self_guild 服务器忙，请稍后再试", ERR_SELF_GUILD_FAIL_SERVER_ERROR) // 12-1
)

// leave_guild
var (
	ErrLeaveGuildFailNotInGuild      = newMsgError("leave_guild 你没有联盟", ERR_LEAVE_GUILD_FAIL_NOT_IN_GUILD)                  // 15-1
	ErrLeaveGuildFailLeader          = newMsgError("leave_guild 你是盟主，不能退出", ERR_LEAVE_GUILD_FAIL_LEADER)                    // 15-2
	ErrLeaveGuildFailXiongNuDefender = newMsgError("leave_guild 匈奴入侵防守队员，且当前活动已开启", ERR_LEAVE_GUILD_FAIL_XIONG_NU_DEFENDER) // 15-4
	ErrLeaveGuildFailIsMcWarAtk      = newMsgError("leave_guild 名城战进攻盟不能解散", ERR_LEAVE_GUILD_FAIL_IS_MC_WAR_ATK)            // 15-5
	ErrLeaveGuildFailIsMcWarDef      = newMsgError("leave_guild 名城战防守盟不能解散", ERR_LEAVE_GUILD_FAIL_IS_MC_WAR_DEF)            // 15-6
	ErrLeaveGuildFailInMcWarFight    = newMsgError("leave_guild 名城战战斗阶段不能退出", ERR_LEAVE_GUILD_FAIL_IN_MC_WAR_FIGHT)         // 15-7
	ErrLeaveGuildFailAssembly        = newMsgError("leave_guild 参与集结不能退出联盟", ERR_LEAVE_GUILD_FAIL_ASSEMBLY)                 // 15-8
	ErrLeaveGuildFailServerError     = newMsgError("leave_guild 服务器忙，请稍后再试", ERR_LEAVE_GUILD_FAIL_SERVER_ERROR)             // 15-3
)

// kick_other
var (
	ErrKickOtherFailNotInGuild       = newMsgError("kick_other 你没有联盟", ERR_KICK_OTHER_FAIL_NOT_IN_GUILD)                  // 19-1
	ErrKickOtherFailNpc              = newMsgError("kick_other Npc联盟不允许操作", ERR_KICK_OTHER_FAIL_NPC)                      // 19-5
	ErrKickOtherFailDeny             = newMsgError("kick_other 你没有权限操作", ERR_KICK_OTHER_FAIL_DENY)                        // 19-2
	ErrKickOtherFailImpeachLeader    = newMsgError("kick_other 盟主弹劾期间不能踢人", ERR_KICK_OTHER_FAIL_IMPEACH_LEADER)           // 19-7
	ErrKickOtherFailLimit            = newMsgError("kick_other 超出每日踢人上限", ERR_KICK_OTHER_FAIL_LIMIT)                      // 19-6
	ErrKickOtherFailTargetNotInGuild = newMsgError("kick_other 目标不在你的联盟", ERR_KICK_OTHER_FAIL_TARGET_NOT_IN_GUILD)        // 19-3
	ErrKickOtherFailXiongNuDefender  = newMsgError("kick_other 匈奴入侵防守队员，且当前活动已开启", ERR_KICK_OTHER_FAIL_XIONG_NU_DEFENDER) // 19-8
	ErrKickOtherFailInMcWarFight     = newMsgError("kick_other 名城战战斗阶段不能踢人", ERR_KICK_OTHER_FAIL_IN_MC_WAR_FIGHT)         // 19-9
	ErrKickOtherFailAssembly         = newMsgError("kick_other 参与集结成员不能踢出联盟", ERR_KICK_OTHER_FAIL_ASSEMBLY)               // 19-10
	ErrKickOtherFailServerError      = newMsgError("kick_other 服务器忙，请稍后再试", ERR_KICK_OTHER_FAIL_SERVER_ERROR)             // 19-4
)

// update_text
var (
	ErrUpdateTextFailTextTooLong    = newMsgError("update_text 内容太长", ERR_UPDATE_TEXT_FAIL_TEXT_TOO_LONG)      // 22-4
	ErrUpdateTextFailNotInGuild     = newMsgError("update_text 你没有联盟", ERR_UPDATE_TEXT_FAIL_NOT_IN_GUILD)      // 22-1
	ErrUpdateTextFailDeny           = newMsgError("update_text 你没有权限操作", ERR_UPDATE_TEXT_FAIL_DENY)            // 22-2
	ErrUpdateTextFailNpc            = newMsgError("update_text Npc联盟不允许操作", ERR_UPDATE_TEXT_FAIL_NPC)          // 22-5
	ErrUpdateTextFailSensitiveWords = newMsgError("update_text 输入包含敏感词", ERR_UPDATE_TEXT_FAIL_SENSITIVE_WORDS) // 22-6
	ErrUpdateTextFailServerError    = newMsgError("update_text 服务器忙，请稍后再试", ERR_UPDATE_TEXT_FAIL_SERVER_ERROR) // 22-3
)

// update_internal_text
var (
	ErrUpdateInternalTextFailTextTooLong    = newMsgError("update_internal_text 内容太长", ERR_UPDATE_INTERNAL_TEXT_FAIL_TEXT_TOO_LONG)      // 67-1
	ErrUpdateInternalTextFailNotInGuild     = newMsgError("update_internal_text 你没有联盟", ERR_UPDATE_INTERNAL_TEXT_FAIL_NOT_IN_GUILD)      // 67-2
	ErrUpdateInternalTextFailDeny           = newMsgError("update_internal_text 你没有权限操作", ERR_UPDATE_INTERNAL_TEXT_FAIL_DENY)            // 67-3
	ErrUpdateInternalTextFailNpc            = newMsgError("update_internal_text Npc联盟不允许操作", ERR_UPDATE_INTERNAL_TEXT_FAIL_NPC)          // 67-5
	ErrUpdateInternalTextFailSensitiveWords = newMsgError("update_internal_text 输入包含敏感词", ERR_UPDATE_INTERNAL_TEXT_FAIL_SENSITIVE_WORDS) // 67-6
	ErrUpdateInternalTextFailServerError    = newMsgError("update_internal_text 服务器忙，请稍后再试", ERR_UPDATE_INTERNAL_TEXT_FAIL_SERVER_ERROR) // 67-4
)

// update_class_names
var (
	ErrUpdateClassNamesFailNotInGuild       = newMsgError("update_class_names 你没有联盟", ERR_UPDATE_CLASS_NAMES_FAIL_NOT_IN_GUILD)             // 25-1
	ErrUpdateClassNamesFailDeny             = newMsgError("update_class_names 你没有权限操作", ERR_UPDATE_CLASS_NAMES_FAIL_DENY)                   // 25-2
	ErrUpdateClassNamesFailInvalidCount     = newMsgError("update_class_names 阶级个数无效", ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_COUNT)           // 25-3
	ErrUpdateClassNamesFailInvalidDuplicate = newMsgError("update_class_names 阶级名称无效（空或重名）", ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_DUPLICATE) // 25-5
	ErrUpdateClassNamesFailNpc              = newMsgError("update_class_names Npc联盟不允许操作", ERR_UPDATE_CLASS_NAMES_FAIL_NPC)                 // 25-6
	ErrUpdateClassNamesFailSensitiveWords   = newMsgError("update_class_names 输入包含敏感词", ERR_UPDATE_CLASS_NAMES_FAIL_SENSITIVE_WORDS)        // 25-7
	ErrUpdateClassNamesFailServerError      = newMsgError("update_class_names 服务器忙，请稍后再试", ERR_UPDATE_CLASS_NAMES_FAIL_SERVER_ERROR)        // 25-4
)

// update_class_title
var (
	ErrUpdateClassTitleFailInvalidProto    = newMsgError("update_class_title 无效的proto", ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_PROTO)      // 124-5
	ErrUpdateClassTitleFailInvalidTitleId  = newMsgError("update_class_title 无效的系统职称id", ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_TITLE_ID)  // 124-6
	ErrUpdateClassTitleFailInvalidMemberId = newMsgError("update_class_title 无效的联盟成员id", ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_MEMBER_ID) // 124-7
	ErrUpdateClassTitleFailNotInGuild      = newMsgError("update_class_title 不在联盟中", ERR_UPDATE_CLASS_TITLE_FAIL_NOT_IN_GUILD)          // 124-1
	ErrUpdateClassTitleFailDeny            = newMsgError("update_class_title 没有权限", ERR_UPDATE_CLASS_TITLE_FAIL_DENY)                   // 124-2
	ErrUpdateClassTitleFailCountLimit      = newMsgError("update_class_title 自定义职称个数无效", ERR_UPDATE_CLASS_TITLE_FAIL_COUNT_LIMIT)       // 124-8
	ErrUpdateClassTitleFailNameExist       = newMsgError("update_class_title 职称名字已经被使用了", ERR_UPDATE_CLASS_TITLE_FAIL_NAME_EXIST)       // 124-3
	ErrUpdateClassTitleFailNpc             = newMsgError("update_class_title Npc联盟不允许操作", ERR_UPDATE_CLASS_TITLE_FAIL_NPC)              // 124-9
	ErrUpdateClassTitleFailSensitiveWords  = newMsgError("update_class_title 输入包含敏感词", ERR_UPDATE_CLASS_TITLE_FAIL_SENSITIVE_WORDS)     // 124-10
	ErrUpdateClassTitleFailServerError     = newMsgError("update_class_title 服务器忙，请稍后再试", ERR_UPDATE_CLASS_TITLE_FAIL_SERVER_ERROR)     // 124-4
)

// update_flag_type
var (
	ErrUpdateFlagTypeFailNotInGuild  = newMsgError("update_flag_type 你没有联盟", ERR_UPDATE_FLAG_TYPE_FAIL_NOT_IN_GUILD)      // 28-1
	ErrUpdateFlagTypeFailDeny        = newMsgError("update_flag_type 你没有权限操作", ERR_UPDATE_FLAG_TYPE_FAIL_DENY)            // 28-2
	ErrUpdateFlagTypeFailInvalidType = newMsgError("update_flag_type 旗帜类型无效", ERR_UPDATE_FLAG_TYPE_FAIL_INVALID_TYPE)     // 28-3
	ErrUpdateFlagTypeFailNpc         = newMsgError("update_flag_type Npc联盟不允许操作", ERR_UPDATE_FLAG_TYPE_FAIL_NPC)          // 28-5
	ErrUpdateFlagTypeFailServerError = newMsgError("update_flag_type 服务器忙，请稍后再试", ERR_UPDATE_FLAG_TYPE_FAIL_SERVER_ERROR) // 28-4
)

// update_member_class_level
var (
	ErrUpdateMemberClassLevelFailInvalidClassLevel    = newMsgError("update_member_class_level 无效的阶级", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_INVALID_CLASS_LEVEL)          // 31-7
	ErrUpdateMemberClassLevelFailNotInGuild           = newMsgError("update_member_class_level 你没有联盟", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NOT_IN_GUILD)                 // 31-1
	ErrUpdateMemberClassLevelFailDeny                 = newMsgError("update_member_class_level 你没有权限操作", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_DENY)                       // 31-2
	ErrUpdateMemberClassLevelFailTargetNotInGuild     = newMsgError("update_member_class_level 目标不在你的联盟", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_TARGET_NOT_IN_GUILD)       // 31-3
	ErrUpdateMemberClassLevelFailTargetSameClassLevel = newMsgError("update_member_class_level 目标当前就是这个阶级", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_TARGET_SAME_CLASS_LEVEL) // 31-8
	ErrUpdateMemberClassLevelFailDenyTarget           = newMsgError("update_member_class_level 权限不足（目标权限不比你低）", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_DENY_TARGET)         // 31-4
	ErrUpdateMemberClassLevelFailClassFull            = newMsgError("update_member_class_level 目标阶级已经满员", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_CLASS_FULL)                // 31-5
	ErrUpdateMemberClassLevelFailNpc                  = newMsgError("update_member_class_level Npc联盟不允许操作", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NPC)                     // 31-9
	ErrUpdateMemberClassLevelFailServerError          = newMsgError("update_member_class_level 服务器忙，请稍后再试", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_SERVER_ERROR)            // 31-6
	ErrUpdateMemberClassLevelFailChangeKingInCd       = newMsgError("update_member_class_level 禅让国王，在CD中", ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_CHANGE_KING_IN_CD)        // 31-10
)

// cancel_change_leader
var (
	ErrCancelChangeLeaderFailNotInGuild      = newMsgError("cancel_change_leader 你没有联盟", ERR_CANCEL_CHANGE_LEADER_FAIL_NOT_IN_GUILD)          // 82-1
	ErrCancelChangeLeaderFailDeny            = newMsgError("cancel_change_leader 你不是盟主，不能取消", ERR_CANCEL_CHANGE_LEADER_FAIL_DENY)             // 82-2
	ErrCancelChangeLeaderFailNotChangeLeader = newMsgError("cancel_change_leader 当前没有禅让倒计时", ERR_CANCEL_CHANGE_LEADER_FAIL_NOT_CHANGE_LEADER) // 82-3
	ErrCancelChangeLeaderFailServerError     = newMsgError("cancel_change_leader 服务器忙，请稍后再试", ERR_CANCEL_CHANGE_LEADER_FAIL_SERVER_ERROR)     // 82-4
)

// update_join_condition
var (
	ErrUpdateJoinConditionFailNotInGuild          = newMsgError("update_join_condition 你没有联盟", ERR_UPDATE_JOIN_CONDITION_FAIL_NOT_IN_GUILD)               // 70-1
	ErrUpdateJoinConditionFailDeny                = newMsgError("update_join_condition 你没有权限操作", ERR_UPDATE_JOIN_CONDITION_FAIL_DENY)                     // 70-2
	ErrUpdateJoinConditionFailInvalidHeroLevel    = newMsgError("update_join_condition 无效的君主等级", ERR_UPDATE_JOIN_CONDITION_FAIL_INVALID_HERO_LEVEL)       // 70-4
	ErrUpdateJoinConditionFailInvalidJunXianLevel = newMsgError("update_join_condition 无效的百战军衔等级", ERR_UPDATE_JOIN_CONDITION_FAIL_INVALID_JUN_XIAN_LEVEL) // 70-5
	ErrUpdateJoinConditionFailNpc                 = newMsgError("update_join_condition Npc联盟不允许操作", ERR_UPDATE_JOIN_CONDITION_FAIL_NPC)                   // 70-6
	ErrUpdateJoinConditionFailServerError         = newMsgError("update_join_condition 服务器忙，请稍后再试", ERR_UPDATE_JOIN_CONDITION_FAIL_SERVER_ERROR)          // 70-3
)

// update_guild_name
var (
	ErrUpdateGuildNameFailInvalidName     = newMsgError("update_guild_name 无效的联盟名字", ERR_UPDATE_GUILD_NAME_FAIL_INVALID_NAME)      // 73-1
	ErrUpdateGuildNameFailInvalidFlagName = newMsgError("update_guild_name 无效的联盟旗号", ERR_UPDATE_GUILD_NAME_FAIL_INVALID_FLAG_NAME) // 73-2
	ErrUpdateGuildNameFailNotInGuild      = newMsgError("update_guild_name 你没有联盟", ERR_UPDATE_GUILD_NAME_FAIL_NOT_IN_GUILD)        // 73-3
	ErrUpdateGuildNameFailDeny            = newMsgError("update_guild_name 你没有权限操作", ERR_UPDATE_GUILD_NAME_FAIL_DENY)              // 73-4
	ErrUpdateGuildNameFailExistName       = newMsgError("update_guild_name 联盟名字已存在", ERR_UPDATE_GUILD_NAME_FAIL_EXIST_NAME)        // 73-5
	ErrUpdateGuildNameFailExistFlagName   = newMsgError("update_guild_name 联盟旗号已存在", ERR_UPDATE_GUILD_NAME_FAIL_EXIST_FLAG_NAME)   // 73-6
	ErrUpdateGuildNameFailCostNotEnough   = newMsgError("update_guild_name 改名消耗不足", ERR_UPDATE_GUILD_NAME_FAIL_COST_NOT_ENOUGH)    // 73-7
	ErrUpdateGuildNameFailCooldown        = newMsgError("update_guild_name 改名CD中", ERR_UPDATE_GUILD_NAME_FAIL_COOLDOWN)            // 73-8
	ErrUpdateGuildNameFailNpc             = newMsgError("update_guild_name Npc联盟不允许操作", ERR_UPDATE_GUILD_NAME_FAIL_NPC)            // 73-10
	ErrUpdateGuildNameFailSensitiveWords  = newMsgError("update_guild_name 输入包含敏感词", ERR_UPDATE_GUILD_NAME_FAIL_SENSITIVE_WORDS)   // 73-11
	ErrUpdateGuildNameFailServerError     = newMsgError("update_guild_name 服务器忙，请稍后再试", ERR_UPDATE_GUILD_NAME_FAIL_SERVER_ERROR)   // 73-9
)

// update_guild_label
var (
	ErrUpdateGuildLabelFailNotInGuild     = newMsgError("update_guild_label 你没有联盟", ERR_UPDATE_GUILD_LABEL_FAIL_NOT_IN_GUILD)      // 77-1
	ErrUpdateGuildLabelFailDeny           = newMsgError("update_guild_label 你没有权限操作", ERR_UPDATE_GUILD_LABEL_FAIL_DENY)            // 77-2
	ErrUpdateGuildLabelFailCountLimit     = newMsgError("update_guild_label 标签个数超出上限", ERR_UPDATE_GUILD_LABEL_FAIL_COUNT_LIMIT)    // 77-3
	ErrUpdateGuildLabelFailCharLimit      = newMsgError("update_guild_label 标签字数超出上限", ERR_UPDATE_GUILD_LABEL_FAIL_CHAR_LIMIT)     // 77-7
	ErrUpdateGuildLabelFailDuplicate      = newMsgError("update_guild_label 标签重名", ERR_UPDATE_GUILD_LABEL_FAIL_DUPLICATE)          // 77-5
	ErrUpdateGuildLabelFailNpc            = newMsgError("update_guild_label Npc联盟不允许操作", ERR_UPDATE_GUILD_LABEL_FAIL_NPC)          // 77-6
	ErrUpdateGuildLabelFailSensitiveWords = newMsgError("update_guild_label 输入包含敏感词", ERR_UPDATE_GUILD_LABEL_FAIL_SENSITIVE_WORDS) // 77-8
	ErrUpdateGuildLabelFailServerError    = newMsgError("update_guild_label 服务器忙，请稍后再试", ERR_UPDATE_GUILD_LABEL_FAIL_SERVER_ERROR) // 77-4
)

// donate
var (
	ErrDonateFailNotInGuild       = newMsgError("donate 你没有联盟", ERR_DONATE_FAIL_NOT_IN_GUILD)                   // 85-4
	ErrDonateFailInvalidSequence  = newMsgError("donate 无效的序号", ERR_DONATE_FAIL_INVALID_SEQUENCE)               // 85-1
	ErrDonateFailMaxTimes         = newMsgError("donate 已经达到最大捐献次数", ERR_DONATE_FAIL_MAX_TIMES)                 // 85-2
	ErrDonateFailCostNotEnough    = newMsgError("donate 消耗不足", ERR_DONATE_FAIL_COST_NOT_ENOUGH)                 // 85-3
	ErrDonateFailDonateTimesLimit = newMsgError("donate 捐献次数已经达到外使院捐献次数上限", ERR_DONATE_FAIL_DONATE_TIMES_LIMIT) // 85-6
	ErrDonateFailLevelNotEnough   = newMsgError("donate 君主等级不够，无法捐献", ERR_DONATE_FAIL_LEVEL_NOT_ENOUGH)         // 85-7
	ErrDonateFailServerError      = newMsgError("donate 服务器忙，请稍后再试", ERR_DONATE_FAIL_SERVER_ERROR)              // 85-5
)

// upgrade_level
var (
	ErrUpgradeLevelFailNotInGuild    = newMsgError("upgrade_level 你没有联盟", ERR_UPGRADE_LEVEL_FAIL_NOT_IN_GUILD)      // 92-1
	ErrUpgradeLevelFailDeny          = newMsgError("upgrade_level 你没有权限操作", ERR_UPGRADE_LEVEL_FAIL_DENY)            // 92-2
	ErrUpgradeLevelFailUpgrading     = newMsgError("upgrade_level 正在升级中", ERR_UPGRADE_LEVEL_FAIL_UPGRADING)         // 92-3
	ErrUpgradeLevelFailMaxLevel      = newMsgError("upgrade_level 帮派已经达到最高级", ERR_UPGRADE_LEVEL_FAIL_MAX_LEVEL)     // 92-4
	ErrUpgradeLevelFailCostNotEnough = newMsgError("upgrade_level 建设值不足", ERR_UPGRADE_LEVEL_FAIL_COST_NOT_ENOUGH)   // 92-6
	ErrUpgradeLevelFailNpc           = newMsgError("upgrade_level Npc联盟不允许操作", ERR_UPGRADE_LEVEL_FAIL_NPC)          // 92-7
	ErrUpgradeLevelFailServerError   = newMsgError("upgrade_level 服务器忙，请稍后再试", ERR_UPGRADE_LEVEL_FAIL_SERVER_ERROR) // 92-5
)

// reduce_upgrade_level_cd
var (
	ErrReduceUpgradeLevelCdFailNotInGuild    = newMsgError("reduce_upgrade_level_cd 你没有联盟", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NOT_IN_GUILD)        // 95-1
	ErrReduceUpgradeLevelCdFailDeny          = newMsgError("reduce_upgrade_level_cd 你没有权限操作", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_DENY)              // 95-2
	ErrReduceUpgradeLevelCdFailNoUpgrading   = newMsgError("reduce_upgrade_level_cd 联盟没有在升级，不能加速", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NO_UPGRADING) // 95-6
	ErrReduceUpgradeLevelCdFailMaxTimes      = newMsgError("reduce_upgrade_level_cd 帮派已经达到最大加速次数", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_MAX_TIMES)    // 95-7
	ErrReduceUpgradeLevelCdFailCostNotEnough = newMsgError("reduce_upgrade_level_cd 建设值不足", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_COST_NOT_ENOUGH)     // 95-8
	ErrReduceUpgradeLevelCdFailNpc           = newMsgError("reduce_upgrade_level_cd Npc联盟不允许操作", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NPC)            // 95-9
	ErrReduceUpgradeLevelCdFailServerError   = newMsgError("reduce_upgrade_level_cd 服务器忙，请稍后再试", ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_SERVER_ERROR)   // 95-5
)

// impeach_leader
var (
	ErrImpeachLeaderFailNotInGuild        = newMsgError("impeach_leader 你没有联盟", ERR_IMPEACH_LEADER_FAIL_NOT_IN_GUILD)            // 98-1
	ErrImpeachLeaderFailConditionNotReach = newMsgError("impeach_leader 弹劾条件未满足", ERR_IMPEACH_LEADER_FAIL_CONDITION_NOT_REACH)   // 98-2
	ErrImpeachLeaderFailDeny              = newMsgError("impeach_leader 你没有权限操作", ERR_IMPEACH_LEADER_FAIL_DENY)                  // 98-3
	ErrImpeachLeaderFailImpeachExist      = newMsgError("impeach_leader 已经存在弹劾盟主", ERR_IMPEACH_LEADER_FAIL_IMPEACH_EXIST)        // 98-5
	ErrImpeachLeaderFailChangingLeader    = newMsgError("impeach_leader 正在禅让盟主，不允许弹劾", ERR_IMPEACH_LEADER_FAIL_CHANGING_LEADER)  // 98-7
	ErrImpeachLeaderFailInvalidTime       = newMsgError("impeach_leader 今日弹劾盟主时间已过，请明日再试", ERR_IMPEACH_LEADER_FAIL_INVALID_TIME) // 98-6
	ErrImpeachLeaderFailServerError       = newMsgError("impeach_leader 服务器忙，请稍后再试", ERR_IMPEACH_LEADER_FAIL_SERVER_ERROR)       // 98-4
)

// impeach_leader_vote
var (
	ErrImpeachLeaderVoteFailNotInGuild      = newMsgError("impeach_leader_vote 你没有联盟", ERR_IMPEACH_LEADER_VOTE_FAIL_NOT_IN_GUILD)         // 101-1
	ErrImpeachLeaderVoteFailInvalidTarget   = newMsgError("impeach_leader_vote 无效的投票目标", ERR_IMPEACH_LEADER_VOTE_FAIL_INVALID_TARGET)     // 101-4
	ErrImpeachLeaderVoteFailImpeachNotExist = newMsgError("impeach_leader_vote 当前没有弹劾盟主", ERR_IMPEACH_LEADER_VOTE_FAIL_IMPEACH_NOT_EXIST) // 101-2
	ErrImpeachLeaderVoteFailOldLeader       = newMsgError("impeach_leader_vote 无法投票给原盟主", ERR_IMPEACH_LEADER_VOTE_FAIL_OLD_LEADER)        // 101-5
	ErrImpeachLeaderVoteFailServerError     = newMsgError("impeach_leader_vote 服务器忙，请稍后再试", ERR_IMPEACH_LEADER_VOTE_FAIL_SERVER_ERROR)    // 101-3
)

// list_guild_by_ids
var (
	ErrListGuildByIdsFailInvalidId    = newMsgError("list_guild_by_ids 无效的id", ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_ID)        // 104-2
	ErrListGuildByIdsFailInvalidCount = newMsgError("list_guild_by_ids 无效的id个数", ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_COUNT)   // 104-3
	ErrListGuildByIdsFailServerError  = newMsgError("list_guild_by_ids 服务器忙，请稍后再试", ERR_LIST_GUILD_BY_IDS_FAIL_SERVER_ERROR) // 104-1
)

// user_request_join
var (
	ErrUserRequestJoinFailInvalidId       = newMsgError("user_request_join 无效的联盟id", ERR_USER_REQUEST_JOIN_FAIL_INVALID_ID)                  // 42-2
	ErrUserRequestJoinFailSelfFull        = newMsgError("user_request_join 已达申请上限，请取消其他申请", ERR_USER_REQUEST_JOIN_FAIL_SELF_FULL)            // 42-4
	ErrUserRequestJoinFailSelfGuild       = newMsgError("user_request_join 不能申请加入自己的联盟", ERR_USER_REQUEST_JOIN_FAIL_SELF_GUILD)              // 42-7
	ErrUserRequestJoinFailCondition       = newMsgError("user_request_join 要求未达到", ERR_USER_REQUEST_JOIN_FAIL_CONDITION)                     // 42-5
	ErrUserRequestJoinFailFull            = newMsgError("user_request_join 申请的帮派已经满员", ERR_USER_REQUEST_JOIN_FAIL_FULL)                      // 42-3
	ErrUserRequestJoinFailDuplicate       = newMsgError("user_request_join 这个联盟已经申请了，不要重复申请", ERR_USER_REQUEST_JOIN_FAIL_DUPLICATE)          // 42-8
	ErrUserRequestJoinFailLeader          = newMsgError("user_request_join 盟主不能申请加入其它帮派", ERR_USER_REQUEST_JOIN_FAIL_LEADER)                 // 42-9
	ErrUserRequestJoinFailNpc             = newMsgError("user_request_join 这个联盟是纯Npc联盟，不允许操作", ERR_USER_REQUEST_JOIN_FAIL_NPC)               // 42-10
	ErrUserRequestJoinFailXiongNuDefender = newMsgError("user_request_join 匈奴入侵防守队员，且当前活动已开启", ERR_USER_REQUEST_JOIN_FAIL_XIONG_NU_DEFENDER) // 42-11
	ErrUserRequestJoinFailLeaveCd         = newMsgError("user_request_join 离开此联盟不足4小时，不能加入", ERR_USER_REQUEST_JOIN_FAIL_LEAVE_CD)            // 42-12
	ErrUserRequestJoinFailCountry         = newMsgError("user_request_join 不能加入其他国家的联盟", ERR_USER_REQUEST_JOIN_FAIL_COUNTRY)                 // 42-13
	ErrUserRequestJoinFailServerError     = newMsgError("user_request_join 服务器忙，请稍后再试", ERR_USER_REQUEST_JOIN_FAIL_SERVER_ERROR)             // 42-6
)

// user_cancel_join_request
var (
	ErrUserCancelJoinRequestFailInTheGuild  = newMsgError("user_cancel_join_request 在联盟中", ERR_USER_CANCEL_JOIN_REQUEST_FAIL_IN_THE_GUILD)       // 45-1
	ErrUserCancelJoinRequestFailInvalidId   = newMsgError("user_cancel_join_request 无效的id", ERR_USER_CANCEL_JOIN_REQUEST_FAIL_INVALID_ID)        // 45-2
	ErrUserCancelJoinRequestFailServerError = newMsgError("user_cancel_join_request 服务器忙，请稍后再试", ERR_USER_CANCEL_JOIN_REQUEST_FAIL_SERVER_ERROR) // 45-3
)

// guild_reply_join_request
var (
	ErrGuildReplyJoinRequestFailNoGuild         = newMsgError("guild_reply_join_request 不在联盟中", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_NO_GUILD)                      // 57-1
	ErrGuildReplyJoinRequestFailDeny            = newMsgError("guild_reply_join_request 没有权限", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_DENY)                           // 57-2
	ErrGuildReplyJoinRequestFailInvalidRequest  = newMsgError("guild_reply_join_request 玩家已经取消了申请", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_INVALID_REQUEST)           // 57-4
	ErrGuildReplyJoinRequestFailXiongNuDefender = newMsgError("guild_reply_join_request 匈奴入侵防守队员，且当前活动已开启", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_XIONG_NU_DEFENDER) // 57-5
	ErrGuildReplyJoinRequestFailLeaveCd         = newMsgError("guild_reply_join_request 玩家离开本联盟不足4小时，不能加入", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_LEAVE_CD)          // 57-6
	ErrGuildReplyJoinRequestFailServerError     = newMsgError("guild_reply_join_request 服务器忙，请稍后再试", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_SERVER_ERROR)             // 57-3
	ErrGuildReplyJoinRequestFailFullMember      = newMsgError("guild_reply_join_request 联盟已经满员", ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_FULL_MEMBER)                  // 57-7
)

// guild_invate_other
var (
	ErrGuildInvateOtherFailInvalidId   = newMsgError("guild_invate_other 无效的玩家id", ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID)          // 111-8
	ErrGuildInvateOtherFailNotInGuild  = newMsgError("guild_invate_other 不在联盟中", ERR_GUILD_INVATE_OTHER_FAIL_NOT_IN_GUILD)          // 111-7
	ErrGuildInvateOtherFailDeny        = newMsgError("guild_invate_other 没有权限", ERR_GUILD_INVATE_OTHER_FAIL_DENY)                   // 111-2
	ErrGuildInvateOtherFailInvated     = newMsgError("guild_invate_other 目标已经在邀请队列", ERR_GUILD_INVATE_OTHER_FAIL_INVATED)           // 111-5
	ErrGuildInvateOtherFailFull        = newMsgError("guild_invate_other 邀请列表已满，取消之前申请", ERR_GUILD_INVATE_OTHER_FAIL_FULL)          // 111-6
	ErrGuildInvateOtherFailGuildMember = newMsgError("guild_invate_other 邀请的玩家已经在自己的联盟中", ERR_GUILD_INVATE_OTHER_FAIL_GUILD_MEMBER) // 111-3
	ErrGuildInvateOtherFailServerError = newMsgError("guild_invate_other 服务器忙，请稍后再试", ERR_GUILD_INVATE_OTHER_FAIL_SERVER_ERROR)     // 111-4
	ErrGuildInvateOtherFailFullMember  = newMsgError("guild_invate_other 联盟已经满员", ERR_GUILD_INVATE_OTHER_FAIL_FULL_MEMBER)          // 111-10
)

// guild_cancel_invate_other
var (
	ErrGuildCancelInvateOtherFailInvalidId   = newMsgError("guild_cancel_invate_other 无效的玩家id", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_INVALID_ID)          // 114-7
	ErrGuildCancelInvateOtherFailNotInGuild  = newMsgError("guild_cancel_invate_other 不在联盟中", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_NOT_IN_GUILD)          // 114-6
	ErrGuildCancelInvateOtherFailDeny        = newMsgError("guild_cancel_invate_other 没有权限", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_DENY)                   // 114-2
	ErrGuildCancelInvateOtherFailGuildMember = newMsgError("guild_cancel_invate_other 邀请的玩家已经在自己的联盟中", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_GUILD_MEMBER) // 114-3
	ErrGuildCancelInvateOtherFailIdNotExist  = newMsgError("guild_cancel_invate_other 玩家不在邀请列表中", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_ID_NOT_EXIST)      // 114-4
	ErrGuildCancelInvateOtherFailServerError = newMsgError("guild_cancel_invate_other 服务器忙，请稍后再试", ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_SERVER_ERROR)     // 114-5
)

// user_reply_invate_request
var (
	ErrUserReplyInvateRequestFailInvalidId       = newMsgError("user_reply_invate_request 无效的id", ERR_USER_REPLY_INVATE_REQUEST_FAIL_INVALID_ID)                    // 50-2
	ErrUserReplyInvateRequestFailLeader          = newMsgError("user_reply_invate_request 你是盟主，请先卸任盟主再接受邀请", ERR_USER_REPLY_INVATE_REQUEST_FAIL_LEADER)             // 50-4
	ErrUserReplyInvateRequestFailXiongNuDefender = newMsgError("user_reply_invate_request 匈奴入侵防守队员，且当前活动已开启", ERR_USER_REPLY_INVATE_REQUEST_FAIL_XIONG_NU_DEFENDER) // 50-5
	ErrUserReplyInvateRequestFailLeaveCd         = newMsgError("user_reply_invate_request 离开此联盟不足4小时，不能加入", ERR_USER_REPLY_INVATE_REQUEST_FAIL_LEAVE_CD)            // 50-6
	ErrUserReplyInvateRequestFailServerError     = newMsgError("user_reply_invate_request 服务器忙，请稍后再试", ERR_USER_REPLY_INVATE_REQUEST_FAIL_SERVER_ERROR)             // 50-3
	ErrUserReplyInvateRequestFailFullMember      = newMsgError("user_reply_invate_request 联盟已满员", ERR_USER_REPLY_INVATE_REQUEST_FAIL_FULL_MEMBER)                   // 50-7
)

// list_invite_me_guild
var (
	ErrListInviteMeGuildFailServerError = newMsgError("list_invite_me_guild 服务器忙，请稍后再试", ERR_LIST_INVITE_ME_GUILD_FAIL_SERVER_ERROR) // 195-2
)

// update_friend_guild
var (
	ErrUpdateFriendGuildFailTextTooLong    = newMsgError("update_friend_guild 内容太长", ERR_UPDATE_FRIEND_GUILD_FAIL_TEXT_TOO_LONG)      // 127-1
	ErrUpdateFriendGuildFailNotInGuild     = newMsgError("update_friend_guild 你没有联盟", ERR_UPDATE_FRIEND_GUILD_FAIL_NOT_IN_GUILD)      // 127-2
	ErrUpdateFriendGuildFailDeny           = newMsgError("update_friend_guild 你没有权限操作", ERR_UPDATE_FRIEND_GUILD_FAIL_DENY)            // 127-3
	ErrUpdateFriendGuildFailNpc            = newMsgError("update_friend_guild Npc联盟不允许操作", ERR_UPDATE_FRIEND_GUILD_FAIL_NPC)          // 127-4
	ErrUpdateFriendGuildFailSensitiveWords = newMsgError("update_friend_guild 输入包含敏感词", ERR_UPDATE_FRIEND_GUILD_FAIL_SENSITIVE_WORDS) // 127-6
	ErrUpdateFriendGuildFailServerError    = newMsgError("update_friend_guild 服务器忙，请稍后再试", ERR_UPDATE_FRIEND_GUILD_FAIL_SERVER_ERROR) // 127-5
)

// update_enemy_guild
var (
	ErrUpdateEnemyGuildFailTextTooLong    = newMsgError("update_enemy_guild 内容太长", ERR_UPDATE_ENEMY_GUILD_FAIL_TEXT_TOO_LONG)      // 130-1
	ErrUpdateEnemyGuildFailNotInGuild     = newMsgError("update_enemy_guild 你没有联盟", ERR_UPDATE_ENEMY_GUILD_FAIL_NOT_IN_GUILD)      // 130-2
	ErrUpdateEnemyGuildFailDeny           = newMsgError("update_enemy_guild 你没有权限操作", ERR_UPDATE_ENEMY_GUILD_FAIL_DENY)            // 130-3
	ErrUpdateEnemyGuildFailNpc            = newMsgError("update_enemy_guild Npc联盟不允许操作", ERR_UPDATE_ENEMY_GUILD_FAIL_NPC)          // 130-4
	ErrUpdateEnemyGuildFailSensitiveWords = newMsgError("update_enemy_guild 输入包含敏感词", ERR_UPDATE_ENEMY_GUILD_FAIL_SENSITIVE_WORDS) // 130-6
	ErrUpdateEnemyGuildFailServerError    = newMsgError("update_enemy_guild 服务器忙，请稍后再试", ERR_UPDATE_ENEMY_GUILD_FAIL_SERVER_ERROR) // 130-5
)

// update_guild_prestige
var (
	ErrUpdateGuildPrestigeFailTargetNotFound = newMsgError("update_guild_prestige 声望目标没找到", ERR_UPDATE_GUILD_PRESTIGE_FAIL_TARGET_NOT_FOUND)  // 133-1
	ErrUpdateGuildPrestigeFailNotInGuild     = newMsgError("update_guild_prestige 你没有联盟", ERR_UPDATE_GUILD_PRESTIGE_FAIL_NOT_IN_GUILD)        // 133-2
	ErrUpdateGuildPrestigeFailDeny           = newMsgError("update_guild_prestige 你没有权限操作", ERR_UPDATE_GUILD_PRESTIGE_FAIL_DENY)              // 133-3
	ErrUpdateGuildPrestigeFailNpc            = newMsgError("update_guild_prestige Npc联盟不允许操作", ERR_UPDATE_GUILD_PRESTIGE_FAIL_NPC)            // 133-4
	ErrUpdateGuildPrestigeFailCostNotEnough  = newMsgError("update_guild_prestige 消耗不足", ERR_UPDATE_GUILD_PRESTIGE_FAIL_COST_NOT_ENOUGH)      // 133-6
	ErrUpdateGuildPrestigeFailCountdown      = newMsgError("update_guild_prestige 修改倒计时", ERR_UPDATE_GUILD_PRESTIGE_FAIL_COUNTDOWN)           // 133-7
	ErrUpdateGuildPrestigeFailSameTarget     = newMsgError("update_guild_prestige 修改的目标跟现在的目标一样", ERR_UPDATE_GUILD_PRESTIGE_FAIL_SAME_TARGET) // 133-8
	ErrUpdateGuildPrestigeFailServerError    = newMsgError("update_guild_prestige 服务器忙，请稍后再试", ERR_UPDATE_GUILD_PRESTIGE_FAIL_SERVER_ERROR)   // 133-5
)

// place_guild_statue
var (
	ErrPlaceGuildStatueFailNoGuild     = newMsgError("place_guild_statue 没有联盟", ERR_PLACE_GUILD_STATUE_FAIL_NO_GUILD)           // 136-1
	ErrPlaceGuildStatueFailNotLeader   = newMsgError("place_guild_statue 不是盟主", ERR_PLACE_GUILD_STATUE_FAIL_NOT_LEADER)         // 136-2
	ErrPlaceGuildStatueFailHasPlaced   = newMsgError("place_guild_statue 有放置了，请先取回", ERR_PLACE_GUILD_STATUE_FAIL_HAS_PLACED)    // 136-3
	ErrPlaceGuildStatueFailMapNotFound = newMsgError("place_guild_statue 地图没找到", ERR_PLACE_GUILD_STATUE_FAIL_MAP_NOT_FOUND)     // 136-4
	ErrPlaceGuildStatueFailXInvalid    = newMsgError("place_guild_statue x非法", ERR_PLACE_GUILD_STATUE_FAIL_X_INVALID)           // 136-5
	ErrPlaceGuildStatueFailYInvalid    = newMsgError("place_guild_statue y非法", ERR_PLACE_GUILD_STATUE_FAIL_Y_INVALID)           // 136-6
	ErrPlaceGuildStatueFailServerError = newMsgError("place_guild_statue 服务器忙，请稍后再试", ERR_PLACE_GUILD_STATUE_FAIL_SERVER_ERROR) // 136-7
)

// take_back_guild_statue
var (
	ErrTakeBackGuildStatueFailNoGuild     = newMsgError("take_back_guild_statue 没有联盟", ERR_TAKE_BACK_GUILD_STATUE_FAIL_NO_GUILD)           // 141-1
	ErrTakeBackGuildStatueFailNotLeader   = newMsgError("take_back_guild_statue 不是盟主", ERR_TAKE_BACK_GUILD_STATUE_FAIL_NOT_LEADER)         // 141-2
	ErrTakeBackGuildStatueFailNotPlace    = newMsgError("take_back_guild_statue 没有放置", ERR_TAKE_BACK_GUILD_STATUE_FAIL_NOT_PLACE)          // 141-3
	ErrTakeBackGuildStatueFailServerError = newMsgError("take_back_guild_statue 服务器忙，请稍后再试", ERR_TAKE_BACK_GUILD_STATUE_FAIL_SERVER_ERROR) // 141-4
)

// collect_first_join_guild_prize
var (
	ErrCollectFirstJoinGuildPrizeFailNoGuild   = newMsgError("collect_first_join_guild_prize 没有加入联盟", ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_NO_GUILD)    // 145-3
	ErrCollectFirstJoinGuildPrizeFailCollected = newMsgError("collect_first_join_guild_prize 奖励已经被领取了", ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_COLLECTED) // 145-2
)

// seek_help
var (
	ErrSeekHelpFailDisable    = newMsgError("seek_help 当前求助状态不可用", ERR_SEEK_HELP_FAIL_DISABLE)         // 149-1
	ErrSeekHelpFailNotInGuild = newMsgError("seek_help 自己不在联盟中，不能求助", ERR_SEEK_HELP_FAIL_NOT_IN_GUILD) // 149-2
	ErrSeekHelpFailWaiShiYuan = newMsgError("seek_help 自己还没有外使院建筑", ERR_SEEK_HELP_FAIL_WAI_SHI_YUAN)   // 149-3
)

// help_guild_member
var (
	ErrHelpGuildMemberFailIdNotFound = newMsgError("help_guild_member 求助id没找到", ERR_HELP_GUILD_MEMBER_FAIL_ID_NOT_FOUND)          // 153-1
	ErrHelpGuildMemberFailHelped     = newMsgError("help_guild_member 这条求助你已经帮助过了", ERR_HELP_GUILD_MEMBER_FAIL_HELPED)            // 153-2
	ErrHelpGuildMemberFailNotInGuild = newMsgError("help_guild_member 你不在联盟中，不能帮助盟友的求助", ERR_HELP_GUILD_MEMBER_FAIL_NOT_IN_GUILD) // 153-3
)

// help_all_guild_member
var (
	ErrHelpAllGuildMemberFailNotInGuild = newMsgError("help_all_guild_member 你不在联盟中，不能帮助盟友的求助", ERR_HELP_ALL_GUILD_MEMBER_FAIL_NOT_IN_GUILD) // 160-1
)

// collect_guild_event_prize
var (
	ErrCollectGuildEventPrizeFailIdNotFound = newMsgError("collect_guild_event_prize id不存在", ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_ID_NOT_FOUND)      // 165-1
	ErrCollectGuildEventPrizeFailNotInGuild = newMsgError("collect_guild_event_prize 你不在联盟中", ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_NOT_IN_GUILD)     // 165-2
	ErrCollectGuildEventPrizeFailVipLimit   = newMsgError("collect_guild_event_prize vip等级不够，不能一键领", ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_VIP_LIMIT) // 165-3
)

// collect_full_big_box
var (
	ErrCollectFullBigBoxFailNotInGuild    = newMsgError("collect_full_big_box 你不在联盟中", ERR_COLLECT_FULL_BIG_BOX_FAIL_NOT_IN_GUILD)           // 169-1
	ErrCollectFullBigBoxFailLocked        = newMsgError("collect_full_big_box 宝箱尚未解锁", ERR_COLLECT_FULL_BIG_BOX_FAIL_LOCKED)                 // 169-2
	ErrCollectFullBigBoxFailTimeNotEnough = newMsgError("collect_full_big_box 进入联盟时间不足，不能领取", ERR_COLLECT_FULL_BIG_BOX_FAIL_TIME_NOT_ENOUGH) // 169-4
	ErrCollectFullBigBoxFailServerError   = newMsgError("collect_full_big_box 服务器忙，请稍后再试", ERR_COLLECT_FULL_BIG_BOX_FAIL_SERVER_ERROR)       // 169-3
)

// upgrade_technology
var (
	ErrUpgradeTechnologyFailNotInGuild    = newMsgError("upgrade_technology 你没有联盟", ERR_UPGRADE_TECHNOLOGY_FAIL_NOT_IN_GUILD)      // 174-7
	ErrUpgradeTechnologyFailInvalidGroup  = newMsgError("upgrade_technology 无效的科技组", ERR_UPGRADE_TECHNOLOGY_FAIL_INVALID_GROUP)    // 174-1
	ErrUpgradeTechnologyFailMaxLevel      = newMsgError("upgrade_technology 科技已经达到最大等级", ERR_UPGRADE_TECHNOLOGY_FAIL_MAX_LEVEL)    // 174-2
	ErrUpgradeTechnologyFailRequired      = newMsgError("upgrade_technology 联盟等级不足", ERR_UPGRADE_TECHNOLOGY_FAIL_REQUIRED)         // 174-3
	ErrUpgradeTechnologyFailCostNotEnough = newMsgError("upgrade_technology 联盟建设值不足", ERR_UPGRADE_TECHNOLOGY_FAIL_COST_NOT_ENOUGH) // 174-4
	ErrUpgradeTechnologyFailUpgrading     = newMsgError("upgrade_technology 存在正在升级的科技", ERR_UPGRADE_TECHNOLOGY_FAIL_UPGRADING)     // 174-5
	ErrUpgradeTechnologyFailDeny          = newMsgError("upgrade_technology 你没有权限操作", ERR_UPGRADE_TECHNOLOGY_FAIL_DENY)            // 174-6
	ErrUpgradeTechnologyFailServerError   = newMsgError("upgrade_technology 服务器忙，请稍后再试", ERR_UPGRADE_TECHNOLOGY_FAIL_SERVER_ERROR) // 174-8
)

// reduce_technology_cd
var (
	ErrReduceTechnologyCdFailNotInGuild    = newMsgError("reduce_technology_cd 你没有联盟", ERR_REDUCE_TECHNOLOGY_CD_FAIL_NOT_IN_GUILD)        // 177-1
	ErrReduceTechnologyCdFailDeny          = newMsgError("reduce_technology_cd 你没有权限操作", ERR_REDUCE_TECHNOLOGY_CD_FAIL_DENY)              // 177-2
	ErrReduceTechnologyCdFailNoUpgrading   = newMsgError("reduce_technology_cd 科技没有在升级，不能加速", ERR_REDUCE_TECHNOLOGY_CD_FAIL_NO_UPGRADING) // 177-3
	ErrReduceTechnologyCdFailMaxTimes      = newMsgError("reduce_technology_cd 帮派已经达到最大加速次数", ERR_REDUCE_TECHNOLOGY_CD_FAIL_MAX_TIMES)    // 177-4
	ErrReduceTechnologyCdFailCostNotEnough = newMsgError("reduce_technology_cd 建设值不足", ERR_REDUCE_TECHNOLOGY_CD_FAIL_COST_NOT_ENOUGH)     // 177-5
	ErrReduceTechnologyCdFailNpc           = newMsgError("reduce_technology_cd Npc联盟不允许操作", ERR_REDUCE_TECHNOLOGY_CD_FAIL_NPC)            // 177-6
	ErrReduceTechnologyCdFailServerError   = newMsgError("reduce_technology_cd 服务器忙，请稍后再试", ERR_REDUCE_TECHNOLOGY_CD_FAIL_SERVER_ERROR)   // 177-7
)

// help_tech
var (
	ErrHelpTechFailNotInGuild      = newMsgError("help_tech 你没有联盟", ERR_HELP_TECH_FAIL_NOT_IN_GUILD)         // 186-1
	ErrHelpTechFailCantHelp        = newMsgError("help_tech 你不能协助", ERR_HELP_TECH_FAIL_CANT_HELP)            // 186-2
	ErrHelpTechFailNoTechUpgrading = newMsgError("help_tech 没有升级中的科技", ERR_HELP_TECH_FAIL_NO_TECH_UPGRADING) // 186-3
	ErrHelpTechFailServerError     = newMsgError("help_tech 服务器忙，请稍后再试", ERR_HELP_TECH_FAIL_SERVER_ERROR)    // 186-4
)

// recommend_invite_heros
var (
	ErrRecommendInviteHerosFailServerError = newMsgError("recommend_invite_heros 服务器忙，请稍后再试", ERR_RECOMMEND_INVITE_HEROS_FAIL_SERVER_ERROR) // 189-1
)

// search_no_guild_heros
var (
	ErrSearchNoGuildHerosFailInvalidArg  = newMsgError("search_no_guild_heros 参数错误", ERR_SEARCH_NO_GUILD_HEROS_FAIL_INVALID_ARG)        // 192-2
	ErrSearchNoGuildHerosFailServerError = newMsgError("search_no_guild_heros 服务器忙，请稍后再试", ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR) // 192-1
	ErrSearchNoGuildHerosFailTooFast     = newMsgError("search_no_guild_heros 请求太频繁", ERR_SEARCH_NO_GUILD_HEROS_FAIL_TOO_FAST)          // 192-3
)

// view_mc_war_record
var (
	ErrViewMcWarRecordFailNoGuild = newMsgError("view_mc_war_record 没有联盟", ERR_VIEW_MC_WAR_RECORD_FAIL_NO_GUILD) // 201-1
)

// update_guild_mark
var (
	ErrUpdateGuildMarkFailInvalidIndex   = newMsgError("update_guild_mark 无效的序号", ERR_UPDATE_GUILD_MARK_FAIL_INVALID_INDEX)       // 198-1
	ErrUpdateGuildMarkFailInvalidPos     = newMsgError("update_guild_mark 无效的坐标", ERR_UPDATE_GUILD_MARK_FAIL_INVALID_POS)         // 198-2
	ErrUpdateGuildMarkFailInvalidMsg     = newMsgError("update_guild_mark 无效的标记内容", ERR_UPDATE_GUILD_MARK_FAIL_INVALID_MSG)       // 198-3
	ErrUpdateGuildMarkFailSensitiveWords = newMsgError("update_guild_mark 标记内容包含敏感词", ERR_UPDATE_GUILD_MARK_FAIL_SENSITIVE_WORDS) // 198-4
	ErrUpdateGuildMarkFailNotInGuild     = newMsgError("update_guild_mark 你没有联盟", ERR_UPDATE_GUILD_MARK_FAIL_NOT_IN_GUILD)        // 198-5
	ErrUpdateGuildMarkFailDeny           = newMsgError("update_guild_mark 你没有权限操作", ERR_UPDATE_GUILD_MARK_FAIL_DENY)              // 198-6
	ErrUpdateGuildMarkFailServerError    = newMsgError("update_guild_mark 服务器忙，请稍后再试", ERR_UPDATE_GUILD_MARK_FAIL_SERVER_ERROR)   // 198-7
)

// view_yinliang_record
var (
	ErrViewYinliangRecordFailNoGuild = newMsgError("view_yinliang_record 没有联盟", ERR_VIEW_YINLIANG_RECORD_FAIL_NO_GUILD) // 204-1
)

// send_yinliang_to_other_guild
var (
	ErrSendYinliangToOtherGuildFailDeny          = newMsgError("send_yinliang_to_other_guild 没有权限", ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_DENY)           // 207-1
	ErrSendYinliangToOtherGuildFailNoGuild       = newMsgError("send_yinliang_to_other_guild 对方联盟不存在", ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_NO_GUILD)    // 207-2
	ErrSendYinliangToOtherGuildFailNotEnough     = newMsgError("send_yinliang_to_other_guild 钱不够", ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_NOT_ENOUGH)      // 207-3
	ErrSendYinliangToOtherGuildFailInvalidAmount = newMsgError("send_yinliang_to_other_guild 参数非法", ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_INVALID_AMOUNT) // 207-4
	ErrSendYinliangToOtherGuildFailServerErr     = newMsgError("send_yinliang_to_other_guild 服务器错误", ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_SERVER_ERR)    // 207-5
)

// send_yinliang_to_member
var (
	ErrSendYinliangToMemberFailDeny          = newMsgError("send_yinliang_to_member 没有权限", ERR_SEND_YINLIANG_TO_MEMBER_FAIL_DENY)           // 210-1
	ErrSendYinliangToMemberFailNoMember      = newMsgError("send_yinliang_to_member 成员不存在", ERR_SEND_YINLIANG_TO_MEMBER_FAIL_NO_MEMBER)     // 210-2
	ErrSendYinliangToMemberFailNotEnough     = newMsgError("send_yinliang_to_member 钱不够", ERR_SEND_YINLIANG_TO_MEMBER_FAIL_NOT_ENOUGH)      // 210-3
	ErrSendYinliangToMemberFailInvalidAmount = newMsgError("send_yinliang_to_member 参数非法", ERR_SEND_YINLIANG_TO_MEMBER_FAIL_INVALID_AMOUNT) // 210-4
	ErrSendYinliangToMemberFailServerErr     = newMsgError("send_yinliang_to_member 服务器错误", ERR_SEND_YINLIANG_TO_MEMBER_FAIL_SERVER_ERR)    // 210-5
)

// pay_salary
var (
	ErrPaySalaryFailDeny          = newMsgError("pay_salary 没有权限", ERR_PAY_SALARY_FAIL_DENY)           // 213-1
	ErrPaySalaryFailNotEnough     = newMsgError("pay_salary 钱不够", ERR_PAY_SALARY_FAIL_NOT_ENOUGH)      // 213-2
	ErrPaySalaryFailInvalidAmount = newMsgError("pay_salary 参数非法", ERR_PAY_SALARY_FAIL_INVALID_AMOUNT) // 213-3
	ErrPaySalaryFailServerErr     = newMsgError("pay_salary 服务器错误", ERR_PAY_SALARY_FAIL_SERVER_ERR)    // 213-4
)

// set_salary
var (
	ErrSetSalaryFailDeny          = newMsgError("set_salary 没有权限", ERR_SET_SALARY_FAIL_DENY)           // 216-1
	ErrSetSalaryFailNoMember      = newMsgError("set_salary 成员不存在", ERR_SET_SALARY_FAIL_NO_MEMBER)     // 216-2
	ErrSetSalaryFailInvalidSalary = newMsgError("set_salary 工资非法", ERR_SET_SALARY_FAIL_INVALID_SALARY) // 216-3
	ErrSetSalaryFailServerErr     = newMsgError("set_salary 服务器错误", ERR_SET_SALARY_FAIL_SERVER_ERR)    // 216-4
)

// view_send_yinliang_to_guild
var (
	ErrViewSendYinliangToGuildFailServerErr = newMsgError("view_send_yinliang_to_guild 服务器错误", ERR_VIEW_SEND_YINLIANG_TO_GUILD_FAIL_SERVER_ERR) // 220-1
)

// convene
var (
	ErrConveneFailInvalidTarget    = newMsgError("convene 无效的目标id", ERR_CONVENE_FAIL_INVALID_TARGET)        // 230-1
	ErrConveneFailNotInGuild       = newMsgError("convene 你没有联盟", ERR_CONVENE_FAIL_NOT_IN_GUILD)            // 230-2
	ErrConveneFailDeny             = newMsgError("convene 权限不足", ERR_CONVENE_FAIL_DENY)                     // 230-3
	ErrConveneFailTargetNotInGuild = newMsgError("convene 目标不在你的联盟中", ERR_CONVENE_FAIL_TARGET_NOT_IN_GUILD) // 230-4
	ErrConveneFailCooldown         = newMsgError("convene 盟友召集CD中，请稍后再试", ERR_CONVENE_FAIL_COOLDOWN)        // 230-5
	ErrConveneFailServerError      = newMsgError("convene 服务器忙，请稍后再试", ERR_CONVENE_FAIL_SERVER_ERROR)       // 230-6
)

// collect_daily_guild_rank_prize
var (
	ErrCollectDailyGuildRankPrizeFailNoGuild     = newMsgError("collect_daily_guild_rank_prize 没有联盟", ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD)           // 233-1
	ErrCollectDailyGuildRankPrizeFailNoGuildRank = newMsgError("collect_daily_guild_rank_prize 联盟未上榜", ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD_RANK)     // 233-2
	ErrCollectDailyGuildRankPrizeFailCollected   = newMsgError("collect_daily_guild_rank_prize 无法领取", ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_COLLECTED)          // 233-3
	ErrCollectDailyGuildRankPrizeFailServerError = newMsgError("collect_daily_guild_rank_prize 服务器忙，请稍后再试", ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_SERVER_ERROR) // 233-4
)

// view_daily_guild_rank
var (
	ErrViewDailyGuildRankFailNoCountry   = newMsgError("view_daily_guild_rank 还未加入国家", ERR_VIEW_DAILY_GUILD_RANK_FAIL_NO_COUNTRY)       // 236-1
	ErrViewDailyGuildRankFailServerError = newMsgError("view_daily_guild_rank 服务器忙，请稍后再试", ERR_VIEW_DAILY_GUILD_RANK_FAIL_SERVER_ERROR) // 236-2
)

// add_recommend_mc_build
var (
	ErrAddRecommendMcBuildFailDeny            = newMsgError("add_recommend_mc_build 没有权限", ERR_ADD_RECOMMEND_MC_BUILD_FAIL_DENY)                 // 242-1
	ErrAddRecommendMcBuildFailInvalidMcId     = newMsgError("add_recommend_mc_build 名城不存在", ERR_ADD_RECOMMEND_MC_BUILD_FAIL_INVALID_MC_ID)       // 242-2
	ErrAddRecommendMcBuildFailMcIsRecommended = newMsgError("add_recommend_mc_build 名城已经被推荐", ERR_ADD_RECOMMEND_MC_BUILD_FAIL_MC_IS_RECOMMENDED) // 242-3
	ErrAddRecommendMcBuildFailNoGuild         = newMsgError("add_recommend_mc_build 没有联盟", ERR_ADD_RECOMMEND_MC_BUILD_FAIL_NO_GUILD)             // 242-5
	ErrAddRecommendMcBuildFailServerErr       = newMsgError("add_recommend_mc_build 服务器错误", ERR_ADD_RECOMMEND_MC_BUILD_FAIL_SERVER_ERR)          // 242-4
)

// view_task_progress
var (
	ErrViewTaskProgressFailNoGuild         = newMsgError("view_task_progress 没有联盟", ERR_VIEW_TASK_PROGRESS_FAIL_NO_GUILD)            // 245-1
	ErrViewTaskProgressFailGuildLevelLimit = newMsgError("view_task_progress 联盟等级不足", ERR_VIEW_TASK_PROGRESS_FAIL_GUILD_LEVEL_LIMIT) // 245-2
)

// collect_task_prize
var (
	ErrCollectTaskPrizeFailInvalidValue = newMsgError("collect_task_prize 无效的数据", ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_VALUE) // 249-1
	ErrCollectTaskPrizeFailNoGuild      = newMsgError("collect_task_prize 没有联盟", ERR_COLLECT_TASK_PRIZE_FAIL_NO_GUILD)       // 249-2
	ErrCollectTaskPrizeFailNoPrize      = newMsgError("collect_task_prize 阶段奖励未激活", ERR_COLLECT_TASK_PRIZE_FAIL_NO_PRIZE)    // 249-3
	ErrCollectTaskPrizeFailCollected    = newMsgError("collect_task_prize 已领取", ERR_COLLECT_TASK_PRIZE_FAIL_COLLECTED)       // 249-4
)

// guild_change_country
var (
	ErrGuildChangeCountryFailInvalidCountry = newMsgError("guild_change_country 无效的国家", ERR_GUILD_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY)    // 252-1
	ErrGuildChangeCountryFailNotInGuild     = newMsgError("guild_change_country 你没有联盟", ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)       // 252-2
	ErrGuildChangeCountryFailNotLearder     = newMsgError("guild_change_country 你不是盟主，不能联盟转国", ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER) // 252-3
	ErrGuildChangeCountryFailCostNotEnough  = newMsgError("guild_change_country 消耗不足", ERR_GUILD_CHANGE_COUNTRY_FAIL_COST_NOT_ENOUGH)     // 252-4
	ErrGuildChangeCountryFailCooldown       = newMsgError("guild_change_country 转国cd中", ERR_GUILD_CHANGE_COUNTRY_FAIL_COOLDOWN)           // 252-5
	ErrGuildChangeCountryFailExist          = newMsgError("guild_change_country 当前正在转国中", ERR_GUILD_CHANGE_COUNTRY_FAIL_EXIST)            // 252-6
	ErrGuildChangeCountryFailSameCountry    = newMsgError("guild_change_country 你的联盟已经是这个国家", ERR_GUILD_CHANGE_COUNTRY_FAIL_SAME_COUNTRY) // 252-7
	ErrGuildChangeCountryFailIsKing         = newMsgError("guild_change_country 国王不能转国", ERR_GUILD_CHANGE_COUNTRY_FAIL_IS_KING)           // 252-8
)

// cancel_guild_change_country
var (
	ErrCancelGuildChangeCountryFailNotInGuild = newMsgError("cancel_guild_change_country 你没有联盟", ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)         // 255-1
	ErrCancelGuildChangeCountryFailNotLearder = newMsgError("cancel_guild_change_country 你不是盟主，不能取消联盟转国", ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER) // 255-2
	ErrCancelGuildChangeCountryFailCooldown   = newMsgError("cancel_guild_change_country 转国cd中", ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_COOLDOWN)             // 255-3
	ErrCancelGuildChangeCountryFailNotExist   = newMsgError("cancel_guild_change_country 联盟没有转国中", ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_EXIST)          // 255-4
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
