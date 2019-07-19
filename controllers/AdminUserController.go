package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"youdidi/commonLib"
	"youdidi/models"
	"youdidi/redisClient"
)

type AdminUserController struct {
	beego.Controller
}

type AdminUserLoginInfo struct {
	Id int
	Name string
	Token string
	Type int
}

var (
	AdminLoginPeriod = 1* 24 * 60 * 60 //用户登陆有效期1天
	AdminLoginPrefix = "ADMIN_LOGIN_INFO_"  //缓存在redis中的用户数据key前缀
)


// @router /AdminLogin/ [GET]
func (this *AdminUserController) AdminLogin (){
	this.TplName = "adminLogin.html"
}

// @router /admin/ [GET]
func (this *AdminUserController) Admin (){
	this.TplName = "adminHomepage.html"
}

// @router /AdmindoLogin [POST]
func (this *AdminUserController) AdminDoLogin (){
	logs.Debug("wo zai zhe")
	inputName := this.GetString("name")
	inputPasswd := this.GetString("passwd")

	var dbAdminUser models.Admin_user
	var userInfo []*models.Admin_user


	num := dbAdminUser.GetUserInfoFromName(inputName, &userInfo)
	if (num < 1) {
		logs.Debug("can not find user name=%v", inputName)
		this.Redirect("/AdminLogin/", 302)
		return
	}
	h := md5.New()
	h.Write([]byte(inputPasswd))

	passwdMd5 := hex.EncodeToString(h.Sum(nil))

	if (passwdMd5 != userInfo[0].Passwd) {
		logs.Error("input passwd is not correct inputpasswd=%v md5=%v dbvalue=%v", inputPasswd, passwdMd5, userInfo[0].Passwd)
		this.Redirect("/AdminLogin/", 302)
		return
	}

	token := getToken(inputName , inputPasswd)

	userInfoRedis := &AdminUserLoginInfo{}
	userInfoRedis.Id = userInfo[0].Id
	userInfoRedis.Name = userInfo[0].Name
	userInfoRedis.Type= userInfo[0].Type
	userInfoRedis.Token = token

	data, _ := json.Marshal(userInfoRedis)
	idStr := strconv.Itoa(userInfo[0].Id)

	redisClient.SetKey(AdminLoginPrefix+idStr , string(data))
	redisClient.Setexpire(AdminLoginPrefix+idStr , AdminLoginPeriod)

	this.Ctx.SetSecureCookie("qyt","qyt_admin_id" , idStr) //注入用户id，后续所有用户id都从cookie里获取
	this.Ctx.SetSecureCookie("qyt","qyt_admin_token" , token)

	this.TplName = "adminHomepage.html"
}

// @router /admin/dconfirm [GET]
func (this *AdminUserController) DriverConfirm (){
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetNoConfirm(&dcInfo)

	for i, v := range dcInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		dcInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = dcInfo

	this.TplName = "adminDriverConfirm.html"
}

// @router /admin/confirmDriverDetail/:id [GET]
func (this *AdminUserController) ConfirmDriverDetail (){
	id := this.GetString(":id")
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetOrderFromId(id, &dcInfo)

	for i, v := range dcInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		dcInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = dcInfo[0]

	this.TplName = "adminConfirmDriverDetail.html"
}

// @router /admin/userwithdrew [GET]
func (this *AdminUserController) UserWithdrew () {
	var dbCashFlow models.Cash_flow
	var cfWithdrew []*models.Cash_flow
	var cfWithdrewRefund []*models.Cash_flow
	var cfWithdrewError []*models.Cash_flow
	var cfWithdrewRefundError []*models.Cash_flow
	var cfWithdrewProcess []*models.Cash_flow

	tm, _ := commonLib.GetTodayBeginTime()
	endTime := tm - (24 * 60 * 60 * 1)
	this.Data["endTime"] = time.Unix(endTime,0).Format("2006-01-02 15:04")

	sum1 := 0.00
	sum2 := 0.00


	num1 := dbCashFlow.GetWithdrewOrder(&cfWithdrew, endTime, 1, 0)
	num2 := dbCashFlow.GetWithdrewOrder(&cfWithdrewRefund, endTime, 2, 0)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewError, endTime, 1, 2)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewRefundError, endTime, 2, 2)
	dbCashFlow.GetWithdrewOrder(&cfWithdrewProcess, endTime, 2, 4)

	for i, v := range cfWithdrew {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrew[i].Time = tm.Format("2006-01-02 15:04")
		sum1 += v.Money
	}

	for i, v := range cfWithdrewRefund {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrewRefund[i].Time = tm.Format("2006-01-02 15:04")
		sum2 += v.Money
	}

	for i, v := range cfWithdrewError {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		finishTime64, _ := strconv.ParseInt(v.FinishTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		tm1 := time.Unix(finishTime64, 0)
		cfWithdrewError[i].Time = tm.Format("2006-01-02 15:04")
		cfWithdrewError[i].FinishTime = tm1.Format("2006-01-02 15:04")
	}

	for i, v := range cfWithdrewRefundError {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		finishTime64, _ := strconv.ParseInt(v.FinishTime, 10, 64)
		tm := time.Unix(launchTime64, 0)
		tm1 := time.Unix(finishTime64, 0)
		cfWithdrewRefundError[i].Time = tm.Format("2006-01-02 15:04")
		cfWithdrewRefundError[i].FinishTime = tm1.Format("2006-01-02 15:04")
	}

	for i, v := range cfWithdrewProcess {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cfWithdrewProcess[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num1"] = num1
	this.Data["num2"] = num2
	this.Data["sum1"] = sum1
	this.Data["sum2"] = sum2

	this.Data["listw"] = cfWithdrew
	this.Data["listwr"] = cfWithdrewRefund
	this.Data["listwe"] = cfWithdrewError
	this.Data["listwre"] = cfWithdrewRefundError
	this.Data["listwp"] = cfWithdrewProcess

	this.TplName = "adminUserWithdrew.html"
}

// @router /admin/doconfirmdriver [POST]
func (this *AdminUserController) DoConfirmDriver () {
	oid := this.GetString("oid")
	aType := this.GetString("type")
	mark := this.GetString("mark")

	var dbDc models.Driver_confirm

	code := 0
	msg := ""

	if (! dbDc.DoConfirmDriver(oid, aType, mark)) {
		code = 1
		msg = "系统异常，请重试"
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

// @router /admin/showcomplain [GET]
func (this *AdminUserController) ShowComplain () {
	var dbC models.Complain
	var cInfo []*models.Complain

	num, _ := dbC.GetNoComplain(&cInfo)

	for i, v := range cInfo {
		this.Data["launchTime"] = v.Time;
		launchTime64, _ := strconv.ParseInt(v.Time, 10, 64)
		tm := time.Unix(launchTime64, 0)
		cInfo[i].Time = tm.Format("2006-01-02 15:04")
	}

	this.Data["num"] = num
	this.Data["list"] = cInfo
	this.TplName = "adminShowComplain.html"
}

// @router /admin/dealwithdrew [POST]
func (this *AdminUserController) DealWithDrew () {
	oids := this.GetStrings("oid")
	logs.Debug("withdrew list %v", oids)
	var dbCf models.Cash_flow

	code := 0
	msg := ""

	for _, oid := range oids {
		logs.Debug("withdrew id %v", oid)
		var cfInfo []* models.Cash_flow
		_, num := dbCf.GetOrderInfo(oid, &cfInfo)
		if (num != 1) {
			logs.Error("withdrew order id has something wrong oid=%v reNum=%v", oid, num)
			continue
		}
		if (cfInfo[0].Type == 1) {
			wxId, err := WxEnpTransfers(int64(cfInfo[0].Money * 100), cfInfo[0].User.OpenId, oid, "192.168.0.1", "长庆出行账户提现")
			if (err != nil) {
				_, err := dbCf.UpdateWithDrewResult(false, oid, "", err.Error(), 2)
				if (err != nil) {
					logs.Error("update withdrew result fail oid=%v result=fail err=%v", oid, err.Error())
				}
			} else {
				_, err := dbCf.UpdateWithDrewResult(true, oid, wxId, "", 1)
				if (err != nil) {
					logs.Error("update withdrew result fail oid=%v result=success err=%v", oid, err.Error())
				}
				moneyStr := strconv.FormatFloat(cfInfo[0].Money, 'G' , -1,64)
				balanceStr := strconv.FormatFloat(cfInfo[0].User.Balance, 'G' , -1,64)
				commonLib.SendMsg5(cfInfo[0].User.OpenId,
					4, "", "#173177", "", "",
					"#173177", "",
					"#22c32e","提现确认",
					"#22c32e", "提现成功",
					"#173177", moneyStr,
					"#173177", balanceStr)
			}
		} else if (cfInfo[0].Type == 2) {
			var investInfo []*models.Cash_flow
			_, num1 := dbCf.GetOrderInfo(cfInfo[0].InvestOid, &investInfo)
			if (num1 != 1) {
				logs.Error("get invest info from refund fail refundid=%v investid=%v", oid, cfInfo[0].InvestOid)
				continue
			}
			if (cfInfo[0].Money != investInfo[0].Money) {
				logs.Error("refund money:%v not equal invset money:%v", cfInfo[0].Money, investInfo[0].Money)
				continue
			}
			wxId, err := WxRefund(cfInfo[0].InvestOid, investInfo[0].WechatOrderId, oid, "长庆出行账户提现", cfInfo[0].Money)
			if (err != nil) {
				_, err := dbCf.UpdateWithDrewResult(false, oid, "", err.Error(), 2)
				if (err != nil) {
					logs.Error("update withdrew result fail oid=%v result=fail err=%v", oid, err.Error())
				}
			} else {
				_, err := dbCf.UpdateWithDrewResult(true, oid, wxId, "", 4)
				if (err != nil) {
					logs.Error("update withdrew result fail oid=%v result=success err=%v", oid, err.Error())
				}
			}
		}

	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()
}

// @router /admin/report [GET]
func (this *AdminUserController) Report () {
	var dbCashFlow models.Cash_flow
	var dbUser models.User
	var dbDc models.Driver_confirm
	var dbCom models.Complain
	var dbOrder models.Order

	userNum, _ := dbUser.GetUserNum()
	userNum = userNum - 10
	this.Data["userNnum"] = userNum //平台当前的用户数量

	driverNum, _ := dbUser.GetDriverNum()
	this.Data["driverNum"] = driverNum //平台当前的车主数量


	var cfWithdrew []*models.Cash_flow
	var cfWithdrewRefund []*models.Cash_flow
	var cfWithdrewError []*models.Cash_flow
	var cfWithdrewRefundError []*models.Cash_flow
	var cfWithdrewProcess []*models.Cash_flow

	tm, _ := commonLib.GetTodayBeginTime()
	tm1 := tm - (24 * 60 * 60)
	this.Data["today"] = time.Unix(tm,0).Format("2006-01-02 15:04") //结束时间
	this.Data["yesterday"] = time.Unix(tm1,0).Format("2006-01-02 15:04") //开始时间

	sum_withdrew := 0.00
	sum_withdrew_refund := 0.00
	sum_withdrew_error := 0.00
	sum_withdrew_refund_error := 0.00
	sum_withdrew_refund_process := 0.00


	num_withdrew := dbCashFlow.GetWithdrewOrder(&cfWithdrew, tm, 1, 0)
	num_withdrew_refund := dbCashFlow.GetWithdrewOrder(&cfWithdrewRefund, tm, 2, 0)
	num_withdrew_error := dbCashFlow.GetWithdrewOrder(&cfWithdrewError, tm, 1, 2)
	num_withdrew_refund_error := dbCashFlow.GetWithdrewOrder(&cfWithdrewRefundError, tm, 2, 2)
	num_withdrew_refund_process := dbCashFlow.GetWithdrewOrder(&cfWithdrewProcess, tm, 2, 4)

	for _, v := range cfWithdrew {
		sum_withdrew += v.Money
	}

	for _, v := range cfWithdrewRefund {
		sum_withdrew_refund += v.Money
	}

	for _, v := range cfWithdrewError {
		sum_withdrew_error += v.Money
	}

	for _, v := range cfWithdrewRefundError {
		sum_withdrew_refund_error += v.Money
	}

	for _, v := range cfWithdrewProcess {
		sum_withdrew_refund_process += v.Money
	}



	this.Data["sum_withdrew"] = sum_withdrew //通过支付方式提现的金额
	this.Data["sum_withdrew_refund"] = sum_withdrew_refund //通过退款方式提现的金额
	this.Data["sum_withdrew_error"] = sum_withdrew_error //通过支付方式提现失败的金额
	this.Data["sum_withdrew_refund_error"] = sum_withdrew_refund_error //通过退款方式提现失败的金额
	this.Data["sum_withdrew_refund_process"] = sum_withdrew_refund_process //还没有收到微信确认退款的金额
	this.Data["num_withdrew"] = num_withdrew //通过支付方式提现的单数
	this.Data["num_withdrew_refund"] = num_withdrew_refund //通过退款方式提现的单数
	this.Data["num_withdrew_error"] = num_withdrew_error //通过支付方式提现失败的单数
	this.Data["num_withdrew_refund_error"] = num_withdrew_refund_error //通过退款方式提现失败的单数
	this.Data["num_withdrew_refund_process"] = num_withdrew_refund_process //微信还没有确认退款的单数

	investMoney, investNum := dbCashFlow.GetInvestMoney(tm1, tm, 1)
	this.Data["investMoney"] = investMoney //昨天充值的总共多少钱
	this.Data["investNum"] = investNum //昨天充值的总共多少单

	investMoneyCNF, investNumCNF := dbCashFlow.GetInvestMoney(tm1, tm, 0)
	this.Data["investMoneyCNF"] = investMoneyCNF //昨天发起充值但未完成的金额
	this.Data["investNumCNF"] = investNumCNF //昨天发起充值但未完成的单子数

	investMoneyF, investNumF := dbCashFlow.GetInvestMoney(tm1, tm, 2)
	this.Data["investMoneyF"] = investMoneyF //昨天充值失败的金额
	this.Data["investNumF"] = investNumF //昨天充值的单子数

	var tmp []*models.Driver_confirm
	dcNum := dbDc.GetNoConfirm(&tmp)
	this.Data["dcNum"] = dcNum //待处理的司机身份验证单数

	var tmp1 []*models.Complain
	cNUm, _ := dbCom.GetNoComplain(&tmp1)
	this.Data["cNUm"] = cNUm //待处理投诉单数

	orderNum, requestNum, confirmNum := dbOrder.GetOrderNum(tm1, tm)
	this.Data["orderNum"] = orderNum //昨天车主发起的行程数
	this.Data["requestNum"] = requestNum //昨天请求拼车的乘客数
	this.Data["confirmNum"] = confirmNum //昨天车主确认请求的乘客数

	this.TplName = "adminReport.html"

}