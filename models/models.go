package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Name string `orm:"column(name)" json:"name"`
	Passwd string `orm:"column(passwd)" json:"passwd"`
	Nickname string `orm:"column(nickname)" json:"nickname"`
	Sex int `orm:"column(sex)" json:"sex"`
	City string `orm:"column(city)" json:"city"`
	Province string `orm:"column(province)" json:"province"`
	Subscribe_time	string `orm:"column(subscribe_time)" json:"subscribe_time"`
	Unionid	string `orm:"column(unionid)" json:"unionid"`
	Subscribe bool	`orm:"column(subscribe)" json:"subscribe"`
	OpenId	string `orm:"column(openId)" json:"openId"`
	WechatImg string `orm:"column(wechatImg)" json:"wechatImg"`
	Phone string `orm:"column(phone)" json:"phone"`
	IsPhoneVer bool `orm:"column(isPhoneVer)" json:"isPhoneVer"`
	IsDriver int `orm:"column(isDriver)" json:"isDriver"`
	SfzNum string `orm:"column(sfzNum)" json:"sfzNum"`
	SfzImg string `orm:"column(sfzImg)" json:"sfzImg"`
	DriverLiceseImg string `orm:"column(driverLiceseImg)" json:"driverLiceseImg"`
	CarNum string `orm:"column(carNum)" json:"carNum"`
	CarLiceseImg string `orm:"column(carLiceseImg)" json:"carLiceseImg"`
	OrderNumWV int `orm:"column(orderNumWV)" json:"orderNumWV"`
	IsVer bool `orm:"column(isVer)" json:"isVer"`
	RealName string `orm:"column(realName)" json:"realName"`
	IsStaff bool `orm:"column(isStaff)" json:"isStaff"`
	OnRoadType int `orm:"column(onRoadType)" json:"onRoadType"`
	OrderNumAsP int `orm:"column(orderNumAsP)" json:"orderNumAsP"`
	OrderNumAsD int `orm:"column(orderNumAsD)" json:"orderNumAsD"`
	StarAsP int `orm:"column(starAsP)" json:"starAsP"`
	StarAsD int `orm:"column(starAsD)" json:"starAsD"`
	CancleOasP int `orm:"column(cancleOasP)" json:"cancleOasP"`
	CancleOasD int `orm:"column(cancleOasD)" json:"cancleOasD"`
	Balance float64 `orm:"column(balance)" json:"balance"`
	IsVip bool `orm:"column(isVip)" json:"isVip"`
	VipDate string `orm:"column(vipDate)" json:"vipDate"`
	Orders []*Order `orm:"reverse(many)"`
	Order_details []*Order_detail `orm:"reverse(many)"`
}

type Order struct {
	Id string `orm:"pk;column(id);" json:"id"`
	User *User `json:"user" orm:"rel(fk)"`
	CreateTime string `column(createTime);" json:"createTime"`
	LaunchTime string `column(launchTime);" json:"launchTime"`
	PNum int `column(pNum);" json:"pNum"`
	SrcL string `column(srcL);" json:"srcL"`
	ThroughL string `column(throughL);" json:"throughL"`
	DestL string `column(destL);" json:"destL"`
	Status int `column(status);" json:"status"`
	RequestPnum int `column(requestPnum);" json:"requestPnum"`
	ConfirmPnum int `column(confirmPnum);" json:"confirmPnum"`
	RefusePnum int `column(refusePnum);" json:"refusePnum"`
	CanclePnum int `column(canclePnum);" json:"canclePnum"`
	OnroadPnum int `column(onroadPnum);" json:"onroadPnum"`
	PayedPnum int `column(payedPnum);" json:"payedPnum"`
	Price float64 `column(price);" json:"price"`
	marks string `column(marks);" json:"marks"`
	CancleReason string `column(cancleReason);" json:"cancleReason"`
	Order_locations []*Order_location `orm:"reverse(many)"`
	Order_details []*Order_detail `orm:"reverse(many)"`
}

type Location struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Name string `column(name);" json:"name"`
	Level1 string `column(level1);" json:"level1"`
	Level2 string `column(level2);" json:"level2"`
	Order_locations []*Order_location `orm:"reverse(many)"`
}

type Order_location struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Order *Order `json:"order" orm:"rel(fk)"`
	Location *Location `json:"location" orm:"rel(fk)"`
	LocationCustom string `column(locationCustom);" json:"locationCustom"`
	Type int `column(type);" json:"type"`
	LauchTime string `column(lauchTime);" json:"lauchTime"`
	Status int `column(status);" json:"status"`
}

type Order_detail struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Order *Order `json:"order" orm:"rel(fk)"`
	Driver *User `json:"driver" orm:"rel(fk)"`
	Passage *User `json:"passage" orm:"rel(fk)"`
	ModifyPrice float64 `column(ModifyPrice);" json:"ModifyPrice"`
	isModifyConfirm bool `column(isModifyConfirm);" json:"isModifyConfirm"`
	DStarNum int `column(dStarNum);" json:"dStarNum"`
	DCommit string `column(dCommit);" json:"dCommit"`
	IsDcommit bool `column(isDcommit);" json:"isDcommit"`
	PStarNum int `column(pStarNum);" json:"pStarNum"`
	PCommit string `column(pCommit);" json:"pCommit"`
	IsPcommit bool `column(isPcommit);" json:"isPcommit"`
	Status int `column(status);" json:"status"`
	IsPayed bool `column(isPayed);" json:"isPayed"`
	CancleReason string `column(cancleReason);" json:"cancleReason"`
	Chat string `column(chat);" json:"chat"`
}


func init () {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")

	//所有的数据表需要在这里注册
	orm.RegisterModel(new(User),new(Order),new(Location),new(Order_detail),new(Order_location))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlurls+")/"+mysqldb+"?charset=utf8&loc=Asia%2FShanghai")
	orm.RunSyncdb("default", false, true)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}


}