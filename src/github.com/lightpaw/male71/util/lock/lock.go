package lock

import (
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/pkg/errors"
	"runtime/debug"
	"sync"
	"time"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util"
)

const (
	mustSaveInterval     = 5 * time.Minute  // 定时保存间隔
	evictNotAccessedTime = 10 * time.Minute // 多久没访问就删除
)

var (
	ErrInvalidated = errors.New("entry invalidated")
	ErrEmpty       = errors.New("entry empty")
	ErrCreateExist = errors.New("create exist entry")
	errPanic       = errors.New("provider panic")
)

type LockService struct {
	provider   LockProvider
	entries    *entrymap
	shouldQuit chan struct{}

	evictNotAccessedDuration time.Duration
}

func NewLockService(provider LockProvider, evictNotAccessedDuration time.Duration) *LockService {
	result := &LockService{
		provider:                 provider,
		entries:                  Newentrymap(),
		shouldQuit:               make(chan struct{}),
		evictNotAccessedDuration: evictNotAccessedDuration,
	}

	go result.saveLoop()
	go result.evictLoop()
	go result.checkLoop()
	return result
}

// 提供加载及保存数据的方法
type LockProvider interface {
	Name() string                                         // Service的名字
	GetObject(context.Context, int64) (LockObject, error) // 获取某id的数据, 一般从数据库中读取并解析
	SaveObject(context.Context, int64, LockObject) error  // 保存某id的数据, 一般调用LockObject.Marshal方法, 并把数据保存在数据库中
	CreateObject(context.Context, LockObject) error       // 存入db中
}

// 被LockService保护的数据
type LockObject interface {
	Marshal() ([]byte, error)
}

// 解锁的方法
type Unlocker func()

// 所有玩家都已确保下线, 并且所有会用到这里东西的功能模块都已正常退出之后再调用
func (s *LockService) Close() {
	close(s.shouldQuit)

	logrus.Infof("LockService.%s关闭保存中", s.getProviderName())

	// 关闭时候，尝试保存3次
	for i := 0; i < 3; i++ {
		s.doEvictLoop(0, 60*time.Second) // evict everything
		if count := s.entries.Count(); count != 0 {
			logrus.WithField("count", count).Errorf("LockService.%s关闭删除了之后里面竟然还有对象... 应该其他会用到LockService的其他功能模块都完全退出之后再调用Close", s.getProviderName())
		} else {
			break
		}
	}

	logrus.Infof("LockService.%s已关闭", s.getProviderName())
}

// 必须在hold住对象的lock时调用, 把对象标记为已删除
func (s *LockService) Invalidate(id int64) {
	entry, has := s.entries.Get(id)
	if !has {
		logrus.Errorf("LockServer.%s竟然要invalidate一个不存在的entry. 必须要在Lock住对象的时候调用啊")
		return
	}

	entry.invalidated = true
}

func (s *LockService) Create(id int64, data LockObject) (Unlocker, error) {
	return s.create(id, data, true)
}

func (s *LockService) Put(id int64, data LockObject) (Unlocker, error) {
	return s.create(id, data, false)
}

func (s *LockService) create(id int64, data LockObject, insert bool) (Unlocker, error) {
	entry, has := s.entries.Get(id)
	if has {
		return nil, ErrCreateExist
	}
	// 先创建lockEntry, lock, 加入map, 在db中创建, 再unlock
	entry = &lockEntry{}
	entry.Lock()
	_, ok := s.entries.SetIfAbsent(id, entry)
	if !ok {
		// 已经存在旧的, 这么巧?
		entry.Unlock() // 废弃
		return nil, ErrCreateExist
	}

	// put ok. insert into db
	if insert {
		// load 一次，看下有没有这个数据
		if _, err := s.loadObject(id); err != nil {
			if err != ErrEmpty {
				if ok := s.entries.RemoveIfSame(id, entry); !ok {
					logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
				}
				entry.destroyed = true
				entry.Unlock()
				return nil, err
			}
		} else {
			if ok := s.entries.RemoveIfSame(id, entry); !ok {
				logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()
			return nil, ErrCreateExist
		}

		err := s.createObject(data)
		if err != nil {
			if ok := s.entries.RemoveIfSame(id, entry); !ok {
				logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()

			return nil, errors.Wrapf(err, "lockService.%s.createObject(%d)失败", s.getProviderName(), id)
		}
	}

	entry.data = data
	now := time.Now()
	entry.lastAccessTime, entry.lastSaveTime = now, now
	return entry.Mutex.Unlock, nil
}

// LockService关闭后还可以调用, 并没有限制. 但是修改的内容并不会保存到数据库中
// 用完之后一定要调用LockData.Unlock(), 并且在那之后不要再访问里面的内容(read也不行)
// 不unlock的话所有线程统统卡死, 也不会保存, 也不能关闭, 所有请求同一个对象的线程也卡死
func (s *LockService) Lock(id int64) (LockObject, Unlocker, error) {
	entry, has := s.entries.Get(id)
	if has {
		entry.Lock()
		if entry.destroyed {
			entry.Unlock()
			return s.Lock(id) // recursive
		}

		if entry.invalidated {
			entry.Unlock()
			return nil, nil, ErrInvalidated
		}
		entry.lastAccessTime = time.Now()
		// 使用
		return entry.data, entry.Mutex.Unlock, nil
	}
	// 先创建lockEntry, lock, 加入map, 去db取, 使用, 再unlock
	entry = &lockEntry{}
	entry.Lock()
	old, ok := s.entries.SetIfAbsent(id, entry)
	if ok {
		// put ok. get from database and return
		data, err := s.loadObject(id)
		if err != nil {
			if ok := s.entries.RemoveIfSame(id, entry); !ok {
				logrus.Errorf("LockService.%s.Lock要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()

			if err != ErrEmpty {
				err = errors.Wrapf(err, "lockService.%s.GetObject(%d)失败", s.getProviderName(), id)
			}

			return nil, nil, err
		}

		entry.data = data

		now := time.Now()
		entry.lastAccessTime, entry.lastSaveTime = now, now
		return data, entry.Mutex.Unlock, nil
	}
	// 已经存在旧的, 这么巧?
	entry.Unlock() // 废弃
	old.Lock()
	if old.destroyed {
		// 这么巧? 刚拿出来的时候没有, 要放进去的时候已经有了, 再lock又已经destroy了?
		old.Unlock()
		return s.Lock(id) // recursive
	}
	if old.invalidated {
		old.Unlock()
		return nil, nil, ErrInvalidated
	}
	old.lastAccessTime = time.Now()
	return old.data, old.Mutex.Unlock, nil
}

// 定时删除很久没有人访问的
func (s *LockService) evictLoop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.evictLoop recovered from panic. SEVERE!!!", s.getProviderName())
			metrics.IncPanic()
		}

		logrus.Debugf("LockService.%s.evictLoop exited", s.getProviderName())
	}()

	evictDuration := timeutil.MaxDuration(s.evictNotAccessedDuration, 10*time.Second)

	select {
	case <-time.After(1 * time.Minute):
		// 和保存loop错开运行
	case <-s.shouldQuit:
		return
	}

	tick := time.Tick(timeutil.MinDuration(evictDuration/2, 2*time.Minute))

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doEvictLoop(evictDuration, 3*time.Second)
		}
	}
}

func (s *LockService) doEvictLoop(evictTime, saveTimeout time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.doEvictLoop recovered from panic. SEVERE!!!", s.getProviderName())
			metrics.IncPanic()
		}
	}()

	logrus.Debugf("LockService.%s.evictLoop扫描中", s.getProviderName())
	entries := s.entries.Iter()
	for ch := range entries {
		for _, item := range ch {
			entry := item.Val
			entry.Lock()
			if entry.destroyed {
				entry.Unlock()
				continue
			}

			if entry.invalidated {
				if ok := s.entries.RemoveIfSame(item.Key, entry); !ok {
					logrus.Errorf("LockService.%s.remove要删除的竟然不是在map中的", s.getProviderName())
				} else {
					// 已删除
					entry.destroyed = true
				}
			} else if time.Since(entry.lastAccessTime) >= evictTime || evictTime == 0 {
				// 删除时, 先lock, 检查是否destroy, 保存, 设置已destroy, unlock
				if err := s.saveObject(item.Key, entry.data, saveTimeout); err != nil {
					logrus.WithError(err).WithField("id", item.Key).Errorf("LockService.%s要删除LockData, 保存时出错", s.getProviderName())
				} else {
					// 已保存
					if ok := s.entries.RemoveIfSame(item.Key, entry); !ok {
						logrus.Errorf("LockService.%s.remove要删除的竟然不是在map中的", s.getProviderName())
					} else {
						// 已删除
						entry.destroyed = true
					}
				}
			}
			entry.Unlock()
		}
	}
	logrus.Debugf("LockService.%s.evictLoop扫描完成", s.getProviderName())
}

// 定时保存
func (s *LockService) saveLoop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.saveLoop recovered from panic. SEVERE!!!", s.getProviderName())
			metrics.IncPanic()
		}

		logrus.Debugf("LockService.%s.saveLoop exited", s.getProviderName())
	}()

	tick := time.Tick(2 * time.Minute)

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doSaveLoop(mustSaveInterval)
		}
	}
}

func (s *LockService) doSaveLoop(saveInterval time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.doSaveLoop recovered from panic. SEVERE!!!", s.getProviderName())
			metrics.IncPanic()
		}
	}()
	// 不直接取这里的时间. 可能db很慢很慢, 一个循环花很久, 导致不停在保存, 更增加了db压力
	// 全量扫描保存
	logrus.Debugf("LockService.%s保存检查中", s.getProviderName())
	entries := s.entries.Iter()
	for ch := range entries {
		for _, item := range ch {
			entry := item.Val
			entry.Lock()
			if entry.destroyed || entry.invalidated {
				entry.Unlock()
				continue
			}
			if entry.lastAccessTime.After(entry.lastSaveTime.Add(-time.Second)) && time.Since(entry.lastSaveTime) >= saveInterval {
				// 必须是上次被人接触过的时间 > (上次保存的时间 - 1秒) && 超出保存间隔
				// 上次保存之后的1秒内有用过, 也重新保存
				if err := s.saveObject(item.Key, entry.data, 3*time.Second); err != nil {
					logrus.WithError(err).WithField("id", item.Key).Errorf("LockService.%s定时保存出错", s.getProviderName())
				} else {
					// 已保存
					entry.lastSaveTime = time.Now()
				}
			}
			entry.Unlock()
		}
	}
	logrus.Debugf("LockService.%s保存检查完成", s.getProviderName())
}

func (s *LockService) getProviderName() string {
	return s.provider.Name()
}

func (s *LockService) createObject(o LockObject) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.createObject recovered from panic. SEVERE!!!", s.getProviderName())

			err = errPanic
			metrics.IncPanic()
		}
	}()

	return ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		return s.provider.CreateObject(ctx, o)
	})
}

func (s *LockService) loadObject(key int64) (o LockObject, err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.loadObject recovered from panic. SEVERE!!!", s.getProviderName())

			err = errPanic
			metrics.IncPanic()
		}
	}()

	err = ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		o, err = s.provider.GetObject(ctx, key)
		return err
	})
	return
}

func (s *LockService) saveObject(key int64, data LockObject, timeout time.Duration) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.saveObject recovered from panic. SEVERE!!!", s.getProviderName())
			err = errPanic
			metrics.IncPanic()
		}
	}()

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return ctxfunc.Timeout(timeout, func(ctx context.Context) (err error) {
		return s.provider.SaveObject(ctx, key, data)
	})
}

// 检查所有的entry是否有长时间没有unlock的
func (s *LockService) checkLoop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("LockService.%s.checkLoop recovered from panic. SEVERE!!!", s.getProviderName())
			metrics.IncPanic()
		}
	}()

	<-time.After(30 * time.Second) // 和saveLoop & evictLoop 错开
	tick := time.Tick(2 * time.Minute)

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doCheckLoop()
		}
	}
}

func (s *LockService) doCheckLoop() {
	shouldQuit := make(chan struct{})

	var checkEntry *lockEntry
	entries := s.entries.Iter()
	timeout := time.After(5 * time.Second)

	go call.CatchPanic(func() {
		select {
		case <-shouldQuit:
			return
		case <-timeout:
		}
		if checkEntry := checkEntry; checkEntry != nil {
			for i := 0; i < 100; i++ {
				logrus.WithField("object", checkEntry).Errorf("LockService.%s 检测到长时间没有unlock的对象!!!! SEVERE!!!", s.getProviderName())
			}
			util.DumpStacks(s.getProviderName())
		}
	}, "LockService.doCheckLoop")
	for ch := range entries {
		for _, item := range ch {
			checkEntry = item.Val
			checkEntry.Lock()
			checkEntry.Unlock()
		}
	}
	close(shouldQuit)
}

type lockEntry struct {
	sync.Mutex
	lastAccessTime time.Time  // 上次有人请求的时间
	lastSaveTime   time.Time  // 上次自动保存的时间
	destroyed      bool       // 是否已废弃失效
	invalidated    bool       // 是否已被标记为从db中删除
	data           LockObject // 真正的内容
}

func (l lockEntry) String() string {
	return fmt.Sprintf("%v", l.data)
}
