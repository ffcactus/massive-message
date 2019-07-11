package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"massive-message/notification/repository"
	messageService "massive-message/notification/service/message"
	"os"
)

func main() {

	if err := repository.Init(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification] Initialize repository failed.")
		os.Exit(-1)
	}
	log.Info("[Notification] Initialize repository done.")
	if err := repository.DropTablesIfExist(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification] Drop tables failed.")
		os.Exit(-1)
	}
	log.Info("[Notification] Drop tables done.")
	if err := repository.CreateTables(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification] Create tables failed.")
		os.Exit(-1)
	}
	log.Info("[Notification] Create tables done.")

	go messageService.StartHealthTracker()
	messageService.StartReceiver()
}
