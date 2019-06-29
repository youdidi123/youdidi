package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
	"youdidi/models"
)

type ImgConfirmController struct {
	beego.Controller
}

var (
	dcStatus = []struct {
		Text string
	}{{"审核中"}, {"审核通过"}, {"审核失败"}}
)

// @router /Portal/driverconfirminput [GET]
func (c *ImgConfirmController) DriverConfirmInput() {
	userId, _ := c.Ctx.GetSecureCookie("qyt", "qyt_id")
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetUserOrder(userId, &dcInfo)
	c.Data["num"] = num
	if (num > 0) {
		c.Data["list"] = dcInfo[0]
		c.Data["dcStatus"] = dcStatus
	}

	c.Data["tabIndex"] = 3
	c.TplName = "driverConfirm.html"
}

// @router /Portal/dodriverconfirm [POST]
func (this *ImgConfirmController) DoDriverConfirm() {
	userId, _ := this.Ctx.GetSecureCookie("qyt", "qyt_id")
	userIdInt, _ := strconv.Atoi(userId)
	name := this.GetString("name")
	idNum := this.GetString("idNum")
	carType := this.GetString("carType")
	carNum := this.GetString("carNum")
	sfzSrc := this.GetString("sfzSrc")
	driverLiscen := this.GetString("driverLiscen")
	carLiscen := this.GetString("carLiscen")

	code := 0
	msg := ""
	currentTime := strconv.FormatInt(time.Now().Unix(),10)

	var confirmInfo models.Driver_confirm
	confirmInfo.User = &models.User{Id:userIdInt}
	confirmInfo.Status = 0
	confirmInfo.Time = currentTime
	confirmInfo.CarNum = carNum
	confirmInfo.CarType = carType
	confirmInfo.RealName = name
	confirmInfo.SfzNum = idNum
	var dbDc models.Driver_confirm
	var dcInfo []*models.Driver_confirm

	num := dbDc.GetUserOrder(userId, &dcInfo)

	if (num != 0) {
		if (dcInfo[0].Status == 0) {
			code = 5
			msg = "您有正在处理的申请，请勿重复提交"
			this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
			this.ServeJSON()
			return
		} else if (dcInfo[0].Status == 1) {
			code = 6
			msg = "您已是认证车主啦，无需再次验证"
			this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
			this.ServeJSON()
			return
		}
	}

	sfzFileName := "static/img/confirmImg/sfz_img_" + userId + "_" + currentTime
	driverLiscenFileName := "static/img/confirmImg/dl_img_" + userId + "_" + currentTime
	carLiscenFileName := "static/img/confirmImg/cl_img_" + userId + "_" + currentTime

	if (! StoreImg(sfzSrc, sfzFileName)) {
		code = 1
		msg = "资料上传异常，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	if (! StoreImg(driverLiscen, driverLiscenFileName)) {
		code = 2
		msg = "资料上传异常，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}
	if (! StoreImg(carLiscen, carLiscenFileName)) {
		code = 3
		msg = "资料上传异常，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	confirmInfo.CarLiceseImg = "/img/confirmImg/cl_img_" + userId + "_" + currentTime + ".png"
	confirmInfo.DriverLiceseImg = "/img/confirmImg/dl_img_" + userId + "_" + currentTime + ".png"
	confirmInfo.SfzImg = "/img/confirmImg/sfz_img_" + userId + "_" + currentTime + ".png"

	if (! confirmInfo.CreateDriverConfirm(userId, num)) {
		code = 4
		msg = "系统异常，请重试"
		this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"code":code, "msg":msg};
	this.ServeJSON()

}

func StoreImg (src string, fileName string) bool {
	re, _ := regexp.Compile("data:image/(png|jpg|jpeg);base64,");
	imgType := re.Find([]byte(src))

	//所有图片都后台存储位png的格式
	imgFileName := fileName+".png"
	//imgFileName := "static/driverConfirm.png"

	src = strings.Replace(src , string(imgType) , "" , 1)

	ddd, err1 := base64.StdEncoding.DecodeString(src) //成图片文件并把文件写入到buffer
	err2 := ioutil.WriteFile(imgFileName, ddd, 0666) //需要给图片文件一个可读权限

	if (err1 != nil || err2 != nil) {
		logs.Emergency("store user img fail imgName=%v", fileName)
		return false
	} else{
		return true
	}
}