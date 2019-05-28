package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "net/url"
    "net/http"
    "io/ioutil"
)

type WxLoginController struct {
	beego.Controller
}

// @router /WxLogin/ [POST,GET]
func (c *WxLoginController) WxLogin () {
	appId := beego.AppConfig.String("weixin::AppId")

	oauth2Url := "https://open.weixin.qq.com/connect/oauth2/authorize"
	backUrl := "http://www.youdidi.vip/GetWechatGZAccessToken"

	redirectURL := oauth2Url +
		"?appid=" + appId +
		"&redirect_uri=" + backUrl +
		"&response_type=code" +
		"&scope=snsapi_userinfo" +
		"&state=STATE#wechat_redirect"

	c.Ctx.Redirect(302, redirectURL)
	logs.Notice("redirectURL is :%s", redirectURL)
}

// @router /GetWechatGZAccessToken/ [POST,GET]
func (c *WxLoginController) GetWechatGZAccessToken () {
	appId := beego.AppConfig.String("weixin::AppId")
	appSecret := beego.AppConfig.String("weixin::AppSecret")
	//token := beego.AppConfig.String("weixin::Token")

    code := c.GetString("code")

    u := url.Values{}
    u.Set("appid", appId)
    u.Set("secret", appSecret)
    u.Set("code", code)
    u.Set("grant_type", "authorization_code")


	accessTokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?" + u.Encode()
	//userInfoUrl := "https://api.weixin.qq.com/sns/userinfo"

    c.Ctx.WriteString(accessTokenUrl)
	logs.Notice("accessTokenUrl is :%s", accessTokenUrl)

    resp, err := http.Get(accessTokenUrl)
    if err != nil {
        logs.Error("Get accessTokenUrl Failed:%s", err)
        //c.Ctx.WriteString('Has an error')
        return
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logs.Error("read resp Body failed:%s", err)
        //c.Ctx.WriteString('Has an error')
        return
    }

    c.Ctx.WriteString(string(body))
}
