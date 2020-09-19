package realmerr

import "github.com/pkg/errors"

var (
	ErrLockHeroErr = errors.New("lock英雄失败")
)

// AddBase
var (
	ErrAddBaseHomeNotAlive         = errors.New("AddBase时, 主城是流亡状态")
	ErrAddBaseAlreadyHasRealm      = errors.New("AddBase时, 已经有归属的realm")
	ErrAddBaseHomeAlive            = errors.New("AddBase时, 已经有城了")
	ErrAddBaseHomeAndTentSameRealm = errors.New("AddBase时, 主城和行营同时在一张地图")
)

// Invasion
var (
	ErrInvasionInvalidTarget     = errors.New("Invasion时, target非法")
	ErrInvasionTargetNotExist    = errors.New("Invasion时, target不存在")
	ErrInvasionSelfNoBase        = errors.New("Invasion时, 自己已经没有基地了. 可能刚流亡?")
	ErrInvasionEmptyGeneral      = errors.New("Invasion时, 没有武将")
	ErrInvasionGeneralOutside    = errors.New("Invasion时, 武将已经出征了")
	ErrInvasionNoSoldier         = errors.New("Invasion时, 没有士兵")
	ErrInvasionInvalidRelation   = errors.New("Invasion时, 关系和行动不一致")
	ErrInvasionInvalidTroopIndex = errors.New("Invasion时, 队伍序号无效")
	ErrInvasionMian              = errors.New("Invasion时, target免战")
	ErrInvasionTodayJoinXiongNu  = errors.New("Invasion时, 今日已经参与过反击匈奴了")
)

// upgrade base
var (
	ErrUpgradeBaseAlreadyMax          = errors.New("UpgradeBase时, 已经满级")
	ErrUpgradeBaseNotEnoughProsperity = errors.New("UpgradeBase时, 繁荣度不够")
	ErrUpgradeBaseHomeNotAlive        = errors.New("UpgradeBase时, 流亡状态")
	ErrUpgradeBaseNotMyRealm          = errors.New("UpgradeBase时, 不归我这个realm管")
)

// add home npc
var (
	ErrAddHomeNpcNotHome = errors.New("AddHomeNpc时, 不是主城")
)

// cancel invasion
var (
	ErrCancelInvasionTroopNotFound        = errors.New("CancelInvasion时, 没找到队伍id")
	ErrCancelInvasionTroopAlreadyBacking  = errors.New("CancelInvasion时, 队伍正在回家")
	ErrCancelInvasionTroopAlreadyHome     = errors.New("CancelInvasion时, 队伍本来就在家里防守")
	ErrCancelInvasionTroopAssemblyStarted = errors.New("CancelInvasion时, 集结队伍已经出发，不能召回")
)

// Repatriate
var (
	ErrRepatriateTroopNotFound      = errors.New("Repatriate时, 没找到队伍id")
	ErrRepatriateTroopNoDefending   = errors.New("Repatriate时, 只能遣返驻守部队")
	ErrRepatriateAssemblyStarted    = errors.New("Repatriate时, 集结已出发")
	ErrRepatriateNotAssemblyCreater = errors.New("Repatriate时, 不是集结创建者")
)

// BaozRepatriate
var (
	ErrBaozRepatriateTroopNotFound = errors.New("BaozRepatriate时, 没找到队伍id")
	ErrBaozRepatriateBaozNotKeep   = errors.New("BaozRepatriate时, 不是你控制的宝藏")
)

// SpeedUp
var (
	ErrSpeedUpTroopNotFound      = errors.New("SpeedUp时, 没找到队伍id")
	ErrSpeedUpTroopNoMoving      = errors.New("SpeedUp时, 只能加速行军中的部队")
	ErrSpeedUpOtherTroopNotFound = errors.New("SpeedUp时, 没找到目标队伍id")
	ErrSpeedUpAssemblyWait       = errors.New("SpeedUp时, 集结等待中")
)

// SlowMoveBase
var (
	ErrSlowMoveBaseSelfNoBase = errors.New("SlowMoveBase时, 自己没有基地，可能刚流亡")
	ErrSlowMoveBaseTent       = errors.New("SlowMoveBase时, 行营不能缓慢移动")
	ErrSlowMoveBaseInvalidPos = errors.New("SlowMoveBase时, 无效的目标位置")
)

// FastMoveBase
var (
	ErrFastMoveBaseSelfNoBase = errors.New("FastMoveBase时, 自己没有基地，可能刚流亡")
	ErrFastMoveBasePosChanged = errors.New("FastMoveBase时, 自己基地位置变更了")
	ErrFastMoveBaseOutside    = errors.New("FastMoveBase时, 部队出征中")
)

// cancel SlowMoveBase
var (
	ErrCancelSlowMoveBaseSelfNoBase = errors.New("CancelSlowMoveBase时, 自己没有基地，可能刚流亡")
	ErrCancelSlowMoveBaseTent       = errors.New("CancelSlowMoveBase时, 行营不能缓慢移动")
)

// MoveBase
var (
	ErrMoveBaseInvalidPos   = errors.New("MoveBase时, 无效的目标位置")
	ErrMoveBaseHomeConflict = errors.New("MoveBase时,目标位置已经被其他人占了")
)

// Expel
var (
	ErrExpelSelfNoBase        = errors.New("Expel时，自己没有基地，可能刚流亡")
	ErrExpelTroopsNotFound    = errors.New("Expel时，目标部队不存在")
	ErrExpelTroopsNoRobbing   = errors.New("Expel时，目标部队不是持续掠夺中")
	ErrExpelFightError        = errors.New("Expel时，计算战斗结果失败")
	ErrExpelCaptainOutside    = errors.New("Expel时，武将已出征")
	ErrExpelNoSoldier         = errors.New("Expel时，武将没有带士兵")
	ErrExpelInvalidTroopIndex = errors.New("Expel时，队伍序号无效")
)

// Remove
var (
	ErrRemoveBaseSelfNoBase     = errors.New("Remove时，自己没有基地，可能刚流亡")
	ErrRemoveBaseOutside        = errors.New("Remove时，部队出征中")
	ErrRemoveBaseUnkownBaseType = errors.New("Remove时，未知的基地类型")
)

// Remove tent validtime
var (
	ErrRemoveValidTimeSelfNoBase  = errors.New("RemoveValidTime时，自己没有基地，可能刚流亡")
	ErrRemoveValidTimeNoTent      = errors.New("RemoveValidTime时，没有行营")
	ErrRemoveValidTimeNoValidTime = errors.New("RemoveValidTime时，当前没有在建造")
)

// AddProsperity
var (
	ErrAddProsperitySelfNoBase      = errors.New("AddProsperity时，自己没有基地，可能刚流亡")
	ErrAddProsperityInvalidBaseType = errors.New("AddProsperity时，基地类型错误")
)

// ChangeGuild
var (
	ErrChangeGuildSelfNoBase = errors.New("ChangeGuild时，自己没有基地，可能刚流亡")
)

// Mian
var (
	ErrMianSelfNoBase    = errors.New("Mian时，自己没有基地，可能刚流亡")
	ErrMianTent          = errors.New("Mian时，行营不能免战")
	ErrMianExist         = errors.New("Mian时，当前已经存在免战状态")
	ErrMianCantOverwrite = errors.New("Mian时，当前免战时间更大，无法覆盖")
)

// Query
var (
	ErrGetMilitaryBaseNotFound  = errors.New("GetMilitary时，base不存在")
	ErrGetMilitaryTroopNotFound = errors.New("GetMilitary时，troop不存在")
)

// CreateAssembly
var (
	ErrCreateAssemblyInvalidInput      = errors.New("CreateAssembly时, input非法")
	ErrCreateAssemblyInvalidTarget     = errors.New("CreateAssembly时, target非法")
	ErrCreateAssemblyTargetNotExist    = errors.New("CreateAssembly时, target不存在")
	ErrCreateAssemblySelfNoBase        = errors.New("CreateAssembly时, 自己已经没有基地了. 可能刚流亡?")
	ErrCreateAssemblySelfNoGuild       = errors.New("CreateAssembly时, 自己没有联盟")
	ErrCreateAssemblyEmptyGeneral      = errors.New("CreateAssembly时, 没有武将")
	ErrCreateAssemblyGeneralOutside    = errors.New("CreateAssembly时, 武将已经出征了")
	ErrCreateAssemblyNoSoldier         = errors.New("CreateAssembly时, 没有士兵")
	ErrCreateAssemblyInvalidRelation   = errors.New("CreateAssembly时, 关系和行动不一致")
	ErrCreateAssemblyMian              = errors.New("CreateAssembly时, 目标免战")
	ErrCreateAssemblyInvalidTroopIndex = errors.New("CreateAssembly时, 队伍序号无效")
	ErrCreateAssemblyTodayJoinXiongNu  = errors.New("CreateAssembly时, 今日已经参与过反击匈奴了")
)

// JoinAssembly
var (
	ErrJoinAssemblyInvalidTarget     = errors.New("JoinAssembly时, target非法")
	ErrJoinAssemblyTargetNotExist    = errors.New("JoinAssembly时, target不存在")
	ErrJoinAssemblySelfNoBase        = errors.New("JoinAssembly时, 自己已经没有基地了. 可能刚流亡?")
	ErrJoinAssemblyEmptyGeneral      = errors.New("JoinAssembly时, 没有武将")
	ErrJoinAssemblyGeneralOutside    = errors.New("JoinAssembly时, 武将已经出征了")
	ErrJoinAssemblyNoSoldier         = errors.New("JoinAssembly时, 没有士兵")
	ErrJoinAssemblyInvalidRelation   = errors.New("JoinAssembly时, 关系和行动不一致")
	ErrJoinAssemblyInvalidTroopIndex = errors.New("JoinAssembly时, 队伍序号无效")
	ErrJoinAssemblyTodayJoinXiongNu  = errors.New("JoinAssembly时, 今日已经参与过反击匈奴了")
	ErrJoinAssemblyFull              = errors.New("JoinAssembly时, 集结已满")
	ErrJoinAssemblyMultiJoin         = errors.New("JoinAssembly时, 不能多个队伍加入同一个集结")
	ErrJoinAssemblyStarted           = errors.New("JoinAssembly时, 集结已经出发")
)
