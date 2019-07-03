package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func (u *Driver_confirm) TableName() string {
	return "Driver_confirm"
}


func (u *Driver_confirm) CreateDriverConfirm(userId string, num int64) (bool, int) {
	o := orm.NewOrm()
	o.Begin()
	var dbUser User
	oid := 0

	_, err2 := o.QueryTable(dbUser).Filter("Id", userId).Update(orm.Params{
		"IsDriver": 1,
	})
	if (err2 != nil) {
		o.Rollback()
		logs.Error("update user info fail uid=%v", userId)
		return false, 0
	}

	if (num > 0) {
		_, err3 := o.QueryTable(u).Filter("user_id", userId).Update(orm.Params{
			"Status": 0,
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
			return false, 0
		}
	} else {
		_, err3 := o.Insert(u)
		if (err3 != nil) {
			o.Rollback()
			logs.Error("insert driver confirm info fail", userId)
			return false, 0
		}
		oid = u.Id
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail")
		o.Rollback()
		return false, 0
	}
	return true, oid
}

func (u *Driver_confirm) UpdateImgFile(oid string, fileName string, iType string) (int64, error) {
	col := ""
	if (iType == "sfz") {
		col = "SfzImg"
	} else if (iType == "jsz") {
		col = "DriverLiceseImg"
	} else {
		col = "CarLiceseImg"
	}
	return orm.NewOrm().QueryTable(u).Filter("Id", oid).Update(orm.Params{
		col: fileName,
	})
}

func (u *Driver_confirm) GetUserOrder(userId string, list *[]*Driver_confirm) (int64, error) {
	num, err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("user_id", userId).OrderBy("-time").All(list)

	if(err != nil) {
		logs.Error("get driver confirm order fail uid=&v", userId)
	}
	return num, err
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

func (u *Driver_confirm) DoConfirmDriver(oid string, aType string, mark string) bool {
	o := orm.NewOrm()
	o.Begin()

	if (aType == "0") { //agree
		_, err1 := o.QueryTable(u).Filter("Id", oid).Update(orm.Params{
			"Status": 1,
		})
		if (err1 != nil) {
			logs.Error("update driver confirm order fail id=%v", oid)
			o.Rollback()
			return false
		}
		var dcInfo []*Driver_confirm

		num2, err2 := o.QueryTable(u).Filter("Id", oid).All(&dcInfo)
		if (num2 < 1 || err2 != nil) {
			logs.Error("get driver confirm order fail id=%v", oid)
			o.Rollback()
			return false
		}

		var dbUser User
		_, err3 := o.QueryTable(dbUser).Filter("Id", dcInfo[0].User.Id).Update(orm.Params{
			"IsDriver": 2,
			"SfzNum": dcInfo[0].SfzNum,
			"SfzImg": dcInfo[0].SfzImg,
			"DriverLiceseImg": dcInfo[0].DriverLiceseImg,
			"CarNum": dcInfo[0].CarNum,
			"CarLiceseImg": dcInfo[0].CarLiceseImg,
			"RealName": dcInfo[0].RealName,
			"CarType": dcInfo[0].CarType,
		})
		if (err3 != nil) {
			logs.Error("update user info fail id=%v uid=%v", oid, dcInfo[0].User.Id)
			o.Rollback()
			return false
		}
	} else {
		_, err1 := o.QueryTable(u).Filter("Id", oid).Update(orm.Params{
			"Status": 2,
			"RejectReason": mark,
		})
		if (err1 != nil) {
			logs.Error("update driver confirm order fail id=%v", oid)
			o.Rollback()
			return false
		}
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail id=%v" , oid)
		o.Rollback()
		return false
	}

	return true
}