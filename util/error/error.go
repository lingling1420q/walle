package error

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UnknowError = EvaError{
		AppId:   "walle",
		Code:    1001,
		Message: "未知错误",
		Status:  codes.Unknown,
	}
	PanicError = EvaError{
		AppId:   "walle",
		Code:    1002,
		Message: "异常错误",
		Status:  codes.Unknown,
	}
	HttpError = EvaError{
		AppId:   "walle",
		Code:    2001,
		Message: "http错误",
		Status:  codes.Internal,
	}
	TypeError = EvaError{
		AppId:   "walle",
		Code:    3001,
		Message: "类型转换错误",
		Status:  codes.Internal,
	}
	ContextDieError = EvaError{
		AppId:   "walle",
		Code:    4001,
		Message: "context已到期",
		Status:  codes.DeadlineExceeded,
	}
	GRpcError = EvaError{
		AppId:   "walle",
		Code:    5001,
		Message: "内部错误",
		Status:  codes.Internal,
	}
	NotFoundError = EvaError{
		AppId:   "walle",
		Code:    404,
		Message: "未找到",
		Status:  codes.NotFound,
	}
)

type EvaError struct {
	AppId   string     `json:"appId"`   //错误发生的服务
	Code    int32      `json:"code"`    //错误码 业务的code
	Message string     `json:"message"` //错误消息
	Detail  string     `json:"detail"`  //更详细的错误消息 不对外展示的
	Status  codes.Code `json:"status"`  //grpc的错误码
}

func (m EvaError) SetDetail(detail string) EvaError {
	m.Detail = detail
	return m
}
func (m EvaError) SetCode(code int32) EvaError {
	m.Code = code
	return m
}
func (m EvaError) SetMessage(msg string) EvaError {
	m.Message = msg
	return m
}
func (e EvaError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(message, detail string, code int32, status codes.Code) error {
	return &EvaError{
		AppId:   "walle",
		Code:    code,
		Detail:  detail,
		Message: message,
		Status:  status,
	}
}

func Parse(err string) *EvaError {
	e := new(EvaError)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		a := UnknowError.SetDetail(err)
		return &a
	}
	return e
}

func FromError(err error) *EvaError {
	if verr, ok := err.(*EvaError); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

func EncodeStatus(e *EvaError) *status.Status {
	status := status.New(e.Status, e.Error())

	return status
}

func DecodeStatus(e error) *EvaError {

	status, ok := status.FromError(e)

	if !ok {
		return Parse(e.Error())
	} else {
		return Parse(status.Message())
	}
}
