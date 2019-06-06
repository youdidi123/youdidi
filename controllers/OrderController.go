package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/models"
	"github.com/astaxie/beego/logs"
)
type OrderController struct {
	beego.Controller
}


//订单ID生成规则：精确到秒的日志字符串+用户ID（保证是6位，100000用户）+4位随机数字
func genOrderId (uid int) string {
	oid := ""
	uidS := ""
	ltime := time.Now().Format("200601021504")
	//uidInt, _ := strconv.Atoi(uid)
	if (uid < 100000) {
		uid = uid + 100000
		uidS = strconv.Itoa(uid)
	}

	oid = ltime + uidS + GetRandomCode()

	return oid
}

//点击底部导航栏"车主行程"进入的页面
// @router /Portal/showdriverorder/ [GET]
func (this *OrderController) ShowDriverOrder () {
	uid , _ := this.Ctx.GetSecureCookie("qyt","qyt_id")

	var dbOrder models.Order
	var orderInfo []*models.Order

	num := dbOrder.GetOrderInfoFromUserId(uid, &orderInfo)

	statusText := []struct {
		Text string
	}{{"接单中"}, {"司机到达"}, {"行程中"}, {"到达"}, {"结束"}, {"取消"}, {"无效"}}

	for i , v := range orderInfo{
		launchTime64 , _ := strconv.ParseInt(v.LaunchTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04:05")
	}

	this.Data["orderNum"] = num
	this.Data["orderInfo"] =  orderInfo
	this.Data["StatusText"] = statusText

	onRoadType := GetOnroadTypeFromId(uid)
	if (onRoadType == 0) {
		this.Data["buttonHref"] = "/Portal/createorder"
		this.Data["buttonValue"] = "发布行程"
		this.Data["buttonFunc"] = "createOrder()"
		this.Data["buttonClass"] = "weui-btn weui-btn_primary"
	} else if (onRoadType == 1) {
		this.Data["buttonHref"] = "#"
		this.Data["buttonValue"] = "发布行程"
		this.Data["buttonFunc"] = ""
		this.Data["buttonClass"] = "weui-btn weui-btn_primary weui-btn_disabled"
	} else {
		var currentOrder []*models.Order
		dbOrder.GetCurrentOrderFromUserId(uid, &currentOrder)
		currentOrderId := currentOrder[0].Id
		this.Data["buttonHref"] = "/Portal/driverorderdetail/" + currentOrderId
		this.Data["buttonValue"] = "进入当前行程"
		this.Data["buttonFunc"] = "gotoDetail(" + currentOrderId + ")"
		this.Data["buttonClass"] = "weui-btn weui-btn_primary"
	}
	this.Data["onRoadType"] = onRoadType
	this.Data["tabIndex"] = 2
	this.TplName = "driverOrder.html"
}

// 点击发布行程进入这里过滤
// 车主行程里点击发布行程进入这个逻辑
// 公众号里点击发布行程直接进入这个
// @router /Portal/createorder [GET]
func (this *OrderController) CreateOrder () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	onRoadType := GetOnroadTypeFromId(uid)
	if (onRoadType != 0) {
		logs.Debug("user %v onRoadTyep %v" , uid , onRoadType)
		//如果用户正在一个行程中，则跳转回查看行程页面
		this.ShowDriverOrder()
		return
	}
	userInfo := GetUserInfoFromRedis(uid)
	this.Data["tabIndex"] = 2
	if (userInfo.IsDriver != 2) {
		//如果用户还没有认证司机，则进入认证提醒页面
		orderNumWithoutVer,_ := strconv.Atoi(beego.AppConfig.String("orderNumWithoutVer"))
		this.TplName = "createOrderFilter.html"
		//如果num=0 ，继续发单按钮为disable
		this.Data["num"] = orderNumWithoutVer - userInfo.OrderNumWV
		this.Data["orderNumWithoutVer"] = orderNumWithoutVer
	} else {
		this.TplName = "createOrder.html"
		this.Data["uid"] = uid
	}
}

// 在发单的车主提示拦截中，如果用户点击继续发单，则计入这个逻辑
// @router /Portal/createorderforce [GET]
func (this *OrderController) CreateOrderForce () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	this.Data["tabIndex"] = 2
	this.TplName = "createOrder.html"
	this.Data["uid"] = uid
}

// @router /Portal/docreateorder [POST]
func (this *OrderController) DoCreateOrder () {
	code := 0
	msg := ""

	userIdS := this.GetString("uid")

	onRoadType :=  GetOnroadTypeFromId(userIdS)

	if (onRoadType != 0) {
		code = 2
		msg = "<br>已在行程中，请勿重复发单"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	userId , _ := strconv.Atoi(userIdS)
	launchTime := this.GetString("launchTime")
	startCode, _ := strconv.Atoi(this.GetString("startCode"))
	endCode, _ := strconv.Atoi(this.GetString("endCode"))
	charge := this.GetString("charge")
	siteNum := this.GetString("siteNum")
	travelExplain := this.GetString("travelExplain")
	travelCommit := this.GetString("travelCommit")


	orderId := genOrderId(userId)

	var dbOrder models.Order

	dbOrder.Id = orderId
	dbOrder.User  = &models.User{Id:userId}
	dbOrder.LaunchTime = launchTime
	dbOrder.CreateTime = strconv.FormatInt(time.Now().Unix(),10)
	dbOrder.SrcId  = &models.Location{Id:int64(startCode)}
	dbOrder.DestId  = &models.Location{Id:int64(endCode)}
	dbOrder.PNum , _ = strconv.Atoi(siteNum)
	dbOrder.ThroughL = travelExplain
	dbOrder.Marks = travelCommit
	dbOrder.Price , _ = strconv.ParseFloat(charge,64)

	_ , err := dbOrder.Insert()

	if (err != nil) {
		code = 1
		msg = "发布行程失败，请重试"
		logs.Error("insert order fail , err=" , msg)
	} else {
		var dbUser models.User
		userId64  := int64(userId)
		//这里是车主发单，所有onroadtype等于2
		dbUser.UpdateInfo(userId64 , "onRoadType" , "2")
	}

	//fmt.Println(userId)

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/driverorderdetail/:oid [GET]
func (this *OrderController) DriverOrderDetail () {
	//test := this.Ctx.Request.RequestURI
	oid := this.GetString(":oid")
	this.Data["url"] = oid
	this.TplName = "driverOrderDetail.html"
}

// @router /Portal/searchorder [GET]
func (this *OrderController) SearchOrder () {
	day := this.GetString("uid")
	startCode := this.GetString("startCode")
	endCode := this.GetString("startCode")
	logs.Debug("search order day=%v start=%v end=%v", day , startCode , endCode)
	this.Data["tabIndex"] = 0
	this.TplName = "searchOrder.html"
}