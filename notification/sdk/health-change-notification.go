package sdk

const (
	// HealthChangeExchangeName is the name of the exchange of health change notification.
	HealthChangeExchangeName = "HealthChangeExchange"
)

// HealthChangeNotification represents the health change to the URL.
type HealthChangeNotification struct {
	URL       string // The target object that should pay attention to this notification.
	Criticals int    // The critical alerts that still in effect state.
	Warnings  int    // the warning alerts that still in effect state.
}
