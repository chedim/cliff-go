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

func (m *QueryMessage) Apply() {

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
    value := t.Value()
    qm.value = &value
  }

  return &qm
}

type SubscribeMessage struct {
  QueryMessage
}

func ReadSubscribeMessage(s string) *SubscribeMessage {
  return &SubscribeMessage{
    *ReadQueryMessage(s),
  }
}

type UnsubscribeMessage struct {
  QueryMessage
}

func ReadUnsubscribeMessage(m string) *UnsubscribeMessage {
  return &UnsubscribeMessage{
    *ReadQueryMessage(m),
  }
}

type UnsetMessage struct {
  QueryMessage
}

func ReadUnsetMessage(m string) *UnsetMessage {
  return &UnsetMessage{
    *ReadQueryMessage(m),
  }
}
