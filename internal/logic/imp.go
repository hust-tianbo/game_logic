package logic

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hust-tianbo/game_logic/config"
	"github.com/hust-tianbo/game_logic/internal/mdb"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var BoxDb *gorm.DB
var UserAssetDb *gorm.DB

func InitImp() {
	cf := config.GetConfig()
	// 随机数种子
	rand.Seed(time.Now().Unix())

	// 连接db
	var err error
	BoxDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/box?charset=utf8&parseTime=True&loc=Local", cf.DBUser, cf.DBSecret, cf.DBIP))
	if err != nil {
		panic(err)
	}

	UserAssetDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/user_asset?charset=utf8&parseTime=True&loc=Local", cf.DBUser, cf.DBSecret, cf.DBIP))
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
