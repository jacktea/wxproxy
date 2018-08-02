package api

import (
	"fmt"
	"github.com/gomydodo/wxencrypter"
	"github.com/jacktea/wxproxy/model"
)

//获取三方应用基本信息
func (s *ApiServiceImpl) CacheFindAppBaseInfo(appId string) (info *model.AppBaseInfo, ok bool) {
	ok = false
	if info, ok = appBaseInfos[appId]; !ok {
		if info, ok = s.Repo.FindAppBaseInfo(appId); ok {
			appBaseInfos[appId] = info
		}
	}
	return
}

//获取授权应用授权信息
func (s *ApiServiceImpl) CacheFindAuthorizationAccessInfo(componentAppid string, appid string) (info *model.AuthorizationAccessInfo, ok bool) {
	ok = false
	key := fmt.Sprintf("%s_%s", componentAppid, appid)
	if info, ok = authAccessInfos[key]; !ok {
		if info, ok = s.Repo.FindAuthorizationAccessInfo(componentAppid, appid); ok {
			authAccessInfos[key] = info
		}
	}
	return
}

//获取授权应用基本信息
func (s *ApiServiceImpl) CacheFindAuthorizationInfo(appid string) (info *model.AuthorizationInfo, ok bool) {
	ok = false
	if info, ok = authInfos[appid]; !ok {
		if info, ok = s.Repo.FindAuthorizationInfo(appid); ok {
			authInfos[appid] = info
		}
	}
	return
}

func (s *ApiServiceImpl) CacheFindAuthorizationAppInfo(componentAppid string, appid string) (info *model.AuthorizationAppInfo, ok bool) {
	ok = false
	if info, ok = authAppInfos[appid]; !ok {
		if info, ok = s.Repo.FindAuthorizationAppInfo(componentAppid, appid); ok {
			authAppInfos[appid] = info
		}
	}
	return
}

func (s *ApiServiceImpl) CacheGetWxEncrypter(componentAppid string) (info *wxencrypter.Encrypter, err error) {
	var (
		ok  = false
		abi *model.AppBaseInfo
	)
	if info, ok = wxEncrypters[componentAppid]; !ok {
		if abi, ok = s.CacheFindAppBaseInfo(componentAppid); ok {
			if info, err = wxencrypter.NewEncrypter(abi.Token, abi.EncodingAesKey, abi.AppId); err == nil {
				wxEncrypters[componentAppid] = info
			}
		}
	}
	return
}
