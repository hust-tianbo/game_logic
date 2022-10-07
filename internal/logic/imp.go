package logic

import (
	"math/rand"
	"time"

	"github.com/hust-tianbo/game_logic/internal/mdb"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var BoxDb *gorm.DB
var UserAssetDb *gorm.DB

func init() {
	// 随机数种子
	rand.Seed(time.Now().Unix())

	// 连接db
	var err error
	BoxDb, err = gorm.Open("mysql", "root:1023564552tbd@tcp(172.16.0.8:3306)/box?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	UserAssetDb, err = gorm.Open("mysql", "root:1023564552tbd@tcp(172.16.0.8:3306)/user_asset?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	initSuccess := mdb.LoopRefresh(BoxDb)
	if !initSuccess {
		panic("fresh db info failed")
	}
	BoxDb.SingularTable(true)
	UserAssetDb.SingularTable(true)
}
