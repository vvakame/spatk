package sqb

// NewBuilder can use for Spanner statement building.
func NewBuilder() Builder {
	return &builder{}
}

// Builder provides functionality of Spanner statement building.
type Builder interface {
	VoidBuilder
	Select() SelectBuilder
	Delete() DeleteBuilder
	From() FromBuilder
	Where() WhereBuilder
	OrderBy() OrderByBuilder
	Limit(limit string) VoidBuilder
}

// SelectBuilder provides functionality of SELECT statement building.
type SelectBuilder interface {
	VoidBuilder
	AsStruct() SelectBuilder
	C(name string, at ...string) SelectBuilder
	CS(names ...string) SelectBuilder
	From() FromBuilder
}

// DeleteBuilder provides functionality of DELETE statement building.
type DeleteBuilder interface {
	VoidBuilder
	From() FromBuilder
}

// FromBuilder provides functionality of FROM statement building.
type FromBuilder interface {
	VoidBuilder
	Name(tableName string, at ...string) FromBuilder
	Where() WhereBuilder
}

// WhereBuilder provides functionality of WHERE statement building.
type WhereBuilder interface {
	VoidBuilder
	E(token ...string) WhereBuilder
	OrderBy() OrderByBuilder
}

// OrderByBuilder provides functionality of ORDER BY statement building.
type OrderByBuilder interface {
	VoidBuilder
	O(token ...string) OrderByBuilder
	Limit(limit string) VoidBuilder
}

// VoidBuilder provides Build termination method.
type VoidBuilder interface {
	Build() (string, error)
}
