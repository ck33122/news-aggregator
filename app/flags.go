package app

import (
	"flag"
)

type Flags struct {
	ConfigPath string
	ShowHelp   bool
}

func NewFlags() *Flags {
	flags := &Flags{}
	flag.BoolVar(&flags.ShowHelp, "help", false, "show help message and exit")
	flag.StringVar(&flags.ConfigPath, "config", "app.yaml", "path to configuration file")
	flag.Parse()
	return flags
}
