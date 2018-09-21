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
	"strconv"
)

/**
 * 保存函数定义
 * params	body 			内容
 * params 	contentType		内容类型
 * params 	componentAppid	三方应用Appid(目录)
 * params 	appid		    小程序appid(目录)
 * params 	fName		    文件名
 * params 	force		    强制更新
 */
type saveAction func(body []byte,contentType string,componentAppid,appid,fName string,force bool) error

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
func (this *MimiApiServiceImpl) GetQrCode(componentAppid ,appid ,path string,force bool) *QrCodeResp {
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		return handleError(err)
	}
	reqUrl := GET_MINI_QRCODE + "?access_token=" + token
	if path != "" {
		reqUrl += "&path="+url.QueryEscape(path)
	}
	if !WXConf.OssConf.Enabled {
		return localSave(reqUrl,componentAppid,appid,force)
	}else {
		return ossSave(reqUrl,componentAppid,appid,force)
	}
	//return upload(OssDefaultBucket,reqUrl,componentAppid,appid)
}

//获取小程序二维码
func (this *MimiApiServiceImpl) GetWxACode(componentAppid ,appid ,path string,force bool) *QrCodeResp {
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		return handleError(err)
	}
	reqUrl := GET_MINI_WXACODE + "?access_token=" + token
	params := make(map[string]interface{},0)
	if path != "" {
		params["path"] = path
	}
	if !WXConf.OssConf.Enabled {
		return localPostSave(reqUrl,params,componentAppid,appid,force)
	}else {
		return ossPostSave(reqUrl,params,componentAppid,appid,force)
	}
}

//获取小程序二维码
func (this *MimiApiServiceImpl) GetWxACodeUnlimit(componentAppid ,appid ,page,scene string,force bool) *QrCodeResp {
	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
	if err != nil {
		return handleError(err)
	}
	reqUrl := GET_MINI_WXACODEUNLIMIT + "?access_token=" + token
	params := make(map[string]interface{},0)
	if scene == "" {
		scene = string(utils.Random(10))
	}
	params["scene"] = scene
	if page != "" {
		params["page"] = page
	}
	if !WXConf.OssConf.Enabled {
		return localPostSave(reqUrl,params,componentAppid,appid,force)
	}else {
		return ossPostSave(reqUrl,params,componentAppid,appid,force)
	}
}

//获取小程序二维码
//func (this *MimiApiServiceImpl) GetWxQrCode(componentAppid ,appid ,page,scene string,force bool) *QrCodeResp {
//	token,err := this.Svr.GetAppAccessToken(componentAppid,appid)
//	if err != nil {
//		return handleError(err)
//	}
//	reqUrl := GET_MINI_WXQRCODE + "?access_token=" + token
//	params := make(map[string]interface{},0)
//	if scene == "" {
//		scene = string(utils.Random(10))
//	}
//	params["scene"] = scene
//	if page != "" {
//		params["page"] = page
//	}
//	if !WXConf.OssConf.Enabled {
//		return localPostSave(reqUrl,params,componentAppid,appid,force)
//	}else {
//		return ossPostSave(reqUrl,params,componentAppid,appid,force)
//	}
//}

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

//检测文件是否存在
func checkFileExists(componentAppid,appid,fName string,isOss bool) bool {
	if isOss {
		objectName := getObjectName(componentAppid,appid,fName)
		client := DefaultOssBucket()
		ok,_ := client.IsObjectExist(objectName)
		return ok
	}else {
		path := localFilePath(componentAppid,appid,fName)
		return utils.FileExists(path)
	}
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

func postSave(url string, params map[string]interface{}, componentAppid, appid string,force bool, action saveAction) *QrCodeResp {
	jsonBytes , err := json.Marshal(params)
	if err != nil {
		return handleError(err)
	}
	fName := utils.Md5(string(jsonBytes))

	//如果文件存在并且不是强制刷新，则直接返回
	if checkFileExists(componentAppid,appid,fName,WXConf.OssConf.Enabled) && !force {
		ret := new(QrCodeResp)
		ret.Errcode = 0
		ret.Errmsg = "ok"
		ret.Url = fName
		return ret
	}

	header,body,err := utils.HttpPostProxy(url,"application/json",jsonBytes)
	if err != nil {
		return handleError(err)
	}

	contentType := header.Get("Content-Type")
	if strings.HasPrefix(contentType,"image/") {
		idx := strings.Index(contentType,";")
		if idx != -1 {
			contentType = contentType[:idx]
		}
		err := action(body,contentType,componentAppid,appid,fName,force)
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

//保存图片
func save(url, componentAppid, appid string,force bool, action saveAction) *QrCodeResp {
	fName := strconv.FormatInt(time.Now().UnixNano()/1e6,10)
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
		err := action(body,contentType,componentAppid,appid,fName,force)
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
func ossSave(url, componentAppid, appid string,force bool) *QrCodeResp {
	return save(url,componentAppid,appid,force, ossSaveAction)
}

//保存在OSS
func ossPostSave(url string,params map[string]interface{}, componentAppid, appid string,force bool) *QrCodeResp {
	return postSave(url,params,componentAppid,appid,force, ossSaveAction)
}

//本地保存
func localSave(url, componentAppid, appid string,force bool) *QrCodeResp {
	return save(url,componentAppid,appid,force, localSaveAction)
}

//本地保存
func localPostSave(url string,params map[string]interface{}, componentAppid, appid string,force bool) *QrCodeResp {
	return postSave(url,params,componentAppid,appid,force,localSaveAction)
}

func ossSaveAction(body []byte, contentType string, componentAppid, appid, fName string,force bool) error {
	objectName := getObjectName(componentAppid,appid,fName)
	client := DefaultOssBucket()
	ok,_ := client.IsObjectExist(objectName)
	if !ok {
		return client.PutObject(objectName,bytes.NewReader(body),oss.ContentType(contentType))
	}else if force{
		client.DeleteObject(objectName)
		return client.PutObject(objectName,bytes.NewReader(body),oss.ContentType(contentType))
	}
	return nil
}

func localSaveAction(body []byte, contentType string, componentAppid, appid, fName string,force bool) error {
	path := localFilePath(componentAppid,appid,fName)
	if utils.FileExists(path) && !force {
		return nil
	}
	dir := filepath.Dir(path)
	if !utils.FileExists(dir) {
		if err := os.MkdirAll(dir,0751) ; err != nil {
			return err
		}
	}
	//pthSep := string(os.PathSeparator)
	//now := time.Now().UnixNano()/1e6
	//filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
	//	if info == nil {
	//		return err
	//	}
	//	if !info.IsDir() {
	//		//删除1h之前的文件
	//		if v,err := strconv.ParseInt(info.Name(), 10, 64);err == nil && (now-v > 3600000) {
	//			log.Error(os.Remove(dir+pthSep+info.Name()))
	//		}
	//	}
	//	return nil
	//})
	return ioutil.WriteFile(path,body,0660)
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
