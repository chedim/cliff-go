package cliff

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
  if len(target) != 2 || target[0] != "application" || target[1] != "output" {
    t.Errorf("invalid target: '%s'", strings.Join(target, "."))
  }

  if len(statement.definitions) != 1 {
    t.Errorf("invalid number of definitions: %d", len(statement.definitions))
  }

  def := statement.definitions[0]
  var expected AValue = "Hello, world"
  if *def.Value() != expected {
    t.Errorf("invalid value in definition: %s", *def.Value())
  }

  if def.Condition == nil {
    t.Errorf("empty condition")
  }

  cond := *def.Condition.Value()
  if cv, ok := cond.(bool); !cv || !ok {
    if !ok {
      t.Errorf("invalid condition value type: %T", cond)
    }
  }
}
