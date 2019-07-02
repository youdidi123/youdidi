package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)


func (u *Cash_flow) TableName() string {
	return "Cash_flow"
}

func (u *Cash_flow) GetReadyOrder(list *[]*Cash_flow) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Status", 0).Filter("Type", 1).All(list)
	if (err != nil) {
		logs.Debug("get ready cash flow order fail")
	}
	return num
}