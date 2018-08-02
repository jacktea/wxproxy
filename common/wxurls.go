package common

const (
	//三方应用刷新TOKEN
	API_COMPONENT_TOKEN 	= "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	//三方应用刷新预授权码
	API_CREATE_PREAUTHCODE 	= "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode"
	//公众号/小程序授权
	API_QUERY_QUTH 			= "https://api.weixin.qq.com/cgi-bin/component/api_query_auth"
	//获取公众号/小程序信息
	API_GET_AUTHORIZER_INFO = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info"
	//刷新公众号/小程序Token
	API_AUTHORIZER_TOKEN    = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token"
	//发送客服短信
	SEND_CUSTOM_MSG			= "https://api.weixin.qq.com/cgi-bin/message/custom/send"
	//发送模板消息
	SEND_TPL_MSG			= "https://api.weixin.qq.com/cgi-bin/message/template/send"
	//创建带参二维码
	CREATE_PARAM_QRCODE = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
	//获取用户基本信息
	GET_USER_BASE_INFO = "https://api.weixin.qq.com/cgi-bin/user/info"

	AUTHORIZE_URL = "https://open.weixin.qq.com/connect/oauth2/authorize"

	ACCESS_TOKEN = "https://api.weixin.qq.com/sns/oauth2/component/access_token"

	USERINFO_URL = "https://api.weixin.qq.com/sns/userinfo"

	GZH_JSTICKET = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)
