package api

import (
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/model"
	"github.com/jacktea/wxproxy/etcd"
	"fmt"
	"strconv"
	"github.com/jacktea/wxproxy/utils"
)

func (s *ApiServiceImpl) doAuthorize(componentAppid string, authorizationCode string) (string, error) {
	var appid string
	//获取第三方应用信息
	info, err := s.UpdateAccessToken(componentAppid, false)
	if err != nil {
		return appid, err
	}
	var (
		uri    = API_QUERY_QUTH + "?component_access_token=" + info.ComponentAccessToken
		params = map[string]interface{}{
			"component_appid":    componentAppid,
			"authorization_code": authorizationCode,
		}
		resp authorizerInfoResp
	)
	err = executeApi(uri, params, &resp)
	log.Debug("执行授权!", err)
	if err == nil {
		aInfo := resp.AuthorizationInfo
		appid = aInfo.AuthorizerAppid
		accessTokenExpire := utils.ParseExpire(aInfo.ExpiresIn)
		aai := &model.AuthorizationAccessInfo{
			ComponentAppid:    componentAppid,
			Appid:             appid,
			AccessToken:       aInfo.AuthorizerAccessToken,
			AccessTokenExpire: accessTokenExpire,
			RefreshToken:      aInfo.AuthorizerRefreshToken,
			FuncInfo:          toStr(aInfo.FuncInfo),
			AuthorizationCode: authorizationCode,
			Status:            "1",
		}
		//更新数据库
		s.Repo.MergeAuthorizationAccessInfo(aai)
		//添加凭证监控
		etcd.Put(REFRESH_TOKEN_PREFIX+componentAppid+"/"+appid, aInfo.AuthorizerAccessToken, aInfo.ExpiresIn-30)
		s.UpdateAuthorizerInfo(componentAppid, appid)
		defer s.UpdatePreAuthCode(componentAppid, true)
	}
	return appid, err
}

//更新托管公众号(小程序)信息
func (s *ApiServiceImpl) UpdateAuthorizerInfo(componentAppid string, authorizerAppid string) (err error) {
	info, err := s.UpdateAccessToken(componentAppid, false)
	if err != nil {
		return err
	}
	var (
		uri    = API_GET_AUTHORIZER_INFO + "?component_access_token=" + info.ComponentAccessToken
		params = map[string]interface{}{
			"component_appid":  componentAppid,
			"authorizer_appid": authorizerAppid,
		}
		resp authorizerInfoResp
	)
	err = executeApi(uri, params, &resp)
	if err == nil {
		aInfo := resp.AuthorizerInfo
		miniProgram := 0
		if aInfo.MiniProgramInfo != nil {//当前为小程序
			miniProgram = 1
		}
		dbInfo := &model.AuthorizationInfo{
			Appid:           authorizerAppid,
			NickName:        aInfo.NickName,
			HeadImg:         aInfo.HeadImg,
			ServiceTypeInfo: strconv.Itoa(aInfo.ServiceTypeInfo["id"]),
			VerifyTypeInfo:  strconv.Itoa(aInfo.VerifyTypeInfo["id"]),
			UserName:        aInfo.UserName,
			PrincipalName:   aInfo.PrincipalName,
			Alias:           aInfo.Alias,
			BusinessInfo:    toJsonStr(aInfo.BusinessInfo),
			QrcodeUrl:       aInfo.QrcodeUrl,
			Signature:       aInfo.Signature,
			Miniprogram:	 miniProgram,
		}
		defer delete(authInfos, authorizerAppid)
		s.Repo.MergeAuthorizationInfo(dbInfo)
	}
	return
}

func (s *ApiServiceImpl) UpdateAuthorizationAppNotifyUrl(componentAppid string,
	appid string,
	notifyUrl string,
	mode int,
	debugNotifyUrl string) error {
	item := model.AuthorizationAppInfo{
		ComponentAppid:componentAppid,
		Appid:appid,
		NotifyUrl:notifyUrl,
		Mode:mode,
		DebugNotifyUrl:debugNotifyUrl,
	}
	_,err := s.Repo.MergeAuthorizationAppInfo(&item)
	//清除缓存
	defer delete(authAppInfos,appid)
	return err
}

//时间通知时的应用授权
func (s *ApiServiceImpl) DoAuthorize(componentAppid string, authorizerAppid string, authorizationCode string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	asInfo, b := s.CacheFindAuthorizationAccessInfo(componentAppid, authorizerAppid)
	if b && (asInfo.Status == "1" && authorizationCode == asInfo.AuthorizationCode) { //已经授权不进行重复授权
		return asInfo.Appid, nil
	}
	return s.doAuthorize(componentAppid, authorizationCode)
}

//网页回调时应用授权
func (s *ApiServiceImpl) DoAuthorizerInfo(componentAppid string, authorizationCode string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	asInfo, b := s.Repo.FindAuthAccessInfoByAuthCode(componentAppid, authorizationCode)
	if b && (asInfo.Status == "1" && authorizationCode == asInfo.AuthorizationCode) { //已经授权不进行重复授权
		return asInfo.Appid, nil
	}
	return s.doAuthorize(componentAppid, authorizationCode)
}

//更新授权状态为授权成功
func (s *ApiServiceImpl) AuthorizedSuccess(componentAppid string, authorizerAppid string) bool {
	//清除缓存
	defer delete(authAccessInfos, fmt.Sprintf("%s_%s", componentAppid, authorizerAppid))
	return s.Repo.UpdateAuthorizationAccessInfoStatus(componentAppid, authorizerAppid, "1")
	//model.UpdateAuthorizationInfoStatus(authorizerAppid,"1")
}

//取消授权
func (s *ApiServiceImpl) AuthorizedCancel(componentAppid string, authorizerAppid string) bool {
	//清除缓存
	defer delete(authAccessInfos, fmt.Sprintf("%s_%s", componentAppid, authorizerAppid))
	return s.Repo.UpdateAuthorizationAccessInfoStatus(componentAppid, authorizerAppid, "2")
	//model.UpdateAuthorizationInfoStatus(authorizerAppid,"2")
}
