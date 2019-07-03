package models

import (
	"github.com/astaxie/beego/orm"
)

func (u *Complain) TableName() string {
	return "Complain"
}

func (u *Complain) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Complain) GetComplainFromUser (uid string, list *[]*Complain) (int64, error){
	return orm.NewOrm().QueryTable(u).Filter("user_id", uid).OrderBy("-Time").All(list)
}

func (u *Complain) GetComplainFromId (id string, list *[]*Complain) (int64, error){
	return orm.NewOrm().QueryTable(u).RelatedSel().Filter("Id", id).All(list)
}

func (u *Complain) UpdateComplain (id string, content string, status int) (int64, error) {
	return orm.NewOrm().QueryTable(u).Filter("Id", id).Update(orm.Params{
		"Content":        content,
		"Status":         status,
	})
}

func (u *Complain) GetNoComplain (list *[]*Complain) (int64, error){
	return orm.NewOrm().QueryTable(u).Filter("Status", 0).OrderBy("-Time").All(list)
}