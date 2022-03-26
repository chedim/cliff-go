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

  target := statement.Target
  if len(target.names) != 2 || target.names[0] != "application" || target.names[1] != "output" {
    t.Errorf("invalid target: '%s'", target.names)
  }

  if len(statement.Definitions) != 1 {
    t.Errorf("invalid number of definitions: %d", len(statement.Definitions))
  }

  def := statement.Definitions[0]
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

func TestPluralStatements(t *testing.T) {
  r := strings.NewReader("apples are fruits")
  s := NewCliffScanner(r)

  statement, parserErr := ReadStatement(s)
  if parserErr != nil {
    t.Errorf("Failed to read statement: %s", parserErr)
  }
  target := statement.Target
  if len(target.names) != 1 || target.names[0] != "apple" {
    t.Errorf("invalid target: %s", target.names)
  }

  definitions := statement.Definitions
  if len(definitions) != 1 {
    t.Errorf("Invalid number of definitions: %d", len(definitions))
  }
  def := definitions[0]
  defval := def.Definition().(*Reference)
  if defval == nil {
    t.Errorf("Definition value is nil: %+v", def.value)
  }
  if len(defval.names) != 1 || defval.names[0] != "fruit" {
    t.Errorf("Definition reference is invalid: %v", defval.names)
  }
  if !defval.plural {
    t.Errorf("Definition value is not plural")
  }

  r = strings.NewReader("apples are the next 10 fruits")
  s = NewCliffScanner(r)

  statement, parserErr = ReadStatement(s)
  if parserErr != nil {
    t.Errorf("Failed to read statement: %s", parserErr)
  }
}
