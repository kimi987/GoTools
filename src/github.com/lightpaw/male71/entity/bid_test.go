package entity

import (
	"testing"
	. "github.com/onsi/gomega"
	"time"
	"fmt"
)

func TestBidInfo_Bidding(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()

	endBid := NewBidInfo(ctime.Add(-time.Minute))
	v, succ := endBid.Bidding(1, 2, ctime)
	Ω(succ).Should(Equal(false))

	succBid := NewBidInfo(ctime.Add(time.Hour))
	v, succ = succBid.Bidding(1, 10, ctime)
	Ω(succ).Should(Equal(true))
	Ω(v).Should(Equal(uint64(0)))

	v, succ = succBid.Bidding(2, 20, ctime)
	Ω(succ).Should(Equal(true))
	Ω(v).Should(Equal(uint64(0)))

	v, succ = succBid.Bidding(1, 100, ctime)
	Ω(succ).Should(Equal(true))
	Ω(v).Should(Equal(uint64(10)))

	ft, err := time.Parse("Monday 15:04:05", "Sunday 20:00:08")

	fmt.Println(err)
	fmt.Println(ft)
	fmt.Println(time.Time{})
	fmt.Println(ft.AddDate(1, 0, 0).Sub(time.Time{}))
}

func TestBidInfo_GetBid(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()

	b := NewBidInfo(ctime.Add(time.Hour))
	b.Bidding(1, 10, ctime)

	Ω(b.GetBid(1)).Should(Equal(uint64(10)))

	b.Bidding(1, 100, ctime)
	Ω(b.GetBid(1)).Should(Equal(uint64(100)))
}

func TestBidInfo_Winner(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()
	b := NewBidInfo(ctime.Add(time.Hour))
	b.Bidding(1, 10, ctime)
	b.Bidding(2, 20, ctime)
	b.Bidding(3, 30, ctime)
	b.Bidding(4, 30, ctime.Add(-time.Second))

	k, v, succ := b.Winner(ctime.Add(time.Hour))
	Ω(succ).Should(Equal(true))
	Ω(k).Should(Equal(int64(4)))
	Ω(v).Should(Equal(uint64(30)))
}
