package model

import "time"

// 用户奖品信息
type UserPrize struct {
	PersonID    string    `gorm:"column:person_id"`
	UserPrizeID string    `gorm:"column:user_prize_id"`
	PrizeID     int       `gorm:"column:prize_id"`
	CTime       time.Time `gorm:"column:c_time"`
	MTime       time.Time `gorm:"column:m_time"`
	BuyChannel  int       `gorm:"column:buy_channel"`
	Status      int       `gorm:"column:status"`
}
