package entity

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestNewClientDatas(t *testing.T) {
	RegisterTestingT(t)

	datas := newClientDatas(newHeroMaps())

	datas.unmarshal(nil)

	cp1 := datas.encodeClient()
	Ω(len(cp1.IntValue)).Should(BeEquivalentTo(0))

	sp1 := datas.encodeServer()
	for _, value := range sp1.IntValue {
		Ω(value).Should(BeEquivalentTo(0))
	}

	size := len(datas.int32Array) * Int32BitCount
	for i := 0; i < size; i++ {
		datas.SetBool(i, true)
		Ω(datas.Bool(i)).Should(BeTrue())
	}

	cp2 := datas.encodeClient()

	Ω(len(cp2.IntValue)).Should(BeEquivalentTo(size))

	mp := map[int32]struct{}{}
	for _, v := range cp2.IntValue {
		mp[v] = struct{}{}
		Ω(int(v) < size && v >= 0).Should(BeTrue())
	}

	Ω(len(mp)).Should(BeEquivalentTo(size))

	var v uint32 = 0xffffffff
	sp2 := datas.encodeServer()
	for _, value := range sp2.IntValue {
		Ω(value).Should(BeEquivalentTo(int32(v)))
	}

	datas.unmarshal(sp2)

	cp4 := datas.encodeClient()

	Ω(len(cp4.IntValue)).Should(BeEquivalentTo(size))

	mp = map[int32]struct{}{}
	for _, v := range cp4.IntValue {
		mp[v] = struct{}{}
		Ω(int(v) < size && v >= 0).Should(BeTrue())
	}

	Ω(len(mp)).Should(BeEquivalentTo(size))

	sp4 := datas.encodeServer()
	for _, value := range sp4.IntValue {
		Ω(value).Should(BeEquivalentTo(int32(v)))
	}

	for i := 0; i < size; i++ {
		datas.SetBool(i, false)
		Ω(datas.Bool(i)).Should(BeFalse())
	}

	for i := 0; i < size; i++ {
		datas.SetBool(i, false)
		Ω(datas.Bool(i)).Should(BeFalse())
	}

	cp3 := datas.encodeClient()
	Ω(len(cp3.IntValue)).Should(BeEquivalentTo(0))

	sp3 := datas.encodeServer()
	for _, value := range sp3.IntValue {
		Ω(value).Should(BeEquivalentTo(0))
	}

	for i := size; i < size*2; i++ {
		datas.SetBool(i, true)
		Ω(datas.Bool(i)).Should(BeFalse())
		datas.SetBool(i, false)
		Ω(datas.Bool(i)).Should(BeFalse())
	}
}
