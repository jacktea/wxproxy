package notify

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/jacktea/wxproxy/bootstrap"
	. "github.com/jacktea/wxproxy/config"
	"github.com/jacktea/wxproxy/service"
	"github.com/jacktea/wxproxy/service/api"
	. "github.com/jacktea/wxproxy/wxmsg"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//var httpNotify = h.NewHttpNotify()
var log = golog.Default

type NotifyAction struct {
	Svr api.ApiService `inject:""`
}

func (this *NotifyAction) InitRouter(app *bootstrap.Bootstrapper) {
	app.Post(WXConf.HttpConf.ContextPath+"/{componentAppid}/notify", this.Notify)
	app.Post(WXConf.HttpConf.ContextPath+"/{componentAppid}/{appid}/notify", this.NotifyApp)

	party := app.Party(WXConf.HttpConf.ContextPath + "/notify")
	party.Post("/global/{componentAppid}", this.Notify)
	party.Post("/app/{componentAppid}/{appid}", this.NotifyApp)
	party.Post("/app/{appid}", this.NotifyApp)
	log.Info("init notify router...")
}

func (this *NotifyAction) Notify(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid")
	msgSignature := c.FormValue("msg_signature")
	timestamp := c.FormValue("timestamp")
	nonce := c.FormValue("nonce")
	log.Debug(componentAppid, msgSignature, timestamp, nonce)

	encypter, err := this.Svr.CacheGetWxEncrypter(componentAppid)
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	data, err := Decrypt(encypter, msgSignature, timestamp, nonce, c.Request().Body)
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	oMsg := new(AppMsg)
	if err = xml.Unmarshal(data, oMsg); err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	switch oMsg.InfoType {
	case "component_verify_ticket": //票据通知
		go this.Svr.UpdateTicket(componentAppid, oMsg.ComponentVerifyTicket)
	case "authorized": //授权成功通知
		go this.Svr.DoAuthorize(componentAppid, oMsg.AuthorizerAppid, oMsg.AuthorizationCode)
	case "updateauthorized": //更新授权通知
		go this.Svr.DoAuthorize(componentAppid, oMsg.AuthorizerAppid, oMsg.AuthorizationCode)
	case "unauthorized": //取消授权通知
		go this.Svr.AuthorizedCancel(componentAppid, oMsg.AuthorizerAppid)
	}
	go this.notifyApp(componentAppid, oMsg.AuthorizerAppid, oMsg)

	c.WriteString("success")
}

func (this *NotifyAction) notifyApp(componentAppid, appid string, msg interface{}) {
	if appid == "" {
		return
	}
	if info, b := this.Svr.CacheFindAuthorizationAppInfo(componentAppid, appid); b {
		if data, err := json.Marshal(msg); err == nil {
			postJson(info.NotifyUrl, data)
			if info.Mode == 2 {
				go postJson(info.DebugNotifyUrl, data)
			}
		}
	}
}

func (this *NotifyAction) NotifyApp(c iris.Context) {
	componentAppid := c.Params().Get("componentAppid") //第三方应用APPID
	appid := c.Params().Get("appid")                   //授权方APPID
	msgSignature := c.FormValue("msg_signature")       //消息签名
	timestamp := c.FormValue("timestamp")              //时间戳
	nonce := c.FormValue("nonce")                      //随机码

	encypter, err := this.Svr.CacheGetWxEncrypter(componentAppid)
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	//解密消息报文
	ddata, err := Decrypt(encypter, msgSignature, timestamp, nonce, c.Request().Body)
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	wxMsg := new(WxMsg)
	if err = xml.Unmarshal(ddata, wxMsg); err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}

	if "wx570bc396a51b8ff8" == appid || "wxd101a85aa106f53e" == appid {
		//全网发布
		//return deploy(componentAppid,appid,ddata,wxMsg.MsgType,c)
		this.deploy(componentAppid, appid, wxMsg, c)
		return
	}

	go (func(componentAppid string, appid string, wxMsg *WxMsg) {
		if info, b := this.Svr.CacheFindAuthorizationAppInfo(componentAppid, appid); b {
			if data, err := json.Marshal(wxMsg); err == nil {

				//多次推送成功(返回success)为止,间隔15/15/30/180/1800/1800/1800/1800/3600s
				//httpNotify.AddTask(info.NotifyUrl,data)

				/*
					res,err := http.Post(info.NotifyUrl,"application/json",bytes.NewBuffer(data))
					if err != nil {
						Logger.Errorf("请求: %v 失败，请求内容: %v，错误信息:%v",info.NotifyUrl,string(data),err)
					}else {
						d,err := ioutil.ReadAll(res.Body)
						if err != nil {
							Logger.Errorf("请求: %v 成功，请求内容: %v，解析响应数据失败，错误信息:%v",info.NotifyUrl,string(data),err)
						}else {
							Logger.Infof("请求: %v 成功，请求内容:%v，响应信息: %v",info.NotifyUrl,string(data),string(d))
						}
					}
				*/

				postJson(info.NotifyUrl, data)

				if info.Mode == 2 {
					go postJson(info.DebugNotifyUrl, data)
				}

			}
		}
	})(componentAppid, appid, wxMsg)

	/*
		switch wxMsg.MsgType {
		case "event":
		case "text":
		case "image":
		case "voice":
		case "video","shortvideo":
		case "location":
		case "link":
		}
	*/
	c.WriteString("success")
}

func postJson(url string, data []byte) {
	if url == "" {
		log.Errorf("请求: %v 失败，请求内容: %v，错误信息:%v", url, string(data), "请求地址为空")
		return
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Errorf("请求: %v 失败，请求内容: %v，错误信息:%v", url, string(data), err)
	} else {
		d, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Errorf("请求: %v 成功，请求内容: %v，解析响应数据失败，错误信息:%v", url, string(data), err)
		} else {
			log.Infof("请求: %v 成功，请求内容:%v，响应信息: %v", url, string(data), string(d))
		}
	}
}

func (this *NotifyAction) deploy(componentAppid string,
	appid string,
	wxMsg *WxMsg,
	c iris.Context) {
	wx, err := this.Svr.CacheGetWxEncrypter(componentAppid) //微信加密解密器
	if err != nil {
		log.Error(err)
		c.JSON(service.NewServerErrorResp(err))
		return
	}
	//replyMsg := new(msg.ReplyMsg)
	switch wxMsg.MsgType {
	case "event":
		replyMsg := NewTextReplyMsg(wxMsg.Event + "from_callback")
		//replyMsg.MsgType 			= "text"
		//replyMsg.Content 			= msg.CDATA(wxMsg.Event+"from_callback")
		replyMsg.ToUserName = wxMsg.FromUserName
		replyMsg.FromUserName = wxMsg.ToUserName
		replyMsg.CreateTime = time.Now().Unix()
		edata, err := WXEncode(wx, replyMsg)
		if err != nil {
			c.JSON(service.NewServerErrorResp(err))
			return
		}
		c.Write(edata)
		return
	case "text":
		if wxMsg.Content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
			//replyMsg.MsgType		= "text"
			//replyMsg.Content 		= msg.CDATA(wxMsg.Content+"_callback")
			replyMsg := NewTextReplyMsg(wxMsg.Content + "_callback")
			replyMsg.ToUserName = wxMsg.FromUserName
			replyMsg.FromUserName = wxMsg.ToUserName
			replyMsg.CreateTime = time.Now().Unix()
			edata, err := WXEncode(wx, replyMsg)
			if err != nil {
				c.JSON(service.NewServerErrorResp(err))
				return
			}
			c.Write(edata)
			return
		} else if strings.HasPrefix(wxMsg.Content, "QUERY_AUTH_CODE") {
			s := strings.Split(wxMsg.Content, ":")
			cMsg := NewCustomTextMsg(s[1] + "_from_api")
			cMsg.Touser = wxMsg.FromUserName
			//cMsg.Msgtype = "text"
			//cMsg.Text.Content = s[1]+"_from_api"
			go this.Svr.SendCustomMsg(componentAppid, appid, cMsg)
			c.WriteString("")
			return
		}
	}
	c.WriteString("success")
}

/*
func deploy(componentAppid string,
		appid string,
		ddata []byte,
		msgType string,
		c echo.Context) error  {
	wx,err			:= service.GetWxEncrypter(componentAppid)	//微信加密解密器
	if err != nil{
		Logger.Error(err)
		return err
	}
	switch msgType {
	case "event":
		oMsg := new(msg.EventMsg)
		if err := xml.Unmarshal(ddata,oMsg);err != nil{
			Logger.Error(err)
			return err
		}
		textMsg 				:= new(msg.TextMsg)
		textMsg.MsgType			= "text"
		textMsg.Content 		= oMsg.Event+"from_callback"
		textMsg.ToUserName 		= oMsg.FromUserName
		textMsg.FromUserName 	= oMsg.ToUserName
		textMsg.CreateTime 		= strconv.FormatInt(time.Now().Unix(),10)
		edata,err 				:= msg.WXEncode(wx,textMsg)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK,string(edata))
	case "text":
		oMsg := new(msg.TextMsg)
		if err := xml.Unmarshal(ddata,oMsg);err!=nil{
			return err
		}
		if oMsg.Content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
			textMsg := new(msg.TextMsg)
			textMsg.MsgType			= "text"
			textMsg.Content 		= oMsg.Content+"_callback"
			textMsg.ToUserName 		= oMsg.FromUserName
			textMsg.FromUserName 	= oMsg.ToUserName
			textMsg.CreateTime 		= strconv.FormatInt(time.Now().Unix(),10)
			edata,err := msg.WXEncode(wx,textMsg)
			if err != nil {
				return err
			}
			return c.String(http.StatusOK,string(edata))
		}else if strings.HasPrefix(oMsg.Content,"QUERY_AUTH_CODE"){
			s := strings.Split(oMsg.Content,":")
			cMsg := msg.CustomTextMsg{
			}
			cMsg.Touser = oMsg.FromUserName
			cMsg.Msgtype = "text"
			cMsg.Text.Content = s[1]+"_from_api"
			go service.SendCustomMsg(componentAppid,appid,cMsg)
			return c.String(http.StatusOK,"")
		}
	}
	return c.String(http.StatusOK,"success")
}
*/
