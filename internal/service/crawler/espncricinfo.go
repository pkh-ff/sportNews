package crawler

import (
	"net/http"
	"sportNews/internal/helper"
	crawlerModel "sportNews/internal/model/crawler"
	"sportNews/internal/repository"
	"sportNews/pkg/log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type EspncricinfoServ struct {
	Serv   *Serv
	Source string
	Domain string
}

func NewEspncricinfoServ(db *xorm.EngineGroup) *EspncricinfoServ {
	repo := repository.New(db)

	return &EspncricinfoServ{
		Serv: &Serv{
			Repo:   repo,
			Source: "espncricinfo",
		},
		Domain: "",
	}
}

func (s *EspncricinfoServ) Crawler() {
	log.Info("Crawler: Starting news crawl", zap.String("source", s.Source))
	s.Serv.CrawlerNewsTemplate(s)
}

func (s *EspncricinfoServ) detail(url string) (string, error) {
	return "", nil

}

func (s *EspncricinfoServ) list(page int) ([]crawlerModel.News, error) {
	log.Info("list: Fetching news list", zap.Int("page", page))
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8",
	}

	queryParams := map[string]string{
		"page": strconv.Itoa(page),
	}

	url := "https://www.espncricinfo.com/cricket-news"

	resp, err := helper.SendHTTPRequest(url, http.MethodGet, headers, queryParams)
	if err != nil {
		log.Error("list: Failed to send HTTP request", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("list: Failed to parse document", zap.String("url", url), zap.Error(err))
		return nil, err
	}

	data := make([]crawlerModel.News, 0)
	// 開始解析網頁元素
	doc.Find("div.ds-border-b.ds-border-line.ds-p-4").Each(func(i int, s *goquery.Selection) {
		var m crawlerModel.News
		// 描述
		desc := s.Find("div").Not("[class]")
		m.Description = desc.Text()

		// 詳情連結
		href, exists := s.Find("a").Attr("href")
		if exists {
			m.Link = href
		}

		// 標題
		title := s.Find("h2.ds-text-title-s.ds-font-bold.ds-text-typo")
		m.Title = title.Text()

		// 時間
		time := s.Find("span.ds-text-compact-xs").First()
		parts := strings.Split(time.Text(), "•")
		date := parts[0]
		t, err := helper.ConverseToTimestamp(date, "02-Jan-2006")
		if err != nil {
			log.Error("list: Failed to parse date", zap.String("date", date), zap.Error(err))
		} else {
			m.Time = t
		}

		data = append(data, m)
	})

	return data, nil
}
