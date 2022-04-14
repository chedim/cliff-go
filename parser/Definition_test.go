package parser

import (
	"strings"
	"testing"
)

func TestReadDefinition(t *testing.T) {
  text := "true when true"
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
  if condition == nil || condition.Value() != Bool(true) {
    t.Errorf("invalid condition value: %s", condition)
  }

  text = "30 when something other"
  scanner = NewCliffScanner(strings.NewReader(text))
  def, err = ReadDefinition(scanner)
  if err != nil {
    t.Errorf("unexpected parser error at %s", err)
  }
}

func TestReversedDefinitions(t *testing.T) {
  text := "when name is set: 2"
  scanner := NewCliffScanner(strings.NewReader(text))
  if def, err := ReadDefinition(scanner); err != nil {
    t.Errorf("unexpected parser error at %s", err)
  } else if def.Value().String() != "2" {
    t.Errorf("Invalid definition value: %s", def.Value().String())
  }
}
