package main

import (
	_ "dgServer/routers"
	"net/http"
	_ "net/http/pprof"

	"dgServer/controllers"
	db "dgServer/db"

	"github.com/astaxie/beego"
)

func main1() {
	controllers.GetFieldValue()
}

func main() {

	go func() {
		http.ListenAndServe("localhost:7777", nil)
	}()
	beego.InsertFilter("/*", beego.BeforeRouter, controllers.FilterUser)
	beego.BConfig.WebConfig.Session.SessionOn = true
	closeSign := make(chan bool)
	go db.OnRun(closeSign)
	beego.SetStaticPath("/download", "download")
	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/static/vendors", "/verdors")
	beego.Run()

	close(closeSign)
}
