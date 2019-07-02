package routers

import (
	"github.com/astaxie/beego/context"
	"encoding/json"
	"youdidi/controllers"
	"github.com/astaxie/beego"
	"youdidi/redisClient"
	"github.com/astaxie/beego/logs"
)

func init() {
	var AdminLoginFilter = func(ctx *context.Context)() {
		id, isId := ctx.GetSecureCookie("qyt","qyt_admin_id")
		if (! isId) {
			ctx.Redirect(302, "/AdminLogin/")
			logs.Debug("can not get id from cookie")
		} else {
			logs.Debug("id of cookis : %v" , id)
			token , isToken := ctx.GetSecureCookie("qyt" , "qyt_admin_token")
			if (! isToken) {
				logs.Debug("can not get token from cookie")
				ctx.Redirect(302, "/AdminLogin")
			} else {
				content := redisClient.GetKey(controllers.AdminLoginPrefix+id)
				if (content == "nil") {
					logs.Debug("cache is empty")
					ctx.Redirect(302, "/AdminLogin")
				} else {
					info := &controllers.AdminUserLoginInfo{}
					err := json.Unmarshal([]byte(content), &info)
					if (err != nil) {

						ctx.Redirect(302, "/AdminLogin")
					} else {
						if (token != info.Token) {
							logs.Debug("token did not match of cookie and cache")
							ctx.Redirect(302, "/AdminLogin")
						} else {
							redisClient.Setexpire(controllers.AdminLoginPrefix+id , controllers.AdminLoginPeriod)
						}
					}
				}
			}
		}
	}

	var LoginFilter = func(ctx *context.Context)() {
		id, isId := ctx.GetSecureCookie("qyt","qyt_id")
		logs.Debug("debug cookie id=%v isId=%v",id, isId)
		if (! isId || id == "") {
			ctx.Redirect(302, "/Login")
			logs.Debug("can not get id from cookie")
		} else {
			logs.Debug("id of cookis : %v" , id)
			token , isToken := ctx.GetSecureCookie("qyt" , "qyt_token")
			if (! isToken) {
				logs.Debug("can not get token from cookie")
				ctx.Redirect(302, "/Login")
			} else {
				content := redisClient.GetKey(controllers.LoginPrefix+id)
				if (content == "nil") {
					logs.Debug("cache is empty")
					ctx.Redirect(302, "/Login")
				} else {
					info := &controllers.UserLoginInfo{}
					err := json.Unmarshal([]byte(content), &info)
					if (err != nil) {
						ctx.Redirect(302, "/Login")
					} else {
						if (token != info.Token) {
							logs.Debug("token did not match of cookie and cache")
							ctx.Redirect(302, "/Login")
						} else if(! info.IsPhoneVer){
							ctx.Redirect(302, "/Ver/phonever")
						}else {
							redisClient.Setexpire(controllers.LoginPrefix+id , controllers.LoginPeriod)
						}
					}
				}
			}
		}
	}

	/*var PhoneVerFilter = func(ctx *context.Context)() {
		logs.Debug("into phonever filter")
		id, isId := ctx.GetSecureCookie("qyt","qyt_id")
		if (! isId) {
			ctx.Redirect(302, "/Login")
			logs.Debug("can not get id from cookie")
		}

		content := redisClient.GetKey(controllers.LoginPrefix+id)
		if (content == "nil") {
			logs.Debug("cache is empty")
			ctx.Redirect(302, "/Login")
		}
		info := &controllers.UserLoginInfo{}
		err := json.Unmarshal([]byte(content), &info)
		if (err != nil) {
			ctx.Redirect(302, "/Login")
		}
		if (! info.IsPhoneVer) {
			ctx.Redirect(302, "/Ver/phonever")
		}

	}*/

	var ResetInfoFilter = func(ctx *context.Context)() {
		id, _ := ctx.GetSecureCookie("qyt","qyt_id")
		redisClient.Setexpire(controllers.LoginPrefix+id , controllers.LoginPeriod)
	}

	beego.InsertFilter("/admin/*", beego.BeforeExec, AdminLoginFilter)
	beego.InsertFilter("/portal/*", beego.BeforeExec, LoginFilter)
	//beego.InsertFilter("/portal/*", beego.BeforeExec, PhoneVerFilter)
	beego.InsertFilter("/portal/*", beego.AfterExec, ResetInfoFilter)
	beego.InsertFilter("/", beego.BeforeExec, LoginFilter)
	//beego.InsertFilter("/", beego.BeforeExec, PhoneVerFilter)


    beego.Router("/wxverifytest", &controllers.WxVerifyTestController{})
    //beego.Router("/wxlogin", &controllers.WxLoginController{})
	beego.Router("/", &controllers.OrderController{}, "GET:SearchInput")
	beego.Include(&controllers.WxLoginController{})
	beego.Include(&controllers.UserCenterController{})
	beego.Include(&controllers.MainController{})
	beego.Include(&controllers.ImgConfirmController{})
	beego.Include(&controllers.OrderController{})
	beego.Include(&controllers.ChatController{})
	beego.Include(&controllers.AccountFlowController{})
	beego.Include(&controllers.AdminUserController{})
}
