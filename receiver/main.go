package main

import (
	"bytes"
	"encoding/gob"
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

func initListener() (*gosnmp.TrapListener, error) {
	listener := gosnmp.NewTrapListener()
	listener.OnNewTrap = handler
	listener.Params = gosnmp.Default
	return listener, nil
}

func initMessageQueue() error {
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init MQ service failed, dail failed.")
		return err
	}

	channel, err = connection.Channel()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init MQ service failed, create channel failed.")
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
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init MQ service failed, create exchange failed.")
	}
	log.WithFields(log.Fields{"exchange": exchangerName, "type": "topic"}).Info("[Receiver] MQ service initialized.")
	return nil
}

func handler(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	var (
		network bytes.Buffer
	)

	count++

	log.Info("[Receiver] SNMP trap received = ", count)

	payload := WrapedSnmpPacket{}
	payload.Address = addr
	payload.ReceivedAt = time.Now()
	for _, v := range packet.Variables {
		payload.Variables = append(payload.Variables, SnmpVariable{
			Name:  v.Name,
			Type:  v.Type,
			Value: v.Value,
		})
	}
	encoder := gob.NewEncoder(&network)
	if err := encoder.Encode(payload); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Encoding the SNMP trap message failed.")
		return
	}
	if err := channel.Publish(exchangerName, "Snmp.New", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        network.Bytes(),
	}); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Publish event failed.")
	}
}

func main() {
	listener, err := initListener()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init SNMP trap listener failed.")
		os.Exit(-1)
	}
	err = initMessageQueue()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init MQ service failed.")
		listener.Close()
		os.Exit(-1)
	}
	if err := listener.Listen("0.0.0.0:162"); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Start listener failed.")
	}
	listener.Close()
	channel.Close()
	connection.Close()
}
