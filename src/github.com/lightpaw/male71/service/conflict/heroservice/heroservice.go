package heroservice

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/util/lock"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/config/kv"
	"runtime/debug"
)

//gogen:iface
type HeroDataService struct {
	db                  iface.DbService
	heroSnapshotService iface.HeroSnapshotService
	world               iface.WorldService
	lockService         *lock.LockService
}

func NewHeroDataService(serverConfig *kv.IndividualServerConfig, db iface.DbService, heroSnapshotService iface.HeroSnapshotService, world iface.WorldService) *HeroDataService {
	lockService := lock.NewLockService(&lockProvider{db: db}, serverConfig.HeroEvictNotAccessedDuration)
	return &HeroDataService{
		db:                  db,
		lockService:         lockService,
		heroSnapshotService: heroSnapshotService,
		world:               world,
	}
}

//type HeroUnlocker func(**sharedherodata.SharedHeroData)
//func(c **sharedherodata.SharedHeroData) { *c = nil; unlock() }

type hero_locker struct {
	s      *HeroDataService
	heroId int64
}

func (s *hero_locker) Func(f herolock.Func) {
	s.s.Func(s.heroId, f)
}

func (s *hero_locker) FuncNotError(f herolock.FuncNotError) bool {
	return s.s.FuncNotError(s.heroId, f)
}

func (s *hero_locker) FuncWithSend(f herolock.SendFunc, sender sender.ClosableSender) (hasError bool) {
	return s.s.funcWithSend(s.heroId, f, sender)
}

func (s *HeroDataService) NewHeroLocker(id int64) herolock.HeroLocker {

	return &hero_locker{
		s:      s,
		heroId: id,
	}
}

func (s *HeroDataService) Exist(id int64) (bool, error) {

	_, unlock, err := s.lock(id)
	if err != nil {
		if err == lock.ErrEmpty {
			return false, nil
		}

		logrus.WithError(err).Errorf("HeroDataService.Exist() lockSharedHero error, %v", id)
		return false, err
	}
	defer unlock()

	return true, nil
}

var errNpc = errors.Errorf("Id居然是个Npc")

func (s *HeroDataService) lock(id int64) (*entity.Hero, lock.Unlocker, error) {
	if npcid.IsNpcId(id) {
		return nil, nil, errNpc
	}

	result, unlock, err := s.lockService.Lock(id)
	if err != nil {
		return nil, nil, err
	}
	return result.(*entity.Hero), unlock, nil
}

func (s *HeroDataService) Put(hero *entity.Hero) error {
	unlock, err := s.lockService.Put(hero.Id(), hero)
	if err != nil {
		return err
	}
	unlock()
	return nil
}

func (s *HeroDataService) Create(hero *entity.Hero) error {
	unlock, err := s.lockService.Create(hero.Id(), hero)
	if err != nil {
		return err
	}
	unlock()

	return nil
}

func (s *HeroDataService) Func(id int64, f herolock.Func) {
	sharedHero, unlock, err := s.lock(id)
	if err != nil {
		logrus.WithField("stack", string(debug.Stack())).WithError(err).Errorf("HeroDataService.Func() lockSharedHero error, %v", id)
		f(nil, err)
		return
	}

	var newSnapshot *snapshotdata.HeroSnapshot
	func() {
		defer unlock()
		if f(sharedHero, nil) {
			// 英雄数据有改变，生成新的snapshot
			newSnapshot = s.heroSnapshotService.NewSnapshot(sharedHero)
		}
	}()

	if newSnapshot != nil {
		s.heroSnapshotService.Update(newSnapshot)
	}
}

func (s *HeroDataService) FuncNotError(id int64, f herolock.FuncNotError) (hasError bool) {
	sharedHero, unlock, err := s.lock(id)
	if err != nil {
		logrus.WithError(err).Errorf("HeroDataService.FuncNotError() lockSharedHero error, %v", id)
		return true
	}

	var newSnapshot *snapshotdata.HeroSnapshot
	func() {
		defer unlock()
		if f(sharedHero) {
			// 英雄数据有改变，生成新的snapshot
			newSnapshot = s.heroSnapshotService.NewSnapshot(sharedHero)
		}
	}()

	if newSnapshot != nil {
		s.heroSnapshotService.Update(newSnapshot)
	}

	return false
}

func (s *HeroDataService) sendFunc(id int64, f herolock.SendFunc, msgSender *MsgSender) error {
	sharedHero, unlock, err := s.lock(id)
	if err != nil {
		logrus.WithError(err).Errorf("HeroDataService.SendFunc() lockSharedHero error, %v", id)
		return err
	}

	var newSnapshot *snapshotdata.HeroSnapshot
	func() {
		defer unlock()
		f(sharedHero, msgSender)
		if msgSender.changed {
			// 英雄数据有改变，生成新的snapshot
			newSnapshot = s.heroSnapshotService.NewSnapshot(sharedHero)
		}
	}()

	if newSnapshot != nil {
		s.heroSnapshotService.Update(newSnapshot)
	}

	return nil
}

func (s *HeroDataService) FuncWithSend(id int64, f herolock.SendFunc) (hasError bool) {
	hasError, _ = s.FuncWithSendError(id, f)
	return
}

func (s *HeroDataService) funcWithSend(id int64, f herolock.SendFunc, sender sender.ClosableSender) (hasError bool) {
	hasError, _ = s.funcWithSendError(id, f, sender)
	return
}

func (s *HeroDataService) FuncWithSendError(id int64, f herolock.SendFunc) (bool, error) {
	return s.funcWithSendError(id, f, s.world.GetUserCloseSender(id))
}

func (s *HeroDataService) funcWithSendError(id int64, f herolock.SendFunc, sender sender.ClosableSender) (bool, error) {

	msgSender := getMsgSender()
	defer msgSender.clear()
	if err := s.sendFunc(id, f, msgSender); err != nil {
		// 出错，直接T下线
		if sender != nil {
			errMsg := misc.ErrDisconectReasonFailLockFail
			sender.Disconnect(errMsg)
		}
		return true, err
	} else {
		if sender != nil {
			msgSender.send(sender)
		}
		// 顺便发个系统广播
		msgSender.sendBroadcast(s.world)
		// 执行 tlog
		return !msgSender.ok, nil
	}
}

func (s *HeroDataService) Close() {
	s.lockService.Close()
}

// --- provider ---
type lockProvider struct {
	db iface.DbService
}

func (p lockProvider) Name() string {
	return "HeroData"
}

func (p *lockProvider) GetObject(ctx context.Context, id int64) (lock.LockObject, error) {
	hero, err := p.db.LoadHero(ctx, id)
	if err != nil {
		return nil, err
	}

	if hero == nil {
		return nil, lock.ErrEmpty
	}

	return hero, nil
}

func (p *lockProvider) SaveObject(ctx context.Context, id int64, obj lock.LockObject) error {
	// db save
	return p.db.SaveHero(ctx, obj.(*entity.Hero))
}

func (p *lockProvider) CreateObject(ctx context.Context, obj lock.LockObject) error {
	// db create
	ok, err := p.db.CreateHero(ctx, obj.(*entity.Hero))
	if err != nil {
		return errors.Wrapf(err, "HeroData.CreateObject")
	}

	if !ok {
		return lock.ErrCreateExist
	}

	return nil
}
