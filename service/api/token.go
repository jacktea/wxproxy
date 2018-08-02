package api

import (
	"github.com/jacktea/wxproxy/model"
	. "github.com/jacktea/wxproxy/common"
	"fmt"
	"github.com/jacktea/wxproxy/etcd"
	"github.com/jacktea/wxproxy/utils"
)

//更新三方应用票据
func (this *ApiServiceImpl) UpdateTicket(appId string, ticket string) error {
	var err error
	info, ok := this.CacheFindAppBaseInfo(appId)
	if !ok {
		return NO_DATA
	}
	delete(appBaseInfos, appId)
	this.Repo.UpdateTicket(appId, ticket)

	if info.TokenIsExpired() {
		_, err = this.UpdateAccessToken(appId, false)
	}

	if info.AuthCodeIsExpired() {
		err = this.UpdatePreAuthCode(appId, false)
	}

	return err
}

//更新三方应用访问Token
func (this *ApiServiceImpl) UpdateAccessToken(appId string, force bool) (info *model.AppBaseInfo, err error) {
	info, ok := this.CacheFindAppBaseInfo(appId)
	if !ok {
		return info, NO_DATA
	}
	if info.TokenIsExpired() || force { //token过期或强制更新token
		log.Info("update token")
		var (
			uri    = API_COMPONENT_TOKEN
			params = map[string]interface{}{
				"component_appid":         info.AppId,
				"component_appsecret":     info.AppSecret,
				"component_verify_ticket": info.ComponentVerifyTicket,
			}
			ret compTokenResp
		)
		err = executeApi(uri, params, &ret)
		if err == nil {
			defer delete(appBaseInfos, appId)
			this.Repo.UpdateAccessToken(appId, ret.ComponentAccessToken, ret.ExpiresIn)
			info.ComponentAccessToken = ret.ComponentAccessToken
			info.ComponentAccessTokenExpire = utils.ParseExpire(ret.ExpiresIn)
			etcd.Put(COMPONENT_ACCESS_TOKEN_PREFIX+appId, ret.ComponentAccessToken, ret.ExpiresIn-30)
		}
	}
	return
}

//更新第三方应用预授权码
func (this *ApiServiceImpl) UpdatePreAuthCode(appId string, force bool) (err error) {
	info, ok := this.CacheFindAppBaseInfo(appId)
	if !ok {
		return NO_DATA
	}
	if info.AuthCodeIsExpired() || force {
		log.Info("update preAuthCode")
		var (
			uri    = API_CREATE_PREAUTHCODE + "?component_access_token=" + info.ComponentAccessToken
			params = map[string]interface{}{
				"component_appid": info.AppId,
			}
			ret preAuthCodeResp
		)
		err = executeApi(uri, params, &ret)
		if err == nil {
			defer delete(appBaseInfos, appId)
			this.Repo.UpdatePreAuthCode(appId, ret.PreAuthCode, ret.ExpiresIn)
			etcd.Put(PRE_AUTH_CODE_PREFIX+appId, ret.PreAuthCode, ret.ExpiresIn-30)
		}
	}
	return
}

//刷新托管公众号(小程序)访问Token
func (this *ApiServiceImpl) RefreshAuthorizationToken(componentAppid string, appid string, force bool) (dbAai *model.AuthorizationAccessInfo, err error) {
	var accessToken string
	dbAai, b := this.CacheFindAuthorizationAccessInfo(componentAppid, appid)
	if !b {
		log.Error(NO_AUTH_ACCESS_TOKEN)
		return nil, NO_AUTH_ACCESS_TOKEN
	}
	if !(dbAai.TokenIsExpired() || force) {
		return
		//return dbAai,nil
	}
	info, err := this.UpdateAccessToken(componentAppid, false)
	if err != nil {
		log.Error(err)
		return
	}
	var (
		uri    = API_AUTHORIZER_TOKEN + "?component_access_token=" + info.ComponentAccessToken
		params = map[string]interface{}{
			"component_appid":          componentAppid,
			"authorizer_appid":         appid,
			"authorizer_refresh_token": dbAai.RefreshToken,
		}
		resp refreshAccessTokenResp
	)
	err = executeApi(uri, params, &resp)
	if err == nil {
		accessToken = resp.AuthorizerAccessToken
		accessTokenExpire := utils.ParseExpire(resp.ExpiresIn)
		dbAai = &model.AuthorizationAccessInfo{
			ComponentAppid:    componentAppid,
			Appid:             appid,
			AccessToken:       accessToken,
			AccessTokenExpire: accessTokenExpire,
			RefreshToken:      resp.AuthorizerRefreshToken,
		}
		//清除缓存
		defer delete(authAccessInfos, fmt.Sprintf("%s_%s", componentAppid, appid))
		//更新数据库
		this.Repo.MergeAuthorizationAccessInfo(dbAai)
		//添加凭证监控
		log.Debug("update etcd ", REFRESH_TOKEN_PREFIX+componentAppid+"/"+appid)
		etcd.Put(REFRESH_TOKEN_PREFIX+componentAppid+"/"+appid, resp.AuthorizerAccessToken, resp.ExpiresIn-30)
	} else {
		log.Error("RefreshAuthorizationToken", err)
	}
	return
}

func (this *ApiServiceImpl) GetAppAccessToken(componentAppid string, appid string) (string, error) {
	info, err := this.RefreshAuthorizationToken(componentAppid, appid, false)
	if err != nil {
		//log.Error(err)
		return "", err
	}
	accessToken := info.AccessToken
	return accessToken, err
}
