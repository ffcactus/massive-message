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

var (
	// global variables in this package.
	connection *amqp.Connection
	channel    *amqp.Channel
)

// Init creates the exchange and queue for receiving notification message.
func Init() error {
	return initMessageQueue()
}

// Release the resources used in this service.
func Release() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}

func initMessageQueue() error {
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, dail failed.")
		return err
	}

	channel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, create channel failed.")
		connection.Close()
		return err
	}

	if err := channel.ExchangeDeclare(
		sdk.NotificationExchangeName, // name
		"topic", // type
		true,    // duarable
		false,   // auto-deleted
		false,   // internal,
		false,   // no-wait,
		nil,     // args
	); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, create exchange failed.")
	}
	log.WithFields(log.Fields{"exchange": sdk.NotificationExchangeName, "type": "topic"}).Info("[Notification-Service] MQ service initialized.")
	return nil
}

// Start the notification process.
// This function won't return.
func Start() {
	subscribe([]string{"*.*"}, handler)
}

// Subscribe the topics.
// The handler will process each of the delivery.
// You can call this method mutiple times to use other handlers to process other topics.
func subscribe(topices []string, handler func(d *amqp.Delivery)) error {
	if channel == nil {
		log.Warn("[Notification-Service] Subscribe event failed, no channel, forgot to init event service?")
		return fmt.Errorf("no channel")
	}
	q, err := channel.QueueDeclare(
		fmt.Sprintf("%s to notification", sdk.NotificationExchangeName),
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topices, "error": err}).Warn("[Notification-Service] Subscribe event failed, declare queue failed.")
		return err
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Notification-Service] Event queue created.")
	for _, topice := range topices {
		if err := channel.QueueBind(q.Name, topice, sdk.NotificationExchangeName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topice, "error": err}).Warn("[Notification-Service] Subscribe event failed, bind queue failed.")
			return err
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
	return nil
}

func handler(delivery *amqp.Delivery) {
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	notification := sdk.Notification{}
	if err := decoder.Decode(&notification); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification-Service] Decode payload failed.")
		return
	}
	repository.SaveNotification(&notification)
}
