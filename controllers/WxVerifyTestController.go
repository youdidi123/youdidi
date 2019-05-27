package controllers

import (
	"github.com/astaxie/beego"
    "sort"
    "strings"
)

type WxVerifyTestController struct {
    beego.Controller
}

func (c *WxVerifyTestController) Get() {
    signature := c.GetString("signature")
    timestamp := c.GetString("timestamp")
    nonce := c.GetString("nonce")
    token := "snow2019"

    List := []string{token, timestamp, nonce}
    sort.Strings(List)

    hashcode := Sha1(strings.Join(List, ""))

    c.Ctx.WriteString(hashcode+"\n")
    c.Ctx.WriteString(signature+"\n")
}
