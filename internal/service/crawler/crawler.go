package crawler

import (
	"fmt"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
	"sportNews/internal/model/crawler"
	"sportNews/internal/repository"
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
	fmt.Println("====== Crawler Start ======")
	list, err := c.List(0)
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
		content, err := c.Detail(v.Link)
		if err != nil {
			continue
		}
		v.Content = content
		data = append(data, v)
	}

	s.storeToDB(data)
	fmt.Println("====== Crawler End ======")
}

func newServ(db *xorm.EngineGroup, source string) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo:   &repo,
		Source: source,
	}
}

// story to db
func (s *Serv) storeToDB(data []crawler.News) {
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
			log.Errorf("storeToDB(), query count error:", err)
		}
	}
}

type NewsCrawler interface {
	Crawler()
	List(page int) ([]crawler.News, error)
	Detail(url string) (string, error)
}
