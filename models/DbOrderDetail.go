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

func (u *Order_detail) TableName() string {
	return "Order_detail"
}

func (u *Order_detail) TableIndex() [][]string {
	return [][]string{
		[]string{"Driver", "Passage"},
	}
}

func (u *Order_detail) TableUnique() [][]string {
	return [][]string{
		[]string{"Order", "Driver"},{"Order", "Passage"},
	}
}



func (u *Order_detail) GetOrderDetailFromPassengerId (id string , list *[]*Order_detail) int64 {
	num , err := orm.NewOrm().QueryTable(u).RelatedSel().Filter("passage_id", id).OrderBy("Status").All(list)

	if (err != nil) {
		logs.Error("query order detail from passenger fail pid=%v" , id)
	}
	return num
}

func (u *Order_detail) GetCurrentPassengerOrderFromUserId (id string, list *[]*Order_detail) int64{
	num , _ := orm.NewOrm().QueryTable(u).
		Filter("passage_id", id).
		Filter("status__lt",4). //乘客当前的行程只会有比4小的
		Limit(1).
		All(list)
	return num
}

func (u *Order_detail) GetOrderedOrderFromPassengerId (oid string , uid string , list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("order_id" , oid).
		Filter("passage_id" , uid).
		Filter("status__lt" , 5).
		All(list)
	return num
}

func (u *Order_detail) GetOrderDetailFromOrderId (id string , list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("order_id" , id).OrderBy("status").
		All(list)
	return num
}

func (u *Order_detail) GetOrderDetailFromId (id string ,list *[]*Order_detail) int64 {
	num , _ := orm.NewOrm().QueryTable(u).RelatedSel().
		Filter("id" , id).
		All(list)
	return num
}

func (u *Order_detail) AgreeRequest(odid string , oid string , confirmPnum int , siteNum int) bool {
	var orderInfo Order
	o := orm.NewOrm()
	o.Begin()

	num2 , err2 := o.QueryTable(u).Filter("id" , odid).Filter("Status", 0).Update(orm.Params{
		"Status": 1,
	})
	if (num2 <1 || err2 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	_ , err1 := o.QueryTable(orderInfo).Filter("id", oid).Update(orm.Params{
		"ConfirmPnum": confirmPnum+siteNum,
	})
	if (err1 != nil) {
		logs.Error("add confirmPnum fail oid=%v confirmPnum=%v" , oid , confirmPnum+1)
		o.Rollback()
		return false
	}



	err3 := o.Commit()

	if (err3 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	return true
}

func (u *Order_detail) RefuseRequest(odid string , oid string , pid string , requestNum int ,
	refuseNum int , siteNum int , price float64) bool {
	var orderInfo Order
	var userInfo User
	var userInfos []*User
	o := orm.NewOrm()
	o.Begin()

	num2 , err2 := o.QueryTable(u).Filter("id" , odid).Filter("Status__lt", 4).ForUpdate().Update(orm.Params{
		"Status": 5,
	})
	if (num2 < 1 || err2 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}

	_ , err1 := o.QueryTable(orderInfo).Filter("id", oid).Update(orm.Params{
		"RefusePnum": refuseNum+siteNum,"RequestPnum":requestNum-siteNum,
	})
	if (err1 != nil) {
		logs.Error("add confirmPnum fail oid=%v RefusePnum=%v RequestPnum=%v" , oid , refuseNum+1 , requestNum-1)
		o.Rollback()
		return false
	}

	numuser, erruser := o.QueryTable(userInfo).Filter("id", pid).ForUpdate().All(&userInfos)
	if (numuser < 1 || erruser != nil) {
		logs.Error("get user info fail pid=%v" , u.Id , pid)
		o.Rollback()
		return false
	}

	balance, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", userInfos[0].Balance + price * float64(siteNum)), 64)
	logs.Debug("user balance=%v price=%v siteNum=%v",userInfos[0].Balance, price, siteNum)
	_ , err3 := o.QueryTable(userInfo).Filter("id", pid).ForUpdate().Update(orm.Params{
		"OnRoadType": 0,"Balance":balance,
	})
	if (err3 != nil) {
		logs.Error("update passenger info fail id=%v" , pid)
		o.Rollback()
		return false
	}

	var accountFlow Account_flow
	accountFlow.Type = 4
	accountFlow.User = &User{Id:userInfos[0].Id}
	accountFlow.Oid = oid
	accountFlow.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", price * float64(siteNum)), 64)
	accountFlow.Time = strconv.FormatInt(time.Now().Unix(),10)
	accountFlow.Balance = balance

	_, err8 := o.Insert(&accountFlow)
	if (err8 != nil) {
		logs.Error("record account flow fail oid=%v pid=%v", u.Id, userInfos[0].Id)
		o.Rollback()
		return false
	}

	err4 := o.Commit()

	if (err4 != nil) {
		logs.Error("set orderdetail status fail oid=%v odid=%v" , oid , odid)
		o.Rollback()
		return false
	}



	moneyStr := strconv.FormatFloat(accountFlow.Money, 'G' , -1,64)
	balanceStr := strconv.FormatFloat(accountFlow.Balance, 'G' , -1,64)
	commonLib.SendMsg5(userInfos[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#22c32e","车费退回",
		"#22c32e", "退回成功",
		"#173177", moneyStr,
		"#173177", balanceStr)

	return true
}

func (u *Order_detail) PassengerConfirm(odid string) bool {
	o := orm.NewOrm()
	o.Begin()

	var odInfo []*Order_detail
	num, err1 := o.QueryTable(u).Filter("Id", odid).Filter("Status__lt", 4).ForUpdate().RelatedSel().All(&odInfo)
	if (err1 != nil || num < 1) {
		logs.Error("get order detail fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	driverId := odInfo[0].Driver.Id
	passengerId := odInfo[0].Passage.Id

	var dbUser User
	var driverInfo []*User
	var passengerInfo []*User

	num2, err2 := o.QueryTable(dbUser).Filter("Id", driverId).ForUpdate().All(&driverInfo)
	num3, err3 := o.QueryTable(dbUser).Filter("Id", passengerId).ForUpdate().All(&passengerInfo)

	if (num2 < 1 || num3 < 1 || err2 != nil || err3 != nil) {
		logs.Error("get userinfo fail fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	_, err4 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
		"Status": 4,
	})
	if (err4 != nil) {
		logs.Error("set order detail status to 4 fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	//修改乘客数据
	orderNumAsP := passengerInfo[0].OrderNumAsP + 1
	_, err5 := o.QueryTable(dbUser).Filter("Id", passengerId).Update(orm.Params{
		"OrderNumAsP": orderNumAsP,
		"OnRoadType": 0,
	})
	if (err5 != nil) {
		logs.Error("update passanger info fail pid=%v" , passengerId)
		o.Rollback()
		return false
	}

	//修改车主数据
	infoCostRatio, _ := beego.AppConfig.Float("infoCostRatio")
	allCost := float64(odInfo[0].SiteNum) * odInfo[0].Order.Price
	infoCost := allCost * infoCostRatio
	driverCost := allCost - infoCost


	balance := driverInfo[0].Balance + driverCost
	balance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", balance), 64)
	_, err6 := o.QueryTable(dbUser).Filter("Id", driverId).Update(orm.Params{
		"Balance": balance,
	})
	if (err6 != nil) {
		logs.Error("update driver info fail pid=%v" , driverId)
		o.Rollback()
		return false
	}

	//添加account_flow数据
	var afDriver Account_flow
	var afPassenger Account_flow

	afDriver.Balance = balance
	afDriver.Oid = odInfo[0].Order.Id
	afDriver.User = driverInfo[0]
	afDriver.Type = 5
	afDriver.Time = strconv.FormatInt(time.Now().Unix(),10)
	afDriver.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", driverCost), 64)

	afPassenger.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(odInfo[0].SiteNum) * odInfo[0].Order.Price), 64)
	afPassenger.Time = strconv.FormatInt(time.Now().Unix(),10)
	afPassenger.Oid = odInfo[0].Order.Id
	afPassenger.Type = 3
	afPassenger.User = passengerInfo[0]
	afPassenger.Balance = passengerInfo[0].Balance

	_, err7 := o.Insert(&afDriver)
	_, err8 := o.Insert(&afPassenger)

	if (err7 != nil || err8 != nil) {
		logs.Error("insert account flow fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	var afDriver2 Account_flow
	afDriver2.Type = 9
	afDriver2.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", infoCost), 64)
	afDriver2.User = driverInfo[0]
	afDriver2.Time = strconv.FormatInt(time.Now().Unix(),10)
	afDriver2.Balance = balance

	_, err9 := o.Insert(&afDriver2)
	if (err9 != nil) {
		logs.Error("insert account flow fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail odid=%v" , odid)
		o.Rollback()
		return false
	}
	zhifukuan := strconv.FormatFloat(afPassenger.Money, 'G' , -1,64)
	shouxufei := strconv.FormatFloat(infoCost, 'G' , -1,64)
	pbstr := strconv.FormatFloat(afPassenger.Balance, 'G' , -1,64)
	dbstr := strconv.FormatFloat(balance, 'G' , -1,64)
	dincome := strconv.FormatFloat(driverCost, 'G' , -1,64)

	commonLib.SendMsg5(passengerInfo[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#ff0000","预充车费确认支付",
		"#22c32e", "支付成功",
		"#173177", zhifukuan,
		"#173177", pbstr)
	commonLib.SendMsg5(driverInfo[0].OpenId,
		4, "", "#173177", "", "",
		"#173177", "",
		"#22c32e",passengerInfo[0].Nickname+":支付车费",
		"#22c32e", "支付成功",
		"#173177", dincome + "(信息费：" + shouxufei + ")",
		"#173177", dbstr)

	return true
}

func (u *Order_detail) PassengerCancle(odid string) bool {
	passangerCancleTime, _ := beego.AppConfig.Int64("passangerCancleTime")
	passangerCancleRatio, _ := beego.AppConfig.Float("passangerCancleRatio")
	infoCostRatio, _ := beego.AppConfig.Float("infoCostRatio")

	var afDriver Account_flow
	var afPassenger Account_flow
	var afDriver2 Account_flow
	var afPassenger2 Account_flow

	var dbOrder Order
	var orderInfo []*Order

	isAllBack := false
	is80Back := false
	is0Back := false
	driverIncome := ""
	passengerIncome := ""
	weiyuejin := ""
	shouxufei := ""
	passengerB := ""
	driverB := ""

	o := orm.NewOrm()
	o.Begin()

	var odInfo []*Order_detail
	num, err1 := o.QueryTable(u).Filter("Id", odid).Filter("Status__lt", 4).ForUpdate().RelatedSel().All(&odInfo)
	if (err1 != nil || num < 1) {
		logs.Error("get order detail fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	var dbUser User
	var passengerInfo []*User
	var driverInfo []*User
	num2, err2 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Passage.Id).ForUpdate().All(&passengerInfo)
	num3, err3 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Driver.Id).ForUpdate().All(&driverInfo)
	if (err2 != nil || num2 < 1 || err3 != nil || num3 < 1) {
		logs.Error("get passenger or driver info fail pid=%v did=%v" , odInfo[0].Passage.Id, odInfo[0].Driver.Id)
		o.Rollback()
		return false
	}

	_, err10 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
		"Status": 6,
	})
	if (err10 != nil) {
		logs.Error("update od info fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	num11, err11 := o.QueryTable(dbOrder).Filter("Id", odInfo[0].Order.Id).ForUpdate().All(&orderInfo)
	if (num11 < 1 || err11 != nil) {
		logs.Error("get order info fail id=%v" , odInfo[0].Order.Id)
		o.Rollback()
		return false
	}

	confirmPnum := orderInfo[0].ConfirmPnum
	if (odInfo[0].Status != 0) {
		confirmPnum = confirmPnum - odInfo[0].SiteNum
	}

	_, err12 := o.QueryTable(dbOrder).Filter("Id", odInfo[0].Order.Id).Update(orm.Params{
		"RequestPnum": orderInfo[0].RequestPnum - odInfo[0].SiteNum,
		"CanclePnum" : orderInfo[0].CanclePnum + odInfo[0].SiteNum,
		"ConfirmPnum": confirmPnum,
	})
	if (err12 != nil) {
		logs.Error("update order detail fail id=%v" , odInfo[0].Order.Id)
		o.Rollback()
		return false
	}

	afDriver.User = driverInfo[0]
	afDriver.Oid = odInfo[0].Order.Id
	afPassenger.User = passengerInfo[0]
	afPassenger.Oid = odInfo[0].Order.Id
	afDriver2.User = driverInfo[0]
	afDriver2.Oid = odInfo[0].Order.Id
	afPassenger2.User = passengerInfo[0]
	afPassenger2.Oid = odInfo[0].Order.Id

	allCost := float64(odInfo[0].SiteNum) * odInfo[0].Order.Price
	allCost, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", allCost), 64)

	passengerBalance := passengerInfo[0].Balance
	driverBalance := driverInfo[0].Balance
	orderNumAsP := passengerInfo[0].OrderNumAsP + 1

	currentTime := time.Now().Unix()



	if (odInfo[0].Status == 0) {
		passengerBalance += allCost
		passengerBalance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", passengerBalance), 64)
		_, err4 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Passage.Id).Update(orm.Params{
			"Balance": passengerBalance,
			"OnRoadType": 0,
			"OrderNumAsP":orderNumAsP,
		})
		if (err4 != nil) {
			logs.Error("update passenger info fail pid=%v" , odInfo[0].Passage.Id)
			o.Rollback()
			return false
		}
		afPassenger.Money = allCost
		afPassenger.Balance = passengerBalance
		afPassenger.Type = 4
		afPassenger.Time = strconv.FormatInt(currentTime, 10)
		_, err5 := o.Insert(&afPassenger)
		if (err5 != nil) {
			logs.Error("set account flow fail")
			o.Rollback()
			return false
		}
		isAllBack = true
		passengerIncome = strconv.FormatFloat(afPassenger.Money, 'G' , -1,64)
		passengerB = strconv.FormatFloat(afPassenger.Balance, 'G' , -1,64)
	} else {
		launchTime,_ := strconv.ParseInt(odInfo[0].Order.LaunchTime, 10, 64)
		if (launchTime - currentTime > passangerCancleTime) {
			// 取消时间在30min以外，直接退款
			passengerBalance += allCost
			passengerBalance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", passengerBalance), 64)
			_, err4 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Passage.Id).Update(orm.Params{
				"Balance": passengerBalance,
				"OnRoadType": 0,
				"OrderNumAsP":orderNumAsP,
			})
			if (err4 != nil) {
				logs.Error("update passenger info fail pid=%v" , odInfo[0].Passage.Id)
				o.Rollback()
				return false
			}
			afPassenger.Money = allCost
			afPassenger.Balance = passengerBalance
			afPassenger.Type = 4
			afPassenger.Time = strconv.FormatInt(currentTime, 10)
			_, err5 := o.Insert(&afPassenger)
			if (err5 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			isAllBack = true
			passengerIncome = strconv.FormatFloat(afPassenger.Money, 'G' , -1,64)
			passengerB = strconv.FormatFloat(afPassenger.Balance, 'G' , -1,64)
		} else if (currentTime > launchTime && odInfo[0].Status > 1) {
			//在出发时间之后，且车主以确认到达，扣除100%
			infoCost := allCost * infoCostRatio
			infoCost, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", infoCost), 64)

			driverBalance = driverBalance + (allCost - infoCost)
			//给司机加钱
			driverBalance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", driverBalance), 64)
			_, err4 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Driver.Id).Update(orm.Params{
				"Balance": driverBalance,
			})
			if (err4 != nil) {
				logs.Error("update dirver info fail pid=%v" , odInfo[0].Driver.Id)
				o.Rollback()
				return false
			}
			afDriver.Money = allCost - infoCost
			afDriver.Balance = driverBalance
			afDriver.Type = 7
			afDriver.Time = strconv.FormatInt(currentTime, 10)
			_, err5 := o.Insert(&afDriver)
			if (err5 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			afDriver2.Balance = driverBalance
			afDriver2.Time = strconv.FormatInt(currentTime, 10)
			afDriver2.Money = infoCost
			afDriver2.Type = 9
			_, err8 := o.Insert(&afDriver2)
			if (err8 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			//修改乘客信息
			_, err6 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Passage.Id).Update(orm.Params{
				"OnRoadType": 0,
				"OrderNumAsP":orderNumAsP,
				"CancleOasP": passengerInfo[0].CancleOasP + 1,
			})
			if (err6 != nil) {
				logs.Error("update passanger info fail pid=%v" , odInfo[0].Passage.Id)
				o.Rollback()
				return false
			}
			afPassenger.Money = allCost
			afPassenger.Type = 6
			afPassenger.Balance = passengerBalance
			afPassenger.Time = strconv.FormatInt(currentTime, 10)
			_, err7 := o.Insert(&afPassenger)
			if (err7 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			is0Back = true
			weiyuejin = strconv.FormatFloat(afPassenger.Money, 'G' , -1,64)
			shouxufei = strconv.FormatFloat(infoCost, 'G' , -1,64)
			driverIncome = strconv.FormatFloat(afDriver.Money, 'G' , -1,64)
			passengerIncome = "0"
			driverB = strconv.FormatFloat(driverBalance, 'G' , -1,64)
			passengerB = strconv.FormatFloat(passengerBalance, 'G' , -1,64)
		} else {
			//其余情况都扣20%给车主
			kCost := allCost * passangerCancleRatio
			infoCost := kCost * infoCostRatio
			kCost, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", kCost), 64)
			infoCost, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", infoCost), 64)

			passengerBalance = passengerBalance + (allCost - kCost)
			driverBalance = driverBalance + (kCost - infoCost)

			_, err4 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Driver.Id).Update(orm.Params{
				"Balance": driverBalance,
			})
			if (err4 != nil) {
				logs.Error("update dirver info fail pid=%v" , odInfo[0].Driver.Id)
				o.Rollback()
				return false
			}
			_, err5 := o.QueryTable(dbUser).Filter("Id", odInfo[0].Passage.Id).Update(orm.Params{
				"OnRoadType": 0,
				"Balance": passengerBalance,
				"OrderNumAsP":orderNumAsP,
				"CancleOasP": passengerInfo[0].CancleOasP + 1,
			})
			if (err5 != nil) {
				logs.Error("update passenger info fail pid=%v" , odInfo[0].Passage.Id)
				o.Rollback()
				return false
			}
			afPassenger.Money = (allCost - kCost)
			afPassenger.Type = 4
			afPassenger.Balance = passengerBalance
			afPassenger.Time = strconv.FormatInt(currentTime, 10)
			_, err6 := o.Insert(&afPassenger)
			if (err6 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			afPassenger2.Money = kCost
			afPassenger2.Type = 6
			afPassenger2.Balance = passengerBalance
			afPassenger2.Time = strconv.FormatInt(currentTime, 10)
			_, err7 := o.Insert(&afPassenger2)
			if (err7 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}

			afDriver.Money = kCost - infoCost
			afDriver.Type = 7
			afDriver.Balance = driverBalance
			afDriver.Time = strconv.FormatInt(currentTime, 10)
			_, err8 := o.Insert(&afDriver)
			if (err8 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			afDriver2.Money = infoCost
			afDriver2.Type = 9
			afDriver2.Balance = driverBalance
			afDriver2.Time = strconv.FormatInt(currentTime, 10)
			_, err9 := o.Insert(&afDriver2)
			if (err9 != nil) {
				logs.Error("set account flow fail")
				o.Rollback()
				return false
			}
			is80Back = true
			weiyuejin = strconv.FormatFloat(kCost, 'G' , -1,64)
			shouxufei = strconv.FormatFloat(infoCost, 'G' , -1,64)
			driverIncome = strconv.FormatFloat(afDriver.Money, 'G' , -1,64)
			passengerIncome = strconv.FormatFloat(afPassenger.Money, 'G' , -1,64)
			driverB = strconv.FormatFloat(driverBalance, 'G' , -1,64)
			passengerB = strconv.FormatFloat(passengerBalance, 'G' , -1,64)
		}
	}



	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	if (isAllBack) {
		commonLib.SendMsg5(passengerInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#22c32e","车费退回",
			"#22c32e", "退回成功",
			"#173177", passengerIncome,
			"#173177", passengerB)
	} else if (is80Back) {
		commonLib.SendMsg5(passengerInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#22c32e","车费退回",
			"#22c32e", "退回成功",
			"#173177", passengerIncome,
			"#173177", passengerB)
		commonLib.SendMsg5(passengerInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#ff0000","扣除违约金20%",
			"#22c32e", "扣除成功",
			"#173177", weiyuejin,
			"#173177", passengerB)
		commonLib.SendMsg5(driverInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#22c32e","收款违约金20%",
			"#22c32e", "收款成功",
			"#173177", driverIncome+"(信息费：" + shouxufei + ")",
			"#173177", driverB)
	} else if (is0Back) {
		commonLib.SendMsg5(passengerInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#ff0000","扣除违约金100%",
			"#22c32e", "扣除成功",
			"#173177", weiyuejin,
			"#173177", passengerB)
		commonLib.SendMsg5(driverInfo[0].OpenId,
			4, "", "#173177", "", "",
			"#173177", "",
			"#22c32e","收款违约金100%",
			"#22c32e", "收款成功",
			"#173177", driverIncome+"(信息费：" + shouxufei + ")",
			"#173177", driverB)
	}

	commonLib.SendMsg5(odInfo[0].Driver.OpenId, 3, "http://www.youdidi.vip/Portal/driverorderdetail/"+odInfo[0].Order.Id,
		"#ff0000", "抱歉，乘客已操作取消行程", "系统以释放作为，请关注新乘客的预约",
		"#173177", odInfo[0].Passage.Nickname,
		"#173177", odInfo[0].Order.SrcId.Level1 + "-" + odInfo[0].Order.SrcId.Level2 + "-" + odInfo[0].Order.SrcId.Name,
		"#173177", odInfo[0].Order.DestId.Level1 + "-" + odInfo[0].Order.DestId.Level2 + "-" + odInfo[0].Order.DestId.Name,
		"#173177", "抱歉，行程临时有变",
		"#173177", time.Now().Format("2006-01-02 15:04"))

	return true
}

func (u *Order_detail) Recommand(odid string, uType string, starNum int, mark string, userId int) bool {
	o := orm.NewOrm()
	o.Begin()
	var dbUser User
	var userInfo []*User

	num, err := o.QueryTable(dbUser).Filter("Id", userId).ForUpdate().All(&userInfo)
	if (num < 1 || err != nil) {
		logs.Error("get userinfo fail uid=%v" , userId)
		o.Rollback()
		return false
	}

	if (uType == "0") {
		_, err1 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
			"IsDcommit": true,
			"DStarNum": starNum,
			"DCommit":mark,
		})
		if (err1 != nil) {
			logs.Error("update recommand fail odid=%v" , odid)
			o.Rollback()
			return false
		}
		_, err2 := o.QueryTable(dbUser).Filter("Id", userId).Update(orm.Params{
			"StarAsD": userInfo[0].StarAsD + starNum,
		})
		if (err2 != nil) {
			logs.Error("update userinfo fail uid=%v" , userId)
			o.Rollback()
			return false
		}
	} else {
		_, err1 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
			"IsPcommit": true,
			"PStarNum": starNum,
			"PCommit":mark,
		})
		if (err1 != nil) {
			logs.Error("update recommand fail odid=%v" , odid)
			o.Rollback()
			return false
		}
		_, err2 := o.QueryTable(dbUser).Filter("Id", userId).Update(orm.Params{
			"StarAsP": userInfo[0].StarAsP + starNum,
		})
		if (err2 != nil) {
			logs.Error("update userinfo fail uid=%v" , userId)
			o.Rollback()
			return false
		}
	}

	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail odid=%v" , odid)
		o.Rollback()
		return false
	}
	return true
}

func (u *Order_detail) CancleSingP(odid string, pid string) bool {
	o := orm.NewOrm()
	o.Begin()
	var dbUser User
	var userInfo []*User
	var odidInfo []*Order_detail
	var dbOrder Order
	var orderInfo []*Order

	num1, err1 := o.QueryTable(u).RelatedSel().Filter("Id", odid).Filter("Status__lt", 4).ForUpdate().All(&odidInfo)
	if (num1 < 1 || err1 != nil) {
		logs.Error("get odid info fail odid=%V" , odid)
		o.Rollback()
		return false
	}

	num2, err2 := o.QueryTable(dbOrder).RelatedSel().Filter("Id", odidInfo[0].Order.Id).ForUpdate().All(&orderInfo)
	if (num2 < 1 || err2 != nil) {
		logs.Error("get order info fail oid=%v" , odidInfo[0].Order.Id)
		o.Rollback()
		return false
	}

	num3, err3 := o.QueryTable(dbUser).RelatedSel().Filter("Id", odidInfo[0].Passage.Id).ForUpdate().All(&userInfo)
	if (num3 < 1 || err3 != nil) {
		logs.Error("get user info fail uid=%v" , odidInfo[0].Passage.Id)
		o.Rollback()
		return false
	}

	balance := userInfo[0].Balance
	money, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", orderInfo[0].Price * float64(odidInfo[0].SiteNum)), 64)
	balance = balance + money

	_,err4 := o.QueryTable(dbUser).Filter("Id", odidInfo[0].Passage.Id).Update(orm.Params{
		"Balance": balance,
		"OnRoadType" : 0,
	})
	if (err4 != nil) {
		logs.Error("update user info fail uid=%v" , odidInfo[0].Passage.Id)
		o.Rollback()
		return false
	}

	_, err5 := o.QueryTable(u).Filter("Id", odid).Update(orm.Params{
		"Status": 7,
	})
	if (err5 != nil) {
		logs.Error("update odid info fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	confirmNum := orderInfo[0].ConfirmPnum

	if (odidInfo[0].Status > 0) {
		confirmNum = confirmNum - odidInfo[0].SiteNum
	}

	_, err6 := o.QueryTable(dbOrder).Filter("Id", odidInfo[0].Order.Id).Update(orm.Params{
		"RequestPnum": orderInfo[0].RequestPnum - odidInfo[0].SiteNum,
		"ConfirmPnum": confirmNum,
		"CanclePnum" : orderInfo[0].CanclePnum + odidInfo[0].SiteNum,
	})
	if (err6 != nil) {
		logs.Error("update order info fail odid=%v" , odid)
		o.Rollback()
		return false
	}

	var dbAf Account_flow
	dbAf.Balance = balance
	dbAf.Money = money
	dbAf.User = odidInfo[0].Passage
	dbAf.Oid = orderInfo[0].Id
	dbAf.Type = 4
	dbAf.Time = strconv.FormatInt(time.Now().Unix(),10)

	_, err7 := o.Insert(&dbAf)
	if (err7 != nil) {
		logs.Error("insert account flow fail")
		o.Rollback()
		return false
	}


	errcommit := o.Commit()

	if (errcommit != nil) {
		logs.Error("commit fail odid=%v" , odid)
		o.Rollback()
		return false
	}
	return true
}