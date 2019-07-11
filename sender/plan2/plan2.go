// Package plan2 sents alerts to servers.
// User can specify the vendor, the server count, the severity and alert count.
// After all the active alerts have been sent , it will send the deactive alerts.
package plan2

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"time"
)

// Do perform the work in this plan.
func Do(snmp *gosnmp.GoSNMP, vendor string, serverCount int, severity string, alertCount int, interval int) {
	for i := 0; i < alertCount; i++ {
		for j := 0; j < serverCount; j++ {
			notification := generateEvent(
				fmt.Sprintf("sn-%s-%d", vendor, j),
				fmt.Sprintf("key-%d", i),
				fmt.Sprintf("versus-key-%d", i),
				"Alert",
				severity,
				fmt.Sprintf("CPU %d state %s", i, severity),
			)
			time.Sleep(2000 * time.Microsecond)
			snmp.SendTrap(*notification)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	time.Sleep(time.Duration(interval) * time.Second)
	for i := 0; i < alertCount; i++ {
		for j := 0; j < serverCount; j++ {
			notification := generateEvent(
				fmt.Sprintf("sn-%s-%d", vendor, j),
				fmt.Sprintf("versus-key-%d", i),
				fmt.Sprintf("key-%d", i),
				"Alert",
				"OK",
				fmt.Sprintf("CPU %d state %s", i, "OK"),
			)
			time.Sleep(2000 * time.Microsecond)
			snmp.SendTrap(*notification)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func generateEvent(sn, key, versusKey, trapType, severity, description string) *gosnmp.SnmpTrap {

	pdu := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  gosnmp.ObjectIdentifier,
		Value: ".1.3.6.1.6.3.1.1.5.1",
	}

	snPDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.1",
		Type:  gosnmp.OctetString,
		Value: sn,
	}

	keyPDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.2",
		Type:  gosnmp.OctetString,
		Value: key,
	}

	versusKeyPDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.3",
		Type:  gosnmp.OctetString,
		Value: versusKey,
	}

	typePDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.4",
		Type:  gosnmp.OctetString,
		Value: trapType,
	}

	severityPDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.5",
		Type:  gosnmp.OctetString,
		Value: severity,
	}

	descriptionPDU := gosnmp.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0.6",
		Type:  gosnmp.OctetString,
		Value: description,
	}

	return &gosnmp.SnmpTrap{
		Variables: []gosnmp.SnmpPDU{pdu, snPDU, keyPDU, versusKeyPDU, typePDU, severityPDU, descriptionPDU},
	}
}
