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
		Filter("status__lt",4). //乘客当前的行程只会有比4小的
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
