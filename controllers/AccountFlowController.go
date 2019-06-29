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
func (this *UserCenterController) GetAccountFlow (){
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
