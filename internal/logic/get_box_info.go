package logic

import "github.com/hust-tianbo/game_logic/internal/mdb"

type GetBoxInfoReq struct {
	PersonID string `json:"personid"`
}

// 盒子信息
type BoxInfo struct {
	BoxID          int                 `json:"box_id"`
	BoxName        string              `json:"box_name"`        // 盒子名字
	BoxPic         string              `json:"box_Pic"`         // 盒子封面
	BoxDescription string              `json:"box_description"` // 盒子描述
	BoxPrizes      []PrizeInfoWithRate `json:"box_prizes"`      // 盒子中的奖品信息
}

type PrizeInfoWithRate struct {
	Info  PrizeInfo `json:"info"`  // 奖品信息
	Rate  int       `json:"rate"`  // 中奖比例
	Level int       `json:"level"` // 档位，区分1，2，3
}

// 奖品信息
type PrizeInfo struct {
	PrizeID     int    `json:"prize_id"`
	PrizeName   string `json:"prize_name"`
	BeforeMoney int    `json:"before_money"` // 划线价
	AfterMoney  int    `json:"after_money"`  // 成交价
}

type GetBoxInfoRsp struct {
	Ret     int       `json:"ret"`      // 错误码
	Msg     string    `json:"msg"`      // 错误信息
	BoxList []BoxInfo `json:"box_list"` // 盒子列表，后端排序
}

func GetBoxInfo(req GetBoxInfoReq) GetBoxInfoRsp {
	var rsp GetBoxInfoRsp
	rsp.BoxList = make([]BoxInfo, 0)
	boxList, prizeList, boxToPrize := mdb.GetInfo()
	for _, eleBox := range boxList {
		var tempBox = BoxInfo{
			BoxID:          eleBox.BoxID,
			BoxName:        eleBox.BoxName,
			BoxPic:         eleBox.BoxPic,
			BoxDescription: eleBox.BoxDescription,
			BoxPrizes:      make([]PrizeInfoWithRate, 0),
		}
		for _, eleRela := range boxToPrize {
			if eleRela.BoxID == eleBox.BoxID {
				if elePrize, exist := prizeList[eleRela.PrizeID]; exist {
					tempBox.BoxPrizes = append(tempBox.BoxPrizes, PrizeInfoWithRate{
						Info: PrizeInfo{
							PrizeID:     elePrize.PrizeID,
							PrizeName:   elePrize.PrizeName,
							BeforeMoney: elePrize.BeforeMoney,
							AfterMoney:  elePrize.AfterMoney,
						},
						Rate:  eleRela.Rate,
						Level: eleRela.Level,
					})
				}
			}
		}
		rsp.BoxList = append(rsp.BoxList, tempBox)
	}
	return rsp
}
