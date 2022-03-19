package installer

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInstall(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "starter")
	if err != nil {
		t.Fatalf("Unexpected create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	s := &Starter{
		Home: dir,
	}
	err = s.Install()
	if err != nil {
		t.Fatalf("Unexpected install starter: %v", err)
	}

	installed := s.Command()
	binpath := path.Join(dir, "bin", "starter")
	if installed != binpath {
		t.Fatalf("Unexpected binary file, got %s instead of %s", installed, binpath)
	}

	fi, err := os.Stat(installed)
	if err != nil {
		t.Fatalf("Unexpected find starter command: %v", err)
	}

	if !fi.Mode().IsRegular() {
		t.Fatalf("Unpexected binary file")
	}
}
