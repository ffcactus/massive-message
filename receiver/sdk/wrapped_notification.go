package sdk

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"net"
	"time"
)

const (
	// TrapExchangeName defines the exchange name for trap message.
	TrapExchangeName = "TrapExchange"
)

// SnmpVariable represents the SNMP OID value.
type SnmpVariable struct {
	// Name is an oid in string format eg ".1.3.6.1.4.9.27"
	Name string
	// The type of the value eg Integer
	Type gosnmp.Asn1BER
	// The value to be set by the SNMP set, or the value when
	// sending a trap
	Value interface{}
}

// String returns the string representation of SnmpVariable
func (v SnmpVariable) String() string {
	switch v.Type {
	case gosnmp.OctetString:
		return fmt.Sprintf("%v", string(v.Value.([]byte)))
	case gosnmp.Integer:
		return fmt.Sprintf("%v", v.Value.(int))
	case gosnmp.TimeTicks:
		return fmt.Sprintf("%v", v.Value.(uint))
	default:
		return fmt.Sprintf("%v", v.Value)
	}
}

// WrapedSnmpPacket includes both the raw SNMP packet but also some other useful information for processing it later.
// Since we are using encoding/gob, it's OK to use point here.
type WrapedSnmpPacket struct {
	Address     *net.UDPAddr
	GeneratedAt time.Time
	Variables   []SnmpVariable
}
