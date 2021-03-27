package notify

import (
	"sync"
	"time"
)

// SingleUseBroadcast is a multi-use notification intended for one goroutine to broadcast an event to one or more others.
// Listener goroutines get a notification from a signal channel (a <-chan struct{}) and then must throw the channel away
// and request a new one if they would like to get another signal.
//
// Broadcasters are similar to cancellable contexts, except they can be re-registered to and contexts cannot
//
// Create a new Broadcaster with NewBroadcaster.
type Broadcaster struct {
	chsM *sync.RWMutex
	chs  []chan struct{}
}

// NewBroadcaster creates a new Broadcaster ready for callers to subscribe to (Register) or broadcast
// on (Notify)
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		chsM: new(sync.RWMutex),
		chs:  nil,
	}
}

// CloseableBroadcaster is a broadcaster type that can be closed. Most commonly, these are created by
// calls to NewTickBroadcaster.
type CloseableBroadcaster struct {
	*Broadcaster
	t *time.Ticker
}

// Close immedaitely calls Notify(c.Broadcaster) and then releases any background resources in use
func (c *CloseableBroadcaster) Close() {
	c.Notify()
	c.t.Stop()
}

// NewTickBroadcaster returns a broadcaster that notifies every dur. After each time the returned
// CloseableBroadcaster notifies, callers must re-Register() and get a new channel to listen on
func NewTickBroadcaster(dur time.Duration) *CloseableBroadcaster {
	ticker := time.NewTicker(dur)
	b := NewBroadcaster()
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			b.Notify()
		}
	}()
	return &CloseableBroadcaster{
		Broadcaster: b,
		t:           ticker,
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

// Notify notifies all listeners on r. After Notify returns, all channels previously issued by Register
// will be closed and all listeners will be notified
func (b *Broadcaster) Notify() {
	b.chsM.Lock()
	defer b.chsM.Unlock()
	for _, ch := range b.chs {
		close(ch)
	}
	b.chs = nil
}
