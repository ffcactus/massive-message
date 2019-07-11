package main

import (
	"flag"
	"fmt"
	"github.com/soniah/gosnmp"
	"log"
	"massive-message/sender/plan1"
	"massive-message/sender/plan2"
	"os"
	"strconv"
	"strings"
	"time"
)

// A test program that simulate the SNMP trap notification from the devices.
func main() {

	planFlag := flag.String("plan", "", "The plan to use.")
	sn := flag.String("sn", "sn-huawei-0", "The serial number in the trap message.")
	description := flag.String("description", "CPU temperature too high", "The description of the trap message.")
	key := flag.String("key", "key", "The key of the trap message.")
	versusKey := flag.String("versuskey", "versuskey", "The versus key of the trap message.")
	trapType := flag.String("type", "Alert", "The type of the trap message. Can be Alert or Event")
	serverity := flag.String("serverity", "Critical", "The serverity of the trap message. Can be OK, Warning or Critical")
	flag.Parse()
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

	// if plan is specified, use the plan.
	if *planFlag != "" {
		plan := strings.Split(*planFlag, ",")
		switch plan[0] {
		case "plan1":
			if len(plan) != 4 {
				fmt.Println("Use plan1 with format plan1,serverPerVendor,eventPerServer,interval")
				os.Exit(-1)
			}
			serverPerVendor, _ := strconv.ParseInt(plan[1], 10, 32)
			eventPerServer, _ := strconv.ParseInt(plan[2], 10, 32)
			interval, _ := strconv.ParseInt(plan[3], 10, 32)
			fmt.Printf("Using plan1 with serverPerVendor %d, eventPerServe %d, interval %d\n", int(serverPerVendor), int(eventPerServer), int(interval))
			plan1.Do(gosnmp.Default, int(serverPerVendor), int(eventPerServer), int(interval))
			return
		case "plan2":
			if len(plan) != 6 {
				fmt.Println("Use plan2 with format plan2,vendor,serverCount,severity,alertCount,interval")
				os.Exit(-1)
			}
			vendor := plan[1]
			serverCount, _ := strconv.ParseInt(plan[2], 10, 32)
			severity := plan[3]
			alertCount, _ := strconv.ParseInt(plan[4], 10, 32)
			interval, _ := strconv.ParseInt(plan[5], 10, 32)
			plan2.Do(gosnmp.Default, vendor, int(serverCount), severity, int(alertCount), int(interval))
			return
		}
	}

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
