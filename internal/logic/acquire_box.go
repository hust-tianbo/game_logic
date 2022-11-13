package logic

import (
	"fmt"
	"time"

	"github.com/hust-tianbo/go_lib/log"

	"github.com/hust-tianbo/game_logic/internal/mdb"
	"github.com/hust-tianbo/game_logic/internal/model"
	"github.com/hust-tianbo/game_logic/lib"
)

type AcquireBoxReq struct {
	PersonID      string `json:"personid"`
	BoxID         int    `json:"boxid"`
	InternalToken string `json:"internal_token"` // 如果已经有内部票据，则携带
}

type AcquireBoxRsp struct {
	Ret   int    `json:"ret"`   // 错误码
	Msg   string `json:"msg"`   // 错误信息
	PayID string `json:"payid"` // 支付id
}

// 获取盒子的价格信息
func getBoxMoney(info *model.BoxInfo) int {
	return info.BoxPrice
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
		log.Errorf("[initAcquireBox]init failed:%+v,%+v,%+v,%+v", personId, userBoxID, boxID, dbRes.Error)
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

// 生成支付订单
func genePayOrder(money int, payID string) error {
	return nil
}

// 支付订单状态确认
func payOrderCheck(payID string) error {
	return nil
}

// 获取盒子的第一阶段
func AcquireBox(req *AcquireBoxReq) AcquireBoxRsp {
	var rsp AcquireBoxRsp

	// 校验登录态
	if !lib.CheckToken(req.PersonID, req.InternalToken) {
		rsp.Ret = lib.RetTokenNotValid
		return rsp
	}

	// 查询盒子的信息
	boxInfo, boxExist := mdb.GetOneBoxInfo(req.BoxID)
	if !boxExist {
		log.Errorf("[AcquireBox]init failed:%+v,%+v", boxInfo, boxExist)
		rsp.Ret = lib.RetNotFindBox
		return rsp
	}

	boxMoney := getBoxMoney(&boxInfo)

	// 生成订单id
	user_box_Id := lib.GeneID(req.PersonID)

	// 初始化订单
	initErr := initAcquireBox(req.PersonID, user_box_Id, req.BoxID)
	if initErr != nil {
		log.Errorf("[AcquireBox]initAcquireBox failed:%+v,%+v", req, initErr)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	// 生成付款单据
	payErr := genePayOrder(boxMoney, user_box_Id)
	if payErr != nil {
		log.Errorf("[AcquireBox]genePayOrder failed:%+v,%+v", req, payErr)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	rsp.Ret = lib.RetSuccess
	rsp.PayID = user_box_Id

	return rsp
}

type AcquireBoxCheckReq struct {
	PersonID      string `json:"personid"`
	BoxID         int    `json:"boxid"`
	PayID         string `json:"payid"`          // 支付id
	InternalToken string `json:"internal_token"` // 如果已经有内部票据，则携带
}

type AcquireBoxCheckRsp struct {
	Ret int    `json:"ret"` // 错误码
	Msg string `json:"msg"` // 错误信息
}

// 获取盒子的确认阶段，需要查询支付状态
func AcquireBoxCheck(req *AcquireBoxCheckReq) AcquireBoxCheckRsp {
	var rsp AcquireBoxCheckRsp

	// 校验登录态
	if !lib.CheckToken(req.PersonID, req.InternalToken) {
		rsp.Ret = lib.RetTokenNotValid
		return rsp
	}

	// 校验订单状态
	checkErr := payOrderCheck(req.PayID)
	if checkErr != nil {
		log.Errorf("[AcquireBoxCheck]check failed:%+v,%+v,%+v", req.PersonID, req.PayID, req.PayID)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	// 完成订单，获得盒子
	finishErr := finishAcquireBox(req.PersonID, req.PayID)
	if finishErr != nil {
		log.Errorf("[AcquireBox]finishAcquireBox failed:%+v,%+v", req, finishErr)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	rsp.Ret = lib.RetSuccess

	return rsp
}
