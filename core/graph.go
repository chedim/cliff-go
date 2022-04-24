package core

import (
	"cliff/parser"
	"crypto/sha1"
	"sync"
)

var data sync.Map

func init() {
  data = sync.Map{}
}

func GetDatapointByName(n string) *parser.Datapoint {
  return GetDatapoint(GetDatapointHash(n))
}

func GetDatapointHash(n string) []byte {
  h := sha1.New()
  h.Write([]byte(n))
  sha := h.Sum(nil)
  return sha
}

func GetDatapoint(sha []byte) *parser.Datapoint {
  if v, ok := data.Load(sha); ok {
    return &v
  }
  return nil
}

func SetDatapoint(dp *parser.Datapoint) {
  data.Store(GetDatapointHash(dp.Name()), *dp)
}
