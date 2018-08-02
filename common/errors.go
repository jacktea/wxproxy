package common

import "errors"

var (
	NO_AUTH_ACCESS_TOKEN 	= errors.New("没找到授权TOKEN")
	NO_DATA					= errors.New("没有数据")
	INVALID_URL				= errors.New("无效的链接")
)
