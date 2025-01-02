package sig

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

const markerComment = "+sig"

type PackageInfo struct {
	pkg *packages.Package

	CommandArgs []string
	ImportSpecs []*ImportSpec
	Structs     []*StructInfo
}

type ImportSpec struct {
	Path  string
	Ident string
}

type StructInfo struct {
	parent  *PackageInfo
	Name    string
	Private bool
	Fields  []*FieldInfo
}

type FieldInfo struct {
	parent *StructInfo

	Name                string
	defaultMinValueFunc string
	defaultMaxValueFunc string
	Tag                 *TagInfo
}

type TagInfo struct {
	field *FieldInfo

	TableName    string // e.g. `sig:"table=FooTable"`
	Name         string
	MinValueFunc string // e.g. `sig:"minValue=TimestampMinValue"`
	MaxValueFunc string // e.g. `sig:"maxValue=TimestampMaxValue"`
	Ignore       bool   // e.g. Secret string `spanner:"-"`
}

func Parse(directoryPath string) (*PackageInfo, error) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  directoryPath,
		Fset: token.NewFileSet(),
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, errors.New("1 package expected")
	}

	pkg := pkgs[0]

	if len(pkg.Errors) != 0 {
		var err error
		for _, pkgErr := range pkg.Errors {
			err = errors.Join(err, pkgErr)
		}
		return nil, err
	}

	packageInfo := &PackageInfo{
		pkg: pkg,

		ImportSpecs: []*ImportSpec{
			{Path: "fmt"},
			{Path: "github.com/vvakame/spatk/scur"},
		},
	}

	var retErr error
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			genDecl, ok := n.(*ast.GenDecl)
			if !ok {
				return true
			}

			if !strings.Contains(genDecl.Doc.Text(), markerComment) {
				return true
			}

			structInfos, err := packageInfo.parseStruct(genDecl)
			if err != nil {
				retErr = errors.Join(retErr, err)
				return false
			}

			packageInfo.Structs = append(packageInfo.Structs, structInfos...)

			return true
		})
	}
	if retErr != nil {
		return nil, retErr
	}

	return packageInfo, nil
}

func (pkgInfo *PackageInfo) parseStruct(decl *ast.GenDecl) ([]*StructInfo, error) {
	var structInfos []*StructInfo
	var retErr error
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			retErr = errors.Join(
				retErr,
				fmt.Errorf(
					"%s: non type spec has +sig comment",
					pkgInfo.pkg.Fset.Position(spec.Pos()).String(),
				),
			)
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			retErr = errors.Join(
				retErr,
				fmt.Errorf(
					"%s: non struct type has +sig comment",
					pkgInfo.pkg.Fset.Position(typeSpec.Pos()).String(),
				),
			)
			continue
		}

		st := &StructInfo{
			parent: pkgInfo,
			Name:   typeSpec.Name.Name,
		}

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				// embedded struct in outer struct or multiply field declarations
				// https://play.golang.org/p/bcxbdiMyP4
				continue
			}

			fieldInfos, err := pkgInfo.parseField(st, field)
			if err != nil {
				retErr = errors.Join(retErr, err)
				continue
			}

			st.Fields = append(st.Fields, fieldInfos...)
		}

		structInfos = append(structInfos, st)
	}
	if retErr != nil {
		return nil, retErr
	}

	return structInfos, nil
}

func (pkgInfo *PackageInfo) parseField(st *StructInfo, field *ast.Field) ([]*FieldInfo, error) {
	var fieldInfos []*FieldInfo
	var retErr error
OUTER:
	for _, name := range field.Names {
		fieldInfo := &FieldInfo{
			parent: st,
			Name:   name.Name,
		}

		typ := pkgInfo.pkg.TypesInfo.TypeOf(field.Type)
		switch {
		case isBasicKindOrUnderlying(typ, types.String):
			fieldInfo.defaultMinValueFunc = "scur.StringMinValue"
			fieldInfo.defaultMaxValueFunc = "scur.StringMaxValue"
		case isTimeTime(typ):
			fieldInfo.defaultMinValueFunc = "scur.TimestampMinValue"
			fieldInfo.defaultMaxValueFunc = "scur.TimestampMaxValue"
		default:
			var minValueFunc string
			var maxValueFunc string
			if namedType, ok := typ.(*types.Named); ok {
				typeName := namedType.Obj().Name()
				minValueFunc = typeName + "SpannerMinValue"
				maxValueFunc = typeName + "SpannerMaxValue"
			} else if pointerType, ok := typ.(*types.Pointer); ok {
				if namedType, ok := pointerType.Elem().(*types.Named); ok {
					typeName := namedType.Obj().Name()
					minValueFunc = typeName + "PointerSpannerMinValue"
					maxValueFunc = typeName + "PointerSpannerMaxValue"
				}
			}
			if minValueFunc == "" || maxValueFunc == "" {
				continue
			}

			for ident, def := range pkgInfo.pkg.TypesInfo.Defs {
				// function can be min/max value generator.
				fn, ok := def.(*types.Func)
				if !ok {
					continue
				}
				if fn.Signature().Recv() != nil {
					continue
				}

				rets := fn.Signature().Results()
				if rets.Len() != 1 {
					continue
				}
				ret := rets.At(0)

				switch ident.Name {
				case minValueFunc:
					fieldInfo.defaultMinValueFunc = minValueFunc
				case maxValueFunc:
					fieldInfo.defaultMaxValueFunc = maxValueFunc
				default:
					continue
				}

				if !types.AssignableTo(ret.Type(), typ) {
					retErr = errors.Join(
						retErr,
						fmt.Errorf(
							"%s: function %s return type is not same as field type: %s",
							pkgInfo.pkg.Fset.Position(ret.Pos()).String(),
							ident.Name,
							typ.String(),
						),
					)
					continue OUTER
				}
			}
		}

		tagInfo, err := pkgInfo.parseTag(fieldInfo, field)
		if err != nil {
			retErr = errors.Join(retErr, err)
			continue
		}

		fieldInfo.Tag = tagInfo

		fieldInfos = append(fieldInfos, fieldInfo)
	}
	if retErr != nil {
		return nil, retErr
	}

	return fieldInfos, nil
}

func (pkgInfo *PackageInfo) parseTag(bf *FieldInfo, field *ast.Field) (*TagInfo, error) {
	tagInfo := &TagInfo{
		field: bf,
	}

	if field.Tag == nil {
		return tagInfo, nil
	}

	// remove back quote
	tagBody := field.Tag.Value[1 : len(field.Tag.Value)-1]
	structTag := reflect.StructTag(tagBody)

	{
		tagText := structTag.Get("spanner")
		if tagText == "-" {
			tagInfo.Ignore = true
		} else {
			tagInfo.Name = tagText
		}
	}
	{
		tagText := structTag.Get("sig")
		if tagText != "" {
			attrs := strings.Split(tagText, ",")
			for _, attr := range attrs {
				attr = strings.TrimSpace(attr)
				switch {
				case strings.HasPrefix(attr, "table="):
					tagInfo.TableName = attr[len("table="):]

				case strings.HasPrefix(attr, "min="):
					tagInfo.MinValueFunc = attr[len("min="):]

				case strings.HasPrefix(attr, "max="):
					tagInfo.MaxValueFunc = attr[len("max="):]

				case strings.HasPrefix(attr, "minmax="):
					prefix := attr[len("minmax="):]
					tagInfo.MinValueFunc = prefix + "MinValue"
					tagInfo.MaxValueFunc = prefix + "MaxValue"

				default:
					return nil, fmt.Errorf("unsupported attribute: %s", attr)
				}
			}
		}
	}

	return tagInfo, nil
}

func (pkgInfo *PackageInfo) Emit() ([]byte, error) {
	tmpl := template.New("file")
	tmpl, err := tmpl.Parse(fileTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, pkgInfo)
	if err != nil {
		return nil, err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}

	return src, nil
}

func (pkgInfo *PackageInfo) CommandArguments() string {
	return strings.Join(pkgInfo.CommandArgs, " ")
}

func (pkgInfo *PackageInfo) PackageIdent() string {
	return pkgInfo.pkg.Name
}

func (structInfo *StructInfo) VarPrefix() string {
	if structInfo.Private {
		return "spannerInfo"
	}
	return "SpannerInfo"
}

func (structInfo *StructInfo) TableTypeSuffix() string {
	return "Table"
}

func (structInfo *StructInfo) ColumnTypeSuffix() string {
	return "Column"
}

func (structInfo *StructInfo) TableNameMethod() string {
	return "TableName"
}

func (structInfo *StructInfo) ColumnNamesMethod() string {
	return "ColumnNames"
}

// SimpleName returns struct type name.
func (structInfo *StructInfo) SimpleName() string {
	return structInfo.Name
}

// TableName returns table name from struct.
func (structInfo *StructInfo) TableName() string {
	for _, field := range structInfo.Fields {
		if field.Tag.TableName != "" {
			return field.Tag.TableName
		}
	}
	return structInfo.SimpleName()
}

func (structInfo *StructInfo) EnabledFields() []*FieldInfo {
	var fields []*FieldInfo
	for _, f := range structInfo.Fields {
		if f.Tag.Ignore {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

// ColumnName returns column name from field.
func (fieldInfo *FieldInfo) ColumnName() string {
	if fieldInfo.Tag != nil && fieldInfo.Tag.Name != "" {
		return fieldInfo.Tag.Name
	}
	return fieldInfo.Name
}

func (fieldInfo *FieldInfo) MinValueFunc() string {
	if fieldInfo.Tag != nil && fieldInfo.Tag.MinValueFunc != "" {
		return fieldInfo.Tag.MinValueFunc
	}
	return fieldInfo.defaultMinValueFunc
}

func (fieldInfo *FieldInfo) MaxValueFunc() string {
	if fieldInfo.Tag != nil && fieldInfo.Tag.MaxValueFunc != "" {
		return fieldInfo.Tag.MaxValueFunc
	}
	return fieldInfo.defaultMaxValueFunc
}
