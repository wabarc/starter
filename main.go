package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path"

	"github.com/dustin/go-humanize"
)

var (
	workspace  string
	cacheSize  string
	remotePort int
	desktop    bool
	debug      bool
)

func init() {
	flag.StringVar(&workspace, "workspace", "", "A directory for storing Chrome preferences and extension sources")
	flag.StringVar(&cacheSize, "cache-size", "512MB", "Forces the Chrome maximum disk space to be used by the disk cache, in humanize sizes")
	flag.BoolVar(&debug, "debug", false, "If debug mode is enabled, installed extensions will be captured as an image")
	flag.BoolVar(&desktop, "desktop", false, "If you're working on a desktop, Xvfb isn't required")
	flag.IntVar(&remotePort, "remote-debugging-port", 9222, "Remote debugging port for Chrome")
	flag.Parse()

	if workspace == "" {
		flag.Usage()
		os.Exit(0)
	}

	// Environment for Xvfb
	if os.Getenv("DISPLAY") == "" {
		os.Setenv("DISPLAY", ":99")
	}
	if os.Getenv("XVFB_WHD") == "" {
		os.Setenv("XVFB_WHD", "1280x1024x16")
	}
}

// TODO: handle gracefully
func main() {
	ext := &extension{
		base: path.Join(workspace, extEntry),
	}
	if err := ext.store(); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if !desktop {
		if err := runXvfb(ctx); err != nil {
			log.Fatalln(err)
		}
	}

	cacheSize, err := humanize.ParseBytes(cacheSize)
	if err != nil {
		log.Fatalln(err)
	}

	browser := &browser{
		userDataDir: path.Join(workspace, "UserDataDir"),
		remotePort:  remotePort,
		cacheSize:   cacheSize,
	}
	defer os.RemoveAll(browser.userDataDir)
	go browser.exit(cancel)

	stopped := make(chan bool, 1)
	if err := browser.init(ctx, ext, stopped); err != nil {
		log.Fatalln(err)
	}
	<-stopped

	// exit with unexpected
	log.Println("Browser exit unexpected")
	os.Exit(1)
}
