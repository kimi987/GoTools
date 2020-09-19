package heroid

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestHeroId(t *testing.T) {
	RegisterTestingT(t)

	for i := 0; i < 1000; i++ {
		sid := uint32(i)
		for i := 0; i < 1000; i++ {
			accountId := int64(i)

			heroId := NewHeroId(accountId, sid)
			Ω(heroId >= 0).Should(BeTrue())
			Ω(GetAccountId(heroId)).Should(BeEquivalentTo(accountId))
			Ω(GetSid(heroId)).Should(BeEquivalentTo(sid))
		}
	}

}
