package model

import (
	"time"
)

type SportRank struct {
	Id          int       `xorm:"id"`
	Type        string    `xorm:"type"`
	Description string    `xorm:"description"`
	Data        string    `xorm:"data"`
	Date        time.Time `xorm:"date"`
	CreateAt    time.Time `xorm:"create_at"`
}

type RankDetail struct {
	Team     string `json:"team"`
	Position int    `json:"position"`
	Matches  int    `json:"matches"`
	Points   int    `json:"points"`
	Rating   int    `json:"rating"`
	Icon     string `json:"icon"`
}

func (SportRank) TableName() string {
	return "sport_rank"
}
