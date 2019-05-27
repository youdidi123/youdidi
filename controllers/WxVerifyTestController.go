package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "sort"
    "strings"
)

type WxVerifyTestController struct {
    beego.Controller
}

func (c *WxVerifyTestController) Get() {
    signature := c.GetString("signature")
    timestamp := c.GetString("timestamp")
    echostr := c.GetString("echostr")
    nonce := c.GetString("nonce")
    token := "snow2019"

    List := []string{token, timestamp, nonce}
    sort.Strings(List)

    hashcode := Sha1(strings.Join(List, ""))

    if(hashcode == signature) {
        c.Ctx.WriteString(echostr)
    } else {
        c.Ctx.WriteString("")
    }

	logs.Notice("Hashcode is %s and signature is %s.", hashcode, signature)
}
