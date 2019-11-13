package main

import (
	"log"
	"strings"

	"github.com/nlnwa/veidemann-reset/pkg/file"
	"github.com/nlnwa/veidemann-reset/pkg/rethinkdb"
	"github.com/nlnwa/veidemann-reset/pkg/version"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// configuration defaults

	rethinkdbUser := "admin"
	rethinkdbPassword := "rethinkdb"
	rethinkdbName := "veidemann"
	rethinkdbHost := "localhost"
	rethinkdbPort := 28015

	directories := []string{
		"/warcs",
		"/validwarcs",
		"/invalidwarcs",
		"/delivery",
		"/backup/oos",
	}

	tables := map[string][]string{
		"veidemann": {
			"crawl_host_group",
			"crawl_log",
			"crawled_content",
			"events",
			"executions",
			"extracted_text",
			"job_executions",
			"locks",
			"page_log",
			"storage_ref",
			"uri_queue",
		},
		"report": {
			"invalid_warcs",
			"valid_warcs",
		},
	}

	flag.StringVar(&rethinkdbHost, "rethinkdb-host", rethinkdbHost, "RethinkDb hostname")
	flag.IntVar(&rethinkdbPort, "rethinkdb-port", rethinkdbPort, "RethinkDb port")
	flag.StringVar(&rethinkdbName, "rethinkdb-name", rethinkdbName, "RethinkDb database name")
	flag.StringVar(&rethinkdbUser, "rethinkdb-user", rethinkdbUser, "RethinkDb user")
	flag.StringVar(&rethinkdbPassword, "rethinkdb-password", rethinkdbPassword, "RethinkDb password")

	flag.StringSliceVar(&directories, "directories", directories, "directories to clean")

	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rethinkdbUser = viper.GetString("rethinkdb-user")
	rethinkdbPassword = viper.GetString("rethinkdb-password")
	rethinkdbName = viper.GetString("rethinkdb-name")
	rethinkdbHost = viper.GetString("rethinkdb-host")
	rethinkdbPort = viper.GetInt("rethinkdb-port")

	rethinkdbOptions := rethinkdb.Options{
		Name:     rethinkdbName,
		Host:     rethinkdbHost,
		Port:     rethinkdbPort,
		Username: rethinkdbUser,
		Password: rethinkdbPassword,
	}

	directories = viper.GetStringSlice("directories")

	log.Printf("Veidemann reset, version: %s\n", version.Version)

	file.NewDirectoryCleaner(directories).RemoveFiles()

	err = rethinkdb.NewClient(rethinkdbOptions).Clean(tables)
	if err != nil {
		log.Fatalf("failed to reset database: %s\n", err.Error())
	} else {
		log.Printf("Veidemann reset completed successfully")
	}
}
