package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var flagDir    string

func init() {
	flag.StringVar(&flagDir, "dir", defaultDir(), "target directory")
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	fset := token.NewFileSet()
	if err := filepath.Walk(flagDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "_test.go") {
			astFile, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}

			for _, d := range astFile.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if strings.HasPrefix(fn.Name.Name, "Test") {
						fmt.Printf("%s in file: %s\n", fn.Name.Name, path)
					}
				}
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to Walk file: %w", err)
	}

	return nil
}


func defaultDir() string {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return root
}
