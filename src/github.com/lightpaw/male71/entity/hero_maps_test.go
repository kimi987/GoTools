package entity

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/pb/server_proto"
	"strings"
)

func TestCategory(t *testing.T) {
	RegisterTestingT(t)

	categoryCountMap := make(map[int32]int32)
	for _, v := range dailyResetCategory {
		categoryCountMap[int32(v)]++
	}

	Ω(len(categoryCountMap)).Should(Equal(len(dailyResetCategory)))

	for name, category := range server_proto.HeroMapCategory_value {
		if strings.HasPrefix(strings.ToLower(name), "daily") {
			_, exist := categoryCountMap[category]
			Ω(exist).Should(BeTrue())
		}
	}
}
