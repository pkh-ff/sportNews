package model

import "time"

type BaseResp struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"totalCount"`
	TotalPage  int         `json:"totalPage"`
}

type NewsResp struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Cover       string    `json:"cover"`
	CoverSource string    `json:"coverSource"`
	PubDate     time.Time `json:"pubDate"`
}

type VideoResp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Cover       string `json:"cover"`
}
