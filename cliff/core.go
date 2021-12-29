package cliff

type Parseable interface {
  Parse()
}

type CliffFolder interface {
  Name() string
  Location() string
  SubFolders() *[]CliffFolder
  Files() *[]CliffFile
}

type CliffFile interface {
  Folder() *CliffFolder
  Name() string
  Location() string
	Statements() *[]Statement
}

type Datapoint interface {
	Path() []*string
  Definitions() []*Definition
  // future values
  Values() <-chan Value
  // returns the last value
  Value() Value
}

type Definition interface {
	Evaluate() *Value
  Active() bool
  Dependencies() []*Datapoint
}

type Value interface {
  Type() *Datapoint
}

type Expression interface {
  Type() *Datapoint
  Dependencies() []*Datapoint
  Values() <-chan Value
  Value() Value
}
