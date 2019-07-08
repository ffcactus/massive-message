// Copyright

// This file includes the functionlity that tracking the state change of the target. The target here means the object with the URL in the notification.

package message

import (
	"fmt"
	"massive-message/notification/repository"
	"time"
)

// StartStateTracker should be used as a co-routing. It find out all the targets that has already sent the alerts.
// And then, for each of the target, calculates the statistic of the alerts.
// The statistic result includes the alerts that is still in effect state.
func StartStateTracker() {
	for {
		fmt.Println(repository.GetTargetsHaveAlert())
		time.Sleep(10 * time.Second)
	}
}
