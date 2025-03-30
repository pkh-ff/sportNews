package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sportNews/internal/helper"
	crawlerModel "sportNews/internal/model/crawler"
	"strings"
	"xorm.io/xorm"
)

type NDTVServ struct {
	Serv   *Serv
	Source string
	Domain string
}

func NewNDTVServ(db *xorm.EngineGroup) *NDTVServ {
	return &NDTVServ{
		Serv:   newServ(db, "ndtv"),
		Domain: "https://sports.ndtv.com",
	}
}

func (s *NDTVServ) Crawler() {
	s.Serv.CrawlerNewsTemplate(s)
}

func (s *NDTVServ) List(page int) ([]crawlerModel.News, error) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8",
	}

	url := "https://sports.ndtv.com/cricket/news"
	resp, err := helper.SendHTTPRequest(url, http.MethodGet, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// 開始解析網頁元素
	data := make([]crawlerModel.News, 0)
	ul := doc.Find("ul#container_listing").First()
	ul.Find("div.lst-pg-a").Each(func(i int, element *goquery.Selection) {
		var m crawlerModel.News

		// 標題
		title := element.Find("a.lst-pg_ttl")
		m.Title = title.Text()

		// 詳細頁
		href, exists := title.Attr("href")
		if exists {
			m.Link = s.Domain + href
		}

		// 描述
		desc := element.Find("p.lst-pg_txt.txt_tct.txt_tct-three")
		m.Description = desc.Text()

		// 封面
		img, exists := element.Find("img.lz_img.crd_img-full").Attr("src")
		if exists {
			m.Cover = img
		}

		// 發布時間
		time := element.Find("span.lst-a_pst_lnk").First()
		t, err := helper.ConverseToTimestamp(time.Text(), "Jan 2, 2006")
		if err == nil {
			m.Time = t
		}

		data = append(data, m)
	})

	return data, nil
}

func (s *NDTVServ) Detail(url string) (string, error) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8",
	}

	resp, err := helper.SendHTTPRequest(url, http.MethodGet, headers, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	c := doc.Find("div.story__content")

	content := strings.Builder{}
	c.Find("p").Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) > 0 {
			content.WriteString("<p>")
			content.WriteString(s.Text())
			content.WriteString("</p>")
		}
	})

	return content.String(), nil
}
