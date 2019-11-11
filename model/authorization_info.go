package model

import (
	"time"
)

type AuthorizationInfo struct {
	Id              int64  `xorm:"pk autoincr"`
	Appid           string //授权方appid
	NickName        string //授权方昵称
	HeadImg         string //授权方头像地址
	ServiceTypeInfo string //授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号
	VerifyTypeInfo  string //授权方认证类型
	UserName        string //授权方公众号的原始ID
	PrincipalName   string //公众号的主体名称
	Alias           string //授权方公众号所设置的微信号，可能为空
	BusinessInfo    string //业务功能开通情况
	QrcodeUrl       string //二维码图片的URL，开发者最好自行也进行保存
	Signature       string //帐号介绍
	//Status string							//状态 0 无效 1 授权成功 2 取消授权
	Miniprogram int       //是否为小程序 0 不是 1 是
	CreateTime  time.Time `xorm:"created"`
	ModifyTime  time.Time `xorm:"updated"`
}

func (a *AuthorizationInfo) TableName() string {
	return "wx_authorization_info"
}

func (m *proxyModel) FindAuthorizationInfo(appid string) (*AuthorizationInfo, bool) {
	items := make([]*AuthorizationInfo, 0)
	if err := m.Engine.Where("appid = ?", appid).Find(&items); err != nil {
		log.Error(err)
		return nil, false
	}
	if len(items) == 0 {
		return nil, false
	} else {
		return items[0], true
	}
}

func (m *proxyModel) MergeAuthorizationInfo(aai *AuthorizationInfo) (ret bool, err error) {
	items := make([]*AuthorizationInfo, 0)
	m.Engine.Where("appid = ?", aai.Appid).Find(&items)
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

/*
func UpdateAuthorizationInfoStatus(appid string,status string) bool {
	cnt,err := engine.
	Where("appid = ?",appid).
		Update(&AuthorizationInfo{
		Status:status,
	})
	return nil==err && cnt > 0
}
*/
