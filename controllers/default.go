package controllers

import (
	"github.com/astaxie/beego"
	"youdidi/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	var list []*models.User
	var dbUser models.User

	dbUser.Query().All(&list)

	c.Data["list"] = list
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}
