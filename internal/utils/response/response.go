package response

import (
	"net/http"

	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int    `json:"-"`
	Status     bool   `json:"status"`
	Message    string `json:"message"`
	Error      any    `json:"error,omitempty"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(message string, data any, meta ...pagination.Meta) Response {
	res := Response{
		StatusCode: http.StatusOK,
		Status:     true,
		Message:    message,
		Data:       data,
	}

	if len(meta) > 0 {
		res.Meta = meta
	}

	return res
}

func BuildResponseFailed(message string, err error, data ...any) Response {
	res := Response{
		StatusCode: http.StatusInternalServerError,
		Status:     false,
		Message:    message,
		Error:      err,
	}

	if myErr, ok := err.(MyError.Error); ok {
		res.StatusCode = myErr.StatusCode
		res.Error = myErr
	}

	if len(data) > 0 {
		res.Data = data
	}

	return res
}

// Menangani jika status code belum di set di error contohnya 201 ketika create something
func (r Response) SetStatus(statusCode int) Response {
	res := r
	res.StatusCode = statusCode
	return res
}

// sama seperti ctx.Json tetapi ini untuk response
// yang sudah di set status code nya
func (r Response) Send(ctx *gin.Context) {
	sendStatus := r.StatusCode
	ctx.JSON(sendStatus, r)
}

// sama seperti ctx.Json tetapi ini untuk response
// yang sudah di set status code nya
// dan langsung di abort
func (r Response) SendWithAbort(ctx *gin.Context) {
	sendStatus := r.StatusCode
	ctx.AbortWithStatusJSON(sendStatus, r)
}
