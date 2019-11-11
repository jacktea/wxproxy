package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kataras/golog"
	"io/ioutil"
	"net/http"
	"time"
)

var log = golog.Default

func HttpGetBody(uri string) (retBody []byte, err error) {
	var (
		resp *http.Response
	)
	start := time.Now()
	if resp, err = http.Get(uri); err != nil {
		log.Errorf("请求:%s , 错误:%s , 耗时:%d ns", uri, err, time.Now().Sub(start))
		return
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != 200 {
		err = errors.New("http status not ok." + resp.Status)
		return
	}
	retBody, err = ioutil.ReadAll(resp.Body)
	log.Debugf("请求:%s , 响应:%s , 耗时:%d ns", uri, string(retBody), time.Now().Sub(start))
	return
}

func HttpGetJson(uri string, ret interface{}) (err error) {
	var retBody []byte
	if retBody, err = HttpGetBody(uri); err == nil {
		err = json.Unmarshal(retBody, ret)
	}
	return
}

func HttpPostBody(uri string, contextType string, bodyData []byte) (retBody []byte, err error) {
	var (
		resp *http.Response
	)
	start := time.Now()
	if resp, err = http.Post(uri, contextType, bytes.NewBuffer(bodyData)); err != nil {
		log.Errorf("请求:%s , 参数:%s , 错误:%s , 耗时:%d ns", uri, string(bodyData), err, time.Now().Sub(start))
		return
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != 200 {
		err = errors.New("http status not ok." + resp.Status)
		return
	}
	retBody, err = ioutil.ReadAll(resp.Body)
	log.Debugf("请求:%s , 参数:%s , 响应:%s , 耗时:%d ns", uri, string(bodyData), string(retBody), time.Now().Sub(start))
	return
}

func HttpPostJson(uri string, params map[string]interface{}, ret interface{}) error {
	if jsonBytes, err := json.Marshal(params); err != nil {
		return err
	} else {
		return HttpPostBodyJson(uri, jsonBytes, ret)
	}

	//if err != nil {
	//	return
	//}
	//if retBody, err = HttpPostBody(uri,"application/json",jsonBytes) ; err==nil {
	//	err = json.Unmarshal(retBody, ret)
	//}
	//return
}

func HttpPostBodyJson(uri string, bodyData []byte, ret interface{}) error {
	if retBody, err := HttpPostBody(uri, "application/json", bodyData); err != nil {
		return err
	} else {
		return json.Unmarshal(retBody, ret)
	}
}

func HttpGetProxy(uri string) (header http.Header, retBody []byte, err error) {
	var (
		resp *http.Response
	)
	start := time.Now()
	if resp, err = http.Get(uri); err != nil {
		log.Errorf("请求:%s , 错误:%s , 耗时:%d ns", uri, err, time.Now().Sub(start))
		return
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != 200 {
		err = errors.New("http status not ok." + resp.Status)
		return
	}
	header = resp.Header
	retBody, err = ioutil.ReadAll(resp.Body)
	log.Debugf("请求:%s , 响应:%s , 耗时:%d ns", uri, string(retBody), time.Now().Sub(start))
	return
}

func HttpPostProxy(uri string, contextType string, bodyData []byte) (header http.Header, retBody []byte, err error) {
	var (
		resp *http.Response
	)
	start := time.Now()
	if resp, err = http.Post(uri, contextType, bytes.NewBuffer(bodyData)); err != nil {
		log.Errorf("请求:%s , 参数:%s , 错误:%s , 耗时:%d ns", uri, string(bodyData), err, time.Now().Sub(start))
		return
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != 200 {
		err = errors.New("http status not ok." + resp.Status)
		return
	}
	retBody, err = ioutil.ReadAll(resp.Body)
	header = resp.Header
	log.Debugf("请求:%s , 参数:%s , 响应:%s , 耗时:%d ns", uri, string(bodyData), string(retBody), time.Now().Sub(start))
	return
}

func HttpPostRequestBody(uri string, contextType string, req *http.Request) (header http.Header, retBody []byte, err error) {
	reqData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	return HttpPostProxy(uri, contextType, reqData)
}

/*
func HttpPostJson(uri string, bodyData []byte) (retJson map[string]interface{}, err error) {
	var (
		retBody []byte
	)
	if retBody, err = HttpPostBody(uri, bodyData); err == nil {
		err = json.Unmarshal(retBody, &retJson)
	}
	return
}
*/
