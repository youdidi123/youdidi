package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
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
		Filter("status__lt",3). //乘客当前的行程只会有比3小的
		Limit(1).
		All(list)
	return num
}

func (u *Order_detail) GetOrderInfoFromPassengerId (oid string , uid string , list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("order_id" , oid).
		Filter("passage_id" , uid).
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
	refuseNum int , balance float64, siteNum int , price float64) bool {
	var orderInfo Order
	var userInfo User
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

	_ , err2 := o.QueryTable(u).Filter("id" , odid).Update(orm.Params{
		"Status": 5,
	})
	if (err2 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	_ , err3 := o.QueryTable(userInfo).Filter("id", pid).Update(orm.Params{
		"OnRoadType": 0,"Balance":balance+float64(siteNum)*price,
	})
	if (err3 != nil) {
		logs.Error("update passenger info fail id=%v" , pid)
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
