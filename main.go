package main

import (
	"cliff/parser"
	"fmt"
	"strings"
)

func main() {
  fmt.Printf("Hello, world!")
  r := strings.NewReader("message is \"Hello\"\n")
  s := parser.NewCliffScanner(r)

  if stt, err := parser.ReadStatement(s); err != nil {
    fmt.Printf("Failed to parse statement: %v\n", err)
  } else {
    fmt.Printf("Statement: %s\n", stt)
  }
}
