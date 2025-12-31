package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vvakame/spatk/internal/sig"
)

var (
	command     = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	output      = command.String("output", "", "output file name; default srcdir/<type>_spanner_info.go")
	private     = command.Bool("private", false, "generated type name; export or unexport")
	allowErrors = command.Bool("allow-errors", false, "allow code generation even if source has compile errors (methods for unknown types will panic)")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprint(os.Stderr, "\tsig [flags] [directory]\n")
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

	dir := args[0]
	if !isDirectory(dir) {
		log.Fatal("first argument must be a directory")
	}

	pkgInfo, err := sig.Parse(dir, &sig.ParseOptions{
		AllowErrors: *allowErrors,
	})
	if err != nil {
		log.Fatal(err)
	}

	pkgInfo.CommandArgs = os.Args[1:]
	for _, st := range pkgInfo.Structs {
		st.Private = *private
	}

	// Format the output.
	src, err := pkgInfo.Emit()
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		outputName = filepath.Join(dir, "spanner_info_gen.go")
	}
	err = os.WriteFile(outputName, src, 0644)
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
