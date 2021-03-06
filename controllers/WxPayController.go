package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math"
	"strconv"
	"strings"
	"time"
	"youdidi/controllers/wxpay"
	"youdidi/models"
	"github.com/nanjishidu/gomini/gocrypto"
)

type WxPayController struct {
	beego.Controller
}

type PotalReturn struct {
	Code int
	Msg  string
	Data *wxpay.Params
}

// @router /Portal/WxInvest/ [POST, GET]
func (c *WxPayController) WxInvest() {

	// 获取充值金额
	moneyInt, err := moneyCHeck(c.GetString("money"))
	if err != nil {
		logs.Notice("Money input error :%s", err)
		c.jsonPotalReturn(-1, err.Error(), nil)
		return
	}

	//Get 用户信息
	userLoginInfo, err := GetUserLoginInfoByCookie(c.Ctx)
	if err != nil {
		logs.Notice("Get User Info by Cookie Failed for:%s", err)
		err := fmt.Errorf("Get User Info by Cookie Failed for:%s", err)
		c.jsonPotalReturn(-2, err.Error(), nil)
		return
	}
	logs.Debug("WxInvest userLoginInfo is:%s", userLoginInfo)

	//Get 订单信息
	userId, _ := strconv.Atoi(userLoginInfo.IdStr)
	investOrderId := genOrderId(userId)
	timeNow := time.Now()
	timeStart := timeNow.Format("20060102150401")
	hh, _ := time.ParseDuration("1h")
	timeExpire := timeNow.Add(hh).Format("20060102150401")

	params := make(wxpay.Params)
	params.SetString("body", "长庆出行").
		SetString("device_info", "WEB").
		SetString("openid", userLoginInfo.OpenId).
		SetString("out_trade_no", investOrderId).
		SetString("time_start", timeStart).
		SetString("time_expire", timeExpire).
		SetInt64("total_fee", moneyInt).
		SetString("spbill_create_ip", c.Ctx.Input.IP()).
		SetString("notify_url", "http://www.youdidi.vip/WxInvestSuccess").
		SetString("trade_type", "JSAPI")

	// 日志打印订单信息
	logs.Info("Order Info:%s", params)

	// 调用微信统一下单
	jsapiParams, err := WxUnifiedOrder(params)
	if err != nil {
		logs.Notice("Get WxUnifiedOrder Failed:%s", err)
		err := fmt.Errorf("Get WxUnifiedOrder Failed:%s", err)
		c.jsonPotalReturn(-3, err.Error(), nil)
		return
	}

	// 商户系统录入订单
	var cashFlowOrder models.Cash_flow
	cashFlowOrder.Id = investOrderId
	cashFlowOrder.Type = 0 // 0:充值
	cashFlowOrder.Money = float64(moneyInt) / 100
	cashFlowOrder.Status = 0 // 0:发起 1:成功 2:失败 3:拒绝
	cashFlowOrder.RefuseReason = ""
	cashFlowOrder.Time = strconv.Itoa(int(timeNow.Unix()))
	cashFlowOrder.WechatOrderId = jsapiParams.GetString("prepay_id")
	cashFlowOrder.User = &models.User{Id: userId}
	_, err = cashFlowOrder.Insert()
	if err != nil {
		err := fmt.Errorf("Insert cashFlowOrder to DB error:%s", err)
		logs.Notice(err.Error())
		c.jsonPotalReturn(-3, err.Error(), nil)
		return
	}

	c.jsonPotalReturn(0, "", jsapiParams)
	return
}

// 微信退款
// investOrderId 充值商户订单号
// wxInvestOrderId充值 微信订单号
// refundOrderId 退款订单号
// refundDesc 退款原因描述
func WxRefund(investOrderId string, wxInvestOrderId string,
	refundOrderId string, refundDesc string, refundFee float64) (string, error) {
	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId")
	apiKey := beego.AppConfig.String("weixin::apiKey")
	apiCert := beego.AppConfig.String("weixin::apiCert")

	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, false)

	// 设置证书
	account.SetCertData(apiCert)

	// 新建微信支付客户端
	client := wxpay.NewClient(account)

	// 获取退款订单号

	investOrder, err := InvestOrderCheck(investOrderId, wxInvestOrderId, refundOrderId)
	if (err != nil) {
		err = fmt.Errorf("InvestOrderCheck error for:%s", err)
		logs.Notice(err.Error())
		//c.jsonPotalReturn(-1, err.Error(), nil)
		return "", err
	}

	// 退款金额比较检查
	if ( investOrder.Money <  refundFee) {
		err := fmt.Errorf("The amount of refund:%f exceeds the total amount:%f",
			refundFee, investOrder.Money)
		logs.Notice(err.Error())
		//c.jsonPotalReturn(-1, err.Error(), nil)
		return "", err
	}

	// 获取退款原因
	//refundDesc := c.GetString("refundDesc")

	// 退款
	params := make(wxpay.Params)
	params.SetInt64("total_fee", int64(investOrder.Money*100)).
		SetInt64("refund_fee", int64(refundFee*100)).
		SetString("out_trade_no", investOrderId).
		SetString("out_refund_no", refundOrderId).
		SetString("transaction_id", wxInvestOrderId).
		SetString("notify_url", "http://www.youdidi.vip/WxRefundSuccess").
		SetString("refund_desc", refundDesc)

	responParam, err := client.Refund(params)
	if err != nil {
		//logs.Debug("xml response %s", res)
		err = fmt.Errorf("Refund request weixin err:%s", err)
		logs.Notice(err.Error())
		//c.jsonPotalReturn(-1, err.Error(), nil)
		return "", err
	}

	// debug 微信返回数据
	logs.Debug("WxUnifiedOrder Info:%s", responParam)

	if (responParam.GetString("return_code") == wxpay.Success) {
		if (responParam.GetString("result_code") == wxpay.Success) {
			return responParam.GetString("refund_id"), nil
		} else {
			return "", fmt.Errorf(responParam.GetString("err_code"))
		}
	} else {
		return "", fmt.Errorf("refund err return_code not success, return_msg=%v", responParam.GetString("return_msg"))
	}
}

// 充值成功微信回调调用接口
// @router /WxInvestSuccess [POST,GET]
func (c *WxPayController) WxInvestSuccess() {
	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId")
	apiKey := beego.AppConfig.String("weixin::apiKey")

	resPonse := make(wxpay.Params)

	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, false)

	// 新建微信支付客户端
	client := wxpay.NewClient(account)

	params, err := client.ProcessResponseXml(string(c.Ctx.Input.RequestBody))
	logs.Info("wxpay callback:%s", params)
	if err != nil {
		//校验签名失败
		logs.Info("wxpay callback:%s err:%v", params)
		resPonse.SetString("return_code", "FAIL").
			SetString("return_msg", "sign invaild")
		c.Ctx.WriteString(wxpay.MapToXml(resPonse))
		return
	}

	if (params.GetString("return_code") == wxpay.Success) {
		reAppId := params.GetString("appid")
		reMchId := params.GetString("mch_id")
		result_code := params.GetString("result_code")
		err_code := params.GetString("err_code")
		err_code_des := params.GetString("err_code_des")
		openid := params.GetString("openid")
		trade_type := params.GetString("trade_type")
		wxId := params.GetString("transaction_id") //微信订单号
		cfId := params.GetString("out_trade_no") //自己的订单号
		total_fee := params.GetInt64("total_fee") //金额
		transaction_id := params.GetString("transaction_id") //微信支付订单号

		if (reAppId != appId || reMchId != mchId || trade_type != "JSAPI") {
			logs.Info("appid or mchid or trade_typ is not mach appId=%v re=%v mchId=%v re=%v trade_typ=%v",
				appId,
				reAppId,
				reMchId,
				mchId,
				trade_type)
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "appid or mchid invaild")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

		var dbCf models.Cash_flow

		if (! dbCf.DealWxPayRe(result_code, err_code, err_code_des, openid, wxId, cfId, total_fee, transaction_id)) {
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "update database fail, please retry")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

	} else {
		//进入api主动验证
		resPonse.SetString("return_code", "FAIL").
			SetString("return_msg", "appid or mchid invaild")
		c.Ctx.WriteString(wxpay.MapToXml(resPonse))
		return
	}


	//在这里入库

	//if
	//investOrderId := params.GetString("appid	")

	resPonse.SetString("return_code", "SUCCESS").
		SetString("return_msg", "OK")
	c.Ctx.WriteString(wxpay.MapToXml(resPonse))
	return
}

// 退款成功的微信回调接口
// @router /WxRefundSuccess [POST,GET]
func (c *WxPayController) WxRefundSuccess() {
	logs.Info("WxRefund callback:%s", string(c.Ctx.Input.RequestBody))

	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId")
	apiKey := beego.AppConfig.String("weixin::apiKey")

	resPonse := make(wxpay.Params)

	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, false)

	// 新建微信支付客户端
	client := wxpay.NewClient(account)

	params, err := client.ProcessResponseXmlEnp(string(c.Ctx.Input.RequestBody))
	logs.Info("WxRefund callback:%s", params)
	if err != nil {
		//校验签名失败
		logs.Info("WxRefund callback:%s err:%v", params)
		resPonse.SetString("return_code", "FAIL").
			SetString("return_msg", "sign invaild")
		c.Ctx.WriteString(wxpay.MapToXml(resPonse))
		return
	}

	if (params.GetString("return_code") == wxpay.Success) {
		reAppId := params.GetString("appid")
		reMchId := params.GetString("mch_id")

		req_info := params.GetString("req_info")

		req_info64, err := base64.StdEncoding.DecodeString(req_info)
		if (err != nil) {
			logs.Error("wx refund req_info decode base64 fail err=%v", err.Error())
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "appid or mchid invaild")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

		gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(apiKey)))

		plaintext, err := gocrypto.AesECBDecrypt(req_info64)
		if err != nil {
			logs.Error("decode AesECBDecrypt fail err=%v", err.Error())
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "appid or mchid invaild")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

		paramsJM := wxpay.XmlToMap(string(plaintext))
		refund_id := paramsJM.GetString("refund_id")
		out_refund_no := paramsJM.GetString("out_refund_no")
		refund_status := paramsJM.GetString("refund_status")
		success_time := paramsJM.GetString("success_time")
		settlement_refund_fee := paramsJM.GetInt64("settlement_refund_fee")


		if (reAppId != appId || reMchId != mchId ) {
			logs.Info("appid or mchid or trade_typ is not mach appId=%v re=%v mchId=%v re=%v",
				appId,
				reAppId,
				reMchId,
				mchId,
				)
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "appid or mchid invaild")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

		var dbCf models.Cash_flow

		if (! dbCf.DealWxRefundRe(refund_id, out_refund_no, refund_status, success_time, settlement_refund_fee)) {
			resPonse.SetString("return_code", "FAIL").
				SetString("return_msg", "update database fail, please retry")
			c.Ctx.WriteString(wxpay.MapToXml(resPonse))
			return
		}

	} else {
		reMsg := params.GetString("return_msg")
		logs.Debug("wx refund retrun msg fail err=%v", reMsg)
		//进入api主动验证
		resPonse.SetString("return_code", "FAIL").
			SetString("return_msg", "appid or mchid invaild")
		c.Ctx.WriteString(wxpay.MapToXml(resPonse))
		return
	}

	resPonse.SetString("return_code", "SUCCESS").
		SetString("return_msg", "OK")
	c.Ctx.WriteString(wxpay.MapToXml(resPonse))
	return

}

// 企业付款提现的接口
// cashOutAmount 提现金额人民币 单位 分 partnerTradeNo 商户用户提现订单
// clientIP  用户IP desc 提现原因描述
func WxEnpTransfers(cashOutAmount int64, OpenId string, partnerTradeNo string,
	clientIP string, desc string) (string, error) {
	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId1")
	apiKey := beego.AppConfig.String("weixin::apiKey1")
	apiCert := beego.AppConfig.String("weixin::apiCert1")


	if cashOutAmount <= 0 {
		return "", fmt.Errorf("cashOutAmount can't be less than zero!")
	}

	// 订单号校验
	if (!IsNum(partnerTradeNo)) {
		err := fmt.Errorf("Illegal Partner Trade Number:%s", partnerTradeNo)
		logs.Notice(err.Error())
		return "", err
	}

	// 企业号提现到个人零钱账户
	params := make(wxpay.Params)
	params.SetInt64("amount", cashOutAmount).
		SetString("partner_trade_no", partnerTradeNo).
		SetString("openid", OpenId).
		SetString("check_name", "NO_CHECK").
		SetString("desc", desc).
		SetString("spbill_create_ip", clientIP)


	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, false)

	// 设置证书
	account.SetCertData(apiCert)

	// 新建微信支付客户端
	client := wxpay.NewClient(account)

	rParams, err := client.EnpTransfers(params)
	if err != nil {
		return "", fmt.Errorf("EnpTransfers err:%s", err)
	}

	if (rParams.GetString("return_code") == wxpay.Success) {
		if (rParams.GetString("result_code") == wxpay.Success) {
			return rParams.GetString("payment_no"), nil
		} else {
			return "", fmt.Errorf(rParams.GetString("err_code"))
		}
	} else {
		return "", fmt.Errorf("EnpTransfers err return_code not success, return_msg=%v", rParams.GetString("return_msg"))
	}

}

// 退款时充值订单检查
func InvestOrderCheck(investOrderId string, wxInvestOrderId string,
	refundOrderId string) (*models.Cash_flow, error) {
	var cashFlowOrder models.Cash_flow
	var cashFlowOrders []*models.Cash_flow

	_, num := cashFlowOrder.GetOrderInfo(investOrderId, &cashFlowOrders)
	if num == 0 {
		err := fmt.Errorf("InvestOrderId not exist and can not get order info from db ")
		return nil, err
	}

	if (num > 1) {
		err := fmt.Errorf("InvestOrderId:%s repeat in db", investOrderId)
		return nil, err
	}

	if (cashFlowOrders[0].WechatOrderId != wxInvestOrderId){
		err := fmt.Errorf("wxInvestOrderId:%s not equal to WechatOrderId:5s in cashFlowOrders",
			wxInvestOrderId, cashFlowOrders[0].WechatOrderId)
		return nil, err
	}

	return  cashFlowOrders[0], nil
}

// 数字字符串判断，检查订单号
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}



// return Json result
func (c *WxPayController) jsonPotalReturn(code int, msg string,
	mapData *wxpay.Params) {
	returnJson := &PotalReturn{code, msg, mapData}
	c.Data["json"] = returnJson
	c.ServeJSON()
}

// 支付金额校验转换
func moneyCHeck(moneyStr string) (int64, error) {
	moneyFloat, err := strconv.ParseFloat(moneyStr, 64)
	if err != nil {
		return 0, err
	}
	moneyInt := int64(math.Floor(moneyFloat*100 + 0.5))
	if moneyInt <= 0 {
		return 0, fmt.Errorf("Money can't be less than zero!")
	}

	return moneyInt, err
}

// 微信支付统一下单接口
func WxUnifiedOrder(params wxpay.Params) (*wxpay.Params, error) {
	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId")
	apiKey := beego.AppConfig.String("weixin::ApiKey")

	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, false)

	// 新建微信支付客户端
	client := wxpay.NewClient(account)

	// 设置http请求超时时间
	client.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	client.SetHttpReadTimeoutMs(1000)

	// 统一下单
	prepayParams, err := client.UnifiedOrder(params)
	if err != nil {
		//logs.Debug("xml response %s", res)
		return nil, fmt.Errorf("UnifiedOrder err:%s", err)
	}

	//logs.Debug("xml response %s", res)

	// debug 微信返回数据
	logs.Debug("WxUnifiedOrder Info:%s", prepayParams)

	// 校验 prepay_id 是否存在
	if !prepayParams.ContainsKey("prepay_id") {
		return nil, fmt.Errorf("No prepay_id return in UnifiedOrder err:%s", err)
	}

	// 组装json数据 返回前端
	prepayId := prepayParams.GetString("prepay_id")
	jsapiParams := make(wxpay.Params)
	jsapiParams.SetString("appId", appId).
		SetString("timeStamp", strconv.FormatInt(time.Now().Unix(), 10)).
		SetString("nonceStr", strconv.FormatInt(time.Now().UTC().UnixNano(), 10)).
		SetString("package", "prepay_id="+prepayId).
		SetString("signType", "MD5") //统一使用MD5 签名

	// 数据签名
	sign := client.Sign(jsapiParams)
	jsapiParams.SetString("paySign", sign)

	return &jsapiParams, nil
}
