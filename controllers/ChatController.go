package controllers

import (
	"github.com/astaxie/beego"
)

type ChatController struct {
	beego.Controller
}

// @router /chat [GET]
func (this *ChatController) GoChat () {
	this.TplName = "chat.html"
}
