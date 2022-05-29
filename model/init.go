package model

import (
	dbLib "github.com/lyj0309/jwc-lib/db"
	"github.com/lyj0309/jwc-lib/lib"
	cacheLib "github.com/patrickmn/go-cache"
	"github.com/silenceper/wechat/v2"
	wxcache "github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
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

	Mini = NewWxMini()
	Official = NewOfficial()
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

func NewWxMini() *miniprogram.MiniProgram {
	wxcache.NewMemory()
	wc := wechat.NewWechat()
	mini := wc.GetMiniProgram(&miniConfig.Config{
		AppID:     lib.Config.MiniAppId,
		AppSecret: lib.Config.MiniAppSecret,
		Cache:     wxcache.NewMemory(),
	})
	return mini
}

func NewOfficial() *officialaccount.OfficialAccount {
	wc := wechat.NewWechat()

	official := wc.GetOfficialAccount(&offConfig.Config{
		AppID:          lib.Config.OffAppId,
		AppSecret:      lib.Config.OffAppSecret,
		Token:          lib.Config.OffToken,
		EncodingAESKey: lib.Config.OffEncodingAESKey,
		Cache:          wxcache.NewMemory(),
	})
	return official
}
