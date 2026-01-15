package process

import (
	"context"
	"sportNews/conf"
	"sportNews/conf/aws"
	"sportNews/internal/service/crawler"
	"sportNews/pkg/log"
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type Process struct {
	DB *xorm.EngineGroup
}

func NDTVProcess(ctx context.Context, conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewNDTVServ(db)

	interval := time.Duration(conf.Process.News) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("NDTVProcess: receive stop signal, exit")
			return

		case t := <-ticker.C:
			log.Info("NDTVProcess: Starting NDTV crawler", zap.Time("startTime", t))
			serv.Crawler()
		}
	}

}

func BCCIProcess(ctx context.Context, conf conf.Conf, awsClient *aws.S3Client, db *xorm.EngineGroup) {
	serv := crawler.NewBCCIServ(db, awsClient)

	interval := time.Duration(conf.App.Process.Ranking) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("BCCIProcess: receive stop signal, exit")
			return

		case t := <-ticker.C:
			log.Info("BCCIProcess: Starting BCCI crawler", zap.Time("startTime", t))
			serv.Crawler()
		}
	}
}

func PictureSyncProcess(ctx context.Context, conf conf.Conf, s3Client *aws.S3Client, db *xorm.EngineGroup) {
	serv := crawler.NewStoryFile(db, s3Client)

	interval := time.Duration(conf.App.Process.SyncPicture) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("PictureSyncProcess: receive stop signal, exit")
			return

		case t := <-ticker.C:
			log.Info("PictureSyncProcess: Starting sync news cover", zap.Time("startTime", t))
			serv.StoryNewsCover()
		}
	}
}
