package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func (u *LocationRank) TableName() string {
	return "Location_rank"
}

func  (u *LocationRank) InsertData(src int, dest int) {
	var info []*LocationRank
	num, err := orm.NewOrm().QueryTable(u).Filter("Src", src).Filter("Dest",dest).All(&info)
	if (err != nil) {
		logs.Error("seacher location ranker fail startcode=%v, endcode=%v, err=%v", src, dest, err.Error())
	} else {
		if (num == 0) {
			var i LocationRank
			i.Dest = &Location{Id:int64(dest)}
			i.Src = &Location{Id:int64(src)}
			i.Num = 1
			orm.NewOrm().Insert(&i)
		} else {
			orm.NewOrm().QueryTable(u).Filter("Src", src).Filter("Dest",dest).Update(orm.Params{"Num":info[0].Num + 1})
		}
	}
}

func  (u *LocationRank) GetTop5(list *[]*LocationRank)(int64, error) {
	return orm.NewOrm().QueryTable(u).RelatedSel().OrderBy("-Num").Limit(5).All(list)
}
