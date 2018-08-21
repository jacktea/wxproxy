package api

import (
	"github.com/jacktea/wxproxy/model"
	"sync"
	"github.com/kataras/golog"
	"github.com/gomydodo/wxencrypter"
)

var log = golog.Default

type ApiService interface {

	CacheFindAppBaseInfo(appId string) (info *model.AppBaseInfo, ok bool)
	CacheFindAuthorizationAccessInfo(componentAppid string, appid string) (info *model.AuthorizationAccessInfo, ok bool)
	CacheFindAuthorizationInfo(appid string) (info *model.AuthorizationInfo, ok bool)
	CacheFindAuthorizationAppInfo(componentAppid string, appid string) (info *model.AuthorizationAppInfo, ok bool)
	CacheGetWxEncrypter(componentAppid string) (info *wxencrypter.Encrypter, err error)

	//开启Etcd监控
	StartEtcdWatch()

	//更新三方应用票据
	UpdateTicket(appId string,ticket string) error
	//更新三方应用访问Token
	UpdateAccessToken(appId string,force bool) (*model.AppBaseInfo,error)
	//更新第三方应用预授权码
	UpdatePreAuthCode(appId string,force bool) error
	//获取第三方应用访问token
	GetComponentAppAccessToken(componentAppid string) (string, error)
	//刷新托管公众号(小程序)访问Token
	RefreshAuthorizationToken(componentAppid string,appid string,force bool) (*model.AuthorizationAccessInfo,error)
	//获取托管公众号(小程序)访问Token
	GetAppAccessToken(componentAppid string, appid string) (string, error)
	//更新托管公众号(小程序)信息
	UpdateAuthorizerInfo(componentAppid string,authorizerAppid string)  error
	//更新托管公众号(小程序)回调地址
	UpdateAuthorizationAppNotifyUrl(componentAppid string,appid string,notifyUrl string,mode int,debugNotifyUrl string) error

	//公众号(小程序)授权时通知回调
	DoAuthorize(componentAppid string,authorizerAppid string,authorizationCode string) (string,error)
	//公众号(小程序)授权时网页回调
	DoAuthorizerInfo(componentAppid string,authorizationCode string) (string,error)
	//更新授权状态为授权成功
	AuthorizedSuccess(componentAppid string,authorizerAppid string) bool
	//取消授权
	AuthorizedCancel(componentAppid string,authorizerAppid string) bool

	//发送客服消息
	SendCustomMsg(componentAppid string,authorizerAppid string,v interface{}) error
	//发送模板消息
	SendTplMsg(componentAppid string,authorizerAppid string,v interface{}) error

	//创建关注二维码
	CreateParamQrcode(componentAppid string,authorizerAppid string,identity string,expire_seconds int64,forever bool) (*QrcodeResp,error)

	//获取用户基本信息
	GetUserBaseInfo(componentAppid string,authorizerAppid string,openid string) (*UserBaseInfoResp,error)

}

type ApiServiceImpl struct {
	Repo model.ModelRepository `inject:""`
	lock sync.RWMutex
}

func NewApiService() ApiService {
	ret := new(ApiServiceImpl)
	ret.StartEtcdWatch()
	return ret
}



