package parser

import (
	"cliff/pubsub"
	"container/list"
	"math"
	"reflect"
)

type MemoryDatapoint struct {
	name      string
	values    list.List
	maxValues *Stack

	children map[string]Datapoint

	definition            []Definition
  beforeChangeBus       pubsub.ABus[DatapointChange, bool]
  afterChangeBus        pubsub.ABus[DatapointChange, any]
}

func NewDatapoint(name string) Datapoint {
  return MemoryDatapoint{
    name: name,
  }
}

func (d MemoryDatapoint) Path(path []string) (r Datapoint) {
  r = &d
  for _, subname := range path {
    r = r.Child(subname)
  }
  return
}

func (d MemoryDatapoint) BeforeChange(l pubsub.BusListener[DatapointChange, bool]) {
  d.beforeChangeBus.Subscribe(l)
}

func (d MemoryDatapoint) AfterChange(l pubsub.BusListener[DatapointChange, any]) {
  d.afterChangeBus.Subscribe(l)
}

func (d MemoryDatapoint) Capacity() int {
  return d.maxValues.Peek().(int)
}

func (d MemoryDatapoint) EnsureCapacity(request int) {
  if d.Capacity() < request {
    d.maxValues.Push(request)
  }
}

func (d MemoryDatapoint) Value() AValue {
  if d.values.Len() == 0 {
    return nil
  }
  return d.values.Front().Value.(AValue)
}

func (d MemoryDatapoint) Set(v AValue) RejectionReason {
  d.beforeChangeBus.Send(func() (Datapoint, AValue) { return d, v })

  ov := d.Value()
  d.values.PushFront(v)
  if d.values.Len() > d.Capacity() {
    d.values.Remove(d.values.Back())
  }

  d.afterChangeBus.Send(func() (Datapoint, AValue) { return d, ov })

  return nil
}

func (d MemoryDatapoint) Slice(max int) (result []AValue) {
  count := int(math.Min(float64(d.Capacity()), float64(max)))
  el := d.values.Front()
  for i := 0; i < count && el != nil; i++ {
    result = append(result, el.Value.(AValue))
    el = el.Next()
  }
  return
}

func (d MemoryDatapoint) Child(name string) Datapoint {
	result, present := d.children[name]
	if !present {
		d.children[name] = &MemoryDatapoint{
      name: name,
      children: make(map[string]Datapoint, 0),
    }
		result = d.children[name]
	}
	return result
}

func (d MemoryDatapoint) Type() Type {
  if len(d.children) != 0 {
    return Type(reflect.Map)
  }

  if d.values.Len() != 0 {
    return Type(reflect.Array)
  }
  v := d.Value()
  if v == nil {
    return Type(reflect.Invalid)
  }

  return v.Type()
}

func (d MemoryDatapoint) Name() string {
  return d.name
}

func (d MemoryDatapoint) String() string {
  return d.name
}
