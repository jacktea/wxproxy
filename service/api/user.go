package api

import (
	"encoding/json"
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/utils"
)

/**
获取用户基本信息
*/
func (a *ApiServiceImpl) GetUserBaseInfo(componentAppid string, authorizerAppid string, openid string) (*UserBaseInfoResp, error) {
	accessToken, err := a.GetAppAccessToken(componentAppid, authorizerAppid)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	uri := GET_USER_BASE_INFO + "?access_token=" + accessToken + "&openid=" + openid + "&lang=zh_CN"
	data, err := utils.HttpGetBody(uri)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp := new(UserBaseInfoResp)
	if err := json.Unmarshal(data, resp); err != nil {
		log.Error(err)
		return nil, err
	}
	return resp, nil
}
