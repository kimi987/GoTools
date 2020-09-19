package routers

import (
	"resCollect/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.ResourceShowController{})
	beego.Router("/ResourceShow", &controllers.ResourceShowController{})
	beego.Router("/ResourceCol", &controllers.ResourceController{})
}
