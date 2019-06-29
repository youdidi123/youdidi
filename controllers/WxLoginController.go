package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "fmt"
	"youdidi/models"
)

type WxLoginController struct {
	beego.Controller
}

type  AccessToken struct {
    Access_token  string
    Expires_in    int
    Refresh_token string
    Openid        string
    Scope         string
    Errcode       int
    Errmsg        string
}

type UserInfo struct {
    Openid        string
    Nickname      string
    Sex           int
    Province      string
    City          string
    Country       string
    Headimgurl    string
    Priviledge    []interface{}
    Unionid       string
    Errcode       int
    Errmsg        string
}

// @router /Wxtest/ [POST,GET]
func (c *WxLoginController) Wxtest () {
    var userInfo UserInfo
	userInfo.Openid = "Openid-test_111111"
	userInfo.Nickname = "Nickname-test_111111"
	userInfo.Sex = 0
	userInfo.Province = "Province-test_111111"
	userInfo.City = "City-test_111111"
	userInfo.Headimgurl = "http://thirdwx.qlogo.cn/mmopen/vi_32/0NNL244MVpxDRwPj3gScx"+
		"6UbLCVjmqPtaHbkKIicxFplEkicOLuwyz42Ip40bP8Lw2ibwA4Vu9LBZvtKn70AicR3cg/132"
	userInfo.Unionid = "Unionid-test_111111"
	userInfo.Province = "Province-test_111111"

	//c.WxDologon(&userInfo)
    c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html;charset=utf-8")
	c.Ctx.WriteString("<img src=\"http://thirdwx.qlogo.cn/mmopen/vi_32/0NNL244MVpxDRwPj3gScx6UbLCVjmqPtaHbkKIicxFplEkicOLuwyz42Ip40bP8Lw2ibwA4Vu9LBZvtKn70AicR3cg/132\" alt=\"test\" />")
}

// @router /WxLogin/ [POST,GET]
func (c *WxLoginController) WxLogin () {
	appId := beego.AppConfig.String("weixin::AppId")

	oauth2Url := "https://open.weixin.qq.com/connect/oauth2/authorize"
	backUrl := "http://www.youdidi.vip/UserInfoCheck"

	redirectURL := oauth2Url +
		"?appid=" + appId +
		"&redirect_uri=" + backUrl +
		"&response_type=code" +
		"&scope=snsapi_userinfo" +
		"&state=STATE#wechat_redirect"

	c.Ctx.Redirect(302, redirectURL)
	logs.Notice("redirectURL is :%s", redirectURL)
}

// @router /UserInfoCheck/ [POST,GET]
func (c *WxLoginController) UserInfoCheck () {
    code := c.GetString("code")

    accessToken, err := GetWechatGZAccessToken(code)
    if err != nil {
        logs.Error("Get accessToken Failed:%s", err)
        c.Abort("401")
        return
    }
    //var accessToken AccessToken
    //accessToken.Access_token = "22_uKo3_E_UxGVOlfAaMR-vz_fL8BlkmZU09f3J-WFh06wPkHaa5GrVKGQp1QUVnwvuD-1K723rIAGZgJj-QhkLAxPLtdZqPMYV49jvUYRYzHI"
    //accessToken.Openid = "ooafc5o6_Jkfgk8BH9VobbfQzz6U"
    userInfo, err := WxGetUserInfo(accessToken)
    if err != nil {
        logs.Error("Get userInfo Failed:%s", err)
        c.Abort("401")
        return
    }

    //数据库注册微信登陆信息
	c.WxDologon(userInfo)
	c.Ctx.Redirect(302, "http://www.youdidi.vip/Portal/showdriverorder/")
    //c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html;charset=utf-8")
	//c.Ctx.WriteString("<img src=\""+userInfo.Headimgurl+"\" alt=\"test\" /><br />")
	//c.Ctx.WriteString(userInfo.Nickname+"\n")

    //refreshTokenUrl := "https://api.weixin.qq.com/sns/oauth2/refresh_token?" +
    //                  "appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN"
}

//根据微信的用户信息进行注册，并返回用户的cookie，
// 如果已经注册过OpenId已经存在则直接返回用户的cookie
func (c *WxLoginController) WxDologon(userInfo *UserInfo) error {
	var dbUser models.User
	var list []*models.User

	num, err := dbUser.GetUserInfoFormOpenId(userInfo.Openid, &list)
	if (err != nil) {
		logs.Error("OpenId %s get user info from db error：%s!", userInfo.Openid, err)
	}
	if (num == 0) {
		logs.Notice("OpenId %s not resgisted：%s!", userInfo.Openid, err)
		var newUser models.User
		newUser.Nickname = userInfo.Nickname
		newUser.Province = userInfo.Province
		newUser.OpenId = userInfo.Openid
		newUser.Sex = userInfo.Sex
		newUser.Unionid = userInfo.Unionid
		newUser.WechatImg = userInfo.Headimgurl
		_, err = newUser.Insert()
		if (err != nil){
			logs.Error("OpenId %s creat user failed!", userInfo.Openid)
		}

	} else if (num > 1) {
		logs.Error("OpenId %s multiple registration!", userInfo.Openid)
	}

	logs.Notice("%s", list[0])


	token, idStr, err := CacheUserLoginInfo(list[0])
	if (err != nil) {
		logs.Warn("Cache user login info failed!")
	}

	c.Ctx.SetSecureCookie("qyt", "qyt_id", idStr) //注入用户id，后续所有用户id都从cookie里获取
	c.Ctx.SetSecureCookie("qyt", "qyt_token", token)


	return nil

}


func GetWechatGZAccessToken(code string)  (*AccessToken, error) {
	appId := beego.AppConfig.String("weixin::AppId")
	appSecret := beego.AppConfig.String("weixin::AppSecret")

    u := url.Values{}
    u.Set("appid", appId)
    u.Set("secret", appSecret)
    u.Set("code", code)
    u.Set("grant_type", "authorization_code")

	accessTokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?" + u.Encode()
	logs.Notice("accessTokenUrl is :%s", accessTokenUrl)

    resp, err := http.Get(accessTokenUrl)
    if err != nil {
        return nil, fmt.Errorf("Get accessTokenUrl Failed:%s", err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read resp Body failed:%s", err)
    }

    var accessToken AccessToken
    err = json.Unmarshal([]byte(body), &accessToken)
    if err != nil {
        return nil, fmt.Errorf("JsonToMapDemo err: ", err)
    }

    if accessToken.Errcode != 0 {
        return nil, fmt.Errorf("Access_token API error:", accessToken.Errcode,
                   accessToken.Errmsg)
    }
    return &accessToken, nil
}

func WxGetUserInfo (accessToken *AccessToken) (*UserInfo, error) {
    u := url.Values{}
    u.Set("lang", "zh_CN")
    u.Set("openid", accessToken.Openid)
    u.Set("access_token", accessToken.Access_token)

    userInfoUrl := "https://api.weixin.qq.com/sns/userinfo?" + u.Encode()
	logs.Notice("userInfoUrl is :%s", userInfoUrl)

    resp, err := http.Get(userInfoUrl)
    if err != nil {
        return nil, fmt.Errorf("Get userInfoUrl Failed:%s", err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read resp Body failed:%s", err)
    }

    var userInfo UserInfo
    err = json.Unmarshal([]byte(body), &userInfo)
    if err != nil {
        return nil, fmt.Errorf("JsonToMapDemo err: ", err)
    }

    if userInfo.Errcode != 0 {
        return nil, fmt.Errorf("Userinfo API error:", userInfo.Errcode,
                   userInfo.Errmsg)
    }
    return &userInfo, nil
}

