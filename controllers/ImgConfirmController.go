package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"regexp"
	"strings"
)

type ImgConfirmController struct {
	beego.Controller
}

// @router /Test [GET,POST]
func (c *ImgConfirmController) DriverConfirmInput() {
	c.TplName = "driverConfirm.html"
}

// @router /Imgloader [POST]
func (c *ImgConfirmController) Imgloader() {
	imgId := c.GetString("imgId")
	content := c.GetString("base64")

	logs.Notice("user id xxx submit driverConf id is %v" , imgId)

	re, _ := regexp.Compile("data:image/(png|jpg|jpeg);base64,");
	imgType := re.Find([]byte(content))

	imgFileName := "static/driverConfirm_"+imgId+".png"
	//imgFileName := "static/driverConfirm.png"

	content = strings.Replace(content , string(imgType) , "" , 1)
	
	ddd, _ := base64.StdEncoding.DecodeString(content) //成图片文件并把文件写入到buffer
	ioutil.WriteFile(imgFileName, ddd, 0666)

	c.Data["json"]=map[string]interface{}{"code":"1"};
	c.ServeJSON();
}