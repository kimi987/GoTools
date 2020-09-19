package i18n

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestExtractParams(t *testing.T) {
	RegisterTestingT(t)

	shouldEmptyString := []string{
		"", "哈哈",
		"{{", "{{哈哈", "哈哈{{",
		"}}", "哈哈}}", "}}哈哈",
		"{}}", "}{}", "}}{", "}}{{",
		"{哈哈}}", "}{}哈哈", "哈哈}}{", "}}哈哈{{",
	}

	for _, v := range shouldEmptyString {
		params := extractParams(v)
		Ω(params).Should(BeEmpty())
	}

	params := extractParams("{{}}")
	Ω(params).Should(Equal([]string{""}))

	params = extractParams("{{哈哈}}")
	Ω(params).Should(Equal([]string{"哈哈"}))

	params = extractParams("呵呵{{搞笑{{哈哈}}嘿}}嘿")
	Ω(params).Should(Equal([]string{"哈哈"}))

	params = extractParams("搞笑的飞机{{flag_name}}{{hero_name}} 打死不说{{搞笑}}，什么玩意{{what}}")
	Ω(params).Should(Equal([]string{"flag_name", "hero_name", "搞笑", "what"}))

	// s.文字.Text.1 = "{{hero_name}}击败了[color="sss"]{{target}}[]"
	// s.匈奴大营 = "匈奴大营"
	// {{JSON}}{"i18nkey":"s.文字.Text.1","hero_name":"李先生","target":"[color="sss"]{{target}}[]"}

}

func TestName(t *testing.T) {
	RegisterTestingT(t)

	Ω(newRefKey("文字/文本.txt", "text", "RealmTroopSpeedUp")).Should(Equal("s.文字/文本-txt-text-RealmTroopSpeedUp"))

}

func TestReplace(t *testing.T) {
	RegisterTestingT(t)

	s := "搞笑\\n我们是认真的"
	Ω(replaceNewLine(s)).Should(Equal("搞笑\n我们是认真的"))
}
