package main

import (
	"flag"
	"os"

	app "github.com/trchopan/cardano-postgre-bigquery-export/pkg/application"
	repo "github.com/trchopan/cardano-postgre-bigquery-export/pkg/repository"
)

func main() {
	var configFp string
	flag.StringVar(&configFp, "config", "", "Configuration of tables to be export")

	var table string
	flag.StringVar(&table, "table", "", "Interval of check for export (seconds).")

	flag.Parse()

	if configFp == "" || table == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	config, err := repo.ParseConfiguration(configFp)
	app.CheckError(err)

    foundExport, err := repo.FindExportConfig(config.Exports, table)
    app.CheckError(err)

	repo.BQTransferJobCreate(config, foundExport)
}
