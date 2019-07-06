package notification

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	notificationSDK "massive-message/notification/sdk"
	receiverSDK "massive-message/receiver/sdk"
)

const (
	trapExchange         = "SnmpTrapExchanger"
	notificationExchange = "NotificationExchange"
)

var (
	// global variables in this package.
	connection *amqp.Connection
	channel    *amqp.Channel

	// notificationConverterMapping is a mapping from notification to it's converter.
	// The Address in the notification will be taken as the key here, since only this information is contained in the notification for sure.
	notificationConverterMapping = make(map[string]Converter)
	// snURLMapping is a mapping from server's serial number to server's URL.
	snURLMapping = make(map[string]string)

	huaweiNotificationConverter = &HuaweiNotificationConverter{}
	dellNotificationConverter   = &DellNotificationConverter{}
	hpeNotificationConverter    = &HpeNotificationConverter{}
	ibmNotificationConverter    = &IBMNotificationConverter{}
	lenovoNotificationConverter = &LenovoNotificationConverter{}
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
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = huaweiNotificationConverter
		snURLMapping[sn] = "/redfish/v1/servers/" + uuid.New().String()
	}
	for i := 10000; i < 20000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = dellNotificationConverter
		snURLMapping[sn] = "/redfish/v1/servers/" + uuid.New().String()
	}
	for i := 20000; i < 30000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = hpeNotificationConverter
		snURLMapping[sn] = "/redfish/v1/servers/" + uuid.New().String()
	}
	for i := 30000; i < 40000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = ibmNotificationConverter
		snURLMapping[sn] = "/redfish/v1/servers/" + uuid.New().String()
	}
	for i := 40000; i < 50000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		notificationConverterMapping[sn] = lenovoNotificationConverter
		snURLMapping[sn] = "/redfish/v1/servers/" + uuid.New().String()
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
		receiverSDK.TrapExchangeName, // name
		"topic", // type
		true,    // duarable
		false,   // auto-deleted
		false,   // internal,
		false,   // no-wait,
		nil,     // args
	); err != nil {
		log.WithFields(log.Fields{"exchange": receiverSDK.TrapExchangeName, "err": err}).Error("[Server] Init MQ service failed, create exchange failed.")
	}

	if err := channel.ExchangeDeclare(
		notificationSDK.NotificationExchangeName, // name
		"topic", // type
		true,    // duarable
		false,   // auto-deleted
		false,   // internal,
		false,   // no-wait,
		nil,     // args
	); err != nil {
		log.WithFields(log.Fields{"exchange": notificationSDK.NotificationExchangeName, "err": err}).Error("[Server] Init MQ service failed, create exchange failed.")
	}
	log.Info("[Server] MQ service initialized.")
	return nil
}
