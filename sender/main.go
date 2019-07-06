package main

import (
	"github.com/soniah/gosnmp"
	"log"
	//	"os"
	"flag"
	"time"
)

// A test program that simulate the SNMP trap notification from the devices.
func main() {

	sn := flag.String("sn", "sn-huawei-0", "The serial number in the trap message.")
	description := flag.String("description", "CPU temperature too high", "The description of the trap message.")
	key := flag.String("key", "key", "The key of the trap message.")
	versusKey := flag.String("versuskey", "versuskey", "The versus key of the trap message.")
	trapType := flag.String("type", "Alert", "The type of the trap message. Can be Alert or Event")
	serverity := flag.String("serverity", "Critical", "The serverity of the trap message. Can be OK, Warning or Critical")
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

	for i := 0; i < 1; i++ {
		time.Sleep(100 * time.Microsecond)
		pdu := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0",
			Type:  gosnmp.ObjectIdentifier,
			Value: ".1.3.6.1.6.3.1.1.5.1",
		}

		value1 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.1",
			Type:  gosnmp.OctetString,
			Value: *sn,
		}

		value2 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.2",
			Type:  gosnmp.OctetString,
			Value: *key,
		}

		value3 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.3",
			Type:  gosnmp.OctetString,
			Value: *versusKey,
		}

		value4 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.4",
			Type:  gosnmp.OctetString,
			Value: *trapType,
		}

		value5 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.5",
			Type:  gosnmp.OctetString,
			Value: *serverity,
		}

		value6 := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0.6",
			Type:  gosnmp.OctetString,
			Value: *description,
		}

		trap := gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{pdu, value1, value2, value3, value4, value5, value6},
		}

		_, err = gosnmp.Default.SendTrap(trap)
		if err != nil {
			log.Fatalf("SendTrap() err: %v", err)
		}
	}

}
