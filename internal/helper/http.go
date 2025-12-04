package helper

import (
	"net/http"
	"net/url"
	"sportNews/pkg/log"
	"time"

	"go.uber.org/zap"
)

// BuildHttpRequest
// 建立HttpRequest實例
func BuildHttpRequest(baseURL, method string, queryParams, hdrs map[string]string) (*http.Request, error) {
	log.Info("BuildHttpRequest: Constructing HTTP request",
		zap.String("baseURL", baseURL),
		zap.String("method", method),
		zap.Any("queryParams", queryParams),
		zap.Any("headers", hdrs),
	)
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Error("BuildHttpRequest: Failed to create HTTP request", zap.String("method", method), zap.String("url", u.String()), zap.Error(err))
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

	log.Info("BuildHttpRequest: HTTP request constructed successfully", zap.String("url", req.URL.String()), zap.String("method", method))

	return req, nil
}

// SendHTTPRequest
// 發送HTTP Request
func SendHTTPRequest(url, method string, headers, queryParams map[string]string) (*http.Response, error) {
	log.Info("SendHTTPRequest: Sending HTTP request",
		zap.String("url", url),
		zap.String("method", method),
		zap.Any("headers", headers),
		zap.Any("queryParams", queryParams),
	)

	req, err := BuildHttpRequest(url, method, queryParams, headers)
	if err != nil {
		log.Error("SendHTTPRequest: Failed to build HTTP request", zap.Error(err))
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // 設定超時時間
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("SendHTTPRequest: HTTP request failed", zap.Error(err))
		return nil, err
	}

	log.Info("SendHTTPRequest: Received HTTP response", zap.String("url", url), zap.Int("statusCode", resp.StatusCode))
	// 檢查 HTTP 回應碼
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil
	}

	return resp, nil
}
