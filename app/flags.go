package app

import (
	"flag"
	"os"
)

var (
	configPath string
	showHelp   bool
)

func parseFlags() {
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(1)
	}
}

func init() {
	flag.BoolVar(&showHelp, "help", false, "show help message and exit")
	flag.StringVar(&configPath, "config", "app.ini", "path to configuration file")
}
