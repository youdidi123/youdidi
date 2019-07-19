package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"youdidi/commonLib"
)


func (u *Cash_flow) TableName() string {
	return "Cash_flow"
}

func (u *Cash_flow) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Cash_flow) GetOrderInfo(orderId string, list *[]*Cash_flow) (success string , num int64){
	num, error := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Id", orderId).All(list)
	if (error != nil) {
		logs.Error("can not get order info from db orderId=%s ,error=%s" , orderId , error)
		return "false" , 0
	}
	return "true" , num
}

func (u *Cash_flow) GetWithdrewOrder(list *[]*Cash_flow, endTime int64, oType int, status int) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("Status", status).
		Filter("Type", oType).
		Filter("Time__lt", endTime).All(list)
	if (err != nil) {
		logs.Debug("get ready cash flow order fail")
	}
	return num
}

func (u *Cash_flow) UpdateWithDrewResult(succ bool, oid string, wxId string, reason string, status int) (int64, error) {
	if (succ) {
		return orm.NewOrm().QueryTable(u).Filter("Id", oid).Update(orm.Params{
			"Status": status,
			"WechatOrderId": wxId,
			"FinishTime": commonLib.GetCurrentTime(),
		})
	} else {
		return orm.NewOrm().QueryTable(u).Filter("Id", oid).Update(orm.Params{
			"Status": status,
			"FinishTime": commonLib.GetCurrentTime(),
			"RefuseReason": reason,
		})
	}
}

func (u *Cash_flow) DealWxPayRe(result_code string, err_code string, err_code_des string, openid string, wxId string, cfId string, total_fee int64, transaction_id string) bool {
	o := orm.NewOrm()
	o.Begin()

	var orderInfo []*Cash_flow
	var dbUser User
	var userInfo []*User
	var afInfo Account_flow
	currentTime := strconv.FormatInt(time.Now().Unix(),10)
	num1, err1 := o.QueryTable(u).RelatedSel().Filter("Id", cfId).ForUpdate().All(&orderInfo)
	if (num1 < 1 && err1 != nil) {
		logs.Info("get cashflow info fail id=%v", cfId)
		o.Rollback()
		return false
	}

	if (orderInfo[0].Status != 0) {
		logs.Info("order is already deal id=%v", cfId)
		o.Rollback()
		return true
	}

	if (orderInfo[0].User.OpenId != openid) {
		logs.Info("openid is not match reopenId=%v dbopenId=%v oid=%v", openid, orderInfo[0].User.OpenId, cfId)
		o.Rollback()
		return false
	}

	if (float64(total_fee) != (orderInfo[0].Money * 100)) {
		logs.Info("money is not match remoney=%v dbmoney=%v oid=%v", total_fee, orderInfo[0].Money, cfId)
		o.Rollback()
		return false
	}

	num2, err2 := o.QueryTable(dbUser).Filter("Id", orderInfo[0].User.Id).ForUpdate().All(&userInfo)
	if (num2 < 1 && err2 != nil) {
		logs.Info("get user info fail uid=%v oid=%v", orderInfo[0].User.Id, cfId)
		o.Rollback()
		return false
	}
	balance := userInfo[0].Balance

	afInfo.User = userInfo[0]
	afInfo.Money = orderInfo[0].Money
	afInfo.Time = currentTime
	afInfo.Oid = cfId

	if (result_code == "SUCCESS") {
		balance = balance + orderInfo[0].Money
		balance, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", balance), 64)

		_, err3 := o.QueryTable(dbUser).Filter("Id", orderInfo[0].User.Id).Update(orm.Params{
			"Balance": balance,
		})
		if (err3 != nil) {
			logs.Info("update user balance fail uid=%v oid=%v", orderInfo[0].User.Id, cfId)
			o.Rollback()
			return false
		}
		_, err4 := o.QueryTable(u).Filter("Id", cfId).Update(orm.Params{
			"Status": 1,
			"FinishTime":currentTime,
			"WechatOrderId":transaction_id,
		})
		if (err4 != nil) {
			logs.Info("update cf status to 1 fail oid=%v", cfId)
			o.Rollback()
			return false
		}
		afInfo.Type = 0
	} else if (result_code == "FAIL") {
		_, err4 := o.QueryTable(u).Filter("Id", cfId).Update(orm.Params{
			"Status": 2,
			"RefuseReason": err_code_des,
			"FinishTime":currentTime,
			"WechatOrderId":transaction_id,
		})
		if (err4 != nil) {
			logs.Info("update cf status to 1 fail oid=%v", cfId)
			o.Rollback()
			return false
		}
		afInfo.Type = 10
	} else {
		logs.Info("return code not SUCCESS or FAIL is %v oid=%v", result_code, cfId)
		o.Rollback()
		return false
	}
	afInfo.Balance = balance

	_, err5 := o.Insert(&afInfo)
	if (err5 != nil) {
		logs.Info("insert account flow fail oid=%v", cfId)
		o.Rollback()
		return false
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit wx pay re fail oid=%v" , cfId)
		o.Rollback()
		return false
	}
	moneyStr := strconv.FormatFloat(orderInfo[0].Money, 'G' , -1,64)
	balanceStr := strconv.FormatFloat(balance, 'G' , -1,64)
	if (result_code == "SUCCESS") {
		commonLib.SendMsg5(userInfo[0].OpenId, 4, "",
			"#173177", "", "",
			"#173177", "",
			"#22c32e", "账户充值",
			"#22c32e", "充值成功",
			"#173177", moneyStr,
			"#173177", balanceStr)
	} else {
		commonLib.SendMsg5(userInfo[0].OpenId, 4, "",
			"#173177", "", "",
			"#173177", "",
			"#22c32e","账户充值",
			"#ff0000","充值失败",
			"#173177", moneyStr,
			"#173177", balanceStr)
	}

	return true
}

func (u *Cash_flow) IsFirstWithdrew (uid string, tmBegin string) bool {
	re := false
	var tmp []*Cash_flow

	logs.Debug("today begin time %v", tmBegin)
	num, err := orm.NewOrm().QueryTable(u).Filter("user_id", uid).Filter("Time__gte", tmBegin).Filter("type__in", 1,2).All(&tmp)
	if (err != nil) {
		logs.Info("get user withdrew order fail uid=%v", uid)
	}
	if (num == 0) {
		re = true
	}

	return re
}

func (u *Cash_flow) DoWithDrew (uid string, oid string, money string) bool {
	o := orm.NewOrm()
	o.Begin()

	var invest []*Cash_flow
	var accountFlow Account_flow
	var cashFlow Cash_flow
	var dbUser User
	var userInfo []*User

	num1, err1 := o.QueryTable(dbUser).Filter("Id", uid).ForUpdate().All(&userInfo)

	if (num1 < 1 || err1 != nil) {
		logs.Error("get userInfo fail uid=%v", uid)
		o.Rollback()
		return false
	}

	money64, _ := strconv.ParseFloat(money, 64)

	if (money64 > userInfo[0].Balance) {
		logs.Error("balance is not enough for withdrew money=%v balance=%v", money64, userInfo[0].Balance)
		o.Rollback()
		return false
	}

	balance := userInfo[0].Balance - money64
	balance = commonLib.FormatMoney(balance)

	_, err2 := o.QueryTable(dbUser).Filter("Id", uid).Update(orm.Params{
		"Balance": balance,
	})

	if (err2 != nil) {
		logs.Error("update userInfo fail uid=%v", uid)
		o.Rollback()
		return false
	}

	currentTime := commonLib.GetCurrentTime()

	accountFlow.Balance = balance
	accountFlow.User = userInfo[0]
	accountFlow.Type = 1
	accountFlow.Money = money64
	accountFlow.Time = currentTime

	cashFlow.Money = money64
	cashFlow.Status = 0
	cashFlow.User = userInfo[0]
	cashFlow.Time = currentTime

	num3, err3 := o.QueryTable(u).RelatedSel().
		Filter("user_id", uid).
		Filter("Type", 0).
		Filter("Money", money).
		Filter("IsRefund", false).All(&invest)

	if (err3 != nil) {
		logs.Info("get invest history fail")
		o.Rollback()
		return false
	}

	if (num3 > 0) {
		//走退款逻辑
		investId := invest[0].Id
		oid = "t-" + investId
		cashFlow.Id = oid
		cashFlow.Type = 2
		cashFlow.InvestOid = investId
		_, err4 := o.QueryTable(u).Filter("Id", investId).Update(orm.Params{
			"IsRefund": true,
		})
		if (err4 != nil) {
			logs.Error("update cashflow fail id=%v", investId)
			o.Rollback()
			return false
		}
	}else {
		//走正常提现
		cashFlow.Type = 1
		cashFlow.Id = oid
	}

	_, err5 := o.Insert(&cashFlow)
	if (err5 != nil) {
		logs.Error("insert cashflow fail uid=%v", uid)
		o.Rollback()
		return false
	}

	accountFlow.Oid = oid

	_, err6 := o.Insert(&accountFlow)
	if (err6 != nil) {
		logs.Error("insert account fail uid=%v", uid)
		o.Rollback()
		return false
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("create withdrew order fail oid=%v" , oid)
		o.Rollback()
		return false
	}
	return true
}

func (u *Cash_flow) DealWxRefundRe (refund_id string, out_refund_no string, refund_status string, success_time string, settlement_refund_fee int64) bool {
	var cfInfo []*Cash_flow
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Id", out_refund_no).All(&cfInfo)
	if (num < 1 || err !=nil) {
		logs.Error("get refund order info fail oid=%v", out_refund_no)
		return false
	}
	if (int64(cfInfo[0].Money * 100) != settlement_refund_fee) {
		logs.Error("wx return money %v != refund money %v", settlement_refund_fee, cfInfo[0].Money)
		return false
	}

	if (refund_status == "SUCCESS") {
		_, err1 := orm.NewOrm().QueryTable(u).Filter("Id", out_refund_no).Update(orm.Params{
			"Status": 1,
			"FinishTime":success_time,
			"WechatOrderId":refund_id,
		})
		if (err1 != nil) {
			logs.Error("update success refund info fail err=%v", err.Error())
			return false
		}
		moneyStr := strconv.FormatFloat(cfInfo[0].Money, 'G' , -1,64)
		balanceStr := strconv.FormatFloat(cfInfo[0].User.Balance, 'G' , -1,64)
		commonLib.SendMsg5(cfInfo[0].User.OpenId, 4, "",
			"#173177", "", "",
			"#173177", "",
			"#22c32e","提现确认",
			"#22c32e","提现成功",
			"#173177", moneyStr,
			"#173177", balanceStr)
	} else {
		_, err1 := orm.NewOrm().QueryTable(u).Filter("Id", out_refund_no).Update(orm.Params{
			"Status": 2,
			"FinishTime":success_time,
			"WechatOrderId":refund_id,
			"RefuseReason":refund_status,
		})
		if (err1 != nil) {
			logs.Error("update success refund info fail err=%v", err.Error())
			return false
		}
	}
	return true
}
