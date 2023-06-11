package server_test

import (
    "configBin/pkg/encryptor/aes"
    "configBin/pkg/metrics/fake"
    "configBin/pkg/server"
    "configBin/pkg/server/responder"
    "configBin/pkg/server/templates"
    "configBin/pkg/server/utils"
    "configBin/pkg/storage/sqlite"
    "fmt"
    "github.com/google/uuid"
    "net/http"
    "strings"
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

type formRequestSpec struct {
    method         string
    path           string
    formData       string
    cookieBid      uuid.UUID
    cookiePassword string
}

func formRequest(spec formRequestSpec) *http.Request {
    form := strings.NewReader(spec.formData)
    req, _ := http.NewRequest(spec.method, spec.path, form)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    return req
}

func formRequestWithCookie(spec formRequestSpec) *http.Request {
    req := formRequest(spec)
    req.AddCookie(utils.PassCookie(spec.cookieBid, spec.cookiePassword))
    return req
}
