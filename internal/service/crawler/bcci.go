package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sportNews/internal/enum"
	"sportNews/internal/helper"
	"sportNews/internal/model"
	"strconv"
	"strings"
	"time"
	"xorm.io/xorm"
)

type BCCIServ struct {
	Serv   *Serv
	Source string
	Domain string
}

func NewBCCIServ(db *xorm.EngineGroup) *BCCIServ {
	return &BCCIServ{
		Serv:   newServ(db, "bcci"),
		Domain: "https://sports.ndtv.com",
	}
}

func (s *BCCIServ) Crawler() {
	fmt.Println("====== BCCI rank Crawler Start ======")
	typeList := enum.RankTypeList()
	for i, v := range typeList {
		// TODO 先檢查今天資料存不存在如果再存在就跳過，反之才會抓取資料
		data := s.rank(v)
		s.Serv.storeRankToDB(v, data)
		if i < len(typeList) {
			time.Sleep(5 * time.Second) // 避免頻率過快
		}
	}
	fmt.Println("====== BCCI rank Crawler End ======")
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
		// TODO LOG
		fmt.Println(err)
		//return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// TODO LOG
		fmt.Println(err)
	}

	//fmt.Println(doc.Html())
	data := make([]model.RankDetail, 0)

	// rank 1
	r1 := doc.Find("div.team-ranking-wrapper").First()
	r1Model := getRankTopData(r1)
	data = append(data, r1Model)

	// rank 2 ~
	rOther := doc.Find("div.ranking-data-table.table-responsive").First()
	ranks := getRankOtherData(rOther)
	data = append(data, ranks...)
	fmt.Println(data)
	return data
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
		position = 0 // 如果轉換失敗，預設為 0
	}

	// 國家名稱
	country := s.Find(".rank-number span").Text()

	// matches, points, rating
	var matches, points, rating int
	s.Find(".ranking-top-table table tbody tr td").Each(func(i int, s *goquery.Selection) {
		value, err := strconv.Atoi(s.Find("p").Text())
		if err != nil {
			// TODO LOG
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
			// TODO LOG
			position = 0
		}
		detail.Position = position

		matchesStr := tr.Find("td p").Eq(0).Text()
		matches, err := strconv.Atoi(matchesStr)
		if err != nil {
			matches = 0
			// TODO LOG
		}
		detail.Matches = matches

		pointsStr := tr.Find("td p").Eq(1).Text()
		points, err := strconv.Atoi(pointsStr)
		if err != nil {
			points = 0
			// TODO LOG
		}
		detail.Points = points

		// ratings
		ratingsStr := tr.Find("td p").Eq(2).Text()
		ratings, err := strconv.Atoi(ratingsStr)
		if err != nil {
			ratings = 0
			// TODO LOG
		}
		detail.Rating = ratings

		data = append(data, detail)
	})

	return data
}
