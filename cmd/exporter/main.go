package main

import (
	"flag"
	"fmt"
	"os"

	app "github.com/trchopan/cardano-postgre-bigquery-export/pkg/application"
	repo "github.com/trchopan/cardano-postgre-bigquery-export/pkg/repository"
)

func main() {
	var configFp string
	flag.StringVar(&configFp, "config", "", "Configuration of tables to be export")

	var pgConfFp string
	flag.StringVar(&pgConfFp, "pgconf", "", "PostgreSql config file")

	var interval int
	flag.IntVar(&interval, "inteval", 3600, "Interval of check for export (seconds).")

	var table string
	flag.StringVar(&table, "table", "", "Interval of check for export (seconds).")

	flag.Parse()

	if pgConfFp == "" || configFp == "" || table == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	pgConf, err := repo.LoadPostgreConf(pgConfFp)
	app.CheckError(err)

	db, err := repo.NewPostgreConnection(pgConf)
	app.CheckError(err)

	config, err := repo.ParseConfiguration(configFp)
	app.CheckError(err)

	gcsClient, err := repo.NewGCSClient()
	app.CheckError(err)

	foundExport, err := repo.FindExportConfig(config.Exports, table)
	app.CheckError(err)

	dbLastId, err := repo.GetLastId(db, foundExport.Table)
	app.CheckError(err)

	fileNames, err := repo.ListGCSFileNames(gcsClient, config, foundExport)
	app.CheckError(err)

	fmt.Println(fileNames)

	storageLastId, err := repo.GetStorageLastIdFromFileNames(fileNames)
	app.CheckError(err)

	err = repo.ExportSQLToCSV(
		config.DbInstance,
		pgConf.Db,
		config,
		foundExport,
		fmt.Sprint(dbLastId),
		storageLastId,
		dbLastId,
	)
	app.CheckError(err)
}
