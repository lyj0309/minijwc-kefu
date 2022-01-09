package main

import (
	"github.com/gin-gonic/gin"
	"minijwc-kefu/model"
)

func main() {
	model.Init()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.String(200, c.Query("echostr")) })
	r.POST("/", model.Kefu) //微信
	err := r.Run(":8887")
	if err != nil {
		return
	}
}
