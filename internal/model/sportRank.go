package model

import (
	"sportNews/internal/enum"
	"time"
)

type SportRank struct {
	Id       int           `xorm:"id"`
	Type     enum.RankType `xorm:"type"`
	Data     string        `xorm:"data"`
	Date     time.Time     `xorm:"date"`
	CreateAt time.Time     `xorm:"created create_at"`
}

func (SportRank) TableName() string {
	return "sport_rank"
}

type RankDetail struct {
	Team     string `json:"team"`
	Position int    `json:"position"`
	Matches  int    `json:"matches"`
	Points   int    `json:"points"`
	Rating   int    `json:"rating"`
	Icon     string `json:"icon"`
}
