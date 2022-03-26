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
}
