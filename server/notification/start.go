package notification

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	notificationSDK "massive-message/notification/sdk"
	receiverSDK "massive-message/receiver/sdk"
)

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
		log.Warn("[Server-Notification] Subscribe event failed, no channel, forgot to init event service?")
		return fmt.Errorf("no channel")
	}

	q, err := channel.QueueDeclare(
		fmt.Sprintf("%s to server", receiverSDK.TrapExchangeName),
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topices, "error": err}).Warn("[Server-Notification] Subscribe event failed, declare queue failed.")
		return err
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Server-Notification] Event queue created.")
	for _, topice := range topices {
		if err := channel.QueueBind(q.Name, topice, receiverSDK.TrapExchangeName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topice, "error": err}).Warn("[Server-Notification] Subscribe event failed, bind queue failed.")
			return err
		}
		log.WithFields(log.Fields{"topic": topice}).Info("[Server-Notification] Event queue bind.")
	}
	delivery, err := channel.Consume(q.Name, "my-consume", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topices": topices, "error": err}).Warn("[Server-Notification] Subscribe event failed, consume failed.")
	}

	for each := range delivery {
		handler(&each)
		each.Ack(false)
	}
	log.WithFields(log.Fields{"topices": topices}).Warn("[Server-Notification] Subscribe event exit.")
	return nil
}

func handler(delivery *amqp.Delivery) {
	var (
		snOID = ".1.3.6.1.6.3.1.1.4.1.0.1"
		sn    string
	)
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	payload := receiverSDK.WrapedSnmpPacket{}
	if err := decoder.Decode(&payload); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-Notification] Decode payload failed.")
		return
	}

	for _, v := range payload.Variables {
		if v.Name == snOID {
			sn = v.String()
		}
	}
	converter := notificationConverterMapping[sn]
	if converter == nil {
		log.WithFields(log.Fields{"sn": sn}).Warn("[Server-Notification] Unable to find the converter.")
		return
	}
	standardNotifications, err := converter.Convert(&payload)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server-Notification] Handler notification failed, discard notification.")
	}
	network := bytes.Buffer{}
	encoder := gob.NewEncoder(&network)
	for _, v := range standardNotifications {
		if err := encoder.Encode(v); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Server-Notification] Encoding notification message failed.")
			return
		}
		if err := channel.Publish(notificationSDK.NotificationExchangeName, "Notification.New", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        network.Bytes(),
		}); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("[Server-Notification] Publish notification message failed.")
		}
	}
}

// func printPayload(payload *WrapedSnmpPacket) {
// 	// log.Printf("From: %v", payload.Address)
// 	// log.Printf("ReceivedAt: %v", payload.ReceivedAt)
// 	for _, v := range payload.Variables {
// 		// log.Printf("Name: %v", v.Name)
// 		switch v.Type {
// 		case gosnmp.OctetString:
// 			log.Printf("Value: %v", string(v.Value.([]byte)))
// 		case gosnmp.Integer:
// 			log.Printf("Value: %v", v.Value.(int))
// 		case gosnmp.TimeTicks:
// 			log.Printf("Value: %v", v.Value.(uint))
// 		default:
// 			log.Printf("Value: unknown type %v", Asn1BERToString(v.Type))
// 		}
// 	}
// }
