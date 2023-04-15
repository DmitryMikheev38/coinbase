package models

import (
	"gorm.io/gorm"
	"time"
)

type Tick struct {
	Timestamp int64  `gorm:"column:timestamp"`
	ProductID string `json:"product_id" gorm:"column:symbol"`
	BestBid   string `json:"best_bid" gorm:"column:best_bid"`
	BestAsk   string `json:"best_ask" gorm:"column:best_ask"`
}

func (*Tick) TableName() string {
	return "ticks"
}

func (m *Tick) BeforeCreate(tx *gorm.DB) error {
	if m.Timestamp == 0 {
		m.Timestamp = time.Now().Unix()
	}
	return nil
}
