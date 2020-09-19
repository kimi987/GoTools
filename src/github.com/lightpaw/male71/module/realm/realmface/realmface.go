package realmface

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"math"
	"time"
)

// destroy
type BaseDestroyType int32

const (
	BDTBroken BaseDestroyType = 0 // 被打爆
)

type RandomPointType uint8

const (
	RPTRandom RandomPointType = iota // 随机
	//RPTNewHero                        // 新建英雄
	RPTReborn  // 流亡重建
)

// --- 主城/行营 ---

type BaseType = int32

const (
	BaseTypeHome BaseType = 0 // 玩家主城
	BaseTypeNpc  BaseType = 2 // NPC
)

type Base interface {
	Id() int64       // 英雄id
	IdBytes() []byte // 英雄idBytes
	RegionID() int64 // 地区id
	GuildId() int64
	HeroName() string

	GetBaseLevel() uint64
	BaseX() int
	BaseY() int
	Prosperity() uint64
	ProsperityCapcity() uint64
	BaseType() BaseType

	TroopType() shared_proto.TroopType

	// 战斗力?
}

type Home interface {
	Base
}

// --- 部队 ---

//type ActionType = shared_proto.TroopOperate
//
//func (t ActionType) Int32() int32 {
//	return int32(t)
//}

type MoveType int32

func (t MoveType) Int32() int32 {
	return int32(t)
}

const (
	//ActionInvation ActionType = 0
	//ActionAssist   ActionType = 1
	//ActionAssembly   ActionType = 2

	MoveForward MoveType = 0
	MoveStay    MoveType = 1
	MoveBack    MoveType = 2
	MoveWait    MoveType = 3
)

func NewTroopState(o shared_proto.TroopOperate, moveType MoveType) TroopState {
	// 最低的4个byte预留，兼容原来的值
	return TroopState(uint8(o)<<3 | uint8(moveType))
}

type TroopState uint8

func (s TroopState) Operate() shared_proto.TroopOperate {
	return shared_proto.TroopOperate(uint8(s) >> 3)
}

func (s TroopState) MoveType() MoveType {
	return MoveType(uint8(s) & 7)
}

const (
	Temp             = 0 // 临时部队，驱逐，驻守时使用
	MovingToInvade   = TroopState(uint8(shared_proto.TroopOperate_ToInvasion)<<3 | uint8(MoveForward))
	InvadeMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToInvasion)<<3 | uint8(MoveBack))
	Robbing          = TroopState(uint8(shared_proto.TroopOperate_ToInvasion)<<3 | uint8(MoveStay))
	Assembly         = TroopState(uint8(shared_proto.TroopOperate_ToInvasion)<<3 | uint8(MoveWait))

	MovingToAssist   = TroopState(uint8(shared_proto.TroopOperate_ToAssist)<<3 | uint8(MoveForward))
	AssistMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToAssist)<<3 | uint8(MoveBack))
	Defending        = TroopState(uint8(shared_proto.TroopOperate_ToAssist)<<3 | uint8(MoveStay))

	MovingToAssembly   = TroopState(uint8(shared_proto.TroopOperate_ToAssembly)<<3 | uint8(MoveForward))
	AssemblyMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToAssembly)<<3 | uint8(MoveBack))
	AssemblyArrived    = TroopState(uint8(shared_proto.TroopOperate_ToAssembly)<<3 | uint8(MoveStay))

	MovingToWorkshopBuild   = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopBuild)<<3 | uint8(MoveForward))
	WorkshopBuildMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopBuild)<<3 | uint8(MoveBack))

	MovingToWorkshopProd   = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopProd)<<3 | uint8(MoveForward))
	WorkshopProdMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopProd)<<3 | uint8(MoveBack))

	MovingToWorkshopPrize   = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopPrize)<<3 | uint8(MoveForward))
	WorkshopPrizeMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToWorkshopPrize)<<3 | uint8(MoveBack))

	MovingToInvesigate   = TroopState(uint8(shared_proto.TroopOperate_ToInvestigate)<<3 | uint8(MoveForward))
	InvesigateMovingBack = TroopState(uint8(shared_proto.TroopOperate_ToInvestigate)<<3 | uint8(MoveBack))

	//MovingToAssist     TroopState = 2
	//InvadeMovingBack   TroopState = 3
	//AssistMovingBack   TroopState = 6 // 协助者返回
	//Defending          TroopState = 4
	//Robbing            TroopState = 5
	//Assembly           TroopState = 7  // 集结等待
	//MovingToAssembly   TroopState = 8  // 前往集结
	//AssemblyMovingBack TroopState = 9  // 集结返回
	//AssemblyArrived    TroopState = 10 // 到达集结地
)

func (s TroopState) IsMoving() bool {
	switch s {
	case MovingToInvade, InvadeMovingBack,
		MovingToAssist, AssistMovingBack,
		MovingToAssembly, AssemblyMovingBack,
		MovingToWorkshopBuild, WorkshopBuildMovingBack,
		MovingToWorkshopProd, WorkshopProdMovingBack,
		MovingToWorkshopPrize, WorkshopPrizeMovingBack,
		MovingToInvesigate, InvesigateMovingBack:
		return true
	}

	return false
}

func (s TroopState) IsInvateState() bool {
	switch s {
	case MovingToInvade, Robbing, Assembly, MovingToInvesigate:
		return true
	}

	return false
}

func (s TroopState) IsAssistState() bool {
	switch s {
	case MovingToAssist, Defending:
		return true
	}

	return false
}

type Troop interface {
	Id() int64          // 部队的id
	StartingBase() Base // 所属的基地
	TargetBase() Base   // 目标基地
	State() TroopState  // 当前的状态

	BackHomeTargetX() int
	BackHomeTargetY() int

	TargetIsOwnerCanSee() bool

	CreateTime() time.Time
	MoveStartTime() time.Time            // 开始移动时间, 如果是Defending/Robbing状态则为空值
	MoveArriveTime() time.Time           // 预计到达时间, 如果是Defending/Robbing状态则为空值
	RobbingEndTime() time.Time           // 持续掠夺结束时间，如果不是Robbing状态，则为空值
	NextReduceProsperityTime() time.Time // 下次扣繁荣度时间，如果不是Robbing状态，则为空值
	NextAddHateTime() time.Time          // 下次添加仇恨的时间
	NextRobBaowuTime() time.Time


	//kimi 移除地图上的点
	RemoveRealmPoints()

	//Carrying() (gold uint64, food uint64, wood uint64, stone uint64) // 身上背着的资源
	//ClearCarrying() (gold uint64, food uint64, wood uint64, stone uint64) // 清掉背着的资源, 并返回
	//JadeOre() uint64 // 掠夺的玉石矿

	AccumRobPrize() *shared_proto.PrizeProto
	AccumReduceProsperity() uint64

	// 武将, 士兵
	Captains() []Captain

	AssemblyId() int64       // 发起集结的id，这个是自己的话，说明这个集结是我发起的，0表示非集结队伍
	AssemblyTargetId() int64 // 集结的目标

	Dialogue() uint64

	NpcTimes() uint64
}

// --- 武将 ---

type Captain interface {
	Id() uint64 // 仅对每个英雄唯一
	Index() int // 在部队中的index 1-5
	Proto() *shared_proto.CaptainInfoProto
	// 设置最新的士兵数, 并返回旧的士兵数.
	GetAndSetSoldierCount(uint64) uint64
}

//

// --- hero内保存的部队信息 ---
//type TroopInHero interface {
//	Id() int64
//	RegionID() int64
//	TargetBaseID() int64
//	TargetBaseLevel() uint64
//	State() TroopState
//
//	OriginTargetID() int64
//	BackHomeTargetX() int
//	BackHomeTargetY() int
//
//	TargetIsOwnerCanSee() bool
//
//	CreateTime() time.Time
//	MoveStartTime() time.Time
//	MoveArriveTime() time.Time
//	RobbingEndTime() time.Time
//	NextReduceProsperityTime() time.Time
//	NextAddHateTime() time.Time
//	NextRobBaowuTime() time.Time
//
//	//Carrying() (gold uint64, food uint64, wood uint64, stone uint64) // 身上背着的资源
//	//JadeOre() uint64                                                 // 掠夺的玉石矿
//	AccumRobPrize() *shared_proto.PrizeProto // 累积抢到的资源
//	AccumReduceProsperity() uint64
//	Captains() []uint64     // 第一个是captain id, 第二个是captain的index
//	CaptainXIndex() []int32 // captain对应的xindex
//
//	AssemblyId() int64 // 发起集结的id，这个是自己的话，说明这个集结是我发起的，0表示非集结队伍
//}

// --- add base type
type AddBaseType uint8

const (
	AddBaseHomeNewHero  AddBaseType = iota // 新建角色, 调用时 baseLevel和繁荣度是0
	AddBaseHomeReborn                      // 流亡后重生, 调用时 baseLevel和繁荣度为0
	AddBaseHomeTransfer                    // 老家迁城, 调用时必须有baseLevel和繁荣度
)

func GetRealmId(level, sequence uint64, regionType shared_proto.RegionType) int64 {
	return int64((sequence << 12) |
		(level << 4) | // 8bit
		uint64(regionType)) // 4bit
}

func ParseRealmId(id int64) (level, sequence uint64, regionType shared_proto.RegionType) {

	regionType = ParseRegionType(id) // 4bit
	level = ParseRegionLevel(id)     // 8bit
	sequence = ParseRegionSequence(id)
	return
}

func ParseRegionType(id int64) shared_proto.RegionType {
	return shared_proto.RegionType(id & 0xf) // 4bit
}

func ParseRegionLevel(id int64) uint64 {
	return (uint64(id) >> 4) & math.MaxUint8 // 8bit
}

func ParseRegionSequence(id int64) (sequence uint64) {
	sequence = uint64(id) >> 12
	return
}

func GetGuildRealmId(guildId int64) int64 {
	return GetRealmId(1, uint64(guildId), shared_proto.RegionType_GUILD)
}

type ConflictResourcePoint interface {
	LayoutId() uint64             // 布局id
	OutputStartTime() time.Time   // 资源产出开始时间
	ConflictStartTime() time.Time // 资源冲突开始时间
}

//func GetHeroTroopId(heroId int64, index uint64) int64 {
//	return int64(uint64(heroId)<<4 | index)
//}
//
//func GetNpcTroopId(baseId int64, sequence uint64) int64 {
//
//}
//
//func GetTroopHeroId(troopsId int64) int64 {
//	return int64(uint64(troopsId) >> 4)
//}
//
//const indexMask = 1<<4 - 1
//
//func GetTroopIndex(troopsId int64) uint64 {
//	return uint64(troopsId) & indexMask
//}

type ViewArea struct {
	// 中心点坐标
	CenterX, CenterY int

	// 最小坐标
	MinX, MinY int

	// 最大坐标
	MaxX, MaxY int

	UpdateTime time.Time
}

//GetCenterDistant 获取中心点的距离的平方 这里不求根号的原因是 求根运算消耗比较高
func (va *ViewArea) GetCenterDistant(vaTo *ViewArea) int {
	x, y := vaTo.CenterX-va.CenterX, vaTo.CenterY-va.CenterY
	return x*x + y*y
}

func (va *ViewArea) CanSeePos(x, y int) bool {
	return x >= va.MinX && x <= va.MaxX && y >= va.MinY && y <= va.MaxY
}

func (va *ViewArea) CanSeeLine(x1, y1, x2, y2 int) bool {

	// 先简单判断，看得到点的，都能看到线
	if va.CanSeePos(x1, y1) || va.CanSeePos(x2, y2) {
		return true
	}

	// 快速排除，线段的包围盒与矩形是否相交（不相交，则线段与矩形也肯定不相交）
	minX := imath.Min(x1, x2)
	minY := imath.Min(y1, y2)
	maxX := imath.Max(x1, x2)
	maxY := imath.Max(y1, y2)
	if !isRectIntersect(minX, minY, maxX, maxY, va.MinX, va.MinY, va.MaxX, va.MaxY) {
		return false
	}

	// 线与矩形的4条边是否相交，是则相交，否则不相交
	if isSegmentIntersect(x1, y1, x2, y2, va.MinX, va.MinY, va.MinX, va.MaxY) {
		return true
	}
	if isSegmentIntersect(x1, y1, x2, y2, va.MinX, va.MinY, va.MaxX, va.MinY) {
		return true
	}
	if isSegmentIntersect(x1, y1, x2, y2, va.MinX, va.MaxY, va.MaxX, va.MaxY) {
		return true
	}
	if isSegmentIntersect(x1, y1, x2, y2, va.MaxX, va.MinY, va.MaxX, va.MaxY) {
		return true
	}

	return false
}

// 矩形相交
func isRectIntersect(minx1, miny1, maxx1, maxy1, minx2, miny2, maxx2, maxy2 int) bool {
	//假定矩形是用一对点表达的(minx, miny) (maxx, maxy)，那么两个矩形
	//rect1{(minx1, miny1)(maxx1, maxy1)}
	//rect2{(minx2, miny2)(maxx2, maxy2)}
	//相交的结果一定是个矩形，构成这个相交矩形rect{(minx, miny) (maxx, maxy)}的点对坐标是：
	minx := imath.Max(minx1, minx2)
	miny := imath.Max(miny1, miny2)
	maxx := imath.Min(maxx1, maxx2)
	maxy := imath.Min(maxy1, maxy2)

	//如果两个矩形不相交，那么计算得到的点对坐标必然满足：
	//（ minx  >  maxx ） 或者 （ miny  >  maxy ）
	return minx <= maxx && miny <= maxy
}

// 线段相交
func isSegmentIntersect(ax, ay, bx, by, cx, cy, dx, dy int) bool {
	// 线段的两个端点分别在另一条线段的两侧
	return cross(ax, ay, bx, by, cx, cy)*cross(ax, ay, bx, by, dx, dy) <= 0 &&
		cross(cx, cy, dx, dy, ax, ay)*cross(cx, cy, dx, dy, bx, by) <= 0
}

// 求出向量 AB 与向量 AC 的向量积,返回 0 代表共线
// 叉积的概念： 设向量 a(x1, y1) 、 b(x2, y2)
// a x b = x1*y2 - x2*y1
// <0 左侧， >0 右侧， =0 同一直线
func cross(ax, ay, bx, by, cx, cy int) int {
	x1 := bx - ax
	y1 := by - ay
	x2 := cx - ax
	y2 := cy - ay

	return x1*y2 - x2*y1
}
