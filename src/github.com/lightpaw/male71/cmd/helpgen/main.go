package main

import (
	conf "github.com/lightpaw/male7/config"
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/i18n"
	"sort"
	"text/template"
	"bytes"
	"path"
	"os"
	"io/ioutil"
	"go/format"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

func main() {

	// 读取配置，加载

	configPath := confpath.GetConfigPath()
	gos, err := config.NewConfigGameObjects(configPath)
	if err != nil {
		logrus.WithError(err).Panic("加载配置文件失败")
	}

	c, err := conf.ParseConfigDatas(gos)
	if err != nil {
		logrus.WithError(err).Panic("加载配置文件失败")
	}

	basePath := path.Dir(configPath)

	// i18n_text.help.go
	generateI18nHelp(basePath, c)

	// text.help.go
	generateTextHelp(basePath, c)

	// mail.help.go
	generateMailHelp(basePath, c)

	// broadcast.help.go
	generateBroadcastHelp(basePath, c)
}

func generateI18nHelp(basePath string, c *conf.ConfigDatas) {
	data := c.GetI18nData(i18n.DefaultLanguage)

	paramMap := make(map[string]int)
	for _, pair := range data.Pair {
		params := i18n.ExtractParams(pair.Value)
		for _, v := range params {
			paramMap[v]++
		}
	}

	var keys []string
	for k := range paramMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//fmt.Println(len(data.Pair), keys)

	t := template.Must(template.New("i18n").Funcs(template.FuncMap{"hump": HumpName,}).Parse(i18nTempStr))

	b := &bytes.Buffer{}
	if err := t.Execute(b, &keys); err != nil {
		logrus.WithError(err).Panic("template解析生成help.go失败")
	}

	filename := path.Join(basePath, "config", "i18n", "i18n_text.help.go")
	if err := WriteSource(filename, b.Bytes()); err != nil {
		logrus.WithError(err).Panic("写入help.go失败")
	}

	fmt.Println("生成成功", path.Base(filename))
}

const i18nTempStr = `
package i18n
{{range .}}
func (f *Fields) With{{hump .}}({{.}} interface{}) *Fields {
	return f.WithFields("{{.}}", {{.}})
}
{{end}}
`

func generateTextHelp(basePath string, c *conf.ConfigDatas) {

	var keys []string
	for _, t := range c.GetTextArray() {
		keys = append(keys, t.Id)
	}
	sort.Strings(keys)

	t := template.Must(template.New("text").Parse(textHelpTempStr))

	b := &bytes.Buffer{}
	if err := t.Execute(b, &keys); err != nil {
		logrus.WithError(err).Panic("template解析生成help.go失败")
	}

	filename := path.Join(basePath, "config", "data", "text.help.go")
	if err := WriteSource(filename, b.Bytes()); err != nil {
		logrus.WithError(err).Panic("写入help.go失败")
	}

	fmt.Println("生成成功", path.Base(filename))
}

const textHelpTempStr = `
package data

//gogen:config
type TextHelp struct {` +
	"	_ struct{} `singleton:\"true\"`\n" +
	"	_ struct{} `file:\"文字/文本.txt\"`\n" +
	`{{range .}}
		{{.}} *Text
	{{end}}
	}
	`

func generateMailHelp(basePath string, c *conf.ConfigDatas) {

	var keys []string
	for _, t := range c.GetMailDataArray() {
		keys = append(keys, t.Id)
	}
	sort.Strings(keys)

	t := template.Must(template.New("mail").Parse(mailHelpTempStr))

	b := &bytes.Buffer{}
	if err := t.Execute(b, &keys); err != nil {
		logrus.WithError(err).Panic("template解析生成help.go失败")
	}

	filename := path.Join(basePath, "config", "maildata", "mail.help.go")
	if err := WriteSource(filename, b.Bytes()); err != nil {
		logrus.WithError(err).Panic("写入help.go失败")
	}

	fmt.Println("生成成功", path.Base(filename))
}

const mailHelpTempStr = `
package maildata

//gogen:config
type MailHelp struct {` +
	"	_ struct{} `singleton:\"true\"`\n" +
	"	_ struct{} `file:\"文字/邮件.txt\"`\n" +
	`{{range .}}
		{{.}} *MailData
	{{end}}
	}
	`

func generateBroadcastHelp(basePath string, c *conf.ConfigDatas) {

	var keys []string
	for _, t := range c.GetBroadcastDataArray() {
		keys = append(keys, t.Id)
	}
	sort.Strings(keys)

	t := template.Must(template.New("mail").Parse(broadcastHelpTempStr))

	b := &bytes.Buffer{}
	if err := t.Execute(b, &keys); err != nil {
		logrus.WithError(err).Panic("template解析生成help.go失败")
	}

	filename := path.Join(basePath, "config", "data", "broadcast.help.go")
	if err := WriteSource(filename, b.Bytes()); err != nil {
		logrus.WithError(err).Panic("写入help.go失败")
	}

	fmt.Println("生成成功", path.Base(filename))
}

const broadcastHelpTempStr = `
package data

//gogen:config
type BroadcastHelp struct {` +
	"	_ struct{} `singleton:\"true\"`\n" +
	"	_ struct{} `file:\"文字/广播.txt\"`\n" +
	`{{range .}}
		{{.}} *BroadcastData
	{{end}}
	}
	`

func WriteSource(filename string, data []byte) error {
	source, err := format.Source(data)
	if err != nil {
		WriteFile(filename, data) // 把文件写进去
		return errors.Wrapf(err, "源代码有错: %s \n %s", string(data), filename)
	}
	return WriteFile(filename, source)
}

func WriteFile(filename string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}

func HumpName(in string) string {
	return strings.Replace(strings.Title(strings.Replace(in, "_", " ", -1)), " ", "", -1)
}
