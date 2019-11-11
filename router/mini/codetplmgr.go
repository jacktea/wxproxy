package mini

import (
	. "github.com/jacktea/wxproxy/common"
	"github.com/kataras/iris/v12"
)

// 获取草稿箱内的所有临时代码草稿
func (this *MiniAction) Gettemplatedraftlist(c iris.Context) {
	this.getCmpTransparentJson(c, GET_MINI_TEMPLATEDRAFTLIST)
}

// 获取代码模版库中的所有小程序代码模版
func (this *MiniAction) Gettemplatelist(c iris.Context) {
	this.getCmpTransparentJson(c, GET_MINI_TEMPLATELIST)
}

// 将草稿箱的草稿选为小程序代码模版
func (this *MiniAction) Addtotemplate(c iris.Context) {
	this.postCmpTransparentJson(c, ADD_MINI_TEMPLATE)
}

// 删除指定小程序代码模版
func (this *MiniAction) Deletetemplate(c iris.Context) {
	this.postCmpTransparentJson(c, DELETE_MINI_TEMPLATE)
}
