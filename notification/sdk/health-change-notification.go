package sdk

// HealthChangeNotification represents the health change to the URL.
type HealthChangeNotification struct {
	URL      string // The target object that should pay attention to this notification.
	Critical int    // The critical alerts that still in effect state.
	Warning  int    // the warning alerts that still in effect state.
}
