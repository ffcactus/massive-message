package main

import (
	"bytes"
	"encoding/gob"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/streadway/amqp"
	"massive-message/receiver/sdk"
	"net"
	"os"
	"time"
)

var (
	count      int64
	connection *amqp.Connection
	channel    *amqp.Channel
)

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

	if err := channel.ExchangeDeclare(sdk.TrapExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Init MQ service failed, create exchange failed.")
	}
	log.WithFields(log.Fields{"exchange": sdk.TrapExchangeName, "type": "topic"}).Info("[Receiver] MQ service initialized.")
	return nil
}

func handler(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	var (
		network bytes.Buffer
	)

	count++

	log.Info("[Receiver] SNMP trap received = ", count)

	payload := sdk.WrapedSnmpPacket{}
	payload.Address = addr
	payload.GeneratedAt = time.Now()
	for _, v := range packet.Variables {
		payload.Variables = append(payload.Variables, sdk.SnmpVariable{
			Name:  v.Name,
			Type:  v.Type,
			Value: v.Value,
		})
	}
	encoder := gob.NewEncoder(&network)
	if err := encoder.Encode(payload); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Encoding SNMP trap message failed.")
		return
	}
	if err := channel.Publish(sdk.TrapExchangeName, "Trap.New", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "application/octet-stream",
		Body:        network.Bytes(),
	}); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Receiver] Publish SNMP trap message failed.")
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
