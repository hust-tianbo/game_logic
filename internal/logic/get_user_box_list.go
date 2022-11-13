package logic

import (
	"github.com/hust-tianbo/game_logic/internal/mdb"
	"github.com/hust-tianbo/game_logic/internal/model"
	"github.com/hust-tianbo/game_logic/lib"
	"github.com/hust-tianbo/go_lib/log"
	"github.com/jinzhu/gorm"
)

type GetUserBoxListReq struct {
	PersonID string `json:"personid"`

	InternalToken string `json:"internal_token"` // 如果已经有内部票据，则携带
}

type GetUserBoxListRsp struct {
	Ret int    `json:"ret"` // 错误码
	Msg string `json:"msg"` // 错误信息

	UserBox    []model.UserBox `json:"user_box"`     // 用户盒子
	AllBoxList []BoxInfo       `json:"all_box_list"` // 盒子列表
}

func getUserBox(db *gorm.DB, personId string) ([]model.UserBox, error) {
	var userBox []model.UserBox
	dbRes := db.Table(model.UserBoxTable).Where("person_id=?", personId).Find(&userBox)

	// 如果没有查到UserBox
	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[getUserBox]query box failed:%+v", dbRes.Error)
		return userBox, dbRes.Error
	}

	log.Debugf("[refreshBoxToPrize]query box to prize success:%+v", userBox)
	return userBox, nil
}

func GetUserBoxList(req *GetUserBoxListReq) GetUserBoxListRsp {
	var rsp GetUserBoxListRsp
	// 校验登录态
	if !lib.CheckToken(req.PersonID, req.InternalToken) {
		rsp.Ret = lib.RetTokenNotValid
		return rsp
	}

	rsp.UserBox = make([]model.UserBox, 0)
	rsp.AllBoxList = make([]BoxInfo, 0)
	userBoxList, getErr := getUserBox(UserAssetDb, req.PersonID)
	if getErr != nil {
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	// 所有的box信息
	allBoxInfo := make([]BoxInfo, 0)
	boxList, prizeList, boxToPrize := mdb.GetInfo()
	allBoxInfo = convertRsp(boxList, prizeList, boxToPrize)

	rsp.UserBox = userBoxList
	rsp.AllBoxList = allBoxInfo

	return rsp
}
