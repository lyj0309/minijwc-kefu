package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	esLib "github.com/lyj0309/jwc-lib/elastic"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"strings"
)

var sendChan chan sendChanT

type sendChanT struct {
	UserId  string
	Message string
}

func WxOfficial(c *gin.Context) {

	// 传入request和responseWriter
	server := Official.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) (resp *message.Reply) {
		resp = &message.Reply{}

		sendText := func(t string) {
			resp.MsgType = message.MsgTypeText
			resp.MsgData = message.NewText(parseN(t))
		}

		switch msg.MsgType {
		case message.MsgTypeText:
			//resp.MsgData = message.NewText(msg.Content)
			if strings.Contains("人工", msg.Content) {
				sendText("联系人工客服请前往小程序-设置-联系客服，回复“人工”即可")
			} else {
				m := checkNumMessage(msg.Content, string(msg.FromUserName))
				if m != "" {
					sendText(m)
					return
				}

				ans := esLib.GetEsAns(EsClient, msg.Content)
				for _, qa := range *ans {
					fmt.Println(qa)
				}
				if len(*ans) == 0 {
					sendText(OffNoAnswer)
					return
				}

				sendText(`问题：` + (*ans)[0].Question + "\n\n回答：" + (*ans)[0].Answer)

				guess := geneGuess(ans)

				storageQuestion(ans, string(msg.FromUserName))

				sendChan <- sendChanT{
					UserId:  string(msg.FromUserName),
					Message: guess,
				}

			}
		case message.MsgTypeEvent:
			//用户订阅公众号
			if msg.Event == message.EventSubscribe {
				sendText(OffHello)
			}
		}

		return

	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logrus.Error("回复消息错误", err)
		return
	}
	//send(string(msg.FromUserName))

	//发送回复的消息
	err = server.Send()
	if err != nil {
		logrus.Error("发送消息错误", err)
		return
	}

}

func init() {
	sendChan = make(chan sendChanT)
	go func() {
		for {
			select {
			case d := <-sendChan:
				fmt.Println("通道收到", d)
				m := message.NewCustomerTextMessage(d.UserId, parseN(d.Message))
				err := Official.GetCustomerMessageManager().Send(m)
				if err != nil {
					logrus.Error("公众号客服消息发送失败", err)
				}
			}
		}
	}()
}
