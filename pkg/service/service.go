package service

import (
    "configBin/pkg"
    "configBin/pkg/http"
    "fmt"
    "github.com/google/uuid"
)

type Service struct {
    store   http.Storage
    metrics http.Metrics
}

func New(store http.Storage, metrics http.Metrics) *Service {
    return &Service{
        store:   store,
        metrics: metrics,
    }
}

func (s *Service) CreateBin(id uuid.UUID, pass string, unencryptedData string) error {
    err := s.store.CreateBin(id, pass, unencryptedData)
    if err != nil {
        return fmt.Errorf("failed to create bin: %v", err)
    }
    return nil
}

func (s *Service) GetBin(id uuid.UUID, pass string) (*pkg.Bin, error) {
    b, err := s.store.GetBin(id, pass)
    if err != nil {
        return nil, fmt.Errorf("failed to get bin: %v", err)
    }
    return b, nil
}

func (s *Service) UpdateBin(id uuid.UUID, pass string, unencryptedData string) error {
    err := s.store.UpdateBin(id, pass, unencryptedData)
    if err != nil {
        return fmt.Errorf("failed to get bin: %v", err)
    }
    return nil
}
