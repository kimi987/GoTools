package config

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

type GameObjects struct {
	dataMap map[string][]byte

	gbkDecoder    mahonia.Decoder
	decoderLocker sync.Mutex
}

func (g *GameObjects) Data(filename string) string {
	return g.bytes2string(g.dataMap[filename])
}

func (g *GameObjects) Bytes(filename string) []byte {
	return g.dataMap[filename]
}

func (g *GameObjects) bytes2string(data []byte) string {
	if len(data) <= 0 || utf8.Valid(data) {
		return string(data)
	}

	g.decoderLocker.Lock()
	defer g.decoderLocker.Unlock()
	return g.gbkDecoder.ConvertString(string(data))
}

func (g *GameObjects) LoadFile(filename string) ([]*ObjectParser, error) {
	return ParseList(filename, g.Data(filename))
}

func NewConfigGameObjects(dir string) (*GameObjects, error) {

	fmt.Println(filepath.Abs(dir))
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	dir = strings.Replace(dir, "\\", "/", -1) // 更新windows支持

	gos := &GameObjects{}
	gos.dataMap = make(map[string][]byte)
	gos.gbkDecoder = mahonia.NewDecoder("gbk")

	fmt.Println(filepath.Abs(dir))

	err1 := filepath.Walk(dir, func(path0 string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err1 := ioutil.ReadFile(path0)
		if err1 != nil {
			return errors.Wrapf(err1, "read config fail, %s", path0)
		}

		dp := strings.Replace(path0, "\\", "/", -1) // 更新windows支持
		dp = strings.Replace(dp, dir, "", -1)

		gos.dataMap[dp] = data

		return nil
	})

	return gos, err1
}

func NewKeyValueGameObjects(name, content string, kv ...string) (*GameObjects, error) {

	if name == "" {
		return nil, errors.Errorf("name empty")
	}

	if content == "" {
		return nil, errors.Errorf("content empty")
	}

	gos := &GameObjects{dataMap: map[string][]byte{
		name: []byte(content),
	}}

	n := len(kv)
	if n%2 != 0 {
		return nil, errors.Errorf("len(kv) %2 != 0, kv must be pair, len: %s, %s", len(kv), kv)
	}

	n = n / 2
	for i := 0; i < n; i++ {
		k := kv[i*2]
		if k == "" {
			return nil, errors.Errorf("key empty, index: %s", i*2)
		}

		v := kv[i*2+1]
		if v == "" {
			return nil, errors.Errorf("value empty, index: %s", i*2+1)
		}

		if _, ok := gos.dataMap[k]; ok {
			return nil, errors.Errorf("duplicate key, index: %s, key: %s", i*2, k)
		}

		gos.dataMap[k] = []byte(v)
	}

	return gos, nil
}
