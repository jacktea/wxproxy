package api

import (
	. "github.com/jacktea/wxproxy/common"
)

//发送客服消息
func (s *ApiServiceImpl) SendCustomMsg(componentAppid string, authorizerAppid string, v interface{}) error {
	token, err := s.GetAppAccessToken(componentAppid, authorizerAppid)
	if err != nil {
		return err
	}
	return sendJsonMsg(SEND_CUSTOM_MSG, token, v)
}

//发送模板消息
func (s *ApiServiceImpl) SendTplMsg(componentAppid string, authorizerAppid string, v interface{}) error {
	token, err := s.GetAppAccessToken(componentAppid, authorizerAppid)
	if err != nil {
		return err
	}
	return sendJsonMsg(SEND_TPL_MSG, token, v)
}
