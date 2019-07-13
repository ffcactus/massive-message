package snmp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	notificationSDK "massive-message/notification/sdk"
	receiverSDK "massive-message/receiver/sdk"
)

// LenovoNotificationConverter is the converter for Lenovo's notification.
type LenovoNotificationConverter struct {
}

// Convert implements the NotificationConverter interface.
func (LenovoNotificationConverter) Convert(packet *receiverSDK.WrapedSnmpPacket) ([]notificationSDK.Notification, error) {
	var (
		sn               string
		originalKey      string
		versusKey        string
		notificationType string
		serverity        string
		description      string
	)

	single := notificationSDK.Notification{}
	// For Lenovo's notification, take the event number as the original key.
	for _, v := range packet.Variables {
		switch v.Name {
		case ".1.3.6.1.6.3.1.1.4.1.0.1":
			sn = v.String()
		case ".1.3.6.1.6.3.1.1.4.1.0.2":
			originalKey = v.String()
		case ".1.3.6.1.6.3.1.1.4.1.0.3":
			versusKey = v.String()
		case ".1.3.6.1.6.3.1.1.4.1.0.4":
			notificationType = v.String()
		case ".1.3.6.1.6.3.1.1.4.1.0.5":
			serverity = v.String()
		case ".1.3.6.1.6.3.1.1.4.1.0.6":
			description = v.String()
		}
	}
	if sn == "" || originalKey == "" || versusKey == "" || notificationType == "" || serverity == "" || description == "" {
		log.WithFields(log.Fields{"vender": "Lenovo", "Address": packet.Address.IP.String()}).Error("[Server-Notification] Convert SNMP notification failed, drop this notification.")
		return nil, fmt.Errorf("no original key")
	}
	single.URL = snURLMapping[sn]
	single.Key = generateKey("Lenovo", originalKey)
	single.VersusKey = generateKey("Lenovo", versusKey)
	single.GeneratedAt = packet.GeneratedAt
	single.Type = notificationType
	single.Severity = serverity
	single.Description = description

	return []notificationSDK.Notification{single}, nil
}
