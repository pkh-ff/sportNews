package crawler

import (
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
	"sportNews/internal/model/crawler"
	"sportNews/internal/repository"
	"xorm.io/xorm"
)

type Serv struct {
	Repo   *repository.Repository
	Source string
}

func New(db *xorm.EngineGroup, source string) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo:   &repo,
		Source: source,
	}
}

type NewsCrawler interface {
	Crawler()
	List(page int) ([]crawler.News, error)
	Detail(url string) (string, error)
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
