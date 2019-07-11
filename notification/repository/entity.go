package repository

import (
	"github.com/google/uuid"
	"massive-message/notification/sdk"
	"time"
)

// Event represents the table for event objects.
type Event struct {
	ID          string `gorm:"column:ID;primary_key"`
	Key         string `gorm:"index"`
	URL         string `gorm:"index"`
	VersusKey   string
	Type        string
	GeneratedAt time.Time
	Severity    string
	Description string
}

// Alert represents the table for alert objects.
type Alert struct {
	ID          string `gorm:"primary_key"`
	Key         string
	URL         string
	VersusKey   string
	Type        string
	GeneratedAt time.Time
	Severity    string
	Description string
}

// URLHealth represents the table saving the warning and critical count to a URL.
type URLHealth struct {
	URL       string `gorm:"primary_key"`
	Warnings  int
	Criticals int
}

func newAlert(o *sdk.Notification) *Alert {
	ret := Alert{}
	ret.ID = uuid.New().String()
	ret.Key = o.Key
	ret.VersusKey = o.VersusKey
	ret.URL = o.URL
	ret.Type = o.Type
	ret.GeneratedAt = o.GeneratedAt
	ret.Severity = o.Severity
	ret.Description = o.Description
	return &ret
}

func newEvent(o *sdk.Notification) *Event {
	ret := Event{}
	ret.ID = uuid.New().String()
	ret.Key = o.Key
	ret.VersusKey = o.VersusKey
	ret.URL = o.URL
	ret.Type = o.Type
	ret.GeneratedAt = o.GeneratedAt
	ret.Severity = o.Severity
	ret.Description = o.Description
	return &ret
}
