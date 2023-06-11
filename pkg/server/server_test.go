package server_test

import (
    "configBin/pkg/encryptor/aes"
    "configBin/pkg/metrics/fake"
    "configBin/pkg/server"
    "configBin/pkg/server/responder"
    "configBin/pkg/server/templates"
    "configBin/pkg/storage/sqlite"
    "fmt"
)

func NewTestingServer(sqlitePath string) (*server.Server, server.Storage, error) {
    enc := aes.NewAESEncryptor()

    store, err := sqlite.NewSqliteStorage(sqlitePath, enc)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create sqlite storage: %w", err)
    }
    err = store.ApplySchema()
    if err != nil {
        return nil, nil, fmt.Errorf("failed to apply db schema: %w", err)
    }

    tplProvider, err := templates.Build()
    if err != nil {
        return nil, nil, fmt.Errorf("failed to build templates: %w", err)
    }

    resp := responder.New(tplProvider)
    srv := server.New(store, resp, tplProvider, fake.Fake{})

    return srv, store, nil
}
