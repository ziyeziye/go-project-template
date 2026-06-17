package errno

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 定义一些常用的错误 可接入多语言支持
var (
	ServerError  = NewError(10500, "Server internal error")
	Unauthorized = NewError(100402, "Unauthorized")
)

func OtherError(message string) Error {
	return NewError(10403, message)
}

func ErrOf(err error) Error {
	var newErr Error
	if e, ok := err.(Error); ok {
		newErr = e
	} else {
		newErr = OtherError(err.Error())
	}
	return newErr
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(Code int, msg string) Error {
	return Error{
		Code: Code,
		Msg:  msg,
	}
}
