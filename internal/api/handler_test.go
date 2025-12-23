package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sportNews/internal/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newMockTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	e := gin.New()
	App{}.router(e)

	return e
}

func mockRequest(t *testing.T, e http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	return w
}

func TestPageNotFound(t *testing.T) {
	e := newMockTestRouter(t)

	w := mockRequest(t, e, http.MethodGet, "/notExistPath")

	var resp response.HttpResp
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
	require.Nil(t, resp.Data)
}

func TestMethodNotAllowed(t *testing.T) {
	e := newMockTestRouter(t)

	w := mockRequest(t, e, http.MethodPost, "/healthz")

	var resp response.HttpResp
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
	require.Equal(t, response.NO_ALLOW_METHDO, resp.Code)
	require.Equal(t, http.StatusText(http.StatusMethodNotAllowed), resp.Msg)
	require.Nil(t, resp.Data)
}

func TestHealthzRouter(t *testing.T) {
	e := newMockTestRouter(t)

	w := mockRequest(t, e, http.MethodGet, "/healthz")

	var resp struct {
		Status string `json:"status"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, struct {
		Status string `json:"status"`
	}{Status: "ok"}, resp)
}
