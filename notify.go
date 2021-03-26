package notify

import "sync"

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

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		chsM: new(sync.RWMutex),
		chs:  nil,
	}
}

// Register indicates that the caller would like to listen to r. The returned channel will be closed when
// a notification arrives. After the notification, this channel is closed and nothing will happen
// to it again. If you want to re-register, call this function again and use the new return channel.
func Register(r *Broadcaster) <-chan struct{} {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	newCh := make(chan struct{})
	r.chs = append(r.chs, newCh)
	return newCh
}

// Notify notifies all listeners on r. After Notify returns, all channels previously issued by Register
// will be closed and all listeners will be notified
func Notify(r *Broadcaster) {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	for _, ch := range r.chs {
		close(ch)
	}
	r.chs = nil
}
