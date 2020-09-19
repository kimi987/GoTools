package fight

import (
	"github.com/nelsonken/cos-go-sdk-v5/cos"
	"bytes"
	"time"
	"fmt"
	"github.com/pkg/errors"
)

const (
	cosConfigPath = "/m7/config/cos"
)

func NewCosUploader(appid, secretid, secretKey, region, bucketName string, timeout time.Duration, prefix string) *CosUploader {
	client := cos.New(appid, secretid, secretKey, region)
	return NewCosClientUploader(client, bucketName, timeout, prefix)
}

func NewCosClientUploader(client *cos.Client, bucketName string, timeout time.Duration, prefix string) *CosUploader {
	bucket := client.Bucket(bucketName)

	return &CosUploader{
		timeout: timeout,
		prefix:  prefix,
		bucket:  bucket,
	}
}

type CosUploader struct {
	client  *cos.Client
	timeout time.Duration

	prefix string

	bucketName string
	bucket     *cos.Bucket
}

var errEmptyData = errors.Errorf("empty data")

func (c *CosUploader) Upload(filename string, data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errEmptyData
	}

	key := fmt.Sprintf("/%s/%s", time.Now().Format("2006-01-02"), filename)
	err := c.bucket.UploadObject(cos.GetTimeoutCtx(c.timeout), key,
		bytes.NewReader(data), &cos.AccessControl{})
	if err != nil {
		return "", errors.Wrapf(err, "upload fail")
	}

	return c.prefix + key, nil
}
