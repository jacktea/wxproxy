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

## 2. 功能

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

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upcmptoken/三方应用appid`

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.3 刷新三方应预授权码

__请求方式：__ `POST`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upcmpcode/三方应用appid`

__响应内容：__

| 参数名称       | 类型     | 描述             |
| :---------- | :------ | :----------------- |
| errcode | int | 响应码 0 成功 非0 失败 |  
| errmsg | string | 错误信息描述 |  

#### 2.2.4 刷新公众号/小程序访问toekn

__请求方式：__ `POST`

__访问地址：__ `http://www.ishanshan.com/wxproxy/api/token/upauthtoken/三方应用appid/公众号appid`

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

__访问地址：__ `http://www.ishanshan.com/wxproxy/api//account/getauthappinfo/三方应用appid/公众号appid`

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

