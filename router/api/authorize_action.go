package api

import (
	"github.com/kataras/iris"
	"github.com/jacktea/wxproxy/service"
)

//更新托管公众号/小程序信息
func (a *ApiAction) UpdateAppInfo(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	if err := a.Svr.UpdateAuthorizerInfo(componentAppid, appid); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}
