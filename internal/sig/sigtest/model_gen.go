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
	name       string
	alias      string
	forceIndex string
	columns    []spannerInfoModelAColumn
}

type spannerInfoModelAColumn struct {
	name  string
	alias string
}

func (table spannerInfoModelATable) TableName() string {
	tableName := table.name
	if table.forceIndex != "" {
		tableName = fmt.Sprintf("%s@{FORCE_INDEX=%s}", tableName, table.forceIndex)
	}
	if table.alias != "" {
		tableName = fmt.Sprintf("%s AS %s", tableName, table.alias)
	}
	return tableName
}

func (table spannerInfoModelATable) As(aliasName string) spannerInfoModelATable {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}

func (table spannerInfoModelATable) ForceIndex(indexName string) spannerInfoModelATable {
	copied := table.copy()
	copied.forceIndex = indexName
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
		Name:     table.ID(),
		Order:    order,
		MinValue: scur.StringMinValue(),
		MaxValue: scur.StringMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.Name(),
		Order:    order,
		MinValue: scur.StringMinValue(),
		MaxValue: scur.StringMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.UpdatedAt(),
		Order:    order,
		MinValue: scur.TimestampMinValue(),
		MaxValue: scur.TimestampMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.CreatedAt(),
		Order:    order,
		MinValue: scur.TimestampMinValue(),
		MaxValue: scur.TimestampMaxValue(),
		ToValue: func(obj any) any {
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
	name       string
	alias      string
	forceIndex string
	columns    []spannerInfoModelBColumn
}

type spannerInfoModelBColumn struct {
	name  string
	alias string
}

func (table spannerInfoModelBTable) TableName() string {
	tableName := table.name
	if table.forceIndex != "" {
		tableName = fmt.Sprintf("%s@{FORCE_INDEX=%s}", tableName, table.forceIndex)
	}
	if table.alias != "" {
		tableName = fmt.Sprintf("%s AS %s", tableName, table.alias)
	}
	return tableName
}

func (table spannerInfoModelBTable) As(aliasName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}

func (table spannerInfoModelBTable) ForceIndex(indexName string) spannerInfoModelBTable {
	copied := table.copy()
	copied.forceIndex = indexName
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
		Name:     table.ID(),
		Order:    order,
		MinValue: scur.StringMinValue(),
		MaxValue: scur.StringMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.Name(),
		Order:    order,
		MinValue: scur.StringMinValue(),
		MaxValue: scur.StringMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.UpdatedAt(),
		Order:    order,
		MinValue: TimestampMinValue(),
		MaxValue: TimestampMaxValue(),
		ToValue: func(obj any) any {
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
		Name:     table.CreatedAt(),
		Order:    order,
		MinValue: TimestampMinValue(),
		MaxValue: TimestampMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelBar)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.CreatedAt
		},
	}
}

var spannerInfoModelC = spannerInfoModelCTable{
	name: "ModelC",
	columns: []spannerInfoModelCColumn{
		{name: "ModelCID"},
		{name: "OwnTimeType"},
		{name: "UUID"},
		{name: "LocalType1"},
		{name: "LocalType2"},
	},
}

type spannerInfoModelCTable struct {
	name       string
	alias      string
	forceIndex string
	columns    []spannerInfoModelCColumn
}

type spannerInfoModelCColumn struct {
	name  string
	alias string
}

func (table spannerInfoModelCTable) TableName() string {
	tableName := table.name
	if table.forceIndex != "" {
		tableName = fmt.Sprintf("%s@{FORCE_INDEX=%s}", tableName, table.forceIndex)
	}
	if table.alias != "" {
		tableName = fmt.Sprintf("%s AS %s", tableName, table.alias)
	}
	return tableName
}

func (table spannerInfoModelCTable) As(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) ForceIndex(indexName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.forceIndex = indexName
	return copied
}

func (table spannerInfoModelCTable) ColumnNames() []string {
	return []string{
		table.ID(),
		table.OwnTimeType(),
		table.UUID(),
		table.LocalType1(),
		table.LocalType2(),
	}
}

func (table spannerInfoModelCTable) copy() spannerInfoModelCTable {
	copied := table
	columns := make([]spannerInfoModelCColumn, len(table.columns))
	copy(columns, table.columns)
	copied.columns = columns
	return copied
}
func (table spannerInfoModelCTable) ID() string {
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

func (table spannerInfoModelCTable) IDAs(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.columns[0].alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) IDCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:     table.ID(),
		Order:    order,
		MinValue: scur.StringMinValue(),
		MaxValue: scur.StringMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelC)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.ID
		},
	}
}

func (table spannerInfoModelCTable) OwnTimeType() string {
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

func (table spannerInfoModelCTable) OwnTimeTypeAs(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.columns[1].alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) OwnTimeTypeCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:     table.OwnTimeType(),
		Order:    order,
		MinValue: TimeSpannerMinValue(),
		MaxValue: TimeSpannerMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelC)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.OwnTimeType
		},
	}
}

func (table spannerInfoModelCTable) UUID() string {
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

func (table spannerInfoModelCTable) UUIDAs(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.columns[2].alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) UUIDCursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:     table.UUID(),
		Order:    order,
		MinValue: UUIDSpannerMinValue(),
		MaxValue: UUIDSpannerMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelC)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.UUID
		},
	}
}

func (table spannerInfoModelCTable) LocalType1() string {
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

func (table spannerInfoModelCTable) LocalType1As(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.columns[3].alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) LocalType1Cursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:     table.LocalType1(),
		Order:    order,
		MinValue: localTypeSpannerMinValue(),
		MaxValue: localTypeSpannerMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelC)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.LocalType1
		},
	}
}

func (table spannerInfoModelCTable) LocalType2() string {
	column := table.columns[4]
	columnName := column.name
	if table.alias != "" {
		columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
	}
	if column.alias != "" {
		return fmt.Sprintf("%s AS %s", columnName, column.alias)
	}
	return columnName
}

func (table spannerInfoModelCTable) LocalType2As(aliasName string) spannerInfoModelCTable {
	copied := table.copy()
	copied.columns[4].alias = aliasName
	return copied
}

func (table spannerInfoModelCTable) LocalType2Cursor(order scur.Order) *scur.CursorParameter {
	return &scur.CursorParameter{
		Name:     table.LocalType2(),
		Order:    order,
		MinValue: localTypePointerSpannerMinValue(),
		MaxValue: localTypePointerSpannerMaxValue(),
		ToValue: func(obj any) any {
			v, ok := obj.(*ModelC)
			if !ok || v == nil {
				panic(fmt.Sprintf("unexpected cursor object type: %T", obj))
			}
			return v.LocalType2
		},
	}
}
