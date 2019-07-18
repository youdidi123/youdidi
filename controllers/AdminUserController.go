package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"youdidi/commonLib"
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
	var cfWithdrew []*models.Cash_flow
	var cfWithdrewRefund []*models.Cash_flow
	var cfWithdrewError []*models.Cash_flow
	var cfWithdrewRefundError []*models.Cash_flow
	var cfWithdrewProcess []*models.Cash_flow

	tm, _ := commonLib.GetTodayBeginTime()
	endTime := tm - (24 * 60 * 60 * 1)
	this.Data["endTime"] = time.Unix(endTime,0).Format("2006-01-02 15:04")

	sum1 := 0.00
	sum2 := 0.00


	num1 := dbCashFlow.GetWithdrewOrder(&cfWithdrew, endTime, 1, 0)
	num2 := dbCashFlow.GetWithdrewOrder(&cfWithdrewRefund, endTime, 2, 0)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewError, endTime, 1, 2)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewRefundError, endTime, 2, 2)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewProcess, endTime, 2, 4)

	for i, v := range cfWithdrew {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrew[i].Time = tm.Format("2006-01-02 15:04")
		sum1 += v.Money
	}

	for i, v := range cfWithdrewRefund {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrewRefund[i].Time = tm.Format("2006-01-02 15:04")
		sum2 += v.Money
	}

	for i, v := range cfWithdrewError {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		finishTime64, _ := strconv.ParseInt(v.FinishTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		tm1 := time.Unix(finishTime64, 0)
		cfWithdrewError[i].Time = tm.Format("2006-01-02 15:04")
		cfWithdrewError[i].FinishTime = tm1.Format("2006-01-02 15:04")
	}

	for i, v := range cfWithdrewRefundError {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		finishTime64, _ := strconv.ParseInt(v.FinishTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		tm1 := time.Unix(finishTime64, 0)
		cfWithdrewRefundError[i].Time = tm.Format("2006-01-02 15:04")
		cfWithdrewRefundError[i].FinishTime = tm1.Format("2006-01-02 15:04")
	}

	for i, v := range cfWithdrewProcess {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrewProcess[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num1"] = num1
	this.Data["num2"] = num2
	this.Data["sum1"] = sum1
	this.Data["sum2"] = sum2

	this.Data["listw"] = cfWithdrew
	this.Data["listwr"] = cfWithdrewRefund
	this.Data["listwe"] = cfWithdrewError
	this.Data["listwre"] = cfWithdrewRefundError
	this.Data["listwp"] = cfWithdrewProcess

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

// @router /admin/dealwithdrew [POST]
func (this *AdminUserController) DealWithDrew () {
	oids := this.GetStrings("oid")
	logs.Debug("withdrew list %v", oids)
	var dbCf models.Cash_flow

	code := 0
	msg := ""

	for _, oid := range oids {
		logs.Debug("withdrew id %v", oid)
		var cfInfo []* models.Cash_flow
		_, num := dbCf.GetOrderInfo(oid, &cfInfo)
		if (num != 1) {
			logs.Error("withdrew order id has something wrong oid=%v reNum=%v", oid, num)
			continue
		}
		wxId, err := WxEnpTransfers(int64(cfInfo[0].Money * 100), cfInfo[0].User.OpenId, oid, "192.168.0.1", "平台账户提现")
		if (err != nil) {
			_, err := dbCf.UpdateWithDrewResult(false, oid, "", err.Error())
			if (err != nil) {
				logs.Error("update withdrew result fail oid=%v result=fail err=%v", oid, err.Error())
			}
		} else {
			_, err := dbCf.UpdateWithDrewResult(true, oid, wxId, "")
			if (err != nil) {
				logs.Error("update withdrew result fail oid=%v result=success err=%v", oid, err.Error())
			}
			moneyStr := strconv.FormatFloat(cfInfo[0].Money, 'G' , -1,64)
			balanceStr := strconv.FormatFloat(cfInfo[0].User.Balance, 'G' , -1,64)
			commonLib.SendMsg5(cfInfo[0].User.OpenId,
				4, "", "#173177", "", "",
				"#173177", "",
				"#ff0000","预扣车费",
				"#22c32e", "预扣成功",
				"#173177", moneyStr,
				"#173177", balanceStr)
		}
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}