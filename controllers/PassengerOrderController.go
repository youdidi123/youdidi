package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/commonLib"
	"youdidi/models"
	"github.com/astaxie/beego/logs"
)

type PassengerOrderController struct {
	beego.Controller
}

// @router /Portal/pcreateorder [GET]
func (this *PassengerOrderController) PcreateOrder () {
	this.Data["tabIndex"] = 0
	this.TplName = "pCreateOrder.html"
}

// @router /Portal/rcreateorder [POST]
func (this *PassengerOrderController) RcreateOrder () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	code := 0
	msg := ""
	url := "/Portal/pcreateorder/"

	var dbUser models.User
	var userInfo []*models.User

	succ, num := dbUser.GetUserInfoFromId(uid, &userInfo)

	if (num < 1 || succ == "false") {
		code = 1
		msg = "网络开小差了哦，请稍后重试"
		url = "#"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg, "url":url};
		this.ServeJSON()
		return
	}
	if (userInfo[0].Balance <= 0) {
		code = 2
		msg = "您的账户内无可用余额，请及时充值"
		url = "/Portal/invest"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg, "url":url};
		this.ServeJSON()
		return
	}
	if (userInfo[0].OnRoadType == 1 || userInfo[0].OnRoadType == 3) {
		code = 3
		msg = "您有尚未完成的行程，请及时处理"
		url = "/Portal/showpassengerorder"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg, "url":url};
		this.ServeJSON()
		return
	}
	if (userInfo[0].OnRoadType == 2) {
		code = 4
		msg = "您有尚未完成的车主行程，请及时处理"
		url = "/Portal/showdriverorder"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg, "url":url};
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg, "url":url};
	this.ServeJSON()
}

// @router /Portal/dopcreateorder [POST]
func (this *PassengerOrderController) PdoCreateOrder () {
	var dbPo models.PassengerOrder
	msg := ""

	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	launchTime := this.GetString("launchTime")
	startCode, _ := this.GetInt64("startCode")
	endCode, _ := this.GetInt64("endCode")
	charge, _ := this.GetFloat("charge")
	siteNum, _ := this.GetInt("siteNum")
	travelExplain := this.GetString("travelExplain")
	travelCommit := this.GetString("travelCommit")

	code, err := dbPo.CreateOrder(
		uid, launchTime,
		startCode, endCode,
		charge, siteNum,
		travelExplain, travelCommit)

	if (err != nil) {
		msg = err.Error()
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

// @router /Portal/canclepo [POST]
func (this *PassengerOrderController) CancleOrder () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	oid := this.GetString("id")

	var dbPo models.PassengerOrder
	msg := ""

	code, err := dbPo.CancleOrder(uid, oid)

	if (err != nil) {
		msg = err.Error()
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/searchpinput [GET]
func (this *PassengerOrderController) SearchPassengerInput () {
	var dbPo models.PassengerOrder
	var poList []*models.PassengerOrder

	num, err := dbPo.GetTop20Orders(&poList)

	if (err != nil) {
		logs.Error("get top 20 passenger order fail err=%v", err.Error())
	}

	currentTime, _ := strconv.ParseInt(commonLib.GetCurrentTime(), 10, 64)

	for i, v :=  range poList {
		launchTime, _ := strconv.ParseInt(v.LaunchTime, 10, 64)
		poList[i].LaunchTime = commonLib.FormatUnixToStr(v.LaunchTime)
		if (currentTime - launchTime > (5 * 60) && v.Status == 0) {
			poList[i].Status = 4
		}
	}

	this.Data["list"] = poList
	this.Data["num"] = num
	this.Data["tabIndex"] = 2
	this.TplName = "searchPinput.html"
}

// @router /Portal/qiangdanbefore [POST]
func (this *PassengerOrderController) LockPOrderBefore () {
	porderId := this.GetString("oid")
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")

	var dbPo models.PassengerOrder
	msg := ""

	code, ops, err := dbPo.LockOrderBefore(uid, porderId)

	if (err != nil) {
		msg = err.Error()
	}

	this.Data["json"] = map[string]interface{}{"code":code, "ops":ops, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/createandconfirm [POST]
func (this *PassengerOrderController) CreateAndConfirm () {
	porderId := this.GetString("oid")
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	uidInt, _ := strconv.Atoi(uid)

	var dbPo models.PassengerOrder
	msg := ""

	code, url, err := dbPo.CreateAndConfirm(uidInt, porderId, genOrderId(uidInt))

	if (err != nil) {
		msg = err.Error()
	}

	this.Data["json"] = map[string]interface{}{"code":code, "url":url, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/driverqiangdan [POST]
func (this *PassengerOrderController) LockOrder () {
	porderId := this.GetString("oid")
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	uidInt, _ := strconv.Atoi(uid)

	var dbPo models.PassengerOrder
	msg := ""

	code, url, err := dbPo.LockAndConfirm(uidInt, porderId)

	if (err != nil) {
		msg = err.Error()
	}

	this.Data["json"] = map[string]interface{}{"code":code, "url":url, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/searchporder/:startcode/:endcode/:launchtime [GET]
func (this *PassengerOrderController) SearchPOrder () {
	startCode, _ := strconv.Atoi(this.GetString(":startcode"))
	endCode, _ := strconv.Atoi(this.GetString(":endcode"))
	//startCode64 , _ := strconv.ParseInt(startCode, 10, 64)
	//endCode64 , _ := strconv.ParseInt(endCode, 10, 64)
	startCodeLocation := startCode % 1000000
	endCodeLocation := endCode % 1000000

	launchTime := this.GetString(":launchtime")
	launchTime = launchTime + " 00:00:00"
	tmStart, _ := time.ParseInLocation("2006-01-02 15:04:05", launchTime, time.Local)
	tmEnd := tmStart.Unix() + (1*24*60*60)

	logs.Debug("search order launchTime=%v start=%v end=%v", launchTime , startCode , endCode)

	var dbOrder models.PassengerOrder
	var orderInfo []*models.PassengerOrder

	num := dbOrder.GetReadyOrders(&orderInfo , startCode , endCode ,
		tmStart.Unix(), tmEnd , startCodeLocation , endCodeLocation)

	currentTime, _ := strconv.ParseInt(commonLib.GetCurrentTime(), 10, 64)

	for i , v := range orderInfo{
		launchTime64 , _ := strconv.ParseInt(v.LaunchTime, 10, 64)
		if (currentTime - launchTime64 > (5 * 60) && v.Status == 0) {
			orderInfo[i].Status = 4
		}
		tm := time.Unix(launchTime64, 0)
		orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04")

	}

	this.Data["num"] = num
	this.Data["list"] = orderInfo
	this.Data["tabIndex"] = 2
	this.TplName = "searchPOrder.html"
}