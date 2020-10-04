package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [--debug|--trace] [--help]\n\n", os.Args[0])
		fmt.Println("Options:")
		args := [][]string{
			{"--help -h", "print this help and exit"},
			{"--debug -v", "enable debug logging"},
			{"--trace -vv", "enable trace logging"},
		}
		for _, arg := range args {
			fmt.Printf("    %-20s %s\n", arg[0], arg[1])
		}

		hidden := map[string]int{
			"debug": 0,
			"v":     0,
			"trace": 0,
			"vv":    0,
		}
		flag.VisitAll(func(f *flag.Flag) {
			if _, ok := hidden[f.Name]; !ok {
				fmt.Printf("    --%-18s %s\n", f.Name, f.Usage)
			}
		})
	}

	loggingArgs()

	flag.Parse()
	loggingInit()

	log.Infoln("Starting dphoto - delegates")
	startUDiskListener()
}
