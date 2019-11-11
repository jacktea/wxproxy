package model

import (
	"time"
)

//小程序预览信息
type MiniPrgmPrevInfo struct {
	Id         int64     `xorm:"pk autoincr"`
	Appid      string    //小程序ID
	JsonMd5    string    //上次上传的数据指纹
	QrCodeUrl  string    //预览图片地址
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

func (a *MiniPrgmPrevInfo) TableName() string {
	return "wx_mini_prgm_prev_info"
}

//查找小程序的预览信息
func (m *proxyModel) FindMiniPrgmPrevInfo(appid string) (*MiniPrgmPrevInfo, bool) {
	items := make([]*MiniPrgmPrevInfo, 0)
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

//更新小程序的预览信息
func (m *proxyModel) MergeMiniPrgmPrevInfo(aai *MiniPrgmPrevInfo) (ret bool, err error) {
	var cnt int64
	if info, ok := m.FindMiniPrgmPrevInfo(aai.Appid); ok {
		cnt, err = m.Engine.ID(info.Id).Update(aai)
	} else {
		cnt, err = m.Engine.Insert(aai)
	}
	return cnt > 0, err
}
