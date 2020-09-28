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

	app.Init("news-aggregator-migrate")
	defer app.Destroy()
	log := app.GetLog().Sugar()

	oldVersion, newVersion, err := migrations.Run(app.GetDB(), flag.Args()...)
	if err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	if newVersion != oldVersion {
		log.Infof("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		log.Infof("version is %d (not changed)", oldVersion)
	}
}
