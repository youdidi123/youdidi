package models

import (
	"errors"
	"strconv"
	"youdidi/commonLib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func (u *PassengerOrder) TableName() string {
	return "Passenger_order"
}

func (u * PassengerOrder) CreateOrder(uid string, launchTime string,
	startCode int64, endCode int64,
	charge float64, siteNum int,
	travelExplain string , travelCommit string) (int, error) {

	currentTime, _ := strconv.ParseInt(commonLib.GetCurrentTime(), 10, 64)
	launchTimeInt64 := commonLib.TransTimeStoInt64(launchTime)

	if (launchTimeInt64 - currentTime < 10 * 60) {
		return 1, errors.New("请选择至少10分钟后的出发时间")
	}

	o := orm.NewOrm()
	o.Begin()

	var dbUser User
	var userInfo []*User

	num, err := o.QueryTable(dbUser).Filter("Id", uid).ForUpdate().All(&userInfo)

	if (err != nil || num < 1) {
		logs.Error("get user info fail uid=%v", uid)
		o.Rollback()
		return 2, errors.New("网络繁忙，请重试")
	}

	if (userInfo[0].OnRoadType != 0) {
		o.Rollback()
		return 3, errors.New("您有尚未完成的行程，请及时处理")
	}

	cost := charge * float64(siteNum)
	cost = commonLib.FormatMoney(cost)

	if (cost > userInfo[0].Balance) {
		o.Rollback()
		return 4, errors.New("账户余额不足，请充值")
	}

	var po PassengerOrder
	po.User = userInfo[0]
	po.CreateTime = commonLib.GetCurrentTime()
	po.LaunchTime = strconv.FormatInt(launchTimeInt64,10)
	po.Price = charge
	po.PNum = siteNum
	po.ThroughL = travelExplain
	po.Mark = travelCommit
	po.SrcId = &Location{Id:startCode}
	po.DestId = &Location{Id:endCode}
	po.Status = 0
	po.SrcLocationId = startCode % 1000000
	po.DestLocationId = endCode % 1000000

	_, err = o.Insert(&po)

	if (err != nil) {
		logs.Error("insert passenger order fail err=%v", err.Error())
		o.Rollback()
		return 5, errors.New("网络繁忙，请重试")
	}

	balance := commonLib.FormatMoney(userInfo[0].Balance -  cost)

	_, err = o.QueryTable(dbUser).Filter("Id", uid).Update(orm.Params{
		"Balance":balance,
		"OnRoadType":3,
	})
	if (err != nil) {
		logs.Error("update user info fail err=%v", err.Error())
		o.Rollback()
		return 6, errors.New("网络繁忙，请重试")
	}

	var af Account_flow
	af.Balance = balance
	af.Money = cost
	af.User = userInfo[0]
	af.Time = commonLib.GetCurrentTime()
	af.Type = 2

	_, err = o.Insert(&af)
	if (err != nil) {
		logs.Error("insert account flow fail err=%v", err.Error())
		o.Rollback()
		return 7, errors.New("网络繁忙，请重试")
	}

	err = o.Commit()
	if (err != nil) {
		logs.Error("commit fail err=%v", err.Error())
		o.Rollback()
		return 8, errors.New("网络繁忙，请重试")
	}

	moneyStr := strconv.FormatFloat(cost, 'G' , -1,64)
	balanceStr := strconv.FormatFloat(balance, 'G' , -1,64)
	commonLib.SendMsg5(userInfo[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#ff0000","预扣车费",
		"#22c32e", "预扣成功",
		"#173177", moneyStr,
		"#173177", balanceStr)

	return 0, nil

}


func (u *PassengerOrder) GetLastOrdersById (uid string, list *[]*PassengerOrder) (int64, error) {
	return orm.NewOrm().QueryTable(u).RelatedSel().Filter("user_id", uid).Filter("Status", 0).OrderBy("-LaunchTime").Limit(1).All(list)
}

func (u * PassengerOrder) CancleOrder (uid string, oid string) (int, error) {
	o := orm.NewOrm()
	o.Begin()

	var poInfo []*PassengerOrder
	num, err := o.QueryTable(u).Filter("Id", oid).ForUpdate().All(&poInfo)

	if (num < 1 || err != nil) {
		logs.Error("get passenger order fail oid=%v",oid)
		o.Rollback()
		return 1, errors.New("网络繁忙，请重试")
	}

	if (uid != strconv.Itoa(poInfo[0].User.Id)) {
		o.Rollback()
		return 2, errors.New("这个行程不属于您哦")
	}

	if (poInfo[0].Status != 0) {
		o.Rollback()
		return 3, errors.New("请勿重复操作")
	}

	_, err = o.QueryTable(u).Filter("Id", oid).Update(orm.Params{
		"Status":3,
	})
	if (err != nil) {
		logs.Error("update passenger order status 0 to 3 fail err=%v", err.Error())
		o.Rollback()
		return 4, errors.New("网络繁忙，请重试")
	}

	var dbUser User
	var userInfo []*User

	num, err = o.QueryTable(dbUser).Filter("Id", uid).ForUpdate().All(&userInfo)
	if (num < 1 || err != nil) {
		logs.Error("get userinfo fail uid=%v", uid)
		o.Rollback()
		return 5, errors.New("网络繁忙，请重试")
	}

	cost := commonLib.FormatMoney(poInfo[0].Price * float64(poInfo[0].PNum))
	balance := commonLib.FormatMoney(userInfo[0].Balance + cost)
	_, err = o.QueryTable(dbUser).Filter("Id", uid).Update(orm.Params{
		"Balance":balance,
		"OnRoadType":0,
	})
	if (err != nil) {
		logs.Error("update passenger order status 0 to 3 fail err=%v", err.Error())
		o.Rollback()
		return 6, errors.New("网络繁忙，请重试")
	}

	var af Account_flow
	af.Balance = balance
	af.Money = cost
	af.User = userInfo[0]
	af.Time = commonLib.GetCurrentTime()
	af.Type = 4

	_, err = o.Insert(&af)
	if (err != nil) {
		logs.Error("insert account flow fail err=%v", err.Error())
		o.Rollback()
		return 7, errors.New("网络繁忙，请重试")
	}

	err = o.Commit()
	if (err != nil) {
		logs.Error("commit fail err=%v", err.Error())
		o.Rollback()
		return 8, errors.New("网络繁忙，请重试")
	}

	moneyStr := strconv.FormatFloat(cost, 'G' , -1,64)
	balanceStr := strconv.FormatFloat(balance, 'G' , -1,64)
	commonLib.SendMsg5(userInfo[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#22c32e","车费退回",
		"#22c32e", "退回成功",
		"#173177", moneyStr,
		"#173177", balanceStr)

	return 0, nil
}

func (u *PassengerOrder) GetTop20Orders (list *[]*PassengerOrder) (int64, error) {
	return orm.NewOrm().QueryTable(u).RelatedSel().OrderBy("-LaunchTime").Limit(20).All(list)
}

func (u *PassengerOrder) LockOrderBefore (uid string, oid string) (int, int, error) {
	var dbUser User
	var userInfo []*User

	o := orm.NewOrm()

	num, err := o.QueryTable(dbUser).Filter("Id", uid).All(&userInfo)
	if (num < 1 || err != nil) {
		logs.Error("get user info fail uid=%v", uid)
		return 1, 0, errors.New("网络繁忙，请重试")
	}
	if (userInfo[0].IsDriver != 2) {
		return 2, 0, errors.New("仅有认证车主才可直接抢单")
	}
	if (userInfo[0].DisableTime > commonLib.GetCurrentTime()) {
		return 3, 0, errors.New("您在" + commonLib.FormatUnixToStr(userInfo[0].DisableTime) + "前不允许接单")
	}

	var dbOrder Order
	var orderList []*Order

	num, err = o.QueryTable(dbOrder).Filter("user_id", uid).Filter("Status", 1).All(&orderList)
	if (err != nil) {
		logs.Error("get user info fail uid=%v", uid)
		return 4, 0, errors.New("网络繁忙，请重试")
	}

	if (num > 0) {
		return 5, 0, errors.New("您有已在行程中的车主行程尚未结束，请结束此行程后再继续抢单")
	}

	if (userInfo[0].OnRoadType == 1 ||  userInfo[0].OnRoadType == 3) {
		return 6, 0, errors.New("您有已在行程中的乘客行程尚未结束，请完成此行程后再继续抢单")
	}

	num, err = o.QueryTable(dbOrder).Filter("user_id", uid).Filter("Status", 0).All(&orderList)
	if (err != nil) {
		logs.Error("get user info fail uid=%v", uid)
		return 7, 0, errors.New("网络繁忙，请重试")
	}

	if (num < 1) {
		return 0, 1, errors.New("") //车主当前没有行程，需前端提示确认是否自动创建
	}

	var pOrderInfo []*PassengerOrder
	num, err = o.QueryTable(u).Filter("Id", oid).All(&pOrderInfo)
	if (num < 1 || err != nil) {
		logs.Error("get passenger order info fail  oid=%v", oid)
		return 8, 0, errors.New("网络繁忙，请重试")
	}

	if (pOrderInfo[0].PNum > orderList[0].PNum - orderList[0].RequestPnum) {
		return 9, 0, errors.New("您当前行程中的剩余座位不足")
	}

	match := true
	str := ""

	if (pOrderInfo[0].Price != orderList[0].Price) {
		match = false
		str = str + "【定价】"
	}

	if (pOrderInfo[0].LaunchTime != orderList[0].LaunchTime) {
		match = false
		str = str + "【出发时间】"
	}

	if (pOrderInfo[0].SrcLocationId != orderList[0].SrcLocationId || pOrderInfo[0].DestLocationId != orderList[0].DestLocationId) {
		match = false
		str = str + "【路线】"
	}

	if (match) {
		return 0, 2, errors.New("")
	} else {
		return 0, 3, errors.New(str)
	}

}

func (u *PassengerOrder) CreateAndConfirm (uid int, poid string, oid string) (int, string, error) {
	o := orm.NewOrm()
	o.Begin()

	var poInfo []*PassengerOrder
	num, err := o.QueryTable(u).RelatedSel().Filter("Id", poid).ForUpdate().All(&poInfo)
	if (num < 1 || err != nil) {
		logs.Error("get passenger order info fail  oid=%v", poid)
		o.Rollback()
		return 1, "", errors.New("网络繁忙，请重试")
	}

	if (poInfo[0].Status != 0) {
		o.Rollback()
		return 2, "", errors.New("来晚一步哦，改乘客行程已被抢")
	}

	var dbUser User

	_, err = o.QueryTable(dbUser).Filter("Id", uid).Update(orm.Params{
		"OnRoadType":2,
	})
	if (err != nil) {
		logs.Error("update driver onroad type fail err=%v", err.Error())
		o.Rollback()
		return 3, "", errors.New("网络繁忙，请重试")
	}

	_, err = o.QueryTable(dbUser).Filter("Id", poInfo[0].User.Id).Update(orm.Params{
		"OnRoadType":1,
	})
	if (err != nil) {
		logs.Error("update passenger onroad type fail err=%v", err.Error())
		o.Rollback()
		return 4, "", errors.New("网络繁忙，请重试")
	}

	_, err = o.QueryTable(u).Filter("Id", poid).Update(orm.Params{
		"Status":1,
	})
	if (err != nil) {
		logs.Error("update passenger order status fail err=%v", err.Error())
		o.Rollback()
		return 5, "", errors.New("网络繁忙，请重试")
	}


	var orderInfo Order
	orderInfo.Id = oid
	orderInfo.LaunchTime = poInfo[0].LaunchTime
	orderInfo.CreateTime = commonLib.GetCurrentTime()
	orderInfo.Price = poInfo[0].Price
	orderInfo.PNum = 4
	orderInfo.Status = 0
	orderInfo.RequestPnum = poInfo[0].PNum
	orderInfo.ConfirmPnum = poInfo[0].PNum
	orderInfo.SrcId = poInfo[0].SrcId
	orderInfo.DestId = poInfo[0].DestId
	orderInfo.SrcLocationId = poInfo[0].SrcLocationId
	orderInfo.DestLocationId = poInfo[0].DestLocationId
	orderInfo.ThroughL = poInfo[0].ThroughL
	orderInfo.User = &User{Id:uid}

	_, err = o.Insert(&orderInfo)
	if (err != nil) {
		logs.Error("insert order fail err=%v", err.Error())
		o.Rollback()
		return 6, "", errors.New("网络繁忙，请重试")
	}

	var odInfo Order_detail
	odInfo.Status = 1
	odInfo.Price = poInfo[0].Price
	odInfo.SiteNum = poInfo[0].PNum
	odInfo.Order = &orderInfo
	odInfo.Passage = poInfo[0].User
	odInfo.Driver = &User{Id:uid}

	_, err = o.Insert(&odInfo)
	if (err != nil) {
		logs.Error("insert order  detail fail err=%v", err.Error())
		o.Rollback()
		return 7, "", errors.New("网络繁忙，请重试")
	}

	err = o.Commit()
	if (err != nil) {
		logs.Error("commit fail err=%v", err.Error())
		o.Rollback()
		return 8, "", errors.New("网络繁忙，请重试")
	}

	var driverInfo []*User
	_, err = o.QueryTable(dbUser).Filter("Id", uid).All(&driverInfo)

	if (err == nil) {
		commonLib.SendMsg5(poInfo[0].User.OpenId, 5, "http://www.youdidi.vip/Portal/passengerorderdetail/"+strconv.Itoa(odInfo.Id),
			"#22c32e", "车主确认行程通知", "出发前30分钟内取消将会收取违约金，若行程发生变动，请及时操作变更",
			"#173177", poInfo[0].SrcId.Name + " - " + poInfo[0].DestId.Name,
			"#173177", commonLib.FormatUnixToStr(poInfo[0].LaunchTime),
			"#173177", driverInfo[0].CarNum,
			"#173177",driverInfo[0].CarType,
			"#173177", driverInfo[0].Nickname + "(" + driverInfo[0].Phone + ")")
	}

	return 0, "/Portal/driverorderdetail/"+oid, nil
}

func (u *PassengerOrder) LockAndConfirm (uid int, poid string) (int, string, error) {
	o := orm.NewOrm()
	o.Begin()

	var dbUser User
	var poInfo []*PassengerOrder
	num, err := o.QueryTable(u).RelatedSel().Filter("Id", poid).ForUpdate().All(&poInfo)
	if (num < 1 || err != nil) {
		logs.Error("get passenger order info fail  oid=%v", poid)
		o.Rollback()
		return 1, "", errors.New("网络繁忙，请重试")
	}

	if (poInfo[0].Status != 0) {
		o.Rollback()
		return 2, "", errors.New("来晚一步哦，改乘客行程已被抢")
	}

	var dbOrder Order
	var orderInfo []*Order

	num, err = o.QueryTable(dbOrder).RelatedSel().Filter("user_id", uid).Filter("Status", 0).ForUpdate().All(&orderInfo)
	if (num < 1 || err != nil) {
		logs.Error("get dirver order info fail  uid=%v", uid)
		o.Rollback()
		return 3, "", errors.New("网络繁忙，请重试")
	}

	if (poInfo[0].PNum > orderInfo[0].PNum - orderInfo[0].RequestPnum) {
		o.Rollback()
		return 4, "", errors.New("当前行程座位数不足")
	}

	_, err = o.QueryTable(dbUser).Filter("Id", poInfo[0].User.Id).Update(orm.Params{
		"OnRoadType":1,
	})
	if (err != nil) {
		logs.Error("update passenger onroad type fail err=%v", err.Error())
		o.Rollback()
		return 5, "", errors.New("网络繁忙，请重试")
	}

	_, err = o.QueryTable(u).Filter("Id", poid).Update(orm.Params{
		"Status":1,
	})
	if (err != nil) {
		logs.Error("update passenger order status fail err=%v", err.Error())
		o.Rollback()
		return 6, "", errors.New("网络繁忙，请重试")
	}

	_, err = o.QueryTable(dbOrder).Filter("Id", orderInfo[0].Id).Update(orm.Params{
		"RequestPnum":orderInfo[0].RequestPnum + poInfo[0].PNum,
		"ConfirmPnum": orderInfo[0].ConfirmPnum + poInfo[0].PNum,
	})
	if (err != nil) {
		logs.Error("update order info fail err=%v", err.Error())
		o.Rollback()
		return 7, "", errors.New("网络繁忙，请重试")
	}

	var odInfo Order_detail
	odInfo.Status = 1
	odInfo.Price = poInfo[0].Price
	odInfo.SiteNum = poInfo[0].PNum
	odInfo.Order = orderInfo[0]
	odInfo.Passage = poInfo[0].User
	odInfo.Driver = &User{Id:uid}

	_, err = o.Insert(&odInfo)
	if (err != nil) {
		logs.Error("insert order  detail fail err=%v", err.Error())
		o.Rollback()
		return 8, "", errors.New("网络繁忙，请重试")
	}


	err = o.Commit()
	if (err != nil) {
		logs.Error("commit fail err=%v", err.Error())
		o.Rollback()
		return 9, "", errors.New("网络繁忙，请重试")
	}

	var driverInfo []*User
	_, err = o.QueryTable(dbUser).Filter("Id", uid).All(&driverInfo)

	if (err == nil) {
		commonLib.SendMsg5(poInfo[0].User.OpenId, 5, "http://www.youdidi.vip/Portal/passengerorderdetail/"+strconv.Itoa(odInfo.Id),
			"#22c32e", "车主确认行程通知", "出发前30分钟内取消将会收取违约金，若行程发生变动，请及时操作变更",
			"#173177", poInfo[0].SrcId.Name + " - " + poInfo[0].DestId.Name,
			"#173177", commonLib.FormatUnixToStr(poInfo[0].LaunchTime),
			"#173177", driverInfo[0].CarNum,
			"#173177",driverInfo[0].CarType,
			"#173177", driverInfo[0].Nickname + "(" + driverInfo[0].Phone + ")")
	}

	return 0, "/Portal/driverorderdetail/"+orderInfo[0].Id, nil
}

func (u *PassengerOrder) GetReadyOrders (list *[]*PassengerOrder , startCode int , endCode int , tmStart int64 , tmEnd int64 , startLocation int , endLocation int) int64 {
	var jqInfo []*PassengerOrder
	var mhInfo []*PassengerOrder
	var resultInfo []*PassengerOrder
	var enableInfo []*PassengerOrder
	var disableInfo []*PassengerOrder
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