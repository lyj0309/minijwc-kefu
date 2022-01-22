package model

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	dbLib "github.com/lyj0309/jwc-lib/db"
	esLib "github.com/lyj0309/jwc-lib/elastic"
	wxLib "github.com/lyj0309/jwc-lib/wx"
	"github.com/olivere/elastic/v7"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/officialaccount"
	"gorm.io/gorm"
	"time"
)

var (
	Mini     *miniprogram.MiniProgram
	Official *officialaccount.OfficialAccount
	EsClient *elastic.Client
	rdb      *redis.Client
	db       *gorm.DB
)

//Donate 表结构
type Donate struct {
	ID         string
	Amount     int    `json:"amount"`
	Code       string `json:"code" gorm:"-"` //忽略此字段
	Note       string `json:"note"`
	CreateTime time.Time
	PayTime    time.Time
	User       string //openid
}

func init() {
	EsClient = esLib.NewElastic()
	rdb = dbLib.NewRedis()
	db = dbLib.NewDB()

	fmt.Println("新建表")
	err := db.AutoMigrate(Donate{})
	if err != nil {
		fmt.Println(err)
		return
	}
	//esLib.InsertCsv(EsClient)

	Mini = wxLib.NewWxMini()
	Official = wxLib.NewOfficial()
}
