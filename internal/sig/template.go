package sig

const structTemplate = `
{{- $st := . -}}
var {{ $st.VarPrefix }}{{ $st.TableName }} = {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}{
	name: "{{ $st.TableName }}",
	columns: []{{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.ColumnTypeSuffix }}{
{{- range $idx, $f := $st.Fields}}
		{name: "{{$f.ColumnName}}"},
{{- end}}
	},
}

type {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }} struct {
	name    string
	alias   string
	columns []{{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.ColumnTypeSuffix }}
}

type {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.ColumnTypeSuffix }} struct {
	name  string
	alias string
}

func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) {{ $st.TableNameMethod }}() string {
	if table.alias != "" {
		return fmt.Sprintf("%s AS %s", table.name, table.alias)
	}
	return table.name
}
func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) As(aliasName string) {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }} {
	copied := table.copy()
	copied.alias = aliasName
	return copied
}
func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) {{ $st.ColumnNamesMethod }}() []string {
	return []string{
		{{- range $idx, $f := $st.Fields}}
			table.{{$f.Name}}(),
		{{- end}}
	}
}
func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) copy() {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }} {
	copied := table
	columns := make([]{{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.ColumnTypeSuffix }}, len(table.columns))
	copy(columns, table.columns)
	copied.columns = columns
	return copied
}


{{- range $idx, $f := $st.Fields}}
	func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) {{$f.Name}}() string {
		column := table.columns[{{ $idx }}]
		columnName := column.name
		if table.alias != "" {
			columnName = fmt.Sprintf("%s.%s", table.alias, columnName)
		}
		if column.alias != "" {
			return fmt.Sprintf("%s AS %s", columnName, column.alias)
		}
		return columnName
	}
	func (table {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }}) {{$f.Name}}As(aliasName string) {{ $st.VarPrefix }}{{ $st.TableName }}{{ $st.TableTypeSuffix }} {
		copied := table.copy()
		copied.columns[{{ $idx }}].alias = aliasName
		return copied
	}
{{- end}}
`
