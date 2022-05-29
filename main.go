package main

import (
	"github.com/gin-gonic/gin"
	"minijwc-kefu/model"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.String(200, c.Query("echostr")) })
	r.GET("/official", func(c *gin.Context) { c.String(200, c.Query("echostr")) })
	r.POST("/", model.Kefu)               //客服
	r.POST("/official", model.WxOfficial) //公众号
	err := r.Run(":9000")
	if err != nil {
		return
	}
}
