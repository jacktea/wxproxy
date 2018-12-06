package api

import (
	"github.com/jacktea/wxproxy/bootstrap"
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"net/http"
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
	//添加模板
	r.Post("/msg/addmsgtpl/{componentAppid}/{appid}",this.AddTemplage)

	log.Info("init api router")

}

func transPost(urlPrefix,token,contentType string,req *http.Request,c iris.Context) ([]byte,error) {
	header,ret,err := utils.HttpPostRequestBody(urlPrefix+"?access_token="+token,contentType,req)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_,err = c.Write(ret)
	return ret,err
}

func transGet(urlPrefix,token string,c iris.Context) ([]byte,error) {
	header,ret,err := utils.HttpGetProxy(urlPrefix+"?access_token="+token)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_,err = c.Write(ret)
	return ret,err
}

func (this *ApiAction) postTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transPost(urlPrefix,token,context.ContentJSONHeaderValue,c.Request(),c)
}

func (this *ApiAction) getTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transGet(urlPrefix,token,c)
}

func (this *ApiAction) postCmpTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	token,err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transPost(urlPrefix,token,context.ContentJSONHeaderValue,c.Request(),c)
}

func (this *ApiAction) getCmpTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	token,err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transGet(urlPrefix,token,c)
}

