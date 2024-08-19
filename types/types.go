package types

type Package struct {
	Name       string
	IsExplicit bool
	Version    string
	Size       uint
	Deps       []string // NOTE: change this to []*Package or map[string]*Packge ?
}
