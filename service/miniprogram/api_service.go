package miniprogram

import (
	"github.com/jacktea/wxproxy/model"
	"sync"
	"github.com/kataras/golog"
)

var log = golog.Default

type MiniApiService interface {

}

type MimiApiServiceImpl struct {
	Repo model.ModelRepository `inject:""`
	lock sync.RWMutex
}

func New() MiniApiService {
	ret := new(MimiApiServiceImpl)
	return ret
}

