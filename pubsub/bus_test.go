package pubsub

import (
	"fmt"
	"testing"
)

func TestPubSub(t *testing.T) {
  bus := new(SimpleBus[string, bool])
  received := false

  bus.Subscribe(func(msg string) bool {
    fmt.Printf("Message: %s", msg)
    received = true
    return true
  })

  responses := bus.Send("Hello, World!")

  bus.Close()

  if resplen := len(responses); resplen != 1 {
    t.Errorf("Invalid number of responses: %d", resplen)
  }
  if !responses[0] {
    t.Errorf("Invalid response received")
  }
  if !received {
    t.Error("Not received")
  }
  if bus.Opened() {
    t.Error("Bus left opened")
  }
}
