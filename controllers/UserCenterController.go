package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"youdidi/models"
	"youdidi/redisClient"
	"github.com/astaxie/beego/httplib"
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
	PhoneVerPrefix = "RANDOM_CODE_"
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
				this.SetSession("qyt_id" , idStr)

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
	//id, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	id := this.GetSession("qyt_id")
	this.Data["userId"] = id
	this.TplName = "phoneVer.html"
}

// @router /Ver/getvercode [GET,POST]
func (this *UserCenterController) GetVerCode() {
	expireMin := 5
	userId := this.GetString("userId")
	phoneNum := this.GetString("phoneNum")

	baseUrl := beego.AppConfig.String("smsBaseUrl")
	//connectTimeout , _ := strconv.Atoi(beego.AppConfig.String("connectTimeout"))
	//readWriteTimeout , _ := strconv.Atoi(beego.AppConfig.String("readWriteTimeout"))

	accountSid := "8aaf07086ab0c082016ab465923401a3"
	token := "868401e59c874a87874b9d9d028c3e17"
	appId := "8aaf07086ab0c082016ab465928a01a9"

	sig, auth := getSig(accountSid, token)
	randomCode := GetRandomCode()

	var cacheClient redisClient.CacheClient
	cacheClient.GetConnet()
	cacheClient.SetKey(PhoneVerPrefix+userId , randomCode)
	cacheClient.Setexpire(PhoneVerPrefix+userId , expireMin * 60)

	baseUrl = baseUrl + "/2013-12-26/Accounts/" + accountSid + "/SMS/TemplateSMS?sig=" + sig

	req := httplib.Post(baseUrl)

	//req.SetTimeout(connectTimeout , readWriteTimeout * time.Second)
	req.Header("Accept","application/json")
	req.Header("Content-Type","application/json;charset=utf-8")
	req.Header("Content-Length","256")
	req.Header("Authorization",auth)

	body := "{\"to\":\"" + phoneNum + "\",\"appId\":\"" + appId + "\",\"templateId\":\"1\",\"datas\":[\"" + randomCode + "\",\"" + strconv.Itoa(expireMin) + "\"]}"
	fmt.Println(body)

	//在这里发请求，发送验证码，钱不够，测试的时候再取消注释
	req.Body(body)
	result , err := req.String()

	fmt.Println(result , err)

	this.Data["json"] = "[{\"code\": 1, \"userId\": " + userId + ", \"phoneNum\": " + phoneNum + "}]"
	this.ServeJSON()
	}

func getSig (id string , token string) (string , string){
	ltime := time.Now().Format("20060102150405")
	fmt.Println(ltime)

	sig := md5.New()
	sig.Write([]byte(id+token+ltime))

	auth := base64.StdEncoding.EncodeToString([]byte(id+":"+ltime))

	return strings.ToUpper(hex.EncodeToString(sig.Sum(nil))),auth
}

func GetRandomCode () string{
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	min := 100000
	code := r1.Intn(1000000)
	if (code < min) {
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

	var cacheClient redisClient.CacheClient
	cacheClient.GetConnet()
	content := cacheClient.GetKey(PhoneVerPrefix+userId)

	if (content != verCode) {
		fmt.Println(verCode , content)
		logs.Error("input code %v is not equal redis code %v " , verCode , content)
		this.Ctx.Redirect(302, "/Ver/phonever")
	}

	logs.Debug("input code %v is equal redis code %v " , verCode , content)

	var dbUser models.User
	dbUser.UpdateInfo(userIdInt64 , "phone" , phoneNum)
	dbUser.UpdateInfo(userIdInt64 , "IsPhoneVer" , "1")

	userinfo := cacheClient.GetKey(LoginPrefix+userId)

	info := &UserLoginInfo{}
	err := json.Unmarshal([]byte(userinfo), &info)
	if (err != nil) {
		logs.Error("get userinfo from redis fail %v " , err)
		this.Ctx.Redirect(302, "/Login")
	}

	info.Phone = phoneNum
	info.IsPhoneVer = true

	data, _ := json.Marshal(info)

	cacheClient.SetKey(LoginPrefix+userId , string(data))
	cacheClient.Setexpire(LoginPrefix+userId , LoginPeriod)

	this.Ctx.Redirect(302, "/Portal/home")
}