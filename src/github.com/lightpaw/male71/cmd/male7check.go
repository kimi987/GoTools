package main

import (
	"fmt"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/build"
	conf "github.com/lightpaw/male7/config"
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/confpath"
)

func main() {

	// 初始化配置
	generateI18n()

	service.Init()
	fmt.Println("服务器配置检查通过，版本号：", build.GetConfigVersion())
}

func generateI18n() {
	basePath := confpath.GetConfigPath()
	gos, err := config.NewConfigGameObjects(basePath)
	if err != nil {
		logrus.WithError(err).Panic("加载配置文件失败")
	}

	c, err := conf.ParseConfigDatas(gos)
	if err != nil {
		logrus.WithError(err).Panic("加载配置文件失败")
	}

	for _, n := range c.I18nData().Array {
		n.Generate(gos, basePath)
	}
}
