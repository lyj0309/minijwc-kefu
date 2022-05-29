package model

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

type KefuMessage struct {
	ID       int
	Message  string
	Time     *time.Time
	To       ToType //谁发谁，1是客服发我，2是我发客服
	UserId   string
	Platform string //平台，公众号还是小程序
}

type ToType int
type PlatformType string

const (
	ToTypeMe             = 1
	ToTypeBack           = 2
	PlatformTypeMini     = "mini"
	PlatformTypeOfficial = "official"
)

func (m KefuMessage) storage() {
	t := time.Now()
	m.Time = &t
	db.Create(&m)
}

func parseN(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func geneGuess(ans *[]QA) (guess string) {

	guess = "猜你想问：\n"
	for i, qa := range *ans {
		if i != 0 {
			guess += strconv.Itoa(i) + ":" + qa.Question + "\n"
		}
	}
	if guess == "猜你想问：\n\n" {
		return ""
	}

	guess = guess + "\n回复对应数字查看解答"
	return
}

func checkNumMessage(msg string, id string) (res string) {
	res = ""
	num, err := strconv.Atoi(msg)
	if err != nil {
		return
	}

	r, found := cache.Get(id)
	if !found {
		return "消息已经过期，请重新提问"
	}
	if err != nil {
		logrus.Error(err)
		return
	}

	var qa []QA
	err = json.Unmarshal([]byte(r.(string)), &qa)
	if err != nil {
		return
	}

	num--
	if num > len(qa) || num < 0 {
		return
	}

	return qa[num].Question + "\n\n" + qa[num].Answer
}

func storageQuestion(ans *[]QA, id string) {
	*ans = (*ans)[1:]
	b, err := json.Marshal(ans)
	if err != nil {
		return
	}
	cache.Set(id, string(b), time.Minute*10)
}
