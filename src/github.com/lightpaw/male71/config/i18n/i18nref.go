package i18n

import (
	"fmt"
	"strings"
	"github.com/lightpaw/logrus"
	"sort"
)

var i18nMap = make(map[string]*I18nRef)

func resetI18nMap() {
	for k := range i18nMap {
		delete(i18nMap, k)
	}
}

func getI18nKeys() []string {
	var keys []string
	for k := range i18nMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func NewI18nRef(filename, fieldName string, key interface{}, value string) *I18nRef {
	refKey := newRefKey(filename, fieldName, key)
	value = strings.TrimSpace(value)

	c := &I18nRef{
		Key:      refKey,
		Value:    value,
		paramMap: extractParamMap(value),
		asKey:    keyString(refKey),
		asJson:   keysOnlyJson(refKey),
	}

	if conf := i18nMap[c.Key]; conf != nil {
		logrus.WithField("key", c.Key).Panicf("国际化配置存在重复的key")
	}
	i18nMap[c.Key] = c

	return c
}

var refKeyPrefix = "s."

func GetConfigKey(refKey string) (ok bool, subKey string) {
	if !strings.HasPrefix(refKey, refKeyPrefix) {
		return
	}
	i := strings.LastIndex(refKey, "-")
	if i < 0 {
		return
	}

	subKey = refKey[i+1:]
	if len(subKey) <= 0 {
		return
	}
	ok = true
	return
}

func newRefKey(tablename, fieldname, key interface{}) string {
	// 拼接一个唯一id, table-fieldname-key
	return refKeyPrefix + strings.Replace(fmt.Sprintf("%v-%v-%v", tablename, fieldname, key), ".", "-", -1)
}

type I18nRef struct {
	Key string

	Value string

	paramMap map[string]int

	asKey  string
	asJson string
}

func (c *I18nRef) Encode() string {
	return c.Key
}

func ExtractParams(value string) (params []string) {
	return extractParams(value)
}

func extractParamMap(value string) (map[string]int) {
	array := extractParams(value)

	if len(array) > 0 {
		dataMap := make(map[string]int)
		for _, k := range array {
			dataMap[k]++
		}
		return dataMap
	}
	return nil
}

func extractParams(value string) (params []string) {
	if len(value) <= 0 {
		return
	}

	// 搞笑的飞机{{flag_name}}{{hero_name}} 打死不说{{搞笑}}，什么玩意{{what}}
	// return []string{"flag_name", "hero_name", "搞笑", "what", }

	// 根据{{split，然后找到每个值里面存在 }}的，把}} 去掉，剩下的就是参数
	n := strings.Count(value, "{{")
	s := value
	for i := 0; i < n; i++ {

		array := strings.SplitN(s, "{{", 2)
		if len(array) < 2 {
			// 到头了
			break
		}

		// 取后半部分，截取 }}之前的部分
		pa := strings.SplitN(array[1], "}}", 2)
		if len(pa) < 2 {
			// 到头了
			break
		}

		param := pa[0]
		s = pa[1]

		idx := strings.LastIndex(param, "{{")
		if idx > 0 {
			param = param[idx+2:]
		}

		params = append(params, param)
	}

	return
}
