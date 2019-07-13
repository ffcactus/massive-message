package snmp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	notificationSDK "massive-message/notification/sdk"
	receiverSDK "massive-message/receiver/sdk"
)

// SnmpVariable represents the variable in the snmp notification.
type Variable struct {
	// Name is an oid in string format eg ".1.3.6.1.4.9.27"
	Name string
	// The type of the value eg Integer
	Type gosnmp.Asn1BER
	// The value to be set by the SNMP set, or the value when
	// sending a trap
	Value interface{}
}

func (v Variable) String() string {
	switch v.Type {
	case gosnmp.OctetString:
		return string(v.Value.([]byte))
	case gosnmp.Integer:
		return fmt.Sprintf("%d", v.Value.(int))
	case gosnmp.TimeTicks:
		return fmt.Sprintf("%d", v.Value.(uint))
	default:
		log.WithFields(log.Fields{"type": Asn1BERToString(v.Type)}).Warn("Convert SnmpVariable to string failed.")
		return ""
	}
}

// Converter represents the method that a converter should have.
type Converter interface {
	Convert(packet *receiverSDK.WrapedSnmpPacket) ([]notificationSDK.Notification, error)
}

// generateKey generates the standard notifications' key
func generateKey(vendor, originalKey string) string {
	return fmt.Sprintf("%s-%s", vendor, originalKey)
}

// Asn1BERToString converts the Asn1BER type to string
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
