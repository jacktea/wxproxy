package api

import (
	"github.com/jacktea/wxproxy/service"
	"github.com/kataras/iris/v12"
)

//获取用户基本信息
func (a *ApiAction) GetUserBaseInfo(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	openid := c.PostValue("openid")
	info, err := a.Svr.GetUserBaseInfo(componentAppid, appid, openid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(info)
}
