package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sportNews/internal/helper"
	"sportNews/internal/log"
	crawlerModel "sportNews/internal/model/crawler"
	"sportNews/internal/repository"
	"strings"
	"time"
	"xorm.io/xorm"
)

type NDTVServ struct {
	Serv   *Serv
	Source string
	Domain string
}

func NewMDTVServ(db *xorm.EngineGroup) *NDTVServ {
	repo := repository.New(db)

	return &NDTVServ{
		Serv: &Serv{
			Repo:   &repo,
			Source: "ndtv",
		},
		Domain: "https://sports.ndtv.com",
	}
}

func (s *NDTVServ) Crawler() {
	fmt.Println("====== Crawler Start ======")
	list, err := s.List(0)
	data := make([]crawlerModel.News, 0)
	if err != nil {
		// TODO ERROR
	} else {
		for _, v := range list {
			//先用title與source去DB檢查該新聞是否存在
			count, err := s.Serv.Repo.GetCountByTitle(v.Title, s.Source)
			if err != nil {
				// TODO
				log.Errorf("storeToDB(), query count error:", err)
			}

			// 如果存在不處理該篇新聞
			if count > 0 {
				continue
			}

			time.Sleep(5 * time.Second)
			content, err := s.Detail(v.Link)
			if err != nil {
				continue // 或者 continue
			}
			v.Content = content
			data = append(data, v)
		}
	}

	s.Serv.storeToDB(data)
	fmt.Println("====== Crawler End ======")
}

func (s *NDTVServ) List(page int) ([]crawlerModel.News, error) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8",
	}

	req, err := helper.BuildHttpRequest("https://sports.ndtv.com/cricket/news", http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}

	// 發送 HTTP 請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 檢查http code 是否為200
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	// 使用 goquery 解析 HTML
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

		time := element.Find("span.lst-a_pst_lnk").First()
		t, err := helper.ConverseToTimestamp(time.Text(), "Jan 2, 2006")
		if err != nil {
			// TODO
		} else {
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

	req, err := helper.BuildHttpRequest(url, http.MethodGet, nil, headers)
	if err != nil {
		return "", err
	}

	// 發送 HTTP 請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 檢查http code 是否為200
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	c := doc.Find("div.story__content ")

	content := strings.Builder{}
	c.Find("p").Each(func(i int, s *goquery.Selection) {
		content.WriteString("<p>")
		content.WriteString(s.Text())
		content.WriteString("</p>")
	})

	return content.String(), nil
}
