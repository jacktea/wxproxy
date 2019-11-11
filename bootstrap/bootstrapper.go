package bootstrap

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/facebookgo/inject"
	"github.com/gorilla/securecookie"
	"github.com/jacktea/wxproxy/config"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/thoas/stats"
)

type Router interface {
	InitRouter(app *Bootstrapper)
}

type SafeDestroy interface {
	Destroy() error
}

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	Conf         *config.Config
	AppName      string
	AppOwner     string
	ContextPath  string
	AppSpawnDate time.Time
	Sessions     *sessions.Sessions
	Beans        []interface{}
}

//// New returns a new Bootstrapper.
//func New(appName, appOwner string) *Bootstrapper {
//	b := &Bootstrapper{
//		AppName:      appName,
//		AppOwner:     appOwner,
//		AppSpawnDate: time.Now(),
//		Application:  iris.New(),
//		Beans:        make([]interface{}, 0),
//	}
//	return b
//}

func New(conf *config.Config) *Bootstrapper {
	return &Bootstrapper{
		AppName:      conf.CommonConf.AppName,
		AppOwner:     conf.CommonConf.AppOwner,
		ContextPath:  conf.HttpConf.ContextPath,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
		Beans:        make([]interface{}, 0),
		Conf:         conf,
	}
}

func (b *Bootstrapper) InitMiddleware() *Bootstrapper {
	b.Use(recover.New())
	b.Use(logger.New())
	st := stats.New()
	b.Get(b.ContextPath+"/stats", func(ctx iris.Context) {
		ctx.JSON(st.Data())
	})
	b.Use(iris.FromStd(st.ServeHTTP))
	return b
}

func (b *Bootstrapper) RegistBeans(beans ...interface{}) {
	for _, bean := range beans {
		b.Beans = append(b.Beans, bean)
	}
}

func (b *Bootstrapper) AutoInject() *Bootstrapper {
	arr := make([]*inject.Object, 0)
	for _, bean := range b.Beans {
		switch reflect.TypeOf(bean).Kind() {
		case reflect.Struct:
			arr = append(arr, &inject.Object{Value: &bean})
		case reflect.Ptr:
			arr = append(arr, &inject.Object{Value: bean})
		}
	}
	var g inject.Graph
	g.Provide(arr...)
	g.Populate()
	return b
}

func (b *Bootstrapper) InitRouter() *Bootstrapper {
	for _, bean := range b.Beans {
		r, ok := bean.(Router)
		if ok {
			r.InitRouter(b)
		}
	}
	return b
}

// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	b.RegisterView(iris.HTML(viewsDir, ".html"))
}

// SetupSessions initializes the sessions, optionally.
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

// SetupErrorHandlers prepares the http error handlers
// `(context.StatusCodeNotSuccessful`,  which defaults to < 200 || >= 400 but you can change it).
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

const (
	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = "./public/"
	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	Favicon = "favicon.ico"
)

// Configure accepts configurations and runs them inside the Bootstraper's context.
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	//b.Use(recover.New())
	//b.Use(logger.New())
	b.SetupViews("./views")
	b.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	//b.SetupErrorHandlers()
	// static files
	//b.Favicon(StaticAssets + Favicon)
	//b.StaticWeb(StaticAssets[1:len(StaticAssets)-1], StaticAssets)

	// middleware, after static files

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		for _, bean := range b.Beans {
			if b, ok := bean.(SafeDestroy); ok {
				b.Destroy()
			}
		}
		b.Shutdown(ctx)
	})

	return b
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}

func (b *Bootstrapper) Start(cfgs ...iris.Configurator) {
	httpConf := b.Conf.HttpConf
	b.Run(iris.Addr(fmt.Sprintf("%v:%v", httpConf.Host, httpConf.Port)), cfgs...)
}
