package mdb

import (
	"sync"
	"time"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/jinzhu/gorm"

	"github.com/hust-tianbo/game_logic/internal/model"
)

var BoxList []model.BoxInfo
var PrizeList []model.PrizeInfo
var BoxToPrizeList []model.BoxToPrize

var lock sync.Locker

func refreshBox(db *gorm.DB) error {
	var boxList []model.BoxInfo
	dbRes := db.Table(model.BoxInfoTable).Find(&boxList)

	// 如果没有查到则需要更新票据
	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[GetBoxs]query box failed:%+v", dbRes.Error)
		return dbRes.Error
	}

	BoxList = boxList
	log.Debugf("[refreshBox]query box success:%+v", BoxList)
	return nil
}

func refreshPrize(db *gorm.DB) error {
	var prizeList []model.PrizeInfo
	dbRes := db.Table(model.PrizeInfoTable).Find(&prizeList)

	// 如果没有查到则需要更新票据
	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[GetBoxs]query box failed:%+v", dbRes.Error)
		return dbRes.Error
	}

	PrizeList = prizeList
	log.Debugf("[refreshPrize]query prize success:%+v", PrizeList)
	return nil
}

func refreshBoxToPrize(db *gorm.DB) error {
	var boxToPrize []model.BoxToPrize
	dbRes := db.Table(model.BoxToPrizeTable).Find(&boxToPrize)

	// 如果没有查到则需要更新票据
	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[GetBoxs]query box failed:%+v", dbRes.Error)
		return dbRes.Error
	}

	BoxToPrizeList = boxToPrize
	log.Debugf("[refreshBoxToPrize]query box to prize success:%+v", BoxToPrizeList)
	return nil
}

func refresh(db *gorm.DB) bool {
	lock.Lock()
	defer lock.Unlock()

	var sync sync.WaitGroup
	sync.Add(3)
	var boxError error
	go func() {
		defer sync.Done()
		boxError = refreshBox(db)
	}()

	var prizeError error
	go func() {
		defer sync.Done()
		prizeError = refreshPrize(db)
	}()

	var boxPrizeError error
	go func() {
		defer sync.Done()
		boxPrizeError = refreshBoxToPrize(db)
	}()

	sync.Wait()

	var refreshResult = boxError == nil && prizeError == nil && boxPrizeError == nil
	log.Debugf("[refresh]refresh result:%+v", refreshResult)

	return refreshResult
}

func LoopRefresh(db *gorm.DB) bool {
	// 初始化
	var success = refresh(db)
	if !success {
		return false
	}

	// 异步更新
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				refresh(db)
			}
		}
	}()

	return true
}
