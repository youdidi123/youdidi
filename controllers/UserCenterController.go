package controllers

import (
	"crypto/md5"
    "encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/context"
	"strconv"
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
	idStr      string
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

				userLoginInfo, err := GenUserLoginInfo(list[0])
				if (err != nil) {
					logs.Error("Generate user login info failed!")
				}

				err = CacheUserLoginInfo(userLoginInfo)

				if (err != nil) {
					logs.Warn("Cache user login info failed!")
				}

				this.Ctx.SetSecureCookie("qyt", "qyt_id", userLoginInfo.idStr) //注入用户id，后续所有用户id都从cookie里获取
				this.Ctx.SetSecureCookie("qyt", "qyt_token", userLoginInfo.Token)

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

	code := 0
	msg := ""

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

	if (err != nil) {
		code = 1
		msg = "网络开小差了哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	code = 0
	msg = "操作成功"
	logs.Debug(result)

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Ver/verPhone [POST]
func (this *UserCenterController) VerPhone() {
	userId := this.GetString("userId")
	phoneNum := this.GetString("phoneNum")
	verCode := this.GetString("verCode")

	code := 0
	msg := ""

	userIdInt64, _ := strconv.ParseInt(userId, 10, 64)

	content := redisClient.GetKey(PhoneVerPrefix + userId)

	if (content != verCode) {
		logs.Error("input code %v is not equal redis code %v " , verCode , content)
		//！！！这里提示不友好，验证不通过会直接再次跳转验证页面，怎是没有提示
		code = 1
		msg = "验证码错误"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	logs.Debug("input code %v is equal redis code %v ", verCode, content)

	var dbUser models.User
	dbUser.UpdateInfo(userIdInt64, "phone", phoneNum)
	dbUser.UpdateInfo(userIdInt64, "IsPhoneVer", "1")

	userinfo := redisClient.GetKey(LoginPrefix + userId)

	info := &UserLoginInfo{}
	err := json.Unmarshal([]byte(userinfo), &info)
	if (err != nil) {
		logs.Error("get userinfo from redis fail %v " , err)
		code = 2
		msg = "网络开小差了哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	info.Phone = phoneNum
	info.IsPhoneVer = true

	data, _ := json.Marshal(info)

	redisClient.SetKey(LoginPrefix+userId, string(data))
	redisClient.Setexpire(LoginPrefix+userId, LoginPeriod)

	code = 0
	msg = "操作成功"
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
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

// @router /Portal/help [GET]
func (this *UserCenterController) Help() {
	this.Data["tabIndex"] = 3

	this.TplName = "help.html"
}

// @router /Portal/aboutus [GET]
func (this *UserCenterController) AboutUs() {
	this.Data["tabIndex"] = 3

	this.TplName = "aboutUs.html"
}

// @router /Portal/disclaimer [GET]
func (this *UserCenterController) Disclaimer() {
	this.Data["tabIndex"] = 3

	this.TplName = "disclaimer.html"
}

// Get user info from redis
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
	_ , num := dbUser.GetUserInfoFromId(uid, &userInfo)

	if (num == 0 ) {
		logs.Error("get user onroad type err uid=%v" , uid)
		dbUser.GetUserInfoFromId(uid, &userInfo)
		return 0
	}

	return userInfo[0].OnRoadType
}

// generate user login info
func GenUserLoginInfo(userInfo *models.User) (*UserLoginInfo, error) {
	token := getToken(userInfo.Name, userInfo.Passwd)
	userLoginInfo := &UserLoginInfo{}
	userLoginInfo.Name = userInfo.Name
	userLoginInfo.IsPhoneVer = userInfo.IsPhoneVer
	userLoginInfo.IsDriver = userInfo.IsDriver
	userLoginInfo.Token = token
	userLoginInfo.Nickname = userInfo.Nickname
	userLoginInfo.OpenId = userInfo.OpenId
	userLoginInfo.OrderNumWV = userInfo.OrderNumWV
	userLoginInfo.Phone = userInfo.Phone
	userLoginInfo.idStr = strconv.Itoa(userInfo.Id)

	return userLoginInfo, nil
}

// cache user info in redis
func CacheUserLoginInfo(userLoginInfo *UserLoginInfo) error {
	data, _ := json.Marshal(userLoginInfo)
	fmt.Println("data: %v", string(data))

	redisClient.SetKey(LoginPrefix+userLoginInfo.idStr, string(data))
	redisClient.Setexpire(LoginPrefix+userLoginInfo.idStr, LoginPeriod)
	return nil
}



func GetUserLoginInfoByCookie(ctx *context.Context) (*UserLoginInfo, error) {
	info := &UserLoginInfo{}
	id, isId := ctx.GetSecureCookie("qyt","qyt_id")
	if (! isId) {
		return nil, fmt.Errorf("can not get id from cookie:")
	} else {
		logs.Debug("id of cookis : %v" , id)
		token, isToken := ctx.GetSecureCookie("qyt" , "qyt_token")
		if (! isToken) {
			return nil, fmt.Errorf("can not get token from cookie")
		} else {
			content := redisClient.GetKey(LoginPrefix+id)
			if (content == "nil") {
				return nil, fmt.Errorf("cache is empty")
			} else {
				err := json.Unmarshal([]byte(content), &info)
				if (err != nil) {
					return nil, fmt.Errorf("Userinfo Unmarshal error:%s", err)
				} else {
					if (token != info.Token) {
						return nil, fmt.Errorf("token did not match of cookie and cache")
					}
				}
			}
		}
	}
	return info, nil
}