package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"],
        beego.ControllerComments{
            Method: "GetAccountFlow",
            Router: `/Portal/accountflow`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"],
        beego.ControllerComments{
            Method: "DoWithdraw",
            Router: `/Portal/dowithdraw`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"],
        beego.ControllerComments{
            Method: "GetOpenId",
            Router: `/Portal/getOpenId`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"],
        beego.ControllerComments{
            Method: "Invest",
            Router: `/Portal/invest`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AccountFlowController"],
        beego.ControllerComments{
            Method: "Withdraw",
            Router: `/Portal/withdraw`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "AdminLogin",
            Router: `/AdminLogin/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "AdminDoLogin",
            Router: `/AdmindoLogin`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "Admin",
            Router: `/admin/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "ConfirmDriverDetail",
            Router: `/admin/confirmDriverDetail/:id`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "DriverConfirm",
            Router: `/admin/dconfirm`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"] = append(beego.GlobalControllerRouter["youdidi/controllers:AdminUserController"],
        beego.ControllerComments{
            Method: "UserWithdrew",
            Router: `/admin/userwithdrew`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ChatController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ChatController"],
        beego.ControllerComments{
            Method: "GoChat",
            Router: `/Portal/chat/:pid/:oid/:type`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ChatController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ChatController"],
        beego.ControllerComments{
            Method: "RefreshMsg",
            Router: `/Portal/refreshMsg`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ChatController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ChatController"],
        beego.ControllerComments{
            Method: "SetMsg",
            Router: `/Portal/setMsg`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"],
        beego.ControllerComments{
            Method: "DoDriverConfirm",
            Router: `/Portal/dodriverconfirm`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"] = append(beego.GlobalControllerRouter["youdidi/controllers:ImgConfirmController"],
        beego.ControllerComments{
            Method: "DriverConfirmInput",
            Router: `/Portal/driverconfirminput`,
            AllowHTTPMethods: []string{"GET"},
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
            Method: "AgreeRequest",
            Router: `/Portal/agreerequest`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "CreateOrder",
            Router: `/Portal/createorder`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "CreateOrderForce",
            Router: `/Portal/createorderforce`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DoCreateOrder",
            Router: `/Portal/docreateorder`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DoRecommand",
            Router: `/Portal/dorecommand`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DoRequire",
            Router: `/Portal/dorequire`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DriverCancle",
            Router: `/Portal/drivercancle`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DriverOrderDetail",
            Router: `/Portal/driverorderdetail/:oid`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DriverGetEnd",
            Router: `/Portal/getend`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DriverGetStart",
            Router: `/Portal/getstart`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PassengerCancle",
            Router: `/Portal/passengercancle`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PassengerConfirm",
            Router: `/Portal/passengerconfirm`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PassengerOrderDetail",
            Router: `/Portal/passengerorderdetail/:odid`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "Recommand",
            Router: `/Portal/recommand/:odid/:uType`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "RefuseRequest",
            Router: `/Portal/refuserequest`,
            AllowHTTPMethods: []string{"POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "SearchInput",
            Router: `/Portal/searchinput`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "SearchOrder",
            Router: `/Portal/searchorder`,
            AllowHTTPMethods: []string{"GET","POST"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "ShowDriverOrder",
            Router: `/Portal/showdriverorder/`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:OrderController"] = append(beego.GlobalControllerRouter["youdidi/controllers:OrderController"],
        beego.ControllerComments{
            Method: "ShowPassengerOrder",
            Router: `/Portal/showpassengerorder`,
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
            Method: "Dologon",
            Router: `/Dologon/`,
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
            Method: "AboutUs",
            Router: `/Portal/aboutus`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "Disclaimer",
            Router: `/Portal/disclaimer`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "Help",
            Router: `/Portal/help`,
            AllowHTTPMethods: []string{"GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"] = append(beego.GlobalControllerRouter["youdidi/controllers:UserCenterController"],
        beego.ControllerComments{
            Method: "UserInfo",
            Router: `/Portal/userinfo`,
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
            Method: "UserInfoCheck",
            Router: `/UserInfoCheck/`,
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

    beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"] = append(beego.GlobalControllerRouter["youdidi/controllers:WxLoginController"],
        beego.ControllerComments{
            Method: "Wxtest",
            Router: `/Wxtest/`,
            AllowHTTPMethods: []string{"POST","GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["youdidi/controllers:WxPayController"] = append(beego.GlobalControllerRouter["youdidi/controllers:WxPayController"],
        beego.ControllerComments{
            Method: "WxInvest",
            Router: `/Portal/WxInvest/`,
            AllowHTTPMethods: []string{"POST","GET"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
