package main

import (
	log "github.com/sirupsen/logrus"
	"massive-message/server/notification"
	"os"
)

func main() {
	err := notification.Init()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed.")
		os.Exit(-1)
	}
	defer notification.Release()

	notification.Start()
}
