package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"youdidi/models"
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

// @router /AdmindoLogin [POST]
func (this *AdminUserController) AdminDoLogin (){
	logs.Debug("wo zai zhe")
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	var dbAdminUser models.Admin_user
	var userInfo []*models.Admin_user


	num := dbAdminUser.GetUserInfoFromName(inputName, &userInfo)
	if (num < 1) {
		logs.Debug("can not find user name=%v", inputName)
		this.Redirect("/AdminLogin/", 302)
		return
	}
	h := md5.New()
	h.Write([]byte(inputPasswd))

	passwdMd5 := hex.EncodeToString(h.Sum(nil))

	if (passwdMd5 != userInfo[0].Passwd) {
		logs.Error("input passwd is not correct inputpasswd=%v md5=%v dbvalue=%v", inputPasswd, passwdMd5, userInfo[0].Passwd)
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

// @router /admin/userwithdrew [GET]
func (this *AdminUserController) UserWithdrew () {
	var dbCashFlow models.Cash_flow
	var cfInfo []*models.Cash_flow

	num := dbCashFlow.GetReadyOrder(&cfInfo)

	for i, v := range cfInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = cfInfo

	this.TplName = "adminUserWithdrew.html"
}

// @router /admin/doconfirmdriver [POST]
func (this *AdminUserController) DoConfirmDriver () {
	oid := this.GetString("oid")
	aType := this.GetString("type")
	mark := this.GetString("mark")

	var dbDc models.Driver_confirm

	code := 0
	msg := ""

	if (! dbDc.DoConfirmDriver(oid, aType, mark)) {
		code = 1
		msg = "系统异常，请重试"
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

// @router /admin/showcomplain [GET]
func (this *AdminUserController) ShowComplain () {
	var dbC models.Complain
	var cInfo []*models.Complain

	num, _ := dbC.GetNoComplain(&cInfo)

	for i, v := range cInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = cInfo
	this.TplName = "adminShowComplain.html"
}