package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	// vendor使用 godep save 不要使用 godep save ./...
	// 删除vendor中所有空文件夹

	for {
		removeDirCount := 0

		filepath.Walk("vendor", func(path string, info os.FileInfo, err error) error {

			// 找到所有的空文件夹，删除掉

			// 不是文件夹
			if !info.IsDir() {
				return nil
			}

			if !exists(path) {
				return nil
			}

			fs, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}

			// 有内容
			if len(fs) > 0 {
				return nil
			}

			fmt.Println("删除空文件夹", path)
			err = os.Remove(path)
			if err != nil {
				panic(err)
			}

			removeDirCount++

			return nil
		})

		if removeDirCount <= 0 {
			break
		}
	}

	// 删除vendor中的male7文件夹
	for {
		removeCount := 0

		filepath.Walk("vendor/github.com/lightpaw/male7/", func(path string, info os.FileInfo, err error) error {

			// 如果是文件夹，则判断是不是空的，如果是，则删除，否则下一轮再删除

			if !exists(path) {
				return nil
			}

			if !info.IsDir() {
				// 文件，直接删除
				fmt.Println("删除文件", path)
				err = os.Remove(path)
				if err != nil {
					panic(err)
				}

				removeCount++

				return nil
			}

			fs, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}

			// 有内容
			if len(fs) > 0 {
				return nil
			}

			fmt.Println("删除空文件夹", path)
			err = os.Remove(path)
			if err != nil {
				panic(err)
			}

			removeCount++

			return nil
		})

		if removeCount <= 0 {
			break
		}
	}

	fmt.Println("结束")

}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
