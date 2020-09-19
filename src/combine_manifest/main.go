//go:generate msgp
package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
)

const (
	MANIFEST_SIZE_LIMIT = 256 * 256 * 128
)

var (
	outDir = flag.String("out", "_out", "out dir")

	compressAb = flag.Bool("compress", true, "compress asset bundle")
)

type FileList struct {
	Version int `msg:"version"`

	Folders []*Folder `msg:"folders"`
}

type Folder struct {
	Path string `msg:"path"`

	Files []*SingleFle `msg:"files"`
}

type SingleFle struct {
	Ignore bool `msg:"ignore"`

	Name string `msg:"name"`

	Version int `msg:"version"`

	File_length int `msg:"file_length"`

	Full_Path string
}

func main() {
	flag.Parse()

	dir, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("获取当前路径出错: %s\n", err)
		return
	}

	currentVersion, err := strconv.Atoi(filepath.Base(dir))
	if err != nil {
		fmt.Printf("当前路径竟然不是数字的版本号: %s\n", filepath.Base(dir))
		return
	}

	os.RemoveAll(*outDir)

	data, err := ioutil.ReadFile("files.bin")
	if err != nil {
		fmt.Printf("读取files.bin出错: %s\n", err)
		return
	}

	data, err = ungzipBytes(data)
	if err != nil {
		fmt.Printf("解压files.bin出错: %s\n", err)
		return
	}

	list := &FileList{}

	if err := list.DecodeMsg(msgp.NewReader(bytes.NewBuffer(data))); err != nil {
		fmt.Printf("decode 出错: %s\n", err)
		return
	}

	hasError := false

	files := make(map[string]*SingleFle)

	for _, folder := range list.Folders {
		for _, file := range folder.Files {
			var fullPath string
			if folder.Path == "" {
				fullPath = file.Name
			} else {
				fullPath = folder.Path + "/" + file.Name
			}

			file.Full_Path = fullPath

			if _, has := files[fullPath]; has {
				fmt.Printf("文件重复: %s\n", fullPath)
				hasError = true
			} else {
				files[fullPath] = file
			}
		}
	}

	// 检查files里是不是文件都存在
	for path, file := range files {
		if file.Version > currentVersion {
			fmt.Printf("竟然有文件的版本号大于当前的版本号: %d, %s\n", file.Version, path)
			hasError = true
			continue
		}

		realPath := filepath.Join("..", strconv.Itoa(file.Version), path)

		info, err := os.Stat(realPath)
		if err != nil {
			fmt.Printf("stat err. 文件不存在? : %s -> %s\n", realPath, err)
			hasError = true
			continue
		}

		if info.IsDir() {
			fmt.Printf("竟然有文件夹在filelist里: %s\n", path)
			hasError = true
			continue
		}

		// 限制manifest的大小为8m, 可以用3个byte表示长度, 并且第一个bit表示是否有压缩
	}

	// 检查ab和manifest是不是都一一对应, 且版本号相同
	for path, file := range files {
		ext := filepath.Ext(path)
		if ext == ".ab" {
			manifestPath := path + ".manifest"

			if manifestFile, has := files[manifestPath]; !has {
				fmt.Printf("竟然有.ab, 没有对应的.manifest: %s\n", path)
				hasError = true
			} else {
				if manifestFile.Version != file.Version {
					fmt.Printf("竟然有.ab和.manifest的版本号不匹配: %d-%d, %s\n", file.Version, manifestFile.Version, path)
					hasError = true
				}
			}
		} else if ext == ".manifest" {
			abPath := path[0 : len(path)-len(ext)]
			if _, has := files[abPath]; !has {
				fmt.Printf("竟然有.manifest, 没有对应的.ab: %s\n", path)
				hasError = true
			}

			if file.File_length > MANIFEST_SIZE_LIMIT {
				fmt.Printf("manifest文件超过8m大小上限: %s\n", path)
				hasError = true
			}
		}
	}

	// 检查实际文件夹的文件都在filelist中, 且都是当前版本号
	whilelist := map[string]bool{
		"Cores.bin": true,
		"Cores.txt": true,
		"Cores.zip": true,
		"files.bin": true,
		"files.txt": true,
	}

	for index := 0; index < currentVersion; index++ {
		whilelist[fmt.Sprintf("Cores_%d.zip", index)] = true
	}

	// 先获取所有文件列表, 防止生成的文件也处理了
	actualFiles := make([]*filepathAndSingleFile, 0, 100)

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			if file, has := files[path]; !has {
				if _, has := whilelist[path]; !has {
					if base := filepath.Base(path); base != ".DS_Store" {
						fmt.Printf("包含filelist中没有的文件: %s\n", path)
						hasError = true
					}
				}

				actualFiles = append(actualFiles, &filepathAndSingleFile{path: path})
			} else {
				if file.Version != currentVersion {
					fmt.Printf("当前文件夹下, 竟然包含一个不是这个版本的文件: %d: %s\n", file.Version, path)
					hasError = true
				}
				actualFiles = append(actualFiles, &filepathAndSingleFile{path: path, singleFile: file})
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("遍历出错: %s\n", err)
		hasError = true
	}

	// 删掉filelist中的manifest

	fileCount := 0
	for _, folder := range list.Folders {
		newFiles := make([]*SingleFle, 0, len(folder.Files)/2)

		for _, file := range folder.Files {
			if ext := filepath.Ext(file.Name); ext != ".manifest" {
				newFiles = append(newFiles, file)
				fileCount++

				if file.Version < currentVersion {
					file.File_length = getOldVersionFileSize(file.Version, file.Full_Path)
				}
			}
		}

		folder.Files = newFiles
	}

	if fileCount*2 != len(files) {
		fmt.Printf("竟然新文件数量不是老文件数量的一半. 新 %d, 老%d", fileCount, len(files))
		hasError = true
	}

	if hasError {
		fmt.Printf("Abort")

		return
	}

	// 开始实际把ab和manifest都合并起来

	os.MkdirAll(*outDir, os.ModePerm)

	cpuCount := runtime.NumCPU()

	fileChan := make(chan *filepathAndSingleFile, len(actualFiles))

	for _, f := range actualFiles {
		fileChan <- f
	}

	close(fileChan)
	latch := &sync.WaitGroup{}
	latch.Add(len(actualFiles))

	for i := 0; i < cpuCount; i++ {
		go func() {
			for {
				f, ok := <-fileChan
				if !ok {
					return
				}

				processFile(f)
				latch.Done()
			}
		}()
	}

	latch.Wait()

	// 检查所有文件都有大小
	for _, folder := range list.Folders {
		for _, f := range folder.Files {
			if f.File_length == 0 {
				fmt.Printf("竟然有文件没有设置大小: %v\n", f)
			}
		}
	}

	// write final files.bin
	if err := writeFileBin(list); err != nil {
		fmt.Printf("写入files.bin失败: %s\n", err)
		return
	}

	// 进入最终目录检查 TODO

	fmt.Println("done")
}

type filepathAndSingleFile struct {
	path       string
	singleFile *SingleFle
}

func writeFileBin(list *FileList) error {
	binPath := filepath.Join(*outDir, "files.bin")

	fos, err := os.Create(binPath)
	if err != nil {
		return errors.Wrap(err, "打开最终files.bin失败")
	}

	gos := gzip.NewWriter(fos)

	msgpWriter := msgp.NewWriter(gos)

	if err := list.EncodeMsg(msgpWriter); err != nil {
		return errors.Wrap(err, "encode最终files.bin出错")
	}

	if err := msgpWriter.Flush(); err != nil {
		return errors.Wrap(err, "flush最终files.bin出错")
	}

	if err := gos.Close(); err != nil {
		return errors.Wrap(err, "close gzip出错")
	}

	if err := fos.Close(); err != nil {
		return errors.Wrap(err, "close fos出错")
	}

	return nil
}

var (
	cacheFileList map[int]*FileListInfo = make(map[int]*FileListInfo, 10)
)

type FileListInfo struct {
	fileList *FileList
	content  map[string]*SingleFle
}

func doReadFileList(version int) (*FileListInfo, error) {
	path := filepath.Join("..", strconv.Itoa(version), "_out", "files.bin")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("读取files.bin出错: %s\n", err)
		return nil, errors.Wrap(err, "files.bin不存在")
	}

	data, err = ungzipBytes(data)
	if err != nil {
		fmt.Printf("解压files.bin出错: %s\n", err)
		return nil, errors.Wrap(err, "解压files.bin出错")
	}

	list := &FileList{}

	if err := list.DecodeMsg(msgp.NewReader(bytes.NewBuffer(data))); err != nil {
		fmt.Printf("decode 出错: %s\n", err)
		return nil, errors.Wrap(err, "decode files.bin出错")
	}

	files := make(map[string]*SingleFle)

	for _, folder := range list.Folders {
		for _, file := range folder.Files {
			var fullPath string
			if folder.Path == "" {
				fullPath = file.Name
			} else {
				fullPath = folder.Path + "/" + file.Name
			}

			if _, has := files[fullPath]; has {
				fmt.Printf("文件重复: %s\n", fullPath)
			} else {
				files[fullPath] = file
			}
		}
	}

	return &FileListInfo{fileList: list, content: files}, nil
}

func getOldVersionFileSize(version int, path string) int {
	list, has := cacheFileList[version]
	if !has {
		var err error
		list, err = doReadFileList(version)
		if err != nil {
			fmt.Printf("获取文件大小出错: %s\n", err)
			return 0
		}

		cacheFileList[version] = list
	}

	singleFile, has := list.content[path]
	if !has {
		fmt.Printf("文件在%d版本里没找到: %s\n", version, path)
		return 0
	}

	return singleFile.File_length
}

func processFile(pathAndFile *filepathAndSingleFile) (finalErr error) {
	path := pathAndFile.path

	ext := filepath.Ext(path)
	if ext == ".manifest" { // 这个文件特殊写
		// 如果是manifest, 则略过
		return
	}

	dir := filepath.Dir(path)
	fileName := filepath.Base(path)
	finalDir := filepath.Join(*outDir, dir)
	os.MkdirAll(finalDir, os.ModePerm)
	final := filepath.Join(finalDir, fileName)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("读取文件错误%s: %s\n", path, err)
		finalErr = errors.New("")
	}

	// 如果是ab, 找到对应的manifest, 合并
	if ext == ".ab" || path == "AssetBundles" {
		manifestPath := path + ".manifest"

		manifestContent, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			fmt.Printf("读取manifest文件错误: %s\n", err)
			finalErr = errors.New("")
			return
		}

		manifestLen := len(manifestContent)

		isManifestGzipped := false
		isAbGzipped := false

		if manifestLen >= 2048 {
			// 压缩manifest
			if gzipped, err := gzipBytes(manifestContent); err != nil {
				fmt.Printf("压缩 manifest出错: %s\n", err)
				finalErr = err
				return
			} else {
				isManifestGzipped = true
				manifestContent = gzipped
				manifestLen = len(manifestContent)
			}
		}

		abContent := content

		if *compressAb && len(content) > 4096 {
			compressed, err := gzipBytes(content)
			if err != nil {
				fmt.Printf("压缩ab出错: %s\n", err)
				finalErr = err
				return
			}

			if len(content)-len(compressed) > 2048 {
				// 只有压缩的效果超过2k才压缩

				isAbGzipped = true
				abContent = compressed
			}
		}

		finalLen := len(abContent) + manifestLen + 3
		finalContent := make([]byte, 3, finalLen)

		finalContent[0] = byte((manifestLen >> 16) & 0xff)
		finalContent[1] = byte((manifestLen >> 8) & 0xff)
		finalContent[2] = byte(manifestLen & 0xff)

		if isManifestGzipped {
			finalContent[0] = finalContent[0] | 0x80 // set first bit
		}

		if isAbGzipped {
			finalContent[0] = finalContent[0] | 0x40 // set second bit
		}

		finalContent = append(finalContent, manifestContent...)
		finalContent = append(finalContent, abContent...)

		// 检查一下

		if len(finalContent) != finalLen {
			fmt.Printf("最终len不对, 应该是 %d, 实际是 %d\n", finalLen, len(finalContent))
			finalErr = errors.New("")
		}

		isFinalGzipped := false
		isFinalAbGzipped := false
		if finalContent[0]&0x80 != 0 {
			isFinalGzipped = true
		}

		if finalContent[0]&0x40 != 0 {
			isFinalAbGzipped = true
		}

		if isFinalGzipped != isManifestGzipped {
			fmt.Printf("解析出来判断是否压缩不对\n")
			finalErr = errors.New("")
		}

		if isFinalAbGzipped != isAbGzipped {
			fmt.Printf("解析出来判断ab是否压缩不对\n")
			finalErr = errors.New("")
		}

		size := (int(finalContent[0]&0x3f) << 16) | (int(finalContent[1]) << 8) | int(finalContent[2])
		if size != manifestLen {
			fmt.Printf("解析出来manifestLen不对, 应该是 %d, 实际是  %d\n", manifestLen, size)
			finalErr = errors.New("")
		}

		if !Equal(manifestContent, finalContent[3:3+manifestLen]) {
			fmt.Printf("解析出来manifest和原来的不同\n")
			finalErr = errors.New("")
		}

		if !isAbGzipped {
			if !Equal(content, finalContent[3+manifestLen:]) {
				fmt.Printf("解析出来ab content和原来的不同\n")
				finalErr = errors.New("")
			}
		} else {
			unzippedAb, err := ungzipBytes(finalContent[3+manifestLen:])
			if err != nil {
				fmt.Printf("竟然解压ab出错: %s\n")
				finalErr = err
			} else {
				if !Equal(unzippedAb, content) {
					fmt.Printf("解析并解压出来ab content和原来不同\n")
					finalErr = errors.New("")
				}
			}
		}

		content = finalContent
	}

	if pathAndFile.singleFile != nil {
		pathAndFile.singleFile.File_length = len(content)
	}
	// 其他的直接copy到out
	if err := ioutil.WriteFile(final, content, os.ModePerm); err != nil {
		fmt.Printf("写入最终文件失败%s: %s\n", final, err)
		finalErr = errors.New("")
	}
	return
}

func ungzipBytes(data []byte) ([]byte, error) {
	rd, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "创建gzip reader失败")
	}

	result, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, errors.Wrap(err, "读取gzip失败")
	}

	return result, nil
}

func gzipBytes(data []byte) ([]byte, error) {
	final := &bytes.Buffer{}
	wr := gzip.NewWriter(final)

	if n, err := wr.Write(data); n != len(data) || err != nil {
		return nil, errors.Wrap(err, "写writer出错")
	}

	if err := wr.Close(); err != nil {
		return nil, errors.Wrap(err, "close wr出错")
	}

	return final.Bytes(), nil
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
