package parser

import (
	"strings"
	"testing"
)

func TestReadStatement(t *testing.T) {
  r := strings.NewReader("application output is \"Hello, world\" when true")
  s := NewCliffScanner(r)
  statement, err := ReadStatement(s)
  if err != nil {
    t.Errorf("failed to read a statement: %s", err)
  }

  target := statement.Target()
  if len(target.names) != 2 || target.names[0] != "application" || target.names[1] != "output" {
    t.Errorf("invalid target: '%s'", target.names)
  }

  if len(statement.definitions) != 1 {
    t.Errorf("invalid number of definitions: %d", len(statement.definitions))
  }

  def := statement.definitions[0]
  expected := String("Hello, world")
  actual, ok := def.Value().(String)
  if !ok {
    t.Errorf("invalid definition type: %T", def.Value())
  }
  if actual != expected {
    t.Errorf("invalid value in definition: %s", actual)
  }

  if def.Condition() == nil {
    t.Errorf("empty condition")
  }

  cond := def.Condition().Value()
  if cv, ok := cond.(Bool); !cv || !Bool(ok) {
    if !ok {
      t.Errorf("invalid condition value type: %T", cond)
    }
  }
}
