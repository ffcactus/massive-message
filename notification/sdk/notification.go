package sdk

import (
	"time"
)

// Notification represents the notification passing through the modules of ths system.
type Notification struct {
	Key         string
	VersusKey   string
	URL         string
	Type        string
	ReceivedAt  time.Time
	Severity    string
	Description string
}
