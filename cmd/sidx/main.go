package main

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/favclip/genbase"
	"github.com/vvakame/spatk/internal/sidx"
)

var (
	command       = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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

	var pInfo *genbase.PackageInfo
	var dir string
	var err error
	p := &genbase.Parser{SkipSemanticsCheck: true}
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		pInfo, err = p.ParsePackageDir(dir)
		if err != nil {
			return err
		}
	} else {
		dir = filepath.Dir(args[0])
		pInfo, err = p.ParsePackageDir(dir)
		if err != nil {
			return err
		}
	}

	if *varNamePrefix == "" {
		*varNamePrefix = "spannerIndex"
	}

	b, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	b, err = sidx.Build(context.Background(), &sidx.Config{
		PackageIdent:  pInfo.Name(),
		DDL:           string(b),
		VarNamePrefix: *varNamePrefix,
	})
	if err != nil {
		return err
	}

	outputName := *output
	if outputName == "" {
		baseName := "model_spanner_index.go"
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	log.Println(outputName)
	err = ioutil.WriteFile(outputName, b, 0644)
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
