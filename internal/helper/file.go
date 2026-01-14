package helper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sportNews/pkg/log"

	"go.uber.org/zap"
)

// WriteToFile
// 將文字內容寫入檔案中
func WriteToFile(filename, content string) error {
	log.Info("WriteToFile: Opening file for writing", zap.String("filename", filename))
	file, err := os.Create(filename)
	if err != nil {
		log.Error("WriteToFile: Failed to open file", zap.String("filename", filename), zap.Error(err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Error("WriteToFile: Failed to close file", zap.String("filename", filename), zap.Error(cerr))
			err = fmt.Errorf("WriteToFile: failed to close file %s: %w", filename, cerr)
		}
	}()

	// 將內容寫入檔案
	_, err = file.WriteString(content)
	if err != nil {
		log.Error("WriteToFile: Failed to write content", zap.String("filename", filename), zap.Error(err))
		return fmt.Errorf("WriteToFile: failed to write content to file %s: %w", filename, err)
	}

	// 確保內容真正寫入磁碟
	if err := file.Sync(); err != nil {
		log.Error("WriteToFile: Failed to sync file", zap.String("filename", filename), zap.Error(err))
		return fmt.Errorf("WriteToFile: failed to sync file %s: %w", filename, err)
	}

	log.Info("WriteToFile: Successfully wrote content to file", zap.String("filename", filename))
	return nil
}

// GetFileNameFromURL
// 取得檔名(含副檔名)
func GetFileNameFromURL(urlStr string) (string, error) {
	// 解析 URL，去除查詢字符串
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Error("GetFileNameFromURL, Unable to parse file", zap.Error(err))
		return "", err
	}

	// 去除查詢字符串
	fileName := path.Base(parsedURL.Path)

	return fileName, nil
}

// DownloadFileFromUrl
// 從網路連結下載檔案
func DownloadFileFromUrl(url string) ([]byte, error) {
	// 發送 HTTP GET 請求以獲取圖片
	resp, err := http.Get(url)
	if err != nil {
		log.Error("DownloadFileFromUrl, failed to download image", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 讀取檔案數據
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("DownloadFileFromUrl, failed to read image data", zap.Error(err))
		return nil, err
	}

	return imgData, nil
}

// FileURLExists
// 透過http status來判斷網路上檔案是否存在
func FileURLExists(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		log.Error("FileURLExists: failed to perform HEAD request", zap.String("url", url), zap.Error(err))
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetFileExtensionFromURL
// 從URL路徑中取得副檔名
func GetFileExtensionFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Error("GetFileExtensionFromURL, Unable to parse file", zap.Error(err))
		return "", err
	}

	// 取副檔名
	ext := path.Ext(u.Path) // 只看 Path 部分

	return ext, nil
}
