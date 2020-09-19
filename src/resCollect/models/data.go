package models
import (
	"fmt"
	"os"
	"resCollect/conf"
	"strings"
	"sync"
	"time"
)

const (
	// RES_PATH = "/home/samba/temp/hotUpdate/"
	RES_PATH = "D:/waterhunt/client/Assets/AssetBundles/"
)

//GlobalResCol 全局资源收集列表
var GlobalResCol *globalResCol

type globalResCol struct {
	resCheckMap   map[string]int
	reses         []*ResData
	startTime     int64  //当前启动的时间戳
	quality       string //当前品质参数
	platform      string //当前平台
	clientVersion int    //客户端版本
	sync.Mutex
}

type ResData struct {
	ResVer     int //当前资源的版本
	ResName    string
	ResSubPath string
	ResPath    string
	ResTime    int64 //当前时间与初始时间戳的时间差
	ResSize    int   //当前资源的大小
}

//InitData 初始化数据
func InitData() {
	GlobalResCol = &globalResCol{
		resCheckMap: make(map[string]int),
	}
}

//GetReses 获取资源
func (g *globalResCol) GetReses() []*ResData {
	return g.reses
}

//GetClientVerison 获取当前客户端版本
func (g *globalResCol) GetClientVerison() int {
	return g.clientVersion
}

//GetTotalResSize 获取总大小
func (g *globalResCol) GetTotalResSize() int {
	reses := g.reses
	totalSize := 0
	for _, v := range reses {
		if v == nil {
			continue
		}
		totalSize += v.ResSize
	}
	return totalSize
}

//Reset 重置资源
func (g *globalResCol) Reset() {
	g.Lock()
	defer g.Unlock()
	g.resCheckMap = make(map[string]int)
	g.reses = nil
	g.startTime = time.Now().Unix()
	// g.platform = platform
	// g.clientVersion = clientVersion
	// g.clientABVersion = clientABVersion
}

//Reset 重置资源
func (g *globalResCol) ResetWithParams(quality, platform string, clientVersion, isReset int) {
	g.Lock()
	defer g.Unlock()
	if isReset == 0 {
		g.resCheckMap = make(map[string]int)
		g.reses = nil
	}
	g.startTime = time.Now().Unix()
	g.quality = quality
	g.platform = platform
	g.clientVersion = clientVersion
}

//AddRes 添加资源
func (g *globalResCol) AddRes(resVer int, resPath, resName string) {
	g.Lock()
	defer g.Unlock()
	name := fmt.Sprintf("%s/%s", resPath, resName)
	if _, ok := g.resCheckMap[name]; ok {
		return
	}
	// if !strings.Contains(resName, "黑人") {
	// 	//强制使用低模
	// 	resPath = strings.ReplaceAll(resPath, "High/Avatar", "Low/Avatar")
	// 	resPath = strings.ReplaceAll(resPath, "Mid/Avatar", "Low/Avatar")
	// }

	names := strings.Split(resName, "/")

	g.resCheckMap[name] = 1
	resData := &ResData{
		ResVer:     resVer,
		ResName:    names[len(names)-1],
		ResSubPath: strings.Join(names[:len(names)-1], "/"),
		ResPath:    resPath,
		ResTime:    time.Now().Unix() - g.startTime,
	}

	info := g.GetFileInfo(resData)

	if info != nil {
		resData.ResSize = int(info.Size() / 1024)
	}

	g.reses = append(g.reses, resData)
}

//RemoveRes 删除资源
func (g *globalResCol) RemoveRes(path string) bool {
	paths := strings.Split(path, ":")
	if len(paths) < 2 {
		return false
	}
	for n, v := range g.reses {
		if v == nil {
			continue
		}
		if v.ResPath == paths[0] && v.ResName == paths[1] {
			g.reses = append(g.reses[:n], g.reses[n+1:]...)
			break
		}
	}

	return true
}

//GetFileInfo 获取资源信息
func (g *globalResCol) GetFileInfo(d *ResData) os.FileInfo {
	if g.platform == "" {
		return nil
	}

	filePath := fmt.Sprintf("%s%s%s%d/%d/%s/%s", conf.Config.ResPath, g.platform, g.quality, g.clientVersion, d.ResVer, d.ResSubPath, d.ResName)

	info := getFileInfo(fmt.Sprintf("%s.ab", filePath))

	return info
}

//获取文件信息
func getFileInfo(filePath string) os.FileInfo {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		fmt.Println("err = ", err)
	}

	return fileInfo
}

