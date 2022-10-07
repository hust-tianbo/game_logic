package model

import "time"

const UserBoxTable = "user_box"
const (
	BoxStatusInit    = 1 // 初始化
	BoxStatusSuccess = 2 // 购买成功
	BoxStatusConsume = 3 // 已经消费成功
)

const (
	BoxBuyChannelCash = 1 // 现金
	BoxBuyChannelGold = 2 // 金币
)

// 用户盒子信息
type UserBox struct {
	PersonID   string    `gorm:"column:person_id"`
	UserBoxID  string    `gorm:"column:user_box_id"`
	BoxID      int       `gorm:"column:box_id"`
	CTime      time.Time `gorm:"column:c_time"`
	MTime      time.Time `gorm:"column:m_time"`
	BuyChannel int       `gorm:"column:buy_channel"`
	Status     int       `gorm:"column:status"`
}
