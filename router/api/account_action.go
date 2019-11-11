package api

import (
	"github.com/jacktea/wxproxy/service"
	"github.com/kataras/iris/v12"
	"strconv"
)

//创建带参二维码
func (a *ApiAction) CreateParamQrcode(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	identity := c.PostValue("identity")
	expire, err := c.PostValueInt64("expire")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	forever, err := c.PostValueBool("forever")
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	resp, err := a.Svr.CreateParamQrcode(componentAppid, appid, identity, expire, forever)
	if err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(resp)
}

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

//更新应用的通知URL
func (a *ApiAction) UpdateAuthorizationAppNotifyUrl(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	notifyUrl := c.FormValue("notify_url")
	modeStr := c.FormValue("mode")
	debugNotifyUrl := c.FormValue("debug_notify_url")

	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		log.Error(err, "解析mode失败，设为默认值1")
		mode = 1
	}

	if err := a.Svr.UpdateAuthorizationAppNotifyUrl(componentAppid, appid, notifyUrl, mode, debugNotifyUrl); err != nil {
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	c.JSON(service.SUCCESS_RESP)
}

//获取托管应用的应用信息
func (a *ApiAction) GetAuthAppInfo(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	appid := c.Params().Get("appid")
	info, ok := a.Svr.CacheFindAuthorizationInfo(appid)
	if !ok {
		c.JSON(service.NewCommonResp(1000, "未找到相关内容"))
		return
	}
	m := make(map[string]interface{}, 0)
	m["errcode"] = 0
	m["appid"] = info.Appid
	m["nickName"] = info.NickName
	m["headImg"] = info.HeadImg
	m["serviceTypeInfo"] = info.ServiceTypeInfo
	m["verifyTypeInfo"] = info.VerifyTypeInfo
	m["userName"] = info.UserName
	m["principalName"] = info.PrincipalName
	m["alias"] = info.Alias
	m["businessInfo"] = info.BusinessInfo
	m["qrcodeUrl"] = info.QrcodeUrl
	m["signature"] = info.Signature
	m["miniprogram"] = info.Miniprogram
	aaInfo, ok := a.Svr.CacheFindAuthorizationAccessInfo(componentAppid, appid)
	if ok {
		m["accessToken"] = aaInfo.AccessToken
		m["accessTokenExpire"] = aaInfo.AccessTokenExpire.Format("2006-01-02 15:04:05") //访问token过期时间
		m["authorizationStatus"] = aaInfo.Status                                        //授权状态
	}
	aInfo, ok := a.Svr.CacheFindAuthorizationAppInfo(componentAppid, appid)
	if ok {
		m["mode"] = aInfo.Mode                     //模式	1普通 2调试
		m["notifyUrl"] = aInfo.NotifyUrl           //通知URL
		m["debugNotifyUrl"] = aInfo.DebugNotifyUrl //调试通知URL
	}
	c.JSON(m)
}
