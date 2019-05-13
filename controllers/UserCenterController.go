package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"youdidi/models"
)

type UserCenterController struct {
	beego.Controller
}

// @router /Login/ [GET]
func (this *UserCenterController) Login (){
	this.TplName = "login.html"
}

// @router /Dologin/ [POST]
func (this *UserCenterController) Dologin () {
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	msg := ""

	logs.Notice("user named %s begin to login", inputName)
	logs.Debug("name is %s , passwd is %s" , inputName , inputPasswd)

	var dbUser models.User
	var list []*models.User

	success, num := dbUser.GetUserInfo(inputName, &list)

	if (success != "true") {
		logs.Error("get info of %s fail" , inputName)
		msg = "网络异常，请重试"
	} else {
		if (num == 0) {
			msg = "未注册用户，请先注册"
		} else {
			logs.Debug("get info of %s success; pwd:", inputName)
			if (inputPasswd == list[0].Passwd) {
				msg = "登陆成功"
			} else {
				msg = "密码错误"
			}
		}
	}

	this.Data["userName"] = inputName
	this.Data["msg"] = msg
	this.Data["isMsg"] = "1"
	this.TplName = "login.html"
}