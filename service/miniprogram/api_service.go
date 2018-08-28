package miniprogram

import (
	"github.com/jacktea/wxproxy/model"
	"sync"
	"github.com/kataras/golog"
	"github.com/jacktea/wxproxy/service/api"
	"io"
)

var log = golog.Default

type MiniApiService interface {
	GetQrCode(componentAppid ,appid ,path string) *QrCodeResp
	DownLoadQrCode(componentAppid,appid,fName string) (io.ReadCloser,error)
}

type MimiApiServiceImpl struct {
	Repo model.ModelRepository `inject:""`
	Svr  api.ApiService	`inject:""`
	lock sync.RWMutex
}

func NewMiniApiService() MiniApiService {
	ret := new(MimiApiServiceImpl)
	return ret
}

