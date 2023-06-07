package http

import (
    "configBin/pkg"
    "github.com/google/uuid"
)

type Service interface {
    CreateBin(id uuid.UUID, pass string, unencryptedData string) (string, error)
    GetBin(id uuid.UUID, pass string) (*pkg.Bin, error)
    UpdateBin(id uuid.UUID, pass string, unencryptedData string) error
    RollbackBin(id uuid.UUID, pass string, version string) error
}
