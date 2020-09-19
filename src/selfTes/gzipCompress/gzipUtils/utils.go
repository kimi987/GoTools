package gzipCompress

import (
	"os"
	"io/ioutil"
)

//获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
    dir, err := ioutil.ReadDir(dirPth)
    if err != nil {
        return nil, nil, err
    }
    PthSep := string(os.PathSeparator)


    return files, dirs, nil
}
