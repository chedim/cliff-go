package net

import (
	"cliff/parser"
	"encoding/json"
	"strings"
)

type QueryMessage struct {
	targetName *string
  targetHash *[]byte
}

func ReadQueryMessage(message string) *QueryMessage {
	if message[1] == '"' || message[1] == '\'' {
    var target string
    if json.Unmarshal([]byte(message[1:]), &target) == nil {
      return &QueryMessage{
        targetName: &target,
      }
    }
	} else if len(message) == 41 {
    hash := []byte(message[1:])
    return &QueryMessage{
      targetHash: &hash,
    }
  }
  return nil
}


type ValueMessage struct {
  QueryMessage
  value *parser.AValue
}

func ReadValueMessage(m string) *ValueMessage {
  qm := ValueMessage{}

  var s *parser.Scanner
  if m[1] == '"' || m[1] == '\'' {
    if json.Unmarshal([]byte(m[1:]), qm.targetName) != nil {
      return nil
    }
    s = parser.NewCliffScanner(strings.NewReader(m[len(*qm.targetName):]))
  } else {
    hash := []byte(m[1:41])
    qm.targetHash = &hash
    s = parser.NewCliffScanner(strings.NewReader(m[42:]))
  }

  if t, e := parser.ReadExpression(s); e != nil {
    return nil
  } else {

  }
}
