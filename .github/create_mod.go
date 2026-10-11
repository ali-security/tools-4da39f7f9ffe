//go:build ignore
// +build ignore

// create_mod builds the module zip for golang.org/x/tools exactly as the Go
// module proxy lays it out, using golang.org/x/mod/zip.CreateFromDir.
//
// It is built by .github/workflows/seal-build.yml in a throwaway module that
// requires golang.org/x/mod; the ignore build tag (and the .github dot-dir)
// keep it out of `go build ./...` / `go vet ./...` of golang.org/x/tools.
//
// Usage: create_mod <module-path> <version> <source-dir> <output-zip>
package main

import (
	"log"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: create_mod <module-path> <version> <source-dir> <output-zip>")
	}
	m := module.Version{Path: os.Args[1], Version: os.Args[2]}
	f, err := os.Create(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	if err := zip.CreateFromDir(f, m, os.Args[3]); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("created module zip %s for %s@%s", os.Args[4], m.Path, m.Version)
}
