package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"youdidi/commonLib"
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
		OrderBy("Status","-LaunchTime").
		All(list)
	return num
}

func (u *Order) GetCurrentDriverOrderFromUserId (id string, list *[]*Order) int64{
	num , _ := orm.NewOrm().QueryTable(u).
		Filter("user_id", id).
		Filter("status__lt",2). //车主当前的行程只会有比2小的
		OrderBy("-LaunchTime").
		Limit(1).
		All(list)
	return num
}

func (u *Order) GetReadyOrders (list *[]*Order , startCode int , endCode int , tmStart int64 , tmEnd int64 , startLocation int , endLocation int) int64 {
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
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().Filter("id", oid).All(list)
	return num
}

func (u *Order) DoRequire (od *Order_detail, pid string, siteNum int , mark string) bool{
	var userInfo User
	var userInfos []*User
	o := orm.NewOrm()
	o.Begin()

	var orderInfo []*Order
	numorder, errorder := o.QueryTable(u).Filter("Id", u.Id).RelatedSel().ForUpdate().All(&orderInfo)

	if (numorder < 1 || errorder != nil) {
		logs.Error("get order info fail oid=%v" , u.Id)
		o.Rollback()
		return false
	}

	numuser, erruser := o.QueryTable(userInfo).Filter("id", pid).ForUpdate().All(&userInfos)
	if (numuser < 1 || erruser != nil) {
		logs.Error("get user info fail pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}
	balance, _ := strconv.ParseFloat(fmt.Sprintf("%.2f",userInfos[0].Balance - orderInfo[0].Price * float64(siteNum)), 64)
	if (balance < 0) {
		logs.Error("user balance is not enough uid=%v balance=%v" , u.Id , balance)
		o.Rollback()
		return false
	}

	//订单中讲requestPnum+座位数
	if ((orderInfo[0].PNum - orderInfo[0].RequestPnum) < od.SiteNum) {
		logs.Error("siteNum is not enough")
		o.Rollback()
		return false
	}
	num := u.RequestPnum + od.SiteNum
	_ , err2 := o.QueryTable(u).Filter("id", od.Order.Id).Update(orm.Params{"RequestPnum":num})
	if (err2 != nil) {
		logs.Error("requestPnum add 1 fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}


	//想order detail里插入一条记录
	_ , err := o.Insert(od)
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


	var chatInfo Chat
	chatInfo.Order = &Order{Id:u.Id}
	if (mark == "") {
		mark = "希望能和您同行，辛苦通过拼车申请"
	}
	chatInfo.Content = mark
	passengerId, _ := strconv.Atoi(pid)
	chatInfo.Passenger = &User{Id:passengerId}
	chatInfo.Type = 0
	chatInfo.TimeStamp = strconv.FormatInt(time.Now().Unix(),10)

	_, err3 := o.Insert(&chatInfo)
	if (err3 != nil) {
		logs.Error("insert chat info fail but go on")
	}

	var accountFlow Account_flow
	accountFlow.Type = 2
	accountFlow.User = &User{Id:passengerId}
	accountFlow.Oid = u.Id
	accountFlow.Money,_ = strconv.ParseFloat(fmt.Sprintf("%.2f",u.Price * float64(siteNum)), 64)
	accountFlow.Time = strconv.FormatInt(time.Now().Unix(),10)
	accountFlow.Balance = balance

	_, err4 := o.Insert(&accountFlow)
	if (err4 != nil) {
		logs.Error("record account flow fail oid=%v pid=%v", u.Id, pid)
		o.Rollback()
		return false
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail oid=%v pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}
	moneyStr := strconv.FormatFloat(accountFlow.Money, 'G' , -1,64)
	balanceStr := strconv.FormatFloat(accountFlow.Balance, 'G' , -1,64)
	commonLib.SendMsg5(userInfos[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#ff0000","预扣车费",
		"#22c32e", "预扣成功",
		"#173177", moneyStr,
		"#173177", balanceStr)
	return true
}

func (u *Order) DriverGetStart (oid string) bool{
	o := orm.NewOrm()
	var od Order_detail
	o.Begin()
	var odInfo []*Order_detail

	_, err4 := o.QueryTable(od).Filter("order_id", oid).Filter("Status", 1).RelatedSel().All(&odInfo)

	_, err1 := o.QueryTable(u).Filter("id", oid).Update(orm.Params{"Status":1})
	if (err1 != nil) {
		logs.Error("set order status to 1 fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	_, err2 := o.QueryTable(od).Filter("order_id", oid).Filter("Status", 1).Update(orm.Params{"Status":2})
	if (err2 != nil) {
		logs.Error("set order detail status to 2 fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	err3 := o.Commit()

	if (err3 != nil) {
		logs.Error("commit fail oid=%v" , u.Id)
		o.Rollback()
		return false
	}
	tm := time.Now()

	if (err4 == nil) {
		for _, v := range odInfo {
			commonLib.SendMsg4(v.Passage.OpenId, 6, "http://www.youdidi.vip/Portal/passengerorderdetail/" + strconv.Itoa(v.Id),
				"#22c32e", "车主已到达出发地点", "请尽快到达出发地点，以免耽误行程",
				"#173177", "同行拼车",
				"#173177", v.Order.Id,
				"#22c32e", "车主已到达出发地",
				"#173177", tm.Format("2006-01-02 15:04"))
		}
	}

	return true
}

func (u *Order) DriverGetEnd (oid string, uid string) bool{
	o := orm.NewOrm()
	var driver User
	var driverInfo []*User
	var od Order_detail
	var odInfo []*Order_detail
	o.Begin()

	_, err6 := o.QueryTable(od).Filter("order_id", oid).Filter("Status", 2).RelatedSel().All(&odInfo)

	_, err1 := o.QueryTable(u).Filter("id", oid).Update(orm.Params{"Status":2})
	if (err1 != nil) {
		logs.Error("set order status to 2 fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	_, err2 := o.QueryTable(driver).Filter("id", uid).Update(orm.Params{"OnRoadType":0})
	if (err2 != nil) {
		logs.Error("set user onroadtype to 0 fail uid=%v" , uid)
		o.Rollback()
		return false
	}

	_, err5 := o.QueryTable(od).Filter("order_id", oid).Filter("Status", 2).Update(orm.Params{"Status":3})
	if (err5 != nil) {
		logs.Error("set order detail status to 2 fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	num, err3 :=  o.QueryTable(driver).Filter("id", uid).All(&driverInfo)

	if (num < 1 || err3 != nil) {
		logs.Error("get user info fail uid=%v" , uid)
		o.Rollback()
		return false
	}

	orderNumWV := driverInfo[0].OrderNumWV
	if (driverInfo[0].IsDriver < 2) {
		orderNumWV += 1
	}

	_, err4 := o.QueryTable(driver).Filter("id", uid).Update(orm.Params{
		"OrderNumAsD":driverInfo[0].OrderNumAsD + 1,
		"OrderNumWV":orderNumWV})
	if (err4 != nil) {
		logs.Error("add user orderNumAsD fail uid=%v" , uid)
		o.Rollback()
		return false
	}

	tm := time.Now()

	if (err6 == nil) {
		for _, v := range odInfo {
			commonLib.SendMsg4(v.Passage.OpenId, 6, "http://www.youdidi.vip/Portal/passengerorderdetail/" + strconv.Itoa(v.Id),
				"#22c32e", "车主已确认到达目的地", "乘客请尽快确认到达目的地以便车主收款，若以确认请忽略次消息",
				"#173177", "同行拼车",
				"#173177", v.Order.Id,
				"#22c32e", "车主已确认到达目的地",
				"#173177", tm.Format("2006-01-02 15:04"))
		}
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail oid=%v" , u.Id)
		o.Rollback()
		return false
	}



	return true
}

func (u *Order) DriverCancle (oid string, confirmNum int, driverId string) bool{
	o := orm.NewOrm()
	var driver User
	var driverInfo []*User
	var od Order_detail
	var odList []*Order_detail

	o.Begin()

	num1, err1 := o.QueryTable(u).Filter("id", oid).Filter("Status__lt", 2).Update(orm.Params{"Status":3})
	if (num1< 1 || err1 != nil) {
		logs.Error("set order status to 3 fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	num, err2 :=  o.QueryTable(driver).Filter("id", driverId).All(&driverInfo)

	if (num < 1 || err2 != nil) {
		logs.Error("get user info fail uid=%v" , driverId)
		o.Rollback()
		return false
	}

	orderNumAsD := driverInfo[0].OrderNumAsD + 1
	cancleOasD := driverInfo[0].CancleOasD
	orderNumWV := driverInfo[0].OrderNumWV
	disableTime := ""
	if (confirmNum > 0) {
		delayDay, _ := beego.AppConfig.Int64("driverCancleDelay")
		currentTime := time.Now().Unix()
		currentTime = currentTime + delayDay * 24 * 60 * 60
		disableTime = strconv.FormatInt(currentTime,10)
		cancleOasD += 1 //仅当有确定乘客的情况下取消才计入取消次数里
		if (driverInfo[0].IsDriver < 2) {
			orderNumWV += 1
		}
	}

	_, err3 := o.QueryTable(driver).Filter("id", driverId).Update(orm.Params{
		"OrderNumWV":orderNumWV,
		"OrderNumAsD":orderNumAsD,
		"CancleOasD":cancleOasD,
		"DisableTime":disableTime,
		"OnRoadType":0,
	})

	if (err3 != nil) {
		logs.Error("update driver info fail uid=%v" , driverId)
		o.Rollback()
		return false
	}

	_, err4 := o.QueryTable(od).RelatedSel().Filter("order_id", oid).Filter("Status__lt", 4).ForUpdate().All(&odList)

	if (err4 != nil) {
		logs.Error("get order detail fail oid=%v" , oid)
		o.Rollback()
		return false
	}

	for _, v := range odList {
		_, err5 := o.QueryTable(od).Filter("id", v.Id).Update(orm.Params{"Status":7})
		if (err5 != nil) {
			logs.Error("set order detail status to 7 fail odid=%v" , v.Id)
			o.Rollback()
			return false
		}
		var passenger User
		var passengerInfo []*User
		num, err6 := o.QueryTable(passenger).Filter("id", v.Passage.Id).ForUpdate().All(&passengerInfo)
		if (num < 1 || err6 != nil) {
			logs.Error("get passenger info fail id=%v" , v.Passage.Id)
			o.Rollback()
			return false
		}
		balance, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v.Order.Price * float64(v.SiteNum) + passengerInfo[0].Balance), 64)

		_, err7 := o.QueryTable(passenger).Filter("id", v.Passage.Id).Update(orm.Params{"Balance":balance,"OnRoadType":0})
		if (err7 != nil) {
			logs.Error("return money to passager fail uid=%v" , v.Id)
			o.Rollback()
			return false
		}
		var accountFlow Account_flow
		accountFlow.Type = 4
		accountFlow.User = &User{Id:passengerInfo[0].Id}
		accountFlow.Oid = oid
		accountFlow.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", v.Order.Price * float64(v.SiteNum)), 64)
		accountFlow.Time = strconv.FormatInt(time.Now().Unix(),10)
		accountFlow.Balance = balance

		_, err8 := o.Insert(&accountFlow)
		if (err8 != nil) {
			logs.Error("record account flow fail oid=%v pid=%v", u.Id, passengerInfo[0].Id)
			o.Rollback()
			return false
		}
		commonLib.SendMsg5(passengerInfo[0].OpenId, 3, "http://www.youdidi.vip/Portal/passengerorderdetail/"+strconv.Itoa(v.Id),
			"#ff0000", "抱歉，车主已操作取消行程", "为避免对您的影响，请尽快查询其他车主发起的行程",
			"#173177", v.Passage.Nickname,
			"#173177", v.Order.SrcId.Level1 + "-" + v.Order.SrcId.Level2 + "-" + v.Order.SrcId.Name,
			"#173177", v.Order.DestId.Level1 + "-" + v.Order.DestId.Level2 + "-" + v.Order.DestId.Name,
			"#173177", "抱歉，行程临时有变",
			"#173177", time.Now().Format("2006-01-02 15:04"))

		moneyStr := strconv.FormatFloat(accountFlow.Money, 'G' , -1,64)
		balanceStr := strconv.FormatFloat(accountFlow.Balance, 'G' , -1,64)
		commonLib.SendMsg5(passengerInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#22c32e","车费退回",
			"#22c32e", "退回成功",
			"#173177", moneyStr,
			"#173177", balanceStr)

	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail oid=%v" , u.Id)
		o.Rollback()
		return false
	}

	return true
}