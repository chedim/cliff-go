package cliff;

type Package struct {
  location string
}

func (p *Package) Location() string {
  return p.location
}

