package types

type Package struct { // TODO -- make fields exported and immutable
	Name string
	IsExplicit bool
	Version string
	Size uint
	Deps []string // NOTE -- change this to []*Package or map[string]*Packge ?
}

