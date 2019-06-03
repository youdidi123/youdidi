package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
)
type OrderController struct {
	beego.Controller
}

//订单ID生成规则：精确到秒的日志字符串+用户ID（保证是6位，100000用户）+6位随机数字
func genOrderId (uid string) string {
	oid := ""
	ltime := time.Now().Format("20060102150405")
	uidInt, _ := strconv.Atoi(uid)
	if (uidInt < 100000) {
		uidInt = uidInt + 100000
		uid = strconv.Itoa(uidInt)
	}

	oid = ltime + uid + GetRandomCode()

	return oid
}

// @router /Portal/showDriverorder/ [GET]
func (this *OrderController) ShowDriverOrder () {
	this.Data["tabIndex"] = 2
	this.TplName = "driverOrder.html"
}

// @router /Portal/createorder [GET]
func (this *OrderController) CreateOrder () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	userInfo := GetUserInfoFromRedis(uid)
	this.Data["tabIndex"] = 2
	if (userInfo.IsDriver != 2) {
		orderNumWithoutVer,_ := strconv.Atoi(beego.AppConfig.String("orderNumWithoutVer"))
		this.TplName = "createOrderFilter.html"
		this.Data["num"] = orderNumWithoutVer - userInfo.OrderNumWV
		this.Data["orderNumWithoutVer"] = orderNumWithoutVer
	} else {
		this.TplName = "createOrder.html"
		this.Data["uid"] = uid
	}
}

// @router /Portal/createorderforce [GET]
func (this *OrderController) CreateOrderForce () {
	uid, _ := this.Ctx.GetSecureCookie("qyt","qyt_id")
	this.Data["tabIndex"] = 2
	this.TplName = "createOrder.html"
	this.Data["uid"] = uid

}