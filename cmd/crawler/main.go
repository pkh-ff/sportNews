package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sportNews/conf"
	"sportNews/conf/aws"
	"sportNews/conf/database"
	"sportNews/internal/assets"
	"sportNews/internal/process"
	"sportNews/pkg/log"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	confPath string // config path
	wg       sync.WaitGroup
)

func main() {
	flag.StringVar(&confPath, "c", "process.conf.yaml", "default config path")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := conf.New(confPath)
	if err != nil {
		fmt.Printf("setting conf %v\n", err)
		return
	}

	// set log config
	log.InitLogger(config.App.Debug)

	log.Infof("setting conf", "conf:%v", config)
	log.Info("app info", zap.String("Project", config.App.Name))

	// setting database
	db, err := database.New(config.App.Debug, config.DB)
	if err != nil {
		log.Error("init database", zap.Error(err))
		return
	}

	assets.Setup(config.Assets)

	s3Client, err := aws.New(config.Aws)

	crawlerProcess(ctx, *config, db, s3Client)

	shutdown(cancel, config, db)
}

func crawlerProcess(ctx context.Context, conf conf.Conf, db *xorm.EngineGroup, s3Client *s3.Client) {
	// NDTV
	wg.Add(1)
	go func() {
		defer wg.Done()
		process.NDTVProcess(ctx, conf.App, db)
	}()

	// BCCI
	wg.Add(1)
	go func() {
		defer wg.Done()
		process.BCCIProcess(ctx, conf, s3Client, db)
	}()

	// picture
	wg.Add(1)
	go func() {
		defer wg.Done()
		process.PictureSyncProcess(ctx, conf, s3Client, db)
	}()
}

func shutdown(cancel context.CancelFunc, config *conf.Conf, db *xorm.EngineGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	s := <-quit
	log.Infof("service shutdown", "get a signal %s. %s Server is shutting down ...", s.String(), config.App.Name)

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("all crawler goroutines exited")
	case <-time.After(5 * time.Second):
		log.Warn("timeout waiting for goroutines, forcing shutdown")
	}

	if err := db.Close(); err != nil {
		log.Error("close database error", zap.Error(err))
	}

	log.Infof("service shutdown", "%s Server is exit", config.App.Name)
	log.CloseLogger()
}
