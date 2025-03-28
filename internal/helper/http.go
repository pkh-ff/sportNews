package helper

import (
	"net/http"
	"net/url"
)

// BuildHttpRequest
// 建立HttpRequest實例
func BuildHttpRequest(baseURL, method string, queryParams, hdrs map[string]string) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// 設定query string參數
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// 建立請求
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// 設定標頭
	for k, v := range hdrs {
		req.Header.Set(k, v)
	}

	return req, nil
}
