package main

import (
	"net/http"

	"github.com/HarvestStars/YuePaoTalent/db"
	"github.com/HarvestStars/YuePaoTalent/protocol"
	"github.com/HarvestStars/YuePaoTalent/util/common"
	"github.com/HarvestStars/YuePaoTalent/util/server"
	"github.com/gin-gonic/gin"
)

// SignUp 注册账户
func SignUp(c *gin.Context) {
	var req protocol.UserReq
	c.BindJSON(&req)
	// 判断是否已经存在
	var userInfo protocol.UserInfo
	var count int
	db.DataBase.Model(&userInfo).Where("user_name = ?", req.UserName).Count(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": "", "error": "账户已经存在"})
		return
	}

	// 注册
	userInfo.UserName = req.UserName
	userInfo.Salt = common.GetRandomBoth(4)
	userInfo.PassWord = common.Sha1En(req.PassWord + userInfo.Salt)
	db.DataBase.Model(&userInfo).Create(&userInfo)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": "账户创建成功", "error": ""})
}

// LogIn 登录已有账户
func LogIn(c *gin.Context) {
	var req protocol.UserReq
	c.BindJSON(req)
	// 判断是否已经存在
	var userInfo protocol.UserInfo
	var count int
	db.DataBase.Model(&userInfo).Where("name = ?", req.UserName).Count(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": "", "error": "账户不存在，请先注册"})
		return
	}

	// 登录
	db.DataBase.Model(&userInfo).Where("name = ?", req.UserName).Find(&userInfo)
	if common.Sha1En(req.PassWord+userInfo.Salt) != userInfo.PassWord {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": "", "error": "密码错误"})
		return
	}
	// ip组登录，记录在服务器缓存
	ip := c.ClientIP()
	server.ClientIPs[ip] = true
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": "", "error": ""})
}

// StartChat 利用websocket开启即时通讯
func StartChat(c *gin.Context) {
	ip := c.ClientIP()
	if _, ok := server.ClientIPs[ip]; !ok {
		// 未登录
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "data": "", "error": "未登录，请先登录"})
	}
	server.ServeWs(hub, c.Writer, c.Request)
}
