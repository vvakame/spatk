package scur

import "strconv"

// Order means sort order in spanner query.
type Order int

const (
	// OrderAsc represent of ASC order.
	OrderAsc Order = iota
	// OrderDesc represent of DESC order.
	OrderDesc
)

func (order Order) String() string {
	switch order {
	case OrderAsc:
		return "ASC"
	case OrderDesc:
		return "DESC"
	default:
		return strconv.Itoa(int(order))
	}
}
