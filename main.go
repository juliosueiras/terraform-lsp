package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
  oldLog "log"
  log "github.com/sirupsen/logrus"

  "io/ioutil"
	"github.com/creachadair/jrpc2/channel"
	"github.com/juliosueiras/terraform-lsp/langserver"
)

var location = flag.String("log-location", "", "Location of the lsp log")
var debug = flag.Bool("debug", false, "Enable debug output")
var enableLogFile = flag.Bool("enable-log-file", false, "Enable log file")

var Version string
var GitCommit string
var Date string

var version = flag.Bool("version", false, "Show version")

func main() {
	flag.Parse()

  oldLog.SetOutput(ioutil.Discard)
  oldLog.SetFlags(0)


	if *version {
		fmt.Printf("v%s, commit: %s, build on: %s", strings.Trim(Version, "v"), GitCommit, Date)
		return
	}

	Server := langserver.CreateServer()

  log.Infof("Log Level is Debug: %t", *debug)

  if *debug {
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }

	if *enableLogFile {
		f, err := os.OpenFile(fmt.Sprintf("%stf-lsp.log", *location), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// Start the server on a channel comprising stdin/stdout.
	Server.Start(channel.Header("")(os.Stdin, os.Stdout))
	log.Info("Server started")

	// Wait for the server to exit, and report any errors.
	if err := Server.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}

	log.Info("Server Finish")
}
