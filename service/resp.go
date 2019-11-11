package service

type Resp interface {
	IsSuccess() bool
	GetErrcode() int
	GetErrmsg() string
}

type CommonResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg,omitempty"`
}

var (
	SUCCESS_RESP = NewCommonResp(0, "ok")
)

func (c CommonResp) IsSuccess() bool {
	return c.Errcode == 0
}
func (c CommonResp) GetErrcode() int {
	return c.Errcode
}
func (c CommonResp) GetErrmsg() string {
	return c.Errmsg
}

func NewCommonResp(code int, msg string) *CommonResp {
	return &CommonResp{
		Errcode: code,
		Errmsg:  msg,
	}
}

func NewServerErrorResp(err error) *CommonResp {
	return NewCommonResp(1000, err.Error())
}
