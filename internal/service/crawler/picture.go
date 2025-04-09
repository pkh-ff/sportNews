package crawler

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"net/http"
	"sportNews/internal/helper"
	"sportNews/internal/repository"
	"sportNews/pkg/log"
	"xorm.io/xorm"
)

type StoryFile struct {
	Repo     *repository.Repository
	S3Client *s3.Client
	Bucket   string
}

func NewStoryFile(db *xorm.EngineGroup, s3Client *s3.Client, bucket string) *StoryFile {
	repo := repository.New(db)

	return &StoryFile{
		Repo:     &repo,
		S3Client: s3Client,
		Bucket:   bucket,
	}
}

// StoryNewsCover
// 同步DB中新聞封面
func (s *StoryFile) StoryNewsCover() {
	log.Info("StoryNewsCover: Starting sync news cover to s3")
	// 取得沒有封面的新聞
	news, err := s.Repo.GetNoCoverNews()
	if err != nil {
		log.Error("StoryNewsCover, Get News for no cover fail", zap.Error(err))
		return
	}

	for _, v := range news {
		filename, err := helper.GetFileNameFromURL(v.CoverSource)
		if err != nil {
			log.Error("StoryNewsCover, Unable to get file name", zap.Error(err))
			continue
		}

		imgData, err := helper.DownloadFileFromUrl(v.CoverSource)
		if err != nil {
			log.Error("StoryNewsCover, Unable download get file", zap.Error(err))
			continue
		}
		contentType := http.DetectContentType(imgData)
		objectKey := "sport-news/news/" + filename

		err = helper.UploadToS3(s.S3Client, imgData, s.Bucket, objectKey, contentType)
		if err != nil {
			log.Error("StoryNewsCover, upload fil to s3 fail", zap.String("url", v.CoverSource), zap.Error(err))
			continue
		}

		v.Cover = "/" + objectKey
		err = s.Repo.UpdateNews(v)
		if err != nil {
			log.Error("StoryNewsCover, update news data fail", zap.Int32("news", v.Id), zap.Error(err))
			continue
		}
	}
}
