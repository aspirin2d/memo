package memo

import (
	"github.com/gin-gonic/gin"
)

var OKMessage = gin.H{"ok": true}

type WrapError struct {
	ErrorMessage
	code  int
	error error
}

type ErrorMessage struct {
	Message string `json:"msg"`
}

func (w WrapError) Error() string {
	return w.Message
}

func (w WrapError) Code() int {
	return w.code
}

func NewWrapError(code int, err error, msg string) WrapError {
	if msg == "" {
		msg = err.Error()
	}

	return WrapError{
		code:         code,
		error:        err,
		ErrorMessage: ErrorMessage{msg},
	}
}

func (m *Memo) AbortWithError(c *gin.Context, err error) {
	switch e := err.(type) {
	case WrapError:
		c.AbortWithStatusJSON(e.code, ErrorMessage{e.Message})
	default:
		c.AbortWithStatusJSON(500, ErrorMessage{"oops, an unknown error occurred, please try later"})
	}
	if m.Logger != nil {
		m.Logger.Error(err)
	}
}
