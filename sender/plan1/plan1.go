// Package plan1 sends 50000 * 100 events.
// Assume there are 5 vendor and each of them have 10000 servers.
// For each of the server sends 100 events.
package plan1

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"time"
)

// Do perform the work in this plan.
func Do(snmp *gosnmp.GoSNMP, serverPerVendor, eventPerServer, interval int) {
	vendor1 := make([]gosnmp.SnmpTrap, serverPerVendor*eventPerServer)
	vendor2 := make([]gosnmp.SnmpTrap, serverPerVendor*eventPerServer)
	vendor3 := make([]gosnmp.SnmpTrap, serverPerVendor*eventPerServer)
	vendor4 := make([]gosnmp.SnmpTrap, serverPerVendor*eventPerServer)
	vendor5 := make([]gosnmp.SnmpTrap, serverPerVendor*eventPerServer)

	// prepare
	i := 0
	for sn := 0; sn < serverPerVendor; sn++ {
		for key := 0; key < eventPerServer; key++ {
			vendor1[i] = *generateEvent(fmt.Sprintf("sn-huawei-%d", key), fmt.Sprintf("key%d", key))
			vendor2[i] = *generateEvent(fmt.Sprintf("sn-hpe-%d", key), fmt.Sprintf("key%d", key))
			vendor3[i] = *generateEvent(fmt.Sprintf("sn-dell-%d", key), fmt.Sprintf("key%d", key))
			vendor4[i] = *generateEvent(fmt.Sprintf("sn-ibm-%d", key), fmt.Sprintf("key%d", key))
			vendor5[i] = *generateEvent(fmt.Sprintf("sn-lenovo-%d", key), fmt.Sprintf("key%d", key))
			i++
		}
	}

	for i := 0; i < serverPerVendor*eventPerServer; i++ {
		time.Sleep(time.Duration(interval) * time.Microsecond)
		snmp.SendTrap(vendor1[i])
		time.Sleep(time.Duration(interval) * time.Microsecond)
		snmp.SendTrap(vendor2[i])
		time.Sleep(time.Duration(interval) * time.Microsecond)
		snmp.SendTrap(vendor3[i])
		time.Sleep(time.Duration(interval) * time.Microsecond)
		snmp.SendTrap(vendor4[i])
		time.Sleep(time.Duration(interval) * time.Microsecond)
		snmp.SendTrap(vendor5[i])
		fmt.Println(i)
	}
}

func generateEvent(sn, key string) *gosnmp.SnmpTrap {
	trapType := "Event"
	serverity := "OK"

	pdu := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  gosnmp.ObjectIdentifier,
		Value: ".1.3.6.1.6.3.1.1.5.1",
	}

	value1 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.1",
		Type:  gosnmp.OctetString,
		Value: sn,
	}

	value2 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.2",
		Type:  gosnmp.OctetString,
		Value: key,
	}

	value3 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.3",
		Type:  gosnmp.OctetString,
		Value: fmt.Sprintf("versus%s", key),
	}

	value4 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.4",
		Type:  gosnmp.OctetString,
		Value: trapType,
	}

	value5 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.5",
		Type:  gosnmp.OctetString,
		Value: serverity,
	}

	value6 := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.6",
		Type:  gosnmp.OctetString,
		Value: fmt.Sprintf("This is event for key %s", key),
	}

	return &gosnmp.SnmpTrap{
		Variables: []gosnmp.SnmpPDU{pdu, value1, value2, value3, value4, value5, value6},
	}
}
