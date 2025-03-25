package httpError

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleNotFound(c *gin.Context) {
	handleErr := NotFound()
	c.JSON(handleErr.StatusCode, handleErr)
	return
}

func HandleNoAllowMethod(c *gin.Context) {
	handleErr := noAllowMethod()
	c.JSON(handleErr.StatusCode, handleErr)
	return
}

const (
	PARAMETER_ERROR = 4001

	NOT_FOUND = 4041

	NO_ALLOW_METHDO = 4051

	INTERNA_ERROR     = 5000
	UNKNOWN_ERROR     = 5001
	INSERT_DATA_ERROR = 5002
	UPDATE_DATA_ERROR = 5003
	DELETE_DATA_ERROR = 5004
)

type HttpException struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

func (e *HttpException) Error() string {
	return e.Msg
}

type HandlerFunc func(c *gin.Context) error

func ErrHandler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		err = h(c)
		if err != nil {
			var apiException *HttpException
			if h, ok := err.(*HttpException); ok {
				apiException = h
			} else if e, ok := err.(error); ok {
				apiException = unknownError(e.Error())
			} else {
				apiException = InternalError()
			}
			c.JSON(apiException.StatusCode, apiException)
			return
		}
	}
}

func newHttpException(statusCode int, errorCode int, msg string) *HttpException {
	return &HttpException{
		Code:       errorCode,
		StatusCode: statusCode,
		Msg:        msg,
	}
}

var (
	ErrInternal   = InternalError()
	ErrNoRows     = NotFound()
	ErrParameter  = parameterError()
	ErrInsertFail = insertFail()
	ErrUpdateFail = updateFail()
	ErrDeleteFail = deleteFail()
)

// NotFound Not found page error response
func NotFound() *HttpException {
	return newHttpException(http.StatusNotFound, NOT_FOUND, http.StatusText(http.StatusNotFound))
}

// parameterError 參數錯誤
func parameterError() *HttpException {
	return newHttpException(http.StatusBadRequest, PARAMETER_ERROR, "Missing required parameter error or parameter setting error")
}

// InternalError Service Internal Error response
func InternalError() *HttpException {
	return newHttpException(http.StatusInternalServerError, INTERNA_ERROR, http.StatusText(http.StatusInternalServerError))
}

// unknownError Unknown Error response
func unknownError(message string) *HttpException {
	return newHttpException(http.StatusInternalServerError, UNKNOWN_ERROR, message)
}

func insertFail() *HttpException {
	return newHttpException(http.StatusInternalServerError, INSERT_DATA_ERROR, "Data creation failed")
}

func updateFail() *HttpException {
	return newHttpException(http.StatusInternalServerError, UPDATE_DATA_ERROR, "Data update failed")
}

func deleteFail() *HttpException {
	return newHttpException(http.StatusInternalServerError, DELETE_DATA_ERROR, "Data delete failed")
}

// noAllowMethod
// Not Allow Method
func noAllowMethod() *HttpException {
	return newHttpException(http.StatusMethodNotAllowed, NO_ALLOW_METHDO, http.StatusText(http.StatusMethodNotAllowed))
}
