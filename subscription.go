package notify

import (
	"context"
	"time"
)

type Subscription struct {
	ctx  context.Context
	done func()
	c    chan struct{}
}

// SubscriptionTo creates a new Subscription from a broadcaster. The caller should call Stop()
// on the returned Subscription when they're done, to release background resources
func SubscriptionTo(ctx context.Context, b *Broadcaster) *Subscription {

	ctx, done := context.WithCancel(ctx)
	ret := &Subscription{
		ctx:  ctx,
		done: done,
		c:    make(chan struct{}),
	}

	go func() {
		for {
			nCh := b.Register()
			select {
			case <-nCh:
				ret.send()
			case <-ctx.Done():
				return
			}
		}
	}()
	return ret
}

func TickingSubscription(ctx context.Context, tickDur time.Duration) *Subscription {
	ticker := time.NewTicker(tickDur)
	b := NewBroadcaster()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				b.Broadcast()
			}
		}
	}()

	return SubscriptionTo(ctx, b)
}

func (s *Subscription) send() bool {
	select {
	case <-s.ctx.Done():
		return false
	case s.c <- struct{}{}:
		return true
	}
}

// Block will block until a broadcast on the underlying broadcaster happens or Stop() is called
func (s *Subscription) Block() {
	<-s.c
}

// Stop immediately causes all calls to Block() to unblock, and stops listening to
// broadcaster on the underlying broadcaster. Callers should not use the Subscription
// after calling Stop
func (s *Subscription) Stop() {
	s.done()
	close(s.c)
}
