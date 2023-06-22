package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var ErrNoBinIDInPath = errors.New("no bin id in path")

func ExtractBinIDFromPathVar(r *http.Request) (*uuid.UUID, error) {
	urlVars := mux.Vars(r)
	bid, exists := urlVars["bid"]
	if !exists {
		return nil, ErrNoBinIDInPath
	}
	id, err := uuid.Parse(bid)
	if err != nil {
		return nil, fmt.Errorf("invalid bin id: %w", err)
	}
	return &id, nil
}
