package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ck33122/news-aggregator/app"
	"github.com/go-pg/migrations/v8"

	// this import runs all init() funcs in list module, needs to register migrations
	_ "github.com/ck33122/news-aggregator/migrations"
)

const usageText = `Migration tool. Supported commands are:
  init                   - creates version info table in the database
  up                     - runs all available migrations
  up [target]            - runs available migrations up to the target one
  down                   - reverts last migration
  reset                  - reverts all migrations
  version                - prints current db version
  set_version [version]  - sets db version without running migrations
Usage:
  go run main.go <command> [args] [-config <path_to_config_file>]
`

func main() {
	flag.Usage = func() {
		fmt.Print(usageText)
		os.Exit(2)
	}
	flags := app.NewFlags()
	if flags.ShowHelp {
		flag.Usage()
		return
	}
	config, err := app.NewConfig(flags.ConfigPath, "news-aggregator-migrate")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating config object: %v", err)
		os.Exit(1)
	}
	config.Logger.Dir = "" // write log files on migration isn't very good idea
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
	defer db.Close()

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		log.Sugar().Fatalf("error running migrations: %v", err)
	}

	if newVersion != oldVersion {
		log.Sugar().Infof("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		log.Sugar().Infof("version is %d (not changed)", oldVersion)
	}
}
