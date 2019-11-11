package model

import (
	"fmt"
	"github.com/jacktea/wxproxy/utils"
	"time"
)

type AppBaseInfo struct {
	Id                         int64  `xorm:"pk autoincr"`
	AppId                      string `xorm:"unique"`
	AppSecret                  string
	Token                      string
	EncodingAesKey             string
	ComponentVerifyTicket      string
	ComponentAccessToken       string
	ComponentAccessTokenExpire time.Time
	PreAuthCode                string
	PreAuthCodeExpire          time.Time
	CreateTime                 time.Time `xorm:"created"`
	ModifyTime                 time.Time `xorm:"updated"`
}

func (ab *AppBaseInfo) TableName() string {
	return "wx_app_base_info"
}

func (ab *AppBaseInfo) TokenIsExpired() bool {
	return ab.ComponentAccessToken == "" || time.Now().After(ab.ComponentAccessTokenExpire)
}

func (ab *AppBaseInfo) AuthCodeIsExpired() bool {
	return ab.PreAuthCode == "" || time.Now().After(ab.PreAuthCodeExpire)
}

/**
 * 根据ID获取信息
 */
func (m *proxyModel) GetAppBaseInfo(id int64) (*AppBaseInfo, error) {
	app := &AppBaseInfo{}
	_, err := m.Engine.ID(id).Get(app)
	return app, err
}

/**
 * 根据APPID获取信息
 */
func (m *proxyModel) FindAppBaseInfo(appId string) (*AppBaseInfo, bool) {
	apps := make([]*AppBaseInfo, 0)
	if err := m.Engine.Where("app_id = ? ", appId).Find(&apps); err != nil {
		log.Error(err)
		return nil, false
	}
	if len(apps) == 0 {
		return nil, false
	} else {
		return apps[0], true
	}
}

/**
 * 更新票据
 */
func (m *proxyModel) UpdateTicket(appId string, ticket string) bool {
	cnt, err := m.Engine.Where("app_id = ? ", appId).Update(&AppBaseInfo{ComponentVerifyTicket: ticket})
	return err == nil && cnt > 0
}

/**
 * 更新授权token
 */
func (m *proxyModel) UpdateAccessToken(appId string, accessToken string, expiresIn int64) bool {
	e := utils.ParseExpire(expiresIn)
	cnt, err := m.Engine.
		Where("app_id = ? ", appId).
		Update(&AppBaseInfo{ComponentAccessToken: accessToken,
			ComponentAccessTokenExpire: e})
	if err != nil {
		log.Error(err)
	}
	return err == nil && cnt > 0
}

/**
 * 更新预授权码
 */
func (m *proxyModel) UpdatePreAuthCode(appId string, preAuthCode string, expiresIn int64) bool {
	now := time.Now()
	d, err := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	if err != nil {
		log.Error(err)
		return false
	}
	e := now.Add(d)
	cnt, err := m.Engine.
		Where("app_id = ? ", appId).
		Update(&AppBaseInfo{PreAuthCode: preAuthCode,
			PreAuthCodeExpire: e})
	if err != nil {
		log.Error(err)
	}
	return err == nil && cnt > 0
}
