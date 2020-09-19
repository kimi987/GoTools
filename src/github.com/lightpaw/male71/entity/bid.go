package entity

import (
	"time"
	"sync"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
)

/**
	竞价
  */

type BidInfo struct {
	endTime time.Time

	bidMap *sync.Map // k: 出价人id, v: BidObj
}

func NewBidInfo(endTime time.Time) *BidInfo {
	c := &BidInfo{}
	c.bidMap = &sync.Map{}
	c.endTime = endTime

	return c
}

type BidObj struct {
	k       int64     // 出价人id
	v       uint64    // 出价
	bidTime time.Time // 出价时间
}

func (c *BidInfo) Encode() *server_proto.BidInfoProto {
	p := &server_proto.BidInfoProto{}
	p.EndTime = timeutil.Marshal64(c.endTime)

	p.BidMap = make(map[int64]*server_proto.BidObjProto)
	c.WalkBid(func(obj *BidObj) {
		p.BidMap[obj.k] = obj.Encode()
	})

	return p
}

func (c *BidInfo) Unmarshal(p *server_proto.BidInfoProto) {
	c.endTime = timeutil.Unix64(p.EndTime)
	for k, v := range p.BidMap {
		b := &BidObj{}
		b.Unmarshal(v)
		c.bidMap.Store(k, b)
	}
}

func (c *BidObj) Encode() *server_proto.BidObjProto {
	p := &server_proto.BidObjProto{}
	p.K = c.k
	p.V = c.v
	p.BidTime = timeutil.Marshal64(c.bidTime)
	return p
}

func (c *BidObj) Unmarshal(p *server_proto.BidObjProto) {
	c.k = p.K
	c.v = p.V
	c.bidTime = timeutil.Unix64(p.BidTime)
}

func (c *BidObj) K() int64 {
	return c.k
}

func (c *BidObj) V() uint64 {
	return c.v
}

func (c *BidInfo) GetBid(k int64) uint64 {
	if v, ok := c.bidMap.Load(k); ok {
		return v.(*BidObj).v
	}
	return 0
}

func (c *BidInfo) WalkBid(f func(obj *BidObj)) {
	c.bidMap.Range(func(k, v interface{}) bool {
		o := v.(*BidObj)
		f(o)
		return true
	})
}

func (c *BidInfo) Bidding(k int64, v uint64, ctime time.Time) (oldV uint64, succ bool) {
	if ctime.After(c.endTime) {
		return
	}

	if oldObj, ok := c.bidMap.Load(k); ok {
		oldV = oldObj.(*BidObj).v
		if v <= oldV {
			return
		}
		oldObj.(*BidObj).v = v
		return oldV, true
	}

	value := &BidObj{k: k, v: v, bidTime: ctime}
	c.bidMap.Store(k, value)
	return 0, true
}

func (c *BidInfo) EndTime() time.Time {
	return c.endTime
}

func (c *BidInfo) Winner(ctime time.Time) (k int64, v uint64, succ bool) {
	var result *BidObj
	c.bidMap.Range(func(k, v interface{}) bool {
		newObj := v.(*BidObj)
		if result == nil || newObj.v > result.v {
			result = newObj
		} else if newObj.v == result.v && newObj.bidTime.Before(result.bidTime) {
			result = newObj
		}
		return true
	})

	if result == nil {
		return
	}

	return result.k, result.v, true
}

func (c *BidInfo) Losers(ctime time.Time) (losers map[int64]uint64, succ bool) {
	winner, _, _ := c.Winner(ctime)
	losers = make(map[int64]uint64)

	c.bidMap.Range(func(k, v interface{}) bool {
		o := v.(*BidObj)
		if o.K() == winner {
			return true
		}
		losers[o.K()] = o.V()
		return true
	})

	succ = true
	return
}
