package build

import (
	"time"
	"encoding/binary"
	"encoding/hex"
	"math"
	"github.com/lightpaw/male7/config/confpath"
	"path/filepath"
	"os"
	"crypto/md5"
	"io/ioutil"
	"github.com/lightpaw/logrus"
	"strings"
	"strconv"
)

var (
	BuildTime     = ""
	buildUnixTime int64
	GitTag        = "Git Tag No Provided"
	ClientVersion = "-"
	ServerVersion = ""
)

func GetBuildTime() string {
	if t := GetBuildUnixTime(); t != 0 {
		return time.Unix(GetBuildUnixTime(), 0).Format("2006-01-02_15:04:05")
	} else {
		return time.Now().Format("2006-01-02_15:04:05")
	}
}

func GetBuildUnixTime() int64 {
	if buildUnixTime == 0 && len(BuildTime) > 0 {
		if ut, err := strconv.ParseInt(BuildTime, 10, 64); err != nil {
			logrus.Error("解析BuildTime失败，" + BuildTime)
		} else {
			buildUnixTime = ut
		}
	}

	return buildUnixTime
}

func GetClientVersion() string {
	return ClientVersion
}

func GetVersion() string {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(GetBuildUnixTime()&math.MaxUint16))
	return hex.EncodeToString(b)
}

func GetConfigVersion() string {
	// 读取配置的文件夹，每个文件都加载上来，计算

	s, err := confpath.FindConfigPath("conf")
	if err != nil {
		logrus.WithError(err).Error("生成配置版本号，获取配置路径失败")
		return "ffff"
	}

	hash := md5.New()
	err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			// 隐藏文件，跳过
			return nil
		}

		if strings.Contains(path, ".svn") {
			// 跳过svn内容
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if len(b) > 0 {
			hash.Write(b)
		}
		return nil
	})

	if err != nil {
		logrus.WithError(err).Error("生成配置版本号，读取配置失败")
		return "ffff"
	}

	sum := hash.Sum(nil)
	for i := 1; i < len(sum)/2; i++ {
		sum[0] ^= sum[i*2]
		sum[1] ^= sum[i*2+1]
	}

	return hex.EncodeToString(sum[:2])
}
