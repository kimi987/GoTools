package data

import (
	"bytes"
	"fmt"
	"strings"
)

func ParseTextFormatter(s string) (*TextFormatter, error) {
	t := &TextFormatter{}
	t.Text = strings.Split(s, ";")

	if len(t.Text) > 1 {
		b := bytes.Buffer{}
		for i := 0; i < len(t.Text); i++ {
			b.WriteString(t.Text[i])
		}
		t.OneText = b.String()
	} else {
		t.OneText = t.Text[0]
	}

	return t, nil
}

//gogen:config
type TextFormatter struct {
	Text []string

	OneText string `head:"-"`
}

func (t *TextFormatter) Format(args ...interface{}) string {
	return fmt.Sprintf(t.OneText, args...)
}

func (t *TextFormatter) FormatIgnoreEmpty(args ...interface{}) string {

	if len(t.Text) > len(args) {
		// 参数比这内容还少，当成普通类型处理
		return fmt.Sprintf(t.OneText, args...)
	} else {
		b := bytes.Buffer{}

		n := len(t.Text)
		if n < len(args) {
			n--
			for i := 0; i < n; i++ {
				if s, ok := args[i].(string); ok && len(s) == 0 {
					continue
				}

				b.WriteString(fmt.Sprintf(t.Text[i], args[i]))
			}
			b.WriteString(fmt.Sprintf(t.Text[n], args[n:]...))
		} else {
			for i := 0; i < n; i++ {
				if s, ok := args[i].(string); ok && len(s) == 0 {
					continue
				}

				b.WriteString(fmt.Sprintf(t.Text[i], args[i]))
			}
		}

		return b.String()
	}
}
