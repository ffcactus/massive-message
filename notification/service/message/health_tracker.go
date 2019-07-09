// Copyright

// This file includes the functionlity that tracking the state change of the target. The target here means the object with the URL in the notification.

package message

import (
	"bytes"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"massive-message/notification/repository"
	"massive-message/notification/sdk"
	"time"
)

const (
	retryOnErrorInterval = 60
)

// StartHealthTracker should be used as a co-routing. For each of the url, it finds out all effective alerts and broadcast this information.
func StartHealthTracker() {
	// args: exchange, type, durable, auto-deleted, internal, no-wait, args
	if err := channel.ExchangeDeclare(sdk.HealthChangeExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification-Service] Start health tracker failed, create exchange failed.")
		return
	}
	for {
		urls, err := repository.GetTargetsHaveAlert()
		// On error wait 1 minute and retry.
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn(fmt.Sprintf("[Notification-Service] Tracker health change failed, get URLs failed, retry after %d seconds.", retryOnErrorInterval))
			time.Sleep(retryOnErrorInterval * time.Second)
			continue
		}
		for _, url := range urls {
			notification, err := repository.CombineAlertsByURL(url)
			if err != nil {
				log.WithFields(log.Fields{"url": url, "error": err}).Warn("[Notification-Service] Tracker health change for URL failed, combine alerts failed.")
				continue
			}
			sendHealthChangeNotification(notification)
		}
	}
}

// find out the current health from the alerts.
func sendHealthChangeNotification(notification *sdk.HealthChangeNotification) {
	network := bytes.Buffer{}
	encoder := gob.NewEncoder(&network)

	if err := encoder.Encode(notification); err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("[Server-Service] Encoding health change notification message failed.")
		return
	}
	if err := channel.Publish(sdk.HealthChangeExchangeName, "HealthChange.New", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        network.Bytes(),
	}); err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("[Server-Service] Publish health change notification message failed.")
	}
}
