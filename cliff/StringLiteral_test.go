package cliff

import (
  "testing"
  "strings"
)

func TestReadString(t *testing.T) {
  const v = "\"tEst\\\"\""
  r := strings.NewReader(v)
  s := NewCliffScanner(r)

  result, err := ReadString(s, DQUOTE)
  if err != nil {
    t.Errorf("failed to read string: %s", err)
  }

  if result.value != "tEst\""{
    t.Errorf("Read string is invalid: '%s'", result.value)
  }
}
