package mini

import (
	"github.com/kataras/iris"
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/utils"
	"fmt"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/service/miniprogram"
)

// 微信登录
func (this *MiniAction) Jscode2Session(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	jsCode := c.FormValue("js_code")
	token,err := this.Svr.GetComponentAppAccessToken(componentAppid)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	url := fmt.Sprintf("%s?appid=%s&js_code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s",
		USER_MINI_LOGIN,appid,jsCode,componentAppid,token)
	resp := new(miniprogram.ClientLoginResp)
	err = utils.HttpGetJson(url,resp)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(resp)
}
