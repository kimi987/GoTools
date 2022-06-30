package routers

import (
	"dgServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/", &controllers.MainController{})
	beego.Router("/UsersStatus", &controllers.UsersStatus{})
	beego.Router("/TaskStatus", &controllers.TaskStatus{})
	beego.Router("/TaskScore", &controllers.TaskScore{})
	beego.Router("/PayInfo", &controllers.PayInfo{})
	beego.Router("/PayRate", &controllers.PayRate{})
	beego.Router("/TeacherStatus", &controllers.TeacherInfo{})
	beego.Router("/UserDetail", &controllers.UserDetail{})
}
