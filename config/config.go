package config

import (
	"github.com/ghodss/yaml"
	"github.com/jacktea/wxproxy/utils"
	"github.com/kataras/golog"
	"io/ioutil"
	"os"
	"path/filepath"
)

var log = golog.Default

type Config struct {
	CommonConf *CommonConf `json:"common"`
	DBConf     *DBConf     `json:"db"`
	EtcdConf   *EtcdConf   `json:"etcd"`
	HttpConf   *HttpConf   `json:"http"`
	OssConf    *OssConf    `json:"oss"`
}

type CommonConf struct {
	AppName  string `json:"app_name"`
	AppOwner string `json:"app_owner"`
	AppPath  string `json:"app_path"`
	WorkPath string `json:"work_path"`
	ConfPath string `json:"conf_path"` //配置文件路径
	DataPath string `json:"data_path"` //数据目录
	LogLevel string `json:"log_level"` //日志级别
}

type DBConf struct {
	DriveName string `json:"drive_name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	DBName    string `json:"db_name"`
	Charset   string `json:"charset"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type EtcdConf struct {
	Embed      bool     `json:"embed"`
	Name       string   `json:"name"`
	ClientPort int      `json:"client_port"`
	PeerPort   int      `json:"peer_port"`
	DataDir    string   `json:"data_dir"`
	Endpoints  []string `json:"endpoints"`
}

type HttpConf struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	ContextPath string `json:"context_path"`
}

type OssConf struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Enabled         bool   `json:"enabled"`
}

var WXConf *Config

func NewConfig(confPath string) (*Config, error) {
	var err error
	appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	workPath, _ := os.Getwd()
	WXConf = &Config{
		CommonConf: &CommonConf{
			AppName:  "wxproxy",
			AppOwner: "zzg",
			AppPath:  appPath,
			WorkPath: workPath,
			DataPath: "data",
			LogLevel: "info",
		},
		DBConf: &DBConf{
			DriveName: "mysql",
			Host:      "127.0.0.1",
			Port:      3306,
			DBName:    "wxappdb",
			Charset:   "utf8",
			UserName:  "jpss",
			Password:  "jpss541018",
		},
		EtcdConf: &EtcdConf{
			Embed:      true,
			Name:       "embedEtcd",
			ClientPort: 2379,
			PeerPort:   2380,
			DataDir:    "default.etcd",
			Endpoints:  []string{"http://127.0.0.1:2379"},
		},
		HttpConf: &HttpConf{
			Port:        8011,
			ContextPath: "/wxproxy",
		},
		OssConf: &OssConf{
			Enabled: false,
		},
	}

	//解析配置文件路径
	p, b := utils.ValidatePath([]string{
		confPath,
		filepath.Join(workPath, "conf", "wc.yml"),
		filepath.Join(appPath, "conf", "wc.yml"),
	})
	if b {
		data, _ := ioutil.ReadFile(p)
		err = yaml.Unmarshal(data, WXConf)
	}
	return WXConf, err
}
