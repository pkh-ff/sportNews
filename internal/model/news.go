package model

import (
	"sportNews/internal/enum"
	"time"
)

type News struct {
	Id          int         `xorm:"id"`
	Title       string      `xorm:"title"`
	Description string      `xorm:"description"`
	Cover       string      `xorm:"cover"`
	CoverSource string      `xorm:"cover_source"`
	Link        string      `xorm:"link"`
	Content     string      `xorm:"content"`
	Source      string      `xomr:"source"`
	Status      enum.Status `xorm:"status"`
	PubDate     time.Time   `xorm:"pub_date"`
	CreateAt    time.Time   `xorm:"created create_at"`
	UpdateAt    time.Time   `xorm:"updated update_at"`
}

func (News) TableName() string {
	return "news"
}
