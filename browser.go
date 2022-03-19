package starter

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wabarc/helper"
)

type Browser struct {
	UserDataDir string
	RemotePort  int
	CacheSize   uint64
}

// flags: https://peter.sh/experiments/chromium-command-line-switches/
func (b *Browser) flags() []string {
	return []string{
		"--lang=en-US",
		"--disable-gpu",
		"--no-sandbox",
		// "--headless",
		"--disable-background-networking",
		"--enable-features=NetworkService,NetworkServiceInProcess",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-features=site-per-process,Translate,BlinkGenPropertyTrees",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-infobars",
		"--disable-notifications",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--no-first-run",
		`--no-startup-window`,
		"--no-default-browser-check",
		"--force-color-profile=srgb",
		"--metrics-recording-only",
		"--safebrowsing-disable-auto-update",
		"--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
		"--use-fake-device-for-media-stream",
		"--ignore-certificate-errors",
		"--disable-extensions=false",
		"--window-size=1280,1024",
		"--user-data-dir=" + b.UserDataDir,
		"--disk-cache-size=" + strconv.FormatUint(b.CacheSize, 9),
	}
}

// initial a chrome
func (b *Browser) Init(ctx context.Context, ext *Extension, stopped chan bool) error {
	// Remove exists user data
	_ = os.RemoveAll(b.UserDataDir)

	if b.isPortBusy() {
		return fmt.Errorf("Remote debugging port is busy")
	}

	exts := strings.Join(ext.lists(), ",")
	opts := append(b.flags(),
		`--load-extension=`+exts,
		`--disable-extensions-except=`+exts,
		`--remote-debugging-address=0.0.0.0`,
		`--remote-debugging-port=`+strconv.Itoa(b.RemotePort),
		`chrome://extensions`,
	)

	cmd := exec.CommandContext(ctx, helper.FindChromeExecPath(), opts...)
	cmd.Env = append(os.Environ(), cmd.Env...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Run StdoutPipe failed: %w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Start failed: %w", err)
	}
	go func() {
		readOutput(out)
		_ = cmd.Wait()
		stopped <- true
	}()

	return nil
}

func (b *Browser) Exit(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	var once sync.Once
	for {
		sig := <-signalChan
		if sig == os.Interrupt {
			once.Do(func() {
				cancel()
			})
			return
		}
	}
}

func (b *Browser) isPortBusy() bool {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(b.RemotePort))
	conn, _ := net.DialTimeout("tcp", addr, time.Second)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func readOutput(rc io.ReadCloser) {
	for {
		out := make([]byte, 1024)
		_, err := rc.Read(out)
		fmt.Print(string(out))
		if err != nil {
			break
		}
	}
}
