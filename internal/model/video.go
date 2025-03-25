package model

import (
	"sportNews/internal/enum"
	"time"
)

type Video struct {
	Id          int         `xorm:"id"`
	Title       string      `xorm:"title"`
	Description string      `xorm:"description"`
	Cover       string      `xorm:"cover"`
	Link        string      `xorm:"link"`
	Status      enum.Status `xorm:"status"`
	CreateAt    time.Time   `xorm:"create_at"`
	UpdateAt    time.Time   `xorm:"update_at"`
}

func (Video) TableName() string {
	return "video"
}
