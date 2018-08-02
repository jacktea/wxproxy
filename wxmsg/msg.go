package wxmsg

import (
	"encoding/xml"
)

type AppMsg struct {
	XMLName      					xml.Name `xml:"xml"`
	AppId 							string	//第三方平台appid
	CreateTime 						int64	//时间戳
	InfoType 						string	//component_verify_ticket,unauthorized是取消授权，updateauthorized是更新授权，authorized是授权成功通知
	ComponentVerifyTicket 			string	//Ticket内容
	AuthorizerAppid					string	//公众号或小程序
	AuthorizationCode				string 	//可用于换取公众号的接口调用凭据，详细见上面的说明
	AuthorizationCodeExpiredTime	string 	//授权码过期时间
	PreAuthCode						string	//预授权码
}

type BaseMsg struct {
	XMLName      					xml.Name `xml:"xml" json:"-"`
	ToUserName 						string	//开发者微信号
	FromUserName					string	//发送方帐号（一个OpenID）
	CreateTime						int64	//消息创建时间 （整型）
	MsgType							CDATA	//消息类型
}

type WxMsg struct {
	BaseMsg

	MsgId			int64	`xml:"MsgId,omitempty" json:"MsgId,omitempty"`			// 消息id，64位整型
	Content			string	`xml:"Content,omitempty" json:"Content,omitempty"`		// MsgType:text 文本消息内容
	MediaId			string	`xml:"MediaId,omitempty" json:"MediaId,omitempty"`		// MsgType:image,voice,video,shortvideo 消息媒体id
	PicUrl			string	`xml:"PicUrl,omitempty" json:"PicUrl,omitempty"`		// MsgType:image 图片链接（由系统生成）
	Format			string	`xml:"Format,omitempty" json:"Format,omitempty"`		// MsgType:voice 语音格式，如amr，speex等
	Recognition		string 	`xml:"Recognition,omitempty" json:"Recognition,omitempty"`	// MsgType:voice 语音识别结果，UTF8编码
	ThumbMediaId	string	`xml:"ThumbMediaId,omitempty" json:"ThumbMediaId,omitempty"`	// MsgType:video,shortvideo 视频消息缩略图的媒体id
	LocationX		float64	`xml:"Location_X,omitempty" json:"Location_X,omitempty"`	// MsgType:location 地理位置维度
	LocationY		float64	`xml:"Location_Y,omitempty" json:"Location_Y,omitempty"`	// MsgType:location 地理位置经度
	Scale			float64	`xml:"Scale,omitempty" json:"Scale,omitempty"`			// MsgType:location 地图缩放大小
	Label			string	`xml:"Label,omitempty" json:"Label,omitempty"`			// MsgType:location 地理位置信息
	Title			string	`xml:"Title,omitempty" json:"Title,omitempty" `			// MsgType:link 消息标题
	Description		string	`xml:"Description,omitempty" json:"Description,omitempty"`	// MsgType:link 消息描述
	Url				string	`xml:"Url,omitempty" json:"Url,omitempty"`			// MsgType:link 消息链接
	//事件消息
	Event			string	`xml:"Event,omitempty" json:"Event,omitempty"`			// MsgType:event 事件类型 subscribe,unsubscribe,SCAN,LOCATION,CLICK,VIEW
	EventKey		string	`xml:"EventKey,omitempty" json:"EventKey,omitempty"`		// MsgType:event 事件KEY值
	Ticket			string	`xml:"Ticket,omitempty" json:"Ticket,omitempty"`		// MsgType:event 二维码的ticket，可用来换取二维码图片
	Latitude		float64	`xml:"Latitude,omitempty" json:"Latitude,omitempty"`		// MsgType:event 地理位置纬度
	Longitude		float64	`xml:"Longitude,omitempty" json:"Longitude,omitempty"`		// MsgType:event 地理位置经度
	Precision		float64	`xml:"Precision,omitempty" json:"Precision,omitempty"`		// MsgType:event 地理位置精度
	Status 			string	`xml:"Status,omitempty" json:"Status,omitempty"`		// MsgType:event 发送模板消息状态
}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

type ReplyMsg struct {
	BaseMsg
	//文本消息
	Content			CDATA	`xml:"Content,omitempty"`
	//图片消息
	Image			*struct{
		MediaId		CDATA
	}						`xml:"Image,omitempty"`
	//语音消息
	Voice			*struct{
		MediaId		CDATA
	}						`xml:"Voice,omitempty"`
	//视频消息
	Video			*struct{
		MediaId		CDATA
		Title		CDATA
		Description	CDATA
	}						`xml:"Video,omitempty"`
	//音乐消息
	Music			*struct{
		ThumbMediaId	CDATA
		Title			CDATA
		Description		CDATA
		MusicUrl		CDATA
		HQMusicUrl		CDATA

	}						`xml:"Music,omitempty"`
	//图文消息
	ArticleCount	int		`xml:"ArticleCount,omitempty"`
	//Articles		struct{
	//	Item		[]struct{
	//		Title		CDATA
	//		Description	CDATA
	//		PicUrl		CDATA
	//		Url			CDATA
	//	}
	//}						`xml:"Articles,omitempty"`

}

func NewTextReplyMsg(content string) *ReplyMsg {
	msg := new(ReplyMsg)
	msg.MsgType = "text"
	msg.Content = CDATA(content)
	return msg
}

//客服消息
type CustomMsg struct {
	Touser			string			`json:"touser"`	//接收方OPENID
	Msgtype			string			`json:"msgtype"`//消息类型
}
//文本消息
type CustomTextMsg struct {
	CustomMsg
	Text			CustomTextItem	`json:"text"`
}

type CustomTextItem struct {
	Content			string			`json:"content"`//内容
}

func NewCustomTextMsg(content string) *CustomTextMsg  {
	msg := new(CustomTextMsg)
	msg.Msgtype	= "text"
	msg.Text.Content = content
	return msg
}

//图片消息
type CustomImageMsg struct {
	CustomMsg
	Image			CustomMediaItem	`json:"image"`
}
//语音消息
type CustomVoiceMsg struct {
	CustomMsg
	Voice			CustomMediaItem	`json:"voice"`
}

type CustomMediaItem struct {
	MediaId			string			`json:"media_id"`
}

//视频消息
type CustomVideoMsg struct {
	CustomMsg
	Video			CustomVideoItem	`json:"video"`
}

//
type CustomMpnewsMsg struct {
	CustomMsg
	Mpnews			CustomMediaItem	`json:"mpnews"`
}

//视频消息
type CustomVideoItem struct {
	MediaId			string			`json:"media_id"`
	ThumbMediaId	string			`json:"thumb_media_id"`
	Title			string			`json:"title"`
	Description		string			`json:"description"`
}

//音乐消息
type CustomMusicMsg struct {
	CustomMsg
	Music			struct{
		Title			string			`json:"title"`
		Description		string			`json:"description"`
		ThumbMediaId	string			`json:"thumb_media_id"`
		Musicurl		string			`json:"musicurl"`
		Hqmusicurl		string			`json:"hqmusicurl"`
	}	`json:"music"`
}

/*
type CustomMusicMsg struct {
	CustomMsg
	Music			CustomMusicItem	`json:"music"`
}

type CustomMusicItem struct {
	Title			string			`json:"title"`
	Description		string			`json:"description"`
	ThumbMediaId	string			`json:"thumb_media_id"`
	Musicurl		string			`json:"musicurl"`
	Hqmusicurl		string			`json:"hqmusicurl"`
}
*/

//图文消息
type CustomNewsMsg struct {
	CustomMsg
	News			struct{
		Articles	[]struct{
			Title			string			`json:"title"`
			Description		string			`json:"description"`
			Url				string			`json:"url"`
			Picurl			string			`json:"picurl"`
		}	`json:"articles"`
	}	`json:"news"`
}
/*
type CustomNewsMsg struct {
	CustomMsg
	News			CustomNewsData	`json:"news"`
}

type CustomNewsData struct {
	Articles		[]CustomArticleItem	`json:"articles"`
}

type CustomArticleItem struct {
	Title			string			`json:"title"`
	Description		string			`json:"description"`
	Url				string			`json:"url"`
	Picurl			string			`json:"picurl"`
}
*/

//卡券
type CustomWxcardMsg struct {
	CustomMsg
	Wxcard			struct{
		CardId			string			`json:"card_id"`
	}	`json:"wxcard"`
}

/*
type CustomWxcardMsg struct {
	CustomMsg
	Wxcard			CustomCardItem	`json:"wxcard"`
}

type CustomCardItem struct {
	CardId			string			`json:"card_id"`
}
*/
