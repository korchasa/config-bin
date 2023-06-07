package pkg

import (
    "github.com/google/uuid"
    "time"
)

type Bin struct {
    ID             uuid.UUID
    CurrentVersion int
    History        []Configuration
}

type Configuration struct {
    EncryptedData   string
    UnencryptedData string
    Version         int
    CreatedAt       time.Time
    Format          string
}
