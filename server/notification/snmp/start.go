package snmp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	notificationSDK "massive-message/notification/sdk"
	receiverSDK "massive-message/receiver/sdk"
)

var (
	snmpReceiveChannel *amqp.Channel
	notificationSendChannel *amqp.Channel
)

// Listen on the notification.
// This function should be used ad co-routine.
// It first prepare the channel and exchange for both SNMP receiving and notification sending.
// Then it declare queue for SNMP receiving.
// For each of the SNMP message, finding the converter to convert to standard notification message, sending it.
func Listen() {
	// Open connection.
	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Service-SNMP] Listen failed, create connection failed.")
		return
	}
	defer func() {
		if err := connection.Close(); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("[Service-SNMP] Listen failed, Close amqp connection failed.")
		}
	}()
	// Create channel for snmp receive and notification send.
	snmpReceiveChannel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Service-SNMP] Listen failed, create channel failed.")
		return
	}
	defer func() {
		if err := snmpReceiveChannel.Close(); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("[Service-SNMP] Close SNMP receive channel failed.")
		}
	}()
	notificationSendChannel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Service-SNMP] Listen failed, create channel failed.")
		return
	}
	defer func() {
		if err := notificationSendChannel.Close(); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("[Service-SNMP] Close notification send channel failed.")
		}
	}()
	// Create exchange for snmp receive and notification send.
	if err := snmpReceiveChannel.ExchangeDeclare(receiverSDK.TrapExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"exchange": receiverSDK.TrapExchangeName, "err": err}).Error("[Service-SNMP] Listen failed, create exchange failed.")
		return
	}
	if err := snmpReceiveChannel.ExchangeDeclare(notificationSDK.NotificationExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"exchange": notificationSDK.NotificationExchangeName, "err": err}).Error("[Service-SNMP] Listen failed, create exchange failed.")
		return
	}
	// Queue.
	topics := []string{"*.*"}
	q, err := snmpReceiveChannel.QueueDeclare(fmt.Sprintf("%s to server", receiverSDK.TrapExchangeName), true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topics": topics, "error": err}).Warn("[Service-SNMP] Listen failed, declare queue failed.")
		return
	}
	// Topic.
	for _, topic := range topics {
		if err := snmpReceiveChannel.QueueBind(q.Name, topic, receiverSDK.TrapExchangeName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topic, "error": err}).Warn("[Service-SNMP] Listen failed, bind queue to topic failed.")
			return
		}
	}
	delivery, err := snmpReceiveChannel.Consume(q.Name, "server_management_service", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topics": topics, "error": err}).Warn("[Service-SNMP] Listen failed, consume failed.")
	}
	log.Info("[Service-SNMP] Listen on delivery.")
	for each := range delivery {
		handler(&each)
		if err := each.Ack(false); err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("[Service-SNMP] Ack failed.")
		}
	}
	log.WithFields(log.Fields{"topics": topics}).Warn("[Service-SNMP] Stop listen.")
}


func handler(delivery *amqp.Delivery) {
	var (
		snOID = ".1.3.6.1.6.3.1.1.4.1.0.1"
		sn    string
	)
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	payload := receiverSDK.WrapedSnmpPacket{}
	if err := decoder.Decode(&payload); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Service-SNMP] Handler notification failed, decode SNMP payload failed.")
		return
	}

	for _, v := range payload.Variables {
		if v.Name == snOID {
			sn = v.String()
		}
	}
	converter := notificationConverterMapping[sn]
	if converter == nil {
		log.WithFields(log.Fields{"sn": sn}).Warn("[Service-SNMP] Handler notification failed, unable to find the converter.")
		return
	}
	standardNotifications, err := converter.Convert(&payload)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Service-SNMP] Handler notification failed, convert SNMP payload failed.")
	}
	network := bytes.Buffer{}
	encoder := gob.NewEncoder(&network)
	for _, v := range standardNotifications {
		if err := encoder.Encode(v); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Service-SNMP] Handler notification failed, Encode to notification message failed.")
			return
		}
		if err := notificationSendChannel.Publish(notificationSDK.NotificationExchangeName, "Notification.New", false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         network.Bytes(),
		}); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Service-SNMP] Handler notification failed, publish notification message failed.")
		}
	}
}
