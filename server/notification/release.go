package notification

// Release channel and connections used in message queue.
func Release() {
	channel.Close()
	connection.Close()
}
