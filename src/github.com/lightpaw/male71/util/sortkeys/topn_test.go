package sortkeys

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestU64TopN(t *testing.T) {
	RegisterTestingT(t)

	topN := NewU64TopN(5)

	topN.Add(1, "a")
	topN.Add(2, "b")
	topN.Add(4, "d")
	topN.Add(3, "c")

	Ω(topN.Array()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(1, "a"),
		NewU64KV(2, "b"),
		NewU64KV(4, "d"),
		NewU64KV(3, "c"),
	}))

	Ω(topN.SortAsc()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(1, "a"),
		NewU64KV(2, "b"),
		NewU64KV(3, "c"),
		NewU64KV(4, "d"),
	}))

	Ω(topN.SortDesc()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(4, "d"),
		NewU64KV(3, "c"),
		NewU64KV(2, "b"),
		NewU64KV(1, "a"),
	}))

	topN.Add(6, "f")

	Ω(topN.Array()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(1, "a"),
		NewU64KV(2, "b"),
		NewU64KV(4, "d"),
		NewU64KV(3, "c"),
		NewU64KV(6, "f"),
	}))

	topN.Add(5, "e")

	Ω(topN.Array()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(5, "e"),
		NewU64KV(2, "b"),
		NewU64KV(4, "d"),
		NewU64KV(3, "c"),
		NewU64KV(6, "f"),
	}))

	topN.Add(7, "g")

	Ω(topN.Array()).Should(BeEquivalentTo([]*U64KV{
		NewU64KV(5, "e"),
		NewU64KV(7, "g"),
		NewU64KV(4, "d"),
		NewU64KV(3, "c"),
		NewU64KV(6, "f"),
	}))
}
