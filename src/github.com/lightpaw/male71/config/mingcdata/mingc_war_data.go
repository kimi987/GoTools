package mingcdata

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"time"
)

const McWarLoopDuration = time.Duration(200 * time.Millisecond)

//gogen:config
type MingcMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"名城战/杂项.txt"`
	_ struct{} `protogen:"true"`

	FightPrepareDuration                   time.Duration    `desc:"战斗开始后的准备时间" default:"30m"`
	JoinFightDuration                      time.Duration    `desc:"入场时间" default:"180s"`
	ApplyAstLimit                          uint64           `desc:"联盟申请协助最大数量" default:"3"`
	StartAfterServerOpen                   time.Duration    `desc:"开服多久后开始名城战"`
	StartSelfCapitalAfterServerOpen        time.Duration    `desc:"开服多久后开始名城战本国都城"`
	StartOtherCapitalAfterServerOpen       time.Duration    `desc:"开服多久后开始名城战他国都城"`
	StartSelfCapitalNoticeAfterServerOpen  time.Duration    `desc:"开服多久后开始名城战本国都城大事记"`
	StartOtherCapitalNoticeAfterServerOpen time.Duration    `desc:"开服多久后开始名城战他国都城大事记"`
	DestroyProsperityMaxTroop              uint64           `desc:"攻打据点摧毁繁荣度最多叠加队伍" default:"5"`
	PerDestroyProsperity                   uint64           `desc:"每支队伍每x秒摧毁繁荣度"`
	DestroyProsperityDuration              time.Duration    `desc:"摧毁繁荣度间隔"`
	ReliveDuration                         time.Duration    `desc:"补兵持续时间"`
	WallStat                               *data.SpriteStat `default:"nullable" protofield:"-"`           // 城墙属性
	WallFixDamage                          uint64           `default:"0" validator:"uint" protofield:"-"` // 城墙固定伤害
	WallLevel                              uint64           `default:"1"`                                 // 城墙等级
	SceneRecordMaxLen                      uint64           `default:"50" protofield:"-"`
	Speed                                  float64          `desc:"速度/秒" default:"1"`
	CloseDuCheng                           bool             `desc:"是否关闭都城" default:"true"`
	JoinFightHeroMinLevel                  uint64           `desc:"参战最小君主等级" default:"10"`
	RedPointMinGuildLevel                  uint64           `desc:"红点推送最小联盟等级" default:"3" protofield:"-"`
	SaveHeroRecordMaxDays                  uint64           `desc:"保存几天的玩家战斗记录" default:"7" protofield:"-"`
	DailyUpdateMingcTime                   time.Duration    `desc:"每日重置名城的时间" default:"22h"`
	FreeTankSpeed                          float64          `desc:"免费攻城车速度/秒" default:"0.7"`
	FreeTankPerDestroyProsperity           uint64           `desc:"免费攻城车每x秒摧毁繁荣度" default:"300"`
	TouShiBuildingTurnDuration             time.Duration    `desc:"投石机转向时长" default:"5s"`
	TouShiBuildingPrepareDuration          time.Duration    `desc:"投石机装填时长" default:"5s"`
	TouShiBuildingDestroyProsperity        uint64           `desc:"投石机每次摧毁繁荣度" default:"300"`
	TouShiBuildingBaseHurt                 uint64           `desc:"投石机对每支队伍的基础伤害" default:"100"`
	TouShiBuildingHurtPercent              *data.Amount     `desc:"投石机对每支队伍的伤害加成" default:"10%"`
	TouShiBuildingBaseHurtMaxTroop         uint64           `desc:"投石机最多伤害X支队伍" default:"5"`
	TouShiBuildingBombFlyDuration          time.Duration    `desc:"投石机炮弹飞行时长" default:"4s"`
	DurmDuration                           time.Duration    `desc:"击鼓间隔" default:"300s"`
	DrumStopDuration                       time.Duration    `desc:"准备阶段最后多久停止击鼓" default:"10s"`
	DrumMinBaiZhanLevel                    uint64           `desc:"击鼓百战最低等级" default:"2"`
}

func (d *MingcMiscData) Init(filename string) {
	check.PanicNotTrue(d.StartAfterServerOpen <= d.StartSelfCapitalAfterServerOpen, "%v, StartSelfDuAfterServerOpen:%v 必须 >= StartAfterServerOpen:%v", filename, d.StartSelfCapitalAfterServerOpen, d.StartAfterServerOpen)
	check.PanicNotTrue(d.StartSelfCapitalAfterServerOpen <= d.StartOtherCapitalAfterServerOpen, "%v, StartOtherDuAfterServerOpen:%v 必须 >= StartSelfDuAfterServerOpen:%v", filename, d.StartOtherCapitalAfterServerOpen, d.StartSelfCapitalAfterServerOpen)
}

//gogen:config
type MingcTimeData struct {
	_ struct{} `file:"名城战/时间.txt"`
	_ struct{} `protogen:"true"`

	Id uint64 `desc:"id"`

	Desc string `desc:"文字描述" default:""`

	ApplyAtkStart    string        `desc:"申请挑战开始。格式:星期,时间。例: w7,20h0m"`
	ApplyAtkDuration time.Duration `desc:"申请挑战持续时间"`
	ApplyAstStart    string        `desc:"申请协助开始。格式:星期,时间。例: w7,20h0m"`
	ApplyAstDuration time.Duration
	FightStart       string `desc:"战斗开始。格式:星期,时间。例: w7,20h0m"`
	FightDuration    time.Duration

	ApplyAtkTime *timeutil.WeekDurTime `head:"-" protofield:"-"`
	ApplyAstTime *timeutil.WeekDurTime `head:"-" protofield:"-"`
	FightTime    *timeutil.WeekDurTime `head:"-" protofield:"-"`
	WarTime      *timeutil.WeekDurTime `head:"-" protofield:"-"`
}

func (*MingcTimeData) InitAll(filename string, datas interface {
	GetMingcTimeDataArray() []*MingcTimeData
}) {
	for _, d := range datas.GetMingcTimeDataArray() {
		if t, err := timeutil.BuildWeekDurTime(d.ApplyAtkStart, d.ApplyAtkDuration); err != nil {
			logrus.WithError(err).Panicf("%s 解析时间格式错误 ApplyAtkStart：%v", filename, d.ApplyAtkStart)
		} else {
			d.ApplyAtkTime = t
		}
		if t, err := timeutil.BuildWeekDurTime(d.ApplyAstStart, d.ApplyAstDuration); err != nil {
			logrus.WithError(err).Panicf("%s 解析时间格式错误 ApplyAstStart：%v", filename, d.ApplyAstStart)
		} else {
			d.ApplyAstTime = t
		}
		if t, err := timeutil.BuildWeekDurTime(d.FightStart, d.FightDuration); err != nil {
			logrus.WithError(err).Panicf("%s 解析时间格式错误 FightStart：%v", filename, d.FightStart)
		} else {
			d.FightTime = t
		}

		if t, err := timeutil.BuildWeekDurTime(d.FightStart, d.FightDuration); err != nil {
			logrus.WithError(err).Panicf("%s 解析时间格式错误 FightStart：%v", filename, d.FightStart)
		} else {
			d.FightTime = t
		}

		// 总持续时间
		warDur := d.ApplyAtkTime.Dur + d.ApplyAstTime.Dur + d.FightTime.Dur
		if t, err := timeutil.BuildWeekDurTime(d.ApplyAtkStart, warDur); err != nil {
			logrus.Panicf("%s 解析时间格式错误 FightStart：%v", filename, d.FightStart)
		} else {
			d.WarTime = t
		}
	}
}

func (d *MingcTimeData) GetNextSchedule(ctime time.Time) (s map[shared_proto.MingcWarState][]time.Time, state shared_proto.MingcWarState) {
	s = make(map[shared_proto.MingcWarState][]time.Time)
	applyAtkStartTime, applyAtkEndTime := d.ApplyAtkTime.NextTime(ctime)
	s[shared_proto.MingcWarState_MC_T_APPLY_ATK] = []time.Time{applyAtkStartTime, applyAtkEndTime}
	applyAstStartTime, applyAstEndTime := applyAtkEndTime, applyAtkEndTime.Add(d.ApplyAstTime.Dur)
	s[shared_proto.MingcWarState_MC_T_APPLY_AST] = []time.Time{applyAstStartTime, applyAstEndTime}
	fightStartTime, fightEndTime := applyAstEndTime, applyAstEndTime.Add(d.FightTime.Dur)
	s[shared_proto.MingcWarState_MC_T_FIGHT] = []time.Time{fightStartTime, fightEndTime}
	s[shared_proto.MingcWarState_MC_T_FIGHT_END] = []time.Time{fightEndTime, fightEndTime}

	if ctime.Before(applyAtkStartTime) {
		state = shared_proto.MingcWarState_MC_T_NOT_START
	} else if between(ctime, applyAtkStartTime, applyAtkEndTime) {
		state = shared_proto.MingcWarState_MC_T_APPLY_ATK
	} else if between(ctime, applyAstStartTime, fightStartTime) {
		state = shared_proto.MingcWarState_MC_T_APPLY_AST
	} else if between(ctime, fightStartTime, fightEndTime) {
		state = shared_proto.MingcWarState_MC_T_FIGHT
	} else {
		state = shared_proto.MingcWarState_MC_T_FIGHT_END
	}

	return
}

func between(ctime time.Time, startTime time.Time, endTime time.Time) bool {
	return ctime.After(startTime) && ctime.Before(endTime)
}

//gogen:config
type MingcWarSceneData struct {
	_ struct{} `file:"名城战/场景.txt"`
	_ struct{} `protogen:"true"`

	Id   uint64 `desc:"值和名城id 相同"`
	Desc string `desc:"场景说明"`

	AtkReliveName string `desc:"攻方复活点名字"`
	AtkRelivePosX int
	AtkRelivePosY int
	AtkRelivePos  cb.Cube `head:"-" protofield:"-"`

	AtkHomeName string `desc:"攻方大本营名字"`
	AtkHomePosX int
	AtkHomePosY int
	AtkHomePos  cb.Cube `head:"-" protofield:"-"`

	AtkCastleName []string  `desc:"攻方要塞名字"`
	AtkCastlePosX []int     ` validator:"int,duplicate"`
	AtkCastlePosY []int     ` validator:"int,duplicate"`
	AtkCastlePos  []cb.Cube `head:"-" protofield:"-"`

	AtkGateName []string  `desc:"攻方关卡名字"`
	AtkGatePosX []int     `validator:"int,duplicate"`
	AtkGatePosY []int     ` validator:"int,duplicate"`
	AtkGatePos  []cb.Cube `head:"-" protofield:"-"`

	DefReliveName string `desc:"守方复活点名字"`
	DefRelivePosX int
	DefRelivePosY int
	DefRelivePos  cb.Cube `head:"-" protofield:"-"`

	DefHomeName string `desc:"守方大本营名字"`
	DefHomePosX int
	DefHomePosY int
	DefHomePos  cb.Cube `head:"-" protofield:"-"`

	DefCastleName []string  `desc:"守方要塞名字"`
	DefCastlePosX []int     `validator:"int,duplicate"`
	DefCastlePosY []int     `validator:"int,duplicate"`
	DefCastlePos  []cb.Cube `head:"-" protofield:"-"`

	DefGateName []string  `desc:"守方关卡名字"`
	DefGatePosX []int     `validator:"int,duplicate"`
	DefGatePosY []int     `validator:"int,duplicate"`
	DefGatePos  []cb.Cube `head:"-" protofield:"-"`

	AtkTouShiName []string  `desc:"攻方投石机名字" validator:"string,duplicate"`
	AtkTouShiPosX []int     `validator:"int,duplicate"`
	AtkTouShiPosY []int     `validator:"int,duplicate"`
	AtkTouShiPos  []cb.Cube `head:"-" protofield:"-"`

	DefTouShiName []string  `desc:"守方投石机名字" validator:"string,duplicate"`
	DefTouShiPosX []int     `validator:"int,duplicate"`
	DefTouShiPosY []int     `validator:"int,duplicate"`
	DefTouShiPos  []cb.Cube `head:"-" protofield:"-"`

	AtkFullProsperity uint64 `desc:"攻方总繁荣度" head:"-" `
	DefFullProsperity uint64 `desc:"守方总繁荣度" head:"-" `

	mcWarMap          map[cb.Cube][]cb.Cube
	allPoses          map[cb.Cube]struct{}
	touShiBuildingMap map[cb.Cube][]cb.Cube

	drumDatas []*MingcWarDrumStatData
}

func (d *MingcWarSceneData) DrumStat(times uint64) (succ bool, desc string, stat *shared_proto.SpriteStatProto) {
	b := data.NewSpriteStatBuilder()

	isIncre := rand.Int31n(2) == 1
	for _, dd := range d.drumDatas {
		if !dd.DrumRange.ContainsClosed(times) {
			continue
		}

		if isIncre {
			b.AddDamageIncrePer(dd.DrumDamageIncreRange.Rand())
			desc = dd.DrumDamageIncreDesc
		} else {
			b.AddDamageDecrePer(dd.DrumDamageDecreRange.Rand())
			desc = dd.DrumDamageDecreDesc
		}
		succ = true
		break
	}
	stat = b.Build().Encode4Init()
	return
}

func (d *MingcWarSceneData) GetDest(pos cb.Cube) []cb.Cube {
	return d.mcWarMap[pos]
}

func (d *MingcWarSceneData) Init(filename string, config interface {
	GetMingcBaseData(uint64) *MingcBaseData
	GetMingcWarBuildingData(uint64) *MingcWarBuildingData
	GetMingcWarMapDataArray() []*MingcWarMapData
	GetMingcWarTouShiBuildingTargetDataArray() []*MingcWarTouShiBuildingTargetData
	GetMingcWarDrumStatDataArray() []*MingcWarDrumStatData
}) {
	check.PanicNotTrue(config.GetMingcBaseData(d.Id) != nil, "%v id 必须为名城id。id:%v", filename, d.Id)

	d.AtkRelivePos = cb.XYCube(d.AtkRelivePosX, d.AtkRelivePosY)
	d.AtkHomePos = cb.XYCube(d.AtkHomePosX, d.AtkHomePosY)
	size := u64.Min(u64.FromInt(len(d.AtkGatePosX)), u64.FromInt(len(d.AtkGatePosY)))
	for i := 0; i < int(size); i++ {
		d.AtkGatePos = append(d.AtkGatePos, cb.XYCube(d.AtkGatePosX[i], d.AtkGatePosY[i]))
	}
	size = u64.Min(u64.FromInt(len(d.AtkCastlePosX)), u64.FromInt(len(d.AtkCastlePosY)))
	for i := 0; i < int(size); i++ {
		d.AtkCastlePos = append(d.AtkCastlePos, cb.XYCube(d.AtkCastlePosX[i], d.AtkCastlePosY[i]))
	}

	d.DefRelivePos = cb.XYCube(d.DefRelivePosX, d.DefRelivePosY)
	d.DefHomePos = cb.XYCube(d.DefHomePosX, d.DefHomePosY)
	size = u64.Min(u64.FromInt(len(d.DefGatePosX)), u64.FromInt(len(d.DefGatePosY)))
	for i := 0; i < int(size); i++ {
		d.DefGatePos = append(d.DefGatePos, cb.XYCube(d.DefGatePosX[i], d.DefGatePosY[i]))
	}
	size = u64.Min(u64.FromInt(len(d.DefCastlePosX)), u64.FromInt(len(d.DefCastlePosY)))
	for i := 0; i < int(size); i++ {
		d.DefCastlePos = append(d.DefCastlePos, cb.XYCube(d.DefCastlePosX[i], d.DefCastlePosY[i]))
	}

	// 投石机
	d.touShiBuildingMap = make(map[cb.Cube][]cb.Cube)
	size = u64.Min(u64.FromInt(len(d.AtkTouShiPosX)), u64.FromInt(len(d.AtkTouShiPosY)))
	for i := 0; i < int(size); i++ {
		touShiPos := cb.XYCube(d.AtkTouShiPosX[i], d.AtkTouShiPosY[i])
		d.AtkTouShiPos = append(d.AtkTouShiPos, touShiPos)
		d.touShiBuildingMap[touShiPos] = make([]cb.Cube, 0)
	}
	size = u64.Min(u64.FromInt(len(d.DefTouShiPosX)), u64.FromInt(len(d.DefTouShiPosY)))
	for i := 0; i < int(size); i++ {
		touShiPos := cb.XYCube(d.DefTouShiPosX[i], d.DefTouShiPosY[i])
		d.DefTouShiPos = append(d.DefTouShiPos, touShiPos)
		d.touShiBuildingMap[touShiPos] = make([]cb.Cube, 0)
	}
	for _, tsData := range config.GetMingcWarTouShiBuildingTargetDataArray() {
		if d.Id != tsData.Mingc.Id {
			continue
		}

		if targets, ok := d.touShiBuildingMap[tsData.Pos]; ok {
			d.touShiBuildingMap[tsData.Pos] = append(targets, tsData.Targets...)
		} else {
			logrus.Panicf("mcid:%v 投石机建筑坐标 x:%v y:%v 在 %v 中不存在", d.Id, tsData.PosX, tsData.PosY, filename)
		}
	}

	for pos, targets := range d.touShiBuildingMap {
		x, y := pos.XY()
		check.PanicNotTrue(len(targets) > 0, "投石机建筑 x:%v y:%v 没有配置目标", x, y)
	}

	d.AtkFullProsperity = config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_HOME)).Prosperity
	d.AtkFullProsperity += config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_CASTLE)).Prosperity * uint64(len(d.AtkCastlePos))
	d.AtkFullProsperity += config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_GATE)).Prosperity * uint64(len(d.AtkGatePos))

	d.DefFullProsperity = config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_HOME)).Prosperity
	d.DefFullProsperity += config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_CASTLE)).Prosperity * uint64(len(d.DefCastlePos))
	d.DefFullProsperity += config.GetMingcWarBuildingData(uint64(shared_proto.MingcWarBuildingType_MC_B_GATE)).Prosperity * uint64(len(d.DefGatePos))

	d.allPoses = make(map[cb.Cube]struct{})
	if _, ok := d.allPoses[d.AtkRelivePos]; ok {
		logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, d.AtkRelivePosX, d.AtkRelivePosY)
	} else {
		d.allPoses[d.AtkRelivePos] = struct{}{}
	}
	if _, ok := d.allPoses[d.AtkHomePos]; ok {
		logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, d.AtkHomePosX, d.AtkHomePosY)
	} else {
		d.allPoses[d.AtkHomePos] = struct{}{}
	}
	for _, pos := range d.AtkCastlePos {
		if _, ok := d.allPoses[pos]; ok {
			x, y := pos.XY()
			logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, x, y)
		} else {
			d.allPoses[pos] = struct{}{}
		}
	}
	for _, pos := range d.AtkGatePos {
		if _, ok := d.allPoses[pos]; ok {
			x, y := pos.XY()
			logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, x, y)
		} else {
			d.allPoses[pos] = struct{}{}
		}
	}
	if _, ok := d.allPoses[d.DefRelivePos]; ok {
		logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, d.DefRelivePosX, d.DefRelivePosY)
	} else {
		d.allPoses[d.DefRelivePos] = struct{}{}
	}
	if _, ok := d.allPoses[d.DefHomePos]; ok {
		logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, d.DefHomePosX, d.DefHomePosY)
	} else {
		d.allPoses[d.DefHomePos] = struct{}{}
	}
	for _, pos := range d.DefCastlePos {
		if _, ok := d.allPoses[pos]; ok {
			x, y := pos.XY()
			logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, x, y)
		} else {
			d.allPoses[pos] = struct{}{}
		}
	}
	for _, pos := range d.DefGatePos {
		if _, ok := d.allPoses[pos]; ok {
			x, y := pos.XY()
			logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, x, y)
		} else {
			d.allPoses[pos] = struct{}{}
		}
	}

	for pos := range d.touShiBuildingMap {
		if _, ok := d.allPoses[pos]; ok {
			x, y := pos.XY()
			logrus.Panicf("%v id:%v 坐标重复 x:%v, y:%v", filename, d.Id, x, y)
		} else {
			d.allPoses[pos] = struct{}{}
		}
	}

	d.mcWarMap = make(map[cb.Cube][]cb.Cube)
	for _, mapData := range config.GetMingcWarMapDataArray() {
		if mapData.Mingc.Id == d.Id {
			if _, ok := d.allPoses[mapData.Start]; !ok {
				x, y := mapData.Start.XY()
				logrus.Panicf("%v 坐标不存在 id:%v x:%v y:%v", "名城战/地图.txt", mapData.Id, x, y)
			}
			for _, pos := range mapData.Dests {
				if _, ok := d.allPoses[pos]; !ok {
					x, y := pos.XY()
					logrus.Panicf("%v 坐标不存在 id:%v x:%v y:%v", "名城战/地图.txt", mapData.Id, x, y)
				}
			}

			d.mcWarMap[mapData.Start] = append(d.mcWarMap[mapData.Start], mapData.Dests...)
			// 反着写,会重复但不会丢失，防止配错
			for _, desc := range mapData.Dests {
				d.mcWarMap[desc] = append(d.mcWarMap[desc], mapData.Start)
			}
		}
	}
	check.PanicNotTrue(len(d.mcWarMap) > 0, "%v id:%v 没有配置地图", filename, d.Id)

	d.drumDatas = config.GetMingcWarDrumStatDataArray()
	for i, dd := range d.drumDatas {
		if i == 0 {
			check.PanicNotTrue(dd.DrumRange.min == 1, "鼓舞加成.txt，第一行 drum_min 必须从 1 开始。now:%v", dd.DrumRange.min)
			continue
		}

		check.PanicNotTrue(dd.DrumRange.min == d.drumDatas[i-1].DrumRange.max+1, "鼓舞加成.txt，drum_min 必须等于前一行 durm_max+1。now:%v", dd.DrumRange.min)
	}

}

func (d *MingcWarSceneData) McWarMap() map[cb.Cube][]cb.Cube {
	return d.mcWarMap
}

func (d *MingcWarSceneData) GetTouShiTarget(pos cb.Cube) (target []cb.Cube) {
	return d.touShiBuildingMap[pos]
}

func (d *MingcWarSceneData) CanArrive(start, dest cb.Cube) bool {
	if dests, ok := d.McWarMap()[start]; ok {
		for _, d := range dests {
			if d == dest {
				return true
			}
		}
	}
	// 反着查
	if dests, ok := d.McWarMap()[dest]; ok {
		for _, d := range dests {
			if d == start {
				return true
			}
		}
	}
	return false
}

//gogen:config
type MingcWarBuildingData struct {
	_ struct{} `file:"名城战/据点.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"mingc_war.proto"`

	Id          uint64                            `desc:"id，值和 Type 相同"`
	Type        shared_proto.MingcWarBuildingType `desc:"建筑类型"`
	Prosperity  uint64                            `desc:"繁荣度" validator:"uint"`
	AtkModel    string                            `desc:"攻方模型"`
	DefModel    string                            `desc:"守方模型"`
	WallAtk     bool                              `desc:"城墙能不能攻击"`
	CanBeAtked  bool                              `desc:"能不能被攻击"`
	CombatScene *scene.CombatScene                `desc:"战斗场景" protofield:",%s.Id,string"`
}

func (d *MingcWarBuildingData) Init(filename string) {
	check.PanicNotTrue(d.Id == uint64(d.Type), "%v id 必须等于 type. id:%v", filename, d.Id)
}

//gogen:config
type MingcWarMapData struct {
	_ struct{} `file:"名城战/地图.txt"`
	_ struct{} `protogen:"true"`

	Id     uint64
	Mingc  *MingcBaseData `desc:"名城 id" protofield:",config.U64ToI32(%s.Id),int32"`
	StartX int            `desc:"起始点坐标x"`
	StartY int            `desc:"起始点坐标y"`
	Start  cb.Cube        `head:"-" protofield:"-"`

	DestX []int `validator:"int,duplicate" desc:"所有能直接到达的坐标x"`
	DestY []int `validator:"int,duplicate" desc:"所有能直接到达的坐标y"`

	Dests []cb.Cube `head:"-" protofield:"-"`
}

func (d *MingcWarMapData) Init(filename string) {
	d.Start = cb.XYCube(d.StartX, d.StartY)
	check.PanicNotTrue(len(d.DestX) == len(d.DestY), "%v dest_x 和 dest_y 必须一一对应", filename)
	for i, x := range d.DestX {
		y := d.DestY[i]
		dest := cb.XYCube(x, y)
		check.PanicNotTrue(dest != d.Start, "%v 起点不能和终点相同 x:%v y:%v", filename, x, y)
		d.Dests = append(d.Dests, dest)
	}
}

//gogen:config
type MingcWarNpcData struct {
	_ struct{} `file:"名城战/初始城主.txt"`

	Id           uint64
	Mingc        *MingcBaseData
	Npc          *monsterdata.MonsterMasterData `protofield:"-"`
	Def          bool
	AstDef       bool
	Guild        *MingcWarNpcGuildData
	AiType       server_proto.McWarAiType
	BaiZhanLevel uint64 `validator:"uint"`
}

func (d *MingcWarNpcData) GuildId() int64 {
	return d.Guild.guildId()
}

func RecoverMcWarGuildDataId(gid int64) uint64 {
	return u64.FromInt64(-gid)
}

//gogen:config
type MingcWarNpcGuildData struct {
	_ struct{} `file:"名城战/初始城主联盟.txt"`

	Id       uint64
	Name     string
	FlagName string
	Level    uint64
	Country  uint64
}

func (d *MingcWarNpcGuildData) Init(filename string) {
	check.PanicNotTrue(d.Id > 10000000, "%v id 必须> 10000000.id:%v", filename, d.Id)
}

func (d *MingcWarNpcGuildData) guildId() int64 {
	return -int64(d.Id)
}

func (d *MingcWarNpcGuildData) GuildBasicProto() *shared_proto.GuildBasicProto {
	p := &shared_proto.GuildBasicProto{}
	p.Id = int32(d.guildId())
	p.Name = d.Name
	p.FlagName = d.FlagName
	p.Country = u64.Int32(d.Country)

	return p
}

//gogen:config
type MingcWarTroopLastBeatWhenFailData struct {
	_ struct{} `file:"名城战/舍命一击.txt"`
	_ struct{} `protogen:"true"`

	BaiZhanLevel       uint64       `key:"true"`
	SoliderAmount      uint64       `validator:"uint"`
	HurtPercent        *data.Amount `parser:"data.ParseAmount"`
	AtkBackHurtPercent *data.Amount `parser:"data.ParseAmount"`
}

func (d *MingcWarTroopLastBeatWhenFailData) InitAll(filename string, datas interface {
	GetMingcWarTroopLastBeatWhenFailDataArray() []*MingcWarTroopLastBeatWhenFailData
}) {
	var maxLevel uint64
	for _, d := range datas.GetMingcWarTroopLastBeatWhenFailDataArray() {
		if d.BaiZhanLevel > maxLevel {
			maxLevel = d.BaiZhanLevel
		}
	}

	check.PanicNotTrue(int(maxLevel) == len(datas.GetMingcWarTroopLastBeatWhenFailDataArray()), "%v, 百战军衔必须从1开始依次递增", filename)
}

//gogen:config
type MingcWarTouShiBuildingTargetData struct {
	_ struct{} `file:"名城战/投石机目标.txt"`
	_ struct{} `protogen:"true"`

	Id    uint64
	Mingc *MingcBaseData `desc:"名城 id" protofield:",config.U64ToI32(%s.Id),int32"`

	PosX int     `desc:"投石机坐标X"`
	PosY int     `desc:"投石机坐标Y"`
	Pos  cb.Cube `head:"-" protofield:"-"`

	TargetX []int     `validator:"int,duplicate" desc:"所有目标坐标x"`
	TargetY []int     `validator:"int,duplicate" desc:"所有目标坐标x"`
	Targets []cb.Cube `head:"-" protofield:"-"`
}

func (d *MingcWarTouShiBuildingTargetData) Init(filename string) {
	d.Pos = cb.XYCube(d.PosX, d.PosY)
	check.PanicNotTrue(len(d.TargetX) == len(d.TargetY), "%v target_x 和 target_y 必须一一对应", filename)
	for i, x := range d.TargetX {
		y := d.TargetY[i]
		target := cb.XYCube(x, y)
		check.PanicNotTrue(target != d.Pos, "%v 坐标点不能目标点和相同 x:%v y:%v", filename, x, y)
		d.Targets = append(d.Targets, target)
	}

	check.PanicNotTrue(len(d.Targets) == 4, "%v, 投石机建筑 x:%v y:%v 目标必须为4个", filename, d.PosX, d.PosY)
}

//gogen:config
type MingcWarDrumStatData struct {
	_ struct{} `file:"名城战/鼓舞加成.txt"`

	Id uint64

	DrumMin   uint64
	DrumMax   uint64
	DrumRange *U64Range `head:"-" protofield:"-"`

	DrumDamageIncreDesc  string
	DrumDamageIncreMin   uint64
	DrumDamageIncreMax   uint64
	DrumDamageIncreRange *U64Range `head:"-" protofield:"-"`

	DrumDamageDecreDesc  string
	DrumDamageDecreMin   uint64
	DrumDamageDecreMax   uint64
	DrumDamageDecreRange *U64Range `head:"-" protofield:"-"`
}

func (d *MingcWarDrumStatData) Init(filename string) {
	check.PanicNotTrue(d.DrumMin <= d.DrumMax, "%v, 次数区间 min:%v 不能大于 max:%v", filename, d.DrumMin, d.DrumMax)
	d.DrumRange = NewU64Range(d.DrumMin, d.DrumMax)

	check.PanicNotTrue(d.DrumDamageIncreMin <= d.DrumDamageIncreMax, "%v, 伤害加成区间 min:%v 不能大于 max:%v", filename, d.DrumDamageIncreMin, d.DrumDamageIncreMax)
	d.DrumDamageIncreRange = NewU64Range(d.DrumDamageIncreMin, d.DrumDamageIncreMax)

	check.PanicNotTrue(d.DrumDamageDecreMin <= d.DrumDamageDecreMax, "%v, 伤害减免区间 min:%v 不能大于 max:%v", filename, d.DrumDamageDecreMin, d.DrumDamageDecreMax)
	d.DrumDamageDecreRange = NewU64Range(d.DrumDamageDecreMin, d.DrumDamageDecreMax)
}

func NewU64Range(min, max uint64) *U64Range {
	r := &U64Range{}
	r.min = min
	r.max = max
	return r
}

type U64Range struct {
	min uint64
	max uint64
}

func (r *U64Range) Rand() uint64 {
	return r.min + u64.FromInt64(rand.Int63n(int64(r.max-r.min)))
}

func (r *U64Range) ContainsClosed(n uint64) bool {
	return n >= r.min && n <= r.max
}

//gogen:config
type MingcWarMultiKillData struct {
	_ struct{} `file:"名城战/连斩.txt"`

	MultiKill uint64 `key:"true"` // 连斩数
}
