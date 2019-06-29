package models

import (
	"github.com/astaxie/beego/orm"
)

func (u *Chat) TableName() string {
	return "Chat"
}

func (u *Chat) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Chat) GetAllMsg (oid string, pid string, list *[]*Chat) (int64, error){
	return orm.NewOrm().QueryTable(u).Filter("order_id", oid).Filter("passenger_id", pid).OrderBy("TimeStamp").All(list)
}