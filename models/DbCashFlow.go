package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)


func (u *Cash_flow) TableName() string {
	return "Cash_flow"
}

func (u *Cash_flow) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Cash_flow) GetReadyOrder(list *[]*Cash_flow) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Status", 0).Filter("Type", 1).All(list)
	if (err != nil) {
		logs.Debug("get ready cash flow order fail")
	}
	return num
}

func (u *Cash_flow) DealWxPayRe(result_code string, err_code string, err_code_des string, openid string, wxId string, cfId string, total_fee int64) bool {
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
	return true
}