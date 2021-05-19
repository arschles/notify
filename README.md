# Notify

Notify is a small library to safely do basic in-memory notifications between goroutines in Go. It
provides two major modes of functionality, detailed below.

>Official Go documentation can be found at [pkg.go.dev/github.com/arschles/notify].

## One-time broadcasting

One-time broadcasting allows goroutines to broadcast an "event" to many others. Any goroutine
can listen for events, but once they receive an event, they won't receive any more of them on the
`chan` that they're listening on. They must re-register (see "Multi-use broadcasting" below)

Once
the broadcast is complete, all  associated with the broadcaster become no-ops, which 
is why this is called one-time use. Here's an example how to use this:

```go
// create a new broadcaster
bcaster := NewBroadcaster()

// create 10 listeners
for i := 0; i < 10; i++ {
    go func(bcaster *Broadcaster, i int) {
        // ch will be closed when an event happens
        ch := bcaster.Register()
        <-ch
        log.Printf("Listener %d received the event!")
    }(bcaster, i)
}

// send the event to the 10 listeners
bcaster.Broadcast()
```

## Multi-use broadcasting

TODO
