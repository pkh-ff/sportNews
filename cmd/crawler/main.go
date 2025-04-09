package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sportNews/conf"
	"sportNews/conf/aws"
	"sportNews/conf/database"
	"sportNews/internal/process"
	"sportNews/pkg/log"
	"sync"
	"syscall"
	"time"
	"xorm.io/xorm"
)

var (
	confPath string // config path
	wg       sync.WaitGroup
	wgCount  int
)

func main() {
	flag.StringVar(&confPath, "c", "process.conf.yaml", "default config path")
	flag.Parse()

	//ctx := context.Background()
	config, err := conf.New(confPath)
	if err != nil {
		log.Error("setting conf", zap.Error(err))
	}

	log.Infof("setting conf", "conf:%v", config)
	log.Info("app info", zap.String("Project", config.App.Name))

	// setting database
	db, err := database.New(config.App.Debug, config.DB)
	if err != nil {
		log.Error("init database", zap.Error(err))
	}

	s3Client, err := aws.New(config.Aws)

	// set log config
	log.InitLogger(config.App.Debug)

	crawlerProcess(*config, db, s3Client)

	shutdown(config, db)
	wg.Wait()
}

func crawlerProcess(conf conf.Conf, db *xorm.EngineGroup, s3Client *s3.Client) {
	wg.Add(3)
	wgCount = 3

	go func() {
		defer wg.Done()
		process.NDTVProcess(conf.App, db)
	}()

	go func() {
		defer wg.Done()
		process.BCCIProcess(conf.App, db)
	}()

	go func() {
		defer wg.Done()
		process.PictureSyncProcess(conf, s3Client, db)
	}()
}

func shutdown(config *conf.Conf, db *xorm.EngineGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	log.Infof("service shutdown", "get a signal %s. %s Server is shutting down ...", s.String(), config.App.Name)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.Close()

	for i := 0; i < wgCount; i++ {
		wg.Done()
	}

	log.Infof("service shutdown", "%s Server is exit", config.App.Name)
	log.CloseLogger()
}
