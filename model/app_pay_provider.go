package model

type AppPayProvider struct {
	Id int64 `xorm:"pk autoincr"`
	Appid 		string		//公众号ID
	MuchId 		string		//商户号
	Key			string		//API秘钥
	JsDomain	string		//JSAPI支付授权目录(公众号支付)
	ScanNotify	string		//扫码回调链接
	H5Domain	string		//H5支付域名
	IsProvider	bool		//是否是服务商
}

type AppPayClient struct {
	Id int64 `xorm:"pk autoincr"`
	ProviderId 	int64
	SubAppid	string		//子商户公众号
	SubMuchId	string		//子商户ID
	JsDomain	string		//JSAPI支付授权目录(公众号支付)
	ScanNotify	string		//扫码回调链接
}

