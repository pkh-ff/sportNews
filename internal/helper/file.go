package helper

import (
	"fmt"
	"os"
)

// WriteToFile
// 將文字內容寫入檔案中
func WriteToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 將內容寫入檔案
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	return nil
}
