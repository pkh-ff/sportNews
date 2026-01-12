package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	url := server.URL + "/api/test"

	resp, err := SendHTTPRequest(url, http.MethodGet, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendHTTPRequestWithStatusNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	resp, err := SendHTTPRequest(server.URL, http.MethodGet, nil, nil)

	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestSendHTTPRequestWithNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	resp, err := SendHTTPRequest(server.URL, http.MethodGet, nil, nil)

	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestSendHTTPRequest_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := SendHTTPRequest(server.URL, http.MethodGet, nil, nil)
	require.Error(t, err)
}
