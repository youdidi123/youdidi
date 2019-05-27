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
	var LoginFilter = func(ctx *context.Context)() {
		id, isId := ctx.GetSecureCookie("qyt","qyt_id")
		if (! isId) {
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
						} else {
							redisClient.Setexpire(controllers.LoginPrefix+id , controllers.LoginPeriod)
						}
					}
				}
			}
		}
	}

	var PhoneVerFilter = func(ctx *context.Context)() {
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

	}

	var ResetInfoFilter = func(ctx *context.Context)() {
		id, _ := ctx.GetSecureCookie("qyt","qyt_id")
		redisClient.Setexpire(controllers.LoginPrefix+id , controllers.LoginPeriod)
	}

	beego.InsertFilter("/Portal/*", beego.BeforeExec, LoginFilter)
	beego.InsertFilter("/Portal/*", beego.BeforeExec, PhoneVerFilter)
	beego.InsertFilter("/Portal/*", beego.AfterExec, ResetInfoFilter)
	beego.InsertFilter("/", beego.BeforeExec, LoginFilter)
	beego.InsertFilter("/", beego.BeforeExec, PhoneVerFilter)


    beego.Router("/wxverifytest", &controllers.WxVerifyTestController{})
	beego.Router("/", &controllers.MainController{}, "GET:Get")
	beego.Include(&controllers.UserCenterController{})
	beego.Include(&controllers.MainController{})
	beego.Include(&controllers.ImgConfirmController{})
	beego.Include(&controllers.OrderController{})


}
