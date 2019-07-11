package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	log "github.com/sirupsen/logrus"
	"massive-message/server/controller"
	"massive-message/server/notification"
	"massive-message/server/notification/healthchange"
	"massive-message/server/repository"
	"os"
)

func main() {
	repository.Init()
	err := notification.Init()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed.")
		os.Exit(-1)
	}
	defer notification.Release()
	go healthchange.Start()
	go notification.Start()

	serverController := controller.Server{}
	serverNS := beego.NewNamespace(
		"/api/v1/servers",
		beego.NSRouter("/", &serverController, "get:GetServers"),
		beego.NSRouter("/:id", &serverController, "get:GetServerByID"),
	)
	beego.AddNamespace(serverNS)

	beego.BConfig.Listen.HTTPPort = 80
	beego.BConfig.CopyRequestBody = true
	beego.SetLevel(beego.LevelNotice)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	beego.Run()
}
