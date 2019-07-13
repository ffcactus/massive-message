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

// StartReceiver starts the notification receiver.
// This function should be called as co-routine.
func StartReceiver() {

	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, dail failed.")
		return
	}
	defer func() {
		if err := connection.Close(); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Close connection failed.")
		}
	}()
	channel, err := connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Init MQ service failed, create channel failed.")
		return
	}
	defer func() {
		if err := channel.Close(); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Close connection failed.")
		}
	}()
	// args: exchange, type, durable, auto-deleted, internal, no-wait, args
	if err := channel.ExchangeDeclare(sdk.NotificationExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Notification-Service] Start notification process failed, declare exchange failed.")
		return
	}
	topics := []string{"*.*"}
	// args: name, durable, autoDelete, exclusive, noWait, args
	q, err := channel.QueueDeclare(
		fmt.Sprintf("%s to notification", sdk.NotificationExchangeName),
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topics, "error": err}).Warn("[Notification-Service] Subscribe event failed, declare queue failed.")
		return
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Notification-Service] Event queue created.")
	for _, topic := range topics {
		if err := channel.QueueBind(q.Name, topic, sdk.NotificationExchangeName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topic, "error": err}).Warn("[Notification-Service] Subscribe event failed, bind queue failed.")
			return
		}
		log.WithFields(log.Fields{"topic": topic}).Info("[Notification-Service] Event queue bind.")
	}
	delivery, err := channel.Consume(q.Name, "my-consume", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topics": topics, "error": err}).Warn("[Notification-Service] Subscribe event failed, consume failed.")
	}

	for each := range delivery {
		handler(&each)
		if err := each.Ack(false); err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("[Notification-Service] [Notification-Service] Handle notification failed, ack failed.")
		}
	}
	log.WithFields(log.Fields{"topics": topics}).Warn("[Notification-Service] Subscribe event exit.")
}

func handler(delivery *amqp.Delivery) {
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	notification := sdk.Notification{}
	if err := decoder.Decode(&notification); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification-Service] Decode payload failed.")
		return
	}
	if err := repository.SaveNotification(&notification); err != nil {
		log.WithFields(log.Fields{"notification": notification}).Warn("[Notification-Service] Save notification failed.")
	}
	log.WithFields(log.Fields{"url": notification.URL, "key": notification.Key}).Info("[Notification-Service] Notification saved.")
}
