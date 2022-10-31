package rest

import "net/http"

var (
	NotFoundErr       = NewError("not found", "", http.StatusNotFound)
	InternalServerErr = NewError("internal server error", "", http.StatusInternalServerError)
	BadReqErr         = NewError("bad request", "", http.StatusBadRequest)
)

type Error struct {
	Message      string
	DebugMessage string
	httpCode     int
}

func (e *Error) SetMsg(msg string) *Error {
	return NewError(msg, e.DebugMessage, e.httpCode)
}

func (e *Error) SetDebugMsg(msg string) *Error {
	return NewError(e.Message, msg, e.httpCode)
}

func (e *Error) Error() string {
	return e.DebugMessage
}

func NewError(message, debugMessage string, code int) *Error {
	return &Error{
		Message:      message,
		DebugMessage: debugMessage,
		httpCode:     code,
	}
}
