package realmface

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	. "github.com/onsi/gomega"
	"testing"
)

func TestName(t *testing.T) {
	RegisterTestingT(t)

	id := GetGuildRealmId(1)
	Ω(ParseRegionSequence(id)).Should(Equal(uint64(1)))
	Ω(ParseRegionLevel(id)).Should(Equal(uint64(1)))
	Ω(ParseRegionType(id)).Should(Equal(shared_proto.RegionType_GUILD))

	i := -1
	Ω(i & 1).Should(Equal(1))

	id = GetRealmId(5, 6, shared_proto.RegionType_MONSTER)
	Ω(ParseRegionSequence(id)).Should(Equal(uint64(6)))
	Ω(ParseRegionLevel(id)).Should(Equal(uint64(5)))
	Ω(ParseRegionType(id)).Should(Equal(shared_proto.RegionType_MONSTER))

}
