# 微信代理

## 1. 配置

### 1.1 基础配置

```yaml
etcd:						#etcd配置,不必须，默认为内嵌模式
  embed: true				#是否内嵌etcd(单机模式可内嵌)
  client_port: 2379			#内嵌时需要，可选
  peer_port: 2380			#内嵌时需要，可选
  name: embedEtcd			#内嵌时需要，可选
  data_dir: default.etcd	#内嵌时需要，可选
  endpoints:				#非内嵌时必须(集群模式)
   - http://127.0.0.1:2379
db:							#数据库配置	
  drive_name: mysql
  host: 127.0.0.1
  port: 3306
  db_name: wxappdb
  charset: utf8
  user_name: jpss
  password: ******
http:						#http配置
  port: 8011
  context_path: /wxproxy
```

### 1.2 日志配置


### 1.3 SQL

```sql
CREATE TABLE `wx_app_base_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_id` varchar(20) COLLATE utf8mb4_bin NOT NULL,
  `app_secret` varchar(50) COLLATE utf8mb4_bin NOT NULL,
  `component_access_token` varchar(150) COLLATE utf8mb4_bin DEFAULT NULL,
  `component_access_token_expire` datetime DEFAULT NULL,
  `component_verify_ticket` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `encoding_aes_key` varchar(50) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `modify_time` datetime DEFAULT NULL,
  `pre_auth_code` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `pre_auth_code_expire` datetime DEFAULT NULL,
  `token` varchar(50) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_blwxuagw31rc8c6257lq4c6ju` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC

CREATE TABLE `wx_authorization_access_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `access_token` varchar(150) COLLATE utf8mb4_bin NOT NULL,
  `access_token_expire` datetime NOT NULL,
  `appid` varchar(20) COLLATE utf8mb4_bin NOT NULL,
  `create_time` datetime DEFAULT NULL,
  `func_info` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL,
  `modify_time` datetime DEFAULT NULL,
  `refresh_token` varchar(150) COLLATE utf8mb4_bin NOT NULL,
  `status` varchar(2) COLLATE utf8mb4_bin DEFAULT NULL,
  `component_appid` varchar(20) COLLATE utf8mb4_bin NOT NULL,
  `authorization_code` varchar(150) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_a1r7l77111tmg8cl661hl81y4` (`appid`,`component_appid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC

CREATE TABLE `wx_authorization_app_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `component_appid` varchar(20) NOT NULL,
  `appid` varchar(20) NOT NULL,
  `notify_url` varchar(250) NOT NULL,
  `mode`  int NOT NULL,
  `debug_notify_url` varchar(250) NOT NULL,
  `create_time` datetime NOT NULL,
  `modify_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_sdfwiefnas234f2` (`component_appid`,`appid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC

CREATE TABLE `wx_authorization_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `alias` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `appid` varchar(20) COLLATE utf8mb4_bin NOT NULL,
  `business_info` varchar(120) COLLATE utf8mb4_bin DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `head_img` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL,
  `modify_time` datetime DEFAULT NULL,
  `nick_name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `principal_name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `qrcode_url` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL,
  `service_type_info` varchar(2) COLLATE utf8mb4_bin DEFAULT NULL,
  `user_name` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `verify_type_info` varchar(2) COLLATE utf8mb4_bin DEFAULT NULL,
  `signature` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_oeyk5rxi35hkhc53fatmr9m4s` (`appid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC
```

## 2. 公众号

### 2.1 公众号授权

__授权地址：__ `http://www.ishanshan.com/wxproxy/authorize/三方应用appid/回调地址base64`

__三方应用appid：__ wx7c44bb2354440737

__回调地址：__ `http://www.ishanshan.com`

__则请求地址：__ `http://www.ishanshan.com/wxproxy/authorize/wx7c44bb2354440737/aHR0cDovL3d3dy5pc2hhbnNoYW4uY29t`

__回调请求(GET)：__ `http://www.ishanshan.com?appid=你的appid&componentAppid=wxe5bd129decc89caf`

### 2.2 管理API

#### 2.2.1 新增/修改第三方应用信息

#### 2.2.2 刷新三方应用token

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upcmptoken/三方应用appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| force | bool | 否    | 强制重新获取 |  

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.3 刷新三方应预授权码

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upcmpcode/三方应用appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| force | bool | 否    | 强制重新获取 |  

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.4 刷新公众号/小程序访问toekn

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upauthtoken/三方应用appid/公众号appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| force | bool | 否    | 强制重新获取 |  

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.5 更新公众号对应应用的消息通知地址

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/account/upappnotifyurl/三方应用appid/公众号appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| notify_url | string | 是    | 公众号消息通知回调地址通知回调地址 |  
| mode | int | 是    | 模式 1普通 2调试，在调试模式下会同时向notify_url和debug_notify_url推送消息 | 
| debug_notify_url | string | 是    | 调试模式下的推送通知回调地址 | 

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.6 创建带参关注二维码


__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/account/createqrcode/三方应用appid/公众号appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| identity | string | 是    | 参数标识 |  
| expire | int | 是    | 过期时间 | 
| forever | bool | 是    | 是否为永久二维码 | 


__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  
| ticket | string | 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。 |  
| expire_seconds | int | 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。 |  
| url | string | 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片 |  

#### 2.2.7 获取公众号相关信息

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/account/getauthappinfo/三方应用appid/公众号appid`

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  
| appid | string | 授权方appid | 
| nickName | string | 授权方昵称 | 
| headImg | string | 授权方头像地址 | 
| serviceTypeInfo | string | 授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号 | 
| verifyTypeInfo | string | 授权方认证类型 | 
| userName | string | 授权方公众号的原始ID | 
| principalName | string | 公众号的主体名称 | 
| alias | string | 授权方公众号所设置的微信号，可能为空 | 
| businessInfo | string | 业务功能开通情况 | 
| qrcodeUrl | string | 二维码图片的URL，开发者最好自行也进行保存 | 
| accessToken | string | 访问token | 
| accessTokenExpire | string | 访问token过期时间 | 
| authorizationStatus | string | 公众号授权状态 0 无效 1 授权成功 2 取消授权 |
| mode | int | 通知模式 1普通 2调试 |
| notifyUrl | string | 回调通知URL |
| debugNotifyUrl | string | 调试模式下通知URL |


### 2.3 微信授权

#### 2.3.1 微信网页授权

##### 2.3.1.1 回调方式为GET

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxauth/apply/wx/三方应用appid/公众号appid?rd=回调地址&scope=(授权方式snsapi_userinfo|snsapi_base)`

##### 2.3.1.2 回调方式为POST

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxauth/apply/wxex/三方应用appid/公众号appid?rd=回调地址&scope=(授权方式snsapi_userinfo|snsapi_base)`

#### 2.3.2 获取用户信息

##### 2.3.2.1 普通方式

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxauth/wx/userinfo?openid=用户openid`

***普通访问方式加上callback方式也会转换为jsonp方式***

##### 2.3.2.2 jsonp方式

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxauth/wx/jp/userinfo?openid=用户openid&callback=回调函数名`

***jsonp访问方式去掉callback方式也会转换为普通方式***

##### 2.3.2.3 老版本普通方式

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxoauth/wx/userinfo?openid=用户openid`

***普通访问方式加上callback方式也会转换为jsonp方式***

##### 2.3.2.4 老版本jsonp方式

__访问地址：__ `http://www.ishanshan.com/wxproxy/wxoauth/wx/jp/userinfo?openid=用户openid&callback=回调函数名`

***jsonp访问方式去掉callback方式也会转换为普通方式***

#### 2.3.3 js验签

##### 2.3.3.1 普通方式

**访问地址：** `http://www.ishanshan.com/wxproxy/wxauth/wx/jsconfig/三方应用appid/公众号appid`

***普通访问方式加上callback方式也会转换为jsonp方式***

##### 2.3.3.1 jsonp方式

**访问地址：** `http://www.ishanshan.com/wxproxy/wxauth/wx/jp/jsconfig/三方应用appid/公众号appid`

***普通访问方式加上callback方式也会转换为jsonp方式***

### 2.4 消息管理

#### 2.4.1 发送客服消息

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/msg/custommsg/三方应用appid/公众号appid`

__请求内容(参见微信公众平台)：__

```json
{
    "touser":"OPENID",
    "msgtype":"text",
    "text":
    {
         "content":"Hello World123<a href=\"http://www.ishanshan.com\">前往</a>"
    }
}
```

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.4.2 发送模板消息

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/msg/tplmsg/三方应用appid/公众号appid`

__请求内容(参见微信公众平台)：__

```json
{
           "touser":"OPENID",
           "template_id":"ngqIpbwh8bUfcSsECmogfXcV14J0tQlEpBO27izEYtY",
           "url":"http://weixin.qq.com/download",  
           "miniprogram":{
             "appid":"xiaochengxuappid12345",
             "pagepath":"index?foo=bar"
           },          
           "data":{
                   "first": {
                       "value":"恭喜你购买成功！",
                       "color":"#173177"
                   },
                   "keyword1":{
                       "value":"巧克力",
                       "color":"#173177"
                   },
                   "keyword2": {
                       "value":"39.8元",
                       "color":"#173177"
                   },
                   "keyword3": {
                       "value":"2014年9月22日",
                       "color":"#173177"
                   },
                   "remark":{
                       "value":"欢迎再次购买！",
                       "color":"#173177"
                   }
           }
       }
```

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.4.3 发送客服文本消息

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/send/customtextmsg/三方应用appid/公众号appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| openid | string | 是    | 用户openid |  
| content | string | 是    | 发送内容 | 

### 2.5 用户管理

#### 2.5.1 获取用户基本信息

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/x-www-form-urlencoded;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/user/info/三方应用appid/公众号appid`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| openid | string | 是    | 用户openid |  

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 | 
| subscribe | int |	用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。|
| openid	| string |用户的标识，对当前公众号唯一|
| nickname | string |	用户的昵称|
| sex | int |	用户的性别，值为1时是男性，值为2时是女性，值为0时是未知|
| city | string |	用户所在城市|
| country | string |	用户所在国家|
| province | string |	用户所在省份|
| language | string |	用户的语言，简体中文为zh_CN|
| headimgurl | string |	用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。|
| subscribe_time | long |	用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间|
| unionid | string |	只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。|
| remark | string |	公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注|
| groupid | int |	用户所在的分组ID（兼容旧的用户分组接口）|
| tagid_list | array |	用户被打上的标签ID列表|
| subscribe_scene | string |	返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他|
| qr_scene | long |	二维码扫码场景（开发者自定义）|
| qr_scene_str | string |	二维码扫码场景描述（开发者自定义）|

## 3. 小程序

### 3.1 修改服务器地址

#### 3.1.1 设置小程序服务器域名

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/domain/modify/三方应用appid/公众号appid`

__请求参数：__

```json
  {
    "action":"add",
    "requestdomain":["https://www.qq.com","https://www.qq.com"],
    "wsrequestdomain":["wss://www.qq.com","wss://www.qq.com"],
    "uploaddomain":["https://www.qq.com","https://www.qq.com"],
    "downloaddomain":["https://www.qq.com","https://www.qq.com"]
  }
```

| 参数名称       | 描述                |
| :---------- | :----------------- |
| action | add添加, delete删除, set覆盖, get获取。当参数是get时不需要填四个域名字段 |  
| requestdomain | request合法域名，当action参数是get时不需要此字段 | 
| wsrequestdomain | socket合法域名，当action参数是get时不需要此字段 | 
| uploaddomain | uploadFile合法域名，当action参数是get时不需要此字段 | 
| downloaddomain | downloadFile合法域名，当action参数是get时不需要此字段 | 

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok",
    "requestdomain":["https://www.qq.com","https://www.qq.com"],
    "wsrequestdomain":["wss://www.qq.com","wss://www.qq.com"],
    "uploaddomain":["https://www.qq.com","https://www.qq.com"],
    "downloaddomain":["https://www.qq.com","https://www.qq.com"]
  }
```

#### 3.1.2 设置小程序业务域名

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/domain/setwebviewdomain/三方应用appid/公众号appid`

__请求参数：__

```json
  {
    "action":"add",
    "webviewdomain":["https://www.qq.com","https://m.qq.com"]
  }
```

| 参数名称       | 描述                |
| :---------- | :----------------- |
| action | add添加, delete删除, set覆盖, get获取。当参数是get时不需要填 |  
| webviewdomain | 小程序业务域名，当action参数是get时不需要此字段 | 

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok"
  }
```

### 3.2 小程序基本信息设置

#### 3.2.1 获取帐号基本信息

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/info/getaccountbasicinfo/三方应用appid/公众号appid`

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok",
    "appid": "wxdc685123d955453", //帐号appid
    "account_type": 2,            //帐号类型（1：订阅号，2：服务号，3：小程序）
    "principal_type": 1,          //主体类型（1：企业）
    "principal_name": "深圳市腾讯计算机系统有限公司",//主体名称
    "realname_status": 1,
    "wx_verify_info": {           //微信认证信息
        "qualification_verify": 1,
        "naming_verify": 1,
        "annual_review": 1,
        "annual_review_begin_time": 1550490981,
        "annual_review_end_time": 1558266981
    },
    "signature_info": {           //功能介绍信息
        "signature": "功能介绍",
        "modify_used_count": 1,
        "modify_quota": 5
    },
    "head_image_info": {          //头像信息
        "head_image_url": "http://mmbiz.qpic.cn/mmbiz/a5icZrUmbV8p5jb6RZ8aYfjfS2AVle8URwBt8QIu6XbGewB9wiaWYWkPwq4R7pfdsFibuLkic16UcxDSNYtB8HnC1Q/0",
        "modify_used_count": 3,
        "modify_quota": 5
    }
  }
```

#### 3.2.2 小程序名称设置及改名

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/info/setnickname/三方应用appid/公众号appid`

__请求参数：__

```json
  {
    "nick_name": "XXX公司",
    "id_card": "12345678-0",
    "license": "广州市新港中路397号TIT创意园",
    "naming_other_stuff_1": "3LaLzqiTrQcD20DlX_o-OV1-nlYMu7sdVAL7SV2PrxVyjZFZZmB3O6LPGaYXlZWq",
    "naming_other_stuff_2": "",
    "naming_other_stuff_3": "",
    "naming_other_stuff_4": "",
    "naming_other_stuff_5": ""

  }
```

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok",
    "wording": "",    //材料说明
    "audit_id": 12345 //审核单id
  }
```

#### 3.2.3 小程序改名审核状态查询

#### 3.2.4 微信认证名称检测

#### 3.2.5 修改头像

#### 3.2.5 修改功能介绍

### 3.3 成员管理

#### 3.3.1 绑定微信用户为小程序体验者

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/bind/bindtester/三方应用appid/公众号appid`

__请求参数：__

```json
  {
    "wechatid":"testid"   //微信号
  }
```

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok",
    "userstr":"xxxxxxxxx" //人员对应的唯一字符串,解绑时有用
  }
```

#### 3.3.2 解除绑定小程序的体验者

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/bind/unbindtester/三方应用appid/公众号appid`

__请求参数：__

```json
  {
    "wechatid":"testid",   //微信号
    "userstr":"xxxxxx"    //人员对应的唯一字符串（可通过获取体验者api获取已绑定人员的字符串，userstr和wechatid填写其中一个即可）
  }
```

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok"
  }
```

#### 3.3.3 获取体验者列表

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/bind/memberauth/三方应用appid/公众号appid`

__请求参数：__

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok",
    "members":[{              //人员列表
      "userstr" : "xxxxxxxx"  //人员对应的唯一字符串
    },{
      "userstr" : "yyyyyyyy"
    }]
  }
```

### 3.4 代码管理

#### 3.4.1 为授权的小程序帐号上传小程序代码

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/commit/三方应用appid/公众号appid`

__请求参数：__


```json
  {
    "template_id":0,          //代码库中的代码模版ID
    "ext_json":"JSON_STRING", //*ext_json需为string类型，请参考下面的格式*
    "user_version":"V1.0",    //代码版本号，开发者可自定义
    "user_desc":"test"        //代码描述，开发者可自定义
  }
```
** ext_json需为string类型，格式示例如下 ：**
```json
{
    "extAppid":"",            //授权方Appid，可填入商户AppID，以区分不同商户
    "ext":{                   //自定义字段仅允许在这里定义，可在小程序中调用
        "attr1":"value1",
        "attr2":"value2"
    },
    "extPages":{              //页面配置
        "index":{
        },
        "search/index":{
        }
    },
    "pages":["index","search/index"],
    "window":{
    },
    "networkTimeout":{
    },
    "tabBar":{
    }
}
```

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok"
  }
```

#### 3.4.2 为授权的小程序帐号上传小程序代码并获取预览二维码地址

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/preview/三方应用appid/公众号appid?force=false&path=页面路径`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| force | bool | 否    | true是每次重新提交代码 |  
| path | string | 否    | 小程序内页面路径 | 

```json
  {
    "template_id":0,          //代码库中的代码模版ID
    "ext_json":"JSON_STRING", //*ext_json需为string类型，请参考下面的格式*
    "user_version":"V1.0",    //代码版本号，开发者可自定义
    "user_desc":"test"        //代码描述，开发者可自定义
  }
```
** ext_json需为string类型，格式示例如下 ：**
```json
{
    "extAppid":"",            //授权方Appid，可填入商户AppID，以区分不同商户
    "ext":{                   //自定义字段仅允许在这里定义，可在小程序中调用
        "attr1":"value1",
        "attr2":"value2"
    },
    "extPages":{              //页面配置
        "index":{
        },
        "search/index":{
        }
    },
    "pages":["index","search/index"],
    "window":{
    },
    "networkTimeout":{
    },
    "tabBar":{
    }
}
```

__响应内容：__

```json
  {
    "errcode": 0,
    "errmsg": "ok",
    "url": "http://www.ishanshan.com/wxproxy/mini/code/prevqrcode/wx7c44bb2354440737/wxfd1cd3fd79937051/20180829124436"
  }
```

#### 3.4.3 获取体验小程序的体验二维码

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getqrcode/三方应用appid/公众号appid`

__响应内容：__

成功时返回图片
失败时返回JSON格式错误码
```json
{
    "errcode":-1,
    "errmsg":"system error",
    "url": "http://127.0.0.1:8011/wxproxy/mini/code/prevqrcode/wx7c44bb2354440737/wxfd1cd3fd79937051/20180829112437"
}
```

### 3.4.4 获取体验小程序的体验二维码

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getqrcodeex/三方应用appid/公众号appid?path=`

__请求参数：__

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| path | string | 否    | 小程序内页面路径 | 

__响应内容：__

```json
{
    "errcode":0,
    "errmsg":"ok",
    "url": "http://127.0.0.1:8011/wxproxy/mini/code/prevqrcode/wx7c44bb2354440737/wxfd1cd3fd79937051/20180829112437"
}
```


### 3.4.5 获取小程序的二维码(个数限制)

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getwxacode/三方应用appid/公众号appid?path=&force=false`

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| path  | string | 否    | 小程序内页面路径 | 
| force | bool   | 否    | 强制重新获取 | 

__响应内容：__

```json
{
    "errcode":0,
    "errmsg":"ok",
    "url": "http://127.0.0.1:8011/wxproxy/mini/code/prevqrcode/wx7c44bb2354440737/wxfd1cd3fd79937051/20180829112437"
}
```

### 3.4.6 获取小程序的二维码(个数不限)

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getwxacodeunlimit/三方应用appid/公众号appid?page=&scene=&force=false`

| 参数名称       | 类型     | 是否必须 | 描述                |
| :---------- | :------ | :---- | :----------------- |
| path  | string | 否    | 小程序内页面路径 | 
| scene  | string | 否    | 页面参数 | 
| force | bool   | 否    | 强制重新获取 | 

__响应内容：__

```json
{
    "errcode":0,
    "errmsg":"ok",
    "url": "http://127.0.0.1:8011/wxproxy/mini/code/prevqrcode/wx7c44bb2354440737/wxfd1cd3fd79937051/20180829112437"
}
```

#### 3.4.5 获取授权小程序帐号的可选类目

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getcategroy/三方应用appid/公众号appid`

__响应内容：__

```json
{
    "errcode":0,
    "errmsg": "ok",
    "category_list" : [               //可填选的类目列表
        {
            "first_class":"工具",       //一级类目名称
            "second_class":"备忘录",    //二级类目名称
            "first_id":1,              //一级类目的ID编号
            "second_id":2              //二级类目的ID编号
        },{
            "first_class":"教育",
            "second_class":"学历教育",
            "third_class":"高等",       //三级类目名称
            "first_id":3,
            "second_id":4,
            "third_id":5                //三级类目的ID编号
        }
    ]
}
```

#### 3.4.6 获取小程序的第三方提交代码的页面配置

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/getpage/三方应用appid/公众号appid`

__响应内容：__

```json
{
    "errcode":0,
    "errmsg":"ok",
    "page_list":[         //page_list 页面配置列表
        "index",
        "page\/list",
        "page\/detail"
    ]
}
```

#### 3.4.7 将第三方提交的代码包提交审核

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/submitaudit/三方应用appid/公众号appid`

__请求参数：__

```json
{
    "item_list": [                //提交审核项的一个列表（至少填写1项，至多填写5项）
    {
        "address":"index",        //小程序的页面，可通过“获取小程序的第三方提交代码的页面配置”接口获得
        "tag":"学习 生活",         //小程序的标签，多个标签用空格分隔，标签不能多于10个，标签长度不超过20
        "first_class": "文娱",    //一级类目名称，可通过“获取授权小程序帐号的可选类目”接口获得
        "second_class": "资讯",   //二级类目(同上)
        "first_id":1,             //一级类目的ID，可通过“获取授权小程序帐号的可选类目”接口获得
        "second_id":2,            //二级类目的ID(同上)
        "title": "首页"           //小程序页面的标题,标题长度不超过32
    },{
        "address":"page/logs/logs",
        "tag":"学习 工作",
        "first_class": "教育",
        "second_class": "学历教育",
        "third_class": "高等",
        "first_id":3,
        "second_id":4,
        "third_id":5,
        "title": "日志"
    }
  ]
}
```

__响应内容：__

```json
  {
  "errcode":0,
  "errmsg":"ok",
  "auditid":1234567 //审核编号需保存
  }
```

#### 3.4.8 查询某个指定版本的审核状态

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/queryauditstatus/三方应用appid/公众号appid`

__请求参数：__

```json
{
  "auditid":1234567 //送审时返回的编号
}
```

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok",
    "status":1, //审核状态，其中0为审核成功，1为审核失败，2为审核中
    "reason":"帐号信息不合规范"//当status=1，审核被拒绝时，返回的拒绝原因
  }
```

#### 3.4.9 查询最新一次提交的审核状态

__请求方式：__ `GET`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/querylastauditstatus/三方应用appid/公众号appid`

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok",
    "status":1, //审核状态，其中0为审核成功，1为审核失败，2为审核中
    "reason":"帐号信息不合规范"//当status=1，审核被拒绝时，返回的拒绝原因
  }
```

#### 3.4.10 发布已通过审核的小程序

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/dorelease/三方应用appid/公众号appid`

__请求参数：__

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok"
  }
```

#### 3.4.11 修改小程序线上代码的可见状态

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/changevisitstatus/三方应用appid/公众号appid`

__请求参数：__

```json
{
  "action":"close" //设置可访问状态，发布后默认可访问，close为不可见，open为可见
}
```

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok"
  }
```

#### 3.4.12 小程序版本回退

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/revertcoderelease/三方应用appid/公众号appid`

__响应内容：__

```json
  {
    "errcode":0,
    "errmsg":"ok"
  }
```

#### 3.4.13 查询当前设置的最低基础库版本及各版本用户占比

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/queryweappsupportversion/三方应用appid/公众号appid`

__请求参数：__

__响应内容：__

```json
  {
      "errcode": 0,
      "errmsg": "ok",
      "now_version": "1.0.0",       //当前版本
      "uv_info": {                  //受影响用户占比，item参数里面为一系列数组，每个数组带有参数percentage（受影响比例）以及version（版本号）
          "items": [{
                  "percentage": 0,
                  "version": "1.0.0"
              },{
                  "percentage": 0,
                  "version": "1.0.1"
              },{
                  "percentage": 0,
                  "version": "1.1.0"
              },{
              }
          ]
      } 
  }
```

#### 3.4.14 设置最低基础库版本

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/code/setminweappsupportversion/三方应用appid/公众号appid`

__请求参数：__

```json
{
  "version":"1.0.0"
}
```

__响应内容：__

```json
{
"errcode" : 0,
"errmsg" : "ok"
}
```

### 3.5 小程序代码模版库管理

#### 3.5.1 获取草稿箱内的所有临时代码草稿

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/codetplmgr/gettemplatedraftlist/三方应用appid`

__请求参数：__

__响应内容：__

```json
{
  "errcode": 0,
  "errmsg": "ok",
  "draft_template_list":[{ 
    "create_time": 1488965944,  //开发者上传草稿时间
    "user_version": "VVV",      //模版版本号，开发者自定义字段
    "user_desc": "AAS",         //模版描述 开发者自定义字段
    "draft_id":0                //草稿id
  },{ 
    "create_time": 1504790906,
    "user_version": "11",
    "user_desc": "111111",
    "draft_id": 4 
   }]
}
```

#### 3.5.2 获取代码模版库中的所有小程序代码模版

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/codetplmgr/gettemplatelist/三方应用appid`

__请求参数：__

__响应内容：__

```json
{
  "errcode": 0,
  "errmsg": "ok",
  "template_list":[{ 
    "create_time": 1488965944,  //被添加为模版的时间
    "user_version": "VVV",      //模版版本号，开发者自定义字段
    "user_desc": "AAS",         //模版描述 开发者自定义字段
    "template_id":0             //模版id
  },{ 
    "create_time": 1504790906,
    "user_version": "11",
    "user_desc": "111111",
    "template_id": 4 
   }]
}
```

#### 3.5.3 将草稿箱的草稿选为小程序代码模版

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/codetplmgr/addtotemplate/三方应用appid`

__请求参数：__

```json
{
  "draft_id":0 //草稿ID，本字段可通过“ 获取草稿箱内的所有临时代码草稿 ”接口获得
}
```

__响应内容：__

```json
{
"errcode" : 0,
"errmsg" : "ok"
}
```

#### 3.5.4 删除指定小程序代码模版

__请求方式：__ `POST`

__请 求 头：__ `Content-Type: application/json;charset=utf-8`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/codetplmgr/deletetemplate/三方应用appid`

__请求参数：__

```json
{
  "template_id":0 //要删除的模版ID
}
```

__响应内容：__

```json
{
"errcode" : 0,
"errmsg" : "ok"
}
```

### 3.6 微信登录

__请求方式：__ `GET`

__访问地址：__ `http://www.ishanshan.com/wxproxy/mini/user/login/三方应用appid/公众号appid?js_code=JSCODE`

__请求参数：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| js_code | String | 登录时获取的 code |  


__响应内容：__

```json
{
  "openid":"OPENID",          //用户唯一标识的openid
  "session_key":"SESSIONKEY"  //会话密钥
}
```

### 3.7 小程序模板设置

### 3.8 微信开放平台帐号管理

### 3.9 基础信息设置

### 3.10 小程序插件管理权限集