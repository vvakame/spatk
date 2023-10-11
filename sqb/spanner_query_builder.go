package sqb

type Option func(b *builder)

// NewBuilder can use for Spanner statement building.
func NewBuilder(opts ...Option) Builder {
	b := &builder{
		lineHead: true,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Builder provides functionality of Spanner statement building.
type Builder interface {
	VoidBuilder
	Select() SelectBuilder
	Update(tableName string, at ...string) UpdateBuilder
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

// UpdateBuilder provides functionality of UPDATE statement building.
type UpdateBuilder interface {
	VoidBuilder
	Set() SetBuilder
	Where() WhereBuilder
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

// SetBuilder provides functionality of SET clause building.
type SetBuilder interface {
	VoidBuilder
	U(token ...string) SetBuilder
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
