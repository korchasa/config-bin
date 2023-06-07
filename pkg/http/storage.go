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
    RollbackBin(id uuid.UUID, pass string, version int) error
    IsReady() bool
    Close()
}
