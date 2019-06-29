package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)


func (u *Admin_user) TableName() string {
	return "Admin_user"
}

func (u *Admin_user) GetUserInfoFromName(name string, list *[]*Admin_user) int64{
	num, err := orm.NewOrm().QueryTable(u).Filter("Name", name).All(list)
	if (err != nil) {
		logs.Error("get admin user info fail name=%v", name)
	}
	return num
}