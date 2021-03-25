package notify

import "sync"

// Reusable is a multi-use notification intended for one goroutine to broadcast an event to one or more others.
// Listener goroutines must accept the notification within a pre-defined amount of time or the notification will
// be aborted. Once a notification goes out, all listener goroutines are automatically unregistered.
//
// All operations on Reusables are concurrency-safe
type Reusable struct {
	chsM *sync.RWMutex
	chs  []chan struct{}
}

// Register indicates that the caller would like to listen to r. The returned channel will be closed when
// a notification arrives. After the notification, this channel is closed and nothing will happen
// to it again. If you want to re-register, call this function again and use the new return channel.
func Register(r *Reusable) <-chan struct{} {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	newCh := make(chan struct{})
	r.chs = append(r.chs, newCh)
	return newCh
}

// Notify notifies all listeners on r. After Notify returns, there will be no listeners on r.
func Notify(r *Reusable) {
	r.chsM.Lock()
	defer r.chsM.Unlock()
	for _, ch := range r.chs {
		close(ch)
	}
	r.chs = nil
}
