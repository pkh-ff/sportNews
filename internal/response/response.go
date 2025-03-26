package response

import "net/http"

func Success(data interface{}) *HttpResp {
	return &HttpResp{
		Code:       0,
		StatusCode: http.StatusOK,
		Data:       data,
		Msg:        "success",
	}
}

// NotFound Not found page error response
func NotFound() *HttpResp {
	return newHttpException(http.StatusNotFound, NOT_FOUND, http.StatusText(http.StatusNotFound))
}

// ParameterError 參數錯誤
func ParameterError() *HttpResp {
	return newHttpException(http.StatusBadRequest, PARAMETER_ERROR, "Missing required parameter error or parameter setting error")
}

// InternalError Service Internal Error response
func InternalError() *HttpResp {
	return newHttpException(http.StatusInternalServerError, INTERNA_ERROR, http.StatusText(http.StatusInternalServerError))
}

// unknownError Unknown Error response
func unknownError(message string) *HttpResp {
	return newHttpException(http.StatusInternalServerError, UNKNOWN_ERROR, message)
}

func InsertFail() *HttpResp {
	return newHttpException(http.StatusInternalServerError, INSERT_DATA_ERROR, "Data creation failed")
}

func UpdateFail() *HttpResp {
	return newHttpException(http.StatusInternalServerError, UPDATE_DATA_ERROR, "Data update failed")
}

func DeleteFail() *HttpResp {
	return newHttpException(http.StatusInternalServerError, DELETE_DATA_ERROR, "Data delete failed")
}

// noAllowMethod
// Not Allow Method
func noAllowMethod() *HttpResp {
	return newHttpException(http.StatusMethodNotAllowed, NO_ALLOW_METHDO, http.StatusText(http.StatusMethodNotAllowed))
}

func newHttpException(statusCode int, errorCode int, msg string) *HttpResp {
	return &HttpResp{
		Code:       errorCode,
		StatusCode: statusCode,
		Msg:        msg,
	}
}
