package data

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestRate(t *testing.T) {
	RegisterTestingT(t)

	rate, err := ParseRate("abc")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("-1")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("10001")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("101/100")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("-1/100")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("0/0")
	Ω(err).Should(HaveOccurred())

	rate, err = ParseRate("")
	Ω(err).Should(Succeed())
	Ω(rate.Try()).Should(BeFalse())

	rate, err = ParseRate("0")
	Ω(err).Should(Succeed())
	Ω(rate.Try()).Should(BeFalse())

	rate, err = ParseRate("0/100")
	Ω(err).Should(Succeed())
	Ω(rate.Try()).Should(BeFalse())

	rate, err = ParseRate("10000")
	Ω(err).Should(Succeed())
	Ω(rate.Try()).Should(BeTrue())

	rate, err = ParseRate("100/100")
	Ω(err).Should(Succeed())
	Ω(rate.Try()).Should(BeTrue())

	rate, err = ParseRate("50/100")
	Ω(err).Should(Succeed())

}
