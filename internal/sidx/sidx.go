package sidx

import (
	"bytes"
	"context"
	"errors"
	goformat "go/format"
	"text/template"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/MakeNowJust/memefish/pkg/token"
)

type Config struct {
	PackageIdent  string
	DDL           string
	VarNamePrefix string
}

func Build(ctx context.Context, cfg *Config) ([]byte, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	if cfg.PackageIdent == "" {
		return nil, errors.New("cfg.PackageIdent must be required")
	}
	if cfg.VarNamePrefix == "" {
		cfg.VarNamePrefix = "spannerIndex"
	}

	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{
				Buffer: cfg.DDL,
			},
		},
	}

	ddls, err := p.ParseDDLs()
	if err != nil {
		return nil, err
	}

	var indices []*ast.CreateIndex
	for _, ddl := range ddls {
		indexAST, ok := ddl.(*ast.CreateIndex)
		if !ok {
			continue
		}

		indices = append(indices, indexAST)
	}

	type TemplateValue struct {
		Config
		Indices []*ast.CreateIndex
	}
	v := &TemplateValue{
		Config:  *cfg,
		Indices: indices,
	}

	tmpl := template.New("file")
	tmpl, err = tmpl.Parse(fileTemplate)
	if err != nil {
		return nil, err
	}
	tmpl = tmpl.New("packageHeader")
	tmpl, err = tmpl.Parse(packageHeaderTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "file", v)
	if err != nil {
		return nil, err
	}

	return goformat.Source(buf.Bytes())
}
