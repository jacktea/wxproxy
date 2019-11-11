package api

import (
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/service"
	"github.com/kataras/iris/v12"
	"io/ioutil"
)

//发送客服消息
func (a *ApiAction) SendCustomMsg(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	d, _ := ioutil.ReadAll(c.Request().Body)
	content := string(d)
	err := a.Svr.SendCustomMsg(componentAppid, appid, content)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}

//发送模板消息
func (a *ApiAction) SendTplMsg(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	d, _ := ioutil.ReadAll(c.Request().Body)
	content := string(d)
	err := a.Svr.SendTplMsg(componentAppid, appid, content)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}

func (a *ApiAction) AddTemplage(c iris.Context) {
	a.postTransparentJson(c, ADD_TEMPLATE)
}
