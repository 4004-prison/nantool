package errorman

import "fmt"

// Error error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"msg,omitempty"`
}

// New return error with code and error message
func New(code int64, params ...interface{}) Error {
	if len(params) > 0 {
		return Error{
			Code:    code,
			Message: fmt.Sprintf(Translation(code), params...),
		}
	}
	return Error{
		Code: code,
	}
}

// Error return error string
func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return Translation(e.Code)
}

// NewErrorWithString return error with code and string message
func NewErrorWithString(code int64, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
