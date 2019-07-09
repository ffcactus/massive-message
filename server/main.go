package main

import (
	log "github.com/sirupsen/logrus"
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
	notification.Start()
}
