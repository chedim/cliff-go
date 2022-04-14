package parser

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"github.com/gertd/go-pluralize"
	"go.uber.org/zap"
)

var pc = pluralize.NewClient()

func TextArray(t []*Tokenized) []string {
  r := make([]string, len(t))
  for i := 0; i < len(t); i++ {
    r[i] = t[i].Literal
  }
  return r
}

func Text(t []*Tokenized) string {
  r := ""
  for i := 0; i < len(t); i++ {
    r += t[i].Literal
    if (i < len(t)) {
      r += " "
    }
  }
  return r
}

func NormalizedText(t []*Tokenized) string {
  r := ""
  for i := 0; i < len(t); i++ {
    r += pc.Singular(strings.ToLower(t[i].Literal))
  }
  return r;
}

func NormalizedTextArray(t []*Tokenized) []string {
  r := make([]string, len(t))
  for i := 0; i < len(t); i++ {
    r[i] = pc.Singular(strings.ToLower(t[i].Literal))
  }
  return r
}

func NormalizeTextArray(t []string) []string {
  r := make([]string, len(t))
  for i := 0; i < len(t); i++ {
    r[i] = pc.Singular(strings.ToLower(t[i]))
  }
  return r
}
type Stack struct {
  *list.List
  mut sync.Mutex
  len int
}

func NewStack() *Stack {
  return &Stack{List: list.New()}
}

func (s *Stack) String() string {
  r := ""
  i := 0
  for v := s.List.Back(); v != nil; v = v.Prev()  {
    r += fmt.Sprintf("%d = %#v\n", i, v.Value)
    i++
  }
  return r
}

func (s *Stack) Push(x interface{}) {
  s.mut.Lock()
  defer s.mut.Unlock()

  s.List.PushBack(x)
  s.len++
}

func (s *Stack) Pop() interface{} {
  s.mut.Lock()
  defer s.mut.Unlock()

  if s.Len() == 0 {
    return nil
  }
  s.len--
  tail := s.Back()
  val := tail.Value
  s.Remove(tail)
  return val
}

func (s *Stack) Peek() interface{} {
  s.mut.Lock()
  defer s.mut.Unlock()
  if s.Len() == 0 {
    return nil
  }

  return s.Front().Value
}

func (s *Stack) Len() int {
  return s.len
}

func LogDebug(args...interface{}) {
  logger, _ := zap.NewDevelopment()
  l := logger.Sugar()
  l.Debug(args...)
}
