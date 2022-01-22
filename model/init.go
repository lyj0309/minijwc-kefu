package model

import (
	"github.com/go-redis/redis/v8"
	dbLib "github.com/lyj0309/jwc-lib/db"
	esLib "github.com/lyj0309/jwc-lib/elastic"
	wxLib "github.com/lyj0309/jwc-lib/wx"
	"github.com/olivere/elastic/v7"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Mini     *miniprogram.MiniProgram
	Official *officialaccount.OfficialAccount
	EsClient *elastic.Client
	rdb      *redis.Client
	db       *gorm.DB
)

func init() {
	EsClient = esLib.NewElastic()
	rdb = dbLib.NewRedis()
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
