package notify

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionMultipleEvents(t *testing.T) {
	const timeout = 1 * time.Second
	ctx := context.Background()

	r := require.New(t)
	b := NewBroadcaster()

	sub := SubscriptionTo(ctx, b)
	defer sub.Stop()

	const numNotifies = 10
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numNotifies; i++ {
			b.Broadcast()
		}
	}()
	wg.Wait()

	for i := 0; i < numNotifies; i++ {
		select {
		case <-sub.Register(ctx):
			t.Logf("Received notification")
		case <-ctx.Done():
			r.Failf("timeout", "timed out after %s", timeout)
		}
	}
}

func TestSubscriptionClose(t *testing.T) {
	t.Fatal("TODO")
}
