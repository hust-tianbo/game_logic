package logic

import (
	"fmt"
	"time"

	"github.com/hust-tianbo/go_lib/log"

	"github.com/hust-tianbo/game_logic/internal/model"
)

type AcquireBoxReq struct {
	PersonID string `json:"personid"`
}

type AcquireBoxRsp struct {
	Ret int    `json:"ret"` // 错误码
	Msg string `json:"msg"` // 错误信息
}

// 获取盒子的价格信息
func getBoxMoney() int {
	return 0
}

// 初始化一条盒子购买记录
func initAcquireBox(personId string, userBoxID string, boxID int) error {
	now := time.Now()
	dbRes := UserAssetDb.Table(model.UserBoxTable).Create(&model.UserBox{
		PersonID:   personId,
		UserBoxID:  userBoxID,
		BoxID:      boxID,
		CTime:      now,
		MTime:      now,
		BuyChannel: model.BoxBuyChannelCash,
		Status:     model.BoxStatusInit,
	})
	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[initAcquireBox]init failed:%+v,%+v,%+v", personId, userBoxID, boxID)
		return fmt.Errorf("init failed")
	}
	return nil
}

func finishAcquireBox(personId string, userBoxID string) error {
	nowTime := time.Now()

	// 需要更新状态，到已经获取成功的状态
	dbRes := UserAssetDb.Table(model.UserBoxTable).Where(&model.UserBox{PersonID: personId, UserBoxID: userBoxID}).Update(map[string]interface{}{
		"status": model.BoxStatusSuccess, "m_time": nowTime})

	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[finishAcquireBox]acquire box success:%+v,%+v", personId, userBoxID)
		return fmt.Errorf("finish acquire failed")
	}
	return nil
}

func AcquireBox(req *AcquireBoxReq) AcquireBoxRsp {
	//boxMoney := getBoxMoney()

	return AcquireBoxRsp{}
}
