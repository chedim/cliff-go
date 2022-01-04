package parser

import (
	"errors"
	"fmt"
)

type ParserError struct {
	error
	Location Span
}

func NewParserError(location Span, message string) *ParserError {
	e := new(ParserError)
	e.error = errors.New(message)
  e.Location = location
  return e
}

func WrapParserError(err *ParserError, message string) *ParserError {
  e := new(ParserError)
  e.error = err
  e.Location = err.Location
  return e
}

func ExtendParserError(location Span, err error) *ParserError {
  e := new(ParserError)
  e.error = err
  e.Location = location
  return e
}

func (pe *ParserError) Error() string {
  return fmt.Sprint("(", pe.Location.StartLine, ":", pe.Location.StartColumn, ") ", pe.error.Error())
}
