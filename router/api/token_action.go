package api

import (
	"github.com/jacktea/wxproxy/service"
	"github.com/kataras/iris/v12"
)

//更新第三方应用Token
func (a *ApiAction) UpdateCmpToken(c iris.Context) {
	appid := c.Params().Get("appid")
	force, _ := c.PostValueBool("force")
	if _, err := a.Svr.UpdateAccessToken(appid, force); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}

//更新第三方应用预授权码
func (a *ApiAction) UpdateCmpCode(c iris.Context) {
	appid := c.Params().Get("appid")
	force, _ := c.PostValueBool("force")
	if err := a.Svr.UpdatePreAuthCode(appid, force); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}

//刷新托管公众号/小程序授权Token
func (a *ApiAction) RefreshAppAuthorizationToken(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	force, _ := c.PostValueBool("force")
	if _, err := a.Svr.RefreshAuthorizationToken(componentAppid, appid, force); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}
