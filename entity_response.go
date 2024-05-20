package common

import "encoding/json"

type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
type ResponseGeneric[T any] struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
}

func (res *Response) AppendMsg(msg string) *Response {

	if res.Msg != "" {
		msg = res.Msg + "," + msg
	}
	return &Response{
		Status: res.Status,
		Msg:    msg,
		Data:   res.Data,
	}
}
func (res *Response) WithData(data interface{}) *Response {
	return &Response{
		Status: res.Status,
		Msg:    res.Msg,
		Data:   data,
	}
}
func (res *Response) ToBytes() []byte {
	raw, _ := json.Marshal(res)
	return raw
}

func response(code int, msg string) *Response {
	return &Response{
		Status: code,
		Msg:    msg,
		Data:   nil,
	}
}

var (
	OK                 = response(0, "服务调用成功")
	BaseErrorNo        = OK.Status
	ErrNotFound        = response(BaseErrorNo+1, "服务调用成功,但没有找到相应数据")
	ErrDeleteError     = response(BaseErrorNo+2, "无法删除:")
	ErrParam           = response(BaseErrorNo+3, "参数有误")
	ErrSignParam       = response(BaseErrorNo+4, "签名参数有误")
	ErrReqBodyRead     = response(BaseErrorNo+5, "读取body有误")
	ErrReqBodyParse    = response(BaseErrorNo+6, "请求参数反序列化错误")
	ErrPubSubPublish   = response(BaseErrorNo+7, "发布消息错误")
	ErrServiceInvokeDB = response(BaseErrorNo+8, "调用数据库异常")
	ErrAddEdgeToGraph  = response(BaseErrorNo+9, "添加边到图异常")
	ErrService         = response(BaseErrorNo+10, "服务异常")
	ErrListUnMashal    = response(BaseErrorNo+11, "列表反序列化错误")
	ErrModelParse      = response(BaseErrorNo+12, "对象反序列化错误")
	ErrExists          = response(BaseErrorNo+13, "已存在")
	ErrAuthz           = response(BaseErrorNo+14, "没有权限")
)
