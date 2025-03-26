package main

import (
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sportNews/api"
	"sportNews/conf"
	"sportNews/conf/database"
	"sportNews/internal/assets"
	"sportNews/internal/log"
	"syscall"
	"time"
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
	log.Infof("setting conf", "conf:%v", config)
	log.Info("app info", zap.String("Project", config.App.Name))

	// setting database
	db, err := database.New(config.DB)
	if err != nil {
		log.Error("init database", zap.Error(err))
	}

	assets.Setup(config.Assets)

	serv := api.New(ctx, &config.App, db)
	serv.Run()

	shutdown(config, db)
}

func shutdown(config *conf.Conf, db *xorm.EngineGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	log.Debugf("service shutdown", "get a signal %s. %s Server is shutdowning ...", s.String(), config.App.Name)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.Close()

	log.Infof("service shutdown", "%s Server is exit", config.App.Name)
}
