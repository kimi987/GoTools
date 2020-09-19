package util

import (
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"os"
	"io/ioutil"
	"runtime"
)

func DumpStacks(name string) {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	ioutil.WriteFile(name+"_"+time.Now().Format(timeutil.SecondsLayout), buf, os.ModePerm)
}
