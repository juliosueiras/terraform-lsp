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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Version string
var GitCommit string
var Date string

func init() {
	flag.Bool("tcp", false, "Use TCP instead of Stdio(which is default)")
	flag.Int("port", 9900, "Port for TCP Server")
	flag.String("address", "127.0.0.1", "Address for TCP Server")

	flag.String("log-location", "", "Location of the lsp log")
	flag.String("log-jrpc2-location", "", "Location of the lsp log for jrpc2")

	flag.Bool("debug", false, "Enable debug output")
	flag.Bool("debug-jrpc2", false, "Enable debug output for jrpc2")

	flag.Bool("enable-log-file", false, "Enable log file")
	flag.Bool("enable-log-jrpc2-file", false, "Enable log file for JRPC2")

	flag.Bool("version", false, "Show version")

	// Load config from file
	configViper()
}

func main() {
	oldLog.SetOutput(ioutil.Discard)
	oldLog.SetFlags(0)

	version := viper.GetBool("version")
	if version {
		fmt.Printf("v%s, commit: %s, build on: %s\n", strings.Trim(Version, "v"), GitCommit, Date)
		return
	}

	debug := viper.GetBool("debug")
	log.Infof("Log Level is Debug: %t", debug)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	enableLogFile := viper.GetBool("enable-log-file")
	if enableLogFile {
		location := viper.GetString("log-location")
		f, err := os.OpenFile(fmt.Sprintf("%stf-lsp.log", location), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	var oldLogInstance *oldLog.Logger

	debugJRPC2 := viper.GetBool("debug-jrpc2")
	tcp := viper.GetBool("tcp")
	enableLogFileJRPC2 := viper.GetBool("enable-log-jrpc2-file")

	if debugJRPC2 {
		if !tcp && !enableLogFileJRPC2 {
			log.Fatal("Debug for JRPC2 has to be set for log file location if is set to use stdio")
		}

		oldLogInstance = oldLog.New(os.Stdout, "", 0)
		if enableLogFileJRPC2 {
			locationJRPC2 := viper.GetString("log-jrpc2-location")
			f, err := os.OpenFile(fmt.Sprintf("%stf-lsp-jrpc2.log", locationJRPC2), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
			oldLogInstance.SetOutput(f)
		}
	}

	langserver.InitializeServiceMap()

	if tcp {
		address := viper.GetString("address")
		port := viper.GetInt("port")
		langserver.RunTCPServer(address, port, oldLogInstance)
	} else {
		langserver.RunStdioServer(oldLogInstance)
	}
}

func configViper() {
	// Accept CLI arguments
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
