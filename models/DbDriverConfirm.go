package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func (u *Driver_confirm) TableName() string {
	return "Driver_confirm"
}


func (u *Driver_confirm) CreateDriverConfirm(userId string, num int64) bool {
	o := orm.NewOrm()
	o.Begin()
	var dbUser User

	_, err2 := o.QueryTable(dbUser).Filter("Id", userId).Update(orm.Params{
		"IsDriver": 1,
	})
	if (err2 != nil) {
		o.Rollback()
		logs.Error("update user info fail uid=%v", userId)
		return false
	}

	if (num > 0) {
		_, err3 := o.QueryTable(u).Filter("user_id", userId).Update(orm.Params{
			"Status": 0,
			"SfzImg":u.SfzImg,
			"DriverLiceseImg":u.DriverLiceseImg,
			"CarLiceseImg":u.CarLiceseImg,
			"SfzNum":u.SfzNum,
			"RealName":u.RealName,
			"CarType":u.CarType,
			"CarNum":u.CarNum,
			"Time":u.Time,
			"RejectReason":"",
		})
		if (err3 != nil) {
			o.Rollback()
			logs.Error("insert driver confirm info fail", userId)
			return false
		}
	} else {
		_, err3 := o.Insert(u)
		if (err3 != nil) {
			o.Rollback()
			logs.Error("insert driver confirm info fail", userId)
			return false
		}
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail")
		o.Rollback()
		return false
	}
	return true
}

func (u *Driver_confirm) GetUserOrder(userId string, list *[]*Driver_confirm) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("user_id", userId).OrderBy("-time").All(list)

	if(err != nil) {
		logs.Error("get driver confirm order fail uid=&v", userId)
	}
	return num
}

func (u *Driver_confirm) GetNoConfirm(list *[]*Driver_confirm) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Status", 0).OrderBy("Time").All(list)
	if (err != nil) {
		logs.Error("get no confirm driverConfirm fail")
	}
	return num
}

func (u *Driver_confirm) GetOrderFromId(id string , list *[]*Driver_confirm) int64 {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("Id", id).All(list)
	if (err != nil) {
		logs.Error("get no confirm driverConfirm fail")
	}
	return num
}