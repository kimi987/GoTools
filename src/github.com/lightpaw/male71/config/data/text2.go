package data

import (
	"bytes"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

func ParseText2(text string) (*Text2, error) {

	array := strings.Split(text, "%%")

	// 从第二段开始，将第一个字符拿下，跟位置拼起来，位置在后面
	var paramRanks []string
	for i, a := range array {
		if i <= 0 {
			continue
		}

		if len(a) <= 0 {
			return nil, errors.Errorf("ParseText2 无效的字符，%s", text)
		}

		rk := a[0:1] + strconv.Itoa(i-1)
		paramRanks = append(paramRanks, rk)
		array[i] = a[1:]
	}

	sort.Strings(paramRanks)
	paramIndex := make([]int, len(paramRanks))
	for i, v := range paramRanks {
		idx, err := strconv.Atoi(v[len(v)-1:])
		if err != nil {
			return nil, errors.Errorf("ParseText2 解析index失败，%s", text)
		}

		paramIndex[i] = idx
	}

	t := &Text2{}
	t.originText = text
	t.array = array
	t.paramIndex = paramIndex

	return t, nil

}

// 支持参数位置可变
//gogen:config
type Text2 struct {
	originText string

	array []string

	paramIndex []int
}

func (t *Text2) OriginText() string {
	return t.originText
}

func (t *Text2) Format(args ...string) string {
	b := bytes.Buffer{}

	for i, a := range t.array {
		if i > 0 {
			idx := t.paramIndex[i-1]
			if idx < len(args) {
				b.WriteString(args[idx])
			}
		}

		b.WriteString(a)
	}

	return b.String()
}
