package kimono

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

// Tidy runs `go get -u` and `go mod tidy` on all supported Go
// modules in the current git repository.
func Tidy() error {
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return err
	}
	return filepath.WalkDir(filepath.Dir(root), sanitizeWalkDirFn)
}

func sanitizeWalkDirFn(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return nil
	}
	if d.Name() == ".git" || d.Name() == "vendor" {
		return filepath.SkipDir
	}
	if !futil.Exists(filepath.Join(path, "go.mod")) {
		return filepath.SkipDir
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	if !hasDependencies() {
		return filepath.SkipDir
	}
	fmt.Printf("\n%s:\n", path)
	_ = update()
	_ = tidy()
	return nil
}

func hasDependencies() bool {
	out := run.Out(`go`, `list`, `-m`, `all`)
	return len(strings.Split(out, "\n")) > 1
}

func update() error {
	return run.Exec("go", "get", "-u")
}

func tidy() error {
	return run.Exec("go", "mod", "tidy")
}

func ListDependents() ([]string, error) {
	modName := run.Out("go", "list", "-m")
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(
		filepath.Dir(root),
		findDependentsWalkDirFn(modName),
	)
	return nil, nil
}

func ListDependencies() ([]string, error) {
	out, err := exec.Command(`go`, `list`, `-m`, `all`).Output()
	if err != nil {
		return nil, err
	}
	return to.Lines(out), nil
}

func findDependentsWalkDirFn(modPath string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".git" || d.Name() == "vendor" {
			return filepath.SkipDir
		}
		if !futil.Exists(filepath.Join(path, "go.mod")) {
			return nil
		}
		if err := os.Chdir(path); err != nil {
			return err
		}
		if !hasDependencies() {
			return filepath.SkipDir
		}

		return nil
	}
}
