package cliff

type StatementOperator interface {
  Apply(t *Datapoint, e *Expression)
}

type Statement struct {
  location *Span
  abstract bool
	target   *Datapoint
  operator *StatementOperator
  definition *Expression
}

func ReadStatement(source *SourceFile, scanner *Scanner) (statement *Statement, err error){
	tntoks := scanner.scanWords()
  tname := Text(tntoks)
  stttoks := scanner.scanKeywords()
  sttype := Text(stttoks)

  return
}

func (s *Statement) Dependencies() []*Datapoint {
  
}
