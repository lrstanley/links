// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

type Config struct {
	Site       string `env:"SITE_URL" short:"s" long:"site-name" default:"https://links.wtf" description:"site url, used for url generation"`
	SessionDir string `env:"SESSION_DIR" long:"session-dir" description:"optional location to store temporary sessions"`
	Quiet      bool   `env:"QUIET" short:"q" long:"quiet" description:"don't log to stdout"`
	Debug      bool   `env:"DEBUG" long:"debug" description:"enable debugging (pprof endpoints)"`
	HTTP       string `env:"HTTP" short:"b" long:"http" default:":8080" description:"ip:port pair to bind to"`
	Proxy      bool   `env:"PROXY" short:"p" long:"behind-proxy" description:"if X-Forwarded-For headers should be trusted"`
	TLS        struct {
		Enable bool   `env:"TLS_ENABLE" long:"enable" description:"run tls server rather than standard http"`
		Cert   string `env:"TLS_CERT_PATH" short:"c" long:"cert" description:"path to ssl cert file"`
		Key    string `env:"TLS_KEY_PATH" short:"k" long:"key" description:"path to ssl key file"`
	} `group:"TLS Options" namespace:"tls"`
	DBPath          string `env:"DB_PATH" short:"d" long:"db" default:"store.db" description:"path to database file"`
	KeyLength       int64  `env:"KEY_LENGTH" long:"key-length" default:"4" description:"default length of key (uuid) for generated urls"`
	HTTPPreInclude  string `env:"HTTP_PRE_INCLUDE" long:"http-pre-include" description:"HTTP include which is included directly after css is included (near top of the page)"`
	HTTPPostInclude string `env:"HTTP_POST_INCLUDE" long:"http-post-include" description:"HTTP include which is included directly after js is included (near bottom of the page)"`

	ExportFile string `short:"e" long:"export-file" default:"links.export" description:"file to export db to"`
	ExportJSON bool   `long:"export-json" description:"export db to json elements"`

	SafeBrowsing struct {
		APIKey           string        `env:"SAFEBROWSING_API_KEY" long:"api-key" description:"Google API Key used for querying SafeBrowsing, disabled if not provided (see: https://github.com/lrstanley/links#google-safebrowsing)"`
		DBPath           string        `env:"SAFEBROWSING_DB_PATH" long:"db" default:"safebrowsing.db" description:"path to SafeBrowsing database file"`
		UpdatePeriod     time.Duration `env:"SAFEBROWSING_UPDATE_PERIOD" long:"update-period" default:"1h" description:"duration between updates to the SafeBrowsing API local database"`
		RedirectFallback bool          `env:"SAFEBROWSING_REDIRECT_FALLBACK" long:"redirect-fallback" description:"if the SafeBrowsing request fails (local cache, and remote hit), this still lets the redirect happen"`
	} `group:"Safe Browsing Support" namespace:"safebrowsing"`

	Prometheus struct {
		Enabled  bool   `env:"PROM_ENABLED" long:"enabled" description:"enable exposing of prometheus metrics (on std port, or --prometheus.addr)"`
		Addr     string `env:"PROM_ADDR" long:"addr" description:"expose on custom address/port, e.g. ':9001' (all ips) or 'localhost:9001' (local only)"`
		Endpoint string `env:"PROM_ENDPOINT" long:"endpoint" default:"/metrics" description:"endpoint to expose metrics on"`
	} `group:"Prometheus Metrics" namespace:"prom"`

	VersionFlag bool `short:"v" long:"version" description:"display the version of links.wtf and exit"`

	CommandAdd    CommandAdd    `command:"add" description:"add a link"`
	CommandDelete CommandDelete `command:"delete" description:"delete a link, id, or link matching an author"`
}

var conf Config

var debug *log.Logger

func initLogger() {
	debug = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	debug.Print("initialized logger")
}

func main() {
	parser := flags.NewParser(&conf, flags.HelpFlag)
	parser.SubcommandsOptional = true
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if parser.Active != nil {
		os.Exit(0)
	}

	if conf.VersionFlag {
		fmt.Printf("links version: %s [%s] (%s, %s), compiled %s\n", version, commit, runtime.GOOS, runtime.GOARCH, date)
		os.Exit(0)
	}

	// Do some configuration validation.
	if conf.HTTP == "" {
		debug.Fatalf("invalid http flag supplied: %s", conf.HTTP)
	}
	if conf.KeyLength < 4 {
		conf.KeyLength = 4
	}

	conf.Site = strings.TrimRight(conf.Site, "/")

	initLogger()

	// Verify db is accessible.
	verifyDB()

	if conf.ExportJSON {
		dbExportJSON(conf.ExportFile)
		debug.Print("export complete")
		return
	}

	// Google SafeBrowsing.
	initSafeBrowsing()

	// Setup methods to allow signaling to all children methods that we're stopping.
	ctx, closer := context.WithCancel(context.Background())
	errors := make(chan error)
	wg := &sync.WaitGroup{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Initialize the http/https server.
	go httpServer(ctx, wg, errors)

	fmt.Println("listening for signal. CTRL+C to quit.")

	go func() {
		for {
			select {
			case <-signals:
				fmt.Println("\nsignal received, shutting down")
			case <-errors:
				debug.Println(err)
			}

			// Signal to exit.
			closer()
		}
	}()

	// Wait for the context to close, and wait for all goroutines/processes to exit.
	<-ctx.Done()
	wg.Wait()

	if safeBrowser != nil {
		if err = safeBrowser.Close(); err != nil {
			debug.Fatalf("error closing google safebrowsing: %v", err)
		}
	}

	debug.Println("shutdown complete")
}
