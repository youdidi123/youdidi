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

// @router /Portal/driverconfirminput [GET]
func (c *ImgConfirmController) DriverConfirmInput() {
	c.Data["tabIndex"] = 2
	c.TplName = "driverConfirm.html"
}

// @router /Imgloader [POST]
func (c *ImgConfirmController) Imgloader() {
	imgId := c.GetString("imgId")
	content := c.GetString("base64")

	logs.Notice("user id xxx submit driverConf id is %v" , imgId)

	//获取图片数据，只能处理png，jpg，jpeg
	re, _ := regexp.Compile("data:image/(png|jpg|jpeg);base64,");
	imgType := re.Find([]byte(content))

	//所有图片都后台存储位png的格式
	imgFileName := "static/driverConfirm_"+imgId+".png"
	//imgFileName := "static/driverConfirm.png"

	content = strings.Replace(content , string(imgType) , "" , 1)
	
	ddd, _ := base64.StdEncoding.DecodeString(content) //成图片文件并把文件写入到buffer
	ioutil.WriteFile(imgFileName, ddd, 0666) //需要给图片文件一个可读权限

	c.Data["json"]=map[string]interface{}{"code":"1"};
	c.ServeJSON();
}