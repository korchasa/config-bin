package main

import (
    "configBin/pkg/encryptor/aes"
    "configBin/pkg/http"
    "configBin/pkg/http/templates"
    "configBin/pkg/prometheus"
    "configBin/pkg/storage/sqlite"
    log "github.com/sirupsen/logrus"
    "os"
)

func init() {
    log.SetOutput(os.Stdout)
    log.SetReportCaller(false)
    log.SetLevel(log.DebugLevel)
    log.SetFormatter(
        &log.TextFormatter{
            ForceColors: true,
        },
    )
}

func main() {
    enc := aes.NewAESEncryptor()

    store, err := sqlite.NewSqliteStorage("var/db.sqlite", enc)
    if err != nil {
        log.Fatal(err)
    }
    defer store.Close()

    tplProvider, err := templates.Build()
    if err != nil {
        log.Fatal(err)
    }

    metrics := prometheus.New()

    httpServer := http.NewServer(store, metrics, tplProvider)
    err = httpServer.Run("localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
}
