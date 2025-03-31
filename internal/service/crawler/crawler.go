package crawler

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("====== Crawler News Start ======")
	list, err := c.list(0)
	if err != nil {
		return
	}

	data := make([]crawler.News, 0)
	for _, v := range list {
		// 檢查新聞是否已存在
		count, err := s.Repo.GetCountByTitle(v.Title, s.Source)
		if err != nil {
			continue
		}
		if count > 0 {
			continue
		}

		time.Sleep(5 * time.Second) // 避免頻率過快
		content, err := c.detail(v.Link)
		if err != nil {
			continue
		}
		v.Content = content
		data = append(data, v)
	}

	s.storeNewsToDB(data)
	fmt.Println("====== Crawler News End ======")
}

// CrawlerRankDataTemplate
// 排行榜爬蟲模板
func (s *Serv) CrawlerRankDataTemplate(c RankCrawler) {
	fmt.Println("====== Crawler Rank Start ======")
	typeList := enum.RankTypeList()

	date := time.Now().Format("2006-01-02") // 獲取當前時間，並格式化為 "YYYY-MM-DD"
	for i, v := range typeList {
		// 檢查今天資料存不存在，如果再存在就跳過，反之才會抓取資料
		b, err := s.Repo.CheckRankDataExist(date, v)
		if err != nil {
			// TODO
			continue
		}

		if b == true {
			continue
		}

		data := c.rank(v)
		err = s.storeRankToDB(v, data)
		if err != nil {
			// TODO
			continue
		}

		// 每個類型排行榜資料只留最新5筆
		count, err := s.Repo.GetRankDataCountByType(v)
		if err != nil {
			// TODO
			continue
		}
		if count >= 5 {
			rankData, err := s.Repo.GetOldestRankDataByType(v)
			if err != nil {
				continue
			}

			err = s.Repo.DeleteRankData(rankData)
		}

		if i < len(typeList) {
			time.Sleep(5 * time.Second) // 避免頻率過快
		}
	}
	fmt.Println("====== Crawler Rank End ======")
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
			log.Errorf("storeNewsToDB(), query count error:", err)
		}
	}
}

// story rank data to db
func (s *Serv) storeRankToDB(t enum.RankType, data []model.RankDetail) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("storeRankToDB(), query count error:", err)
		return err
	}

	rank := model.SportRank{
		Type: t,
		Date: time.Now(),
		Data: string(jsonData),
	}

	err = s.Repo.InsertRank(rank)
	if err != nil {
		log.Errorf("storeNewsToDB(), query count error:", err)
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
