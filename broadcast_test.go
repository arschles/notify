package notify

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func iter(n int, fn func(i int) error) error {
	for i := 0; i < n; i++ {
		if err := fn(i); err != nil {
			return err
		}
	}
	return nil
}

func parFunc(i int, fn func(i int) error) func() error {
	return func() error {
		return fn(i)
	}
}

func iterPar(ctx context.Context, n int, fn func(i int) error) error {
	errgrp, _ := errgroup.WithContext(ctx)
	for i := 0; i < n; i++ {
		errgrp.Go(parFunc(i, fn))
	}
	return errgrp.Wait()
}

func TestBroadcast(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)

	b := NewBroadcaster()

	// register a bunch of listeners
	const numListeners = 10
	chans := make([]<-chan struct{}, numListeners)
	iter(10, func(i int) error {
		chans[i] = b.Register()
		return nil
	})

	// make sure they all don't receive before calling notify
	iterPar(ctx, 10, func(i int) error {
		t := time.NewTimer(200 * time.Millisecond)
		ch := chans[i]
		select {
		case <-ch:
			return fmt.Errorf("signal %d received before notify", i)
		case <-t.C:
			return nil
		}
	})

	notifiedCh := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Second)
		b.Broadcast()
		close(notifiedCh)
	}()
	<-notifiedCh

	iter(numListeners, func(i int) error {
		ch := chans[i]
		t := time.NewTimer(2 * time.Second)
		defer t.Stop()
		select {
		case <-ch:
		case <-t.C:
			r.Fail("didn't receive on notified channel within 1s")
		}
		return nil
	})
}
