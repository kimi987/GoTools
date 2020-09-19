package data

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestCompareCondition(t *testing.T) {
	RegisterTestingT(t)

	c, err := ParseCompareCondition(">-1")
	Ω(err).Should(HaveOccurred())
	c, err = ParseCompareCondition("a")
	Ω(err).Should(HaveOccurred())

	c, err = ParseCompareCondition(">5")
	Ω(err).Should(Succeed())
	Ω(c.Compare(5)).Should(BeFalse())
	Ω(c.Compare(6)).Should(BeTrue())
	Ω(c.Compare(3)).Should(BeFalse())

	c, err = ParseCompareCondition("<>5")
	Ω(err).Should(Succeed())
	Ω(c.Compare(5)).Should(BeFalse())
	Ω(c.Compare(6)).Should(BeTrue())
	Ω(c.Compare(3)).Should(BeTrue())

	c, err = ParseCompareCondition("<=5")
	Ω(err).Should(Succeed())
	Ω(c.Compare(5)).Should(BeTrue())
	Ω(c.Compare(6)).Should(BeFalse())
	Ω(c.Compare(3)).Should(BeTrue())
}
