package env

import (
	"cliff/parser"
	"testing"
)

func TestDatapointOperations(t *testing.T) {
  dp := parser.NewDatapoint("test")
  SetDatapoint(dp)

  ndp := GetDatapointByName("test")
  if ndp == nil {
    t.Error("Failed to get test datapoint")
  }
}
