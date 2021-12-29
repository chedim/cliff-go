package cliff

import "os"

type SourceFile struct {
	folder     *Package
	name       string
	statements []*Statement
}

func NewSourceFile(folder *Package, name string) (result *SourceFile, err error) {
  result := &SourceFile{folder: folder, name: name}
  file, err := os.Open(result.Location())
  if err == nil {
    scanner := NewCliffScanner(file)
    result.statements = make([]*Statement, 0)
    for statement, serr := ReadStatement(result, scanner); serr == nil; statement, serr = ReadStatement(result, scanner) {
      result.statements = append(result.statements, statement)
    }
  }
  return
}

func (me *SourceFile) Folder() *Package {
  return me.folder
}

func (me *SourceFile) Name() string {
  return me.name
}

func (me *SourceFile) Location() string {
  return me.folder.Location() + os.PathSeparator + me.name
}

func (me *SourceFile) Statements() *[]Statement {
  return &me.statements
}
