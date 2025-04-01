package crawler

import (
	"encoding/json"
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/internal/model/crawler"
	"sportNews/internal/repository"
	"sportNews/pkg/log"
	"time"
	"xorm.io/xorm"
)

type Serv struct {
	Repo   *repository.Repository
	Source string
}

// CrawlerNewsTemplate
// 新聞爬蟲模板
func (s *Serv) CrawlerNewsTemplate(c NewsCrawler) {
	log.Info("CrawlerNewsTemplate: Starting to crawl news", zap.String("source", s.Source))
	list, err := c.list(0)
	if err != nil {
		log.Error("CrawlerNewsTemplate: Failed to get news list", zap.Error(err))
		return
	}

	data := make([]crawler.News, 0)
	for _, v := range list {
		// 檢查新聞是否已存在
		count, err := s.Repo.GetCountByTitle(v.Title, s.Source)
		if err != nil {
			log.Error("CrawlerNewsTemplate: Failed to get news count by title", zap.Error(err))
			continue
		}
		if count > 0 {
			log.Info("CrawlerNewsTemplate: News already exists, skipping", zap.String("title", v.Title))
			continue
		}

		time.Sleep(5 * time.Second) // 避免頻率過快
		content, err := c.detail(v.Link)
		if err != nil {
			log.Error("CrawlerNewsTemplate: Failed to fetch news content", zap.String("link", v.Link), zap.Error(err))
			continue
		}
		v.Content = content
		data = append(data, v)
	}

	s.storeNewsToDB(data)
	log.Info("CrawlerNewsTemplate: Finished crawling news", zap.Int("newsCount", len(data)))
}

// CrawlerRankDataTemplate
// 排行榜爬蟲模板
func (s *Serv) CrawlerRankDataTemplate(c RankCrawler) {
	log.Info("CrawlerRankDataTemplate: Starting to crawl rank data")

	typeList := enum.RankTypeList()

	date := time.Now().Format("2006-01-02") // 獲取當前時間，並格式化為 "YYYY-MM-DD"
	for i, v := range typeList {
		// 檢查今天資料存不存在，如果再存在就跳過，反之才會抓取資料
		b, err := s.Repo.CheckRankDataExist(date, v)
		if err != nil {
			log.Error("CrawlerRankDataTemplate: Failed to check if rank data exists", zap.Error(err))
			continue
		}

		if b == true {
			log.Info("CrawlerRankDataTemplate: Rank data already exists, skipping", zap.String("rankType", string(v)))
			continue
		}

		data := c.rank(v)
		err = s.storeRankToDB(v, data)
		if err != nil {
			log.Error("CrawlerRankDataTemplate: Failed to store rank data", zap.Error(err))
			continue
		}

		// 每個類型排行榜資料只留最新5筆
		count, err := s.Repo.GetRankDataCountByType(v)
		if err != nil {
			log.Error("CrawlerRankDataTemplate: Failed to get rank data count by type", zap.Error(err))
			continue
		}
		if count >= 5 {
			rankData, err := s.Repo.GetOldestRankDataByType(v)
			if err != nil {
				log.Error("CrawlerRankDataTemplate: Failed to get oldest rank data", zap.Error(err))
				continue
			}

			_, err = s.Repo.DeleteRankData(rankData)
		}

		if i < len(typeList) {
			time.Sleep(5 * time.Second) // 避免頻率過快
		}
	}
	log.Info("CrawlerRankDataTemplate: Finished crawling rank data")
}

func newServ(db *xorm.EngineGroup, source string) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo:   &repo,
		Source: source,
	}
}

// story news data to db
func (s *Serv) storeNewsToDB(data []crawler.News) {

	for _, v := range data {
		news := model.News{
			Title:       v.Title,
			CoverSource: v.Cover,
			Link:        v.Link,
			Description: v.Description,
			Status:      enum.Enable,
			Source:      s.Source,
			PubDate:     v.Time,
			Content:     v.Content,
		}
		err := s.Repo.InsertNews(news)
		if err != nil {
			log.Errorf("storeNewsToDB: Failed to insert news into DB", zap.Error(err))
		}
	}
}

// story rank data to db
func (s *Serv) storeRankToDB(t enum.RankType, data []model.RankDetail) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("storeRankToDB: Failed to marshal rank data", zap.Error(err))
		return err
	}

	rank := model.SportRank{
		Type: t,
		Date: time.Now(),
		Data: string(jsonData),
	}

	err = s.Repo.InsertRank(rank)
	if err != nil {
		log.Errorf("storeRankToDB: Failed to insert rank data into DB", zap.Error(err))
		return err
	}

	return nil
}

type NewsCrawler interface {
	Crawler()
	list(page int) ([]crawler.News, error)
	detail(url string) (string, error)
}

type RankCrawler interface {
	Crawler()
	rank(t enum.RankType) []model.RankDetail
}
