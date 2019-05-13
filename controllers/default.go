package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"youdidi/models"
	"youdidi/redisClient"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	var list []*models.User
	var dbUser models.User

	dbUser.Query().All(&list)

	var cacheClient redisClient.CacheClient
	cacheClient.GetConnet()
	cacheClient.SetKey("abc","123")

	test := cacheClient.GetKey("abc")

	fmt.Printf("key value is",test)


	c.Data["list"] = list
	c.TplName = "index.html"
}
