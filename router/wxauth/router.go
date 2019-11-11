package wxauth

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jacktea/wxproxy/bootstrap"
	"github.com/jacktea/wxproxy/cache"
	. "github.com/jacktea/wxproxy/common"
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

const templ = `
	<html>
	<head>
	</head>
	<body onload="document.getElementById('autoForm').submit();">
	<form id="autoForm" method="post" action="{{.Url}}">
		<input type="hidden" name="openid" value="{{.Openid}}">
		<input type="hidden" name="unionid" value="{{.Unionid}}">
		<input type="hidden" name="accessToken" value="{{.AccessToken}}">
	</form>
	</body>
	</html>
	`

var (
	TokenCache    cache.Cache = cache.NewCache(5000, 7000*time.Second)
	UserInfoCache cache.Cache = cache.NewCache(5000, 7000*time.Second)
	JsTicketCache cache.Cache = cache.NewCache(5000, 7000*time.Second)
	tpl                       = template.Must(template.New("wxauth_form").Parse(templ))
	log                       = golog.Default
)

type AuthAction struct {
	Svr api.ApiService `inject:""`
}

func (this *AuthAction) InitRouter(app *bootstrap.Bootstrapper) {
	party := app.Party(WXConf.HttpConf.ContextPath + "/wxauth")
	party.Get("/apply/{action_type}/{component_appid}/{appid}", this.WxApply)
	party.Get("/apply/{action_type}/{component_appid}/{appid}/do/{url}", this.WxDo)
	party.Any("/wx/userinfo", this.Userinfo)
	party.Get("/wx/jp/userinfo", this.Userinfo)
	party.Any("/wx/jsconfig/{component_appid}/{appid}", this.Jsticket)
	party.Get("/wx/jp/jsconfig/{component_appid}/{appid}", this.Jsticket)

	partyOld := app.Party(WXConf.HttpConf.ContextPath + "/wxoauth")
	partyOld.Any("/wx/userinfo", this.UserinfoOld)
	partyOld.Get("/wx/jp/userinfo", this.UserinfoOld)
	log.Info("init wxauth router...")
}

func (this *AuthAction) WxApply(c iris.Context) {
	//actionType := c.Params().Get("action_type")
	componentAppid := c.Params().Get("component_appid")
	appid := c.Params().Get("appid")
	rd := c.FormValue("rd")
	scope := c.FormValue("scope")
	if "" == scope {
		scope = "snsapi_base"
	}
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(rd))

	rd = fmt.Sprintf("%s://%s%s/do/%s", utils.Scheme(c.Request()), c.Host(), c.Path(), urlBase64)
	ret := fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&component_appid=%s#wechat_redirect",
		AUTHORIZE_URL, appid, url.QueryEscape(rd), scope, utils.Random(6), componentAppid)
	//return c.String(http.StatusOK,ret)
	log.Debug("微信服务地址:", ret)
	log.Debug("微信授权成功后的回调地址:", rd)
	c.Redirect(ret, http.StatusMovedPermanently)
	//return c.Redirect(http.StatusMovedPermanently, ret)
}

func (this *AuthAction) WxDo(c iris.Context) {
	actionType := c.Params().Get("action_type")
	componentAppid := c.Params().Get("component_appid")
	appid := c.Params().Get("appid")
	urlBase64 := c.Params().Get("url")
	code := c.FormValue("code")

	info, ok := this.Svr.CacheFindAppBaseInfo(componentAppid)
	if !ok {
		c.JSON(service.NewCommonResp(1000, "三方应用不存在!"))
		return
	}
	url := fmt.Sprintf("%s?appid=%s&code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s",
		ACCESS_TOKEN, appid, code, componentAppid, info.ComponentAccessToken)
	ret := new(api.AuthAccessToken)
	if err := utils.HttpGetJson(url, ret); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	if ret.IsSuccess() {
		TokenCache.Add(ret.Openid, ret) //缓存Token
		//如果是用户授权,则异步获取用户信息进行缓存
		log.Debug(ret.Scope, ",", strings.Contains(ret.Scope, "snsapi_userinfo"), ",", ret.AccessToken)
		if strings.Contains(ret.Scope, "snsapi_userinfo") {
			go userinfo(ret.AccessToken, ret.Openid)
		}

		data, err := base64.URLEncoding.DecodeString(urlBase64)
		if err != nil {
			c.JSON(service.NewServerErrorResp(err))
			return
		}
		rd := string(data)
		switch actionType {
		case "wx":
			if strings.Contains(rd, "?") {
				rd += "&"
			} else {
				rd += "?"
			}
			rd += fmt.Sprintf("openid=%s&unionid=%s&accessToken=%s",
				ret.Openid, ret.Unionid, ret.AccessToken)
			log.Debug("最终回调(GET)URL:", rd)
			c.Redirect(rd, http.StatusMovedPermanently)
			return
		case "wxex":
			buf := new(bytes.Buffer)
			tpl.Execute(buf, map[string]interface{}{
				"Url":         rd,
				"Openid":      ret.Openid,
				"Unionid":     ret.Unionid,
				"AccessToken": ret.AccessToken,
			})
			log.Debugf("最终回调(POST)URL:%s,openid=%s,unionid=%s,accessToken=%s \n", rd, ret.Openid, ret.Unionid, ret.AccessToken)
			c.Write(buf.Bytes())
			return
		default:
			c.JSON(service.NewCommonResp(2000, "un support action!"))
			return
		}
	} else {
		c.JSON(ret)
		return
	}
}

func (this *AuthAction) Userinfo(c iris.Context) {
	openid := c.FormValue("openid")
	callback := c.FormValue("callback")
	var (
		err error
		ret interface{}
	)
	defer func() {
		if err != nil {
			ret = &service.CommonResp{
				Errcode: 101,
				Errmsg:  err.Error(),
			}
			log.Error(err)
		}
		if ret == nil {
			ret = &service.CommonResp{
				Errcode: 101,
				Errmsg:  "请重新授权来获取用户信息",
			}
		}
		if "" == callback {
			c.JSON(ret)
		} else {
			c.JSONP(ret, context.JSONP{Callback: callback})
		}
	}()
	//从缓存获取
	if v, ok := UserInfoCache.Get(openid); ok {
		ret = v
		return
	}
	//如果token未失效，通过token重新获取用户信息
	if v, ok := TokenCache.Get(openid); ok {
		if t, ok := v.(api.AuthAccessToken); ok {
			if !strings.Contains(t.Scope, "snsapi_userinfo") {
				err = errors.New("请重新授权来获取用户信息")
				return
			}
			ret, err = userinfo(t.AccessToken, t.Openid)
		}
	}
}

type UserinfoOldResp struct {
	ErrorCode    int         `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data,omitempty"`
}

func (this *AuthAction) UserinfoOld(c iris.Context) {
	openid := c.FormValue("openid")
	callback := c.FormValue("callback")
	var (
		err error
		ret interface{}
	)
	defer func() {
		if err != nil {
			ret = &UserinfoOldResp{
				ErrorCode:    101,
				ErrorMessage: err.Error(),
			}
			log.Error(err)
		}
		if ret == nil {
			ret = &UserinfoOldResp{
				ErrorCode:    101,
				ErrorMessage: "请重新授权来获取用户信息",
			}
		} else {
			ret = &UserinfoOldResp{
				ErrorCode:    9000,
				ErrorMessage: "成功",
				Data:         ret,
			}
		}
		if "" == callback {
			c.JSON(ret)
		} else {
			c.JSONP(ret, context.JSONP{Callback: callback})
		}
	}()
	//从缓存获取
	if v, ok := UserInfoCache.Get(openid); ok {
		ret = v
		return
	}
	//如果token未失效，通过token重新获取用户信息
	if v, ok := TokenCache.Get(openid); ok {
		if t, ok := v.(api.AuthAccessToken); ok {
			if !strings.Contains(t.Scope, "snsapi_userinfo") {
				err = errors.New("请重新授权来获取用户信息")
				return
			}
			v, e := userinfo(t.AccessToken, t.Openid)
			if err != nil {
				err = e
			}
			ret = v
		}
	}
	return
}

func userinfo(access_token string, openid string) (*api.AuthUserinfo, error) {
	if v, ok := UserInfoCache.Get(openid); ok {
		if v, ok := v.(*api.AuthUserinfo); ok {
			return v, nil
		}
	}
	ret := new(api.AuthUserinfo)
	url := fmt.Sprintf("%s?access_token=%s&openid=%s&lang=%s",
		USERINFO_URL, access_token, openid, "zh_CN")
	if err := utils.HttpGetJson(url, ret); err != nil {
		log.Error(err)
		return nil, err
	}
	if ret.IsSuccess() {
		UserInfoCache.Add(openid, ret)
		return ret, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%d", ret.Errcode))
	}
}

func (this *AuthAction) Jsticket(c iris.Context) {
	var (
		ticket      *api.JsTicket
		accessToken string
		err         error
	)
	componentAppid := c.Params().Get("component_appid")
	appid := c.Params().Get("appid")
	url := c.FormValue("url")
	callback := c.FormValue("callback")

	if accessToken, err = this.Svr.GetAppAccessToken(componentAppid, appid); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	if ticket, err = jsticket(accessToken, appid); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	nonceStr := string(utils.Random(20))
	timestamp := time.Now().Unix()
	signature := _sha1(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket.Ticket, nonceStr, timestamp, url))
	ret := &map[string]interface{}{
		"signature": signature,
		"appId":     appid,
		"nonceStr":  nonceStr,
		"timestamp": timestamp,
		"url":       url,
	}
	if "" == callback {
		c.JSON(ret)
	} else {
		c.JSONP(ret, context.JSONP{Callback: callback})
	}
}

func _sha1(in string) string {
	h := sha1.New()
	io.WriteString(h, in)
	encode := h.Sum(nil)
	return hex.EncodeToString(encode)
}

func jsticket(access_token string, appid string) (*api.JsTicket, error) {
	if v, ok := JsTicketCache.Get(appid); ok {
		if v, ok := v.(*api.JsTicket); ok {
			return v, nil
		}
	}
	ret := new(api.JsTicket)
	url := fmt.Sprintf("%s?access_token=%s&type=jsapi",
		GZH_JSTICKET, access_token)
	if err := utils.HttpGetJson(url, ret); err != nil {
		return nil, err
	}
	if ret.IsSuccess() {
		JsTicketCache.Add(appid, ret)
		return ret, nil
	} else {
		return nil, fmt.Errorf("%d", ret.Errcode)
	}
}
