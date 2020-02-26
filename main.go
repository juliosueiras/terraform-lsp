package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	oldLog "log"
	"os"
	"strings"

	"github.com/juliosueiras/terraform-lsp/langserver"
)

var tcp = flag.Bool("tcp", false, "Use TCP instead of Stdio(which is default)")
var port = flag.Int("port", 9900, "Port for TCP Server")
var address = flag.String("address", "127.0.0.1", "Address for TCP Server")

var location = flag.String("log-location", "", "Location of the lsp log")
var locationJRPC2 = flag.String("log-jrpc2-location", "", "Location of the lsp log for jrpc2")

var debug = flag.Bool("debug", false, "Enable debug output")
var debugJRPC2 = flag.Bool("debug-jrpc2", false, "Enable debug output for jrpc2")

var enableLogFile = flag.Bool("enable-log-file", false, "Enable log file")
var enableLogFileJRPC2 = flag.Bool("enable-log-jrpc2-file", false, "Enable log file for JRPC2")

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

  var oldLogInstance *oldLog.Logger

  if *debugJRPC2 {
    if !*tcp && !*enableLogFileJRPC2 {
      log.Fatal("Debug for JRPC2 has to be set for log file location if is set to use stdio")
    }

    oldLogInstance = oldLog.New(os.Stdout, "", 0)
    if *enableLogFileJRPC2 {
      f, err := os.OpenFile(fmt.Sprintf("%stf-lsp-jrpc2.log", *locationJRPC2), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
      if err != nil {
        log.Fatalf("error opening file: %v", err)
      }
      defer f.Close()
      oldLogInstance.SetOutput(f)
    }
	}

  langserver.InitializeServiceMap()

	if *tcp {
		langserver.RunTCPServer(*address, *port, oldLogInstance)
	} else {
		langserver.RunStdioServer(oldLogInstance)
	}
}
