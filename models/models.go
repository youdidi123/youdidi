package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
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

func init () {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")

	orm.RegisterModel(new(User))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlurls+")/"+mysqldb+"?charset=utf8&loc=Asia%2FShanghai")
	orm.RunSyncdb("default", false, true)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}


}