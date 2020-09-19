package i18n

import (
	"github.com/lightpaw/config"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"bytes"
	"path"
	"os"
	"io/ioutil"
	"github.com/axgle/mahonia"
	"sort"
	"strings"
	"github.com/lightpaw/male7/util/check"
	"reflect"
)

const DefaultLanguage = "cn"

//gogen:config
type I18nData struct {
	_ struct{} `file:"i18n/语言.txt"`
	_ struct{} `proto:"shared_proto.I18NDataProto"`
	_ struct{} `protoconfig:"I18N"`

	Id string

	Language string
	//LanguageRef *I18nRef `head:"language"`

	Display bool

	Pair []*I18nPair `head:"-"`

	itemMap map[string]*item

	defaultData *I18nData
}

type item struct {
	key string

	value string

	translation string
}

const (
	pattern       = "i18n/i18n_%s.txt"
	diffPattern   = "i18n/i18n_%s.diff.txt"
	deletePattern = "i18n/i18n_%s.delete.txt"
)

func (*I18nData) InitAll(fileName string, dataMap map[string]*I18nData) {
	resetI18nMap() // 清掉这个数据

	defaultData := dataMap[DefaultLanguage]
	check.PanicNotTrue(defaultData != nil, "%s 没有配置默认的语言 %s", fileName, DefaultLanguage)

	for _, d := range dataMap {
		d.defaultData = defaultData

		if d.Display {
			if err := d.CheckError(); err != nil {
				logrus.WithError(err).Panic("国际化文件初始化失败")
			}
		}
	}
}

func (d *I18nData) Init99(gos *config.GameObjects, fileName string) {

	if strings.ToLower(d.Id) == DefaultLanguage {
		itemMap := make(map[string]*item)
		for _, ref := range i18nMap {
			itemMap[ref.Key] = &item{
				key:         ref.Key,
				value:       replaceNewLine(ref.Value),
				translation: replaceNewLine(ref.Value),
			}
		}
		d.itemMap = itemMap
		d.Pair = itemMap2Pairs(itemMap)
	} else {
		filePath := fmt.Sprintf(pattern, d.Language)
		logrus.WithField("file", filePath).Debug("加载国际化文件")

		if array, err := gos.LoadFile(filePath); err != nil {
			logrus.WithField("file", filePath).WithError(err).Panic("加载国际化文件失败")
		} else {
			itemMap := make(map[string]*item)
			for _, p := range array {
				key := p.String("key")
				value := p.String("value")
				tr := p.String("translation")
				itemMap[key] = &item{
					key:         key,
					value:       replaceNewLine(value),
					translation: replaceNewLine(tr),
				}
			}
			d.itemMap = itemMap

			d.Pair = itemMap2Pairs(itemMap)
		}
	}
}

func replaceNewLine(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func (d *I18nData) getDefaultI18nItems() []*item {
	itemMap := d.defaultData.itemMap

	var keys []string
	for k := range itemMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var items []*item
	for _, k := range keys {
		items = append(items, itemMap[k])
	}

	return items
}

func (d *I18nData) CheckError() error {

	// 注册的Key，在这里都找得到，而且都翻译了
	refs := d.getDefaultI18nItems()
	for _, ref := range refs {
		k := ref.key
		item := d.itemMap[k]

		if item == nil {
			return errors.Errorf("语言 %s 的翻译项 %s 没找到", d.Language, k)
		}

		if item.value != ref.value {
			return errors.Errorf("语言 %s 的翻译项 %s 有更新，请重新翻译", d.Language, k)
		}

		if len(item.value) > 0 && len(item.translation) <= 0 {
			return errors.Errorf("语言 %s 的翻译项 %s 翻译内容为空", d.Language, k)
		}

		valueParamMap := extractParamMap(item.value)
		transParamMap := extractParamMap(item.translation)
		if !reflect.DeepEqual(valueParamMap, transParamMap) {
			return errors.Errorf("语言 %s 的翻译项 %s 翻译内容的参数不一致，中文：%s 翻译：%s", d.Language, k, item.value, item.translation)
		}

	}

	return nil
}

func (d *I18nData) Generate(gos *config.GameObjects, basePath string) error {
	if strings.ToLower(d.Id) == DefaultLanguage {
		return nil
	}

	filePath := fmt.Sprintf(pattern, d.Language)
	diffFilePath := fmt.Sprintf(diffPattern, d.Language)
	deleteFilePath := fmt.Sprintf(deletePattern, d.Language)

	list, err := gos.LoadFile(diffFilePath)
	if err != nil {
		return errors.Wrapf(err, "加载 %s 配置失败", diffFilePath)
	}

	diffMap := make(map[string]*config.ObjectParser)
	for _, p := range list {
		key := p.String("key")
		if len(key) > 0 {
			diffMap[key] = p
		}
	}

	b := bytes.Buffer{}
	b.WriteString("不要改	中文	翻译\nkey	value	translation\n")

	// 生成i18n
	refs := d.getDefaultI18nItems()
	for _, ref := range refs {
		k := ref.key
		if len(ref.value) <= 0 {
			continue
		}

		item := d.itemMap[k]
		diff := diffMap[k]

		if diff != nil {
			newValue := diff.String("new_value")
			newTranslation := diff.String("new_translation")
			if ref.value == newValue && len(newTranslation) > 0 {
				b.WriteString(fmt.Sprintf("%v	%v	%v\n", k, ref.value, newTranslation))
				continue
			}
		}

		if item != nil {
			if ref.value == item.value {
				b.WriteString(fmt.Sprintf("%v	%v	%v\n", k, ref.value, item.translation))
				continue
			}
		}

		b.WriteString(fmt.Sprintf("%v	%v	%v\n", k, ref.value, ""))
	}

	if err := writeConfigFile(basePath, filePath, b.Bytes()); err != nil {
		return err
	}

	// 生成 i18n.diff
	b.Reset()
	b.WriteString("不要改	旧中文	旧翻译	新中文	新翻译\nkey	value	translation	new_value	new_translation\n")
	n := b.Len()
	for _, ref := range refs {
		k := ref.key
		if len(ref.value) <= 0 {
			continue
		}

		item := d.itemMap[k]
		diff := diffMap[k]

		var value, translation string
		if diff != nil {
			value = diff.String("value")
			translation = diff.String("translation")

			newValue := diff.String("new_value")
			newTranslation := diff.String("new_translation")
			if ref.value == newValue && len(newTranslation) > 0 {
				continue
			}
		}

		if item != nil {
			if ref.value == item.value {
				if len(item.translation) > 0 {
					continue
				}
			} else {
				value = item.value
				translation = item.translation
			}

		}

		b.WriteString(fmt.Sprintf("%v	%v	%v	%v	%v\n", k, value, translation, ref.value, ""))
	}

	if n < b.Len() {
		// 有新的内容写入
		if err := writeConfigFile(basePath, diffFilePath, b.Bytes()); err != nil {
			return err
		}
	} else {
		// 删掉这个文件
		os.Remove(path.Join(basePath, diffFilePath))
	}

	var removeKeys []string
	for k := range d.itemMap {
		ref := d.defaultData.itemMap[k]
		if ref == nil {
			removeKeys = append(removeKeys, k)
		}
	}

	if len(removeKeys) > 0 {
		b.Reset()
		for _, k := range removeKeys {
			item := d.itemMap[k]
			b.WriteString(fmt.Sprintf("%v	%v	%v\n", item.key, item.value, item.translation))
		}

		if data := gos.Bytes(deleteFilePath); len(data) > 0 {
			// 这里原始的就是gbk
			gbkString := gbkEncoder.ConvertString(b.String())
			b.Reset()
			b.Write(data)
			b.WriteString(gbkString)

			if err := writeFile(basePath, deleteFilePath, b.Bytes()); err != nil {
				return err
			}
		} else {
			if err := writeConfigFile(basePath, deleteFilePath, b.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}

var gbkEncoder = mahonia.NewEncoder("gbk")

func writeConfigFile(basePath, fileName string, bytes []byte) error {
	return writeFile(basePath, fileName, []byte(gbkEncoder.ConvertString(string(bytes))))
}

func writeFile(basePath, fileName string, bytes []byte) error {

	fileName = path.Join(basePath, fileName)

	err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "输出国际化文件失败，%v", fileName)
	}

	err = ioutil.WriteFile(fileName, bytes, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "输出国际化文件失败，%v", fileName)
	}

	return nil
}

// pair

//gogen:config
type I18nPair struct {
	_ struct{} `proto:"shared_proto.I18NPairProto"`

	Key []string

	Value string
}

type pairslice []*I18nPair

func (p pairslice) Len() int           { return len(p) }
func (p pairslice) Less(i, j int) bool { return p[i].Key[0] < p[j].Key[0] }
func (p pairslice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func itemMap2Pairs(itemMap map[string]*item) []*I18nPair {

	pairMap := make(map[string]*I18nPair)

	keys := getI18nKeys()
	for _, k := range keys {
		item := itemMap[k]
		if item == nil {
			continue
		}

		pairKey := item.translation
		pair := pairMap[pairKey]
		if pair == nil {
			pair = &I18nPair{
				Value: pairKey,
			}
			pairMap[pairKey] = pair
		}

		pair.Key = append(pair.Key, item.key)
	}

	var pairs []*I18nPair
	for _, v := range pairMap {
		pairs = append(pairs, v)
	}
	sort.Sort(pairslice(pairs))

	return pairs
}
