package main

import (
	"github.com/kataras/iris"
	api "github.com/jacktea/wxproxy/router/api"
	sapi "github.com/jacktea/wxproxy/service/api"
	"github.com/jacktea/wxproxy/router/notify"
	"github.com/jacktea/wxproxy/model"
	"github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/etcd"
	"github.com/jacktea/wxproxy/router/authorize"
	"github.com/jacktea/wxproxy/router/wxauth"
	"github.com/jacktea/wxproxy/bootstrap"

	kp "gopkg.in/alecthomas/kingpin.v2"
	"os/exec"
	"fmt"
	"os"
	"time"
	"github.com/kataras/golog"
	"bufio"
	"os/signal"
	"syscall"
	"github.com/jacktea/wxproxy/router/mini"
)

var (
	app	= kp.New("wxproxy", "微信代理").Author("zzg").Version("1.0.0")
	cmd *exec.Cmd
	log = golog.Default
	daemon = app.Flag("daemon","后台运行").Bool()
	forever = app.Flag("forever", "守护进程方式运行").Default("false").Bool()
	cfgPath = app.Flag("conf","配置文件路径").Short('c').String()
	logfile = app.Flag("log", "日志文件路径").Short('l').String()
)

func initConfig() (master bool,err error) {
	app.Parse(os.Args[1:])

	master = false
	if *logfile != "" {
		f, e := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			log.Fatal(e)
		}
		log.SetOutput(f)
	}
	if *daemon {
		args := make([]string,0)
		for _, arg := range os.Args[1:] {
			if arg != "--daemon" {
				args = append(args, arg)
			}
		}
		cmd = exec.Command(os.Args[0], args...)
		cmd.Start()
		f := ""
		if *forever {
			f = "forever "
		}
		log.Infof("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if *forever {
		args := make([]string,0)
		for _, arg := range os.Args[1:] {
			if arg != "--forever" {
				args = append(args, arg)
			}
		}
		go func() {
			for {
				if cmd != nil {
					cmd.Process.Kill()
				}
				cmd = exec.Command(os.Args[0], args...)
				cmdReaderStderr, err := cmd.StderrPipe()
				if err != nil {
					log.Errorf("ERR:%s,restarting...\n", err)
					continue
				}
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					log.Errorf("ERR:%s,restarting...\n", err)
					continue
				}
				scanner := bufio.NewScanner(cmdReader)
				scannerStdErr := bufio.NewScanner(cmdReaderStderr)
				go func() {
					for scanner.Scan() {
						fmt.Println(scanner.Text())
					}
				}()
				go func() {
					for scannerStdErr.Scan() {
						fmt.Println(scannerStdErr.Text())
					}
				}()
				if err := cmd.Start(); err != nil {
					log.Info("ERR:%s,restarting...\n", err)
					continue
				}
				pid := cmd.Process.Pid
				log.Infof("worker %s [PID] %d running...\n", os.Args[0], pid)
				if err := cmd.Wait(); err != nil {
					log.Errorf("ERR:%s,restarting...", err)
					continue
				}
				log.Infof("worker %s [PID] %d unexpected exited, restarting...\n", os.Args[0], pid)
				time.Sleep(time.Second * 5)
			}
		}()
		master = true
		return
	}
	return
}

func startWeb(conf *config.Config) {
	etcd.InitEtcd(conf.EtcdConf)
	engine := model.NewEngine(conf.DBConf)
	apiModel := model.NewApiModel()
	apiService := sapi.NewApiService()

	b := bootstrap.New(conf)
	b.RegistBeans(
		engine,
		apiModel,
		apiService,
		new(api.ApiAction),
		new(notify.NotifyAction),
		new(authorize.AuthorizeAction),
		new(wxauth.AuthAction),
		new(mini.MiniAction))
	b.AutoInject().
		Bootstrap().
		InitMiddleware().
		InitRouter().
		Start(
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutInterruptHandler)
}

func waitSignal() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping services...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func main() {

	master,err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf,err := config.NewConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	if master {
		waitSignal()
	}else {
		startWeb(conf)
	}
	//conf,_ := config.NewConfig("")



	//var g inject.Graph
	//g.Provide(
	//	&inject.Object{Value:repo},
	//	&inject.Object{Value:apiService},
	//	&inject.Object{Value:api.InitRouter(app)},
	//	&inject.Object{Value:notify.InitRouter(app)},
	//	&inject.Object{Value:authorize.InitRouter(app)},
	//	&inject.Object{Value:wxauth.InitRouter(app)},
	//	)
	//g.Populate()
	//
	////路由初始化
	////api.InitRouter(app,apiService)
	////notify.InitRouter(app)
	////authorize.InitRouter(app)
	////wxauth.InitRouter(app)
	//
	//app.Get(config.WXConf.HttpConf.ContextPath+"/test/{a}", func(c iris.Context) {
	//	url := fmt.Sprintf("%s://%s%s",utils.Scheme(c.Request()),c.Host(),c.Path())
	//	c.WriteString(url)
	//})
	//app.Get(config.WXConf.HttpConf.ContextPath+"/test/{a}/{b}", func(c iris.Context) {
	//	url := fmt.Sprintf("%s://%s%s",utils.Scheme(c.Request()),c.Host(),c.Path())
	//	c.WriteString(url)
	//})
	//
	//
	//
	//app.Run(iris.Addr(":8011"),
	//		iris.WithoutServerError(iris.ErrServerClosed),
	//		iris.WithoutInterruptHandler)
}
