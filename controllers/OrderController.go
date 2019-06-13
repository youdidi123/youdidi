package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/models"
	"github.com/astaxie/beego/logs"
)
type OrderController struct {
	beego.Controller
}

var (
	statusText = []struct {
		Text string
	}{{"接单中"},{"暂停接单"}, {"司机到达"}, {"行程中"}, {"到达"}, {"完成"}, {"取消"}, {"无效"}}
	odStatusText = []struct {
		Text string
	}{{"等待车主确认"},{"等待改价确认"}, {"等待出发"}, {"行程中"}, {"乘客到达"},{"等待评价"},{"完成"},  {"取消"}}
)

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
		msg = "已在行程中，请勿重复发单"
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
	dbOrder.SrcLocationId = dbOrder.SrcId.Id % 1000000
	dbOrder.DestLocationId = dbOrder.DestId.Id % 1000000
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
	this.Data["tabIndex"] = 2
	this.Data["url"] = oid
	this.TplName = "driverOrderDetail.html"
}

// @router /Portal/passengerorderdetail/:oid [GET]
func (this *OrderController) PassengerOrderDetail () {
	uid, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString(":oid")

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	num := dbOrderDetail.GetOrderInfoFromPassengerId(oid , uid , &orderDetailInfo)

	if (num > 0) {
		for i, v := range orderDetailInfo {
			launchTime64, _ := strconv.ParseInt(v.Order.LaunchTime, 10, 64)
			tm := time.Unix(launchTime64, 0)
			orderDetailInfo[i].Order.LaunchTime = tm.Format("2006-01-02 15:04")
		}
		this.Data["List"] = orderDetailInfo[0]
		//if (orderDetailInfo[0].Status )
	}

	this.Data["StatusText"] = odStatusText
	this.Data["num"] = num
	this.Data["tabIndex"] = 1
	this.Data["url"] = oid
	this.TplName = "passengerOrderDetail.html"
}

// @router /Portal/searchorder [GET,POST]
func (this *OrderController) SearchOrder () {
	startCode := this.GetString("startCode")
	endCode := this.GetString("endCode")
	startCode64 , _ := strconv.ParseInt(startCode, 10, 64)
	endCode64 , _ := strconv.ParseInt(endCode, 10, 64)
	startCodeLocation := startCode64 % 1000000
	endCodeLocation := endCode64 % 1000000

	launchTime := this.GetString("launchTime")
	launchTime = launchTime + " 00:00"

	tmStart, _ := time.Parse("2006-01-02 15:04", launchTime)
	tmEnd := tmStart.Unix() + (1*24*60*60)

	logs.Debug("search order launchTime=%v start=%v end=%v", launchTime , startCode , endCode)
	var dbOrder models.Order
	var orderInfo []*models.Order

	num := dbOrder.GetReadyOrders(&orderInfo , startCode64 , endCode64 ,
		tmStart.Unix() , tmEnd , startCodeLocation , endCodeLocation)

	for i , v := range orderInfo{
		launchTime64 , _ := strconv.ParseInt(v.LaunchTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04")
	}

	this.Data["StatusText"] = statusText
	this.Data["num"] = num
	this.Data["orders"] = orderInfo
	this.Data["tabIndex"] = 0
	this.TplName = "searchOrder.html"
}

// @router /Portal/dorequire [POST]
func (this *OrderController) DoRequire () {
	code := 0
	msg := ""

	oid := this.GetString("oid")
	count , _ := strconv.Atoi(this.GetString("count"))
	mark := this.GetString("mark")

	fmt.Println(mark)

	userIdS, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	userId , _ := strconv.Atoi(userIdS)

	onRoadType :=  GetOnroadTypeFromId(userIdS)

	if (onRoadType == 1) {
		code = 1
		msg = "您有以预约行程尚未结束，请勿重复发起预约"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	} else if (onRoadType == 2) {
		code = 2
		msg = "您有车主行程尚未结束"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	numOd := dbOrderDetail.GetOrderInfoFromPassengerId(oid , userIdS , &orderDetailInfo)

	if (numOd > 0) {
		code = 8
		msg = "您以预约过该行程，请勿重复预约"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	//此处给订单加锁，后面必须要释放！！！
	if (! SetOrderLock(oid)) {
		code = 3
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbOrder models.Order
	var orderInfo []*models.Order

	dbOrder.GetOrderFromId(oid, &orderInfo)

	//判断当前订单的空余座位数是不是满足要求
	if (count > (orderInfo[0].PNum - (orderInfo[0].ConfirmPnum + orderInfo[0].RequestPnum))) {
		logs.Error("sitnum do not match require rnum=%v restnum=%v" , count ,
			orderInfo[0].PNum - (orderInfo[0].ConfirmPnum + orderInfo[0].RequestPnum))
		//优先释放锁，不管成功不成功都要继续
		DelOrderLock(oid)
		code = 4
		msg = "来晚一步。。当前剩余座位数不足"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbUser models.User
	var userInfo []*models.User

	_, num := dbUser.GetUserInfoFromId(userIdS, &userInfo)

	if (num < 1) {
		//优先释放锁，不管成功不成功都要继续
		logs.Emergency("get userinfo fail uid=%v" , userIdS)
		DelOrderLock(oid)
		code = 5
		msg = "系统后台异常，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userInfo[0].Balance < orderInfo[0].Price * float64(count)) {
		//优先释放锁，不管成功不成功都要继续
		logs.Notice("user balance is not enough uid=%v balance=%v orderprice=%v requireSite=%v" ,
			userIdS , userInfo[0].Balance, orderInfo[0].Price , count)
		DelOrderLock(oid)
		code = 6
		msg = "账户余额不足，当前余额为:" + strconv.FormatFloat(userInfo[0].Balance, 'G' , -1,64) +
			"元 不足支付行程总价:" + strconv.FormatFloat(orderInfo[0].Price * float64(count), 'G' , -1,64) +
			"元"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var od models.Order_detail
	od.Order = &models.Order{Id:oid}
	od.IsPayed = true
	od.Passage = &models.User{Id:userId}
	od.Driver = &models.User{Id:orderInfo[0].User.Id}
	od.SiteNum = count

	if (orderInfo[0].DoRequire(od, userIdS, count , mark , userInfo[0].Balance - orderInfo[0].Price * float64(count))) {
		DelOrderLock(oid)
		code = 0
		msg = "预约成功"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
	} else {
		DelOrderLock(oid)
		code = 7
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
	}


}

// @router /Portal/searchinput [GET]
func (this *OrderController) SearchInput () {
	this.Data["tabIndex"] = 0
	this.TplName = "searchInput.html"
}

//点击底部导航栏"车主行程"进入的页面
// @router /Portal/showdriverorder/ [GET]
func (this *OrderController) ShowDriverOrder () {
	uid, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")

	var dbOrder models.Order
	var orderInfo []*models.Order

	num := dbOrder.GetOrderInfoFromUserId(uid, &orderInfo)

	if (num > 0) {
		for i, v := range orderInfo {
			launchTime64, _ := strconv.ParseInt(v.LaunchTime, 10, 64)
			tm := time.Unix(launchTime64, 0)
			orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04")
		}
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
		num := dbOrder.GetCurrentDriverOrderFromUserId(uid, &currentOrder)
		if (num > 0) {
			currentOrderId := currentOrder[0].Id
			this.Data["buttonHref"] = "/Portal/driverorderdetail/" + currentOrderId
			this.Data["buttonValue"] = "进入当前行程"
			this.Data["buttonFunc"] = "gotoDetail(" + currentOrderId + ")"
		} else {
			this.Data["buttonHref"] = "/Portal/createorder"
			this.Data["buttonValue"] = "发布行程"
			this.Data["buttonFunc"] = "createOrder()"
		}
		this.Data["buttonClass"] = "weui-btn weui-btn_primary"
	}
	this.Data["onRoadType"] = onRoadType
	this.Data["tabIndex"] = 2
	this.TplName = "driverOrder.html"
}

// @router /Portal/showpassengerorder [GET]
func (this *OrderController) ShowPassengerOrder () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	num := dbOrderDetail.GetOrderDetailFromPassengerId(userId , &orderDetailInfo)

	if (num > 0) {
		for i, v := range orderDetailInfo {
			launchTime64, _ := strconv.ParseInt(v.Order.LaunchTime, 10, 64)
			tm := time.Unix(launchTime64, 0)
			orderDetailInfo[i].Order.LaunchTime = tm.Format("2006-01-02 15:04")
		}
	}

	onRoadType := GetOnroadTypeFromId(userId)
	if (onRoadType != 1) {
		this.Data["buttonHref"] = "/Portal/searchinput"
		this.Data["buttonValue"] = "找拼车"
		this.Data["buttonClass"] = "weui-btn weui-btn_primary"
	} else {
		var currentOrder []*models.Order_detail
		num := dbOrderDetail.GetCurrentPassengerOrderFromUserId(userId, &currentOrder)
		if (num > 0) {
			currentOrderId := currentOrder[0].Order.Id
			this.Data["buttonHref"] = "/Portal/passengerorderdetail/" + currentOrderId
			this.Data["buttonValue"] = "进入当前行程"
		} else {
			this.Data["buttonHref"] = "/Portal/searchinput"
			this.Data["buttonValue"] = "找拼车"
		}
		this.Data["buttonClass"] = "weui-btn weui-btn_primary"
	}

	this.Data["tabIndex"] = 1
	this.Data["orderNum"] = num
	this.Data["orderInfo"] =  orderDetailInfo
	this.Data["StatusText"] = odStatusText
	this.Data["onRoadType"] = onRoadType
	this.TplName = "passengerOrder.html"

}