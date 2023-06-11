package utils

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "net/http"
)

func ExtractBIDFromPathVar(r *http.Request) (*uuid.UUID, error) {
    urlVars := mux.Vars(r)
    bid, exists := urlVars["bid"]
    if !exists {
        return nil, fmt.Errorf("no bin id in path")
    }
    id, err := uuid.Parse(bid)
    if err != nil {
        return nil, fmt.Errorf("invalid bin id: %w", err)
    }
    return &id, nil
}
