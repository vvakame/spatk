package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/favclip/genbase"
	"github.com/vvakame/spatk/internal/sig"
)

var (
	command   = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	typeNames = command.String("type", "", "comma-separated list of type names; must be set")
	output    = command.String("output", "", "output file name; default srcdir/<type>_spanner_info.go")
	private   = command.Bool("private", false, "generated type name; export or unexport")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprint(os.Stderr, "\tsig [flags] [directory]\n")
	_, _ = fmt.Fprint(os.Stderr, "\tsig [flags] files... # Must be a single package\n")
	_, _ = fmt.Fprint(os.Stderr, "Flags:\n")
	command.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("sig: ")
	command.Usage = Usage
	err := command.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var dir string
	var pInfo *genbase.PackageInfo
	p := &genbase.Parser{SkipSemanticsCheck: true}
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		pInfo, err = p.ParsePackageDir(dir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dir = filepath.Dir(args[0])
		pInfo, err = p.ParsePackageFiles(args)
		if err != nil {
			log.Fatal(err)
		}
	}

	var typeInfos genbase.TypeInfos
	if len(*typeNames) == 0 {
		typeInfos = pInfo.CollectTaggedTypeInfos("+sig")
	} else {
		typeInfos = pInfo.CollectTypeInfos(strings.Split(*typeNames, ","))
	}

	if len(typeInfos) == 0 {
		flag.Usage()
	}

	bu := sig.BuildSource{}
	err = bu.Parse(pInfo, typeInfos)
	if err != nil {
		log.Fatal(err)
	}
	for _, st := range bu.Structs {
		st.Private = *private
	}

	// Format the output.
	src, err := bu.Emit(nil)
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Print("warning: compile the package to analyze the error")
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_spanner_info.go", typeInfos[0].Name())
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
