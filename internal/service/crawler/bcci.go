package crawler

import (
	"fmt"
	"net/http"
	"sportNews/internal/assets"
	"sportNews/internal/enum"
	"sportNews/internal/helper"
	"sportNews/internal/model"
	"sportNews/pkg/log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type BCCIServ struct {
	Serv     *Serv
	Source   string
	Domain   string
	S3Client *s3.Client
	Bucket   string
	Acl      bool
}

func NewBCCIServ(db *xorm.EngineGroup, s3Client *s3.Client, bucket string, acl bool) *BCCIServ {
	return &BCCIServ{
		Serv:     newServ(db, "bcci"),
		Domain:   "https://sports.ndtv.com",
		S3Client: s3Client,
		Bucket:   bucket,
		Acl:      acl,
	}
}

func (s *BCCIServ) Crawler() {
	log.Info("Crawler: Starting ranking data crawling")
	s.Serv.CrawlerRankDataTemplate(s)
}

func (s *BCCIServ) rank(t enum.RankType) []model.RankDetail {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8",
	}

	url := fmt.Sprintf("https://www.bcci.tv/international/men/rankings/%s", t)

	resp, err := helper.SendHTTPRequest(url, http.MethodGet, headers, nil)
	if err != nil {
		log.Error("rank: Failed to send HTTP request", zap.String("url", url), zap.Error(err))
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("BCCIServ.rank: Failed to parse HTML document", zap.Error(err))
		return nil
	}

	data := make([]model.RankDetail, 0)

	// rank 1
	r1 := doc.Find("div.team-ranking-wrapper").First()
	r1Model := getRankTopData(r1)
	data = append(data, r1Model)

	// rank 2 ~
	rOther := doc.Find("div.ranking-data-table.table-responsive").First()
	ranks := getRankOtherData(rOther)
	data = append(data, ranks...)

	s.replaceTeamIcon(&data)

	return data
}

func (s *BCCIServ) replaceTeamIcon(data *[]model.RankDetail) {
	for i := 0; i < len(*data); i++ {
		// 檢查檔案是否存在
		ext, err := helper.GetFileExtensionFromURL((*data)[i].Icon)
		if err != nil {
			log.Error("replaceTeamIcon, Unable to get file extension", zap.String("url", (*data)[i].Icon), zap.Error(err))
			continue
		}

		url := assets.FullAssetsPath(getTeamIconPath((*data)[i].Team, ext))

		exist, err := helper.FileURLExists(url)
		if err != nil {
			log.Warn("replaceTeamIcon, failed to check file existence", zap.String("url", url), zap.Error(err))
			continue
		}

		objectKey := helper.TeamIconPrefix + (*data)[i].Team + ext
		// 檢查檔案是否存在serve上，反判斷後續要不要處理檔案上傳
		if exist {
			log.Info("replaceTeamIcon, file is existence", zap.String("url", url))
			(*data)[i].Icon = "/" + objectKey
			continue
		}

		//
		// update file to s3
		//
		imgData, err := helper.DownloadFileFromUrl((*data)[i].Icon)
		if err != nil {
			log.Error("StoryNewsCover, Unable to get file name", zap.Error(err))
			continue
		}

		err = helper.UploadToS3(s.S3Client, imgData, s.Bucket, objectKey, http.DetectContentType(imgData), s.Acl)
		if err != nil {
			log.Error("StoryNewsCover, upload fil to s3 fail", zap.String("url", (*data)[i].Icon), zap.Error(err))
			continue
		}

		(*data)[i].Icon = "/" + objectKey
	}
}

// 取得第一名資料
func getRankTopData(s *goquery.Selection) model.RankDetail {
	// icon
	img, exists := s.Find("img").Attr("data-src")
	if !exists {
		img = ""
	}

	// 取得排名數字並移除 "#"
	rank := strings.TrimPrefix(s.Find(".rank-number h1").Text(), "#")
	position, err := strconv.Atoi(rank)
	if err != nil {
		log.Warn("getRankTopData: Failed to parse rank number", zap.String("rank", rank), zap.Error(err))
		position = 0 // 如果轉換失敗，預設為 0
	}

	// 國家名稱
	country := s.Find(".rank-number span").Text()

	// matches, points, rating
	var matches, points, rating int
	s.Find(".ranking-top-table table tbody tr td").Each(func(i int, s *goquery.Selection) {
		value, err := strconv.Atoi(s.Find("p").Text())
		if err != nil {
			log.Warn("getRankTopData: Failed to parse integer value", zap.String("value", s.Text()), zap.Error(err))
			value = 0
		}

		switch i {
		case 0:
			matches = value
		case 1:
			points = value
		case 2:
			rating = value
		}
	})

	return model.RankDetail{
		Icon:     img,
		Team:     country,
		Position: position,
		Matches:  matches,
		Points:   points,
		Rating:   rating,
	}
}

func getRankOtherData(s *goquery.Selection) []model.RankDetail {
	data := make([]model.RankDetail, 0)

	s.Find("table.table tbody tr").Each(func(i int, tr *goquery.Selection) {
		detail := model.RankDetail{}

		// icon
		img, exist := tr.Find("td img").Attr("data-src")
		if exist {
			detail.Icon = img
		}

		// team
		detail.Team = tr.Find("td h6").Text()

		// position
		positionStr := tr.Find("td h5").Text()
		position, err := strconv.Atoi(strings.TrimSpace(positionStr))
		if err != nil {
			log.Warn("getRankOtherData: Failed to parse position", zap.String("position", positionStr), zap.Error(err))
			position = 0
		}
		detail.Position = position

		matchesStr := tr.Find("td p").Eq(0).Text()
		matches, err := strconv.Atoi(matchesStr)
		if err != nil {
			log.Warn("getRankOtherData: Failed to parse matches", zap.String("matches", matchesStr), zap.Error(err))
			matches = 0
		}
		detail.Matches = matches

		pointsStr := tr.Find("td p").Eq(1).Text()
		points, err := strconv.Atoi(pointsStr)
		if err != nil {
			log.Warn("getRankOtherData: Failed to parse points", zap.String("points", pointsStr), zap.Error(err))
			points = 0
		}
		detail.Points = points

		// ratings
		ratingsStr := tr.Find("td p").Eq(2).Text()
		ratings, err := strconv.Atoi(ratingsStr)
		if err != nil {
			log.Warn("getRankOtherData: Failed to parse ratings", zap.String("ratings", ratingsStr), zap.Error(err))
			ratings = 0
		}
		detail.Rating = ratings

		data = append(data, detail)
	})

	return data
}

func getTeamIconPath(fileName, ext string) string {
	return "/" + helper.TeamIconPrefix + fileName + ext
}
