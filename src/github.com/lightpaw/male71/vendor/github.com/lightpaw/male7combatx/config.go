package combatx

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"math"
)

func NewConfig(proto *shared_proto.CombatConfigProto) (*Config, error) {

	if proto.FramePerSecond <= 0 {
		return nil, errors.Errorf("create combat.Config fail, FramePerSecond<=0")
	}

	if proto.ConfigDenominator <= 0 {
		return nil, errors.Errorf("create combat.Config fail, ConfigDenominator<=0")
	}

	if proto.MinAttackDuration <= 0 {
		return nil, errors.Errorf("create combat.Config fail, MinAttackDuration<=0")
	}

	if proto.MaxAttackDuration <= 0 {
		return nil, errors.Errorf("create combat.Config fail, MaxAttackDuration<=0")
	}

	if proto.MinAttackDuration > proto.MaxAttackDuration {
		return nil, errors.Errorf("create combat.Config fail, proto.MinAttackDuration > proto.MaxAttackDuration")
	}

	if proto.MinMoveSpeed <= 0 {
		return nil, errors.Errorf("create combat.Config fail, MinMoveSpeed<=0")
	}

	if proto.MaxMoveSpeed <= 0 {
		return nil, errors.Errorf("create combat.Config fail, MaxMoveSpeed<=0")
	}

	if proto.MinMoveSpeed > proto.MaxMoveSpeed {
		return nil, errors.Errorf("create combat.Config fail, proto.MinMoveSpeed > proto.MaxMoveSpeed")
	}

	if proto.MinStat == nil {
		return nil, errors.Errorf("create combat.Config fail, proto.MinStat == nil")
	}

	framePerSecond := int(proto.FramePerSecond)
	denominator := float64(proto.ConfigDenominator)

	config := &Config{}

	config.StateMap = make(map[int32]*StateData, len(proto.State))
	for _, s := range proto.State {
		data := newStateData(s, framePerSecond, denominator)
		config.StateMap[data.Id()] = data
	}

	config.SpellMap = make(map[int32]*SpellData, len(proto.Spell))
	for _, s := range proto.Spell {
		data := newSpellData(config.StateMap, s, framePerSecond, denominator)
		config.SpellMap[data.Id()] = data
	}

	config.PassiveSpellMap = make(map[int32]*PassiveSpellData, len(proto.PassiveSpell))
	for _, s := range proto.PassiveSpell {
		data, err := newPassiveSpellData(config.StateMap, config.SpellMap, s, framePerSecond, denominator)
		if err != nil {
			return nil, errors.Wrapf(err, "create combat.Config fail, newPassiveSpellData fail.")
		}
		config.PassiveSpellMap[s.Id] = data
	}

	config.SpellFacadeMap = make(map[int32]*SpellFacadeData, len(proto.SpellIdMap))
	for _, s := range proto.SpellIdMap {
		data := newSpellFacadeData(config.SpellMap, config.PassiveSpellMap, s)
		config.SpellFacadeMap[data.Id] = data
	}

	config.RaceMap = make(map[shared_proto.Race]*RaceData, len(proto.Race))
	for _, race := range proto.Race {
		data, err := newRaceData(config.SpellFacadeMap, race, framePerSecond)
		if err != nil {
			return nil, errors.Wrapf(err, "create combat.Config fail, newRaceData fail.")
		}
		config.RaceMap[race.Race] = data
	}

	config.CaptainSpellMap = make(map[int32]*CaptainSpellData)

	config.CaptainMap = make(map[int32]*CaptainData, len(proto.Captain))
	for _, c := range proto.Captain {
		data, err := newCaptainData(config.SpellFacadeMap, config.RaceMap, c)
		if err != nil {
			return nil, errors.Wrapf(err, "create combat.Config fail, newCaptain fail.")
		}
		config.CaptainMap[c.Id] = data

		for _, star := range data.stars {
			for _, v := range star.UnlockCountSpell {
				config.CaptainSpellMap[v.Id] = v
			}
		}
	}

	config.NamelessCaptainSpellMap = make(map[int32]*CaptainSpellCountData, len(proto.NamelessCaptain))
	for _, c := range proto.NamelessCaptain {
		data, err := newNamelessCaptainSpellData(config.SpellFacadeMap, config.RaceMap, c)
		if err != nil {
			return nil, errors.Wrapf(err, "create combat.Config fail, newNamelessCaptainSpellData fail.")
		}
		config.NamelessCaptainSpellMap[c.Id] = data

		for _, v := range data.UnlockCountSpell {
			config.CaptainSpellMap[v.Id] = v
		}
	}

	config.FramePerSecond = framePerSecond
	config.ConfigDenominator = denominator
	config.IConfigDenominator = int(proto.ConfigDenominator)

	config.MinAttackFrame = IMax(int(proto.MinAttackDuration)*framePerSecond/1000, 1)
	config.MaxAttackFrame = IMax(int(proto.MaxAttackDuration)*framePerSecond/1000, 1)

	config.MinMoveSpeedPerFrame = IMax(int(proto.MinMoveSpeed)/framePerSecond, 5)
	config.MaxMoveSpeedPerFrame = IMax(int(proto.MaxMoveSpeed)/framePerSecond, 5)

	config.MinStat = toF64Stat(proto.MinStat)

	config.MaxFrame = int(proto.MaxDuration) * framePerSecond

	config.scorePercent = proto.GetScorePercent()

	config.CheckMoveFrame = IMax(int(proto.CheckMoveDuration)*framePerSecond/1000, 5)

	config.CritRate = int(proto.CritRate)

	config.Coef = float64(proto.Coef) / denominator

	config.CellLen = int(proto.CellLen)

	config.MaxRage = int(proto.MaxRage)

	config.RageRecoverSpeed = int(proto.RageRecoverSpeed)

	config.AddRagePerHint = int(proto.AddRagePerHint)

	config.AddRageLost1Percent = int(proto.AddRageLost1Percent)

	config.WallWaitFrame = int(proto.WallWaitDuration) * framePerSecond / 1000
	config.WallAttackFixDamageTimes = int(proto.WallAttackFixDamageTimes)
	config.WallBeenHurtLostMaxPercent = float64(proto.WallBeenHurtLostMaxPercent) / denominator
	config.WallSpell = config.SpellMap[proto.WallSpell]
	if config.WallSpell == nil {
		return nil, errors.Errorf("create combat.Config fail, config.WallSpell == nil, spell: %v", proto.WallSpell)
	}

	config.ShortMoveDistance = int(proto.ShortMoveDistance)

	config.InitAttackerX = proto.InitAttackerX
	config.InitDefenserX = proto.InitDefenserX
	config.InitWallX = int(proto.InitWallX)

	// 城墙飞行最低时间
	config.WallDelayMinFrame = IMax(int(proto.WallFlyMinDuration)*framePerSecond/1000, 3)

	return config, nil
}

type Config struct {
	StateMap map[int32]*StateData

	SpellMap map[int32]*SpellData

	PassiveSpellMap map[int32]*PassiveSpellData

	SpellFacadeMap map[int32]*SpellFacadeData

	CaptainMap map[int32]*CaptainData

	NamelessCaptainSpellMap map[int32]*CaptainSpellCountData

	CaptainSpellMap map[int32]*CaptainSpellData

	RaceMap map[shared_proto.Race]*RaceData

	FramePerSecond int

	ConfigDenominator  float64
	IConfigDenominator int

	MinAttackFrame int

	MaxAttackFrame int

	MinMoveSpeedPerFrame int

	MaxMoveSpeedPerFrame int

	MinStat *F64Stat

	MaxFrame int

	scorePercent []int32

	CheckMoveFrame int

	CritRate int

	Coef float64

	CellLen int

	MaxRage int

	RageRecoverSpeed int

	AddRagePerHint int

	AddRageLost1Percent int

	WallWaitFrame int

	WallAttackFixDamageTimes int

	WallBeenHurtLostMaxPercent float64

	WallSpell *SpellData

	ShortMoveDistance int

	InitAttackerX int32
	InitDefenserX int32
	InitWallX     int

	// 城墙攻击飞行最小时间
	WallDelayMinFrame int
}

func (c *Config) getScore(percent int32) (score int32) {
	for i := 0; i < len(c.scorePercent); i++ {
		if percent > c.scorePercent[i] {
			score = int32(i + 1)
		} else {
			break
		}
	}

	return
}

func (c *Config) getCaptainSpell(captainId, star, unlockCount int32, isNameless bool) *CaptainSpellData {
	captainSpellId := newCaptainSpellId(captainId, star, unlockCount, isNameless)
	if csd := c.CaptainSpellMap[captainSpellId]; csd != nil {
		return csd
	}

	if isNameless {
		captain := c.NamelessCaptainSpellMap[captainId]
		if captain != nil {
			return captain.GetSpell(int(unlockCount))
		}
	} else {
		captain := c.CaptainMap[captainId]
		if captain != nil {
			return captain.getStarSpellData(int(star), int(unlockCount))
		}
	}

	return nil
}

func newRaceData(spellFacadeMap map[int32]*SpellFacadeData, proto *shared_proto.RaceDataProto, framePerSecond int) (*RaceData, error) {

	// 兵种技能
	soldierSpell := make([]*SpellFacadeData, len(proto.SoldierSpell))
	for i, id := range proto.SoldierSpell {
		spell := spellFacadeMap[id]
		if spell == nil {
			return nil, errors.Errorf("兵种技能没找到，Race: %v, Spell: %v", proto.Race, proto.SoldierSpell)
		}
		soldierSpell[i] = spell
	}

	data := &RaceData{}
	data.proto = proto
	data.soldierSpell = soldierSpell

	data.priorityMap = make(map[shared_proto.Race]int, len(proto.Priority))
	for i, priorityRace := range proto.Priority {
		// 左边的最大
		data.priorityMap[priorityRace] = 1 + len(proto.Priority) - i

		if i == 0 {
			data.firstPriorityRace = priorityRace
		}
	}

	data.raceCoefMap = make(map[shared_proto.Race]float64)
	for i, raceCoef := range proto.RaceCoef {
		targetRace := shared_proto.Race(i + 1)
		data.raceCoefMap[ targetRace] = float64(raceCoef) / Denominator
	}

	data.viewRange = int(proto.ViewRange)

	data.movePerFrame = IMax(int(proto.MoveSpeed)/framePerSecond, 5)

	data.wallCoef = math.Max(float64(proto.WallCoef)/Denominator, 0.1)

	return data, nil
}

func newCaptainSpellId(id, star, count int32, isNameless bool) int32 {
	return newCaptainSpellId0(newCaptainSpellCountId(id, star, isNameless), count)
}

func newCaptainSpellId0(captainSpellId, count int32) int32 {
	return captainSpellId*100 + count
}

func newCaptainSpellCountId(id, star int32, isNameless bool) int32 {
	if isNameless {
		return id*2 + 1
	} else {
		return (id*100 + star) * 2
	}
}

func newNamelessCaptainSpellData(spellFacadeMap map[int32]*SpellFacadeData, raceMap map[shared_proto.Race]*RaceData, proto *shared_proto.NamelessCaptainDataProto) (*CaptainSpellCountData, error) {
	raceData := raceMap[proto.Race]
	if raceData == nil {
		return nil, errors.Errorf("无名武将职业没找到，%v", proto.Race)
	}

	baseSpell := spellFacadeMap[proto.BaseSpell]
	if baseSpell == nil {
		return nil, errors.Errorf("无名武将普攻没找到，%v", proto.BaseSpell)
	}

	captainSpellId := newCaptainSpellCountId(proto.Id, 0, true)

	captainSpellData, err := newCaptainSpellCountData(spellFacadeMap, baseSpell, nil, captainSpellId, proto.Spell, proto.InitRage)
	if err != nil {
		return nil, errors.Wrapf(err, "无名武将星级技能初始化失败，%v", proto.Id)
	}

	return captainSpellData, nil
}

func newCaptainData(spellFacadeMap map[int32]*SpellFacadeData, raceMap map[shared_proto.Race]*RaceData, proto *shared_proto.CaptainDataProto) (*CaptainData, error) {
	raceData := raceMap[proto.Race]
	if raceData == nil {
		return nil, errors.Errorf("武将职业没找到，%v", proto.Race)
	}

	if len(proto.Star) <= 0 {
		return nil, errors.Errorf("武将星级数据没找到，len(proto.Star) <= 0")
	}

	baseSpell := spellFacadeMap[proto.BaseSpell]
	if baseSpell == nil {
		return nil, errors.Errorf("武将普攻没找到，%v", proto.BaseSpell)
	}

	c := &CaptainData{}
	c.proto = proto

	c.stars = make([]*CaptainSpellCountData, len(proto.Star))
	for i, star := range proto.Star {
		captainSpellId := newCaptainSpellCountId(proto.Id, star.Star, false)
		captainSpellData, err := newCaptainSpellCountData(spellFacadeMap, baseSpell, raceData.soldierSpell, captainSpellId, star.Spell, proto.InitRage)
		if err != nil {
			return nil, errors.Wrapf(err, "武将星级技能初始化失败，%v", proto.Id)
		}

		c.stars[i] = captainSpellData
	}

	return c, nil
}

func newCaptainSpellCountData(spellFacadeMap map[int32]*SpellFacadeData,
	baseSpell *SpellFacadeData, soldierSpell []*SpellFacadeData, captainSpellCountId int32,
	spellIds []int32, initRage int32) (*CaptainSpellCountData, error) {

	c := &CaptainSpellCountData{}
	c.Id = captainSpellCountId

	for i := 0; i <= len(spellIds); i++ {
		captainSpellId := newCaptainSpellId0(captainSpellCountId, int32(i))
		data, err := newCaptainSpellData(spellFacadeMap, baseSpell, soldierSpell, captainSpellId,
			spellIds[:i], initRage)
		if err != nil {
			return nil, err
		}
		c.UnlockCountSpell = append(c.UnlockCountSpell, data)
	}

	return c, nil
}

func newCaptainSpellData(spellFacadeMap map[int32]*SpellFacadeData,
	baseSpell *SpellFacadeData, soldierSpell []*SpellFacadeData, captainSpellId int32,
	spellIds []int32, initRage int32) (*CaptainSpellData, error) {

	c := &CaptainSpellData{}
	c.Id = captainSpellId
	c.initRage = int(initRage)

	var spells []*SpellData
	var passiveSpells []*PassiveSpellData
	addSpell := func(s *SpellFacadeData) {
		if s.Spell != nil {
			spells = append(spells, s.Spell)
		}
		passiveSpells = append(passiveSpells, s.PassiveSpell...)
	}

	// 普攻
	c.baseSpell = baseSpell.Spell
	if c.baseSpell == nil {
		return nil, errors.Errorf("武将普攻没配置主动技能，%v", baseSpell.Id)
	}

	addSpell(baseSpell)

	//兵种技能
	for _, s := range soldierSpell {
		addSpell(s)
	}

	// 星级技能
	for _, sid := range spellIds {
		sf := spellFacadeMap[sid]
		if sf == nil {
			return nil, errors.Errorf("武将技能列表中的技能没找到，%v", sid)
		}

		addSpell(sf)
	}

	c.spells = spells
	c.passiveSpells = passiveSpells

	for _, s := range spells {
		if s.proto.RageSpell {
			c.rageSpell = s
			break
		}
	}

	c.initBeenHurtEffectInc = make(map[int32]float64)
	c.initBeenHurtEffectDec = make(map[int32]float64)

	for _, ps := range passiveSpells {
		switch ps.proto.TriggerType {
		case shared_proto.SpellTriggerType_STBeginRelease:
			c.beginReleaseSpell = append(c.beginReleaseSpell, ps)

		case shared_proto.SpellTriggerType_STFirstHit:
			c.firstAttackSpell = append(c.firstAttackSpell, ps)

		case shared_proto.SpellTriggerType_STFirstHitTarget:
			c.firstAttackTargetSpell = append(c.firstAttackTargetSpell, ps)

		case shared_proto.SpellTriggerType_STHitN:
			if ps.triggerHit > 1 {
				c.timesNSpell = append(c.timesNSpell, ps)
			} else {
				c.times1Spell = append(c.times1Spell, ps)
			}

		case shared_proto.SpellTriggerType_STBeenHit:
			c.beenHurtSpell = append(c.beenHurtSpell, ps)

		case shared_proto.SpellTriggerType_STShieldBroken:
			c.shieldBrokenSpell = append(c.shieldBrokenSpell, ps)
		}

		if ps.proto.SpriteStat != nil {
			if c.initStat == nil {
				c.initStat = &shared_proto.SpriteStatProto{}
			}

			addStatTo(c.initStat, ps.proto.SpriteStat)
		}

		if ps.relivePercent > 0 {
			c.reliveSpell = append(c.reliveSpell, ps)
		}

		c.initRage += ps.rage

		if n := IMin(len(ps.proto.BeenHurtEffectIncType), len(ps.BeenHurtEffectInc)); n > 0 {
			for i := 0; i < n; i++ {
				c.initBeenHurtEffectInc[ps.proto.BeenHurtEffectIncType[i]] += ps.BeenHurtEffectInc[i]
			}
		}

		if n := IMin(len(ps.proto.BeenHurtEffectDecType), len(ps.BeenHurtEffectDec)); n > 0 {
			for i := 0; i < n; i++ {
				c.initBeenHurtEffectDec[ps.proto.BeenHurtEffectDecType[i]] += ps.BeenHurtEffectDec[i]
			}
		}
	}

	return c, nil
}

type CaptainData struct {
	proto *shared_proto.CaptainDataProto

	stars []*CaptainSpellCountData
}

func (c *CaptainData) getStarSpellData(star, unlockCount int) *CaptainSpellData {

	idx := star - 1
	if idx >= 0 && idx < len(c.stars) {
		return c.stars[idx].GetSpell(unlockCount)
	}

	if star > 0 {
		return c.stars[len(c.stars)-1].GetSpell(unlockCount)
	} else {
		return c.stars[0].GetSpell(unlockCount)
	}
}

type CaptainSpellCountData struct {
	Id int32

	UnlockCountSpell []*CaptainSpellData
}

func (c *CaptainSpellCountData) GetSpell(unlockCount int) *CaptainSpellData {
	if unlockCount >= 0 && unlockCount < len(c.UnlockCountSpell) {
		return c.UnlockCountSpell[unlockCount]
	}

	if unlockCount > 0 {
		return c.UnlockCountSpell[len(c.UnlockCountSpell)-1]
	}

	return c.UnlockCountSpell[0]
}

type CaptainSpellData struct {
	Id int32

	baseSpell *SpellData

	rageSpell *SpellData

	// 所有主动技能
	spells []*SpellData

	// 被动技能（全）
	passiveSpells []*PassiveSpellData

	// 开局释放
	beginReleaseSpell []*PassiveSpellData

	// 先战技能
	firstAttackSpell []*PassiveSpellData

	// 首次攻击目标技能
	firstAttackTargetSpell []*PassiveSpellData

	// 攻击N次触发技能
	times1Spell []*PassiveSpellData
	timesNSpell []*PassiveSpellData

	// 被攻击触发
	beenHurtSpell []*PassiveSpellData

	// 护盾被击破触发
	shieldBrokenSpell []*PassiveSpellData

	// 复活技能
	reliveSpell []*PassiveSpellData

	// 所有被动附加的属性
	initStat *shared_proto.SpriteStatProto

	// 所有被动附加的初始怒气
	initRage int

	// 所有被动附加的被打更疼
	initBeenHurtEffectInc map[int32]float64

	// 所有被动附加的被打更轻
	initBeenHurtEffectDec map[int32]float64
}

type RaceData struct {
	proto *shared_proto.RaceDataProto

	soldierSpell []*SpellFacadeData

	priorityMap map[shared_proto.Race]int
	raceCoefMap map[shared_proto.Race]float64

	firstPriorityRace shared_proto.Race

	// 视野
	viewRange int

	// 移动每帧
	movePerFrame int

	wallCoef float64
}

func (c *RaceData) getTargetPriority(targetRace shared_proto.Race) int {
	return c.priorityMap[targetRace]
}

func (c *RaceData) getTroopsCoef(targetRace shared_proto.Race) (coef float64) {
	coef, ok := c.raceCoefMap[targetRace]
	if ok {
		return coef
	}

	return 1 // 没找到，系数为1
}

func newSpellData(stateMap map[int32]*StateData, proto *shared_proto.SpellDataProto, framePerSecond int, denominator float64) *SpellData {
	data := &SpellData{proto: proto}
	data.ReleaseRange = int(proto.ReleaseRange)
	data.CooldownFrame = IMax(int(proto.Cooldown)*framePerSecond/1000, 1)
	data.StrongeFrame = IMax(int(proto.StrongeDuration)*framePerSecond/1000, 1)
	data.HurtRange = int(proto.HurtRange)
	data.HurtCount = int(proto.HurtCount)
	data.Coef = float64(proto.Coef) / denominator

	if proto.FlySpeed > 0 {
		data.FlySpeedPerFrame = IMax(int(proto.FlySpeed)/framePerSecond, 5)
	}
	//data.DamageDelayFrame = int(proto.DamageDelay) * framePerSecond / 1000

	data.Target = newSpellTargetData(proto.Target)

	if n := IMin(len(proto.SelfState), len(proto.SelfStateRate)); n > 0 {
		for i := 0; i < n; i++ {
			state := stateMap[proto.SelfState[i]]
			if state != nil {
				sr := &StateDataWithRate{
					Data:        state,
					TriggerRate: int(proto.SelfStateRate[i]),
				}
				data.SelfState = append(data.SelfState, sr)
			}
		}
	}

	if n := IMin(len(proto.TargetState), len(proto.TargetStateRate)); n > 0 {
		for i := 0; i < n; i++ {
			state := stateMap[proto.TargetState[i]]
			if state != nil {
				sr := &StateDataWithRate{
					Data:        state,
					TriggerRate: int(proto.TargetStateRate[i]),
				}
				data.TargetState = append(data.TargetState, sr)
			}
		}
	}

	data.SelfRage = int(proto.SelfRage)
	data.TargetRage = int(proto.TargetRage)

	return data
}

// 英雄技能数据
type SpellData struct {
	proto *shared_proto.SpellDataProto

	Target *SpellTargetData

	ReleaseRange     int
	CooldownFrame    int
	StrongeFrame     int
	HurtRange        int
	HurtCount        int
	Coef             float64
	FlySpeedPerFrame int
	//DamageDelayFrame int

	SelfState   []*StateDataWithRate
	TargetState []*StateDataWithRate

	SelfRage   int
	TargetRage int
}

func (s *SpellData) Id() int32 {
	return s.proto.Id
}

func isInRace(array []shared_proto.Race, race shared_proto.Race, emptyResult bool) bool {
	if len(array) <= 0 {
		return emptyResult
	}

	for _, r := range array {
		if race == r {
			return true
		}
	}
	return false
}

type StateDataWithRate struct {
	Data *StateData

	TriggerRate int
}

func newStateData(proto *shared_proto.StateDataProto, framePerSecond int, denominator float64) *StateData {

	framePerTick := calFrameFromMs(int(proto.TickDuration), framePerSecond, 1)

	data := &StateData{proto: proto}

	data.framePerTick = framePerTick
	data.StackMaxTimes = int(proto.StackMaxTimes)
	data.TickTimes = int(proto.TickTimes)
	data.MoveSpeedRate = float64(proto.MoveSpeedRate) / denominator
	data.AttackSpeedRate = float64(proto.AttackSpeedRate) / denominator
	data.ShieldRate = float64(proto.ShieldRate) / denominator
	if data.ShieldRate > 0 {
		rate := float64(proto.ShieldEffectRate) / denominator
		if rate <= 0 {
			rate = 1
		}
		data.ShieldEffectRate = math.Max(math.Min(rate, 1), 0.2)
	}
	data.BeenHurtEffectInc = i32a2f64a(proto.BeenHurtEffectInc, denominator)
	data.BeenHurtEffectDec = i32a2f64a(proto.BeenHurtEffectDec, denominator)
	data.DamageCoef = float64(proto.DamageCoef) / denominator
	data.Rage = int(proto.Rage)
	data.RageRecoverRate = float64(proto.RageRecoverRate) / denominator

	return data
}

func i32a2f64a(a []int32, denominator float64) []float64 {
	out := make([]float64, len(a))
	for i, v := range a {
		out[i] = float64(v) / denominator
	}
	return out
}

type StateData struct {
	proto *shared_proto.StateDataProto

	// 每次tick需要多少帧
	framePerTick int

	StackMaxTimes     int
	TickTimes         int
	MoveSpeedRate     float64
	AttackSpeedRate   float64
	ShieldRate        float64
	ShieldEffectRate  float64
	BeenHurtEffectInc []float64
	BeenHurtEffectDec []float64
	DamageCoef        float64
	Rage              int
	RageRecoverRate   float64

	//Id int32
	//
	//StackType int32
	//
	//StackMaxTimes int32
	//
	//TickTimes int32
	//
	//TickDuration int32
	//
	//AddStat    *shared_proto.SpriteStatProto
	//ReduceStat *shared_proto.SpriteStatProto
	//
	//// 移动速度
	//MoveSpeed int
	//
	//// 攻速（只影响普攻）
	//AttackSpeed int
	//
	//// 护盾
	//Shield int
	//
	//// 护盾抵扣伤害比例
	//ShieldRate int
	//
	//// 不可走
	//Unmovable bool
	//
	//// 不可攻击（普攻）
	//NotAttackable bool
	//
	//// 沉默（只能普攻）
	//Silence bool
	//
	//// 晕眩
	//Stun bool
	//
	//// 状态类型 0-无 1-流血 2-中毒 3-燃烧
	//EffectType int
	//
	//// 别人打我，打的更疼
	//BeenHurtEffectType int //  0-无 1-流血 2-中毒 3-燃烧
	//BeenHurtEffectInc  float64
	//
	//// 掉血
	//DamageCoef int
}

func (s *StateData) Id() int32 {
	return s.proto.Id
}

func newState(data *StateData, startFrame, framePerSecond int, caster *Troops, tickDamage int) *State {

	s := &State{
		data:       data,
		stackTimes: 1,
	}

	// 刷新计时
	s.reflushTime(startFrame, caster, tickDamage)

	return s
}

// 状态
type State struct {
	data *StateData

	// 施法者（这个状态谁给你加的）
	caster *Troops

	// 状态开始帧
	startFrame int

	// 结束帧
	endFrame int

	// tick次数
	tickTimes int

	// 下次tick的时间
	nextTickFrame int

	// 当前堆叠层数
	stackTimes int

	// 护盾值
	shieldPerStack int // 每层的护盾值，只在初始的时候计算一次
	shield         int

	// 提前算出来的每次tick掉血
	tickDamage int
}

func (s *State) reflushTime(frame int, caster *Troops, tickDamage int) {
	s.caster = caster
	s.startFrame = frame
	s.endFrame = frame + int(s.data.TickTimes)*s.data.framePerTick
	s.tickTimes = 0
	s.nextTickFrame = s.calNextTickFrame()

	if s.data.ShieldRate > 0 {
		s.shieldPerStack = int(float64(caster.totalLife) * s.data.ShieldRate)

		// 刷新护盾，补充护盾到满值
		if newShield := s.shieldPerStack * s.stackTimes; s.shield < newShield {
			s.shield = newShield
		}
	}

	// 每跳多少伤害
	s.tickDamage = tickDamage
}

func (s *State) tick() int {
	s.tickTimes++
	s.nextTickFrame = s.calNextTickFrame()
	return s.tickTimes
}

func (s *State) getRemainTickTimes() int {
	return IMax(s.data.TickTimes-s.tickTimes, 0)
}

func (s *State) getNextTickFrame() int {
	return s.nextTickFrame
}

func (s *State) calNextTickFrame() int {
	return s.startFrame + (s.tickTimes+1)*s.data.framePerTick
}

func calFrameFromMs(millisecond, framePerSecond, minFrame int) int {
	return IMax(millisecond*framePerSecond/1000, minFrame)
}

func newPassiveSpellData(stateMap map[int32]*StateData, spellData map[int32]*SpellData, proto *shared_proto.PassiveSpellDataProto, framePerSecond int, denominator float64) (*PassiveSpellData, error) {

	data := &PassiveSpellData{}
	data.proto = proto
	data.triggerRate = int(proto.TriggerRate)
	data.triggerHit = int(proto.TriggerHit)
	data.triggerTarget = newSpellTargetData(proto.TriggerTarget)

	if data.proto.TargetCooldown > 0 {
		data.targetCooldownFrame = IMax(int(data.proto.TargetCooldown)*framePerSecond/1000, 1)
	}

	data.selfState = make([]*StateData, 0, len(proto.SelfState))
	for _, id := range proto.SelfState {
		state := stateMap[id]
		if state == nil {
			return nil, errors.Errorf("newPassiveSpell selfState id not found, %d", id)
		}
		data.selfState = append(data.selfState, state)
	}

	data.targetState = make([]*StateData, 0, len(proto.TargetState))
	for _, id := range proto.TargetState {
		state := stateMap[id]
		if state == nil {
			return nil, errors.Errorf("newPassiveSpell targetState id not found, %d", id)
		}
		data.targetState = append(data.targetState, state)
	}

	if proto.Spell != 0 {
		data.spell = spellData[proto.Spell]
		if data.spell == nil {
			return nil, errors.Errorf("newPassiveSpell spell id not found, %d", proto.Spell)
		}
	}

	data.rage = int(proto.Rage)

	data.BeenHurtEffectInc = i32a2f64a(proto.BeenHurtEffectInc, denominator)
	data.BeenHurtEffectDec = i32a2f64a(proto.BeenHurtEffectDec, denominator)

	data.relivePercent = float64(proto.RelivePercent) / denominator

	return data, nil
}

type PassiveSpellData struct {
	proto *shared_proto.PassiveSpellDataProto

	// 触发概率
	triggerRate int

	// 触发类型

	// 触发攻击次数
	triggerHit int

	// 触发目标
	triggerTarget *SpellTargetData

	// 目标触发CD
	targetCooldownFrame int

	// 给自己加状态
	selfState []*StateData

	// 给目标加状态
	targetState []*StateData

	// 触发技能
	spell *SpellData

	// 附加怒气
	rage int

	BeenHurtEffectInc []float64
	BeenHurtEffectDec []float64

	relivePercent float64
}

func newSpellTargetData(proto *shared_proto.SpellTargetDataProto) *SpellTargetData {

	data := &SpellTargetData{
		proto: proto,
	}

	if proto != nil {
		data.targetWhenState = proto.TargetUnmovable || proto.TargetNotAttackable ||
			proto.TargetSilence || proto.TargetStun

		data.shouldCheckTarget = len(proto.TargetRace) > 0 || data.targetWhenState || proto.TargetEffectType > 0
	}

	return data
}

type SpellTargetData struct {
	proto *shared_proto.SpellTargetDataProto

	shouldCheckTarget bool

	targetWhenState bool
}

func (data *SpellTargetData) IsValidTroop(target *Troops) bool {

	if !data.shouldCheckTarget {
		return true
	}

	if !isInRace(data.proto.TargetRace, target.getRace(), true) {
		return false
	}

	if !data.isTargetWhen(target) {
		return false
	}

	// 效果状态
	if data.proto.TargetEffectType > 0 &&
		!target.hasEffectState(data.proto.TargetEffectType) {
		return false
	}

	return true
}

func (data *SpellTargetData) isTargetWhen(target *Troops) bool {

	if !data.targetWhenState {
		return true
	}

	if data.proto.TargetUnmovable && target.isUnmovable() {
		return true
	}

	if data.proto.TargetNotAttackable && target.isNotAttackable() {
		return true
	}

	if data.proto.TargetSilence && target.isSilence() {
		return true
	}

	if data.proto.TargetStun && target.isStun() {
		return true
	}

	return false
}

func newSpellFacadeData(spellMap map[int32]*SpellData, passiveSpellMap map[int32]*PassiveSpellData, proto *shared_proto.SpellIdProto) *SpellFacadeData {

	s := &SpellFacadeData{}
	s.Id = proto.Id
	s.Spell = spellMap[proto.Spell]

	for _, v := range proto.PassiveSpell {
		ps := passiveSpellMap[v]
		if ps != nil {
			s.PassiveSpell = append(s.PassiveSpell, passiveSpellMap[v])
		}
	}

	return s
}

type SpellFacadeData struct {
	Id int32

	Spell *SpellData

	PassiveSpell []*PassiveSpellData
}
