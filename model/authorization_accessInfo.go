package model

import (
	"time"
)

type AuthorizationAccessInfo struct {
	Id                int64 `xorm:"pk autoincr"`
	ComponentAppid    string
	AuthorizationCode string
	Appid             string
	AccessToken       string
	AccessTokenExpire time.Time
	RefreshToken      string
	FuncInfo          string
	Status            string    //状态 0 无效 1 授权成功 2 取消授权
	CreateTime        time.Time `xorm:"created"`
	ModifyTime        time.Time `xorm:"updated"`
}

func (a *AuthorizationAccessInfo) TableName() string {
	return "wx_authorization_access_info"
}

func (ab *AuthorizationAccessInfo) TokenIsExpired() bool {
	return ab.AccessToken == "" || time.Now().After(ab.AccessTokenExpire)
}

func (m *proxyModel) FindAuthorizationAccessInfo(componentAppid string, appid string) (*AuthorizationAccessInfo, bool) {
	items := make([]*AuthorizationAccessInfo, 0)
	err := m.Engine.Where("component_appid = ? and appid = ? and status = ? ", componentAppid, appid, "1").Find(&items)
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

func (m *proxyModel) FindAuthAccessInfoByAuthCode(componentAppid string, authorizationCode string) (*AuthorizationAccessInfo, bool) {
	items := make([]*AuthorizationAccessInfo, 0)
	if err := m.Engine.Where("component_appid = ? and authorization_code = ?", componentAppid, authorizationCode).Find(&items); err != nil {
		log.Error(err)
		return nil, false
	}
	if len(items) == 0 {
		return nil, false
	} else {
		return items[0], true
	}
}

func (m *proxyModel) MergeAuthorizationAccessInfo(aai *AuthorizationAccessInfo) (ret bool, err error) {
	items := make([]*AuthorizationAccessInfo, 0)
	err = m.Engine.Where("component_appid = ? and appid = ?", aai.ComponentAppid, aai.Appid).Find(&items)
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

func (m *proxyModel) UpdateAuthorizationAccessInfoStatus(componentAppid string, appid string, status string) bool {
	cnt, err := m.Engine.
		Where("component_appid = ? and appid = ?", componentAppid, appid).
		Update(&AuthorizationAccessInfo{
			Status: status,
		})
	if err != nil {
		log.Error(err)
	}
	return nil == err && cnt > 0
}
