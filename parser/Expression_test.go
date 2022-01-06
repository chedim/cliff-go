package parser

import (
	"strings"
	"testing"
)

func TestSumExpression(t *testing.T) {
  scanner := NewCliffScanner(strings.NewReader("2 + 2"))
  expr, err := ReadExpression(scanner)
  if err != nil {
    t.Errorf("%e", err.error)
  }

  res := expr.Value()
  if res != Integer(4) {
    t.Errorf("invalid result value: %s", res)
  }
}

func TestSubExpression(t *testing.T) {
  scanner := NewCliffScanner(strings.NewReader("2 - 2"))
  expr, err := ReadExpression(scanner)
  if err != nil {
    t.Errorf("%e", err.error)
  }

  res := expr.Value()
  if res != Integer(0) {
    t.Errorf("invalid result value: %s", res)
  }
}
