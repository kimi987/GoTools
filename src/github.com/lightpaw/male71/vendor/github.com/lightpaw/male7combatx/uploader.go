package combatx

import (
	"github.com/nelsonken/cos-go-sdk-v5/cos"
	"time"
	"fmt"
	"bytes"
	"github.com/pkg/errors"
	"path/filepath"
	"os"
	"strings"
	"io/ioutil"
	"path"
	"github.com/lightpaw/logrus"
	"runtime/debug"
)

type Uploader interface {
	Upload(filename string, data []byte) (link, secondLink string, err error)
}

type UploaderFunc func(filename string, data []byte) (link, secondLink string, err error)

func (f UploaderFunc) Upload(filename string, data []byte) (link, secondLink string, err error) {
	return f(filename, data)
}

func NewLocalUploader(dir string) Uploader {
	linkPrefix := NewLinkPrefix(LocalPrefix)
	return UploaderFunc(func(filename string, data []byte) (path, secondLink string, err error) {
		err = writeLocalFile(dir, filename, data)
		if err != nil {
			return
		}
		path = linkPrefix.String(filename)
		return
	})
}

const (
	CosPrefix   = "{{cos}}"
	LocalPrefix = "{{local}}"

	cosDir = "cos"
)

func NewLinkPrefix(prefix string) *LinkPrefix {
	var lp [2]string
	if strings.HasSuffix(prefix, "/") {
		lp[0] = prefix[:len(prefix)-1]
		lp[1] = prefix
	} else {
		lp[0] = prefix
		lp[1] = prefix + "/"
	}

	linkPrefix := LinkPrefix(lp)
	return &linkPrefix
}

type LinkPrefix [2]string // [0]=prefix [1]=prefix/

func (lp *LinkPrefix) String(filename string) string {
	if strings.HasPrefix(filename, "/") {
		return lp[0] + filename
	} else {
		return lp[1] + filename
	}
}

func writeLocalFile(dir, filename string, data []byte) error {

	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	// marshal数据，保存在本地
	// 没找到好的uuid库，先临时使用时间戳作为文件名
	//filename := fmt.Sprintf("%v.txt", time.Now().UnixNano())
	err = ioutil.WriteFile(filepath.Join(dir, filename), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func NewCosUploader(appid, secretid, secretKey, region, bucketName string, timeout time.Duration, prefix, localDir, localLinkPrefix string) *CosUploader {
	client := cos.New(appid, secretid, secretKey, region)
	return NewCosClientUploader(client, bucketName, timeout, prefix, localDir, localLinkPrefix)
}

func NewCosClientUploader(client *cos.Client, bucketName string, timeout time.Duration, prefix, localDir, localLinkPrefix string) *CosUploader {
	bucket := client.Bucket(bucketName)

	cu := &CosUploader{
		timeout: timeout,
		prefix:  prefix,
		bucket:  bucket,
	}

	if len(localDir) > 0 && len(localLinkPrefix) > 0 {
		cu.localDir = localDir
		cu.localCosDir = path.Join(localDir, cosDir)
		cu.localLinkPrefix = NewLinkPrefix(localLinkPrefix)

		cu.localUploadTicker = time.NewTicker(10 * time.Minute)

		go catchPanic(cu.loop, "CosUploader.loop()")
	}

	return cu
}

type CosUploader struct {
	client  *cos.Client
	timeout time.Duration

	prefix string

	localDir          string
	localCosDir       string
	localLinkPrefix   *LinkPrefix
	localUploadTicker *time.Ticker

	bucketName string
	bucket     *cos.Bucket
}

var errEmptyData = errors.Errorf("empty data")

func (c *CosUploader) Upload(filename string, data []byte) (string, string, error) {
	if len(data) <= 0 {
		return "", "", errEmptyData
	}

	key := fmt.Sprintf("/%s/%s", time.Now().Format("2006-01-02"), filename)
	err := c.bucket.UploadObject(cos.GetTimeoutCtx(c.timeout), key,
		bytes.NewReader(data), &cos.AccessControl{})
	if err != nil {
		if len(c.localDir) <= 0 {
			return "", "", errors.Wrapf(err, "upload cos fail")
		}

		// 上传失败，本地处理
		filename := cosDir + key
		err := writeLocalFile(c.localDir, filename, data)
		if err != nil {
			return "", "", errors.Wrapf(err, "upload local fail")
		}

		link := c.localLinkPrefix.String(filename)

		return CosPrefix + key, link, nil
	}

	return CosPrefix + key, "", nil
}

func (c *CosUploader) Close() {
	if c.localUploadTicker != nil {
		c.localUploadTicker.Stop()
	}
}

func (c *CosUploader) loop() {
	for range c.localUploadTicker.C {
		go catchPanic(c.scanUploadFile, "CosUploader.scanUploadFile()")
	}
}

func (c *CosUploader) scanUploadFile() {

	// 遍历文件夹下的文件，然后一个个文件传上去
	filepath.Walk(c.localCosDir, func(path0 string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		key := strings.Replace(path0, "\\", "/", -1) // 更新windows支持
		key = strings.Replace(key, c.localCosDir, "", 1)

		if err := c.tryUploadLocalFile(path0, key); err != nil {
			logrus.WithError(err).Error("上传本地战斗回放到Cos失败")
		}

		return nil
	})

}

func (c *CosUploader) tryUploadLocalFile(fullFilePath, key string) error {

	// 通过fullPath读取文件内容
	data, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		return errors.Wrap(err, "tryUploadLocalFile read file fail")
	}

	// 将文件内容上传到Cos
	err = c.bucket.UploadObject(cos.GetTimeoutCtx(c.timeout), key,
		bytes.NewReader(data), &cos.AccessControl{})
	if err != nil {
		return errors.Wrap(err, "tryUploadLocalFile upload cos fail")
	}

	// 删除fullPath文件
	if err := os.RemoveAll(fullFilePath); err != nil {
		return errors.Wrap(err, "tryUploadLocalFile delete file fail")
	}
	return nil
}

func catchPanic(f func(), name string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", name)
		}
	}()

	f()
}
