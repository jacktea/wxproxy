package model

import (
	"time"
)

//小程序预览信息
type ProxyApi struct {
	Id          int64     `xorm:"pk autoincr"`
	Name        string    // 接口名称
	Url         string    // 接口地址
	Action      string    // 请求方式 1:get 2: post
	Type        string    // 类型 1 公众号 2 小程序 3 三方应用
	ExtraParams string    `xorm:"'extra_params'"`
	CreateTime  time.Time `xorm:"created"`
	ModifyTime  time.Time `xorm:"updated"`
}

func (a *ProxyApi) TableName() string {
	return "wx_proxy_api"
}

// 查找代理API
func (m *proxyModel) FindProxyApi(id int64) (*ProxyApi, bool) {
	items := make([]*ProxyApi, 0)
	if err := m.Engine.Where("id = ?", id).Find(&items); err != nil {
		log.Error(err)
		return nil, false
	}
	if len(items) == 0 {
		return nil, false
	} else {
		return items[0], true
	}
}

// 新增或修改代理API
func (m *proxyModel) MergeProxyApi(aai *ProxyApi) (ret bool, err error) {
	var cnt int64
	if aai.Id == 0 {
		cnt, err = m.Engine.Insert(aai)
	} else if info, ok := m.FindProxyApi(aai.Id); ok {
		cnt, err = m.Engine.ID(info.Id).Update(aai)
	} else {
		cnt, err = m.Engine.Insert(aai)
	}
	return cnt > 0, err
}
