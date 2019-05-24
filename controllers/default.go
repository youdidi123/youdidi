package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}
// @router /Portal/home [GET]
func (c *MainController) Get() {
	c.TplName = "index.html"
}
