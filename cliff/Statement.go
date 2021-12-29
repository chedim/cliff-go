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
  statement = new(Statement)
  statement.location = scanner.Position()
	dtoks := scanner.scanWords()
  dname := NormalizedText(dtoks)
  statement.target = DatapointReference(dname)

  optoks := scanner.scanKeywords()
  optext := Text(optoks)

  return
}

func (s *Statement) Dependencies() []*Datapoint {
  
}
