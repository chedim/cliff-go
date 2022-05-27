package env

import (
	"cliff/parser"
	"crypto/sha1"
	"sync"
)

var data sync.Map

func init() {
  data = sync.Map{}
}

func GetDatapointByName(n string) parser.Datapoint {
  return GetDatapoint(GetDatapointHash(n))
}

func GetDatapointHash(n string) string {
  h := sha1.New()
  h.Write([]byte(n))
  sha := h.Sum(nil)
  return string(sha)
}

func GetDatapoint(sha string) parser.Datapoint {
  if v, ok := data.Load(sha); ok {
    dp := v.(parser.Datapoint)
    return dp
  }
  return nil
}

func SetDatapoint(dp parser.Datapoint) {
  data.Store(GetDatapointHash(dp.Name()), dp)
}
