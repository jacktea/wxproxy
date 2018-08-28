package miniprogram

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"time"
	. "github.com/jacktea/wxproxy/common"
	. "github.com/jacktea/wxproxy/config"
	"net/url"
	"io"
	"github.com/jacktea/wxproxy/utils"
	"bytes"
	"os"
	"path/filepath"
)

type saveAction func(body []byte,contentType string,componentAppid,appid,fName string) error

//var ossConf = WXConf.OssConf

//var endpoint = ossConf.Endpoint
//var accessKeyId = ossConf.AccessKeyId
//var accessKeySecret = ossConf.AccessKeySecret
//var bucketName = ossConf.BucketName
//var ossEnabled = ossConf.Enabled

var ossDefaultClient *oss.Client
var ossDefaultBucket *oss.Bucket

func NewOssClient() *oss.Client {
	client, err := oss.New(WXConf.OssConf.Endpoint, WXConf.OssConf.AccessKeyId, WXConf.OssConf.AccessKeySecret)
	if err != nil {
		log.Error(err)
		return nil
	}
	return client
}

func DefaultOssClient() *oss.Client {
	if ossDefaultClient == nil {
		ossDefaultClient = NewOssClient()
	}
	return ossDefaultClient
}

func DefaultOssBucket() *oss.Bucket {
	if ossDefaultBucket == nil {
		bucket, err := DefaultOssClient().Bucket(WXConf.OssConf.BucketName)
		if err != nil {
			log.Error(err)
			return nil
		}
		ossDefaultBucket = bucket
	}
	return ossDefaultBucket
}

func handleError(err error) *QrCodeResp {
	ret := new(QrCodeResp)
	ret.Errcode = -1
	ret.Errmsg = err.Error()
	return ret
}

//获取预览二维码
func (this *MimiApiServiceImpl) GetQrCode(componentAppid ,appid ,path string) *QrCodeResp {
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		return handleError(err)
	}
	reqUrl := GET_MINI_QRCODE + "?access_token=" + token
	if path != "" {
		reqUrl += "&path="+url.QueryEscape(path)
	}
	if !WXConf.OssConf.Enabled {
		return localSave(reqUrl,componentAppid,appid)
	}else {
		return ossSave(reqUrl,componentAppid,appid)
	}
	//return upload(OssDefaultBucket,reqUrl,componentAppid,appid)
}

//查看预览二维码图片
func (this *MimiApiServiceImpl) DownLoadQrCode(componentAppid,appid,fName string) (io.ReadCloser,error){
	if !WXConf.OssConf.Enabled {
		return localLoad(componentAppid,appid,fName)
	}else {
		return ossLoad(componentAppid,appid,fName)
	}
}

func (this *MimiApiServiceImpl) GetObjectName(componentAppid,appid,fName string) string {
	return getObjectName(componentAppid,appid,fName)
}

//OSS上存放的文件名
func getObjectName(componentAppid,appid,fName string) string {
	objectName := fmt.Sprintf("miniprogram/%v/%v/%v",componentAppid,appid,fName)
	log.Debug("oss objectName is ",objectName)
	return objectName
}

//本地文件存放路径
func localFilePath(componentAppid, appid, fName string) string {
	file := fmt.Sprintf("%v/miniprogram/%v/%v/%v",WXConf.CommonConf.DataPath,componentAppid,appid,fName)
	log.Debug("local file is ",file)
	return file
}

//从本地加载
func localLoad(componentAppid, appid, fName string) (io.ReadCloser,error) {
	objectName := localFilePath(componentAppid,appid,fName)
	return os.Open(objectName)
}

//从OSS服务器加载
func ossLoad(componentAppid, appid, fName string) (io.ReadCloser,error) {
	objectName := getObjectName(componentAppid,appid,fName)
	return DefaultOssBucket().GetObject(objectName)
}

//保存图片
func save(url, componentAppid, appid string,action saveAction) *QrCodeResp {
	fName := time.Now().Format("20060102030405")
	header,body,err := utils.HttpGetProxy(url)
	if err != nil {
		return handleError(err)
	}
	contentType := header.Get("Content-Type")
	if strings.HasPrefix(contentType,"image/") {
		idx := strings.Index(contentType,";")
		if idx != -1 {
			contentType = contentType[:idx]
		}
		err := action(body,contentType,componentAppid,appid,fName)
		if err != nil {
			return handleError(err)
		}else{
			ret := new(QrCodeResp)
			ret.Errcode = 0
			ret.Errmsg = "ok"
			ret.Url = fName
			return ret
		}
	}else {
		ret := new(QrCodeResp)
		if err := json.Unmarshal(body,ret);err != nil {
			return handleError(err)
		}else {
			return ret
		}
	}
}

//保存在OSS
func ossSave(url, componentAppid, appid string) *QrCodeResp {
	return save(url,componentAppid,appid, func(body []byte, contentType string, componentAppid,appid,fName string) error {
		objectName := getObjectName(componentAppid,appid,fName)
		return DefaultOssBucket().PutObject(objectName,bytes.NewReader(body),oss.ContentType(contentType))
	})
}

//本地保存
func localSave(url, componentAppid, appid string) *QrCodeResp {
	return save(url,componentAppid,appid, func(body []byte, contentType string, componentAppid,appid,fName string) error {
		path := localFilePath(componentAppid,appid,fName)
		dir := filepath.Dir(path)
		if !utils.FileExists(dir) {
			if err := os.MkdirAll(dir,0751) ; err != nil {
				return err
			}
		}
		pthSep := string(os.PathSeparator)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if !info.IsDir() {
				log.Error(os.Remove(dir+pthSep+info.Name()))
			}
			return nil
		})
		return ioutil.WriteFile(path,body,0660)
	})
}

//func upload(bucket *oss.Bucket,url,componentAppid ,appid string) *QrCodeResp {
//	fName := time.Now().Format("20060102030405")
//	objectName := getObjectName(componentAppid,appid,fName)
//	resp,err := http.Get(url)
//	if err != nil {
//		return handleError(err)
//	}
//	defer resp.Body.Close()
//	contentType := resp.Header.Get("Content-Type")
//	if strings.HasPrefix(contentType,"image/") {
//		idx := strings.Index(contentType,";")
//		if idx != -1 {
//			contentType = contentType[:idx]
//		}
//		err := bucket.PutObject(objectName,resp.Body,oss.ContentType(contentType))
//		if err != nil {
//			return handleError(err)
//		}else{
//			ret := new(QrCodeResp)
//			ret.Errcode = 0
//			ret.Errmsg = "ok"
//			ret.Url = fName
//			return ret
//		}
//	}else {
//
//		b,err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			return handleError(err)
//		}else {
//			ret := new(QrCodeResp)
//			if err := json.Unmarshal(b,ret);err != nil {
//				return handleError(err)
//			}else {
//				return ret
//			}
//		}
//	}
//}
