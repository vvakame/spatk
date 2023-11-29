package sidx

type Index struct {
	Name         string
	Table        string
	SQL          string
	Unique       bool
	NullFiltered bool
}
