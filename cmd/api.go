package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ck33122/news-aggregator/api"
	"github.com/ck33122/news-aggregator/app"
)

func main() {
	flags := app.NewFlags()
	if flags.ShowHelp {
		flag.Usage()
		return
	}
	config, err := app.NewConfig(flags.ConfigPath, "news-aggregator-api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating config object: %v", err)
		os.Exit(1)
	}
	log, err := app.NewLog(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating log object: %v", err)
		os.Exit(1)
	}
	db, err := app.NewDB(log, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating database object: %v", err)
		os.Exit(1)
	}
	actions, err := app.NewActions(log, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating actions object: %v", err)
		os.Exit(1)
	}
	server := app.NewServer(log, config)
	api.Setup(server, actions)
	server.Run()
}
