package sidx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	goformat "go/format"
	"strconv"
	"text/template"

	"cloud.google.com/go/spanner/spansql"
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

	ddl, err := spansql.ParseDDL("-", cfg.DDL)
	if err != nil {
		return nil, err
	}

	var indices []*spansql.CreateIndex
	for _, ddlStmt := range ddl.List {
		indexAST, ok := ddlStmt.(*spansql.CreateIndex)
		if !ok {
			continue
		}

		indices = append(indices, indexAST)
	}

	type TemplateValue struct {
		Config
		Indices []*spansql.CreateIndex
	}
	v := &TemplateValue{
		Config:  *cfg,
		Indices: indices,
	}

	tmpl := template.New("file").Funcs(map[string]interface{}{
		"quote": func(v interface{}) (string, error) {
			switch v := v.(type) {
			case string:
				return strconv.Quote(v), nil
			case spansql.ID:
				return strconv.Quote(string(v)), nil
			default:
				return "", fmt.Errorf("quote: unsupported type: %T", v)
			}
		},
	})
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
