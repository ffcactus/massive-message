package notification

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// HpeNotificationConverter is the converter for HPE's notification.
type HpeNotificationConverter struct {
}

// Convert implements the NotificationConverter interface.
func (HpeNotificationConverter) Convert(packet *WrapedSnmpPacket) ([]StandardNotification, error) {
	var (
		originalKey      string
		notificationType string
		serverity        string
		description      string
	)

	single := StandardNotification{}
	// For HPE's notification, take the event number as the original key.
	for _, v := range packet.Variables {
		switch v.Name {
		case "EventNumber":
			originalKey = v.String()
		case "Type":
			notificationType = v.String()
		case "Serverity":
			serverity = v.String()
		case "Description":
			description = v.String()
		}
	}
	if originalKey == "" || notificationType == "" || serverity == "" || description == "" {
		log.WithFields(log.Fields{"vender": "HPE", "Address": packet.Address.IP.String()}).Error("[Notification] Convert SNMP notification failed, drop this notification.")
		return nil, fmt.Errorf("no original key")
	}
	single.Key = generateKey("HPE", originalKey)
	single.ReceivedAt = packet.ReceivedAt
	single.Type = notificationType
	single.Severity = serverity
	single.Description = description

	return []StandardNotification{single}, nil
}
