//go:build go1.18
// +build go1.18

package starter

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed extensions/* all:extensions/**/_*
var extensions embed.FS

const (
	ExtEntry = "extensions" // Entry directory of extensions
	fileMode = 0o755
)

type Extension struct {
	// A directory where extensions can be saved for future use in Chrome.
	// If no directory is specified, it defaults to a temporary directory.
	Base string
}

func entries() ([]fs.DirEntry, error) {
	return extensions.ReadDir(ExtEntry)
}

func (e *Extension) Store() error {
	entries, err := entries()
	if err != nil {
		return err
	}

	dest := e.basedir()
	for _, entry := range entries {
		if info, _ := entry.Info(); info.IsDir() {
			err = fs.WalkDir(extensions, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				base := strings.TrimPrefix(filepath.Dir(path), ExtEntry)
				destDir := filepath.Join(dest, base)
				if err := os.MkdirAll(destDir, fileMode); err != nil {
					return err
				}

				if !d.IsDir() {
					buf, err := fs.ReadFile(extensions, path)
					if err != nil {
						return err
					}
					destName := filepath.Join(destDir, filepath.Base(path))
					_ = os.WriteFile(destName, buf, fileMode)
				}

				return nil
			})
		}
	}
	return err
}

func (e *Extension) lists() (ls []string) {
	entries, err := entries()
	if err != nil {
		return
	}

	base := e.basedir()
	for _, entry := range entries {
		if info, _ := entry.Info(); info.IsDir() {
			dir := filepath.Join(base, entry.Name())
			ls = append(ls, dir)
		}
	}

	return
}

func (e *Extension) basedir() string {
	e.Base, _ = filepath.Abs(e.Base)
	fi, err := os.Lstat(e.Base)
	if err == nil && fi.IsDir() {
		return e.Base
	}

	if err != nil && os.IsNotExist(err) {
		// If dir not exist, try to create it.
		if err := os.MkdirAll(e.Base, fileMode); err != nil {
			e.Base, _ = os.MkdirTemp(os.TempDir(), "starter-")
		}
	}

	return e.Base
}
