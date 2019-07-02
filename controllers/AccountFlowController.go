package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"youdidi/models"
)

type AccountFlowController struct {
	beego.Controller
}

var (
	accountFlowType = []struct {
		Text string
	}{{"用户充值"}, {"用户提现"}, {"预付款"}, {"确认付款"}, {"退款"}, {"收款"},
		{"付违约款"}, {"收违约款"}, {"用户提现到账"}, {"平台信息费"}}
)

// @router /Portal/accountflow [GET]
func (this *AccountFlowController) GetAccountFlow (){
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")

	var dbAccountFlow models.Account_flow
	var accountInfo []*models.Account_flow

	num := dbAccountFlow.GetAccountInfoFromUserId(userId, &accountInfo)

	for i, v := range accountInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		accountInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["info"] = accountInfo
	this.Data["type"] = accountFlowType
	this.Data["tabIndex"] = 3
	this.TplName = "userAccountFlow.html"

}


// @router /Portal/invest [GET]
func (this *AccountFlowController) Invest (){
	this.Data["tabIndex"] = 3
	this.TplName = "invest.html"
}

// @router /Portal/withdraw [GET]
func (this *AccountFlowController) Withdraw () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")

	var dbUser models.User
	var userInfo []*models.User

	dbUser.GetUserInfoFromId(userId, &userInfo)

	this.Data["balance"] = userInfo[0].Balance
	this.Data["tabIndex"] = 3
	this.TplName = "withdraw.html"
}

// @router /Portal/dowithdraw [POST]
func (this *AccountFlowController) DoWithdraw () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	userIdInt, _ := strconv.Atoi(userId)
	money := this.GetString("money")

	code := 0
	msg := ""

	var dbAc models.Account_flow
	cacheFlowId := genOrderId(userIdInt)

	if (dbAc.DoWithDrew(userId, money, cacheFlowId)) {
		code = 0
		msg = "操作成功"
	} else {
		code = 1
		msg = "系统错误，请重试"
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /Portal/getOpenId [GET]
func (this *AccountFlowController) GetOpenId () {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")

	var dbUser models.User
	var userInfo []*models.User

	succ,num := dbUser.GetUserInfoFromId(userId, &userInfo)

	code := 0
	msg := ""

	if (succ != "true" || num < 1) {
		code = 1
		msg = "支付失败，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	code = 0
	msg = userInfo[0].OpenId
	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

