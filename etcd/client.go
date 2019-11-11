package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jacktea/wxproxy/config"
	"time"
)

type WatchExecutor interface {
	Match(e *clientv3.Event) bool

	Execute(e *clientv3.Event)
}

var cli *clientv3.Client

func InitEtcd(conf *config.EtcdConf) {
	var endpoints []string
	if conf.Embed {
		endpoints = []string{fmt.Sprintf("localhost:%d", conf.ClientPort)}
		startServer(conf)
	} else {
		endpoints = conf.Endpoints
	}
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second})
	if err != nil {
		goto END
	}

END:
	if err != nil {
		if cli != nil {
			cli.Close()
		}
	}
}

func Put(key string, value string, timeout int64) (ret bool) {
	ret = true
	resp, err := cli.Grant(context.TODO(), timeout)
	if err != nil {
		ret = false
	}
	_, err = cli.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		ret = false
	}
	return
}

func Get(key string) (ret string) {
	resp, err := cli.Get(context.TODO(), key)
	if err != nil {

	}
	if resp.Count > 0 {
		ret = fmt.Sprintf("%s", resp.Kvs[0].Value)
	}
	return
}

func Watch(prefix string, executors ...WatchExecutor) {
	rch := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {

			for _, we := range executors {
				if we.Match(ev) {
					we.Execute(ev)
				}
			}
		}
	}
}
