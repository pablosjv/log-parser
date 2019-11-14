package main

import (
	"flag"
	"fmt"
	"log-parser/internal/version"

	"log-parser/pkg/parser"

	"os"
)

// Configurable variables
var (
	targetHost   string
	originHost   string
	fileName     string
	initTime     int64
	endTime      int64
	tail         bool
	reportPeriod int64

	printVersion bool
	printHelp    bool
)

func main() {
	flag.StringVar(&targetHost, "targetHost", "", "Hostname to wich other connects")
	flag.StringVar(&originHost, "originHost", "", "Hostname that connect to others")
	flag.StringVar(&fileName, "f", "", "Log filename to parse")
	flag.Int64Var(&initTime, "i", 0, "Init time in unix format")
	flag.Int64Var(&endTime, "e", 0, "End time in unix format")
	flag.BoolVar(&tail, "t", false, "Tail the log file and report")
	flag.Int64Var(&reportPeriod, "period", 3600, "Period in seconds to report")

	flag.BoolVar(&printVersion, "v", false, "Print the version and exits")
	flag.BoolVar(&printHelp, "h", false, "Print this help message")

	flag.Parse()

	if printVersion {
		fmt.Println("Version", version.Version)
		os.Exit(0)
	}
	if printHelp {
		flag.Usage()
		os.Exit(0)
	}

	app := parser.NewApp(
		initTime, endTime, targetHost, originHost, fileName, tail, reportPeriod)
	app.Run()
}
