package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"youdidi/models"
)

type ComplainController struct {
	beego.Controller
}

type ComplainContent struct {
	Utype int //0:用户 1:客服
	Content string
	Time string
}

// @router /portal/showcomplain [GET]
func (this *ComplainController) ShowComplain () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")

	var dbComplain models.Complain
	var cpInfo []*models.Complain

	num, err := dbComplain.GetComplainFromUser(userId, &cpInfo)

	if (err != nil) {
		logs.Error("get complain from user fail uid=%v", userId)
		this.Data["num"] = 0
	} else {
		this.Data["num"] = num
	}

	for i, v := range cpInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cpInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["list"] = cpInfo
	this.Data["tabIndex"] = 3

	this.TplName = "showComplain.html"

}

// @router /portal/newcomplain [GET]
func (this *ComplainController) NewComplain () {
	this.Data["tabIndex"] = 3
	this.TplName = "newComplain.html"
}

// @router /portal/donewcomplain [POST]
func (this *ComplainController) DoNewComplain () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	userIdInt, _ := strconv.Atoi(userId)
	content := this.GetString("content")
	title := this.GetString("title")

	contentItem := &ComplainContent{}
	var contentS []*ComplainContent

	contentItem.Time = strconv.FormatInt(time.Now().Unix(),10)
	contentItem.Utype = 0
	contentItem.Content = content

	contentS = append(contentS, contentItem)

	data, _ := json.Marshal(&contentS)

	code := 0
	msg := ""

	var dbC models.Complain

	dbC.User = &models.User{Id:userIdInt}
	dbC.Content = string(data)
	dbC.Title = title
	dbC.Status = 0
	dbC.Time = strconv.FormatInt(time.Now().Unix(),10)

	logs.Debug("content=%v", string(data))

	_,err := dbC.Insert()

	if (err != nil) {
		code = 1
		msg = "网络开小差了哦，请稍后重试"
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

// @router /portal/complaindetail/:id/:utype [GET]
func (this *ComplainController) ComplainDetail () {
	id := this.GetString(":id")
	utype := this.GetString(":utype")

	var dbDc models.Complain
	var dcInfo []*models.Complain

	num, err := dbDc.GetComplainFromId(id, &dcInfo)

	if (err != nil) {
		logs.Error("get complain info fail id=%v", id)
		this.Data["num"] = 0
	} else {
		this.Data["num"] = num

		var content  []*ComplainContent
		err := json.Unmarshal([]byte(dcInfo[0].Content), &content)
		if (err != nil) {
			logs.Error("pares complain json content fail err=%v", err.Error())
			this.Data["num"] = 0
		} else {
			for i, v := range content {
				this.Data["launchTime"] = v.Time;
				launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
				tm := time.Unix(launchTime64, 0)
				content[i].Time = tm.Format("2006-01-02 15:04")
			}
			this.Data["content"] = content
		}
		this.Data["list"] = dcInfo[0]
		this.Data["uType"] = utype
	}
	this.Data["tabIndex"] = 3
	this.TplName = "complainDetail.html"
}


// @router /portal/replycomplain [POST]
func (this *ComplainController) ReplyComplain () {
	replycontent := this.GetString("content")
	id := this.GetString("id")
	uType, _ := this.GetInt("type")

	code := 0
	msg := ""

	var dbDc models.Complain
	var dcInfo []*models.Complain

	num, err := dbDc.GetComplainFromId(id, &dcInfo)

	if (num < 1 || err != nil) {
		logs.Error("get complain info fail id=%v", id)
		code = 1
		msg = "网络开小差了哦，请稍后重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	var contentS []*ComplainContent
	err1 := json.Unmarshal([]byte(dcInfo[0].Content), &contentS)
	if (err1 != nil) {
		logs.Error("pares complain json content fail err=%v", err.Error())
		code = 2
		msg = "网络开小差了哦，请稍后重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	contentItem := &ComplainContent{}
	contentItem.Time = strconv.FormatInt(time.Now().Unix(),10)
	contentItem.Utype = uType
	contentItem.Content = replycontent

	contentS = append(contentS, contentItem)

	data, _ := json.Marshal(&contentS)

	_, err2 := dbDc.UpdateComplain(id, string(data), uType)

	if (err2 != nil) {
		logs.Error("pares complain json content fail err=%v", err.Error())
		code = 3
		msg = "网络开小差了哦，请稍后重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}


	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}