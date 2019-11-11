package miniprogram

import (
	"github.com/jacktea/wxproxy/model"
	"github.com/jacktea/wxproxy/service/api"
	"github.com/kataras/golog"
	"io"
	"sync"
)

var log = golog.Default

type MiniApiService interface {
	GetQrCode(componentAppid, appid, path string, force bool) *QrCodeResp
	GetWxACode(componentAppid, appid, path string, force bool) *QrCodeResp
	GetWxACodeUnlimit(componentAppid, appid, page, scene string, force bool) *QrCodeResp
	//	GetWxQrCode(componentAppid ,appid , page, scene string,force bool) *QrCodeResp
	DownLoadQrCode(componentAppid, appid, fName string) (io.ReadCloser, error)
	//小程序预览
	MiniPreview(componentAppid, appid string, json string, path string, force bool) *QrCodeResp
}

type MimiApiServiceImpl struct {
	Repo model.ModelRepository `inject:""`
	Svr  api.ApiService        `inject:""`
	lock sync.RWMutex
}

func NewMiniApiService() MiniApiService {
	ret := new(MimiApiServiceImpl)
	return ret
}
