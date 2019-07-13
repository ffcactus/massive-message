package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"massive-message/server/controller"
	"massive-message/server/notification/healthchange"
	"massive-message/server/notification/snmp"
	"massive-message/server/repository"
	"os"
)

func main() {
	if err := repository.Init(); err != nil {
		os.Exit(-1)
	}
	snmp.Init()

	go healthchange.Start()
	go snmp.Listen()

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
