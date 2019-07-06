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
	Key         string
	VersusKey   string
	URL         string
	Type        string
	GeneratedAt time.Time
	Severity    string
	Description string
}
