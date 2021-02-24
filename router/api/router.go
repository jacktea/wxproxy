package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/jacktea/wxproxy/bootstrap"
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

var log = golog.Default

type ApiAction struct {
	Svr api.ApiService `inject:""`
}

func (this *ApiAction) InitRouter(app *bootstrap.Bootstrapper) {
	r := app.Party(WXConf.HttpConf.ContextPath + "/api")
	r.Post("/token/upcmptoken/{appid:string}", this.UpdateCmpToken)
	r.Post("/token/upcmpcode/{appid:string}", this.UpdateCmpCode)
	r.Post("/token/upauthtoken/{componentAppid}/{appid}", this.RefreshAppAuthorizationToken)

	//更新托管账号基本信息
	r.Post("/account/upappinfo/{componentAppid}/{appid}", this.UpdateAppInfo)
	//更新托管账号回调地址
	r.Post("/account/upappnotifyurl/{componentAppid}/{appid}", this.UpdateAuthorizationAppNotifyUrl)
	//获取托管公众号的信息
	r.Get("/account/getauthappinfo/{componentAppid}/{appid}", this.GetAuthAppInfo)
	//创建带参关注二维码
	r.Post("/account/createqrcode/{componentAppid}/{appid}", this.CreateParamQrcode)

	//获取用户基本信息
	r.Post("/user/info/{componentAppid}/{appid}", this.GetUserBaseInfo)
	//发送客服消息
	r.Post("/msg/custommsg/{componentAppid}/{appid}", this.SendCustomMsg)
	//发送模板消息
	r.Post("/msg/tplmsg/{componentAppid}/{appid}", this.SendTplMsg)
	//添加模板
	r.Post("/msg/addmsgtpl/{componentAppid}/{appid}", this.AddTemplage)

	// 代理微信接口
	r.Any("/proxyByDb/{componentAppId}/{apiId}", this.proxyByDb)
	r.Any("/proxyByDb/{componentAppId}/{appId}/{apiId}", this.proxyByDb)

	// 代理微信接口
	r.Any("/proxyByDbEx/{componentAppId}/{apiId}", this.proxyByDbEx)
	r.Any("/proxyByDbEx/{componentAppId}/{appId}/{apiId}", this.proxyByDbEx)

	//r.Any("/textProxy", this.textProxy)

	log.Info("init api router")

}

func transPost(urlPrefix, token, contentType string, req *http.Request, c iris.Context) ([]byte, error) {
	header, ret, err := utils.HttpPostRequestBody(urlPrefix+"?access_token="+token, contentType, req)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_, err = c.Write(ret)
	return ret, err
}

func transPostEx(urlPrefix, token string, req *http.Request, c iris.Context) ([]byte, error) {
	header, ret, err := utils.HttpPostRequestBody(urlPrefix+"?access_token="+token, c.GetHeader(context.ContentTypeHeaderKey), req)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_, err = c.Write(ret)
	return ret, err
}

func transGet(urlPrefix, token string, c iris.Context) ([]byte, error) {
	header, ret, err := utils.HttpGetProxy(urlPrefix + "?access_token=" + token)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_, err = c.Write(ret)
	return ret, err
}

func (this *ApiAction) postTransparentJson(c iris.Context, urlPrefix string) ([]byte, error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token, err := this.Svr.GetAppAccessToken(componentAppid, appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	return transPost(urlPrefix, token, context.ContentJSONHeaderValue, c.Request(), c)
}

func (this *ApiAction) getTransparentJson(c iris.Context, urlPrefix string) ([]byte, error) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token, err := this.Svr.GetAppAccessToken(componentAppid, appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	return transGet(urlPrefix, token, c)
}

func (this *ApiAction) postCmpTransparentJson(c iris.Context, urlPrefix string) ([]byte, error) {
	componentAppid := c.Params().Get("componentAppid")
	token, err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	return transPost(urlPrefix, token, context.ContentJSONHeaderValue, c.Request(), c)
}

func (this *ApiAction) getCmpTransparentJson(c iris.Context, urlPrefix string) ([]byte, error) {
	componentAppid := c.Params().Get("componentAppid")
	token, err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return nil, err
	}
	return transGet(urlPrefix, token, c)
}

func (this *ApiAction) proxyByDb(c iris.Context) {
	apiId, err := c.Params().GetInt64("apiId")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	apiInfo, ok := this.Svr.FindProxyApi(apiId)
	if ok {
		//c.JSON(apiInfo)
		var token string
		var err error
		componentAppid := c.Params().Get("componentAppId")
		switch apiInfo.Type {
		case "1", "2": // 公众号，小程序
			appid := c.Params().Get("appId")
			token, err = this.Svr.GetAppAccessToken(componentAppid, appid)
			if err != nil {
				c.JSON(service.NewServerErrorResp(err))
				return
			}
		case "3": // 三方应用
			token, err = this.Svr.GetComponentAppAccessToken(componentAppid)
			if err != nil {
				c.JSON(service.NewServerErrorResp(err))
				return
			}
		}
		switch apiInfo.Action {
		case "1": // 公众号，小程序
			transGet(apiInfo.Url, token, c)
		case "2":
			transPostEx(apiInfo.Url, token, c.Request(), c)
		default:
			c.JSON(service.NewCommonResp(1000, "接口请求类型有误"))
		}
	} else {
		c.JSON(service.NewCommonResp(1000, "接口不存在"))
	}
}

func (this *ApiAction) proxyByDbEx(c iris.Context) {
	apiId, err := c.Params().GetInt64("apiId")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	apiInfo, ok := this.Svr.FindProxyApi(apiId)
	if ok {
		var token string
		var err error
		componentAppid := c.Params().Get("componentAppId")
		switch apiInfo.Type {
		case "1", "2": // 公众号，小程序
			appid := c.Params().Get("appId")
			token, err = this.Svr.GetAppAccessToken(componentAppid, appid)
			if err != nil {
				c.JSON(service.NewServerErrorResp(err))
				return
			}
		case "3": // 三方应用
			token, err = this.Svr.GetComponentAppAccessToken(componentAppid)
			if err != nil {
				c.JSON(service.NewServerErrorResp(err))
				return
			}
		}
		apiUrl := apiInfo.Url
		if strings.Contains(apiInfo.Url, "?") {
			apiUrl += "&access_token=" + token
		} else {
			apiUrl += "?access_token=" + token
		}
		remote, err := url.Parse(apiUrl)
		if err != nil {
			c.JSON(service.NewServerErrorResp(err))
			return
		}
		proxy := NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(c.ResponseWriter(), c.Request())
	}
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawPath = target.RawPath
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		fmt.Println(req.URL.Scheme, req.URL.Host, req.URL.Path, req.URL.RawPath, req.URL.RawQuery)
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

//func (this *ApiAction) textProxy(c iris.Context) {
//	apiUrl := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=42_2ZLfMqWbEZLAlhCpBJjFU5vu727TwYp_P41zEElpaT19vC8BDlfKtiMZVyl2tyh5_E653G_KS8c4bS4EX1GEDKGfrmfJwIoYVRIsqh5TIhMmUQHL7WDwplHMvmSEAv3ZlFWawqEyAo-PrDojUSTjAKDMKL"
//	//apiUrl := "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template"
//
//	remote, err := url.Parse(apiUrl)
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(service.NewServerErrorResp(err))
//		return
//	}
//	proxy := NewSingleHostReverseProxy(remote)
//	proxy.ServeHTTP(c.ResponseWriter(), c.Request())
//}
