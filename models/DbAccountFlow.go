package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
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