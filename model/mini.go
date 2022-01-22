package model

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	esLib "github.com/lyj0309/jwc-lib/elastic"
	"github.com/silenceper/wechat/v2/miniprogram/message"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func Kefu(c *gin.Context) {

	var postForm WxPostForm
	err := c.ShouldBindJSON(&postForm)
	if err != nil {
		logrus.Error(err)
		return
	}

	sendText := func(text string) {
		m := message.NewCustomerTextMessage(postForm.FromUserName, text)
		err = Mini.GetCustomerMessage().Send(m)
		if err != nil {
			logrus.Error(err)
		}
	}

	switch postForm.MsgType {
	case "event":
		if postForm.Event == "user_enter_tempsession" {
			sendText(Hello)
		}

	case "text":
		if strings.Contains("人工", postForm.Content) {
			repTextMsg := TransStuff{
				ToUserName:   postForm.FromUserName,
				FromUserName: postForm.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "transfer_customer_service",
			}
			c.JSON(200, repTextMsg)
			sendText("已经转接人工客服，请耐心等待人工接入~")
		} else {
			m := checkNumMessage(postForm.Content, postForm.FromUserName)
			if m != "" {
				sendText(m)
				return
			}

			ans := esLib.GetEsAns(EsClient, postForm.Content)
			fmt.Println(ans)
			if len(*ans) == 0 {
				sendText(NoAnswer)
				return
			}

			sendText(`问题：` + (*ans)[0].Question + "\n\n回答：" + (*ans)[0].Answer)

			guess := geneGuess(ans)

			sendText(guess)

			storageQuestion(ans, postForm.FromUserName)
		}
	}

	b, err := json.Marshal(&postForm)
	fmt.Println(string(b))

}

// TransStuff 转发人工
type TransStuff struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
}

type WxPostForm struct {
	Content      string  `json:"Content"`
	CreateTime   int64   `json:"CreateTime"`
	FromUserName string  `json:"FromUserName"`
	MsgID        float64 `json:"MsgId"`
	MsgType      string  `json:"MsgType"`
	ToUserName   string  `json:"ToUserName"`
	Event        string  `json:"Event"`
}
