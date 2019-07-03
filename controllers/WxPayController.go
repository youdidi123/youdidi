package controllers

import (
    "fmt"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "github.com/objcoding/wxpay"
    "strconv"
    "time"
    "math"
)

type WxPayController struct {
    beego.Controller
}

type  PotalReturn struct {
    Code  int
    Msg   error
    Data  wxpay.Params
}

// @router /Portal/WxInvest/ [POST,GET]
func (c *WxPayController) WxInvest () {

    // 获取充值金额
    moneyStr := c.GetString("money")
    moneyFloat, err := strconv.ParseFloat(moneyStr,64)
    if (err != nil) {
        logs.Notice("Money input error :%s", err)
        err := fmt.Errorf("Money input error :%s", err)
        returnJson := &PotalReturn{-1, err, nil}
        c.Data["json"] = returnJson
        c.ServeJSON()
        return
    }
    moneyInt := int64(math.Floor(moneyFloat * 100 + 0.5))
    if (moneyInt <=  0) {
        logs.Notice("Money input error :Money can't be less than zero!")
        err := fmt.Errorf("Money can't be less than zero")
        returnJson := &PotalReturn{-1, err, nil}
        c.Data["json"] = returnJson
        c.ServeJSON()
        return
    }

    //Get 用户信息
    userLoginInfo, err := GetUserLoginInfoByCookie(c.Ctx)
    if (err != nil) {
        logs.Notice("Get User Info by Cookie Failed for:%s", err)
        err := fmt.Errorf("Get User Info by Cookie Failed for:%s", err)
        returnJson := &PotalReturn{-2, err, nil}
        c.Data["json"] = returnJson
        c.ServeJSON()
        return
    }
    logs.Debug("WxInvest userLoginInfo is:%s", userLoginInfo)

    //Get 订单信息
    id, _ := strconv.Atoi(userLoginInfo.idStr)
    payOrderId := genOrderId(id)
    timeStart := time.Now().Format("20060102150401")
    hh, _ := time.ParseDuration("1h")
    timeExpire := time.Now().Add(hh).Format("20060102150401")

    params := make(wxpay.Params)
    params.SetString("body", "长庆出行").
        SetString("device_info", "WEB").
        SetString("openid", userLoginInfo.OpenId).
        SetString("out_trade_no", payOrderId).
        SetString("time_start", timeStart).
        SetString("time_expire", timeExpire).
        SetInt64("total_fee", moneyInt).
        SetString("spbill_create_ip", c.Ctx.Input.IP()).
        SetString("notify_url", "http://www.youdidi.vip/Portal/WxInvestSuccess").
        SetString("trade_type", "JSAPI")


    // 日志答应订单信息
    logs.Info("Order Info:%s", params)

    // 调用统一下单
    jsapiParams, err:= WxUnifiedOrder(params)
    if (err != nil) {
        logs.Notice("Get WxUnifiedOrder Failed:%s", err)
        err := fmt.Errorf("Get WxUnifiedOrder Failed:%s", err)
        returnJson := &PotalReturn{-3, err, nil}
        c.Data["json"] = returnJson
        c.ServeJSON()
        return
    }


    returnJson := &PotalReturn{0, nil, *jsapiParams}
    c.Data["json"] = returnJson
    c.ServeJSON()
    return
}

/*
// @router /Portal/WxInvestSuccess/ [POST,GET]
func (c *WxPayController) WxInvestSuccess () {
    appId := beego.AppConfig.String("weixin::AppId")
    mchId := beego.AppConfig.String("weixin::MchId")
    apiKey := beego.AppConfig.String("weixin::apiKey")

    // 创建支付账户
    account := wxpay.NewAccount(appId, mchId, apiKey, true)

    // 新建微信支付客户端
    client := wxpay.NewClient(account)

    params, err := client.processResponseXml()

    //Get 用户信息
    userLoginInfo, err := GetUserLoginInfoByCookie(c.Ctx)
    if (err != nil) {
    }
}


func (c *Client) processResponseXml(xmlStr string) (Params, error) {
    var returnCode string
    params := XmlToMap(xmlStr)
    if params.ContainsKey("return_code") {
        returnCode = params.GetString("return_code")
    } else {
        return nil, errors.New("no return_code in XML")
    }
    if returnCode == Fail {
        return params, nil
    } else if returnCode == Success {
        if c.ValidSign(params) {
            return params, nil
        } else {
            return nil, errors.New("invalid sign value in XML")
        }
    } else {
        return nil, errors.New("return_code value is invalid in XML")
    }
}
*/

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