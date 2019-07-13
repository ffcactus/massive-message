package sdk

import (
	"time"
)

const (
	// NotificationExchangeName is the name of the exchange of notification.
	NotificationExchangeName = "NotificationExchange"
)

// Notification represents the notification passing through the modules of ths system.
type Notification struct {
	Key         string    // The unique ID if this notification.
	VersusKey   string    // If this notification is an alert, the VersusKey means the key to which the notification should be cleaned.
	URL         string    // The target that generates the notification.
	Type        string    // The type of the notification, should only be "Alert" or "Event".
	GeneratedAt time.Time // When this notification be generated.
	Severity    string    // The severity of this notification, should only be "OK", "Warning" and "Critical".
	Description string    // The string representation of this notification.
}

func (o Notification) String() string {
	return o.Key
}
