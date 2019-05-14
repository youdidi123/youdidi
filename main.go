package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "youdidi/routers"
)



func main() {

	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/img", "static/img")

	logs.EnableFuncCallDepth(true)
	logs.Async()
	logs.Async(1e3)

	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"youdidi.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

	beego.Run()
}

