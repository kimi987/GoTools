package transform

import "testing"
import . "github.com/onsi/gomega"

func TestEnumMapKeys(t *testing.T) {
	RegisterTestingT(t)

	in := map[string]int32{
		"1": 100,
		"2": 200,
		"3": 300,
	}

	Ω(EnumMapKeys(in)).Should(ConsistOf([]string{
		"1",
		"2",
		"3",
	}))

	Ω(EnumMapKeys(in, 200)).Should(ConsistOf([]string{
		"1",
		"3",
	}))

	Ω(EnumMapKeys(in, 100, 300)).Should(ConsistOf([]string{
		"2",
	}))

}
