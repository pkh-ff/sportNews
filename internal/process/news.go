package process

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"sportNews/conf"
	"sportNews/internal/service/crawler"
	"sportNews/pkg/log"
	"time"
	"xorm.io/xorm"
)

type Process struct {
	DB *xorm.EngineGroup
}

func NDTVProcess(conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewNDTVServ(db)
	for {
		log.Info("NDTVProcess: Starting NDTV crawler", zap.Time("startTime", time.Now()))
		serv.Crawler()

		time.Sleep(time.Duration(conf.Process.News) * time.Second)
	}

}

func BCCIProcess(conf conf.Conf, s3Client *s3.Client, db *xorm.EngineGroup) {
	serv := crawler.NewBCCIServ(db, s3Client, conf.Aws.Bucket, conf.Aws.Acl)
	for {
		log.Info("BCCIProcess: Starting BCCI crawler", zap.Time("startTime", time.Now()))
		serv.Crawler()

		time.Sleep(time.Duration(conf.App.Process.Ranking) * time.Second)
	}
}

func PictureSyncProcess(conf conf.Conf, s3Client *s3.Client, db *xorm.EngineGroup) {
	serv := crawler.NewStoryFile(db, s3Client, conf.Aws.Bucket, conf.Aws.Acl)
	for {
		log.Info("PictureSyncProcess: Starting sync news cover", zap.Time("startTime", time.Now()))
		serv.StoryNewsCover()

		time.Sleep(time.Duration(conf.App.Process.SyncPicture) * time.Second)
	}
}
