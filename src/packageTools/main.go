package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

const (
	CodePath    = "Codes"
	PackagePath = "Versions"
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
var wg sync.WaitGroup

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("[Error] 需要参数[2], 当前[1]")
		return
	}
	UpdateCodePath := fmt.Sprintf("%s/%s", CodePath, os.Args[1])
	UpdateVersionPath := fmt.Sprintf("%s/%s", PackagePath, os.Args[1])

	fileBinPath := fmt.Sprintf("%s/%s", UpdateCodePath, "Cores.bin")

	data, err := ioutil.ReadFile(fileBinPath)
	if err != nil {
		fmt.Printf("读取Cores.bin出错: %s\n", err)
		return
	}

	data, err = ungzipBytes(data)
	if err != nil {
		fmt.Printf("解压Cores.bin出错: %s\n", err)
		return
	}

	list := &FileList{}

	if err := list.DecodeMsg(msgp.NewReader(bytes.NewBuffer(data))); err != nil {
		fmt.Printf("decode 出错: %s\n", err)
		return
	}
	targetVersion, _ := strconv.Atoi(os.Args[1])

	os.Remove(UpdateVersionPath)
	os.Mkdir(UpdateVersionPath, os.ModePerm)
	
	for index := 0; index < targetVersion; index++ {
		wg.Add(1)
		go genVerisonFiles(list, UpdateCodePath, UpdateVersionPath, index, targetVersion)
	}
	wg.Wait()
}

//genVerisonFiles 生成版本zip
func genVerisonFiles(list *FileList, fromPath, destPath string, versionStart, versionEnd int) {
	moveZipPath := fmt.Sprintf("%s/Cores_%d.zip", destPath, versionStart)
	zipPath := fmt.Sprintf("Cores_%d.zip", versionStart)
	// destPath = fmt.Sprintf("%s/Core_%d", destPath, versionStart)
	destPath = fmt.Sprintf("Cores_%d", versionStart)
	os.Mkdir(destPath, os.ModePerm)


	savePaths := make(map[string]int)
	for _, folder := range list.Folders {
		for _, file := range folder.Files {
			if file.Version > versionStart {
				savePath := fmt.Sprintf("%s/%d",destPath, file.Version)
				if fi, _ := os.Stat(savePath);fi == nil {
					os.Mkdir(savePath, os.ModePerm)
					savePaths[savePath] = 1

					if(versionEnd == file.Version) {
						_, err := CopyFile(fmt.Sprintf("%s/%s", fromPath, "Cores.bin"), fmt.Sprintf("%s/%s", savePath, "Cores.bin"))

						if err != nil {
							fmt.Printf("CopyFile 出错: %s , path = Cores.bin\n", err)
							return
						}
					}
				}

				//对于当前的版本来说 需要更新的版本
				if strings.Contains(file.Name, ".lua.txt") {
					if folder.Path == "" {
						file.Full_Path = "Lua"
					} else {
						file.Full_Path = fmt.Sprintf("%s/%s", "Lua/", folder.Path)
					}
				}

				if strings.Contains(file.Name, ".ab") && !strings.Contains(file.Name, "shaders") {
					//UI文件
					file.Full_Path  = folder.Path
					// if folder.Path == "" {
					// 	file.Full_Path = "ui"
					// } else {
					// 	file.Full_Path = fmt.Sprintf("%s/%s", "ui/", folder.Path)
					// }
				}

				if file.Full_Path == "" {
					_, err := CopyFile(fmt.Sprintf("%s/%s", fromPath, file.Name), fmt.Sprintf("%s/%s", savePath, file.Name))

					if err != nil {
						fmt.Printf("CopyFile 出错: %s , path = %s\n", err, file.Name)
						return
					}
				} else {

					fi, _ := os.Stat(fmt.Sprintf("%s/%s", savePath, file.Full_Path))
					if fi == nil {
						os.MkdirAll(fmt.Sprintf("%s/%s", savePath, file.Full_Path), os.ModePerm) //创建多级目录
					}

					_, err := CopyFile(fmt.Sprintf("%s/%s/%s", fromPath, file.Full_Path, file.Name), fmt.Sprintf("%s/%s/%s", savePath, file.Full_Path, file.Name))

					if err != nil {
						fmt.Printf("CopyFile 出错: %s , path = %s/%s\n", err, file.Full_Path, file.Name)
						return
					}
				}
			}
		}
	}

	fmt.Println("destPath = ", destPath)
	Zip(destPath, savePaths, zipPath)
	os.RemoveAll(destPath)
	os.Rename(zipPath, moveZipPath)
	wg.Done()
}

//压缩文件夹
func Zip(prefx string, srcFiles map[string]int, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	replacePath := fmt.Sprintf("%s/", prefx)

	for srcFile,_ := range srcFiles {
		filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			path = strings.ReplaceAll(path, "\\", "/")

			path1 := strings.ReplaceAll(path, replacePath, "")

			header.Name = strings.TrimPrefix(path1, filepath.Dir(srcFile)+"/")
			// header.Name = path
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
			}

			return err
		})

		os.RemoveAll(srcFile)
	}
	

	return err
}

//封装的文件拷贝方法
func CopyFile(srcName string, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("CopyFile error = ", err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("CopyFile error = ", err)
		return
	}
	defer dst.Close()
	//拷贝文件
	return io.Copy(dst, src)
}

func ungzipBytes(data []byte) ([]byte, error) {
	rd, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建gzip reader失败 [%s]", err)
	}

	result, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, fmt.Errorf("读取gzip失 [%s]", err)
	}

	return result, nil
}

