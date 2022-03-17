package main

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

type browser struct {
	userDataDir string
	remotePort  int
	cacheSize   uint64
}

// flags: https://peter.sh/experiments/chromium-command-line-switches/
func (b *browser) flags() []string {
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
		"--user-data-dir=" + b.userDataDir,
		"--disk-cache-size=" + strconv.FormatUint(b.cacheSize, 9),
	}
}

// initial a chrome
func (b *browser) init(ctx context.Context, ext *extension, stopped chan bool) error {
	// Remove exists user data
	_ = os.RemoveAll(b.userDataDir)

	if b.isPortBusy() {
		return fmt.Errorf("Remote debugging port is busy")
	}

	exts := strings.Join(ext.lists(), ",")
	opts := append(b.flags(),
		`--load-extension=`+exts,
		`--disable-extensions-except=`+exts,
		`--remote-debugging-address=0.0.0.0`,
		`--remote-debugging-port=`+strconv.Itoa(b.remotePort),
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

	// TODO: attach to exists window
	if debug {
		opts = append(opts, `chrome://extensions`)
		cmd = exec.CommandContext(ctx, helper.FindChromeExecPath(), opts...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Open extensions page failed: %w", err)
		}
	}

	return nil
}

func (b *browser) exit(cancel context.CancelFunc) {
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

func (b *browser) isPortBusy() bool {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(b.remotePort))
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
