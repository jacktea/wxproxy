package mini

import (
	. "github.com/jacktea/wxproxy/common"
	"github.com/kataras/iris"
)

// 获取小程序的基本信息
func (this *MiniAction) Getaccountbasicinfo(c iris.Context) {
	this.getTransparentJson(c,GET_MINI_BASEINFO)
}

//设置小程序业务域名
func (this *MiniAction) Setnickname(c iris.Context)  {
	this.postTransparentJson(c,SET_MINI_NICKNAME)
}