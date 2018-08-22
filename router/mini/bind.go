package mini

import (
	"github.com/kataras/iris"
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/service"
	"github.com/kataras/iris/context"
	"github.com/jacktea/wxproxy/utils"
)

// 绑定微信用户为小程序体验者
func (this *MiniAction) BindTester(c iris.Context)  {
	this.postTransparentJson(c,BIND_MINI_TESTER)
}

// 解除绑定小程序的体验者
func (this *MiniAction) UnBindTester(c iris.Context)  {
	this.postTransparentJson(c,UNBIND_MINI_TESTER)
}

// 获取体验者列表
func (this *MiniAction) Memberauth(c iris.Context)  {
	body := `{"action":"get_experiencer"}`
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	header,ret,err := utils.HttpPostProxy(GET_MINI_TESTER+"?access_token="+token,context.ContentJSONHeaderValue,[]byte(body))
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.ContentType(header.Get(context.ContentTypeHeaderKey))
	_,err = c.Write(ret)
}