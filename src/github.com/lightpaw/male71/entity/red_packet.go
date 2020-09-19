package entity

import (
	"encoding/json"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"sync"
	"time"
)

func CheckCanCreateRedPacket(allMoney, minMoney, count uint64) bool {
	return count > 0 && allMoney >= minMoney*count
}

func NewDefaultRedPacket() *RedPacket {
	r := &RedPacket{}
	r.grabbedHeros = make(map[int64]bool)
	return r
}

type RedPacket struct {
	id            int64
	data          *red_packet.RedPacketData
	createTime    time.Time
	createHero    *shared_proto.HeroBasicProto
	text          string
	chatType      shared_proto.ChatType
	nextGrabIndex int

	lock sync.RWMutex

	parts        []*RedPacketPart
	grabbedHeros map[int64]bool

	chatId int64 // 对应的聊天消息 id
}

func (r *RedPacket) BuildChatJson() (jsonStr string, err error) {
	j := &ChatJsonObj{}

	r.lock.RLock()
	defer r.lock.RUnlock()

	j.DataId = r.data.Id
	j.Text = r.text
	j.CreateTime = timeutil.Marshal32(r.createTime)
	j.ExipredTime = timeutil.Marshal32(r.createTime.Add(r.data.ExpiredDuration))
	j.AllGrabbed = r.allGrabbed()

	jsonStr, err = j.BuildJson()

	return
}

func (c *ChatJsonObj) BuildJson() (jsonStr string, err error) {
	var b []byte
	if b, err = json.Marshal(c); err != nil {
		return
	}
	jsonStr = string(b)
	return
}

func ChatJsonObjMarshal(jsonStr string) *ChatJsonObj {
	obj := &ChatJsonObj{}
	if err := json.Unmarshal([]byte(jsonStr), obj); err != nil {
		logrus.WithError(err).Error("marshal ChatJsonObj 异常")
		return nil
	}
	return obj
}

// json: {"text":"xxx","create_time":1545040560,"expired_time":1545040560,"data_id":1,"all_grabbed":true}
type ChatJsonObj struct {
	Text        string `json:"text"`
	CreateTime  int32  `json:"create_time"`
	ExipredTime int32  `json:"expired_time"`
	DataId      uint64 `json:"data_id"`
	AllGrabbed  bool   `json:"all_grabbed"`
}

func CreateRedPacket(id int64, data *red_packet.RedPacketData, count uint64, createTime time.Time, text string, createHero *shared_proto.HeroBasicProto, chatType shared_proto.ChatType) (r *RedPacket, succ bool) {
	if !CheckCanCreateRedPacket(data.Money, data.MinPartMoney, count) {
		return
	}

	succ = true

	r = NewDefaultRedPacket()
	r.id = id
	r.data = data
	r.createTime = createTime
	r.text = text
	r.createHero = createHero
	r.chatType = chatType

	r.parts = createParts(data.Money, data.MinPartMoney, count)

	return
}

// len(parts) == count，如果总钱数不够，部分 part.money == 0
func createParts(all, partMin, count uint64) (parts []*RedPacketPart) {
	currAll := all

	for i := 0; i < u64.Int(count); i++ {
		baseMoney := u64.Min(currAll, partMin)
		currAll = u64.Sub(currAll, baseMoney)
		parts = append(parts, newRedPacketPart(baseMoney))
	}

	for i := 0; i < u64.Int(count); i++ {
		if currAll <= 0 {
			break
		}

		var partMoney uint64
		currMin := u64.Min(partMin, currAll)
		if currAll <= currMin {
			partMoney = currAll
		} else if i == u64.Int(count-1) {
			partMoney = currAll
		} else {
			currCount := u64.Sub(count, u64.FromInt(i))
			currMax := u64.Min(currAll, (currAll / currCount) * 2)
			partMoney = u64.FromInt64(rand.Int63n(int64(currMax)))
		}

		parts[i].money += partMoney
		currAll = u64.Sub(currAll, partMoney)
	}

	mix(parts)

	return
}

func mix(parts []*RedPacketPart) {
	n := len(parts)
	for i := range parts {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		parts[idx], parts[swap] = parts[swap], parts[idx]
	}
}

func (r *RedPacket) ChatId() int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.chatId
}

func (r *RedPacket) SetChatId(chatId int64) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.chatId = chatId
}

func (r *RedPacket) CanGrab(ctime time.Time) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.canGrab(ctime)
}

func (r *RedPacket) ExpiredTime() time.Time {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.createTime.Add(r.data.ExpiredDuration)
}

func (r *RedPacket) canGrab(ctime time.Time) bool {
	if ctime.After(r.createTime.Add(r.data.ExpiredDuration)) {
		return false
	}

	return !r.allGrabbed()
}

func (r *RedPacket) AllGrabbed() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.allGrabbed()
}

func (r *RedPacket) allGrabbed() bool {
	return r.nextGrabIndex >= len(r.parts)
}

func (r *RedPacket) Grab(heroId int64, hero *shared_proto.HeroBasicProto, ctime time.Time) (money uint64, allGarbbed bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if !r.canGrab(ctime) {
		allGarbbed = r.allGrabbed()
		logrus.Debugf("抢红包，红包 id:%v 过期了", r.id)
		return
	}

	if r.grabbed(heroId) {
		logrus.Debugf("抢红包，红包 id:%v hero:%v 抢过了", r.id, heroId)
		return
	}

	if part := r.parts[r.nextGrabIndex]; part != nil {
		part.grabbedTime = ctime
		part.grabbedHero = hero

		r.grabbedHeros[heroId] = true

		money = part.money
	}

	r.nextGrabIndex++
	allGarbbed = r.allGrabbed()

	return
}

func (r *RedPacket) Grabbed(heroId int64) (grabbed bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.grabbed(heroId)
}

func (r *RedPacket) grabbed(heroId int64) (grabbed bool) {
	_, grabbed = r.grabbedHeros[heroId]
	return
}

func (r *RedPacket) Id() int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.id
}

func (r *RedPacket) ChatType() shared_proto.ChatType {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.chatType
}

func (r *RedPacket) GuildId() (gid int64) {
	if r.createHero == nil {
		return
	}
	return int64(r.createHero.GuildId)
}

func (r *RedPacket) Encode() *shared_proto.RedPacketProto {
	r.lock.RLock()
	defer r.lock.RUnlock()

	p := &shared_proto.RedPacketProto{}
	p.Id = idbytes.ToBytes(r.id)
	p.DataId = u64.Int32(r.data.Id)
	p.CreateTime = timeutil.Marshal32(r.createTime)
	p.CreateHero = r.createHero
	p.Text = r.text
	p.ChatType = r.chatType
	p.Count = int32(len(r.parts))
	p.AllGrabbed = r.allGrabbed()

	for _, part := range r.parts {
		if part.grabbed() {
			p.GrabbedParts = append(p.GrabbedParts, part.encode())
		}
	}

	return p
}

func (r *RedPacket) EncodeServer() *server_proto.RedPacketServerProto {
	p := &server_proto.RedPacketServerProto{}
	p.Id = r.id
	p.DataId = r.data.Id
	p.CreateHero = r.createHero
	p.CreateTime = timeutil.Marshal64(r.createTime)
	p.Text = r.text
	p.ChatType = r.chatType
	p.ChatId = r.chatId

	for _, part := range r.parts {
		p.Parts = append(p.Parts, part.encodeServer())
	}

	p.GrabbedHero = r.grabbedHeros

	return p
}

func (r *RedPacket) Unmarshal(p *server_proto.RedPacketServerProto, datas *config.RedPacketDataConfig) {
	r.id = p.Id
	r.data = datas.Must(p.DataId)
	if r.data.Id != p.DataId {
		logrus.Errorf("RedPacket.Unmarshal.找不到 dataId: %v", p.Id)
	}

	r.createHero = p.CreateHero
	r.createTime = timeutil.Unix64(p.CreateTime)
	r.text = p.Text
	r.chatType = p.ChatType
	r.chatId = p.ChatId

	if p.Parts != nil {
		for _, p := range p.Parts {
			part := newRedPacketPart(p.Money)
			part.unmarshal(p)
			r.parts = append(r.parts, part)
			if !timeutil.IsZero(part.grabbedTime) {
				r.nextGrabIndex++
			}
		}
	}

	if p.GrabbedHero != nil {
		r.grabbedHeros = p.GrabbedHero
	}

}

func newRedPacketPart(money uint64) *RedPacketPart {
	part := &RedPacketPart{}
	part.money = money
	part.grabbedTime = time.Time{}
	return part
}

type RedPacketPart struct {
	money       uint64
	grabbedTime time.Time
	grabbedHero *shared_proto.HeroBasicProto
}

func (r *RedPacketPart) grabbed() bool {
	return !timeutil.IsZero(r.grabbedTime)
}

func (r *RedPacketPart) encode() *shared_proto.RedPacketPartProto {
	p := &shared_proto.RedPacketPartProto{}
	p.Money = u64.Int32(r.money)
	p.GrabbedTime = timeutil.Marshal32(r.grabbedTime)
	p.GrabbedHero = r.grabbedHero

	return p
}

func (r *RedPacketPart) encodeServer() *server_proto.RedPacketPartServerProto {
	p := &server_proto.RedPacketPartServerProto{}
	p.Money = r.money
	p.GrabbedTime = timeutil.Marshal64(r.grabbedTime)
	p.GrabbedHero = r.grabbedHero

	return p
}

func (r *RedPacketPart) unmarshal(p *server_proto.RedPacketPartServerProto) {
	r.money = p.Money
	r.grabbedHero = p.GrabbedHero
	r.grabbedTime = timeutil.Unix64(p.GrabbedTime)
}

func NewHeroRedPacket() *HeroRedPacket {
	h := &HeroRedPacket{}
	h.boughtRedPacketCounts = make(map[uint64]uint64)
	return h
}

type HeroRedPacket struct {
	grabbedId             []int64 // 抢过的红包
	boughtRedPacketCounts map[uint64]uint64
}

func (h *HeroRedPacket) encode() *shared_proto.HeroRedPacketProto {
	p := &shared_proto.HeroRedPacketProto{}
	for k, v := range h.boughtRedPacketCounts {
		p.RedPackets = append(p.RedPackets, &shared_proto.BoughtRedPacketProto{DataId: u64.Int32(k), Count: u64.Int32(v)})
	}

	return p
}

func (h *HeroRedPacket) encodeServer() *server_proto.HeroRedPacketServerProto {
	p := &server_proto.HeroRedPacketServerProto{}
	p.BoughtRedPacketCounts = h.boughtRedPacketCounts
	p.GrabbedId = h.grabbedId

	return p
}

func (h *HeroRedPacket) unmarshal(p *server_proto.HeroRedPacketServerProto) {
	if p == nil {
		return
	}

	if p.BoughtRedPacketCounts != nil {
		h.boughtRedPacketCounts = p.BoughtRedPacketCounts
	}

	h.grabbedId = p.GrabbedId
}

func (h *HeroRedPacket) Add(dataId uint64, count uint64) {
	h.boughtRedPacketCounts[dataId] += count
}

func (h *HeroRedPacket) IsBought(dataId uint64) bool {
	return h.boughtRedPacketCounts[dataId] > 0
}

func (h *HeroRedPacket) Reduce(dataId uint64) (succ bool) {
	c := h.boughtRedPacketCounts[dataId]
	if c <= 0 {
		return
	}

	h.boughtRedPacketCounts[dataId] -= 1

	if h.boughtRedPacketCounts[dataId] <= 0 {
		delete(h.boughtRedPacketCounts, dataId)
	}

	succ = true
	return
}

func (h *HeroRedPacket) Count(dataId uint64) uint64 {
	return h.boughtRedPacketCounts[dataId]
}
