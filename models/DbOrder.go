package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func (u *Order) TableName() string {
	return "Order"
}

func (u *Order) Insert() (int64, error) {
	return orm.NewOrm().Insert(u)
}

func (u *Order) GetOrderInfoFromUserId (id string, list *[]*Order) int64{
	num , _ := orm.NewOrm().QueryTable(u).
		RelatedSel("SrcId","DestId").
		Filter("user_id", id).
		OrderBy("-LaunchTime").
		All(list)
	return num
}

func (u *Order) GetCurrentDriverOrderFromUserId (id string, list *[]*Order) int64{
	num , _ := orm.NewOrm().QueryTable(u).
		Filter("user_id", id).
		Filter("status__lt",4). //车主当前的行程只会有比4小的
		OrderBy("-LaunchTime").
		Limit(1).
		All(list)
	return num
}

func (u *Order) GetReadyOrders (list *[]*Order , startCode int64 , endCode int64 , tmStart int64 , tmEnd int64 , startLocation int64 , endLocation int64) int64 {
	var jqInfo []*Order
	var mhInfo []*Order
	var resultInfo []*Order
	var enableInfo []*Order
	var disableInfo []*Order
	o :=orm.NewOrm()
	num1 , err := o.QueryTable(u).
		RelatedSel("SrcId", "DestId" , "User").
		//Filter("status__lt", 1). //查询展示所有的行程
		Filter("SrcId" , startCode).
		Filter("DestId" , endCode).
		Filter("LaunchTime__gt" , tmStart).
		Filter("LaunchTime__lt" , tmEnd).
		OrderBy("Status" , "LaunchTime"). //先通过status排序，把能够预约的排在前面，其他的单子都展示，提现单子多
		All(&jqInfo)
	logs.Debug("query jingque order num=%v" , num1)
	if (err != nil) {
		logs.Error("query jingque order search fail err=%v" , err.Error())
	} else {
		for _ , v := range jqInfo {
			if (v.Status < 1) {
				enableInfo = append(enableInfo, v)
			} else {
				disableInfo = append(disableInfo , v)
			}
		}
	}
	cond := orm.NewCondition()
	cond1 := cond.And("SrcLocationId__gt" , startLocation / 100 * 100)
	cond2 := cond.And("SrcLocationId__lt" , (startLocation / 100 + 1)  * 100)
	cond3 := cond.AndCond(cond1).AndCond(cond2)

	cond4 := cond.And("DestLocationId__gt" , endLocation / 100 * 100)
	cond5 := cond.And("DestLocationId__lt" , (endLocation / 100 + 1)  * 100)
	cond6 := cond.AndCond(cond4).AndCond(cond5)

	cond7 := cond.AndCond(cond3).AndCond(cond6)

	cond8 := cond.And("LaunchTime__gt", tmStart)
	cond9 := cond.And("LaunchTime__lt" , tmEnd)
	cond10 := cond.AndCond(cond8).AndCond(cond9)

	cond11 := cond.AndCond(cond7).AndCond(cond10)

	cond12 := cond.And("SrcId" , startCode)
	cond13 := cond.And("DestId" , endCode)
	cond14 := cond.AndCond(cond12).AndCond(cond13)

	cond15 := cond.AndCond(cond11).AndNotCond(cond14)


	num2 , err1 := orm.NewOrm().QueryTable(u).
		RelatedSel("SrcId", "DestId" , "User").
		SetCond(cond15).
		OrderBy("Status" , "LaunchTime").
		All(&mhInfo)
	logs.Debug("query mohu order num=%v" , num2)
	if (err1 != nil) {
		logs.Error("query mohu order search fail err=%v" , err1.Error())
	} else {
		for _ , v := range mhInfo {
			if (v.Status < 1) {
				enableInfo = append(enableInfo, v)
			} else {
				disableInfo = append(disableInfo , v)
			}
		}
	}

	for _ , v := range enableInfo {
		resultInfo = append(resultInfo , v)
	}

	for _ , v := range disableInfo {
		resultInfo = append(resultInfo , v)
	}

	*list = resultInfo

	return int64(len(resultInfo))
}

func (u *Order) GetOrderFromId (oid string , list *[]*Order) int64 {
	num , _ := orm.NewOrm().QueryTable(u).Filter("id", oid).All(list)
	return num
}

func (u *Order) DoRequire (od Order_detail, pid string, siteNum int , mark string, balance float64) bool{
	var userInfo User
	o := orm.NewOrm()
	o.Begin()

	//想order detail里插入一条记录
	_ , err := o.Insert(&od)
	if (err != nil) {
		logs.Error("insert order detail fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}

	//将乘客的状态改为行程中
	_ , err1 := o.QueryTable(userInfo).Filter("id", pid).Update(orm.Params{
		"onRoadType": 1,"balance":balance,
	})
	if (err1 != nil) {
		logs.Error("set user stauts to 1 fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}

	//订单中讲requestPnum+座位数
	num := u.RequestPnum + od.SiteNum
	_ , err2 := o.QueryTable(u).Filter("id", od.Order.Id).Update(orm.Params{"RequestPnum":num})
	if (err2 != nil) {
		logs.Error("requestPnum add 1 fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}
	err3 := o.Commit()

	if (err3 != nil) {
		logs.Error("commit fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}
	return true
}