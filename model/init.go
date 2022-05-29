package model

import (
	dbLib "github.com/lyj0309/jwc-lib/db"
	wxLib "github.com/lyj0309/jwc-lib/wx"
	cacheLib "github.com/patrickmn/go-cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var (
	Mini     *miniprogram.MiniProgram
	Official *officialaccount.OfficialAccount
	db       *gorm.DB
	cache    *cacheLib.Cache
)

func init() {
	cache = cacheLib.New(5*time.Minute, 10*time.Minute)
	db = dbLib.NewDB()

	//esLib.InsertCsv(EsClient)
	checkAndCreateTable(&KefuMessage{})

	Mini = wxLib.NewWxMini()
	Official = wxLib.NewOfficial()
}
func checkAndCreateTable(table interface{}) {
	if !db.Migrator().HasTable(table) {
		logrus.Info("创建数据表")
		err := db.AutoMigrate(table)
		if err != nil {
			logrus.Fatal("数据表生成错误", err)
		}
	}
}
