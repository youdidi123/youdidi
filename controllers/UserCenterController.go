package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"youdidi/models"
	"youdidi/redisClient"
)

type UserCenterController struct {
	beego.Controller
}

type UserLoginInfo struct {
	Name       string
	Nickname   string
	OpenId     string
	IsPhoneVer bool
	IsDriver   int
	OrderNumWV int
	Token      string
	Phone      string
}

var (
	LoginPeriod    = 30 * 60        //用户登陆有效期
	LoginPrefix    = "LOGIN_INFO_"  //缓存在redis中的用户数据key前缀
	PhoneVerPrefix = "RANDOM_CODE_" //短信验证码key前缀
)

//登陆页首页，后续接入微信用户体系后废弃
// @router /Login/ [GET]
func (this *UserCenterController) Login() {
	this.TplName = "login.html"
}

//执行真正的登陆操作，接入微信后需要改造
// @router /Dologin/ [POST,GET]
func (this *UserCenterController) Dologin() {
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	msg := ""
	reUrl := "login.html"

	logs.Notice("user named %s begin to login", inputName)
	logs.Debug("name is %s , passwd is %s", inputName, inputPasswd)

	var dbUser models.User
	var list []*models.User

	success, num := dbUser.GetUserInfo(inputName, &list)

	if success != "true" {
		logs.Error("get info of %s fail", inputName)
		msg = "网络异常，请重试"
	} else {
		if num == 0 {
			msg = "未注册用户，请先注册"
		} else {
			logs.Debug("get info of %s success; pwd:", inputName)
			if inputPasswd == list[0].Passwd {
				//msg = "登陆成功"
				//reUrl = "index.html"

				token, idStr, err := CacheUserLoginInfo(list[0])
				if (err != nil) {
					logs.Warn("Cache user login info failed!")
				}
				idStr = strconv.Itoa(list[0].Id)

				this.Ctx.SetSecureCookie("qyt", "qyt_id", idStr) //注入用户id，后续所有用户id都从cookie里获取
				this.Ctx.SetSecureCookie("qyt", "qyt_token", token)
				//this.SetSession("qyt_id" , idStr)

				this.Ctx.Redirect(302, "/Portal/showdriverorder/")

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


//执行真正的注册操作，非微信入口的注册
// @router /Dologon/ [POST,GET]
func (this *UserCenterController) Dologon() {

}

	//使用用户名，密码，时间戳生成用户的鉴权token
// 用户cookie和服务redis里都需要存储
func getToken(name string, passwd string) string {
	t := time.Now().Unix()
	str := string(t) + name + passwd
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//加载用户绑定手机页面
// @router /Ver/phonever [GET]
func (this *UserCenterController) PhoneVer() {
	id, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	//id := this.GetSession("qyt_id")
	this.Data["userId"] = id
	this.TplName = "phoneVer.html"
}

//获取短信验证码
// @router /Ver/getvercode [GET,POST]
func (this *UserCenterController) GetVerCode() {
	expireMin := 5 //验证码有效时间
	userId := this.GetString("userId")
	phoneNum := this.GetString("phoneNum")

	baseUrl := beego.AppConfig.String("smsBaseUrl")
	//connectTimeout , _ := strconv.Atoi(beego.AppConfig.String("connectTimeout"))
	//readWriteTimeout , _ := strconv.Atoi(beego.AppConfig.String("readWriteTimeout"))

	//验证码发送平台需要提供的各种token
	//https://www.yuntongxun.com/member/main
	accountSid := beego.AppConfig.String("smsAccountSid")
	token := beego.AppConfig.String("smsToken")
	appId := beego.AppConfig.String("smsAppId")

	sig, auth := getSig(accountSid, token)
	randomCode := GetRandomCode()

	redisClient.SetKey(PhoneVerPrefix+userId, randomCode)
	redisClient.Setexpire(PhoneVerPrefix+userId, expireMin*60)

	baseUrl = baseUrl + "/2013-12-26/Accounts/" + accountSid + "/SMS/TemplateSMS?sig=" + sig

	req := httplib.Post(baseUrl)

	//req.SetTimeout(connectTimeout , readWriteTimeout * time.Second)
	req.Header("Accept", "application/json")
	req.Header("Content-Type", "application/json;charset=utf-8")
	req.Header("Content-Length", "256")
	req.Header("Authorization", auth)

	body := "{\"to\":\"" + phoneNum + "\",\"appId\":\"" + appId + "\",\"templateId\":\"1\",\"datas\":[\"" + randomCode + "\",\"" + strconv.Itoa(expireMin) + "\"]}"
	fmt.Println(body)

	//在这里发请求，发送验证码，钱不够，测试的时候再取消注释
	req.Body(body)
	result, err := req.String()

	fmt.Println(result, err)

	this.Data["json"] = "[{\"code\": 1, \"userId\": " + userId + ", \"phoneNum\": " + phoneNum + "}]"
	this.ServeJSON()
}

//sig:md5(所有字母必须大写) auth:base64
//短信验证码平台鉴权使用
func getSig(id string, token string) (string, string) {
	ltime := time.Now().Format("20060102150405")
	fmt.Println(ltime)

	sig := md5.New()
	sig.Write([]byte(id + token + ltime))

	auth := base64.StdEncoding.EncodeToString([]byte(id + ":" + ltime))

	return strings.ToUpper(hex.EncodeToString(sig.Sum(nil))), auth
}

//公共函数，获取一个以当前时间为sed的6位随机数
func GetRandomCode() string {
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	min := 1000
	code := r1.Intn(10000)
	if code < min {
		code += min
	}
	return strconv.Itoa(code)
}

// @router /Ver/verPhone [POST]
func (this *UserCenterController) VerPhone() {
	userId := this.GetString("userId")
	phoneNum := this.GetString("phoneNum")
	verCode := this.GetString("verCode")

	userIdInt64, _ := strconv.ParseInt(userId, 10, 64)

	content := redisClient.GetKey(PhoneVerPrefix + userId)

	if content != verCode {
		fmt.Println(verCode, content)
		logs.Error("input code %v is not equal redis code %v ", verCode, content)
		//！！！这里提示不友好，验证不通过会直接再次跳转验证页面，怎是没有提示
		this.Ctx.Redirect(302, "/Ver/phonever")
	}

	logs.Debug("input code %v is equal redis code %v ", verCode, content)

	var dbUser models.User
	dbUser.UpdateInfo(userIdInt64, "phone", phoneNum)
	dbUser.UpdateInfo(userIdInt64, "IsPhoneVer", "1")

	userinfo := redisClient.GetKey(LoginPrefix + userId)

	info := &UserLoginInfo{}
	err := json.Unmarshal([]byte(userinfo), &info)
	if err != nil {
		logs.Error("get userinfo from redis fail %v ", err)
		this.Ctx.Redirect(302, "/Login")
	}

	info.Phone = phoneNum
	info.IsPhoneVer = true

	data, _ := json.Marshal(info)

	redisClient.SetKey(LoginPrefix+userId, string(data))
	redisClient.Setexpire(LoginPrefix+userId, LoginPeriod)

	this.Ctx.Redirect(302, "/Portal/showdriverorder/")
}

func GetUserInfoFromRedis(uid string) *UserLoginInfo {
	userinfo := redisClient.GetKey(LoginPrefix + uid)

	info := &UserLoginInfo{}
	err := json.Unmarshal([]byte(userinfo), &info)
	if err != nil {
		logs.Error("get userinfo from redis fail %v ", err)
	}
	return info
}

func GetOnroadTypeFromId(uid string) int {
	var dbUser models.User
	var userInfo []*models.User
	dbUser.GetUserInfoFromId(uid, &userInfo)

	return userInfo[0].OnRoadType
}

// @router /Portal/userinfo [GET]
func (this *UserCenterController) UserInfo() {
	uid, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	var dbUser models.User
	var list []*models.User
	this.Data["success"] = true
	this.Data["tabIndex"] = 3

	success, _ := dbUser.GetUserInfoFromId(uid, &list)

	if success != "true" {
		this.Data["success"] = false
		return
	}
	this.Data["list"] = list[0]
	fmt.Println(list[0].Id)
	this.TplName = "userInfo.html"
}

// cache user info in redis
func CacheUserLoginInfo(userInfo *models.User) (string, string, error) {
	token := getToken(userInfo.Name, userInfo.Passwd)
	info := &UserLoginInfo{}
	info.Name = userInfo.Name
	info.IsPhoneVer = userInfo.IsPhoneVer
	info.IsDriver = userInfo.IsDriver
	info.Token = token
	info.Nickname = userInfo.Nickname
	info.OpenId = userInfo.OpenId
	info.OrderNumWV = userInfo.OrderNumWV
	info.Phone = userInfo.Phone

	data, _ := json.Marshal(info)
	fmt.Println("data: %v", string(data))
	idStr := strconv.Itoa(userInfo.Id)

	redisClient.SetKey(LoginPrefix+idStr, string(data))
	redisClient.Setexpire(LoginPrefix+idStr, LoginPeriod)

	return token, idStr, nil

}