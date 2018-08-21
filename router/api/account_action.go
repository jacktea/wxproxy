package api

import (
	"github.com/kataras/iris"
	"github.com/jacktea/wxproxy/service"
)

//创建带参二维码
func (a *ApiAction) CreateParamQrcode(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	identity := c.PostValue("identity")
	expire,err := c.PostValueInt64("expire")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	forever,err := c.PostValueBool("forever")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	resp,err := a.Svr.CreateParamQrcode(componentAppid,appid,identity,expire,forever)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(resp)
}

