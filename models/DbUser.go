package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func (u *User) TableName() string {
	return "User"
}

func (u *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(u).RelatedSel()
}


func (u *User) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}


func (u *User) GetUserInfo(name string, list *[]*User) (success string , num int64){
	num, error := orm.NewOrm().Raw("SELECT * from User where name = ?" , name).QueryRows(list)
	if (error != nil) {
		logs.Error("can not get user info from db name=%s ,error=%s" , name , error)
		return "false" , 0
	}
	return "true" , num
}

func (u *User) UpdateInfo (id int64, key string , value string)  {
	o := orm.NewOrm()
	num , err := o.QueryTable(u).Filter("id", id).Update(orm.Params{
		key: value,
	})
	if (err != nil || num == 0) {
		logs.Error("update User fail id=%v key=%v value=%v" , id , key , value)
	}
}

func GetUserInfoFormName (name string)([]orm.ParamsList , error) {
	o := orm.NewOrm()
	var lsits []orm.ParamsList
	_ , err := o.Raw("SELECT * from User where name = ?" , name).ValuesList(&lsits)
	return lsits , err

}