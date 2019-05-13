package routers

import (
	"youdidi/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Include(&controllers.UserCenterController{})
}
