package sidx

import _ "embed"

//go:embed package_header.go.tmpl
var packageHeaderTemplate string

//go:embed file.go.tmpl
var fileTemplate string
