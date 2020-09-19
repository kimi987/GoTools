package kv

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type configReader struct {
	data map[string]string
}

func newConfigReader() *configReader {
	s := &configReader{data: make(map[string]string)}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		if err != os.ErrNotExist {
			logrus.WithError(err).Panic("加载server配置文件失败")
		}

		return s
	}

	err = yaml.Unmarshal(data, s.data)
	if err != nil {
		logrus.WithError(err).Panic("解析server配置文件失败")
	}

	return s
}

func (s *configReader) String(key string) string {
	return s.data[key]
}

func (s *configReader) DefString(key string, def string) string {
	v := s.data[key]
	if len(v) > 0 {
		return v
	}
	return def
}

func (s *configReader) Int(key string) (int, error) {
	v := s.data[key]
	if len(v) > 0 {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.Wrapf(err, "Config.DefInt(%s) fail, %s", key, v)
		}
		return i, nil
	}
	return 0, errors.Errorf("Config.Int(%s) not found", key)
}

func (s *configReader) DefInt(key string, def int) (int, error) {
	v := s.data[key]
	if len(v) > 0 {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.Wrapf(err, "Config.DefInt(%s) fail, %s", key, v)
		}
		return i, nil
	}
	return def, nil
}

func (s *configReader) Bool(key string) (bool, error) {
	v, ok := s.data[key]
	if ok {
		switch v {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New(fmt.Sprintf("bool只允许配true或false. %s: %s", key, v))
		}
	}

	return false, errors.New("没有配置" + key)
}

func (s *configReader) DefBool(key string, def bool) (bool, error) {
	v, ok := s.data[key]
	if ok {
		switch v {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New(fmt.Sprintf("bool只允许配true或false. %s: %s", key, v))
		}
	}

	return def, nil
}
