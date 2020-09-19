package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	sourceDir = flag.String("source", ".", "source dir to scan")
	gzipDir   = flag.String("gzip", "./gzip", "output gzip dir")
	normalDir = flag.String("out", "./out", "normal dir")

	sizeLimit = flag.Int("limit", 2048, "size limit to gzip")
)

func main() {
	flag.Parse()

	os.RemoveAll(*gzipDir)
	os.RemoveAll(*normalDir)

	// 先获取所有文件列表, 防止生成的文件也处理了
	files := make([]string, 0, 100)

	err := filepath.Walk(*sourceDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	os.MkdirAll(*gzipDir, os.ModePerm)
	os.MkdirAll(*normalDir, os.ModePerm)

	cpuCount := runtime.NumCPU()

	fileChan := make(chan string, len(files))

	for _, f := range files {
		fileChan <- f
	}

	close(fileChan)
	latch := &sync.WaitGroup{}
	latch.Add(len(files))

	for i := 0; i < cpuCount; i++ {
		go func() {
			f, ok := <-fileChan
			if !ok {
				return
			}

			processFile(f)
			latch.Done()
		}()
	}

	latch.Wait()

	fmt.Println("done")
}
func processFile(f string) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Printf("读取文件错误%s: %s\n", f, err)
		return
	}

	base := path.Dir(f)

	if len(content) <= *sizeLimit {
		dir := path.Join(*normalDir, base)
		os.MkdirAll(dir, os.ModePerm)
		final := path.Join(dir, f)
		fmt.Println(final)
		if err := ioutil.WriteFile(f, content, os.ModePerm); err != nil {
			fmt.Printf("写入最终文件失败%s: %s\n", final, err)
		}

		return
	}

	dir := path.Join(*gzipDir, base)
	os.MkdirAll(dir, os.ModePerm)
	final := path.Join(dir, f)
	fmt.Println(final)
	finalFile, err := os.Create(final)
	if err != nil {
		fmt.Printf("打开压缩文件失败%s: %s\n", final, err)
		return
	}
	defer finalFile.Close()

	gzipS := gzip.NewWriter(finalFile)
	defer gzipS.Close()

	if _, err := gzipS.Write(content); err != nil {
		fmt.Printf("写入压缩文件失败%s: %s\n", finalFile, err)
		return
	}
}

