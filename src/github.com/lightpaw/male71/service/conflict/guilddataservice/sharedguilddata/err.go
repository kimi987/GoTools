package sharedguilddata

import "github.com/pkg/errors"

var (
	ErrNotExist = errors.Errorf("帮派不存在")
	ErrTimeout  = errors.Errorf("超时")
)
