package http

import (
    "configBin/pkg"
    "github.com/google/uuid"
)

type Storage interface {
    CreateBin(id uuid.UUID, pass string, unencryptedData string) error
    // GetBin returns the bin with the given ID.
    GetBin(id uuid.UUID, pass string) (*pkg.Bin, error)
    UpdateBin(id uuid.UUID, pass string, unencryptedData string) error
    IsReady() bool
    Close()
}
