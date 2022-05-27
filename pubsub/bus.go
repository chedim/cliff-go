package pubsub

import "cliff/core"

type BusMessage interface{}
type BusListener[M BusMessage, R BusMessage] func(M) R

type ABus[M BusMessage, R BusMessage] interface {
	Send(M) []R
	Subscribe(BusListener[M, R])
  Close()
	Opened() bool
}

type SimpleBus[M BusMessage, R BusMessage] struct {
	closed    bool
	listeners []BusListener[M, R]
}

func (b SimpleBus[M, R]) Send(v M) []R{
  core.LogDebug("Sending message ", v, " to ", len(b.listeners), " subscribers")
	if b.closed {
		panic("Unable to send on closed bus")
	}
  responses := make([]R, 0)
	for i, listener := range b.listeners {
		if listener != nil {
      core.LogDebug("Sending message to listener #", i)
      responses = append(responses, listener(v))
		}
	}
  return responses
}

func (b *SimpleBus[M, R]) Subscribe(l BusListener[M, R]) {
  core.LogDebug("Creating new subscriber...")
	var target int = -1
	for i, cell := range b.listeners {
		if cell == nil {
			target = i
			b.listeners[i] = l
		}
	}
	if target == -1 {
    b.listeners = append(b.listeners, l)
		target = len(b.listeners)
	}

  core.LogDebug("Listeners: ", b.listeners)
}

func (b *SimpleBus[M, R]) Close() {
	b.closed = true
	for i, _ := range b.listeners {
		b.listeners[i] = nil
	}
}

func (b SimpleBus[M, R]) Opened() bool {
	return !b.closed
}
