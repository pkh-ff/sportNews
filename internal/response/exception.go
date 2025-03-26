package response

import (
	"github.com/gin-gonic/gin"
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
	SUCCESS         = 0
	PARAMETER_ERROR = 4001

	NOT_FOUND = 4041

	NO_ALLOW_METHDO = 4051

	INTERNA_ERROR     = 5000
	UNKNOWN_ERROR     = 5001
	INSERT_DATA_ERROR = 5002
	UPDATE_DATA_ERROR = 5003
	DELETE_DATA_ERROR = 5004
)

type HttpResp struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func (e *HttpResp) Error() string {
	return e.Msg
}

type HandlerFunc func(c *gin.Context) error

func ErrHandler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		err = h(c)
		if err != nil {
			var apiException *HttpResp
			if h, ok := err.(*HttpResp); ok {
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

var (
	ErrInternal   = InternalError()
	ErrNoRows     = NotFound()
	ErrParameter  = ParameterError()
	ErrInsertFail = InsertFail()
	ErrUpdateFail = UpdateFail()
	ErrDeleteFail = DeleteFail()
)
