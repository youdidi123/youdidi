package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)
func (u *Account_flow) TableName() string {
	return "Account_Flow"
}

func (u *Account_flow) GetAccountInfoFromUserId (uid string, flow *[]*Account_flow) int64 {
	num, err := orm.NewOrm().QueryTable(u).Filter("user_id", uid).OrderBy("-time", "-type").All(flow)
	if (err != nil) {
		logs.Error("get account flow fail uid=%v", uid)
	}
	return num
}

func (u *Account_flow) DoWithDrew (uid string, money string, cashFlowId string) bool {
	o := orm.NewOrm()
	o.Begin()

	//money64, _ := strconv.ParseFloat(money,64)

	var dbUser User
	var userInfo []*User
	var accountFlow Account_flow
	var cashFlow Cash_flow

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

	_, err2 := o.QueryTable(dbUser).Filter("Id", uid).Update(orm.Params{
		"Balance": balance,
	})

	if (err2 != nil) {
		logs.Error("update userInfo fail uid=%v", uid)
		o.Rollback()
		return false
	}

	currentTime := strconv.FormatInt(time.Now().Unix(),10)

	accountFlow.Balance = balance
	accountFlow.User = userInfo[0]
	accountFlow.Type = 1
	accountFlow.Money = money64
	accountFlow.Time = currentTime

	cashFlow.Money = money64
	cashFlow.Type = 1
	cashFlow.Status = 0
	cashFlow.User = userInfo[0]
	cashFlow.Time = currentTime
	cashFlow.Id = cashFlowId

	_, err3 := o.Insert(&cashFlow)

	if (err3 != nil) {
		logs.Error("insert cashflow fail uid=%v", uid)
		o.Rollback()
		return false
	}

	accountFlow.Oid = cashFlowId

	_, err4 := o.Insert(&accountFlow)
	if (err4 != nil) {
		logs.Error("insert account fail uid=%v", uid)
		o.Rollback()
		return false
	}


	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail uid=%v" , uid)
		o.Rollback()
		return false
	}
	return true
}