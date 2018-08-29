package mini

import (
	. "github.com/jacktea/wxproxy/config"
	"github.com/kataras/golog"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/bootstrap"
	"github.com/kataras/iris"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/utils"
	"github.com/kataras/iris/context"
	"net/http"
	"github.com/jacktea/wxproxy/service/miniprogram"
)

var log = golog.Default

type MiniAction struct {
	Svr api.ApiService	`inject:""`
	MiniSvr miniprogram.MiniApiService `inject:""`
}

func (this *MiniAction) InitRouter(app *bootstrap.Bootstrapper) {
	r := app.Party(WXConf.HttpConf.ContextPath+"/mini")
	r.Post("/domain/modify/{componentAppid}/{appid}",this.ModifyDomain)
	r.Post("/domain/setwebviewdomain/{componentAppid}/{appid}",this.Setwebviewdomain)
	r.Get("/info/getaccountbasicinfo/{componentAppid}/{appid}",this.Getaccountbasicinfo)
	r.Post("/info/setnickname/{componentAppid}/{appid}",this.Setnickname)

	r.Post("/bind/bindtester/{componentAppid}/{appid}",this.BindTester)
	r.Post("/bind/unbindtester/{componentAppid}/{appid}",this.UnBindTester)
	r.Post("/bind/memberauth/{componentAppid}/{appid}",this.Memberauth)

	r.Post("/code/commit/{componentAppid}/{appid}",this.Commit)
	r.Get("/code/getqrcode/{componentAppid}/{appid}",this.GetQrCode)
	r.Get("/code/getqrcodeex/{componentAppid}/{appid}",this.GetQrCodeEx)
	r.Get("/code/prevqrcode/{componentAppid}/{appid}/{fName}",this.PrevQrCode)
	r.Post("/code/preview/{componentAppid}/{appid}",this.Preview)
	r.Get("/code/getcategroy/{componentAppid}/{appid}",this.GetCategory)
	r.Get("/code/getpage/{componentAppid}/{appid}",this.GetPage)
	r.Post("/code/submitaudit/{componentAppid}/{appid}",this.SubmitAudit)
	r.Post("/code/queryauditstatus/{componentAppid}/{appid}",this.QueryAuditStatus)
	r.Get("/code/querylastauditstatus/{componentAppid}/{appid}",this.QueryLastAuditStatus)
	r.Post("/code/dorelease/{componentAppid}/{appid}",this.DoRelease)
	r.Post("/code/changevisitstatus/{componentAppid}/{appid}",this.ChangeVisitStatus)
	r.Get("/code/revertcoderelease/{componentAppid}/{appid}",this.RevertCodeRelease)
	r.Post("/code/queryweappsupportversion/{componentAppid}/{appid}",this.QueryWeAppSupportVersion)
	r.Post("/code/setminweappsupportversion/{componentAppid}/{appid}",this.SetMinWeAppSupportVersion)

	r.Get("/codetplmgr/gettemplatedraftlist/{componentAppid}",this.Gettemplatedraftlist)
	r.Get("/codetplmgr/gettemplatelist/{componentAppid}",this.Gettemplatelist)
	r.Post("/codetplmgr/addtotemplate/{componentAppid}",this.Addtotemplate)
	r.Post("/codetplmgr/deletetemplate/{componentAppid}",this.Deletetemplate)

	r.Get("/user/login/{componentAppid}/{appid}",this.Jscode2Session)

	log.Info("init mini program router...")
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

//func (this *MiniAction) transparent(c iris.Context,urlPrefix string,contentType string) ([]byte,error) {
//	componentAppid := c.Params().Get("componentAppid")
//	appid := c.Params().Get("appid")
//	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
//	if err != nil {
//		c.JSON(service.NewServerErrorResp(err))
//		return nil,err
//	}
//	return transPost(urlPrefix,token,contentType,c.Request(),c)
//	//header,ret,err := utils.HttpPostRequestBody(urlPrefix+"?access_token="+token,contentType,c.Request())
//	//if err != nil {
//	//	c.JSON(service.NewServerErrorResp(err))
//	//	return nil,err
//	//}
//	//c.ContentType(header.Get(context.ContentTypeHeaderKey))
//	//_,err = c.Write(ret)
//	//return ret,err
//}

//func (this *MiniAction) gettransparent(c iris.Context,urlPrefix string) ([]byte,error) {
//	componentAppid := c.Params().Get("componentAppid")
//	appid := c.Params().Get("appid")
//	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
//	if err != nil {
//		c.JSON(service.NewServerErrorResp(err))
//		return nil,err
//	}
//	return transGet(urlPrefix,token,c)
//	//header,ret,err := utils.HttpGetProxy(urlPrefix+"?access_token="+token)
//	//if err != nil {
//	//	c.JSON(service.NewServerErrorResp(err))
//	//	return nil,err
//	//}
//	//c.ContentType(header.Get(context.ContentTypeHeaderKey))
//	//_,err = c.Write(ret)
//	//return ret,err
//}

func (this *MiniAction) postTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transPost(urlPrefix,token,context.ContentJSONHeaderValue,c.Request(),c)
	//return this.transparent(c,urlPrefix,context.ContentJSONHeaderValue)
}

func (this *MiniAction) getTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transGet(urlPrefix,token,c)
}

func (this *MiniAction) postCmpTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	token,err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transPost(urlPrefix,token,context.ContentJSONHeaderValue,c.Request(),c)
}

func (this *MiniAction) getCmpTransparentJson(c iris.Context,urlPrefix string) ([]byte,error) {
	componentAppid := c.Params().Get("componentAppid")
	token,err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil,err
	}
	return transGet(urlPrefix,token,c)
}
