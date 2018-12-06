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
	// 添加模板
	ADD_TEMPLATE			= "https://api.weixin.qq.com/cgi-bin/template/api_add_template"
	//创建带参二维码
	CREATE_PARAM_QRCODE = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
	//获取用户基本信息
	GET_USER_BASE_INFO = "https://api.weixin.qq.com/cgi-bin/user/info"

	AUTHORIZE_URL = "https://open.weixin.qq.com/connect/oauth2/authorize"

	ACCESS_TOKEN = "https://api.weixin.qq.com/sns/oauth2/component/access_token"

	USERINFO_URL = "https://api.weixin.qq.com/sns/userinfo"

	GZH_JSTICKET = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"

	//修改服务器地址 miniprogram/domain.go
	//1. 设置小程序服务域名 https://api.weixin.qq.com/wxa/modify_domain?access_token=TOKEN
	SET_MINI_SVRDOMAIN = "https://api.weixin.qq.com/wxa/modify_domain"
	//2. 设置小程序业务域名 https://api.weixin.qq.com/wxa/setwebviewdomain?access_token=TOKEN
	SET_MINI_WEBDOMAIN = "https://api.weixin.qq.com/wxa/setwebviewdomain"

	//小程序基本信息设置 miniprogram/info.go
	//1. 获取小程序的基本信息 GET https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=TOKEN
	GET_MINI_BASEINFO = "https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo"
	//2. 小程序名称设置及改名 POST https://api.weixin.qq.com/wxa/setnickname?access_token=TOKEN
	SET_MINI_NICKNAME = "https://api.weixin.qq.com/wxa/setnickname"
	//3. 小程序改名审核状态查询 POST https://api.weixin.qq.com/wxa/api_wxa_querynickname?access_token=TOKEN
	//4. 微信认证名称检测 POST https://api.weixin.qq.com/cgi-bin/wxverify/checkwxverifynickname?access_token=TOKEN
	//5. 修改头像 POST https://api.weixin.qq.com/cgi-bin/account/modifyheadimage?access_token=TOKEN
	//6. 修改功能介绍 POST https://api.weixin.qq.com/cgi-bin/account/modifysignature?access_token=TOKEN
	//7. 换绑小程序管理员接口
	//7.1 从第三方平台跳转至微信公众平台授权注册页面
	//7.2 小程序新旧管理员填写信息，扫码确认提交后跳转回第三方平台
	//7.3 跳转至第三方平台，第三方平台调用快速注册API完成管理员换绑。
	//8 类目相关接口
	//8.1 获取账号可以设置的所有类目
	//8.2 添加类目
	//8.3 删除类目
	//8.4 获取账号已经设置的所有类目
	//8.5 修改类目

	//成员管理 miniprogram/bind.go
	//1. 绑定微信用户为小程序体验者 https://api.weixin.qq.com/wxa/bind_tester?access_token=TOKEN
	BIND_MINI_TESTER = "https://api.weixin.qq.com/wxa/bind_tester"
	//2. 解除绑定微信小程序体验者 https://api.weixin.qq.com/wxa/unbind_tester?access_token=TOKEN
	UNBIND_MINI_TESTER = "https://api.weixin.qq.com/wxa/unbind_tester"
	//3. 获取小程序体验者 https://api.weixin.qq.com/wxa/memberauth?access_token=TOKEN
	GET_MINI_TESTER = "https://api.weixin.qq.com/wxa/memberauth"


	//代码管理 miniprogram/code.go
	//1. 为授权的小程序帐号上传小程序代码 https://api.weixin.qq.com/wxa/commit?access_token=TOKEN
	UPLOAD_MINI_COMMITCODE = "https://api.weixin.qq.com/wxa/commit"
	//2. 获取体验小程序的体验二维码 https://api.weixin.qq.com/wxa/ get_qrcode?access_token=TOKEN&path=page%2Findex%3Faction%3D1
	GET_MINI_QRCODE = "https://api.weixin.qq.com/wxa/get_qrcode"
	//2. 获取小程序的二维码 POST https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=ACCESS_TOKEN
	GET_MINI_WXQRCODE = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	//获取小程序码，适用于需要的码数量较少的业务场景。通过该接口生成的小程序码，永久有效，有数量限制
	//POST https://api.weixin.qq.com/wxa/getwxacode?access_token=ACCESS_TOKEN
	GET_MINI_WXACODE = "https://api.weixin.qq.com/wxa/getwxacode"
	//获取小程序码，适用于需要的码数量极多的业务场景。通过该接口生成的小程序码，永久有效，数量暂无限制。
	//POST https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=ACCESS_TOKEN
	GET_MINI_WXACODEUNLIMIT = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	//3. 获取授权小程序帐号的可选类目 https://api.weixin.qq.com/wxa/get_category?access_token=TOKEN
	GET_MINI_CATEGORY = "https://api.weixin.qq.com/wxa/get_category"
	//4. 获取小程序的第三方提交代码的页面配置（仅供第三方开发者代小程序调用） GET https://api.weixin.qq.com/wxa/get_page?access_token=TOKEN
	GET_MINI_PAGE = "https://api.weixin.qq.com/wxa/get_page"
	//5. 将第三方提交的代码包提交审核（仅供第三方开发者代小程序调用）POST https://api.weixin.qq.com/wxa/submit_audit?access_token=TOKEN
	SUBMIT_MINI_AUDIT = "https://api.weixin.qq.com/wxa/submit_audit"
	//7. 查询某个指定版本的审核状态（仅供第三方代小程序调用）POST https://api.weixin.qq.com/wxa/get_auditstatus?access_token=TOKEN
	QUERY_MINI_AUDITSTATUS = "https://api.weixin.qq.com/wxa/get_auditstatus"
	//8. 查询最新一次提交的审核状态（仅供第三方代小程序调用）GET https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token=TOKEN
	QUERY_MINI_LASTAUDITSTATUS = "https://api.weixin.qq.com/wxa/get_latest_auditstatus"
	//9. 发布已通过审核的小程序（仅供第三方代小程序调用）POST https://api.weixin.qq.com/wxa/release?access_token=TOKEN
	DO_MINI_RELEASE = "https://api.weixin.qq.com/wxa/release"
	//10. 修改小程序线上代码的可见状态（仅供第三方代小程序调用）POST https://api.weixin.qq.com/wxa/change_visitstatus?access_token=TOKEN
	CHANGE_MINI_VISITSTATUS = "https://api.weixin.qq.com/wxa/change_visitstatus"
	//11. 小程序版本回退（仅供第三方代小程序调用）GET https://api.weixin.qq.com/wxa/revertcoderelease?access_token=TOKEN
	DO_MINI_REVERTCODERELEASE = "https://api.weixin.qq.com/wxa/revertcoderelease"
	//12. 查询当前设置的最低基础库版本及各版本用户占比（仅供第三方代小程序调用）POST https://api.weixin.qq.com/cgi-bin/wxopen/getweappsupportversion?access_token=TOKEN
	GET_MINI_WEAPPSUPPORTVERSION = "https://api.weixin.qq.com/cgi-bin/wxopen/getweappsupportversion"
	//13. 设置最低基础库版本（仅供第三方代小程序调用）POST https://api.weixin.qq.com/cgi-bin/wxopen/setweappsupportversion?access_token=TOKEN
	SET_MINI_WEAPPSUPPORTVERSION = "https://api.weixin.qq.com/cgi-bin/wxopen/setweappsupportversion"
	//14. 设置小程序“扫普通链接二维码打开小程序”能力
	//14.1 增加或修改二维码规则 POST https://api.weixin.qq.com/cgi-bin/wxopen/qrcodejumpadd?access_token=TOKEN

	//14.2 获取已设置的二维码规则 POST https://api.weixin.qq.com/cgi-bin/wxopen/qrcodejumpget?access_token=TOKEN

	//14.3 获取校验文件名称及内容 POST https://api.weixin.qq.com/cgi-bin/wxopen/qrcodejumpdownload?access_token=TOKEN

	//14.4 删除已设置的二维码规则 POST https://api.weixin.qq.com/cgi-bin/wxopen/qrcodejumpdelete?access_token=TOKEN

	//14.5 发布已设置的二维码规则 POST https://api.weixin.qq.com/cgi-bin/wxopen/qrcodejumppublish?access_token=TOKEN

	//15 小程序审核撤回 单个帐号每天审核撤回次数最多不超过1次，一个月不超过10次。GET https://api.weixin.qq.com/wxa/undocodeaudit?access_token=TOKEN

	//16 小程序分阶段发布
	//16.1 分阶段发布接口 POST https://api.weixin.qq.com/wxa/grayrelease?access_token=TOKEN
	//16.2 取消分阶段发布 GET https://api.weixin.qq.com/wxa/revertgrayrelease?access_token=TOKEN
	//16.3 查询当前分阶段发布详情 GET https://api.weixin.qq.com/wxa/getgrayreleaseplan?access_token=TOKEN

	//小程序代码模版库管理 miniprogram/codetplmgr.go
	//1.获取草稿箱内的所有临时代码草稿 GET https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=TOKEN
	GET_MINI_TEMPLATEDRAFTLIST = "https://api.weixin.qq.com/wxa/gettemplatedraftlist"
	//2.获取代码模版库中的所有小程序代码模版 GET https://api.weixin.qq.com/wxa/gettemplatelist?access_token=TOKEN
	GET_MINI_TEMPLATELIST = "https://api.weixin.qq.com/wxa/gettemplatelist"
	//3.将草稿箱的草稿选为小程序代码模版 POST https://api.weixin.qq.com/wxa/addtotemplate?access_token=TOKEN
	ADD_MINI_TEMPLATE = "https://api.weixin.qq.com/wxa/addtotemplate"
	//4.删除指定小程序代码模版 POST https://api.weixin.qq.com/wxa/deletetemplate?access_token=TOKEN
	DELETE_MINI_TEMPLATE = "https://api.weixin.qq.com/wxa/deletetemplate"

	//微信登录	miniprogram/user.go
	//code 换取 session_key GET https://api.weixin.qq.com/sns/component/jscode2session?appid=APPID&js_code=JSCODE&grant_type=authorization_code&component_appid=COMPONENT_APPID&component_access_token=ACCESS_TOKEN
	USER_MINI_LOGIN = "https://api.weixin.qq.com/sns/component/jscode2session"

	//小程序模板设置	miniprogram/tpl.go
	//1.获取小程序模板库标题列表 POST https://api.weixin.qq.com/cgi-bin/wxopen/template/library/list?access_token=ACCESS_TOKEN
	//2.获取模板库某个模板标题下关键词库 POST https://api.weixin.qq.com/cgi-bin/wxopen/template/library/get?access_token=ACCESS_TOKEN
	//3.组合模板并添加至帐号下的个人模板库 POST https://api.weixin.qq.com/cgi-bin/wxopen/template/add?access_token=ACCESS_TOKEN
	//4.获取帐号下已存在的模板列表 POST https://api.weixin.qq.com/cgi-bin/wxopen/template/list?access_token=ACCESS_TOKEN
	//5.删除帐号下的某个模板 POST https://api.weixin.qq.com/cgi-bin/wxopen/template/del?access_token=ACCESS_TOKEN

	//微信开放平台帐号管理		miniprogram/open.go
	//1. 创建 开放平台帐号并绑定公众号/小程序 POST https://api.weixin.qq.com/cgi-bin/open/create?access_token=xxxx
	//2. 将 公众号/小程序绑定到开放平台帐号下 POST https://api.weixin.qq.com/cgi-bin/open/bind?access_token=xxxx
	//3. 将公众号/小程序从开放平台帐号下解绑 POST https://api.weixin.qq.com/cgi-bin/open/unbind?access_token=xxxx
	//4. 获取公众号/小程序所绑定的开放平台帐号 POST https://api.weixin.qq.com/cgi-bin/open/unbind?access_token=xxxx

	//基础信息设置	miniprogram/privacy.go
	//1. 设置小程序隐私设置（是否可被搜索）POST https://api.weixin.qq.com/wxa/changewxasearchstatus?access_token=TOKEN
	//2. 查询小程序当前隐私设置（是否可被搜索）POST https://api.weixin.qq.com/wxa/getwxasearchstatus?access_token=TOKEN

	//小程序插件管理权限集 miniprogram/plugins.go
	//1. 申请使用插件接口 POST https://api.weixin.qq.com/wxa/plugin?access_token=TOKEN
	//2. 查询已添加的插件 POST https://api.weixin.qq.com/wxa/plugin?access_token=TOKEN
	//3. 删除已添加的插件 POST https://api.weixin.qq.com/wxa/plugin?access_token=TOKEN
















)
