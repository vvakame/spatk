// Code generated by sig -private -output model_gen.go .; DO NOT EDIT

package sigtest

import (
	"fmt"
	"github.com/vvakame/spatk/scur"
)

var spannerInfoModelA = spannerInfoModelATable{
	name: "ModelA",
	columns: []spannerInfoModelAColumn{
		{name: "ModelAID"},
		{name: "Name"},
		{name: "UpdatedAt"},
		{name: "CreatedAt"},
	},
}

type spannerInfoModelATable struct {
	name    string
	alias   string
	columns []spannerInfoModelAColumn
}

type spannerInfoModelAColumn struct {
	name  string
	alias string
}

func (table spannerInfoModelATable) TableName() string {
	if table.alias != "" {
		return fmt.Sprintf("%s AS %s", table.name, table.alias)
	}
	return table.name
}
func (table spannerInfoModelATable) As(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}
func (table spannerInfoModelATable) ColumnNames() []string {
	return []string{
		table.ID(),
		table.Name(),
		table.UpdatedAt(),
		table.CreatedAt(),
	}
}
func (table spannerInfoModelATable) copy() spannerInfoModelATable {
	copied := table
	columns := make([]spannerInfoModelAColumn, len(table.columns))
	copy(columns, table.columns)
	copied.columns = columns
	return copied
}
func (table spannerInfoModelATable) ID() string {
	column := table.columns[0]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelATable) IDAs(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.columns[0].alias = aliasName
	return copied
}
func (table spannerInfoModelATable) IDCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.ID(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelA)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.ID
		},
	}
}
func (table spannerInfoModelATable) Name() string {
	column := table.columns[1]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelATable) NameAs(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.columns[1].alias = aliasName
	return copied
}
func (table spannerInfoModelATable) NameCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.Name(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelA)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.Name
		},
	}
}
func (table spannerInfoModelATable) UpdatedAt() string {
	column := table.columns[2]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelATable) UpdatedAtAs(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.columns[2].alias = aliasName
	return copied
}
func (table spannerInfoModelATable) UpdatedAtCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.UpdatedAt(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelA)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.UpdatedAt
		},
	}
}
func (table spannerInfoModelATable) CreatedAt() string {
	column := table.columns[3]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelATable) CreatedAtAs(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.columns[3].alias = aliasName
	return copied
}
func (table spannerInfoModelATable) CreatedAtCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.CreatedAt(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelA)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.CreatedAt
		},
	}
}

var spannerInfoModelB = spannerInfoModelBTable{
	name: "ModelB",
	columns: []spannerInfoModelBColumn{
		{name: "ModelBID"},
		{name: "Name"},
		{name: "UpdatedAt"},
		{name: "CreatedAt"},
	},
}

type spannerInfoModelBTable struct {
	name    string
	alias   string
	columns []spannerInfoModelBColumn
}

type spannerInfoModelBColumn struct {
	name  string
	alias string
}

func (table spannerInfoModelBTable) TableName() string {
	if table.alias != "" {
		return fmt.Sprintf("%s AS %s", table.name, table.alias)
	}
	return table.name
}
func (table spannerInfoModelBTable) As(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}
func (table spannerInfoModelBTable) ColumnNames() []string {
	return []string{
		table.ID(),
		table.Name(),
		table.UpdatedAt(),
		table.CreatedAt(),
	}
}
func (table spannerInfoModelBTable) copy() spannerInfoModelBTable {
	copied := table
	columns := make([]spannerInfoModelBColumn, len(table.columns))
	copy(columns, table.columns)
	copied.columns = columns
	return copied
}
func (table spannerInfoModelBTable) ID() string {
	column := table.columns[0]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelBTable) IDAs(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.columns[0].alias = aliasName
	return copied
}
func (table spannerInfoModelBTable) IDCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.ID(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelBar)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.ID
		},
	}
}
func (table spannerInfoModelBTable) Name() string {
	column := table.columns[1]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelBTable) NameAs(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.columns[1].alias = aliasName
	return copied
}
func (table spannerInfoModelBTable) NameCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.Name(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelBar)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.Name
		},
	}
}
func (table spannerInfoModelBTable) UpdatedAt() string {
	column := table.columns[2]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelBTable) UpdatedAtAs(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.columns[2].alias = aliasName
	return copied
}
func (table spannerInfoModelBTable) UpdatedAtCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.UpdatedAt(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelBar)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.UpdatedAt
		},
	}
}
func (table spannerInfoModelBTable) CreatedAt() string {
	column := table.columns[3]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}
func (table spannerInfoModelBTable) CreatedAtAs(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.columns[3].alias = aliasName
	return copied
}
func (table spannerInfoModelBTable) CreatedAtCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:  table.CreatedAt(),
		Order: order,
		ToValue: func(obj interface{}) interface{} {
			v, ok := obj.(*ModelBar)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.CreatedAt
		},
	}
}