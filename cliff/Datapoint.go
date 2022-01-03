package cliff;

type Datapoint struct {
  name string
  values []interface{}

  children map[string]*Datapoint

  definition []Definition
  subscriptions []*DatapointSubscriber
}

type DatapointSubscriber interface {
  OnDatapointChanged(*Datapoint)
}

var root = new(Datapoint)

func DatapointByName(name []string) *Datapoint {
  r := root
  for i := 0; i < len(name); i++ {
    r = r.Child((name)[i])
  }
  return r
}

func (d *Datapoint) Subscribe(s *DatapointSubscriber) {
  d.subscriptions = append(d.subscriptions, s)
}

func (d *Datapoint) Child(name string) *Datapoint {
  result, present := d.children[name]
  if !present {
    d.children[name] = &Datapoint{name: name}
    result = d.children[name]
  }
  return result
}
