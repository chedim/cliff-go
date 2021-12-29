package utils

import (
	"strings"

	"github.com/chewxy/lingo/lexer"
)

func NormalizeText(t string) string {
  lx := lexer.New("cliff", strings.NewReader(t))
}

func NormalizedText(t []*Tokenized) string {
  r := ""
  for i := 0; i < len(t); i++ {
    word := NormalizeText(t[i].Literal)
  }
}

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
