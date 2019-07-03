package notification

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/streadway/amqp"
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
		log.Warn("[Notification] Subscribe event failed, no channel, forgot to init event service?")
		return fmt.Errorf("no channel")
	}
	q, err := channel.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topices, "error": err}).Warn("[Notification] Subscribe event failed, declare queue failed.")
		return err
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Notification] Event queue created.")
	for _, topice := range topices {
		if err := channel.QueueBind(q.Name, topice, exchangerName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topice, "error": err}).Warn("[Notification] Subscribe event failed, bind queue failed.")
			return err
		}
		log.WithFields(log.Fields{"topic": topice}).Info("[Notification] Event queue bind.")
	}
	delivery, err := channel.Consume(q.Name, "my-consume", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topices": topices, "error": err}).Warn("[Notification] Subscribe event failed, consume failed.")
	}

	for each := range delivery {
		handler(&each)
		each.Ack(false)
	}
	log.WithFields(log.Fields{"topices": topices}).Warn("[Notification] Subscribe event exit.")
	return nil
}

func handler(delivery *amqp.Delivery) {
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	payload := WrapedSnmpPacket{}
	if err := decoder.Decode(&payload); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification] Decode payload failed.")
		return
	}
	converter := notificationConverterMapping[payload.Address.IP.String()]
	standardNotifications, err := converter.Convert(&payload)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification] Handler notification failed, discard notification.")
	}
	for _, v := range standardNotifications {
		v = v
	}
}

func printPayload(payload *WrapedSnmpPacket) {
	// log.Printf("From: %v", payload.Address)
	// log.Printf("ReceivedAt: %v", payload.ReceivedAt)
	for _, v := range payload.Variables {
		// log.Printf("Name: %v", v.Name)
		switch v.Type {
		case gosnmp.OctetString:
			log.Printf("Value: %v", string(v.Value.([]byte)))
		case gosnmp.Integer:
			log.Printf("Value: %v", v.Value.(int))
		case gosnmp.TimeTicks:
			log.Printf("Value: %v", v.Value.(uint))
		default:
			log.Printf("Value: unknown type %v", Asn1BERToString(v.Type))
		}
	}
}
