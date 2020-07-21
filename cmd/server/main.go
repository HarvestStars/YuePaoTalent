package main

import (
	"github.com/HarvestStars/YuePaoTalent/conf"
	"github.com/HarvestStars/YuePaoTalent/db"
	"github.com/HarvestStars/YuePaoTalent/util/server"
	"github.com/gin-gonic/gin"
)

var hub *server.Hub

func main() {
	conf.Setup()
	db.Setup(conf.MySQLSetting.User, conf.MySQLSetting.PassWord, conf.MySQLSetting.Host, conf.MySQLSetting.DataBase)
	hub = server.NewHub()
	go hub.Run()
	r := gin.Default()
	r.POST("/signup", SignUp)
	r.POST("/login", LogIn)
	r.GET("/chat", StartChat)
	r.Run(":8080")
}
