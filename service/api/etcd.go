package api

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/jacktea/wxproxy/etcd"
	"strings"
)

const (
	GLOBAL_PREFIX                 = "wx/"
	COMPONENT_ACCESS_TOKEN_PREFIX = GLOBAL_PREFIX + "component_access_token/"
	PRE_AUTH_CODE_PREFIX          = GLOBAL_PREFIX + "pre_auth_code/"
	REFRESH_TOKEN_PREFIX          = GLOBAL_PREFIX + "refresh_token/"
)

type DeleteWatch struct {
	svr ApiService
}

func (w *DeleteWatch) Match(ev *clientv3.Event) bool {
	key := string(ev.Kv.Key)
	return ev.Type == 1 &&
		(strings.HasPrefix(key, COMPONENT_ACCESS_TOKEN_PREFIX) ||
			strings.HasPrefix(key, PRE_AUTH_CODE_PREFIX) ||
			strings.HasPrefix(key, REFRESH_TOKEN_PREFIX))
}

func (w *DeleteWatch) Execute(ev *clientv3.Event) {
	key := string(ev.Kv.Key)
	arr := strings.Split(key, "/")
	appId := arr[len(arr)-1]
	if strings.HasPrefix(key, COMPONENT_ACCESS_TOKEN_PREFIX) {
		strings.Split(key, "/")
		log.Debug("Etcd to UpdateAccessToken ",
			" appid="+appId)
		go w.svr.UpdateAccessToken(appId, true)
	}
	if strings.HasPrefix(key, PRE_AUTH_CODE_PREFIX) {
		log.Debug("Etcd to UpdatePreAuthCode",
			" appid="+appId)
		go w.svr.UpdatePreAuthCode(appId, true)
	}
	if strings.HasPrefix(key, REFRESH_TOKEN_PREFIX) {
		componentAppid := arr[len(arr)-2]
		log.Debug("Etcd to RefreshAuthorizationToken",
			" componentAppid="+componentAppid,
			" appid="+appId)
		go w.svr.RefreshAuthorizationToken(componentAppid, appId, true)
	}
}

func (s *ApiServiceImpl) StartEtcdWatch() {
	go etcd.Watch(GLOBAL_PREFIX, &DeleteWatch{
		svr: s,
	})
}
