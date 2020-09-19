package util

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestGetClientIp(t *testing.T) {
	RegisterTestingT(t)

	Ω(ToU32Ip([4]byte{127, 0, 0, 1})).Should(Equal(uint32(0x100007f)))
}
