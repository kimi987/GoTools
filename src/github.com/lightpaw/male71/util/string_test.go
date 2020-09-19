package util

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestByte2String(t *testing.T) {
	RegisterTestingT(t)
	str := "abcdefg"
	barr := []byte(str)

	newStr := Byte2String(barr)
	fmt.Println(newStr)
	Ω(newStr).Should(Equal(str))

	for i := range barr {
		barr[i] = 0
	}
	Ω(newStr).ShouldNot(Equal(str))
}

func TestString2Byte(t *testing.T) {
	RegisterTestingT(t)

	str := "abcdefg"
	barr := String2Byte(str)

	newStr := string(barr)

	Ω(newStr).Should(Equal(str))
}

func TestGetCharLen(t *testing.T) {
	RegisterTestingT(t)

	n := GetCharLen("")
	Ω(n).Should(Equal(0))

	n = GetCharLen("1a2b3c")
	Ω(n).Should(Equal(6))

	n = GetCharLen("哈哈")
	Ω(n).Should(Equal(4))

	n = GetCharLen("哈哈123")
	Ω(n).Should(Equal(7))

	n = GetCharLen("哈哈哈哈哈1213")
	Ω(n).Should(Equal(14))
}

func TestInvalidChar(t *testing.T) {
	RegisterTestingT(t)

	Ω(HaveInvalidChar("啊ab_1")).Should(BeFalse())
	Ω(HaveInvalidChar("1")).Should(BeFalse()) // number
	Ω(HaveInvalidChar(" ")).Should(BeFalse())

	Ω(HaveInvalidChar("!")).Should(BeTrue())
	Ω(HaveInvalidChar("$")).Should(BeTrue())
	Ω(HaveInvalidChar("#")).Should(BeTrue())
	Ω(HaveInvalidChar("@")).Should(BeTrue())
	Ω(HaveInvalidChar("%")).Should(BeTrue())
	Ω(HaveInvalidChar("~")).Should(BeTrue())
	Ω(HaveInvalidChar("`")).Should(BeTrue())
	Ω(HaveInvalidChar("￥")).Should(BeTrue())
	Ω(HaveInvalidChar("-")).Should(BeTrue())
	Ω(HaveInvalidChar("+")).Should(BeTrue())
	Ω(HaveInvalidChar("=")).Should(BeTrue())
	Ω(HaveInvalidChar("\r")).Should(BeTrue())
	Ω(HaveInvalidChar("\n")).Should(BeTrue())
}

func TestValidName(t *testing.T) {
	RegisterTestingT(t)

	Ω(IsValidName("啊ab_1")).Should(BeTrue())
	Ω(IsValidName("1")).Should(BeTrue())

	Ω(IsValidName(" 1")).Should(BeFalse())
	Ω(IsValidName(" 啊")).Should(BeFalse())
	Ω(IsValidName("啊 ")).Should(BeFalse())
	Ω(IsValidName("_啊")).Should(BeFalse())
	Ω(IsValidName("啊_")).Should(BeFalse())
}
