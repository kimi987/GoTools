package confpath

import "testing"
import (
	. "github.com/onsi/gomega"
	"path/filepath"
)

func TestFindConfigPath(t *testing.T) {
	RegisterTestingT(t)

	s, err := FindConfigPath("conf")
	Ω(err).Should(Succeed())

	es, err := filepath.Abs("../../conf")
	Ω(err).Should(Succeed())

	Ω(s).Should(Equal(es))
}
