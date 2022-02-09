package parser

import "fmt"

type Span struct {
	Start       int
	Length      int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
  Debug       string
}

func (s *Span) Extend(o *Span) *Span {
	if o.Start < s.Start {
		return o.Extend(s)
	}

	var cp Span = *s
	cp.Length += o.Length
	cp.EndLine = o.EndLine
	cp.EndColumn = o.EndColumn
	return &cp
}

func (s *Span) String() string {
  return fmt.Sprintf("[%d:%d] : %s", s.StartLine, s.StartColumn, s.Debug)
}
