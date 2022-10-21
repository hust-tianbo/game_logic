package logic

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/hust-tianbo/game_logic/internal/mdb"
	"github.com/hust-tianbo/game_logic/internal/model"
	"github.com/hust-tianbo/game_logic/lib"
	"github.com/hust-tianbo/go_lib/log"
	"github.com/jinzhu/gorm"
)

type AcquirePrizeReq struct {
	PersonID string `json:"personid"`
	BoxID    int    `json:"boxid"`
	PayID    string `json:"payid"`
}

type AcquirePrizeRsp struct {
	Ret   int       `json:"ret"`   // 错误码
	Msg   string    `json:"msg"`   // 错误信息
	Prize PrizeInfo `json:"prize"` // 奖励信息
}

func genePrize(boxID int) (info PrizeInfo, err error) {
	var tempBox BoxInfo
	boxList, prizeList, boxToPrize := mdb.GetInfo()
	for _, eleBox := range boxList {
		// 只关注该boxID
		if eleBox.BoxID != boxID {
			continue
		}
		tempBox = BoxInfo{
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
		sort.Slice(tempBox.BoxPrizes, func(i, j int) bool {
			return tempBox.BoxPrizes[i].Level < tempBox.BoxPrizes[j].Level
		})
		break
	}

	totalRate := 0
	for _, ele := range tempBox.BoxPrizes {
		totalRate += ele.Rate
	}

	tempTotal := 0
	randInt := rand.Intn(totalRate)
	for _, ele := range tempBox.BoxPrizes {
		tempTotal += ele.Rate
		if randInt < tempTotal {
			// 小于总和的话，说明概率落在这个物品区间
			return ele.Info, nil
		}
	}
	return PrizeInfo{}, fmt.Errorf("not find")
}

func consumeBox(tx *gorm.DB, personId string, userBoxID string) error {
	nowTime := time.Now()

	// 需要更新状态，到已经获取成功的状态
	dbRes := tx.Table(model.UserBoxTable).Where(&model.UserBox{PersonID: personId, UserBoxID: userBoxID, Status: model.BoxStatusSuccess}).Update(map[string]interface{}{
		"status": model.BoxStatusConsume, "m_time": nowTime})

	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[consumeBox]consume box failed:%+v,%+v,%+v", personId, userBoxID, dbRes.Error)
		return fmt.Errorf("finish acquire failed")
	}
	return nil
}

func acquirePrize(tx *gorm.DB, personId string, userBoxID string, prizeID int) error {
	now := time.Now()
	dbRes := tx.Table(model.UserPrizeTable).Create(&model.UserPrize{
		PersonID:    personId,
		UserPrizeID: userBoxID,
		PrizeID:     prizeID,
		CTime:       now,
		MTime:       now,
		BuyChannel:  model.PrizeBuyChannelBox,
		Status:      model.PrizeStatusSuccess,
	})
	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[acquirePrize]init failed:%+v,%+v,%+v,%+v", personId, userBoxID, prizeID, dbRes.Error)
		return fmt.Errorf("init failed")
	}
	return nil
}

// 用户使用box
func AcquirePrize(req *AcquirePrizeReq) AcquirePrizeRsp {
	var rsp AcquirePrizeRsp

	if req.PayID == "" {
		log.Errorf("[AcquirePrize]param invalid:%+v", req)
		rsp.Ret = lib.RetParamError
		return rsp
	}

	// 生成一个奖励物品
	prizeInfo, geneErr := genePrize(req.BoxID)
	if geneErr != nil {
		log.Errorf("[AcquirePrize]genePrize failed:%+v,%+v", req, geneErr)
		rsp.Ret = lib.RetNotFindPrize
		return rsp
	}

	// 消耗盒子同时写入奖励物品
	tx := UserAssetDb.Begin()
	if consumeBox(tx, req.PersonID, req.PayID) != nil ||
		acquirePrize(tx, req.PersonID, req.PayID, prizeInfo.PrizeID) != nil {
		log.Errorf("[AcquirePrize]second failed:%+v", req)
		tx.Rollback()
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	if tx.Commit().Error != nil {
		log.Errorf("[AcquirePrize]third failed:%+v", req)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	rsp.Ret = lib.RetSuccess
	rsp.Prize = prizeInfo
	return rsp
}
