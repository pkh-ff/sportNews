package helper

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"sportNews/pkg/log"
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
