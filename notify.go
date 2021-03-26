package notify

import "sync"

// Reusable is a multi-use notification intended for one goroutine to broadcast an event to one or more others.
// Listener goroutines must accept the notification within a pre-defined amount of time or the notification will
// be aborted. Once a notification goes out, all listener goroutines are automatically unregistered.
//
// All operations on Reusables are concurrency-safe.
//
// Create a new SingleUseBroadcast with NewSingleUseBroadcast.
type SingleUseBroadcast struct {
	chsM *sync.RWMutex
	chs  []chan struct{}
}

func NewSingleUseBroadcast() *SingleUseBroadcast {
	return &SingleUseBroadcast{
		chsM: new(sync.RWMutex),
		chs:  nil,
	}
}

// Register indicates that the caller would like to listen to r. The returned channel will be closed when
// a notification arrives. After the notification, this channel is closed and nothing will happen
// to it again. If you want to re-register, call this function again and use the new return channel.
func Register(r *SingleUseBroadcast) <-chan struct{} {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	newCh := make(chan struct{})
	r.chs = append(r.chs, newCh)
	return newCh
}

// Notify notifies all listeners on r. After Notify returns, all channels previously issued by Register
// will be closed and all listeners will be notified
func Notify(r *SingleUseBroadcast) {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	for _, ch := range r.chs {
		close(ch)
	}
	r.chs = nil
}
