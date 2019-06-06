package models

import (
	"github.com/astaxie/beego/orm"
)

func (u *Order) TableName() string {
	return "Order"
}

func (u *Order) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Order) GetOrderInfoFromUserId (id string, list *[]*Order) int64{
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel("SrcId","DestId").Filter("user_id", id).OrderBy("-LaunchTime").All(list)
	return num
}

func (u *Order) GetCurrentOrderFromUserId (id string, list *[]*Order) int64{
	num , _ := orm.NewOrm().QueryTable(u).Filter("user_id", id).Filter("status__lt",3).OrderBy("-LaunchTime").Limit(1).All(list)
	return num
}