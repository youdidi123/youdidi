package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/models"
	"github.com/astaxie/beego/logs"
	"youdidi/redisClient"
)

type AdminUserController struct {
	beego.Controller
}

type AdminUserLoginInfo struct {
	Id int
	Name string
	Token string
	Type int
}

var (
	AdminLoginPeriod = 1* 24 * 60 * 60 //用户登陆有效期1天
	AdminLoginPrefix = "ADMIN_LOGIN_INFO_"  //缓存在redis中的用户数据key前缀
)


// @router /AdminLogin/ [GET]
func (this *AdminUserController) AdminLogin (){
	this.TplName = "adminLogin.html"
}

// @router /admin/ [GET]
func (this *AdminUserController) Admin (){
	this.TplName = "adminHomepage.html"
}

// @router /Admin/doLogin [POST]
func (this *AdminUserController) AdminDoLogin (){
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	var dbAdminUser models.Admin_user
	var userInfo []*models.Admin_user

	num := dbAdminUser.GetUserInfoFromName(inputName, &userInfo)
	if (num < 1) {
		this.Redirect("/AdminLogin/", 302)
		return
	}
	h := md5.New()
	h.Write([]byte(inputPasswd))

	passwdMd5 := hex.EncodeToString(h.Sum(nil))

	if (passwdMd5 != userInfo[0].Passwd) {
		logs.Error("input passwd is not correct inputpasswd=%v md5=%v", inputPasswd, passwdMd5)
		this.Redirect("/AdminLogin/", 302)
		return
	}

	token := getToken(inputName , inputPasswd)

	userInfoRedis := &AdminUserLoginInfo{}
	userInfoRedis.Id = userInfo[0].Id
	userInfoRedis.Name = userInfo[0].Name
	userInfoRedis.Type= userInfo[0].Type
	userInfoRedis.Token = token

	data, _ := json.Marshal(userInfoRedis)
	idStr := strconv.Itoa(userInfo[0].Id)

	redisClient.SetKey(AdminLoginPrefix+idStr , string(data))
	redisClient.Setexpire(AdminLoginPrefix+idStr , AdminLoginPeriod)

	this.Ctx.SetSecureCookie("qyt","qyt_admin_id" , idStr) //注入用户id，后续所有用户id都从cookie里获取
	this.Ctx.SetSecureCookie("qyt","qyt_admin_token" , token)

	this.TplName = "adminHomepage.html"
}

// @router /admin/dconfirm [GET]
func (this *AdminUserController) DriverConfirm (){
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetNoConfirm(&dcInfo)

	for i, v := range dcInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		dcInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = dcInfo

	this.TplName = "adminDriverConfirm.html"
}

// @router /admin/confirmDriverDetail/:id [GET]
func (this *AdminUserController) ConfirmDriverDetail (){
	id := this.GetString(":id")
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetOrderFromId(id, &dcInfo)

	for i, v := range dcInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		dcInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = dcInfo[0]

	this.TplName = "adminConfirmDriverDetail.html"
}