package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id    int64   `orm:"auto;pk;column(id);" json:"id"`
	WechatId  string  `orm:"size(256);column(wechatId)" json:"wechatId"`
	Name string `orm:"size(255);column(name)" json:"name"`
	Passwd string `orm:"size(255);column(passwd)" json:"passwd"`
	Phone string `orm:"size(255);column(phone)" json:"phone"`
	IsPhoneVer bool `orm:"size(4);column(isPhoneVer)" json:"isPhoneVer"`
	IsDriver bool `orm:"size(4);column(isDriver)" json:"isDriver"`
	IsDriverVer bool `orm:"size(4);column(isDriverVer)" json:"isDriverVer"`
	IsOnRoad bool `orm:"size(4);column(isOnRoad)" json:"isOnRoad"`
	Star int64 `orm:"size(11);column(star)" json:"star"`
	IsInternal bool `orm:"size(4);column(isInternal)" json:"isInternal"`
	Charge float64 `column(charge)" json:"charge"`
}

func (u *User) TableName() string {
	return "User"
}

func (u *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(u).RelatedSel()
}


func (u *User) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *User) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(u, fields...)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserInfo(name string, list *[]*User) (success string , num int64){
	num, error := orm.NewOrm().Raw("SELECT * from User where name = ?" , name).QueryRows(list)
	if (error != nil) {
		logs.Error("can not get user info from db name=%s ,error=%s" , name , error)
		return "false" , 0
	}
	return "true" , num
}