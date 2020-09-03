package model

import (
	"github.com/kataras/golog"
)

var log = golog.Default

type ModelRepository interface {
	//Close() error

	GetAppBaseInfo(id int64) (*AppBaseInfo, error)
	FindAppBaseInfo(appId string) (*AppBaseInfo, bool)
	UpdateTicket(appId string, ticket string) bool
	UpdateAccessToken(appId string, accessToken string, expiresIn int64) bool
	UpdatePreAuthCode(appId string, preAuthCode string, expiresIn int64) bool

	FindAuthorizationAccessInfo(componentAppid string, appid string) (*AuthorizationAccessInfo, bool)
	FindAuthAccessInfoByAuthCode(componentAppid string, authorizationCode string) (*AuthorizationAccessInfo, bool)
	MergeAuthorizationAccessInfo(aai *AuthorizationAccessInfo) (bool, error)
	UpdateAuthorizationAccessInfoStatus(componentAppid string, appid string, status string) bool
	FindAuthorizationAppInfo(componentAppid string, appid string) (*AuthorizationAppInfo, bool)
	MergeAuthorizationAppInfo(aai *AuthorizationAppInfo) (bool, error)

	FindAuthorizationInfo(appid string) (*AuthorizationInfo, bool)
	MergeAuthorizationInfo(aai *AuthorizationInfo) (bool, error)

	//查找小程序的预览信息
	FindMiniPrgmPrevInfo(appid string) (*MiniPrgmPrevInfo, bool)
	//更新小程序的预览信息
	MergeMiniPrgmPrevInfo(aai *MiniPrgmPrevInfo) (ret bool, err error)
	// 查找代理API
	FindProxyApi(id int64) (*ProxyApi, bool)
	// 新增或修改代理API
	MergeProxyApi(aai *ProxyApi) (ret bool, err error)
}

type proxyModel struct {
	Engine *Engine `inject:""`
}

func NewApiModel() ModelRepository {
	return new(proxyModel)
}

//func New(conf *config.DBConf) ModelRepository {
//	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
//		conf.UserName,
//		conf.Password,
//		conf.Host,
//		conf.Port,
//		conf.DBName,
//		conf.Charset)
//	engine, err := xorm.NewEngine(conf.DriveName, dataSourceName)
//	if err != nil {
//		panic(err)
//	}
//	err = engine.Ping()
//	if err != nil {
//		panic(err)
//	}
//	return &proxyModel{
//		engine: engine,
//	}
//}
//
//func (m *proxyModel) Close() error {
//	return m.engine.Close()
//}
//
//func (m *proxyModel) Destroy() error {
//	log.Info("orm engine close...")
//	return m.engine.Close()
//}
