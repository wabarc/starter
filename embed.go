package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed extensions/**
var extensions embed.FS

const (
	extEntry = "extensions"
	fileMode = 0o755
)

type extension struct {
	// A directory where extensions can be saved for future use in Chrome.
	// If no directory is specified, it defaults to a temporary directory.
	base string
}

func entries() ([]fs.DirEntry, error) {
	return extensions.ReadDir(extEntry)
}

func (e *extension) store() error {
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
				base := strings.TrimPrefix(filepath.Dir(path), extEntry)
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

func (e *extension) lists() (ls []string) {
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

func (e *extension) basedir() string {
	e.base, _ = filepath.Abs(e.base)
	fi, err := os.Lstat(e.base)
	if err == nil && fi.IsDir() {
		return e.base
	}

	if err != nil && os.IsNotExist(err) {
		// If dir not exist, try to create it.
		if err := os.MkdirAll(e.base, fileMode); err != nil {
			e.base, _ = os.MkdirTemp(os.TempDir(), "starter-")
		}
	}

	return e.base
}
