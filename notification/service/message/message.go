package message

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"massive-message/notification/repository"
	"massive-message/notification/sdk"
)

// var (
// 	// global variables in this package.
// 	connection *amqp.Connection
// 	channel    *amqp.Channel
// )

// InitConnection initialize the connection and channel that is used for message process.
// ReleaseConnection should be used later.
// func InitConnection() error {
// 	var err error
// 	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
// 	if err != nil {
// 		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, dail failed.")
// 		return err
// 	}

// 	channel, err = connection.Channel()
// 	if err != nil {
// 		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, create channel failed.")
// 		connection.Close()
// 		return err
// 	}

// 	log.WithFields(log.Fields{"exchange": sdk.NotificationExchangeName, "type": "topic"}).Info("[Notification-Service] MQ service initialized.")
// 	return nil
// }

// // CloseConnection closes the connection and channel.
// func CloseConnection() {
// 	if channel != nil {
// 		channel.Close()
// 	}
// 	if connection != nil {
// 		connection.Close()
// 	}
// }

// StartReceiver starts the notification receiver.
// This function should be called as co-routine.
func StartReceiver() {

	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, dail failed.")
		return
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, create channel failed.")
		connection.Close()
		return
	}
	defer channel.Close()
	// args: exchange, type, durable, auto-deleted, internal, no-wait, args
	if err := channel.ExchangeDeclare(sdk.NotificationExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Start notification process failed, declare exchange failed.")
		return
	}
	topices := []string{"*.*"}
	q, err := channel.QueueDeclare(
		fmt.Sprintf("%s to notification", sdk.NotificationExchangeName),
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topices, "error": err}).Warn("[Notification-Service] Subscribe event failed, declare queue failed.")
		return
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Notification-Service] Event queue created.")
	for _, topice := range topices {
		if err := channel.QueueBind(q.Name, topice, sdk.NotificationExchangeName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topice, "error": err}).Warn("[Notification-Service] Subscribe event failed, bind queue failed.")
			return
		}
		log.WithFields(log.Fields{"topic": topice}).Info("[Notification-Service] Event queue bind.")
	}
	delivery, err := channel.Consume(q.Name, "my-consume", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topices": topices, "error": err}).Warn("[Notification-Service] Subscribe event failed, consume failed.")
	}

	for each := range delivery {
		handler(&each)
		each.Ack(false)
	}
	log.WithFields(log.Fields{"topices": topices}).Warn("[Notification-Service] Subscribe event exit.")
}

func handler(delivery *amqp.Delivery) {
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	notification := sdk.Notification{}
	if err := decoder.Decode(&notification); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification-Service] Decode payload failed.")
		return
	}
	repository.SaveNotification(&notification)
	log.WithFields(log.Fields{"url": notification.URL, "key": notification.Key}).Info("[Notification-Service] Notification saved.")
}
