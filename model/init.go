package model

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	Mini *miniprogram.MiniProgram
)

func Init() {
	time.Local = time.FixedZone("CST", 8*3600) // 东八
	logrus.Info("设置时区", time.Now())

	wc := wechat.NewWechat()

	AppID := os.Getenv("AppID")
	AppSecret := os.Getenv("AppSecret")
	if AppID == "" || AppSecret == "" {
		logrus.Fatal("请设置AppID和AppSecret")
	}

	Mini = wc.GetMiniProgram(&miniConfig.Config{
		AppID:     AppID,
		AppSecret: AppSecret,
		Cache:     cache.NewMemory(),
	})
}
