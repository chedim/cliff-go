package parser

import (
	"strings"
	"testing"
)

func TestReadSlices(t *testing.T) {
  r := strings.NewReader("the user")
  s := NewCliffScanner(r)

  if slice, err := ReadSliceExpression(s); err != nil {
    t.Errorf("Failed to read slice expression: %v", err)
  } else if name := slice.Target().Target().Name(); name != "user" {
    t.Errorf("Invalid slice datapoint: %s", name)
  }

  r = strings.NewReader("the last 10 messages")
  s = NewCliffScanner(r)
  if slice, err := ReadSliceExpression(s); err != nil {
    t.Errorf("Failed to read slice expression: %v", err)
  } else if name := slice.Target().Target().Name(); name != "message" {
    t.Errorf("Invalid slice datapoint: %s", name)
  }

  r = strings.NewReader("apples are fruits")
  s = NewCliffScanner(r)

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
