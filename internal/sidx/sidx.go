package sidx

import (
	"bytes"
	"context"
	"errors"
	goformat "go/format"
	"text/template"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
)

type Config struct {
	PackageIdent  string
	DDL           string
	VarNamePrefix string
}

type TemplateValue struct {
	Config
	Indices []*ast.CreateIndex
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

	ddlStmts, err := memefish.ParseDDLs("-", cfg.DDL)
	if err != nil {
		return nil, err
	}

	var indices []*ast.CreateIndex
	for _, ddlStmt := range ddlStmts {
		indexAST, ok := ddlStmt.(*ast.CreateIndex)
		if !ok {
			continue
		}

		indices = append(indices, indexAST)
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
