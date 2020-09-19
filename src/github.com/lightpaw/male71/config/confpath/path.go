package confpath

import (
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
)

func GetConfigPath() string {
	s, err := FindConfigPath("conf")
	if err != nil {
		logrus.WithError(err).Panicf("获取配置文件夹失败")
	}
	return s
}

func FindConfigPath(folderName string) (string, error) {

	path0, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	// 防御性，最多100次
	for i := 0; i < 100; i++ {
		confDir := path.Join(path0, folderName)
		if isDirExist(confDir) {
			// 文件夹存在
			return confDir, nil
		}

		parent := path.Dir(path0)
		if parent == path0 {
			return "", errors.Errorf("配置文件夹 %s 没找到", folderName)
		}
		path0 = parent
	}

	return "", errors.Errorf("配置文件夹 %s 没找到（100次都找不到）", folderName)
}

func isDirExist(path string) bool {
	fs, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).Errorf("os.Stat(%s) 出错", path)
		}

		return false
	}

	return fs.IsDir()
}
