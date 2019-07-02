package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/streadway/amqp"
	"net"
	"os"
	"time"
)

const (
	exchangerName = "SnmpTrapExchanger"
)

var (
	count      int64
	connection *amqp.Connection
	channel    *amqp.Channel
)

func Asn1BERToString(i gosnmp.Asn1BER) string {
	switch i {
	case 0x00:
		return "UnknownType"
	case 0x01:
		return "Boolean"
	case 0x02:
		return "Integer"
	case 0x03:
		return "BitString"
	case 0x04:
		return "OctetString"
	case 0x05:
		return "Null"
	case 0x06:
		return "ObjectIdentifier"
	case 0x07:
		return "ObjectDescription"
	case 0x40:
		return "IPAddress"
	case 0x41:
		return "Counter32"
	case 0x42:
		return "Gauge32"
	case 0x43:
		return "TimeTicks"
	case 0x44:
		return "Opaque"
	case 0x45:
		return "NsapAddress"
	case 0x46:
		return "Counter64"
	case 0x47:
		return "Uinteger32"
	case 0x78:
		return "OpaqueFloat"
	case 0x79:
		return "OpaqueDouble"
	case 0x80:
		return "NoSuchObject"
	case 0x81:
		return "NoSuchInstance"
	case 0x82:
		return "EndOfMibView"
	default:
		return "UnknownType"
	}
}

type SnmpVariable struct {
	// Name is an oid in string format eg ".1.3.6.1.4.9.27"
	Name string
	// The type of the value eg Integer
	Type gosnmp.Asn1BER
	// The value to be set by the SNMP set, or the value when
	// sending a trap
	Value interface{}
}

// WrapedSnmpPacket includes both the raw SNMP packet but also some other useful information for processing it later.
// Since we are using encoding/gob, it's OK to use point here.
type WrapedSnmpPacket struct {
	Address    *net.UDPAddr
	ReceivedAt time.Time
	Variables  []SnmpVariable
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

func releaseMessageQueue() {
	channel.Close()
	connection.Close()
}

// Subscribe the topics.
// The handler will process each of the delivery.
// You can call this method mutiple times to use other handlers to process other topics.
func subscribe(topices []string, handler func(d *amqp.Delivery)) error {
	if channel == nil {
		log.Warn("[Server] Subscribe event failed, no channel, forgot to init event service?")
		return fmt.Errorf("no channel")
	}
	q, err := channel.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topic": topices, "error": err}).Warn("[Server] Subscribe event failed, declare queue failed.")
		return err
	}
	log.WithFields(log.Fields{"Name": q.Name}).Info("[Server] Event queue created.")
	for _, topice := range topices {
		if err := channel.QueueBind(q.Name, topice, exchangerName, false, nil); err != nil {
			log.WithFields(log.Fields{"topic": topice, "error": err}).Warn("[Server] Subscribe event failed, bind queue failed.")
			return err
		}
		log.WithFields(log.Fields{"topic": topice}).Info("[Server] Event queue bind.")
	}
	delivery, err := channel.Consume(q.Name, "my-consume", false, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{"topices": topices, "error": err}).Warn("[Server] Subscribe event failed, consume failed.")
	}

	for each := range delivery {
		handler(&each)
		each.Ack(false)
	}
	log.WithFields(log.Fields{"topices": topices}).Warn("[Server] Subscribe event exit.")
	return nil
}

func handler(delivery *amqp.Delivery) {
	count++
	log.Info("[Receiver] SNMP trap received =", count)
	decoder := gob.NewDecoder(bytes.NewBuffer(delivery.Body))
	payload := WrapedSnmpPacket{}
	if err := decoder.Decode(&payload); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Server] Decode payload failed.")
		return
	}
	printPayload(&payload)
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

func main() {
	err := initMessageQueue()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Server] Init MQ service failed.")
		os.Exit(-1)
	}
	defer releaseMessageQueue()
	subscribe([]string{"*.*"}, handler)
}
