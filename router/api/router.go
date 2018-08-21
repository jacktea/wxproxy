package api

import (
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/bootstrap"
	"github.com/kataras/golog"
)

var log = golog.Default

type ApiAction struct {
	Svr api.ApiService	`inject:""`
}

func (this *ApiAction) InitRouter(app *bootstrap.Bootstrapper) {
	r := app.Party(WXConf.HttpConf.ContextPath+"/api")
	r.Post("/token/upcmptoken/{appid:string}",this.UpdateCmpToken)
	r.Post("/token/upcmpcode/{appid:string}",this.UpdateCmpCode)
	r.Post("/token/upauthtoken/{componentAppid}/{appid}",this.RefreshAppAuthorizationToken)

	//更新托管账号基本信息
	r.Post("/account/upappinfo/{componentAppid}/{appid}",this.UpdateAppInfo)
	//更新托管账号回调地址
	r.Post("/account/upappnotifyurl/{componentAppid}/{appid}",this.UpdateAuthorizationAppNotifyUrl)
	//获取托管公众号的信息
	r.Get("/account/getauthappinfo/{componentAppid}/{appid}",this.GetAuthAppInfo)
	//创建带参关注二维码
	r.Post("/account/createqrcode/{componentAppid}/{appid}",this.CreateParamQrcode)


	//获取用户基本信息
	r.Post("/user/info/{componentAppid}/{appid}",this.GetUserBaseInfo)
	//发送客服消息
	r.Post("/msg/custommsg/{componentAppid}/{appid}",this.SendCustomMsg)
	//发送模板消息
	r.Post("/msg/tplmsg/{componentAppid}/{appid}",this.SendTplMsg)

	log.Info("init api router")

}

