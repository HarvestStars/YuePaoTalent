package main

import (
	"github.com/HarvestStars/YuePaoTalent/util/server"
	"github.com/gin-gonic/gin"
)

// Register 注册账户
func Register(c *gin.Context) {

}

// LogIn 登录已有账户
func LogIn(c *gin.Context) {

}

// StartChat 利用websocket开启即时通讯
func StartChat(c *gin.Context) {
	server.ServeWs(hub, c.Writer, c.Request)
}
