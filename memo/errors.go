package memo

type ErrorMessage struct {
	Message string `json:"msg"`
}

// NewError returns a new ErrorMessage with the given message
// and logs the error if log is true
func (m *Memo) NewError(err error, msg string, log bool) ErrorMessage {
	if log {
		m.logger.Error(err)
	}
	return ErrorMessage{Message: msg}
}
