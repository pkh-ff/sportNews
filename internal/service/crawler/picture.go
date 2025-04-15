package crawler

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"sportNews/internal/assets"
	"sportNews/internal/helper"
	"sportNews/internal/repository"
	"sportNews/pkg/log"
	"strconv"
	"xorm.io/xorm"
)

type StoryFile struct {
	Repo     *repository.Repository
	S3Client *s3.Client
	Bucket   string
	Acl      bool
}

func NewStoryFile(db *xorm.EngineGroup, s3Client *s3.Client, bucket string, acl bool) *StoryFile {
	repo := repository.New(db)

	return &StoryFile{
		Repo:     &repo,
		S3Client: s3Client,
		Bucket:   bucket,
		Acl:      acl,
	}
}

// StoryNewsCover
// 同步DB中新聞封面
func (s *StoryFile) StoryNewsCover() {
	log.Info("StoryNewsCover: Starting sync news cover to s3")
	s.SyncSourceNewsCover()
	s.SyncCustomNewsCover()
}

// SyncSourceNewsCover
// 同步&轉存新聞原始封面圖片
func (s *StoryFile) SyncSourceNewsCover() {
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
		objectKey := helper.NewsCoverPrefix + filename

		err = helper.UploadToS3(s.S3Client, imgData, s.Bucket, objectKey, contentType, s.Acl)
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
		break
	}
}

func (s *StoryFile) SyncCustomNewsCover() {
	// 取得沒有自定義封面圖片的新聞
	newsList, err := s.Repo.GetNoCoverCustomNews()
	if err != nil {
		log.Error("SyncCustomNewsCover, failed to get no cover custom news", zap.Error(err))
		return
	}

	// 取得最後更新自定義封面圖片新聞的圖片位址
	news, err := s.Repo.GetLastUpdateCoverCustomNews()
	if err != nil {
		log.Error("SyncCustomNewsCover, failed to get last updated custom cover", zap.Error(err))
		return
	}

	// 從圖片位址拿出編號
	num := 0
	// 檢查是否有資料
	if news.CoverCustom != "" {
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(news.CoverCustom, -1)

		if len(matches) > 0 {
			numStr := matches[len(matches)-1]
			num, err = strconv.Atoi(numStr)
			if err != nil {
				log.Error("SyncCustomNewsCover, failed to convert image number", zap.String("numStr", numStr), zap.Error(err))
				return
			}
		}
	}

	for _, v := range newsList {
		// 遞增圖片序號
		num++
		url := assets.FullAssetsPath(getCustomCoverPath(num))

		// 檢查檔案是否存在
		exist, err := helper.FileURLExists(url)
		if err != nil {
			log.Warn("SyncCustomNewsCover, failed to check file existence", zap.String("url", url), zap.Error(err))
			continue
		}

		if !exist {
			log.Warn("SyncCustomNewsCover, file not found, resetting num", zap.String("url", url))
			num = 1 // 重置封面圖片計數
		}

		// 更新新聞資料
		v.CoverCustom = getCustomCoverPath(num)
		err = s.Repo.UpdateNews(v)
		if err != nil {
			log.Error("SyncCustomNewsCover, failed to update news", zap.Int32("newsId", v.Id), zap.Error(err))
			continue
		}
		log.Info("SyncCustomNewsCover, successfully updated news", zap.Int32("newsId", v.Id), zap.String("cover", v.CoverCustom))
	}

}

func getCustomCoverPath(num int) string {
	return "/" + helper.NewsCustomCoverPrefix + "sport" + strconv.Itoa(num) + ".png"
}
