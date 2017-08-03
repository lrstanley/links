// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

type Config struct {
	Site  string `short:"s" long:"site-name" default:"https://links.ml" description:"site url, used for url generation"`
	Quiet bool   `short:"q" long:"quiet" description:"don't log to stdout"`
	Debug bool   `long:"debug" description:"enable debugging (pprof endpoints)"`
	HTTP  string `short:"b" long:"http" default:":8080" description:"ip:port pair to bind to"`
	Proxy bool   `short:"p" long:"behind-proxy" description:"if X-Forwarded-For headers should be trusted"`
	TLS   struct {
		Enable bool   `long:"enable" description:"run tls server rather than standard http"`
		Cert   string `short:"c" long:"cert" description:"path to ssl cert file"`
		Key    string `short:"k" long:"key" description:"path to ssl key file"`
	} `group:"TLS Options" namespace:"tls"`
	DBPath    string `short:"d" long:"db" default:"store.db" description:"path to database file"`
	KeyLength int    `long:"key-length" default:"4" description:"default length of key (uuid) for generated urls"`

	ExportFile string `short:"e" long:"export-file" default:"links.export" description:"file to export db to"`
	ExportJSON bool   `long:"export-json" description:"export db to json elements"`

	MigrateFlag bool   `long:"migrate" description:"begin migration from links.ml running MySQL"`
	MigrateInfo string `long:"migrate-info" default:"user:passwd@tcp(localhost:3306)/links_db" description:"connection url used to connect to the old mysql instance"`

	VersionFlag bool `short:"v" long:"version" description:"display the version of links.ml and exit"`
}

var conf Config

var debug *log.Logger

func initLogger() {
	debug = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	debug.Print("initialized logger")
}

func main() {
	_, err := flags.Parse(&conf)
	if err != nil {
		if FlagErr, ok := err.(*flags.Error); ok && FlagErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		// go-flags should print to stderr/stdout as necessary, so we won't.
		os.Exit(1)
	}

	if conf.VersionFlag {
		fmt.Printf("links.ml version: %s [%s] (%s, %s), compiled %s\n", version, commit, runtime.GOOS, runtime.GOARCH, date)
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

	if conf.MigrateFlag {
		migrateDB(conf.MigrateInfo)
		return
	}

	if conf.ExportJSON {
		fmt.Printf("%s", hash("test"))
		dbExportJSON(conf.ExportFile)
		debug.Print("export complete")
		return
	}

	// Initialize the http/https server.
	httpServer()
}
