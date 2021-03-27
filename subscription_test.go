package notify

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionMultipleEvents(t *testing.T) {
	const timeout = 1 * time.Second
	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()
	r := require.New(t)
	b := NewBroadcaster()

	sub := SubscriptionTo(ctx, b)

	const numNotifies = 10
	go func() {
		for i := 0; i < numNotifies; i++ {
			b.Broadcast()
		}
	}()

	for i := 0; i < numNotifies; i++ {
		select {
		case <-sub.Register(ctx):
		case <-ctx.Done():
			r.Failf("timeout", "timed out after %s", timeout)
		}
	}

}
