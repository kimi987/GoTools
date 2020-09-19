package combatx

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/imath"
	"math"
	"sort"
	"math/rand"
)

func newTroopsArray(genIndex func() int32, isAttacker bool, tps []*shared_proto.CombatTroopsProto, config *Config, misc *misc) ([]*Troops, error) {
	out := make([]*Troops, 0, len(tps))
	for i, v := range tps {
		troops, err := newTroops(genIndex(), isAttacker, v, config, misc)
		if err != nil {
			return nil, errors.Wrapf(err, "Troops[%v]", i)
		}

		out = append(out, troops)
	}

	return out, nil
}

func newTroops(index int32, isAttacker bool, p *shared_proto.CombatTroopsProto, config *Config, misc *misc) (*Troops, error) {
	if index <= 0 {
		return nil, errors.Errorf("Troops.index should > 0, v: %v", index)
	}

	c := p.Captain
	if c == nil {
		return nil, errors.Errorf("Troops.Captain == nil")
	}

	captainSpellData := config.getCaptainSpell(c.CaptainId, c.Star, c.UnlockSpellCount, c.IsNameless)
	if captainSpellData == nil {
		return nil, errors.Errorf("Troops.Soldier captainData == nil,  cid: %v star: %v nameless: %v", c.CaptainId, c.Star, c.IsNameless)
	}

	if c.Soldier <= 0 {
		return nil, errors.Errorf("Troops.Soldier should > 0,  v: %v", c.Soldier)
	}

	if c.TotalSoldier < c.Soldier {
		return nil, errors.Errorf("Troops.TotalSoldier should > Troops.Soldier,  v: %v", c.TotalSoldier)
	}

	if c.LifePerSoldier <= 0 {
		return nil, errors.Errorf("Troops.LifePerSoldier should > 0,  v: %v", c.LifePerSoldier)
	}

	raceData := config.RaceMap[c.Race]
	if raceData == nil {
		return nil, errors.Errorf("Troops.race data not found, v: %v", c.Race)
	}

	posX := p.X
	posY := p.Y
	if p.FightIndex > 0 {
		xIndex := I32MinMax(c.XIndex, 0, 2)

		offset := xIndex * int32(config.CellLen)

		// 根据index设置位置
		if isAttacker {
			posX = config.InitAttackerX + offset
		} else {
			posX = config.InitDefenserX - offset
		}

		posY = (p.FightIndex - 1) * int32(config.CellLen)
	} else {
		// 自己设置位置
		if !(0 <= posX && posX < misc.request.MapXLen) {
			return nil, errors.Errorf("Troops.X not in [0 - %v), v: %v", misc.mapXLen, posX)
		}

		if !(0 <= posY && posY < misc.request.MapYLen) {
			return nil, errors.Errorf("Troops.Y not in [0 - %v), v: %v", misc.mapYLen, posY)
		}
	}

	// 技能
	spellList := newSpellList(captainSpellData)

	baseStat := c.TotalStat
	if captainSpellData.initStat != nil {
		baseStat = &shared_proto.SpriteStatProto{}
		addStatTo(baseStat, c.TotalStat)
		addStatTo(baseStat, captainSpellData.initStat)
	}

	//p.Index = index
	initProto := &shared_proto.CombatTroopsInitProto{
		Index: index,
		X:     posX,
		Y:     posY,
	}

	initRage := 0 + captainSpellData.initRage
	rage := newRage(initRage, config.MaxRage, config.RageRecoverSpeed, config.FramePerSecond)
	if spellList.HasRageSpell() {
		initProto.Rage = int32(rage.rage1000)
		initProto.Recover = int32(rage.recoverPerFrame1000)
	}

	out := &Troops{
		TroopState:            newTroopState(),
		TroopRage:             rage,
		index:                 index,
		isAttacker:            isAttacker,
		proto:                 p,
		initProto:             initProto,
		raceData:              raceData,
		movePerFrame:          raceData.movePerFrame,
		spellList:             spellList,
		beenHurtEffectIncCoef: make(map[int32]float64),
		beenHurtEffectDecCoef: make(map[int32]float64),
	}

	for k, v := range captainSpellData.initBeenHurtEffectInc {
		out.beenHurtEffectIncCoef[k] = v
	}

	for k, v := range captainSpellData.initBeenHurtEffectDec {
		out.beenHurtEffectDecCoef[k] = v
	}

	out.lifePerSoldier = int(c.LifePerSoldier)
	out.totalLife = IMax(int(c.TotalSoldier*c.LifePerSoldier), 1)

	out.baseStat = baseStat
	out.setTotalStat(out.baseStat, config.MinStat)

	// 状态

	// 初始血量
	out.setLife(int(c.LifePerSoldier) * int(c.Soldier))
	out.totalLife = IMax(out.totalLife, out.life)

	return out, nil
}

type F64Stat struct {
	// 属性 float64 攻防体敏
	attack           float64
	defense          float64
	strength         float64
	dexterity        float64
	damageIncrePer   float64 // 伤害增加百分比
	damageDecrePer   float64 // 伤害减少百分比
	beenHurtIncrePer float64 // 易伤增加百分比
	beenHurtDecrePer float64 // 易伤减少百分比
}

func (s *F64Stat) copy() *F64Stat {
	ns := &F64Stat{}
	ns.attack = s.attack
	ns.defense = s.defense
	ns.strength = s.strength
	ns.dexterity = s.dexterity
	ns.damageIncrePer = s.damageIncrePer
	ns.damageDecrePer = s.damageDecrePer
	ns.beenHurtIncrePer = s.beenHurtIncrePer
	ns.beenHurtDecrePer = s.beenHurtDecrePer

	return ns
}

func (s *F64Stat) NewAvgStat() *F64Stat {
	avg := (s.attack + s.defense + s.strength + s.dexterity) / 4

	ns := s.copy()
	ns.attack = avg
	ns.defense = avg
	ns.strength = avg
	ns.dexterity = avg

	return ns
}

// 武将战斗单位
type Troops struct {
	*TroopState // 状态数据
	*TroopRage  // 怒气

	index int32

	isAttacker bool

	proto     *shared_proto.CombatTroopsProto
	initProto *shared_proto.CombatTroopsInitProto

	raceData *RaceData

	baseStat  *shared_proto.SpriteStatProto
	totalStat *F64Stat // 总属性
	avgStat   *F64Stat // 平均属性

	// 移动速度（每帧）
	movePerFrame int

	lifePerSoldier int

	// 技能数据
	spellList *SpellList

	// 被打伤害加深系数
	beenHurtEffectIncCoef map[int32]float64

	// 被打伤害加深系数
	beenHurtEffectDecCoef map[int32]float64

	// 下面是战斗过程计算使用的变量
	x, y int

	rushTarget       *Troops   // 对线目标
	rushUpDownTarget []*Troops // 对线目标楼上楼下

	// 下一次检查移动目标的Frame
	nextCheckMoveTargetFrame int

	currentPath *shared_proto.TroopMoveActionProto

	currentPathStartFrame int32

	life int

	markLifePercent int // 记录血量百分比，用于排序，排序前设置

	totalLife int

	soldier int

	killSoldier int

	// 硬直结束帧数
	strongeEndFrame int

	delayDamage          []*DelayDamage
	nextDelayDamageFrame int

	// 释放普攻之后，设置为true
	// 移动时候检查，如果为true，
	// 如果目标位置跟当前位置<=临界距离，则使用shortMove
	// 如果目标位置跟当前位置>临界距离，则将isShortMove设置为false，使用move
	checkShortMove bool

	// 复活技能序号
	reliveIndex int
}

func (t *Troops) isFar() bool {
	return t.raceData.proto.IsFar
}

func (t *Troops) isBaseSpell(data *SpellData) bool {
	return t.spellList.data.baseSpell == data
}

func (t *Troops) MarkLifePercent() {
	t.markLifePercent = t.life * 100 / t.totalLife
}

func (t *Troops) newReleaseSpellEffectProtoIfNil(proto *shared_proto.TroopReleaseSpellEffectProto) *shared_proto.TroopReleaseSpellEffectProto {
	if proto != nil {
		return proto
	}

	return &shared_proto.TroopReleaseSpellEffectProto{
		TargetIndex: t.index,
		TargetX:     int32(t.x),
		TargetY:     int32(t.y),
	}
}

type spellcaster interface {
	addKillSoldier(toAdd int)
}

type DelayDamage struct {
	caster spellcaster

	damage int

	effectFrame int

	spellId int32
}

func (target *Troops) addDelayDamage(caster spellcaster, damage, effectFrame int, spellId int32) {
	proto := &DelayDamage{
		caster:      caster,
		damage:      damage,
		effectFrame: effectFrame,
		spellId:     spellId,
	}

	target.nextDelayDamageFrame = IMin(target.nextDelayDamageFrame, effectFrame)

	for i, v := range target.delayDamage {
		if v == nil {
			target.delayDamage[i] = proto
			return
		}
	}

	target.delayDamage = append(target.delayDamage, proto)
}

func newTroopState() *TroopState {
	return &TroopState{
		stateMap:         make(map[int32]*State),
		effectStateCount: make(map[int32]int),
	}
}

type TroopState struct {
	// 状态列表
	stateMap map[int32]*State

	nextUpdateStateFrame int

	isAttackSpeedChanged bool
	isMoveSpeedChanged   bool
	isStatChanged        bool
	isRageRecoverChanged bool

	unmovableCount     int // 不可走
	notAttackableCount int // 不能普攻（可以放技能）
	silenceCount       int // 不能放技能（可以普攻）
	stunCount          int // 晕眩
	shieldCount        int // 护盾个数

	effectStateCount map[int32]int
}

func (c *Troops) hasEffectState(t int32) bool {
	return c.effectStateCount[t] > 0
}

func (c *Troops) isUnmovable() bool {
	return c.unmovableCount > 0
}

func (c *Troops) isNotAttackable() bool {
	return c.notAttackableCount > 0
}

func (c *Troops) isSilence() bool {
	return c.silenceCount > 0
}

func (c *Troops) isStun() bool {
	return c.stunCount > 0
}

func (c *Troops) hasShield() bool {
	return c.shieldCount > 0
}

func (c *Troops) doAddState(toAdd *StateData) {
	if toAdd.proto.Unmovable {
		c.unmovableCount++
	}

	if toAdd.proto.NotAttackable {
		c.notAttackableCount++
	}

	if toAdd.proto.Silence {
		c.silenceCount++
	}

	if toAdd.proto.Stun {
		c.stunCount++
	}

	if toAdd.proto.ShieldRate > 0 {
		c.shieldCount++
	}

	if toAdd.proto.EffectType > 0 {
		c.effectStateCount[toAdd.proto.EffectType]++
	}

	if n := IMin(len(toAdd.proto.BeenHurtEffectIncType), len(toAdd.BeenHurtEffectInc)); n > 0 {
		for i := 0; i < n; i++ {
			c.beenHurtEffectIncCoef[toAdd.proto.BeenHurtEffectIncType[i]] += toAdd.BeenHurtEffectInc[i]
		}
	}

	if n := IMin(len(toAdd.proto.BeenHurtEffectDecType), len(toAdd.BeenHurtEffectDec)); n > 0 {
		for i := 0; i < n; i++ {
			c.beenHurtEffectDecCoef[toAdd.proto.BeenHurtEffectDecType[i]] += toAdd.BeenHurtEffectDec[i]
		}
	}
}

func (c *Troops) removeState(toRemove *StateData) {
	delete(c.stateMap, toRemove.Id())

	if toRemove.proto.Unmovable {
		c.unmovableCount--
	}

	if toRemove.proto.NotAttackable {
		c.notAttackableCount--
	}

	if toRemove.proto.Silence {
		c.silenceCount--
	}

	if toRemove.proto.Stun {
		c.stunCount--
	}

	if toRemove.proto.ShieldRate > 0 {
		c.shieldCount--
	}

	if toRemove.proto.EffectType > 0 {
		c.effectStateCount[toRemove.proto.EffectType]--
	}

	if n := IMin(len(toRemove.proto.BeenHurtEffectIncType), len(toRemove.BeenHurtEffectInc)); n > 0 {
		for i := 0; i < n; i++ {
			c.beenHurtEffectIncCoef[toRemove.proto.BeenHurtEffectIncType[i]] -= toRemove.BeenHurtEffectInc[i]
		}
	}

	if n := IMin(len(toRemove.proto.BeenHurtEffectDecType), len(toRemove.BeenHurtEffectDec)); n > 0 {
		for i := 0; i < n; i++ {
			c.beenHurtEffectDecCoef[toRemove.proto.BeenHurtEffectDecType[i]] -= toRemove.BeenHurtEffectDec[i]
		}
	}

	if toRemove.AttackSpeedRate != 0 {
		c.isAttackSpeedChanged = true
	}

	if toRemove.MoveSpeedRate != 0 {
		c.isMoveSpeedChanged = true
	}

	if toRemove.proto.ChangeStat != nil {
		c.isStatChanged = true
	}
}

func (c *Troops) addState(toAdd *StateData, frame, framePerSecond int, caster *Troops, tickDamage int) *State {

	state := c.stateMap[toAdd.Id()]
	if state == nil {
		state = newState(toAdd, frame, framePerSecond, caster, tickDamage)
		c.stateMap[toAdd.Id()] = state
		c.doAddState(toAdd)
	} else {
		switch toAdd.proto.StackType {
		case shared_proto.StateStackType_SSNone:
			// 不堆叠
			return nil
		case shared_proto.StateStackType_SSReplace:
			// 刷新时间
			state.reflushTime(frame, caster, tickDamage)

		case shared_proto.StateStackType_SSStack:

			// 堆叠
			if state.stackTimes < int(toAdd.StackMaxTimes) {
				state.stackTimes++
			}

			// 刷新时间
			state.reflushTime(frame, caster, tickDamage)
		}
	}

	if toAdd.AttackSpeedRate != 0 {
		c.isAttackSpeedChanged = true
	}

	if toAdd.MoveSpeedRate != 0 {
		c.isMoveSpeedChanged = true
	}

	if toAdd.proto.ChangeStat != nil {
		c.isStatChanged = true
	}

	return state
}

func (c *Troops) updateAttackSpeed(framePerSecond, min, max int) {

	// 重新计算攻速
	var rate float64 = 1
	for _, state := range c.stateMap {
		rate += state.data.AttackSpeedRate
	}

	cooldown := int(float64(c.spellList.data.baseSpell.proto.Cooldown) * rate)
	c.spellList.baseSpellCooldownFrame = IMinMax(cooldown*framePerSecond/1000, min, max)

}

func (c *Troops) updateMoveSpeed(framePerSecond, min, max int) bool {

	// 重新计算移动速度
	var rate float64 = 1
	for _, state := range c.stateMap {
		rate += state.data.MoveSpeedRate
	}

	moveSpeed := int(float64(c.raceData.proto.MoveSpeed) * rate)
	c.movePerFrame = IMinMax(moveSpeed/framePerSecond, min, max)

	return true
}

func (c *Troops) updateRageRecover(baseSpeed, framePerSecond int) {

	// 重新计算攻速
	var rate float64 = 1
	for _, state := range c.stateMap {
		rate += state.data.RageRecoverRate
	}

	c.updateRecoverPerSecond(baseSpeed, framePerSecond, rate)
}

func (c *Troops) updateTotalStat(minStat *F64Stat) {

	totalStat := &shared_proto.SpriteStatProto{}
	addStatTo(totalStat, c.baseStat)

	for _, state := range c.stateMap {
		if state.data.proto.ChangeStat != nil {
			doChangeStat(totalStat, state.data.proto.ChangeStat, state.data.proto.IsAddStat)
		}
	}

	c.setTotalStat(totalStat, minStat)
}

func (c *Troops) setTotalStat(totalStat *shared_proto.SpriteStatProto, minStat *F64Stat) {
	newStat := toF64Stat(totalStat)
	newStat.attack = math.Max(newStat.attack, minStat.attack)
	newStat.defense = math.Max(newStat.defense, minStat.defense)
	newStat.strength = math.Max(newStat.strength, minStat.strength)
	newStat.dexterity = math.Max(newStat.dexterity, minStat.dexterity)
	newStat.damageIncrePer = math.Max(newStat.damageIncrePer, minStat.damageIncrePer)
	newStat.damageDecrePer = math.Max(newStat.damageDecrePer, minStat.damageDecrePer)
	newStat.beenHurtIncrePer = math.Max(newStat.beenHurtIncrePer, minStat.beenHurtIncrePer)
	newStat.beenHurtDecrePer = math.Max(newStat.beenHurtDecrePer, minStat.beenHurtDecrePer)

	c.totalStat = newStat
	c.avgStat = newStat.NewAvgStat()
}

func toF64Stat(spriteStat *shared_proto.SpriteStatProto) *F64Stat {
	f64Stat := &F64Stat{}
	f64Stat.attack = float64(spriteStat.Attack)
	f64Stat.defense = float64(spriteStat.Defense)
	f64Stat.strength = float64(spriteStat.Strength)
	f64Stat.dexterity = float64(spriteStat.Dexterity)
	f64Stat.damageIncrePer = float64(spriteStat.DamageIncrePer) / Denominator
	f64Stat.damageDecrePer = float64(spriteStat.DamageDecrePer) / Denominator
	f64Stat.beenHurtIncrePer = float64(spriteStat.BeenHurtIncrePer) / Denominator
	f64Stat.beenHurtDecrePer = float64(spriteStat.BeenHurtDecrePer) / Denominator

	return f64Stat
}

func doChangeStat(totalStat, toChange *shared_proto.SpriteStatProto, isAdd bool) {
	if isAdd {
		addStatTo(totalStat, toChange)
	} else {
		reduceStatFrom(totalStat, toChange)
	}
}

func addStatTo(b, toAdd *shared_proto.SpriteStatProto) {
	b.Attack += toAdd.Attack
	b.Defense += toAdd.Defense
	b.Strength += toAdd.Strength
	b.Dexterity += toAdd.Dexterity
	b.DamageIncrePer += toAdd.DamageIncrePer
	b.DamageDecrePer += toAdd.DamageDecrePer
	b.BeenHurtIncrePer += toAdd.BeenHurtIncrePer
	b.BeenHurtDecrePer += toAdd.BeenHurtDecrePer
}

func reduceStatFrom(b, toReduce *shared_proto.SpriteStatProto) {
	b.Attack -= toReduce.Attack
	b.Defense -= toReduce.Defense
	b.Strength -= toReduce.Strength
	b.Dexterity -= toReduce.Dexterity
	b.DamageIncrePer -= toReduce.DamageIncrePer
	b.DamageDecrePer -= toReduce.DamageDecrePer
	b.BeenHurtIncrePer -= toReduce.BeenHurtIncrePer
	b.BeenHurtDecrePer -= toReduce.BeenHurtDecrePer
}

func (c *Troops) clearRushTarget() {
	c.rushTarget = nil
	c.rushUpDownTarget = nil
}

func (c *Troops) resetPosition() {
	c.x = int(c.initProto.X)
	c.y = int(c.initProto.Y)
}

func (c *Troops) getRaceData() *RaceData {
	return c.raceData
}

func (c *Troops) getRace() shared_proto.Race {
	return c.getRaceData().proto.Race
}

func (c *Troops) getFirstPriorityRace() shared_proto.Race {
	return c.raceData.firstPriorityRace
}

func (c *Troops) hasStepToMove() bool {
	return c.currentPath != nil
}

func (c *Troops) reduceLife(toReduce int) int {
	oldSoldier := c.soldier
	c.setLife(c.life - toReduce)
	return oldSoldier - c.soldier
}

func (c *Troops) setLife(toSet int) int {
	c.life = imath.Max(toSet, 0)
	c.soldier = (c.life + c.lifePerSoldier - 1) / c.lifePerSoldier
	return c.life
}

func (c *Troops) getLifePercent() int {
	return c.life * 100 / c.totalLife
}

func (c *Troops) getSoldier() int {
	return c.soldier
}

func (c *Troops) isAlive() bool {
	return c.life > 0
}

func (c *Troops) isInRange(x, y, hurtRange int) bool {
	return IsInRange(c.x, c.y, x, y, hurtRange)
}

func (c *Troops) canSee(x, y int) bool {
	return IsInRange(c.x, c.y, x, y, c.raceData.viewRange)
}

func (c *Troops) addKillSoldier(toAdd int) {
	if c != nil {
		c.killSoldier += toAdd
	}
}

type TroopsSpeedSlice []*Troops

func (p TroopsSpeedSlice) Len() int      { return len(p) }
func (p TroopsSpeedSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p TroopsSpeedSlice) Less(i, j int) bool {
	t1 := p[i]
	t2 := p[j]

	// 机动力从高到低排序
	if t1.movePerFrame != t2.movePerFrame {
		return t1.movePerFrame > t2.movePerFrame
	}

	// 机动力相同时，攻方武将按照编号从小到大排在前面，守方武将按照编号从小到大排在后面
	if t1.isAttacker != t2.isAttacker {
		return t1.isAttacker
	}

	return t1.index < t2.index
}

func newSpellList(data *CaptainSpellData) *SpellList {
	list := &SpellList{}
	list.data = data

	list.baseSpellCooldownFrame = data.baseSpell.CooldownFrame

	list.spells = make([]*Spell, len(data.spells))
	for i, v := range data.spells {
		spell := newSpell(v)
		list.spells[i] = spell

		if data.baseSpell == v {
			list.baseSpell = spell
		}
	}

	list.spellTargetMap = make(map[int32]struct{})
	list.triggerCdMap = make(map[int32]int)

	return list
}

func newRage(initRage, maxRage, recoverPerSecond, framePerSecond int) *TroopRage {

	recoverPerFrame1000 := IMax(recoverPerSecond*I1000/framePerSecond, 1)
	r := &TroopRage{
		rage1000:            initRage * I1000,
		maxRage1000:         maxRage * I1000,
		recoverPerFrame1000: recoverPerFrame1000,
	}
	r.truncate()

	return r
}

const (
	I1000 = 1000
)

type TroopRage struct {
	rage1000            int
	maxRage1000         int
	recoverPerFrame1000 int

	rageMark         int
	rageRecorverMark int
}

func (r *TroopRage) Rage() int {
	return r.rage1000 / I1000
}

func (r *TroopRage) RecoverRage() {
	r.rage1000 += r.recoverPerFrame1000
	r.truncate()
}

func (r *TroopRage) isFullRage() bool {
	return r.rage1000 >= r.maxRage1000
}

func (r *TroopRage) addRage(toAdd int) {
	r.rage1000 += toAdd * I1000
	r.truncate()
}

func (r *TroopRage) truncate() {
	r.rage1000 = IMin(r.rage1000, r.maxRage1000)
}

func (s *TroopRage) updateRecoverPerSecond(toSet, framePerSecond int, recoverRate float64) {
	newAmount := toSet * I1000 / framePerSecond
	if recoverRate != 0 {
		newAmount = int(float64(newAmount) * (1 + recoverRate))
	}

	s.recoverPerFrame1000 = IMax(newAmount, 1)
}

func (s *TroopRage) ClearRage() {
	s.rage1000 = 0
}

func (r *TroopRage) MarkRage() {
	r.rageMark = r.rage1000
	r.rageRecorverMark = r.recoverPerFrame1000
}

func (r *TroopRage) newRageUpdateProtoIfChanged() *shared_proto.TroopRageUpdateProto {
	if r.rageMark != r.rage1000 || r.rageRecorverMark != r.recoverPerFrame1000 {
		return r.newRageUpdateProto()
	}
	return nil
}

func (s *TroopRage) newRageUpdateProto() *shared_proto.TroopRageUpdateProto {
	return &shared_proto.TroopRageUpdateProto{
		Rage:    int32(s.rage1000),
		Recover: int32(s.recoverPerFrame1000),
	}
}

type SpellList struct {
	data *CaptainSpellData

	// 普攻技能
	baseSpellCooldownFrame int

	spells    []*Spell
	baseSpell *Spell

	// 是否使用过普攻（触发先战技能）
	baseSpellUsed bool

	// 使用普攻攻击同一目标的次数（攻击N次触发技能）
	baseSpellHurtTimes int

	// 上一次普攻伤害目标
	baseSpellHurtTarget int32

	// 技能目标列表
	spellTargetMap map[int32]struct{}

	triggerCdMap map[int32]int
}

func newSpell(data *SpellData) *Spell {
	spell := &Spell{}
	spell.data = data

	return spell
}

// 英雄技能数据
type Spell struct {
	data *SpellData

	nextReleaseFrame int

	lastTarget *Troops
}

func (s *SpellList) HasRageSpell() bool {
	return s.data.rageSpell != nil
}

func (s *SpellList) IsInTriggerCd(id int32, frame int) bool {
	nextTriggerFrame := s.triggerCdMap[id]
	return frame < nextTriggerFrame
}

func (s *SpellList) SetTriggerCd(id int32, toSet int) {
	s.triggerCdMap[id] = toSet
}

func sortTroopsByMaxLifePercent(troops []*Troops) {
	sort.Sort(sort.Reverse(LifePercentTroopSlice(troops)))
}

func sortTroopsByMinLifePercent(troops []*Troops) {
	sort.Sort(LifePercentTroopSlice(troops))
}

func sortTroopsByMaxRage(troops []*Troops) {
	sort.Sort(sort.Reverse(RageTroopSlice(troops)))
}

func sortTroopsByMinRage(troops []*Troops) {
	sort.Sort(RageTroopSlice(troops))
}

type LifePercentTroopSlice []*Troops

func (p LifePercentTroopSlice) Less(i, j int) bool { return p[i].markLifePercent < p[j].markLifePercent }
func (p LifePercentTroopSlice) Len() int           { return len(p) }
func (p LifePercentTroopSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type RageTroopSlice []*Troops

func (p RageTroopSlice) Less(i, j int) bool { return p[i].rage1000 < p[j].rage1000 }
func (p RageTroopSlice) Len() int           { return len(p) }
func (p RageTroopSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func mixTroops(list []*Troops, count int, random *rand.Rand) {

	n := len(list)
	count = IMin(n, count)

	for i := range list {
		if i >= count {
			break
		}

		swap := i + random.Intn(n-i)
		list[i], list[swap] = list[swap], list[i]
	}
	sort.Sort(LifePercentTroopSlice(list))
}
