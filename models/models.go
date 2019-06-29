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
	CarType string `orm:"column(carType)" json:"carType"`
	DisableTime string `orm:"column(disableTime)" json:"disableTime"`
	Orders []*Order `orm:"reverse(many)"`
	Order_details []*Order_detail `orm:"reverse(many)"`
	Account_flows []*Account_flow `orm:"reverse(many)"`
}

type Order struct {
	Id string `orm:"pk;column(id);" json:"id"`
	User *User `json:"user" orm:"rel(fk)"`
	CreateTime string `column(createTime);" json:"createTime"`
	LaunchTime string `column(launchTime);" json:"launchTime"`
	PNum int `column(pNum);" json:"pNum"`
	SrcId *Location `json:"srcId" orm:"rel(fk)"`
	DestId *Location `json:"DestId" orm:"rel(fk)"`
	SrcLocationId int64 `column(srcLocationId);" json:"srcLocationId"`
	DestLocationId int64 `column(destLocationId);" json:"destLocationId"`
	ThroughL string `column(throughL);" json:"throughL"`
	Status int `column(status);" json:"status"`
	RequestPnum int `column(requestPnum);" json:"requestPnum"`
	ConfirmPnum int `column(confirmPnum);" json:"confirmPnum"`
	RefusePnum int `column(refusePnum);" json:"refusePnum"`
	CanclePnum int `column(canclePnum);" json:"canclePnum"`
	OnroadPnum int `column(onroadPnum);" json:"onroadPnum"`
	PayedPnum int `column(payedPnum);" json:"payedPnum"`
	Price float64 `column(price);" json:"price"`
	Marks string `column(marks);" json:"marks"`
	CancleReason string `column(cancleReason);" json:"cancleReason"`
	Order_details []*Order_detail `orm:"reverse(many)"`
	Chats []*Chat `orm:"reverse(many)"`
}

type Location struct {
	Id int64 `orm:"auto;pk;column(id);" json:"id"`
	Name string `column(name);" json:"name"`
	Level1 string `column(level1);" json:"level1"`
	Level2 string `column(level2);" json:"level2"`
	Orders []*Order `orm:"reverse(many)"`
}

type Order_detail struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Order *Order `json:"order" orm:"rel(fk)"`
	Driver *User `json:"driver" orm:"rel(fk)"`
	Passage *User `json:"passage" orm:"rel(fk)"`
	SiteNum int `column(siteNum);" json:"siteNum"`
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

type Account_flow struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Time string `column(time);" json:"time"`
	User *User `json:"order" orm:"rel(fk)"`
	Type int `column(type);" json:"type"` //0:用户充值 1:用户提现 2:预付款 3:确认付款 4:退款 5:收款 6:付违约款 7:收违约款 8:用户提现到账 9:平台信息费
	Oid string `column(oid);" json:"oid"` //0，1,8对应Cash_flow ID 其余对应Order ID
	Money float64 `column(money);" json:"money"`
	Balance float64 `column(balance);" json:"balance"` //记录每一次操作后的用户余额；type为3，6，8,9此字段不填
}

type Chat struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Order *Order `json:"order" orm:"rel(fk)"`
	Passenger *User `json:"passenger" orm:"rel(fk)"`
	Content string `column(content);" json:"content"`
	Type int `column(type);" json:"type"` //0:乘客发送的消息 1：司机发送的消息
	TimeStamp string `column(timeStamp);" json:"timeStamp"`
}

type Driver_confirm struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Time string `column(time);" json:"time"`
	RealName string `column(realName);" json:"realName"`
	CarNum string `column(carNum);" json:"carNum"`
	SfzNum string `column(sfzNum);" json:"sfzNum"`
	CarType string `column(carType);" json:"carType"`
	SfzImg string `orm:"column(sfzImg)" json:"sfzImg"`
	DriverLiceseImg string `orm:"column(driverLiceseImg)" json:"driverLiceseImg"`
	CarLiceseImg string `orm:"column(carLiceseImg)" json:"carLiceseImg"`
	Status int `column(status);" json:"status"` //0：审核中 1：审核通过 2：审核失败
	RejectReason string `column(rejectReason);" json:"rejectReason"`
	User *User `json:"user" orm:"rel(fk)"`
}

type Admin_user struct {
	Id int `orm:"auto;pk;column(id);" json:"id"`
	Name string `column(name);" json:"name"`
	Passwd string `column(passwd);" json:"passwd"`
	Phone string `column(phone);" json:"phone"`
	Email string `column(email);" json:"email"`
	Type int `column(type);" json:"type"`

}

func init () {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")

	//所有的数据表需要在这里注册
	orm.RegisterModel(
		new(User),
		new(Order),
		new(Location),
		new(Order_detail),
		new(Chat),
		new(Account_flow),
		new(Driver_confirm),
		new(Admin_user),
		)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlurls+")/"+mysqldb+"?charset=utf8&loc=Asia%2FShanghai")
	orm.RunSyncdb("default", false, true)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}


}