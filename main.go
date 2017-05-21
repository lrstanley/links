package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	HTTP  string `short:"b" long:"http" default:":8080" description:"ip:port pair to bind to"`
	Proxy bool   `short:"p" long:"behind-proxy" description:"if X-Forwarded-For headers should be trusted"`
	TLS   struct {
		Enable bool   `long:"enable" description:"run tls server rather than standard http"`
		Cert   string `short:"c" long:"cert" description:"path to ssl cert file"`
		Key    string `short:"k" long:"key" description:"path to ssl key file"`
	} `group:"TLS Options" namespace:"tls"`
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

	initLogger()

	// Initialize the http/
	httpServer()
}
