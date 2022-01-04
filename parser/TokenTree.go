package parser;

type TokenAnalyzer func (f *SourceFile, s *Scanner) interface{}

type TokenTree struct {
  TokenAnalyzer
  children map[Token]*TokenTree
}

func (tt *TokenTree) Path(ts []Token) (r *TokenTree) {
  for i := 0; i < len(ts); i++ {
    tt = tt.Child(ts[i])
  }

  return tt
}

func (tt *TokenTree) Child(t Token) *TokenTree {
    if _, k := tt.children[t]; !k {
      tt.children[t] = new(TokenTree)
    }
    return tt.children[t]
}

func (tt *TokenTree) AddAnalyzer(a TokenAnalyzer, ts ...Token) (r *TokenTree) {
  r = tt.Path(ts)
  r.TokenAnalyzer = a
  return
}

func (tt *TokenTree) GetAnalyzer(ts []*Tokenized) TokenAnalyzer {
  c := tt
  for i := 0; i < len(ts); i++ {
    c = c.Child(ts[i].Token)
  }

  return c.TokenAnalyzer
}
