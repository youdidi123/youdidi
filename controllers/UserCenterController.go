package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"youdidi/models"
	"youdidi/redisClient"
)

type UserCenterController struct {
	beego.Controller
}

type UserLoginInfo struct {
	Name string
	IsPhoneVer bool
	IsDriver bool
	Token string
	Phone string
}

var (
	LoginPeriod = 30*60 //用户登陆有效期
	LoginPrefix = "LOGIN_INFO_"
)

// @router /Login/ [GET]
func (this *UserCenterController) Login (){
	this.TplName = "login.html"
}

// @router /Dologin/ [POST,GET]
func (this *UserCenterController) Dologin () {
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	msg := ""
	reUrl := "login.html"

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
				//msg = "登陆成功"
				//reUrl = "index.html"
				token := getToken(inputName , inputPasswd)

				info := &UserLoginInfo{}
				info.Name = inputName
				info.IsPhoneVer = list[0].IsPhoneVer
				info.IsDriver = list[0].IsDriver
				info.Token = token
				info.Phone = list[0].Phone

				data, _ := json.Marshal(info)
				fmt.Println("data: %v", string(data))
				idStr := strconv.FormatInt(list[0].Id,10)

				var cacheClient redisClient.CacheClient
				cacheClient.GetConnet()
				cacheClient.SetKey(LoginPrefix+idStr , string(data))
				cacheClient.Setexpire(LoginPrefix+idStr , LoginPeriod)

				this.Ctx.SetSecureCookie("qyt","qyt_id" , idStr)
				this.Ctx.SetSecureCookie("qyt","qyt_token" , token)

				this.Ctx.Redirect(302, "/Portal/home")

			} else {
				msg = "密码错误"
			}
		}
	}

	this.Data["userName"] = inputName
	this.Data["msg"] = msg
	this.Data["isMsg"] = "1"
	this.TplName = reUrl
}


func getToken(name string , passwd string) string{
	t := time.Now().Unix()
	str := string(t)+name+passwd
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// @router /Ver/phonever [GET]
func (this *UserCenterController) PhoneVer() {
	id, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	this.Data["userId"] = id
	this.TplName = "phoneVer.html"
}