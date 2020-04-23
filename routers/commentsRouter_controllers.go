package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["xiaozhang/controllers:DataController"] = append(beego.GlobalControllerRouter["xiaozhang/controllers:DataController"],
        beego.ControllerComments{
            Method: "Count",
            Router: `/count`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xiaozhang/controllers:DataController"] = append(beego.GlobalControllerRouter["xiaozhang/controllers:DataController"],
        beego.ControllerComments{
            Method: "Getone",
            Router: `/getone`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xiaozhang/controllers:DataController"] = append(beego.GlobalControllerRouter["xiaozhang/controllers:DataController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/post`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
