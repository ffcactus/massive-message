package healthchange

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	notificationSDK "massive-message/notification/sdk"
	"massive-message/server/repository"
	"strings"
)

var (
	// global variables in this package.
	connection *amqp.Connection
	channel    *amqp.Channel
)

// Start begins the health change notification process.
func Start() {
	var (
		err          error
		exchangeName = notificationSDK.HealthChangeExchangeName
	)
	// Connect to MQ service.
	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server-HealthChange] Start message process failed, dail to service failed.")
		return
	}
	defer connection.Close()

	// Create a channel.
	channel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server-HealthChange] Start message process failed, create channel failed.")
		connection.Close()
		return
	}
	defer channel.Close()

	// Create exchange in case we started earlier than sender.
	// args: exchange, type, durable, auto-deleted, internal, no-wait, args.
	if err := channel.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"exchange": exchangeName, "err": err}).Error("[Server-HealthChange] Start message process failed, declare exchange failed.")
	}

	// Declare queue.
	q, err := channel.QueueDeclare(
		fmt.Sprintf("%s to server", exchangeName),
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-HealthChange] Start message process failed, declare queue failed.")
		return
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Server-HealthChange] Start message process, event queue created.")

	// Bind the queue to the exchange.
	if err := channel.QueueBind(q.Name, "*.*", exchangeName, false, nil); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-HealthChange] Start message process failed, bind queue failed.")
		return
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Server-HealthChange] Start message process, event queue bind.")

	// Get the delivery channel.
	// args: queue, consumerString, autoAck, exclusive, noLocal, noWait, args
	delivery, err := channel.Consume(q.Name, "server-consumer", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-HealthChange] Start message process failed, consume failed.")
	}

	// Keep getting the payload from the channel.
	for each := range delivery {
		handler(&each)
		each.Ack(false)
	}
}

// Handle each of the payload.
func handler(delivery *amqp.Delivery) {
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	notification := notificationSDK.HealthChangeNotification{}
	if err := decoder.Decode(&notification); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-HealthChange] Decode payload failed.")
		return
	}
	ids := strings.Split(notification.URL, "/")
	if len(ids) == 5 && ids[3] == "servers" {
		repository.UpdateServerHealth(ids[4], notification.Warnings, notification.Criticals)
		log.WithFields(log.Fields{"id": ids[4], "warnings": notification.Warnings, "criticals": notification.Criticals}).Info("[Server-HealthChange] Received notification")
	}
}
