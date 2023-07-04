package memo

type ModelError struct {
	Code    int
	Message string
	Log     string
}

func NewModelError(code int, message string, log string) *ModelError {
	return &ModelError{
		Code:    code,
		Message: message,
		Log:     log,
	}
}

type ErrMsg struct {
	Message string `json:"msg"`
}
