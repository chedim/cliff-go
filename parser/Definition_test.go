package parser

import (
	"strings"
	"testing"
)

func TestReadDefinition(t *testing.T) {
  text := "is true when true"
  scanner := NewCliffScanner(strings.NewReader(text))

  def, err := ReadDefinition(scanner)
  if err != nil {
    t.Errorf("unexpected parser error: %s", err)
  }

  definition := def.Definition()
  if definition.Value() != Bool(true) {
    t.Errorf("invalid definition value: %s", definition.Value())
  }

  condition := def.Condition()
  if condition.Value() != Bool(true) {
    t.Errorf("invalid condition value: %s", condition.Value())
  }

  text = "is 30 when something other"
  scanner = NewCliffScanner(strings.NewReader(text))
  def, err = ReadDefinition(scanner)
  if err != nil {
    t.Errorf("unexpected parser error at %s", err)
  }
}
