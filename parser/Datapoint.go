package parser

import (
	"container/list"
	"math"
	"reflect"
)

type Datapoint struct {
	name   string
	values list.List
  maxValues *Stack

	children map[string]*Datapoint

	definition    []Definition
	beforeChangeListeners []*BeforeChangeListener
  afterChangeListeners []*AfterChangeListener
}

type RejectionReason *string
type BeforeChangeListener func(dp *Datapoint, v AValue) RejectionReason
type AfterChangeListener func(dp *Datapoint)

var root = &Datapoint{
  children: make(map[string]*Datapoint, 0),
}

func DatapointByName(name []string) (r *Datapoint) {
  return root.Path(name)
}

func (d *Datapoint) Path(path []string) (r *Datapoint) {
  r = d
  for _, subname := range path {
    r = r.Child(subname)
  }
  return
}

func (d *Datapoint) BeforeChange(s *BeforeChangeListener) {
	d.beforeChangeListeners = append(d.beforeChangeListeners, s)
}

func (d *Datapoint) AfterChange(s *AfterChangeListener) {
  d.afterChangeListeners = append(d.afterChangeListeners, s)
}

func (d *Datapoint) DeleteBeforeChange(s *BeforeChangeListener) {
  var retained []*BeforeChangeListener
  for _, listener := range d.beforeChangeListeners {
    if listener != s {
      retained = append(retained, listener)
    }
  }
}

func (d *Datapoint) Capacity() int {
  return d.maxValues.Peek().(int)
}

func (d *Datapoint) EnsureCapacity(request int) {
  if d.Capacity() < request {
    d.maxValues.Push(request)
  }
}

func (d *Datapoint) Value() AValue {
  if d.values.Len() == 0 {
    return nil
  }
  return d.values.Front().Value.(AValue)
}

func (d *Datapoint) Set(v AValue) RejectionReason {
  for _, listener := range d.beforeChangeListeners {
    rr := (*listener)(d, v)
    if rr != nil {
      return rr
    }
  }

  d.values.PushFront(v)
  if d.values.Len() > d.Capacity() {
    d.values.Remove(d.values.Back())
  }

  for _, listener := range d.afterChangeListeners {
    (*listener)(d)
  }

  return nil
}

func (d *Datapoint) Slice(max int) (result []AValue) {
  count := int(math.Min(float64(d.Capacity()), float64(max)))
  el := d.values.Front()
  for i := 0; i < count && el != nil; i++ {
    result = append(result, el.Value.(AValue))
    el = el.Next()
  }
  return
}

func (d *Datapoint) Child(name string) *Datapoint {
	result, present := d.children[name]
	if !present {
		d.children[name] = &Datapoint{
      name: name,
      children: make(map[string]*Datapoint, 0),
    }
		result = d.children[name]
	}
	return result
}

func (d *Datapoint) Type() Type {
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
