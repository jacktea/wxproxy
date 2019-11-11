package api

import (
	"bytes"
	"encoding/json"
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/utils"
	"html/template"
)

var tempQrcode = template.
	Must(template.New("temp_qrcode").
		Parse(`{"expire_seconds": {{.ExpireSeconds}}, "action_name": "QR_STR_SCENE", "action_info": {"scene": {"scene_str": "{{.Identity}}"}}}`))
var foreverQrcode = template.
	Must(template.New("forerer_qrcode").
		Parse(`{"action_name": "QR_LIMIT_STR_SCENE", "action_info": {"scene": {"scene_str": "{{.Identity}}"}}}`))

/**
创建临时二维码
*/
func (a *ApiServiceImpl) CreateParamQrcode(
	componentAppid string, authorizerAppid string,
	identity string, expire_seconds int64, forever bool) (*QrcodeResp, error) {

	accessToken, err := a.GetAppAccessToken(componentAppid, authorizerAppid)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	if !forever {
		if expire_seconds > 2592000 {
			expire_seconds = 2592000
		}
		err := tempQrcode.Execute(buf, map[string]interface{}{
			"ExpireSeconds": expire_seconds,
			"Identity":      identity,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		err := foreverQrcode.Execute(buf, map[string]interface{}{
			"Identity": identity,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	uri := CREATE_PARAM_QRCODE + "?access_token=" + accessToken
	data, err := utils.HttpPostBody(uri, "application/json", buf.Bytes())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp := new(QrcodeResp)
	if err := json.Unmarshal(data, resp); err != nil {
		log.Error(err)
		return nil, err
	}
	return resp, nil
}
