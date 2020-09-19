package main

import (
	"github.com/astaxie/beego"
	"resCollect/conf"
	"resCollect/models"
	_ "resCollect/routers"
)

func main() {
	conf.InitConfig()
	models.InitData()

	beego.SetStaticPath("/download", "download")
	beego.SetStaticPath("/static", "views/static")
	beego.SetStaticPath("/static/vendors", "views/verdors")

	beego.Run()
}
