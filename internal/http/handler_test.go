package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sportNews/internal/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newMockTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	App{}.router(e)

	return e
}

func mockRequest(e http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	return w
}

func decodeHttpResp(w *httptest.ResponseRecorder) (response.HttpResp, error) {
	var r response.HttpResp
	err := json.Unmarshal(w.Body.Bytes(), &r)

	return r, err
}

func TestPageNotFound(t *testing.T) {
	t.Helper()

	e := newMockTestRouter()

	w := mockRequest(e, http.MethodGet, "/notExistPath")

	resp, err := decodeHttpResp(w)

	require.NoError(t, err)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
	require.Nil(t, resp.Data)
}

func TestMethodNotAllowed(t *testing.T) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	e := gin.New()
	App{}.router(e)

	w := mockRequest(e, http.MethodPost, "/healthz")

	resp, err := decodeHttpResp(w)

	require.NoError(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
	require.Equal(t, response.NO_ALLOW_METHDO, resp.Code)
	require.Equal(t, http.StatusText(http.StatusMethodNotAllowed), resp.Msg)
	require.Nil(t, resp.Data)
}
