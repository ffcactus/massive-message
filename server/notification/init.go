package notification

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	exchangerName = "SnmpTrapExchanger"
)

var (
	// global variables in this package.
	connection *amqp.Connection
	channel    *amqp.Channel

	// notificationConverterMapping is a map from notification to it's converter.
	// The Address in the notification will be taken as the key here, since only this information is contained in the notification for sure.
	notificationConverterMapping = make(map[string]Converter)
	huaweiNotificationConverter  = &HuaweiNotificationConverter{}
	dellNotificationConverter    = &DellNotificationConverter{}
	hpeNotificationConverter     = &HpeNotificationConverter{}
)

// Init includes all kinds of work that should be done before notification processing, it includes and should be performed in the order below:
// 1. Prepare the message queue so that when the nofication comes we can put it into immediately.
// 2. Prepare the mapping from the notification to the converter.
func Init() error {
	if err := initMessageQueue(); err != nil {
		return err
	}
	generateConverterMapping()
	return nil
}

func generateConverterMapping() {
	for i := 0; i < 100000; i++ {
		notificationConverterMapping[fmt.Sprintf("serial-number-%d", i)] = huaweiNotificationConverter
	}
	for i := 100000; i < 200000; i++ {
		notificationConverterMapping[fmt.Sprintf("serial-number-%d", i)] = dellNotificationConverter
	}
	for i := 200000; i < 300000; i++ {
		notificationConverterMapping[fmt.Sprintf("serial-number-%d", i)] = hpeNotificationConverter
	}
}

func initMessageQueue() error {
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed, dail failed.")
		return err
	}

	channel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed, create channel failed.")
		connection.Close()
		return err
	}

	if err := channel.ExchangeDeclare(
		exchangerName, // name
		"topic",       // type
		true,          // duarable
		false,         // auto-deleted
		false,         // internal,
		false,         // no-wait,
		nil,           // args
	); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed, create exchange failed.")
	}
	log.WithFields(log.Fields{"exchange": exchangerName, "type": "topic"}).Info("[Server] MQ service initialized.")
	return nil
}
