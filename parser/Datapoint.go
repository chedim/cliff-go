package parser

import "cliff/pubsub"

type RejectionReason *string

type DatapointChange func() (Datapoint, AValue)

type Datapoint interface {
  Path([]string) Datapoint
  BeforeChange(pubsub.BusListener[DatapointChange, bool])
  AfterChange(pubsub.BusListener[DatapointChange, any])
  EnsureCapacity(int)
  Capacity() int
  Value() AValue
  Set(AValue) RejectionReason
  Slice(int) []AValue
  Child(string) Datapoint
  Type() Type
  Name() string
  String() string
}

var root = &MemoryDatapoint{
  children: make(map[string]Datapoint, 0),
}

func DatapointByName(name []string) (r Datapoint) {
  return root.Path(name)
}

