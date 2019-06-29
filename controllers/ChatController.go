package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/models"
	"github.com/astaxie/beego/logs"
)

type ChatController struct {
	beego.Controller
}

// @router /Portal/chat/:pid/:oid/:type [GET]
func (this *ChatController) GoChat () {
	oid := this.GetString(":oid")
	pid := this.GetString(":pid")
	cType := this.GetString(":type")

	var chat models.Chat
	var chatInfo []*models.Chat

	_, err := chat.GetAllMsg(oid, pid, &chatInfo)

	if (err != nil) {
		logs.Error("get chat info fail oid=%v pid=%v" , oid, pid)
	}

	var dbOrder models.Order
	var orderInfo []*models.Order

	dbOrder.GetOrderFromId(oid, &orderInfo)

	var dbUser models.User
	var driverInfo []*models.User
	var passengerInfo []*models.User

	dbUser.GetUserInfoFromId(pid, &passengerInfo)
	dbUser.GetUserInfoFromId(strconv.Itoa(orderInfo[0].User.Id), &driverInfo)

	this.Data["pimg"] = passengerInfo[0].WechatImg
	this.Data["dimg"] = driverInfo[0].WechatImg
	this.Data["info"] = chatInfo
	this.Data["type"] = cType
	this.Data["oid"] = oid
	this.Data["pid"] = pid
	if (cType == "driver") {
		this.Data["title"] = "乘客:" + passengerInfo[0].Nickname
		this.Data["rType"] = 1
	} else {
		this.Data["title"] = "车主:" + driverInfo[0].Nickname
		this.Data["rType"] = 0
	}
	this.TplName = "chat.html"
}

// @router /Portal/refreshMsg [POST]
func (this *ChatController) RefreshMsg () {
	oid := this.GetString("oid")
	pid := this.GetString("pid")


	var chat models.Chat
	var chatInfo []*models.Chat

	_, err := chat.GetAllMsg(oid, pid, &chatInfo)

	if (err != nil) {
		this.Data["json"] = map[string]interface{}{"code":1, "msg":err.Error()};
		logs.Error("get chat info fail oid=%v pid=%v" , oid, pid)
	} else {
		this.Data["json"] = map[string]interface{}{"code":0, "data":chatInfo};
	}

	this.ServeJSON()
}

// @router /Portal/setMsg [POST]
func (this *ChatController) SetMsg () {
	oid := this.GetString("oid")
	pid, _ := this.GetInt("pid")
	cType, _ := this.GetInt("type")
	msg := this.GetString("msg")

	var chatInfo models.Chat
	chatInfo.Order = &models.Order{Id:oid}
	chatInfo.Passenger = &models.User{Id:pid}
	chatInfo.Content = msg
	chatInfo.Type = cType
	chatInfo.TimeStamp = strconv.FormatInt(time.Now().Unix(),10)

	_, err := chatInfo.Insert()

	if (err != nil) {
		this.Data["json"] = map[string]interface{}{"code":1, "msg":err.Error()};
	} else {
		this.Data["json"] = map[string]interface{}{"code":0, "msg":"success"};
	}
	this.ServeJSON()
}