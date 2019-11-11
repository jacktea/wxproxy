package etcd

import (
	"fmt"
	"github.com/coreos/etcd/embed"
	"github.com/jacktea/wxproxy/config"
	"log"
	"net/url"
	"time"
)

func startServer(conf *config.EtcdConf) {
	ch := make(chan *struct{})
	go func() {
		lcurl, _ := url.Parse(fmt.Sprintf("http://localhost:%d", conf.ClientPort))
		lpurl, _ := url.Parse(fmt.Sprintf("http://localhost:%d", conf.PeerPort))
		cfg := embed.NewConfig()
		cfg.Dir = conf.DataDir
		cfg.ACUrls = []url.URL{*lcurl}
		cfg.APUrls = []url.URL{*lpurl}
		e, err := embed.StartEtcd(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer e.Close()
		select {
		case <-e.Server.ReadyNotify():
			log.Printf("Server is ready!")
			ch <- &struct{}{}
		case <-time.After(60 * time.Second):
			e.Server.Stop() // trigger a shutdown
			log.Printf("Server took too long to start!")
			ch <- &struct{}{}
		}
		log.Fatal(<-e.Err())
		ch <- &struct{}{}
	}()
	<-ch
}
