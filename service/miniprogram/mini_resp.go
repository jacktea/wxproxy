package miniprogram

import (
	. "github.com/jacktea/wxproxy/service"
)

type ClientLoginResp struct {
	CommonResp
	Openid     string `json:"openid"`      //用户唯一标识的openid。
	SessionKey string `json:"session_key"` //会话密钥。
}

type QrCodeResp struct {
	CommonResp
	Url string `json:"url"` //体验二维码链接地址
}
