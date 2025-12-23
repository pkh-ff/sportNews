package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sportNews/conf"
	"sportNews/conf/database"
	"sportNews/internal/assets"
	"sportNews/internal/http"
	"sportNews/pkg/log"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	confPath string
)

func main() {
	flag.StringVar(&confPath, "c", "app.conf.yaml", "default conf path")
	flag.Parse()

	ctx := context.Background()

	// loading conf
	config, err := conf.New(confPath)
	if err != nil {
		log.Error("setting conf", zap.Error(err))
	}

	// logger設定
	log.InitLogger(config.App.Debug)

	log.Infof("Config loaded: %+v", config)
	log.Info("App Info", zap.String("Project", config.App.Name))

	// 設定資料庫
	db, err := database.New(config.App.Debug, config.DB)
	if err != nil {
		log.Error("init database", zap.Error(err))
	}

	assets.Setup(config.Assets)

	serv := http.New(ctx, &config.App, db)
	serv.Run()

	shutdown(config, db)
}

func shutdown(config *conf.Conf, db *xorm.EngineGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	log.Debugf("get a signal %s. %s Server is shutting down ...", s.String(), config.App.Name)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Close(); err != nil {
		log.Warnf("Database close error: %v", err)
	}

	log.Infof("%s Server is exiting...", config.App.Name)
	log.CloseLogger()
}
