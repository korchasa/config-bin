package main

import (
	"configBin/pkg/encryptor/aes"
	"configBin/pkg/metrics/prometheus"
	"configBin/pkg/server"
	"configBin/pkg/server/responder"
	"configBin/pkg/server/templates"
	"configBin/pkg/storage/sqlite"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(
		&log.TextFormatter{
			ForceColors: true,
		},
	)

	listen := ensureEnv("LISTEN")
	sqlitePath := ensureEnv("SQLITE_PATH")

	enc := aes.NewEncryptor()

	store, err := sqlite.NewSqliteStorage(sqlitePath, enc)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create sqlite storage: %w", err))
	}
	err = store.ApplySchema()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to apply db schema: %w", err))
	}

	tplProvider, err := templates.Build()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to build templates: %w", err))
	}

	resp := responder.New(tplProvider)
	metrics := prometheus.New()
	srv := server.New(store, resp, tplProvider, metrics)

	if err != nil {
		log.Fatal(err)
	}
	err = srv.Run(listen)
	if err != nil {
		srv.Close()
		log.Fatal(err)
	}
	srv.Close()
}

func ensureEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s env is empty", name)
	}
	return value
}
