package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math"
	"strconv"
	"time"
	"youdidi/controllers/wxpay"
	"youdidi/models"
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

// @router /WxInvestSuccess [POST,GET]
func (c *WxPayController) WxInvestSuccess() {
	appId := beego.AppConfig.String("weixin::AppId")
	mchId := beego.AppConfig.String("weixin::MchId")
	apiKey := beego.AppConfig.String("weixin::apiKey")

	resPonse := make(wxpay.Params)

	// 创建支付账户
	account := wxpay.NewAccount(appId, mchId, apiKey, true)

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
	}


	//在这里入库

	//if
	//investOrderId := params.GetString("appid	")

	resPonse.SetString("return_code", "SUCCESS").
		SetString("return_msg", "OK")
	c.Ctx.WriteString(wxpay.MapToXml(resPonse))
	return
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
