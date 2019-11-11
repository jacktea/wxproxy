package api

import (
	. "github.com/jacktea/wxproxy/service"
)

//ComponentAccessToken 响应信息
type compTokenResp struct {
	CommonResp
	ComponentAccessToken string `json:"component_access_token"` //第三方平台access_token
	ExpiresIn            int64  `json:"expires_in"`             //过期时间(2小时)
}

//PreAuthCode 响应信息
type preAuthCodeResp struct {
	CommonResp
	PreAuthCode string `json:"pre_auth_code"` //预授权码
	ExpiresIn   int64  `json:"expires_in"`    //过期时间(10分钟)
}

//第三方应用授权(获取公众号/小程序访问token,获取公众号/小程序基本信息)
type authorizerInfoResp struct {
	CommonResp
	AuthorizerInfo    AuthorizerInfo    `json:"authorizer_info"`    //授权用户信息
	AuthorizationInfo AuthorizationInfo `json:"authorization_info"` //授权信息
}

//获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
type refreshAccessTokenResp struct {
	CommonResp
	AuthorizerAccessToken  string `json:"authorizer_access_token"`  //授权方令牌
	ExpiresIn              int64  `json:"expires_in"`               //有效期，为2小时
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"` //刷新令牌
}

//授权方的选项设置信息
type getAuthorizerOptionResp struct {
	CommonResp
	AuthorizerAppid string `json:"authorizer_appid"` //授权公众号或小程序的appid
	OptionName      string `json:"option_name"`      //选项名称
	OptionValue     string `json:"option_value"`     //选项值
}

type AuthorizerInfo struct {
	NickName        string           `json:"nick_name"`         //授权方昵称
	HeadImg         string           `json:"head_img"`          //授权方头像
	ServiceTypeInfo map[string]int   `json:"service_type_info"` //授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号
	VerifyTypeInfo  map[string]int   `json:"verify_type_info"`  //授权方认证类型，-1代表未认证，0代表微信认证，1代表新浪微博认证，2代表腾讯微博认证，3代表已资质认证通过但还未通过名称认证，4代表已资质认证通过、还未通过名称认证，但通过了新浪微博认证，5代表已资质认证通过、还未通过名称认证，但通过了腾讯微博认证
	UserName        string           `json:"user_name"`         //授权方公众号的原始ID
	PrincipalName   string           `json:"principal_name"`    //公众号的主体名称
	BusinessInfo    map[string]int   `json:"business_info"`     //open_store:是否开通微信门店功能 open_scan:是否开通微信扫商品功能 open_pay:是否开通微信支付功能 open_card:是否开通微信卡券功能 open_shake:是否开通微信摇一摇功能
	Alias           string           `json:"alias"`             //授权方公众号所设置的微信号，可能为空
	QrcodeUrl       string           `json:"qrcode_url"`        //二维码图片的URL，开发者最好自行也进行保存
	Signature       string           `json:"signature"`         //帐号介绍
	MiniProgramInfo *MiniProgramInfo `json:"MiniProgramInfo"`   //小程序信息
}

type AuthorizationInfo struct {
	Appid                  string         `json:"appid"`                    //授权方appid
	AuthorizerAppid        string         `json:"authorizer_appid"`         //授权方appid(与上相同，对应不同接口)
	AuthorizerAccessToken  string         `json:"authorizer_access_token"`  //授权方接口调用凭据(在授权的公众号或小程序具备API权限时，才有此返回值)，也简称为令牌
	ExpiresIn              int64          `json:"expires_in"`               //有效期(在授权的公众号或小程序具备API权限时，才有此返回值)，为2小时
	AuthorizerRefreshToken string         `json:"authorizer_refresh_token"` //刷新令牌
	FuncInfo               []FuncInfoItem `json:"func_info"`
}

type FuncInfoItem struct {
	FuncscopeCategory map[string]int `json:"funcscope_category"`
}

type MiniProgramInfo struct {
	Network     map[string]interface{} `json:"network"`
	Categories  []map[string]string    `json:"categories"`
	VisitStatus int64                  `json:"visit_status"`
}

type AuthAccessToken struct {
	CommonResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"`
	Scope        string `json:"scope"`
}

type AuthUserinfo struct {
	CommonResp
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	Unionid    string `json:"unionid"`
}

type JsTicket struct {
	CommonResp
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

/**
生成带参二维码返回值
*/
type QrcodeResp struct {
	CommonResp
	Ticket        string `json:"ticket"`         //获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	ExpireSeconds int64  `json:"expire_seconds"` //该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。
	Url           string `json:"url"`            //二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

/**
用户基本信息
*/
type UserBaseInfoResp struct {
	CommonResp
	Subscribe      int      `json:"subscribe"`       //用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	Openid         string   `json:"openid"`          //用户的标识，对当前公众号唯一
	Nickname       string   `json:"nickname"`        //用户的昵称
	Sex            int      `json:"sex"`             //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City           string   `json:"city"`            //用户所在城市
	Country        string   `json:"country"`         //用户所在国家
	Province       string   `json:"province"`        //用户所在省份
	Language       string   `json:"language"`        //用户的语言，简体中文为zh_CN
	Headimgurl     string   `json:"headimgurl"`      //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	SubscribeTime  int64    `json:"subscribe_time"`  //用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	Unionid        string   `json:"unionid"`         //只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Remark         string   `json:"remark"`          //公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	Groupid        int      `json:"groupid"`         //用户所在的分组ID（兼容旧的用户分组接口）
	TagidList      []string `json:"tagid_list"`      //用户被打上的标签ID列表
	SubscribeScene string   `json:"subscribe_scene"` //返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	QrScene        int64    `json:"qr_scene"`        //二维码扫码场景（开发者自定义）
	QrSceneStr     string   `json:"qr_scene_str"`    //二维码扫码场景描述（开发者自定义）
}
