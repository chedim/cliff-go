package cliff

type DefinitionStruct struct {
  valueExpression Expression
  activityExpression Expression
  dependencies []*Datapoint
}

func NewDefinition(parser *Parser) Definition {
	return DefinitionStruct{}
}

func (me DefinitionStruct) parse(parser *Parser) {
  valueExpression = NewValueExpression()
  tok, lit := parser.Scan()

  activityDependencies := me.activityExpression.Dependencies()
  valueDependencies := me.valueExpression.Dependencies()

  me.dependencies = make([]*Datapoint, len(activityDependencies) + len(valueDependencies))
  for i := 0; i < len(me.dependencies); i++ {
    if i < len(activityDependencies) {
      me.dependencies[i] = activityDependencies[i]
    } else {
      me.dependencies[i] = valueDependencies[i - len(activityDependencies)]
    }
  }
}

func (me DefinitionStruct) Active() bool {
	return true
}

func (me DefinitionStruct) Evaluate() *Value {
}

func (me DefinitionStruct) Dependencies() []*Datapoint {
  return 
}
