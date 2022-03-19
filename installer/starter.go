package installer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

type Starter struct {
	Home string

	bin string // Installed binary file
}

// Install the compiled binary from GitHub.
func (s *Starter) Install() error {
	if goos := runtime.GOOS; goos != "linux" {
		return fmt.Errorf("Only Linux is supported, your OS is: %s", goos)
	}

	endpoint := "https://github.com/wabarc/starter/raw/main/install.sh"
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(endpoint)
	if err != nil {
		return fmt.Errorf("Download install script failed: %w", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	script := path.Join(s.Home, "install.sh")
	if err := os.WriteFile(script, buf, 0o644); err != nil {
		return fmt.Errorf("Write script file failed: %w", err)
	}

	sh := which("sh")
	if sh == "" {
		return fmt.Errorf("POSIX shell not found")
	}
	cmd := exec.Command(sh, script)
	cmd.Dir = s.Home

	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Install starter failed: %w", err)
	}

	// Set binary file path for callsite
	// path is {s.Home}/bin/starter
	s.bin = path.Join(s.Home, "bin", "starter")

	return nil
}

// Command returns the binary path of the starter, which is an executable file.
func (s *Starter) Command() string {
	return s.bin
}

func which(command string) string {
	found, err := exec.LookPath(command)
	if err != nil {
		return ""
	}
	return found
}
