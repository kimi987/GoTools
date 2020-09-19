package entity

import (
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/server_proto"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/u64"
)

type randomEvent struct {
	id       int          // 事件ID
	cube     cb.Cube
}

func (e *randomEvent) Id() uint64 {
	return i32.Uint64(int32(e.id))
}

func (e *randomEvent) SetId(id uint64) {
	e.id = u64.Int(id)
}

func (e *randomEvent) Cube() cb.Cube {
	return e.cube
}

type hero_random_event struct {
	eventMap          map[cb.Cube] *randomEvent // <坐标, 事件>
	bigRefreshTime    time.Time // 大刷新到期时间
	smallRefreshTime  time.Time // 小刷新到期时间
	// 事件图鉴（点开事件就激活）
	handbooks map[uint64]struct{}
}

func newHeroRandomEvent(ctime time.Time) *hero_random_event {
	return &hero_random_event {
		eventMap: make(map[cb.Cube] *randomEvent),
		bigRefreshTime: ctime.Add(-time.Minute),
		smallRefreshTime: ctime,
		handbooks: make(map[uint64]struct{}),
	}
}

func (e *hero_random_event) EventNum() int {
	return len(e.eventMap)
}

func (e *hero_random_event) CheckBigRefreshTime(t time.Time) bool {
	return t.After(e.bigRefreshTime)
}

func (e *hero_random_event) SetBigRefreshTime(t time.Time) {
	e.bigRefreshTime = t
}

func (e *hero_random_event) CheckSmallRefreshTime(t time.Time) bool {
	return t.After(e.smallRefreshTime)
}

func (e *hero_random_event) SetSmallRefreshTime(t time.Time) {
	e.smallRefreshTime = t
}

func (e *hero_random_event) unmarshal(proto *server_proto.HeroRandomEventServerProto) {

	if proto == nil {
		return
	}
	for _, event := range proto.Events {
		cb := cb.XYCubeI32(event.PosX, event.PosY)
		e.eventMap[cb] = &randomEvent {
			id: int(event.Id),
			cube: cb,
		}
	}
	e.bigRefreshTime = timeutil.Unix64(proto.BigRefreshTime)
	e.smallRefreshTime = timeutil.Unix64(proto.SmallRefreshTime)
	for _, eventId := range proto.Handbooks {
		e.handbooks[eventId] = struct{}{}
	}
}

func (e *hero_random_event) encode() *server_proto.HeroRandomEventServerProto {

	proto := &server_proto.HeroRandomEventServerProto{}
	for cb, event := range e.eventMap {
		x, y := cb.XYI32()
		p := &server_proto.EventPositionServerProto {
			Id: int32(event.id),
			PosX: x,
			PosY: y,
		}
		proto.Events = append(proto.Events, p)
	}
	proto.BigRefreshTime = timeutil.Marshal64(e.bigRefreshTime)
	proto.SmallRefreshTime = timeutil.Marshal64(e.smallRefreshTime)
	for eventId, _ := range e.handbooks {
		proto.Handbooks = append(proto.Handbooks, eventId)
	}
	return proto
}

func (e *hero_random_event) encodeClient() *shared_proto.HeroRandomEventProto {

	proto := &shared_proto.HeroRandomEventProto{}
	for cb, _ := range e.eventMap {
		x, y := cb.XYI32()
		p := &shared_proto.EventPositionProto{x, y}
		proto.Events = append(proto.Events, p)
	}
	for eventId, _ := range e.handbooks {
		proto.Handbooks = append(proto.Handbooks, u64.Int32(eventId))
	}
	return proto
}

func (e *hero_random_event) TrySetHandbooks(id uint64) bool {
	if _, ok := e.handbooks[id]; ok {
		return false
	}
	e.handbooks[id] = struct{}{}
	return true
}

func (e *hero_random_event) AddEvent(id, x, y int, endTime time.Time) {

	cb := cb.XYCube(x, y)
	if _, ok := e.eventMap[cb]; !ok {
		e.eventMap[cb] = &randomEvent {
			id: id,
			cube: cb,
		}
	}
}

func (e *hero_random_event) GetEvent(x, y int32) *randomEvent {

	cb := cb.XYCubeI32(x, y)
	if event, ok := e.eventMap[cb]; ok {
		return event
	}
	return nil
}

func (e *hero_random_event) RemoveEvent(x, y int32) {

	cb := cb.XYCubeI32(x, y)
	delete(e.eventMap, cb)
}

func (e *hero_random_event) GetAllEventsList() []cb.Cube {
	var list []cb.Cube
	for cb, _ := range e.eventMap {
		list = append(list, cb)
	}
	return list
}

func (e *hero_random_event) ClearAllEvents() {
	if len(e.eventMap) > 0 {
		e.eventMap = make(map[cb.Cube] *randomEvent)
	}
	return
}

func (e *hero_random_event) PutEvents(cubes []cb.Cube) {
	for _, cube := range cubes {
		e.eventMap[cube] = &randomEvent {
			id: 0,
			cube: cube,
		}
	}
}

func (e *hero_random_event) ClearEvents(cubes []cb.Cube) (invalidEvents []cb.Cube) {
	if len(cubes) > 0 {
		for _, cube := range cubes {
			if _, ok := e.eventMap[cube]; ok {
				delete(e.eventMap, cube)
				invalidEvents = append(invalidEvents, cube)
			}
		}
	}
	return
}
