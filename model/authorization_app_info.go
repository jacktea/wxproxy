package model

import (
	"time"
)

type AuthorizationAppInfo struct {
	Id             int64     `xorm:"pk autoincr"`
	ComponentAppid string    //授权的三方应用
	Appid          string    //授权方appid
	NotifyUrl      string    //消息通知地址
	Mode           int       //模式 1(普通) 2(调试),调试模式下会向NotifyUrl和DebugNotifyUrl推送消息
	DebugNotifyUrl string    //调试模式下的推送地址
	CreateTime     time.Time `xorm:"created"`
	ModifyTime     time.Time `xorm:"updated"`
}

func (a *AuthorizationAppInfo) TableName() string {
	return "wx_authorization_app_info"
}

func (m *proxyModel) FindAuthorizationAppInfo(componentAppid string, appid string) (*AuthorizationAppInfo, bool) {
	items := make([]*AuthorizationAppInfo, 0)
	err := m.Engine.Where("component_appid = ? and appid = ?", componentAppid, appid).Find(&items)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	if len(items) == 0 {
		return nil, false
	} else {
		return items[0], true
	}
}

func (m *proxyModel) MergeAuthorizationAppInfo(aai *AuthorizationAppInfo) (ret bool, err error) {
	items := make([]*AuthorizationAppInfo, 0)
	if err = m.Engine.Where("component_appid = ? and appid = ?", aai.ComponentAppid, aai.Appid).Find(&items); err != nil {
		return false, err
	}
	var cnt int64
	if len(items) == 0 {
		cnt, err = m.Engine.Insert(aai)
	} else {
		cnt, err = m.Engine.ID(items[0].Id).Update(aai)
	}
	if err != nil {
		log.Error(err)
	}
	return cnt > 0, err
}
