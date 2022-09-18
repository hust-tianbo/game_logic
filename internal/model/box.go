package model

const BoxInfoTable string = "box_info"
const BoxToPrizeTable string = "box_to_prize"
const PrizeInfoTable string = "prize_info"

// 盒子信息
type BoxInfo struct {
	BoxID          int    `gorm:"column:box_id"`
	BoxName        string `gorm:"column:box_name"`        // 盒子名字
	BoxPic         string `gorm:"column:box_pic"`         // 盒子封面
	BoxDescription string `gorm:"column:box_description"` // 盒子描述
}

// 奖品信息
type PrizeInfo struct {
	PrizeID     int    `gorm:"column:prize_id"`
	PrizeName   string `gorm:"column:prize_name"`
	BeforeMoney int    `gorm:"column:before_money"` // 划线价
	AfterMoney  int    `gorm:"column:after_money"`  // 成交价
}

type BoxToPrize struct {
	BoxID   int `gorm:"column:box_id"`
	PrizeID int `gorm:"column:prize_id"`
	Rate    int `gorm:"column:rate"`  // 中奖比例
	Level   int `gorm:"column:level"` // 档位，区分1，2，3
}
