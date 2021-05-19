package notify

import (
	"sync"
)

// Broadcaster is a single-use broadcasting system. One or more goroutines can listen to notifications
// any goroutine (usually not in the aforementioned group) can broadcast a signal to the listeners.
//
// After a broadcast happens, all listeners are automatically woken up and de-registered. If they want to
// continue listening to notifications, they must re-register. In this way, this Broadcaster type is
// similar to a cancellable context, but it can be re-registered to (contexts cannot easily).
//
// If you'd like to establish a Subscription that you do not have to re-register for every time a notification
// comes, see the Subscription type.
//
// Create a new Broadcaster with NewBroadcaster.
type Broadcaster struct {
	chsM *sync.Mutex
	chs  []chan struct{}
}

// NewBroadcaster creates a new Broadcaster ready for callers to subscribe to (Register) or broadcast
// on (Notify)
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		chsM: new(sync.Mutex),
		chs:  nil,
	}
}

// Register indicates that the caller would like to listen to r. The returned channel will be closed when
// a notification arrives. After the notification, this channel is closed and nothing will happen
// to it again. If you want to re-register, call this function again and use the new return channel.
func (b *Broadcaster) Register() <-chan struct{} {
	b.chsM.Lock()
	defer b.chsM.Unlock()
	newCh := make(chan struct{})
	b.chs = append(b.chs, newCh)
	return newCh
}

// Broadcast notifies all listeners on r. After Notify returns, all channels previously issued by Register
// will be closed and all listeners will be notified
func (b *Broadcaster) Broadcast() {
	b.chsM.Lock()
	defer b.chsM.Unlock()
	for _, ch := range b.chs {
		close(ch)
	}
	b.chs = nil
}
