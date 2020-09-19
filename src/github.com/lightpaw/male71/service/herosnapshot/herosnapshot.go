package herosnapshot

import (
	"context"
	"github.com/lightpaw/golang-lru"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/event"
	"sync/atomic"
	"runtime/debug"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/util/imath"
)

const (
	offline_cache_count = 5000
)

var (
	notExist = &snapshotdata.HeroSnapshot{}
)

//gogen:iface
type HeroSnapshotService struct {
	db                   iface.DbService
	guildSnapshotService iface.GuildSnapshotService
	baiZhanService       iface.BaiZhanService
	datas                iface.ConfigDatas
	idCounter            uint64

	online        *snapshotmap
	offline       *lru.Cache
	callbacks     []snapshotdata.SnapshotCallback
	callbackQueue *event.FuncQueue
}

func NewHeroSnapshotService(db iface.DbService, datas iface.ConfigDatas, guildSnapshotService iface.GuildSnapshotService, baiZhanService iface.BaiZhanService, serverConfig *kv.IndividualServerConfig, ) *HeroSnapshotService {
	count := imath.Max(serverConfig.HeroOfflineCacheCount, offline_cache_count)
	return newHeroSnapshotService(db, datas, guildSnapshotService, baiZhanService, count)
}

// for test
func newHeroSnapshotService(db iface.DbService, datas iface.ConfigDatas, guildSnapshotService iface.GuildSnapshotService, baiZhanService iface.BaiZhanService, offlineCacheCount int) *HeroSnapshotService {
	cache, err := lru.New(offlineCacheCount)
	if err != nil {
		logrus.WithError(err).Panic("new cache error")
	}

	result := &HeroSnapshotService{
		db:                   db,
		datas:                datas,
		guildSnapshotService: guildSnapshotService,
		baiZhanService:       baiZhanService,
		online:               Newsnapshotmap(),
		offline:              cache,
		callbackQueue:        event.NewFuncQueue(1024, "HeroSnapshotService.loop()"),
	}

	return result
}

// 注册监听snapshot变化callback
func (s *HeroSnapshotService) RegisterCallback(callback snapshotdata.SnapshotCallback) {
	s.callbacks = append(s.callbacks, callback)
}

// 改变了英雄数据后, 缓存英雄snapshot, 必须是在unlock后调用
// 如果snapshot的版本号低于缓存中的版本号, 不会触发callback
func (s *HeroSnapshotService) Update(hero *snapshotdata.HeroSnapshot) {
	if currentValue, has := s.online.SetIfPresent(hero.Id, hero); !has {
		// 不online才加入offline. 只要在online中, 不管version如何, 都不会加入offline.
		if current := s.offline.Add(hero.Id, hero); current != hero { // 里面会比较version
			// 当前的版本更低, 不触发callback
			return
		}
	} else {
		if currentValue != hero {
			// 当前版本更低, 不触发callback
			return
		}
	}

	s.callbackQueue.TryFunc(func() {
		for _, callback := range s.callbacks {
			safeCallback(callback, hero)
		}
	})
}

func safeCallback(callback snapshotdata.SnapshotCallback, hero *snapshotdata.HeroSnapshot) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("HeroSnapshot.Callback recovered from panic. SEVERE!!!")
			metrics.IncPanic()
		}
	}()

	callback.OnHeroSnapshotUpdate(hero)
}

// 英雄上线时调用, 保存snapshot. 这个snapshot必须是没有变化的数据的, 只是上个线而已. 不会触发callback
// 如果英雄上线导致snapshot中缓存的数据有了变化, 必须再调用一次Update, 把这个变化告知其他系统
func (s *HeroSnapshotService) Online(hero *snapshotdata.HeroSnapshot) {
	s.online.Set(hero.Id, hero)
	s.offline.Remove(hero.Id) // 不必须, 但是方便gc, 并且释放缓存名额
}

// 英雄下线时调用, 把snapshot移动到lru中
func (s *HeroSnapshotService) Offline(id int64) {
	if result, has := s.online.Get(id); has {
		s.offline.Add(id, result) // 里面也会比较version, 应该online的版本一定大于offline里的版本
		s.online.RemoveIfSame(id, result)
	} else {
		logrus.Error("offline时, heroSnapshotService中竟然没有这个人的snapshot")
	}
}

// 获得英雄的snapshot, 并不保证是最新的. 尽量
// 就算英雄要从db中加载, 也不会触发callback.
// 返回nil也可能是数据库报错, 英雄未必不存在
func (s *HeroSnapshotService) Get(id int64) *snapshotdata.HeroSnapshot {
	if result, has := s.online.Get(id); has {
		return result
	}

	if result, has := s.offline.Get(id); has {
		if result == notExist {
			return nil
		}
		return result.(*snapshotdata.HeroSnapshot)
	}

	// not cached, get from db
	var hero *entity.Hero
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		hero, err = s.db.LoadHero(ctx, id)
		return
	})
	if err != nil {
		logrus.WithError(err).Error("herosnapshot load hero error")
		return nil
	} else {
		if hero == nil {
			// 英雄不存在
			s.offline.Add(id, notExist)
			return nil
		}

		snapshot := s.NewSnapshot(hero)
		s.offline.Add(id, snapshot)
		return snapshot
	}
}

// 只从Cache中获取，不读取DB
func (s *HeroSnapshotService) GetFromCache(id int64) *snapshotdata.HeroSnapshot {
	return s.getFromCache(id)
}

// 如果缓存中有, 则返回. 没有, 则返回nil
// for test only
func (s *HeroSnapshotService) getFromCache(id int64) *snapshotdata.HeroSnapshot {
	if result, has := s.online.Get(id); has {
		return result
	}

	if result, has := s.offline.Get(id); has {
		if result == notExist {
			return nil
		}
		return result.(*snapshotdata.HeroSnapshot)
	}

	return nil
}

func (s *HeroSnapshotService) GetBasicSnapshotProto(id int64) *shared_proto.HeroBasicSnapshotProto {
	snapshot := s.Get(id)
	if snapshot != nil {
		return snapshot.EncodeClient()
	}
	return nil
}

func (s *HeroSnapshotService) GetBasicProto(id int64) *shared_proto.HeroBasicProto {
	snapshot := s.Get(id)
	if snapshot != nil {
		return snapshot.EncodeBasic4Client()
	}
	return nil
}

func (s *HeroSnapshotService) GetFlagHeroName(id int64) string {
	hero := s.Get(id)
	if hero == nil {
		return idbytes.PlayerName(id)
	}
	return s.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(hero.GuildFlagName(), hero.Name)
}

func (s *HeroSnapshotService) GetHeroName(id int64) string {
	hero := s.Get(id)
	if hero == nil {
		return idbytes.PlayerName(id)
	}
	return hero.Name
}

func (s *HeroSnapshotService) GetTlogHero(id int64) entity.TlogHero {
	if hero := s.Get(id); hero != nil {
		return hero.GetTlogHero()
	}
	return nil
}

// 创建个新的snapshot, 但是并没有保存. 等unlock后再调用Cache保存
// 必须是lock住Hero的情况下才能调用, 确保此时没有其他人能访问hero对象
func (s *HeroSnapshotService) NewSnapshot(hero *entity.Hero) *snapshotdata.HeroSnapshot {
	return snapshotdata.NewSnapshot(atomic.AddUint64(&s.idCounter, 1), hero, s.guildSnapshotService.GetSnapshot, s.baiZhanService.GetJunXianLevel)
}
