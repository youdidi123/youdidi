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
	}{{"接单中"}, {"车主到达出发地"}, {"完成"}, {"已取消"}}
	odStatusText = []struct {
		Text string
	}{{"发起拼车请求"},{"请求已确认"}, {"行程中"}, {"待乘客确认"},{"完成"},{"拒绝请求"},{"乘客取消"},  {"车主取消"}}
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
	var user models.User
	var userInfo []*models.User

	succ, num := user.GetUserInfoFromId(uid, &userInfo)

	if (succ != "true" || num < 1) {
		this.TplName = "createOrder.html"
		this.Data["uid"] = uid
	}
	onRoadType := userInfo[0].OnRoadType
	if (onRoadType != 0) {
		logs.Debug("user %v onRoadTyep %v" , uid , onRoadType)
		//如果用户正在一个行程中，则跳转回查看行程页面
		this.ShowDriverOrder()
		return
	}

	this.Data["tabIndex"] = 2
	if (userInfo[0].IsDriver != 2) {
		//如果用户还没有认证司机，则进入认证提醒页面
		orderNumWithoutVer,_ := strconv.Atoi(beego.AppConfig.String("orderNumWithoutVer"))
		this.TplName = "createOrderFilter.html"
		//如果num=0 ，继续发单按钮为disable
		this.Data["num"] = orderNumWithoutVer - userInfo[0].OrderNumWV
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
	var dbUser models.User
	var userInfo []*models.User


	userIdS := this.GetString("uid")

	_, num := dbUser.GetUserInfoFromId(userIdS, &userInfo)

	if (num < 1) {
		code = 3
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	onRoadType := userInfo[0].OnRoadType
	disableTime := userInfo[0].DisableTime
	currentTime := time.Now().Unix()

	if (disableTime != "") {
		disableTime64, _ := strconv.ParseInt(disableTime, 10, 64)
		if (disableTime64 > currentTime) {
			tm := time.Unix(disableTime64, 0)
			tmText := tm.Format("2006-01-02 15:04")
			code = 4
			msg = "截止"+tmText+"前，禁止发单"
			this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
			this.ServeJSON()
			return
		}
	}

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

	launchTimeUnix, _ := time.ParseInLocation("2006-01-02 15:04:05", launchTime, time.Local)
	logs.Debug("launchTime=%v launchTime=%v launchTime=%v", launchTimeUnix, launchTime, launchTimeUnix.Unix())

	if (launchTimeUnix.Unix() - currentTime < 10 * 60) {
		logs.Debug("launchTime=%v currentTime=%v", launchTimeUnix, currentTime)
		code = 2
		msg = "请选择至少10分钟后的出发时间"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	orderId := genOrderId(userId)

	var dbOrder models.Order

	dbOrder.Id = orderId
	dbOrder.User  = &models.User{Id:userId}
	dbOrder.LaunchTime = strconv.FormatInt(launchTimeUnix.Unix(),10)
	dbOrder.CreateTime = strconv.FormatInt(currentTime,10)
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

	var dbOrder models.Order
	var orderInfo []*models.Order

	onum := dbOrder.GetOrderFromId(oid, &orderInfo)
	this.Data["onum"] = onum

	if (onum > 0) {
		uid, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
		if (uid != strconv.Itoa(orderInfo[0].User.Id)) {
			this.Data["isDriver"] = 0
		} else {
			this.Data["isDriver"] = 1
			var dbOrderDetail models.Order_detail
			var orderDetailInfo []*models.Order_detail

			odnum := dbOrderDetail.GetOrderDetailFromOrderId(oid, &orderDetailInfo)

			for i, v := range orderInfo {
				this.Data["launchTime"] = v.LaunchTime;
				launchTime64, _ := strconv.ParseInt(v.LaunchTime, 10, 64)
				tm := time.Unix(launchTime64, 0)
				orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04")
			}

			this.Data["odnum"] = odnum
			this.Data["odlist"] = orderDetailInfo
			this.Data["oinfo"] = orderInfo[0]
			this.Data["statustext"] = statusText
			this.Data["odstatustest"] = odStatusText
		}
	}
	this.Data["tabIndex"] = 2
	this.TplName = "driverOrderDetail.html"
}

// @router /Portal/passengerorderdetail/:odid [GET]
func (this *OrderController) PassengerOrderDetail () {
	odid := this.GetString(":odid")

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	num := dbOrderDetail.GetOrderDetailFromId(odid, &orderDetailInfo)

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
	this.TplName = "passengerOrderDetail.html"
}

// @router /Portal/searchorder [GET,POST]
func (this *OrderController) SearchOrder () {
	//startCode := this.GetString("startCode")
	//endCode := this.GetString("endCode")
	startCode, _ := strconv.Atoi(this.GetString("startCode"))
	endCode, _ := strconv.Atoi(this.GetString("endCode"))
	//startCode64 , _ := strconv.ParseInt(startCode, 10, 64)
	//endCode64 , _ := strconv.ParseInt(endCode, 10, 64)
	startCodeLocation := startCode % 1000000
	endCodeLocation := endCode % 1000000

	launchTime := this.GetString("launchTime")
	launchTime = launchTime + " 00:00:00"
	tmStart, _ := time.ParseInLocation("2006-01-02 15:04:05", launchTime, time.Local)
	tmEnd := tmStart.Unix() + (1*24*60*60)

	logs.Debug("search order launchTime=%v start=%v end=%v", launchTime , startCode , endCode)
	var dbOrder models.Order
	var orderInfo []*models.Order

	num := dbOrder.GetReadyOrders(&orderInfo , startCode , endCode ,
		tmStart.Unix(), tmEnd , startCodeLocation , endCodeLocation)

	for i , v := range orderInfo{
		launchTime64 , _ := strconv.ParseInt(v.LaunchTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		orderInfo[i].LaunchTime = tm.Format("2006-01-02 15:04")
		orderInfo[i].User.Phone = CryptionPhoneNum(v.User.Phone)
		orderInfo[i].User.CarNum = CryptionCarNum(v.User.CarNum)
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
		logs.Debug("onroadType=%v",onRoadType)
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	numOd := dbOrderDetail.GetOrderedOrderFromPassengerId(oid , userIdS , &orderDetailInfo)

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
	if (count > (orderInfo[0].PNum - orderInfo[0].RequestPnum)) {
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

	if (orderInfo[0].DoRequire(od, userIdS, count , mark)) {
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
			currentOrderId := strconv.Itoa(currentOrder[0].Id)
			logs.Debug("current odid is %v",currentOrderId)
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

// @router /Portal/agreerequest [POST]
func (this *OrderController) AgreeRequest () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString("oid")
	pid := this.GetString("pid")
	odid := this.GetString("odid")

	code := 0
	msg := ""

	var dbOrder models.Order
	var orderInfo []*models.Order

	if (! SetOrderLock(oid)) {
		code = 1
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	num := dbOrder.GetOrderFromId(oid , &orderInfo)

	if (num != 1) {
		DelOrderLock(oid)
		code = 2
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(orderInfo[0].User.Id)) {
		DelOrderLock(oid)
		code = 3
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	numod := dbOrderDetail.GetOrderDetailFromId(odid , &orderDetailInfo)

	if (numod != 1) {
		DelOrderLock(oid)
		code = 4
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (pid != strconv.Itoa(orderDetailInfo[0].Passage.Id)) {
		DelOrderLock(oid)
		code = 7
		msg = "乘客信息不符"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (orderDetailInfo[0].Status != 0) {
		DelOrderLock(oid)
		code = 5
		msg = "该乘客请求以处理完成，请勿重复操作"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (dbOrderDetail.AgreeRequest(odid , oid , orderInfo[0].ConfirmPnum , orderDetailInfo[0].SiteNum)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 6
		msg = "系统错误，请重试"
	}

	DelOrderLock(oid)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/refuserequest [POST]
func (this *OrderController) RefuseRequest () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString("oid")
	pid := this.GetString("pid")
	odid := this.GetString("odid")

	code := 0
	msg := ""

	var dbOrder models.Order
	var orderInfo []*models.Order

	if (! SetOrderLock(oid)) {
		code = 1
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	num := dbOrder.GetOrderFromId(oid , &orderInfo)

	if (num != 1) {
		DelOrderLock(oid)
		code = 3
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(orderInfo[0].User.Id)) {
		DelOrderLock(oid)
		code = 4
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	var dbOrderDetail models.Order_detail
	var orderDetailInfo []*models.Order_detail

	numod := dbOrderDetail.GetOrderDetailFromId(odid , &orderDetailInfo)

	if (numod != 1) {
		DelOrderLock(oid)
		code = 5
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (pid != strconv.Itoa(orderDetailInfo[0].Passage.Id)) {
		DelOrderLock(oid)
		code = 6
		msg = "乘客信息不符"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (orderDetailInfo[0].Status != 0) {
		DelOrderLock(oid)
		code = 7
		msg = "该乘客请求以处理完成，请勿重复操作"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (dbOrderDetail.RefuseRequest(odid , oid , pid , orderInfo[0].RequestPnum ,
		orderInfo[0].RefusePnum , orderDetailInfo[0].SiteNum , orderInfo[0].Price)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 9
		msg = "系统错误，请重试"
	}

	DelOrderLock(oid)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

// @router /Portal/getstart [POST]
func (this *OrderController) DriverGetStart () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString("oid")

	code := 0
	msg := ""

	var dbOrder models.Order
	var orderInfo []*models.Order

	if (! SetOrderLock(oid)) {
		code = 1
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	num := dbOrder.GetOrderFromId(oid , &orderInfo)

	if (num != 1) {
		DelOrderLock(oid)
		code = 2
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(orderInfo[0].User.Id)) {
		DelOrderLock(oid)
		code = 3
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (dbOrder.DriverGetStart(oid)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 4
		msg = "系统错误，请重试"
	}

	DelOrderLock(oid)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/getend [POST]
func (this *OrderController) DriverGetEnd () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString("oid")

	code := 0
	msg := ""

	var dbOrder models.Order
	var orderInfo []*models.Order

	if (! SetOrderLock(oid)) {
		code = 1
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	num := dbOrder.GetOrderFromId(oid , &orderInfo)

	if (num != 1) {
		DelOrderLock(oid)
		code = 2
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(orderInfo[0].User.Id)) {
		DelOrderLock(oid)
		code = 3
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (dbOrder.DriverGetEnd(oid, userId)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 4
		msg = "系统错误，请重试"
	}

	DelOrderLock(oid)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/drivercancle [POST]
func (this *OrderController) DriverCancle () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	oid := this.GetString("oid")
	confirmNum, _ := this.GetInt("confirmNum")

	code := 0
	msg := ""

	var dbOrder models.Order
	var orderInfo []*models.Order

	if (! SetOrderLock(oid)) {
		code = 1
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	num := dbOrder.GetOrderFromId(oid , &orderInfo)

	if (num != 1) {
		DelOrderLock(oid)
		code = 2
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(orderInfo[0].User.Id)) {
		DelOrderLock(oid)
		code = 3
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (dbOrder.DriverCancle(oid, confirmNum, userId)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 4
		msg = "系统错误，请重试"
	}

	DelOrderLock(oid)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/passengerconfirm [POST]
func (this *OrderController) PassengerConfirm () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	odid := this.GetString("odid")

	code := 0
	msg := ""

	var dbOd models.Order_detail
	var odInfo []*models.Order_detail

	num1 := dbOd.GetOrderDetailFromId(odid, &odInfo)

	if (num1 != 1) {
		code = 1
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(odInfo[0].Passage.Id)) {
		code = 2
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	orderId := odInfo[0].Order.Id

	if (! SetOrderLock(orderId)) {
		code = 3
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	if (dbOd.PassengerConfirm(odid)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 4
		msg = "系统错误，请重试"
	}

	DelOrderLock(orderId)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/passengercancle [POST]
func (this *OrderController) PassengerCancle () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	odid := this.GetString("odid")

	code := 0
	msg := ""

	var dbOd models.Order_detail
	var odInfo []*models.Order_detail

	num1 := dbOd.GetOrderDetailFromId(odid, &odInfo)

	if (num1 != 1) {
		code = 1
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (userId != strconv.Itoa(odInfo[0].Passage.Id)) {
		code = 2
		msg = "这个行程不属于你哦"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	orderId := odInfo[0].Order.Id

	if (! SetOrderLock(orderId)) {
		code = 3
		msg = "请求超时，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	if (dbOd.PassengerCancle(odid)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 4
		msg = "系统错误，请重试"
	}

	DelOrderLock(orderId)
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/recommand/:odid/:uType [GET]
func (this *OrderController) Recommand () {
	odid := this.GetString(":odid")
	uType := this.GetString(":uType")

	var dbOd models.Order_detail
	var odInfo []*models.Order_detail

	num := dbOd.GetOrderDetailFromId(odid, &odInfo)

	if (num < 1) {
		this.Data["isRecommand"] = false
		this.Data["starNum"] = 1
		this.Data["mark"] = ""
	} else {
		if (uType == "0") {
			this.Data["isRecommand"] = odInfo[0].IsDcommit
			this.Data["starNum"] = odInfo[0].DStarNum
			this.Data["mark"] = odInfo[0].DCommit
			this.Data["tabIndex"] = 1
		} else {
			this.Data["isRecommand"] = odInfo[0].IsPcommit
			this.Data["starNum"] = odInfo[0].PStarNum
			this.Data["mark"] = odInfo[0].PCommit
			this.Data["tabIndex"] = 2
		}
	}
	this.Data["odid"] = odid
	this.Data["uType"] = uType
	this.TplName = "recommand.html"
}

// @router /Portal/dorecommand [POST]
func (this *OrderController) DoRecommand () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	starNum, _ := this.GetInt("starNum")
	mark := this.GetString("mark")
	uType := this.GetString("uType") //乘客0 车主1
	odid := this.GetString("odid")

	code := 0
	msg := ""

	var dbOd models.Order_detail
	var odInfo []*models.Order_detail

	num := dbOd.GetOrderDetailFromId(odid, &odInfo)

	if (num != 1) {
		code = 1
		msg = "系统错误，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	if (uType == "0") {
		if (strconv.Itoa(odInfo[0].Passage.Id) != userId) {
			code = 2
			msg = "这个行程不属于你哦"
			this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
			this.ServeJSON()
			return
		}
		if (dbOd.Recommand(odid, uType, starNum, mark, odInfo[0].Driver.Id)) {
			code = 0
			msg = "操作成功"
		} else {
			code = 4
			msg = "系统错误，请重试"
		}

	} else {
		if (strconv.Itoa(odInfo[0].Driver.Id) != userId) {
			code = 3
			msg = "这个行程不属于你哦"
			this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
			this.ServeJSON()
			return
		}
		if (dbOd.Recommand(odid, uType, starNum, mark, odInfo[0].Passage.Id)) {
			code = 0
			msg = "操作成功"
		} else {
			code = 4
			msg = "系统错误，请重试"
		}
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}