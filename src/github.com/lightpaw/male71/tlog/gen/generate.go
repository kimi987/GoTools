package gen

import (
	"bytes"
	"path"
	"fmt"
	"io/ioutil"
	"text/template"
	"strings"
	"os"
	"go/format"
	"github.com/pkg/errors"
)

func GenTlogService(xmlPath, fullFileName string) error {

	xmlObj, err := Unmarshal(xmlPath)
	if err != nil {
		return err
	}

	t := template.Must(template.New("tlog").Funcs(template.FuncMap{"hump": HumpName, "increment": Increment}).Parse(tmpl))

	b := &bytes.Buffer{}
	if err := t.Execute(b, xmlObj); err != nil {
		return err
	}

	if err := WriteSource(fullFileName, b.Bytes()); err != nil {
		return err
	}

	fmt.Println("生成成功", path.Base(fullFileName))

	return nil
}

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

func Increment(num int) int {
	return num + 1
}
