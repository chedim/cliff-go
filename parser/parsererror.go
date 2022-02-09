package parser

import (
	"errors"
	"fmt"
	"strings"
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
  debug := pe.Location.String()
  return fmt.Sprint("\n", debug, "\n", strings.Repeat(" ", len(debug) - 1), "^-- ", pe.error.Error(), "\n")
}
