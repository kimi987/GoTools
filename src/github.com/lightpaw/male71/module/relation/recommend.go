package relation

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/sortkeys"
	"github.com/lightpaw/male7/util/u64"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	cacheCountLimit  = 10000
	countPerLocation = 200
)

// 初始化加载一批数据

// 定时根据lastOnlineTime进行排序

// 新用户加入，添加进去，更新一下（这部分数据使用ReadOnly）

// 玩家更新Location时候，同步更新

func NewLocationHeroCache(time iface.TimeService, datas iface.ConfigDatas) *LocationHeroCache {
	countPerLoc := datas.MiscConfig().RefreshRecommendHeroPageSize *
		datas.MiscConfig().RefreshRecommendHeroPageCount
	minLevel := datas.MiscConfig().RefreshRecommendHeroMinLevel
	expireSeconds := int64(datas.MiscConfig().RecommendHeroOfflineExpireDuration.Seconds())
	return newLocationHeroCache(time, countPerLoc, minLevel, expireSeconds)
}

func newLocationHeroCache(time iface.TimeService, countPerLocation uint64, minLevel uint64, expireSeconds int64) *LocationHeroCache {
	c := &LocationHeroCache{
		time:             time,
		minHeroLevel:     minLevel,
		countPerLocation: u64.Min(countPerLocation, countPerLocation),
		expireSeconds:    expireSeconds,
		heroCount:        atomic.NewUint64(0),
		isAddHero:        atomic.NewBool(false),
	}

	go call.CatchLoopPanic(c.loop, "LocationHeroCache.loop()")

	return c
}

//gogen:iface
type LocationHeroCache struct {
	time iface.TimeService

	minHeroLevel     uint64
	countPerLocation uint64
	expireSeconds    int64

	// 存储英雄数据
	heroMap sync.Map

	heroCount *atomic.Uint64
	isAddHero *atomic.Bool

	// 定时遍历英雄数据，构建数组
	locMap sync.Map

	locLoadMap sync.Map

	// 最小的英雄数据过期时间
	minHeroUpdateTime int64
}

func (c *LocationHeroCache) isLoaded(location uint64) bool {
	if v, _ := c.locLoadMap.Load(location); v != nil {
		return v.(bool)
	}
	return false
}

func (c *LocationHeroCache) setLoaded(location uint64) {
	c.locLoadMap.Store(location, true)
}

func (c *LocationHeroCache) getHero(id int64) *HeroObject {
	if o, _ := c.heroMap.Load(id); o != nil {
		return o.(*HeroObject)
	}
	return nil
}

func (c *LocationHeroCache) UpdateHero(hero *snapshotdata.HeroSnapshot, updateTime int64) {

	if hero == nil {
		// 防御性
		return
	}

	if hero.Level < c.minHeroLevel {
		// 等级太低
		return
	}

	heroId := hero.Id
	if heroObj := c.getHero(heroId); heroObj != nil {
		heroObj.location.Store(hero.Location)
		heroObj.updateTime.Store(updateTime)
	} else {
		// 如果没有达到缓存上限，添加进来
		if c.heroCount.Load() < cacheCountLimit {
			obj := newHeroObject(heroId, hero.EncodeClient(), updateTime)
			if actual, loaded := c.heroMap.LoadOrStore(heroId, obj); loaded {
				heroObj := actual.(*HeroObject)

				heroObj.location.Store(hero.Location)
				heroObj.updateTime.Store(updateTime)
			} else {
				c.heroCount.Inc() // 不是很准，但是没关系
				c.isAddHero.Store(true)
			}
		}
	}

}

func (c *LocationHeroCache) updateLocationHeros(location uint64, heroProtos []*shared_proto.HeroBasicSnapshotProto) []*HeroObject {

	var array []*HeroObject
	for _, heroProto := range heroProtos {
		if heroProto == nil {
			continue
		}

		heroId, ok := idbytes.ToId(heroProto.Basic.Id)
		if !ok {
			continue
		}

		heroObj := c.getHero(heroId)
		if heroObj == nil {
			obj := newHeroObject(heroId, heroProto, int64(heroProto.LastOfflineTime))

			if actual, loaded := c.heroMap.LoadOrStore(heroId, obj); loaded {
				heroObj = actual.(*HeroObject)
			} else {
				heroObj = obj
				c.heroCount.Inc() // 不是很准，但是没关系
			}
		}

		array = append(array, heroObj)
	}

	c.locMap.Store(location, array)
	c.setLoaded(location)

	return array
}

func (c *LocationHeroCache) UpdateLocation(heroId int64, location uint64) {
	if heroObj := c.getHero(heroId); heroObj != nil {
		heroObj.location.Store(location)
	}
}

func (c *LocationHeroCache) loop() {

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			call.CatchPanic(c.update, "LocationHeroCache.update()")
		}
	}

}

func (c *LocationHeroCache) walkHero(f func(heroId int64, obj *HeroObject) bool) {
	c.heroMap.Range(func(key, value interface{}) bool {
		heroId := key.(int64)
		obj := value.(*HeroObject)
		return f(heroId, obj)
	})
}

func (c *LocationHeroCache) removeHero(heroId int64) {
	c.heroMap.Delete(heroId)
}

func (c *LocationHeroCache) update() {

	ctime := c.time.CurrentTime()

	ctimeUnix := ctime.Unix()
	expireTime := ctimeUnix - c.expireSeconds
	if !c.isAddHero.Load() && expireTime <= c.minHeroUpdateTime {
		// 还没有过期
		return
	}
	c.isAddHero.Store(false)

	// 遍历所有的影响数据

	count := 0
	locMap := make(map[uint64][]*sortkeys.I64KV)

	var minHeroUpdateTime int64 = math.MaxInt64
	c.walkHero(func(heroId int64, obj *HeroObject) (toContinue bool) {
		toContinue = true

		heroUpdateTime := obj.updateTime.Load()
		if heroUpdateTime < expireTime {
			c.removeHero(heroId)
			return
		}

		kv := sortkeys.NewI64KV(obj.updateTime.Load(), obj)

		if location := obj.location.Load(); location != 0 {
			locMap[location] = append(locMap[location], kv)
		}
		locMap[0] = append(locMap[0], kv)

		count++
		minHeroUpdateTime = i64.Min(minHeroUpdateTime, heroUpdateTime)
		return
	})
	c.minHeroUpdateTime = minHeroUpdateTime

	// 按倒叙排序，将超出长度的列表*2移除掉（防止某个地区人数过多，导致某个地区不存玩家的bug（因为总数固定））
	for loc, kvs := range locMap {
		sort.Sort(sort.Reverse(sortkeys.I64KVSlice(kvs)))

		countPerLocation := u64.Int(c.countPerLocation) * 2
		if len(kvs) > countPerLocation {
			// 超出限制，删除多余的
			for i := countPerLocation; i < len(kvs); i++ {
				heroObj := kvs[i].V.(*HeroObject)
				c.removeHero(heroObj.heroId)

				count--
			}
			locMap[loc] = kvs[:countPerLocation]
		}
	}

	// 先全部put进去，不存在的，remove掉
	for k, kvs := range locMap {
		var heroArray []*HeroObject
		for _, kv := range kvs {
			obj := kv.V.(*HeroObject)
			heroArray = append(heroArray, obj)
		}

		c.locMap.Store(k, heroArray)
	}

	c.locMap.Range(func(key, value interface{}) bool {
		location := key.(uint64)
		if _, exist := locMap[location]; !exist {
			c.locMap.Delete(key)
		}

		return true
	})

	// 更新一下英雄数量
	c.heroCount.Store(u64.FromInt(count))

}

func (c *LocationHeroCache) getHerosByLocation(location uint64) []*HeroObject {
	if array, _ := c.locMap.Load(location); array != nil {
		return array.([]*HeroObject)
	}
	return nil
}

func (c *LocationHeroCache) rangeHeros(location uint64, startIndex int, f func(i int, o *HeroObject) bool) {

	array := c.getHerosByLocation(location)

	n := len(array)
	for i := 0; i < n; i++ {

		idx := (startIndex + i) % n

		o := array[idx]
		if !f(i, o) {
			break
		}
	}

}

func newHeroObject(heroId int64, heroProto *shared_proto.HeroBasicSnapshotProto, updateTime int64) *HeroObject {
	return &HeroObject{
		heroId:     heroId,
		location:   atomic.NewUint64(u64.FromInt32(heroProto.Basic.Location)),
		updateTime: atomic.NewInt64(updateTime),
		heroProto:  heroProto,
	}
}

type HeroObject struct {
	heroId int64

	location *atomic.Uint64

	updateTime *atomic.Int64

	heroProto *shared_proto.HeroBasicSnapshotProto
}
