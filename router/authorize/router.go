package authorize

import (
	"github.com/jacktea/wxproxy/service/api"
	"fmt"
	"encoding/base64"
	"strings"
	"github.com/kataras/iris"
	"github.com/kataras/golog"
	"net/url"
	"bytes"
	"html/template"
	"net/http"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/utils"
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/bootstrap"
)

const templ = `
	<html>
	<head>
		<script>
			window.location.href = {{.Url}}
		</script>
	</head>
	</html>
	`

var log = golog.Default

var tpl = template.Must(template.New("escape").Parse(templ))

type AuthorizeAction struct {
	Svr api.ApiService	`inject:""`
}

func (this *AuthorizeAction) InitRouter(app *bootstrap.Bootstrapper) {
	party := app.Party(WXConf.HttpConf.ContextPath+"/authorize")
	party.Get("/{component_appid}/{url}",this.Authorize)
	party.Get("/{component_appid}/{url}/cb",this.AuthorizeCb)
	log.Info("init authorize router...")
}

func (this *AuthorizeAction) Authorize(c iris.Context) {
	componentAppid 	:= c.Params().Get("component_appid")
	appInfo, ok 	:= this.Svr.CacheFindAppBaseInfo(componentAppid)
	if !ok {
		c.JSON(service.NewCommonResp(1000,"三方应用不存在!"))
		return
	}
	rd := fmt.Sprintf("%s://%s%s/cb",utils.Scheme(c.Request()),c.Host(),c.Path())
	rd = url.QueryEscape(rd)
	url := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s",componentAppid,appInfo.PreAuthCode,rd)
	log.Debug("请求授权URL:",url)

	buf := new(bytes.Buffer)
	tpl.Execute(buf,map[string]interface{}{
		"Url":url,
	})
	c.Write(buf.Bytes())
}

func (this *AuthorizeAction) AuthorizeCb(c iris.Context) {
	componentAppid 	:= c.Params().Get("component_appid")
	urlBase64		:= c.Params().Get("url")
	authCode 		:= c.FormValue("auth_code")
	appid,err 		:= this.Svr.DoAuthorizerInfo(componentAppid,authCode)
	log.Debug(componentAppid,urlBase64,authCode)
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	data,err := base64.URLEncoding.DecodeString(urlBase64)
	if err != nil {
		log.Error("base64解码失败:",err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	rd				:= string(data)
	if strings.Contains(rd,"?") {
		rd += "&appid="+appid+"&componentAppid="+componentAppid
	}else {
		rd += "?appid="+appid+"&componentAppid="+componentAppid
	}
	c.Redirect(rd,http.StatusMovedPermanently)
	//return c.Redirect(http.StatusMovedPermanently,rd)
}
