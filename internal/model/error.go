package model

type Error struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

func (e *Error) Send() *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Err:     e.Err,
	}
}
