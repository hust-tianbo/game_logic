package model

import "time"

const UserPrizeTable = "user_prize"

const (
	PrizeStatusInit    = 1 // 初始化
	PrizeStatusSuccess = 2 // 购买成功
	PrizeStatusConsume = 3 // 已经消费成功
)

const (
	PrizeBuyChannelCash = 1 // 现金
	PrizeBuyChannelGold = 2 // 金币
	PrizeBuyChannelBox  = 3 // 盒子抽取
)

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
