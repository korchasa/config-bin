package pkg

import (
    "github.com/google/uuid"
    "time"
)

type Bin struct {
    ID        uuid.UUID
    Data      string
    Version   int
    CreatedAt time.Time
    History   []Bin
}

func (b *Bin) GetVersion() int {
    return b.Version
}
