package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

func (u *Order_detail) TableName() string {
	return "Order_detail"
}

func (u *Order_detail) GetOrderDetailFromPassengerId (id string , list *[]*Order_detail) int64 {
	num , err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("passage_id", id).OrderBy("Status").All(list)

	if (err != nil) {
		logs.Error("query order detail from passenger fail pid=%v" , id)
	}
	return num
}

func (u *Order_detail) GetCurrentPassengerOrderFromUserId (id string, list *[]*Order_detail) int64{
	num , _ := orm.NewOrm().QueryTable(u).
		Filter("passage_id", id).
		Filter("status__lt",4). //乘客当前的行程只会有比4小的
		Limit(1).
		All(list)
	return num
}

func (u *Order_detail) GetOrderedOrderFromPassengerId (oid string , uid string , list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("order_id" , oid).
		Filter("passage_id" , uid).
		Filter("status__lt" , 5).
		All(list)
	return num
}

func (u *Order_detail) GetOrderDetailFromOrderId (id string , list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("order_id" , id).OrderBy("status").
		All(list)
	return num
}

func (u *Order_detail) GetOrderDetailFromId (id string ,list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("id" , id).
		All(list)
	return num
}

func (u *Order_detail) AgreeRequest(odid string , oid string , confirmPnum int , siteNum int) bool {
	var orderInfo Order
	o := orm.NewOrm()
	o.Begin()

	_ , err1 := o.QueryTable(orderInfo).Filter("id", oid).Update(orm.Params{
		"ConfirmPnum": confirmPnum+siteNum,
	})
	if (err1 != nil) {
		logs.Error("add confirmPnum fail oid=%v confirmPnum=%v" , oid , confirmPnum+1)
		o.Rollback()
		return false
	}

	_ , err2 := o.QueryTable(u).Filter("id" , odid).Update(orm.Params{
		"Status": 1,
	})
	if (err2 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	err3 := o.Commit()

	if (err3 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	return true
}

func (u *Order_detail) RefuseRequest(odid string , oid string , pid string , requestNum int ,
	refuseNum int , siteNum int , price float64) bool {
	var orderInfo Order
	var userInfo User
	var userInfos []*User
	o := orm.NewOrm()
	o.Begin()

	_ , err1 := o.QueryTable(orderInfo).Filter("id", oid).Update(orm.Params{
		"RefusePnum": refuseNum+siteNum,"RequestPnum":requestNum-siteNum,
	})
	if (err1 != nil) {
		logs.Error("add confirmPnum fail oid=%v RefusePnum=%v RequestPnum=%v" , oid , refuseNum+1 , requestNum-1)
		o.Rollback()
		return false
	}

	_ , err2 := o.QueryTable(u).Filter("id" , odid).ForUpdate().Update(orm.Params{
		"Status": 5,
	})
	if (err2 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	numuser, erruser := o.QueryTable(userInfo).Filter("id", pid).ForUpdate().All(&userInfos)
	if (numuser < 0 || erruser != nil) {
		logs.Error("get user info fail pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}

	balance, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", userInfos[0].Balance + price * float64(siteNum)), 64)
	logs.Debug("user balance=%v price=%v siteNum=%v",userInfos[0].Balance, price, siteNum)
	_ , err3 := o.QueryTable(userInfo).Filter("id", pid).ForUpdate().Update(orm.Params{
		"OnRoadType": 0,"Balance":balance,
	})
	if (err3 != nil) {
		logs.Error("update passenger info fail id=%v" , pid)
		o.Rollback()
		return false
	}

	var accountFlow Account_flow
	accountFlow.Type = 4
	accountFlow.User = &User{Id:userInfos[0].Id}
	accountFlow.Oid = oid
	accountFlow.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", price * float64(siteNum)), 64)
	accountFlow.Time = strconv.FormatInt(time.Now().Unix(),10)
	accountFlow.Balance = balance

	_, err8 := o.Insert(&accountFlow)
	if (err8 != nil) {
		logs.Error("record account flow fail oid=%v pid=%v", u.Id, userInfos[0].Id)
		o.Rollback()
		return false
	}

	err4 := o.Commit()

	if (err4 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	return true
}

func (u *Order_detail) PassengerConfirm(odid string) bool {
	o := orm.NewOrm()
	o.Begin()

	var odInfo []*Order_detail
	num, err1 := o.QueryTable(u).Filter("Id", odid).ForUpdate().RelatedSel().All(&odInfo)
	if (err1 != nil || num < 1) {
		logs.Error("get order detail fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	driverId := odInfo[0].Driver.Id
	passengerId := odInfo[0].Passage.Id

	var dbUser User
	var driverInfo []*User
	var passengerInfo []*User

	num2, err2 := o.QueryTable(dbUser).Filter("Id", driverId).ForUpdate().All(&driverInfo)
	num3, err3 := o.QueryTable(dbUser).Filter("Id", passengerId).ForUpdate().All(&passengerInfo)

	if (num2 < 1 || num3 < 1 || err2 != nil || err3 != nil) {
		logs.Error("get userinfo fail fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	_, err4 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
		"Status": 4,
	})
	if (err4 != nil) {
		logs.Error("set order detail status to 4 fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	//修改乘客数据
	orderNumAsP := passengerInfo[0].OrderNumAsP + 1
	_, err5 := o.QueryTable(dbUser).Filter("Id", passengerId).Update(orm.Params{
		"OrderNumAsP": orderNumAsP,
		"OnRoadType": 0,
	})
	if (err5 != nil) {
		logs.Error("update passanger info fail pid=%v" , passengerId)
		o.Rollback()
		return false
	}

	//修改车主数据
	balance := driverInfo[0].Balance + float64(odInfo[0].SiteNum) * odInfo[0].Order.Price
	balance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", balance), 64)
	_, err6 := o.QueryTable(dbUser).Filter("Id", driverId).Update(orm.Params{
		"Balance": balance,
	})
	if (err6 != nil) {
		logs.Error("update driver info fail pid=%v" , driverId)
		o.Rollback()
		return false
	}

	//添加account_flow数据
	var afDriver Account_flow
	var afPassenger Account_flow

	afDriver.Balance = balance
	afDriver.Oid = odInfo[0].Order.Id
	afDriver.User = driverInfo[0]
	afDriver.Type = 5
	afDriver.Time = strconv.FormatInt(time.Now().Unix(),10)
	afDriver.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(odInfo[0].SiteNum) * odInfo[0].Order.Price), 64)

	afPassenger.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(odInfo[0].SiteNum) * odInfo[0].Order.Price), 64)
	afPassenger.Time = strconv.FormatInt(time.Now().Unix(),10)
	afPassenger.Oid = odInfo[0].Order.Id
	afPassenger.Type = 3
	afPassenger.User = passengerInfo[0]

	_, err7 := o.Insert(&afDriver)
	_, err8 := o.Insert(&afPassenger)

	if (err7 != nil || err8 != nil) {
		logs.Error("insert account flow fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	return true
}