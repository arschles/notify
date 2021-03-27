package notify

import "context"

// Registrant is an interface representing someone that can register for notifications.
//
// Implementations may provide one-time registration, meaning that after the first notification
// receipt, the caller will not be registered anymore, or multi-notification registration, meaning
// that the caller will be registered for N notifications, or possibly until they opt-out
type Registrant interface {
	// Register indicates that the caller would like to listen for a notification/notifications.
	// The returned channel must receive when a notification arrives. Beyond that, the semantics
	// are implementation-specific. For example, one-time registration implementations will close
	// the channel, which will cause it to receive only once.
	Register(context.Context) <-chan struct{}
}
