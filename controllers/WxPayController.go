package controllers

import (
    "fmt"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "youdidi/controllers/wxpay"
    "strconv"
    "time"
    "math"
    "youdidi/models"
)

type WxPayController struct {
    beego.Controller
}

type  PotalReturn struct {
    Code  int
    Msg   string
    Data  wxpay.Params
}

// @router /Portal/WxInvest/ [POST, GET]
func (c *WxPayController) WxInvest () {

    // 获取充值金额
    moneyInt, err := moneyCHeck(c.GetString("money"))
    if (err != nil ) {
        logs.Notice("Money input error :%s", err)
        c.jsonPotalReturn(-1, err.Error(), nil)
        return
    }

    //Get 用户信息
    userLoginInfo, err := GetUserLoginInfoByCookie(c.Ctx)
    if (err != nil) {
        logs.Notice("Get User Info by Cookie Failed for:%s", err)
        err := fmt.Errorf("Get User Info by Cookie Failed for:%s", err)
        c.jsonPotalReturn(-2, err.Error(), nil)
        return
    }
    logs.Debug("WxInvest userLoginInfo is:%s", userLoginInfo)

    //Get 订单信息
    userId, _ := strconv.Atoi(userLoginInfo.idStr)
    investOrderId := genOrderId(userId)
    timeStart := time.Now().Format("20060102150401")
    hh, _ := time.ParseDuration("1h")
    timeExpire := time.Now().Add(hh).Format("20060102150401")

    params := make(wxpay.Params)
    params.SetString("body", "长庆出行").
        SetString("device_info", "WEB").
        SetString("openid", userLoginInfo.OpenId).
        SetString("out_trade_no", investOrderId).
        SetString("time_start", timeStart).
        SetString("time_expire", timeExpire).
        SetInt64("total_fee", moneyInt).
        SetString("spbill_create_ip", c.Ctx.Input.IP()).
        SetString("notify_url", "http://www.youdidi.vip/Portal/WxInvestSuccess").
        SetString("trade_type", "JSAPI")

    // 日志打印订单信息
    logs.Info("Order Info:%s", params)

    // 调用微信统一下单
    jsapiParams, err:= WxUnifiedOrder(params)
    if (err != nil) {
        logs.Notice("Get WxUnifiedOrder Failed:%s", err)
        err := fmt.Errorf("Get WxUnifiedOrder Failed:%s", err)
        c.jsonPotalReturn(-3, err.Error(), nil)
        return
    }

    // 商户系统录入订单
    var cashFlowOrder models.Cash_flow
    cashFlowOrder.Id = investOrderId
    cashFlowOrder.Type = 0 // 0:充值
    cashFlowOrder.Money = float64(moneyInt)/100
    cashFlowOrder.Status = 0 // 0:发起 1:成功 2:失败 3:拒绝
    cashFlowOrder.RefuseReason = ""
    cashFlowOrder.Time = timeStart
    cashFlowOrder.WechatOrderId = jsapiParams.GetString("prepay_id")
    cashFlowOrder.User = &models.User{Id:userId}
    _, err = cashFlowOrder.Insert()
    if (err != nil) {
        err := fmt.Errorf("Insert cashFlowOrder to DB error:%s", err)
        logs.Notice(err.Error())
        c.jsonPotalReturn(-3, err.Error(), nil)
        return
    }

    c.jsonPotalReturn(0, "", jsapiParams)
    return
}

// @router /WxInvestSuccess/ [POST, GET]
func (c *WxPayController) WxInvestSuccess () {
    appId := beego.AppConfig.String("weixin::AppId")
    mchId := beego.AppConfig.String("weixin::MchId")
    apiKey := beego.AppConfig.String("weixin::apiKey")

    // 创建支付账户
    account := wxpay.NewAccount(appId, mchId, apiKey, true)

    // 新建微信支付客户端
    client := wxpay.NewClient(account)

    params, err := client.ProcessResponseXml(string(c.Ctx.Input.RequestBody))
    if (err != nil) {
        logs.Info("wxpay callback:%s", params)
    }

    //if
    //investOrderId := params.GetString("appid	")
    resPonse := make(wxpay.Params)
    resPonse.SetString("return_code", "SUCCESS").
        SetString("return_msg", "OK")
    c.Ctx.WriteString(wxpay.MapToXml(resPonse))
    return
}

// return Json result
func (c *WxPayController) jsonPotalReturn (code int, msg string,
    mapData *wxpay.Params) {
    returnJson := &PotalReturn{code, msg, *mapData}
    c.Data["json"] = returnJson
    c.ServeJSON()
}

// 支付金额校验转换
func moneyCHeck(moneyStr string) (int64, error){
    moneyFloat, err := strconv.ParseFloat(moneyStr, 64)
    if (err != nil) {
        return 0, err
    }
    moneyInt := int64(math.Floor(moneyFloat*100 + 0.5))
    if (moneyInt <= 0) {
        return 0, fmt.Errorf("Money can't be less than zero!")
    }

    return moneyInt, err
}

// 微信支付统一下单接口
func WxUnifiedOrder (params wxpay.Params) (*wxpay.Params, error){
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
    if (err != nil) {
        //logs.Debug("xml response %s", res)
        return nil, fmt.Errorf("UnifiedOrder err:%s", err)
    }

    //logs.Debug("xml response %s", res)

    // debug 微信返回数据
    logs.Debug("WxUnifiedOrder Info:%s", prepayParams)

    // 校验 prepay_id 是否存在
    if ( !prepayParams.ContainsKey("prepay_id") ) {
        return nil, fmt.Errorf("No prepay_id return in UnifiedOrder err:%s", err)
    }

    // 组装json数据 返回前端
    prepayId := prepayParams.GetString("prepay_id")
    jsapiParams :=  make(wxpay.Params)
    jsapiParams.SetString("appId", appId).
        SetString("timeStamp", strconv.FormatInt(time.Now().Unix(),10)).
        SetString("nonceStr", strconv.FormatInt(time.Now().UTC().UnixNano(), 10)).
        SetString("package", "prepay_id="+prepayId).
        SetString("signType", "MD5") //统一使用MD5 签名

    // 数据签名
    sign := client.Sign(jsapiParams)
    jsapiParams.SetString("paySign", sign)

    return &jsapiParams, nil
}