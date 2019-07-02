package main

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"log"
	//	"os"
	"time"
)

func main() {
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	gosnmp.Default.Target = "127.0.0.1"
	gosnmp.Default.Port = 162
	gosnmp.Default.Timeout = time.Duration(10) * time.Second
	gosnmp.Default.Retries = 10
	gosnmp.Default.Version = gosnmp.Version2c
	gosnmp.Default.Community = "public"
	// gosnmp.Default.Logger = log.New(os.Stdout, "", 0)

	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()

	for i := 0; i < 1000; i++ {
		time.Sleep(100 * time.Microsecond)
		pdu := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0",
			Type:  gosnmp.ObjectIdentifier,
			Value: ".1.3.6.1.6.3.1.1.5.1",
		}

		value1 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.1",
			Type:  gosnmp.OctetString,
			Value: fmt.Sprintf("message.%d", i),
		}

		trap := gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{pdu, value1},
		}

		_, err = gosnmp.Default.SendTrap(trap)
		if err != nil {
			log.Fatalf("SendTrap() err: %v", err)
		}
	}

}
