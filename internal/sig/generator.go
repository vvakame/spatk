package sig

import (
	"bytes"
	"fmt"
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

	TableName    string // e.g. `sig:"table=FooTable"`
	Name         string
	MinValueFunc string // e.g. `sig:"minValue=TimestampMinValue"`
	MaxValueFunc string // e.g. `sig:"maxValue=TimestampMaxValue"`
	Ignore       bool   // e.g. Secret string `spanner:"-"`
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
			switch key {
			case "spanner":
				tagText := structTag.Get("spanner")
				if tagText == "-" {
					tag.Ignore = true
					continue
				}
				tag.Name = tagText
			case "sig":
				tagText := structTag.Get("sig")
				attrs := strings.Split(tagText, ",")
				for _, attr := range attrs {
					attr = strings.TrimSpace(attr)
					switch {
					case strings.HasPrefix(attr, "table="):
						tag.TableName = attr[len("table="):]
					case strings.HasPrefix(attr, "minValue="):
						tag.MinValueFunc = attr[len("minValue="):]
					case strings.HasPrefix(attr, "maxValue="):
						tag.MaxValueFunc = attr[len("maxValue="):]
					default:
						return fmt.Errorf("unsupported attribute: %s", attr)
					}
				}
			}
		}
	}

	return nil
}

// Emit generate wrapper code.
func (b *BuildSource) Emit(args *[]string) ([]byte, error) {
	b.g.AddImport("fmt", "")
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

	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, st)
	if err != nil {
		return err
	}

	g.Printf("%s", buf.String())

	return nil
}

func (st *BuildStruct) VarPrefix() string {
	if st.Private {
		return "spannerInfo"
	}
	return "SpannerInfo"
}

func (st *BuildStruct) TableTypeSuffix() string {
	return "Table"
}

func (st *BuildStruct) ColumnTypeSuffix() string {
	return "Column"
}

func (st *BuildStruct) TableNameMethod() string {
	return "TableName"
}

func (st *BuildStruct) ColumnNamesMethod() string {
	return "ColumnNames"
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

func (st *BuildStruct) EnabledFields() []*BuildField {
	var fields []*BuildField
	for _, f := range st.Fields {
		if f.Tag.Ignore {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

// ColumnName returns column name from field.
func (fi *BuildField) ColumnName() string {
	if fi.Tag != nil && fi.Tag.Name != "" {
		return fi.Tag.Name
	}
	return fi.Name
}

func (fi *BuildField) MinValueFunc() string {
	if fi.Tag != nil && fi.Tag.MinValueFunc != "" {
		return fi.Tag.MinValueFunc
	}
	return ""
}

func (fi *BuildField) MaxValueFunc() string {
	if fi.Tag != nil && fi.Tag.MaxValueFunc != "" {
		return fi.Tag.MaxValueFunc
	}
	return ""
}
