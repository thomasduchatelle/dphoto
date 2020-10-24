package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

var debug bool
var trace bool

func init() {
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.BoolVar(&debug, "v", false, "enable debug output")
	flag.BoolVar(&trace, "trace", false, "enable trace output")
	flag.BoolVar(&trace, "vv", false, "enable trace output")
}

func loggingInit() {
	log.SetOutput(os.Stdout)
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)

	if trace {
		log.SetLevel(log.TraceLevel)
	} else if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.WithFields(log.Fields{
		"debug": debug,
		"trace": trace,
	}).Traceln("Logging initiated")
}
