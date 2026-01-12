package service

import (
	"os"
	"sportNews/internal/assets"
	"sportNews/pkg/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.InitLogger(false)
	code := m.Run()
	log.CloseLogger()
	os.Exit(code)
}

func setupAssets(t *testing.T) {
	assets.Setup("https://cdn.test")
	t.Cleanup(func() { assets.Setup("") })
}
