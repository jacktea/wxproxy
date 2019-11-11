package wxmsg

import (
	"encoding/xml"
	"github.com/gomydodo/wxencrypter"
	"github.com/kataras/golog"
	"io"
	"io/ioutil"
)

var log = golog.Default

func DecodeXml(appId string,
	token string,
	encodingAesKey string,
	msgSignature string,
	timestamp string,
	nonce string,
	input []byte) (data []byte, err error) {
	e, err := wxencrypter.NewEncrypter(token, encodingAesKey, appId)
	if err != nil {
		data = nil
	}
	data, err = e.Decrypt(msgSignature, timestamp, nonce, input)
	return
}

func CreateEncrypter(appId string,
	token string,
	encodingAesKey string) (e *wxencrypter.Encrypter, err error) {
	e, err = wxencrypter.NewEncrypter(token, encodingAesKey, appId)
	return
}

func Decrypt(wx *wxencrypter.Encrypter, msgSignature string,
	timestamp string,
	nonce string, in io.Reader) (ret []byte, err error) {
	if err != nil {
		return ret, err
	}
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return ret, err
	}
	log.Debug("原始密文：", string(data))
	ret, err = wx.Decrypt(msgSignature, timestamp, nonce, data)
	if err == nil {
		log.Debug("消息明文：", string(ret))
	}
	return
}

func WXEncode(en *wxencrypter.Encrypter, v interface{}) ([]byte, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	log.Debug("消息明文：", string(data))
	rdata, err := en.Encrypt(data)
	if err != nil {
		return nil, err
	}
	log.Debug("消息密文：", string(rdata))
	return rdata, nil
}
