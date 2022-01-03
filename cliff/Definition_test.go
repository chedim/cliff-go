package cliff

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

  val := def.Value()
  if def.Value() == nil {
    t.Error("definition value is nil")
  }
}
