package i18n

import (
	"github.com/lightpaw/jsoniter"
)

func SetDebugExample(b bool) {
	debugExample = b
}

var debugExample = false

// Fields type, used to pass to `WithFields`.
type Fields struct {
	key     *I18nRef
	dataMap map[string]interface{}
}

const nameKey = "i18nkey"
const exampleKey = "i18nexample"
const keyPrefix = "{{KEY}}"
const JsonPrefix = "{{JSON}}"

func keyString(key string) string {
	return keyPrefix + key
}

func keysOnlyJson(key string) string {
	return "{{JSON}}{\"i18nkey\":\"" + key + "\"}"
}

func (k *I18nRef) keyString() string {
	return k.asKey
}

func (k *I18nRef) KeysOnlyJson() string {
	if len(k.Value) > 0 {
		return k.asJson
	}
	return ""
}

func (k *I18nRef) New() *Fields {
	return newKey(k)
}

func newKey(k *I18nRef) *Fields {
	f := &Fields{
		key: k,
	}
	if len(k.paramMap) > 0 {
		f.dataMap = make(map[string]interface{})

		f.dataMap[nameKey] = k.Key
		if debugExample && len(k.Value) > 0 {
			f.dataMap[exampleKey] = k.Value
		}
	}

	return f
}

func NewKey(k *I18nRef) *Fields {
	return newKey(k)
}

func (f *Fields) WithFields(key string, value interface{}) *Fields {
	if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
		if ref, ok := value.(*I18nRef); ok {
			f.dataMap[key] = ref.keyString()
		} else {
			f.dataMap[key] = value
		}
	}
	return f
}

func (f *Fields) WithClickHeroFields(key string, heroName string, heroId int64) *Fields {
	//if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
	//	f.dataMap["$0" + key] = heroName + "," + strconv.FormatInt(heroId, 10)
	//}
	f.WithFields(key, heroName)
	return f
}

func (f *Fields) WithClickGuildFields(key string, guildName string, guildId int64) *Fields {
	//if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
	//	f.dataMap["$1" + key] = guildName + "," + strconv.FormatInt(guildId, 10)
	//}
	f.WithFields(key, guildName)
	return f
}

func (f *Fields) WithClickEquipFields(key string, equipName string, dataId, level, refined uint64) *Fields {
	//if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
	//	f.dataMap["$2" + key] = equipName + "," + strconv.FormatUint(dataId, 10) + "," + strconv.FormatUint(level, 10) + "," + strconv.FormatUint(refined, 10)
	//}
	f.WithFields(key, equipName)
	return f
}

func (f *Fields) WithClickCaptainFields(key string, captainName string, heroId int64, captainId uint64) *Fields {
	//if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
	//	f.dataMap["$3" + key] = captainName + "," + strconv.FormatInt(heroId, 10) + "," + strconv.FormatUint(captainId, 10)
	//}
	f.WithFields(key, captainName)
	return f
}

func (f *Fields) WithClickMingcFields(key string, mingcName string, mingcId uint64) *Fields {
	//if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
	//	f.dataMap["$4" + key] = mingcName + "," + strconv.FormatUint(mingcId, 10)
	//}
	f.WithFields(key, mingcName)
	return f
}

func (f *Fields) WithKey(key string, value *I18nRef) *Fields {
	if len(f.key.paramMap) > 0 && f.key.paramMap[key] > 0 {
		f.dataMap[key] = value.keyString()
	}
	return f
}

func (f *Fields) JsonString() string {
	if len(f.dataMap) > 0 {
		s, _ := jsoniter.MarshalToString(f.dataMap)
		return JsonPrefix + s
	}
	return f.key.KeysOnlyJson()
}
