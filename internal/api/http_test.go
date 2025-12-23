package api

import (
	"os"
	"sportNews/pkg/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.InitLogger(false)
	code := m.Run()
	log.CloseLogger()
	os.Exit(code)
}
