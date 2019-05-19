package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
)
type OrderController struct {
	beego.Controller
}

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
