package response

import (
	"net/http"
	"strconv"
)

// пакет для работы с ответом на запросы
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: strconv.Itoa(http.StatusForbidden),
		Error:  msg,
	}
}
