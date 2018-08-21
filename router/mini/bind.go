package mini

import (
	"github.com/kataras/iris"
	. "github.com/jacktea/wxproxy/common"
)

// 绑定微信用户为小程序体验者
func (this *MiniAction) BindTester(c iris.Context)  {
	this.postTransparentJson(c,BIND_MINI_TESTER)
}

// 解除绑定小程序的体验者
func (this *MiniAction) UnBindTester(c iris.Context)  {
	this.postTransparentJson(c,UNBIND_MINI_TESTER)
}

// 获取体验者列表
func (this *MiniAction) Memberauth(c iris.Context)  {
	this.postTransparentJson(c,GET_MINI_TESTER)
}