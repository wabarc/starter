package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	workspace string
	debug     bool
)

func init() {
	flag.StringVar(&workspace, "workspace", "", "A directory for storing Chrome preferences and extension sources")
	flag.BoolVar(&debug, "debug", false, "If debug mode is enabled, installed extensions will be captured as an image")
	flag.Parse()

	if workspace == "" {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	ext := &extension{
		base: path.Join(workspace, extEntry),
	}
	if err := ext.store(); err != nil {
		log.Fatalln(err)
	}

	if err := run(ext); err != nil {
		log.Fatalln(err)
	}
}

func run(ext *extension) error {
	// create chrome instance
	exts := strings.Join(ext.lists(), ",")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.UserDataDir(path.Join(workspace, "UserDataDir")),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("flag-switches-begin", true),
		chromedp.Flag("flag-switches-end", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		// chromedp.Flag("auto-open-devtools-for-tabs", true),
		chromedp.Flag("disable-notifications", false),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"),
		// chromedp.UserAgent("APIs-Google (+https://developers.google.com/webmasters/APIs-Google.html)"),
		chromedp.WindowSize(1280, 1024),
		chromedp.Flag("incognito", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("load-extension", exts),
		chromedp.Flag("disable-extensions-except", exts),
		chromedp.Flag("remote-debugging-address", "0.0.0.0"),
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.Flag("lang", "en-US"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	if debug {
		log.Println("path of extensions:", exts)

		if err := chromedp.Run(ctx, chromedp.Navigate(`chrome://extensions`)); err != nil {
			return err
		}

		var buf []byte
		if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100)); err != nil {
			return err
		}

		if err := ioutil.WriteFile("installed-extensions.png", buf, fileMode); err != nil {
			return err
		}
	}

	return nil
}
