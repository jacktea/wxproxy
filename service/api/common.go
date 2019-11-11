package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomydodo/wxencrypter"
	"github.com/jacktea/wxproxy/model"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/utils"
	"strconv"
	"strings"
)

var (
	wxEncrypters = make(map[string]*wxencrypter.Encrypter)
	//第三方应用信息缓存
	appBaseInfos = make(map[string]*model.AppBaseInfo)
	//公众号/小程序基本信息
	authInfos = make(map[string]*model.AuthorizationInfo)
	//公众号/小程序授权信息
	authAccessInfos = make(map[string]*model.AuthorizationAccessInfo)
	//公众号/小程序接入相关配置
	authAppInfos = make(map[string]*model.AuthorizationAppInfo)
)

//执行微信API
func executeApi(uri string, params map[string]interface{}, ret interface{}) (err error) {
	if err = utils.HttpPostJson(uri, params, ret); err != nil {
		log.Error(err)
		return err
	}
	if resp, ok := ret.(service.Resp); ok {
		if !resp.IsSuccess() {
			err = createWxApiError(resp)
			log.Error(err)
		}
	}
	return
}

func createWxApiError(resp service.Resp) error {
	return errors.New(fmt.Sprintf("调用微信服务出错,错误码:%d,错误信息:%s", resp.GetErrcode(), resp.GetErrmsg()))
}

func toStr(items []FuncInfoItem) string {
	if nil == items {
		return ""
	}
	slice := make([]string, len(items))
	for i := 0; i < len(slice); i++ {
		slice[i] = strconv.Itoa(items[i].FuncscopeCategory["id"])
	}
	return strings.Join(slice, ",")
}

func toJsonStr(i interface{}) string {
	if data, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return string(data)
	}
}

func sendJsonMsg(gw, accessToken string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	str, ok := v.(string)
	if ok {
		data = []byte(str)
		if !(len(data) > 0 && data[0] == '{' && data[len(data)-1] == '}') {
			log.Debug(str)
			err = errors.New("消息格式不正确")
			return err
		}
	} else {
		if data, err = json.Marshal(v); err != nil {
			return err
		}
	}
	uri := gw + "?access_token=" + accessToken
	data, err = utils.HttpPostBody(uri, "application/json", data)
	if err != nil {
		return err
	}
	resp := new(service.CommonResp)
	if err = json.Unmarshal(data, resp); err != nil {
		return err
	}
	if !resp.IsSuccess() {
		err = createWxApiError(resp)
		log.Error(err)
	}
	return err
}
