package sweep_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const servicesImportPrefix = "github.com/G-Core/terraform-provider-gcore/internal/services/"

// TestSweeperRegistrationsAreComplete guards the one hole that has already bitten
// this repo (a35c763b).
//
// A sweeper registers from an init() in its own service package, so a package
// this one does not blank-import is never registered and can never run -- and
// the sweep still exits 0, leaving the resources behind. Nothing in the build
// or the type system catches that.
func TestSweeperRegistrationsAreComplete(t *testing.T) {
	t.Parallel()

	onDisk := packagesWithSweepFile(t)
	imported := blankImportedServicePackages(t)

	if missing := difference(onDisk, imported); len(missing) > 0 {
		t.Errorf("these packages have a sweep.go but are not blank-imported here, so their sweepers can never run: %v\n"+
			"add `_ %q` to sweep_test.go", missing, servicesImportPrefix+missing[0])
	}
	if extra := difference(imported, onDisk); len(extra) > 0 {
		t.Errorf("these packages are blank-imported here but have no sweep.go (stale import after a rename?): %v", extra)
	}
}

func packagesWithSweepFile(t *testing.T) map[string]bool {
	t.Helper()

	entries, err := os.ReadDir("../services")
	if err != nil {
		t.Fatalf("reading ../services: %s", err)
	}

	out := map[string]bool{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join("../services", e.Name(), "sweep.go")); err == nil {
			out[e.Name()] = true
		}
	}

	// A comparison that silently finds nothing to compare is itself a false
	// green, which is the exact class of bug this test exists to catch.
	if len(out) == 0 {
		t.Fatal("found no sweep.go files at all; the directory walk is broken, not the repo")
	}
	return out
}

func blankImportedServicePackages(t *testing.T) map[string]bool {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("globbing this package: %s", err)
	}

	out := map[string]bool{}
	for _, f := range files {
		parsed, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parsing %s: %s", f, err)
		}
		for _, imp := range parsed.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			if strings.HasPrefix(path, servicesImportPrefix) {
				out[strings.TrimPrefix(path, servicesImportPrefix)] = true
			}
		}
	}
	return out
}

func difference(a, b map[string]bool) []string {
	var out []string
	for k := range a {
		if !b[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
