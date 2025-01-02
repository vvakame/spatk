package main

import (
	"context"
	"errors"
	"flag"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vvakame/spatk/internal/sidx"
	"golang.org/x/tools/go/packages"
)

var (
	command       = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	packageIdent  = command.String("packageIdent", "", "package ident; default extract from -output value")
	varNamePrefix = command.String("varNamePrefix", "", "var name prefix; default spannerIndex")
	output        = command.String("output", "", "output file name; default srcdir/model_spanner_index.go")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sidx: ")
	err := command.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = realMain()
	if err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	args := command.Args()
	if len(args) == 0 {
		return errors.New("1 argument requires")
	}

	var dir string
	var err error
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	if *output == "" {
		baseName := "model_spanner_index.go"
		*output = filepath.Join(dir, strings.ToLower(baseName))
	}
	log.Println(*output)

	if *packageIdent == "" {
		cfg := &packages.Config{
			Mode: packages.NeedName,
			Dir:  filepath.Dir(*output),
			Fset: token.NewFileSet(),
		}

		pkgs, err := packages.Load(cfg)
		if err != nil {
			log.Fatal(err)
		}

		if len(pkgs) != 1 {
			log.Fatal("1 package expected")
		}

		*packageIdent = pkgs[0].Name
	}

	if *varNamePrefix == "" {
		*varNamePrefix = "spannerIndex"
	}

	b, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	b, err = sidx.Build(context.Background(), &sidx.Config{
		PackageIdent:  *packageIdent,
		DDL:           string(b),
		VarNamePrefix: *varNamePrefix,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(*output, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
