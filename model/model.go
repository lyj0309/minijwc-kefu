package model

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

	sendMsg := func(text string) {
		m := message.NewCustomerTextMessage(postForm.FromUserName, "智能客服开发中，回复‘人工’即可接入人工客服")
		err = Mini.GetCustomerMessage().Send(m)
		if err != nil {
			logrus.Error(err)
		}
	}

	switch postForm.MsgType {
	case "event":
		if postForm.Event == "user_enter_tempsession" {
			sendMsg("Hi~，欢迎咨询掌上教务处，智能客服小掌为您排忧解惑，请问有什么可以为您效劳呢?")
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
			return
		}

		sendMsg("智能客服开发中，回复‘人工’即可接入人工客服")
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
