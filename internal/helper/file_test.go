package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileNameFromURL(t *testing.T) {
	url := "http://example.com/test.png"
	filename, err := GetFileNameFromURL(url)

	require.NoError(t, err)
	require.Equal(t, "test.png", filename)
}

func TestGetFileNameFromURLWithEmptyUrl(t *testing.T) {
	url := ""
	filename, err := GetFileNameFromURL(url)

	require.Error(t, err)
	require.Equal(t, "", filename)
}

func TestGetFileNameFromURWithNoxExtensionName(t *testing.T) {
	url := "http://example.com/test"
	filename, err := GetFileNameFromURL(url)

	require.NoError(t, err)
	require.Equal(t, "test", filename)
}

func TestGetFileNameFromURLWithQueryStringURL(t *testing.T) {
	url := "http://example.com/test.png?a=1"
	filename, err := GetFileNameFromURL(url)

	require.NoError(t, err)
	require.Equal(t, "test.png", filename)
}

func TestGetFileNameFromURLWithInvalidURL(t *testing.T) {
	url := "http://example.com\test.png"
	filename, err := GetFileNameFromURL(url)

	require.Error(t, err)
	require.Equal(t, "", filename)
}

func TestDownloadFileFromUrlWithSuccess(t *testing.T) {
	expectedData := []byte("fake image content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expectedData)
	}))
	defer server.Close()

	data, err := DownloadFileFromUrl(server.URL)

	require.NoError(t, err)
	assert.Equal(t, expectedData, data)
}

func TestDownloadFileFromUrlWithNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	data, err := DownloadFileFromUrl(server.URL)

	require.Error(t, err)
	assert.Nil(t, data)
}

func TestDownloadFileFromUrlWithNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	data, err := DownloadFileFromUrl(server.URL)

	require.Error(t, err)
	assert.Nil(t, data)
}

func TestFileURLExistsWithExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	exists, err := FileURLExists(server.URL + "/file.jpg")

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestFileURLExistsWithNotExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	exists, err := FileURLExists(server.URL + "/file.jpg")

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestFileURLExistsNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close() // 立即關閉

	exists, err := FileURLExists(server.URL + "/file.jpg")

	require.Error(t, err)
	assert.False(t, exists)
}

func TestGetFileExtensionFromURL(t *testing.T) {
	url := "http://test.com/test.png"
	ext, err := GetFileExtensionFromURL(url)

	require.NoError(t, err)
	require.Equal(t, ".png", ext)
}

func TestGetFileExtensionFromURLWithPathNoExt(t *testing.T) {
	url := "http://test.com/test"
	ext, err := GetFileExtensionFromURL(url)

	require.NoError(t, err)
	require.Equal(t, "", ext)
}

func TestGetFileExtensionFromURLWithPathQueryString(t *testing.T) {
	url := "http://test.com/test.png?a=1"
	ext, err := GetFileExtensionFromURL(url)

	require.NoError(t, err)
	require.Equal(t, ".png", ext)
}

func TestGetFileExtensionFromURLWithEmptyPath(t *testing.T) {
	url := ""
	ext, err := GetFileExtensionFromURL(url)

	require.NoError(t, err)
	require.Equal(t, "", ext)
}

func TestGetFileExtensionFromURLWithInvalidURL(t *testing.T) {
	url := "http://test.com\test.png?a=1"
	ext, err := GetFileExtensionFromURL(url)

	require.Error(t, err)
	require.Equal(t, "", ext)
}
