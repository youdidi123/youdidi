package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"],
        beego.ControllerComments{
            Method: "Imgloader",
            Router: `/Imgloader`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"],
        beego.ControllerComments{
            Method: "DriverConfirmInput",
            Router: `/Test`,
            AllowHTTPMethods: []string{"GET","POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:MainController"] = append(beego.GlobalControllerRouter["youdidi/controllers:MainController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/Portal/home`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "ShowDriverOrder",
            Router: `/Portal/showDriverorder/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "Dologin",
            Router: `/Dologin/`,
            AllowHTTPMethods: []string{"POST","GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/Login/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "GetVerCode",
            Router: `/Ver/getvercode`,
            AllowHTTPMethods: []string{"GET","POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "PhoneVer",
            Router: `/Ver/phonever`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "VerPhone",
            Router: `/Ver/verPhone`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"] = append(beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"],
        beego.ControllerComments{
            Method: "GetWechatGZAccessToken",
            Router: `/GetWechatGZAccessToken/`,
            AllowHTTPMethods: []string{"POST","GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"] = append(beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"],
        beego.ControllerComments{
            Method: "WxLogin",
            Router: `/WxLogin/`,
            AllowHTTPMethods: []string{"POST","GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
