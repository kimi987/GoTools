package config

import (
	"testing"
	"github.com/lightpaw/male7/config/confpath"
	. "github.com/onsi/gomega"
)

func TestInit(t *testing.T) {
	RegisterTestingT(t)

	_, err := LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())
}
