package miniprogram

import (
	. "github.com/jacktea/wxproxy/common"
	"github.com/jacktea/wxproxy/utils"
	"github.com/jacktea/wxproxy/model"
)

func (this *MimiApiServiceImpl) MiniPreview(componentAppid, appid string, json string,path string,force bool) *QrCodeResp {
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		return handleError(err)
	}
	if info,ok := this.Repo.FindMiniPrgmPrevInfo(appid);ok {
		jsonMd5 := utils.Md5(json)
		if jsonMd5 == info.JsonMd5 && !force {
			ret := new(QrCodeResp)
			ret.Errcode = 0
			ret.Errmsg = "ok"
			ret.Url = info.QrCodeUrl
			return ret
		}else {
			return this.miniPreview(componentAppid,appid,token,json,path)
		}
	}else {
		return this.miniPreview(componentAppid,appid,token,json,path)
	}
}

func (this *MimiApiServiceImpl) miniPreview(componentAppid, appid,token string, json,path string) *QrCodeResp {
	ret := new(QrCodeResp)
	if err := utils.HttpPostBodyJson(UPLOAD_MINI_COMMITCODE+"?access_token="+token,[]byte(json),ret);err != nil {
		return handleError(err)
	}
	if !ret.IsSuccess() {
		return ret
	}
	ret = this.GetQrCode(componentAppid,appid,path,true)
	info := new(model.MiniPrgmPrevInfo)
	info.Appid = appid
	info.JsonMd5 = utils.Md5(json)
	info.QrCodeUrl = ret.Url
	this.Repo.MergeMiniPrgmPrevInfo(info)
	return ret
}

