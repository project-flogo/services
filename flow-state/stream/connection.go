package stream

// Monitors and keeps the state of the current connection
type ConnectionMonitor interface {
	IsAlive() bool
	ConnectionLost()
}

// The connection monitor
type connectionMonitor struct {
	ConnectionStateChannel chan int
}

var _ ConnectionMonitor = (*connectionMonitor)(nil)

// Constructor for the ConnectionMonitor
func NewConnectionMonitor(connectionStateChannel chan int) ConnectionMonitor {
	return &connectionMonitor{ConnectionStateChannel: connectionStateChannel}
}

// Returns true if the connection still alive, or false otherwise
func (monitor *connectionMonitor) IsAlive() bool {
	var alive bool = true
	// Try to write and read to the channel
	defer func() {
		if x := recover(); x != nil {
			alive = false
		}
	}()
	monitor.ConnectionStateChannel <- 1
	<-monitor.ConnectionStateChannel
	return alive
}

// Marks the connection as lost
func (monitor *connectionMonitor) ConnectionLost() {
	close(monitor.ConnectionStateChannel)
}
