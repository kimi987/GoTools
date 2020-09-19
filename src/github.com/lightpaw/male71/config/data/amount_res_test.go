package data

import (
	"fmt"
	"github.com/lightpaw/male7/pb/shared_proto"
	. "github.com/onsi/gomega"
	"sort"
	"testing"
)

func TestAmount_IsZero(t *testing.T) {

	ia := []int{0, 1, 2, 4, 4, 4, 6, 8}

	sort.Sort(sort.Reverse(sort.IntSlice(ia)))
	fmt.Println(ia)
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= -1 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 0 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 1 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 2 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 3 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 4 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 5 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 6 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 7 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 8 }))
	fmt.Println(sort.Search(len(ia), func(i int) bool { return ia[i] <= 9 }))
	//fmt.Println(sort.SearchInts(ia, 5))
	//fmt.Println(sort.SearchInts(ia, 6))

	//fmt.Println(sort.SearchInts(ia, -1))
	//fmt.Println(sort.SearchInts(ia, 0))
	//fmt.Println(sort.SearchInts(ia, 1))
	//fmt.Println(sort.SearchInts(ia, 2))
	//fmt.Println(sort.SearchInts(ia, 3))
	//fmt.Println(sort.SearchInts(ia, 4))
	//fmt.Println(sort.SearchInts(ia, 5))
	//fmt.Println(sort.SearchInts(ia, 6))
	//fmt.Println(sort.SearchInts(ia, 7))
	//fmt.Println(sort.SearchInts(ia, 8))
	//fmt.Println(sort.SearchInts(ia, 9))
}

func Test1(t *testing.T) {

	ia := make([]int, 5)

	ia2 := ia[:]
	ia2 = append(ia2, 1, 2, 3)
	ia2[0] = 3

	fmt.Println(ia)

	fmt.Println(ia2)
}

func TestAmount(t *testing.T) {
	RegisterTestingT(t)

	rs, err := ParseAmount("100+50%")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&Amount{Amount: 100, Percent: 50}))

	rs, err = ParseAmount("100")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&Amount{Amount: 100}))

	rs, err = ParseAmount("50%")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&Amount{Percent: 50}))

	rs, err = ParseAmount("")
	Ω(err).Should(Succeed())
	Ω(rs).Should(BeNil())
}

func TestResAmount(t *testing.T) {
	RegisterTestingT(t)

	rs, err := ParseResAmount("gold:100+50%")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&ResAmount{Type: shared_proto.ResType_GOLD, Amount: 100, Percent: 50}))

	rs, err = ParseResAmount("fooD:100")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&ResAmount{Type: shared_proto.ResType_FOOD, Amount: 100}))

	rs, err = ParseResAmount("WOOD:50%")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&ResAmount{Type: shared_proto.ResType_WOOD, Percent: 50}))

	rs, err = ParseResAmount("stone:00+50%")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&ResAmount{Type: shared_proto.ResType_STONE, Percent: 50}))

	rs, err = ParseResAmount("stone:")
	Ω(err).Should(Succeed())
	Ω(rs).Should(Equal(&ResAmount{Type: shared_proto.ResType_STONE}))

	// fail

	rs, err = ParseResAmount("gold")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold100+50%")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold:+50%")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold:100+%")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold:100+50")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold:100+5")
	Ω(err).Should(HaveOccurred())

	rs, err = ParseResAmount("gold:100+5")
	Ω(err).Should(HaveOccurred())

}
