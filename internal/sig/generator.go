package sig

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	"github.com/favclip/genbase"
)

type BuildSource struct {
	g         *genbase.Generator
	pkg       *genbase.PackageInfo
	typeInfos genbase.TypeInfos

	Structs []*BuildStruct
}

type BuildStruct struct {
	parent   *BuildSource
	typeInfo *genbase.TypeInfo

	Private bool
	Fields  []*BuildField
}

type BuildField struct {
	parent    *BuildStruct
	fieldInfo *genbase.FieldInfo

	Name string
	Tag  *BuildTag
}

type BuildTag struct {
	field *BuildField

	TableName string // e.g. `sig:"table,FooTable"`
	Name      string
	Ignore    bool // e.g. Secret string `spanner:"-"`
}

// Parse construct *BuildSource from package & type information.
// deprecated. use *BuildSource#Parse instead.
func Parse(pkg *genbase.PackageInfo, typeInfos genbase.TypeInfos) (*BuildSource, error) {
	bu := &BuildSource{
		g:         genbase.NewGenerator(pkg),
		pkg:       pkg,
		typeInfos: typeInfos,
	}

	for _, typeInfo := range typeInfos {
		err := bu.parseStruct(typeInfo)
		if err != nil {
			return nil, err
		}
	}

	return bu, nil
}

// Parse construct *BuildSource from package & type information.
func (b *BuildSource) Parse(pkg *genbase.PackageInfo, typeInfos genbase.TypeInfos) error {
	if b.g == nil {
		b.g = genbase.NewGenerator(pkg)
	}
	b.pkg = pkg
	b.typeInfos = typeInfos

	for _, typeInfo := range typeInfos {
		err := b.parseStruct(typeInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BuildSource) parseStruct(typeInfo *genbase.TypeInfo) error {
	structType, err := typeInfo.StructType()
	if err != nil {
		return err
	}

	st := &BuildStruct{
		parent:   b,
		typeInfo: typeInfo,
	}

	for _, fieldInfo := range structType.FieldInfos() {
		if len := len(fieldInfo.Names); len == 0 {
			// embedded struct in outer struct or multiply field declarations
			// https://play.golang.org/p/bcxbdiMyP4
			continue
		}

		for _, nameIdent := range fieldInfo.Names {
			err := b.parseField(st, typeInfo, fieldInfo, nameIdent.Name)
			if err != nil {
				return err
			}
		}
	}

	b.Structs = append(b.Structs, st)

	return nil
}

func (b *BuildSource) parseField(st *BuildStruct, typeInfo *genbase.TypeInfo, fieldInfo *genbase.FieldInfo, name string) error {
	field := &BuildField{
		parent:    st,
		fieldInfo: fieldInfo,
		Name:      name,
	}
	st.Fields = append(st.Fields, field)

	tag := &BuildTag{
		field: field,
		Name:  name,
	}
	field.Tag = tag

	if fieldInfo.Tag != nil {
		// remove back quote
		tagBody := fieldInfo.Tag.Value[1 : len(fieldInfo.Tag.Value)-1]
		tagKeys := genbase.GetKeys(tagBody)
		structTag := reflect.StructTag(tagBody)
		for _, key := range tagKeys {
			if key == "spanner" {
				tagText := structTag.Get("spanner")
				if tagText == "-" {
					tag.Ignore = true
					continue
				}
				tag.Name = tagText
			} else if key == "sig" {
				tagText := structTag.Get("sig")
				if strings.HasPrefix(tagText, "table,") {
					tag.TableName = tagText[len("table,"):]
				}
			}
		}
	}

	return nil
}

// Emit generate wrapper code.
func (b *BuildSource) Emit(args *[]string) ([]byte, error) {
	b.g.AddImport("fmt", "")
	b.g.AddImport("strings", "")
	b.g.AddImport("github.com/vvakame/spatk/sidx", "")
	b.g.AddImport("github.com/vvakame/spatk/scur", "")

	b.g.PrintHeader("sig", args)

	for _, st := range b.Structs {
		err := st.emit(b.g)
		if err != nil {
			return nil, err
		}
	}

	return b.g.Format()
}

func (st *BuildStruct) emit(g *genbase.Generator) error {
	tmpl := template.New("struct")
	tmpl, err := tmpl.Parse(structTemplate)
	if err != nil {
		return err
	}

	var varPrefix string
	if st.Private {
		varPrefix = "spannerInfo"
	} else {
		varPrefix = "SpannerInfo"
	}
	var fields []*BuildField
	for _, f := range st.Fields {
		if f.Tag.Ignore {
			continue
		}
		fields = append(fields, f)
	}
	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, map[string]interface{}{
		"VarPrefix":         varPrefix,
		"TableTypeSuffix":   "Table",
		"ColumnTypeSuffix":  "Column",
		"TableNameMethod":   "TableName",
		"ColumnNamesMethod": "ColumnNames",
		"SimpleName":        st.SimpleName(),
		"TableName":         st.TableName(),
		"Fields":            fields,
	})
	if err != nil {
		return err
	}

	g.Printf("%s", buf.String())

	return nil
}

// SimpleName returns struct type name.
func (st *BuildStruct) SimpleName() string {
	return st.typeInfo.Name()
}

// TableName returns table name from struct.
func (st *BuildStruct) TableName() string {
	for _, field := range st.Fields {
		if field.Tag.TableName != "" {
			return field.Tag.TableName
		}
	}
	return st.SimpleName()
}

// ColumnName returns column name from field.
func (fi *BuildField) ColumnName() string {
	if fi.Tag != nil && fi.Tag.Name != "" {
		return fi.Tag.Name
	}
	return fi.Name
}
