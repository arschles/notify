package notify

import (
	"context"
	"log"
	"time"
)

// Subscription will maintain a subscription to a broadcaster
type Subscription struct {
	Registrant
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
				log.Printf("received event, sending")
				ret.send()
			case <-ctx.Done():
				log.Printf("context done (in SuscriptioTo)")
				return
			}
		}
	}()
	return ret
}

func TickingSubscription(ctx context.Context, tickDur time.Duration) *Subscription {
	ctx, done := context.WithCancel(ctx)
	ticker := time.NewTicker(tickDur)
	sub := &Subscription{
		ctx:  ctx,
		done: done,
		c:    make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sub.send()
			}
		}
	}()

	return sub
}

func (s *Subscription) send() bool {
	select {
	case <-s.ctx.Done():
		log.Printf("send failed, ctx done")
		return false
	case s.c <- struct{}{}:
		log.Printf("send success")
		return true
	}
}

// Register returns a channel that will receive (not close) every time a subscription happens. Since
// the channel receives once, only one goroutine should plan on using it. Generally speaking,
// if a goroutine wants to subscribe to a Broadcaster, it should plan on creating its own Subscription(s)
func (s *Subscription) Register(ctx context.Context) <-chan struct{} {
	return s.c
}

// Stop immediately causes all calls to Block() to unblock, and stops listening to
// broadcaster on the underlying broadcaster. Callers should not use the Subscription
// after calling Stop
func (s *Subscription) Stop() {
	s.done()
	close(s.c)
}
