package model

import (
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jacktea/wxproxy/config"
	"fmt"
)

type Engine struct {
	*xorm.Engine
}

func NewEngine(conf *config.DBConf) *Engine {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.Charset)
	engine, err := xorm.NewEngine(conf.DriveName, dataSourceName)
	if err != nil {
		panic(err)
	}
	err = engine.Ping()
	if err != nil {
		panic(err)
	}
	return &Engine{
		Engine:engine,
	}
}

func (m *Engine) Destroy() error {
	log.Info("orm engine close...")
	return m.Close()
}
