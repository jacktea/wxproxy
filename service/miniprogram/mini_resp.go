package miniprogram

import (
	. "github.com/jacktea/wxproxy/service"
)

type ClientLoginResp struct {
	CommonResp
	Openid 			string 	`json:"openid"`			//用户唯一标识的openid。
	SessionKey 		int64 	`json:"session_key"`	//会话密钥。
}
