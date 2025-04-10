package api

import "time"

type NewsPageResp struct {
	Records    []NewsList `json:"records"`
	TotalCount int64      `json:"totalCount"`
	TotalPage  int        `json:"totalPage"`
}

type NewsList struct {
	Id          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	CoverSource string    `json:"coverSource"`
	CoverCustom string    `json:"coverCustom"`
	PubDate     time.Time `json:"pubDate"`
}

type NewsDetail struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	CoverSource string    `json:"coverSource"`
	CoverCustom string    `json:"coverCustom"`
	Content     string    `json:"content"`
	PubDate     time.Time `json:"pubDate"`
}
