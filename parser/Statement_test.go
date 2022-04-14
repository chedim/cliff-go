package parser

import (
	"reflect"
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
  } else {
    cond := def.Condition().Value()
    if cv, ok := cond.(Bool); !cv || !Bool(ok) {
      if !ok {
        t.Errorf("invalid condition value type: %T", cond)
      }
    }
  }
}

func TestPluralStatements(t *testing.T) {
  if stmt, err := ParseStatement("values are:\n 2\n 4\n \"6\""); err != nil {
    t.Errorf("Failed to read statement: %s", err)
  } else if name := stmt.Target.Target().Name(); name != "value" {
    t.Errorf("Unexpected target datapoint name: %s", name)
  } else if dfns := stmt.Definitions; len(dfns) != 3 {
    t.Errorf("Unexpected definitons number: %d", len(dfns))
  } else if dfns[0].Value().String() != "2" {
    t.Errorf("Invalid value at index 0: %s", dfns[0].Value().String())
  } else if dfns[0].Value().Type() != Type(reflect.Int64) {
    t.Errorf("Invalid value type at index 0: %+v", dfns[0].Value().Type())
  } else if dfns[1].Value().String() != "4" {
    t.Errorf("Invalid value at index 1: %s", dfns[1].Value().String())
  } else if dfns[1].Value().Type() != Type(reflect.Int64) {
    t.Errorf("Invalid value type at index 1: %+v", dfns[1].Value().Type())
  } else if dfns[2].Value().String() != "6" {
    t.Errorf("Invalid value at index 2: %s", dfns[2].Value().String())
  } else if dfns[2].Value().Type() != Type(reflect.String) {
    t.Errorf("Invalid value type at index 2: %+v", dfns[2].Value().Type())
  }
}
