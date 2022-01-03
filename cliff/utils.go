package cliff;

import (
	"strings"

	"github.com/gertd/go-pluralize"
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
