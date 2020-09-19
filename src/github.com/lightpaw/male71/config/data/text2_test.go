package data

import (
	"bytes"
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func newFmtText2(text *Text2) *fmt_text2 {

	originTextBuffer := bytes.Buffer{}
	for i, v := range text.array {
		if i > 0 {
			originTextBuffer.WriteString("%s")
		}
		originTextBuffer.WriteString(v)
	}

	t := &fmt_text2{}
	t.originText = originTextBuffer.String()
	t.paramIndex = text.paramIndex

	return t
}

type fmt_text2 struct {
	originText string

	paramIndex []int
}

func (t *fmt_text2) formatOrigin(args ...interface{}) string {
	newArgs := make([]interface{}, len(t.paramIndex))
	for i, idx := range t.paramIndex {
		if idx < len(args) {
			newArgs[i] = args[idx]
		} else {
			newArgs[i] = ""
		}
	}

	return fmt.Sprintf(t.originText, newArgs...)
}

func TestText2(t *testing.T) {
	RegisterTestingT(t)

	s := "由于您的持续掠夺，将%%3的城池从%%2级打降至%%1级了。"

	text, err := ParseText2(s)
	Ω(err).Should(Succeed())

	fmt.Println(newFmtText2(text).formatOrigin("3", "5", "SB"))
	fmt.Println(text.Format("3", "5", "SB"))

	Ω(newFmtText2(text).formatOrigin("3", "5", "SB")).Should(Equal("由于您的持续掠夺，将SB的城池从5级打降至3级了。"))
	Ω(text.Format("3", "5", "SB")).Should(Equal("由于您的持续掠夺，将SB的城池从5级打降至3级了。"))

	ft2 := fmtText2.formatOrigin("SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5")
	t2 := text2.Format("SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB")
	Ω(ft2).Should(Equal(t2))
	fmt.Println(ft2)
	fmt.Println(t2)
}

var text1, _ = ParseText2("由于您的持续掠夺，将%%s的城池从%%s级打降至%%s级了。")
var fmtText1 = newFmtText2(text1)

func BenchmarkText2_FormatOrigin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmtText1.formatOrigin("SB", "5", "3")
	}
}

func BenchmarkText2_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		text1.Format("SB", "5", "3")
	}
}

var text2, _ = ParseText2("由%%3于%%4您%%s的%%s持%%7续%%s掠%%8夺，将%%3的%%s城%%2池%%s从%%1级%%s打%%2降%%s至%%s级%%s了。")
var fmtText2 = newFmtText2(text2)

func BenchmarkText2_FormatOrigin2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmtText2.formatOrigin("SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5")
	}
}

func BenchmarkText2_Format2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		text2.Format("SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB", "5", "3", "SB")
	}
}
