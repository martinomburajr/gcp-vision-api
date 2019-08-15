package router

import (
	"log"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

type AppError struct {
	Error   error
	Message string
	Code    int
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *AppError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func AppErrorf(code int, msg string, err error) *AppError {
	return &AppError{
		Error:   err,
		Message: msg,
		Code:    code,
	}
}