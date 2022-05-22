package service

import (
	"context"
	"errors"
	"fmt"
	"io"
)

var (
	ErrStoreOperationFailed = errors.New("the store operation has failed")
)

type Port struct {
	ID      string
	Details PortDetails
}

type PortDetails struct {
	Name        string        `json:"name,omitempty"`
	City        string        `json:"city,omitempty"`
	Country     string        `json:"country,omitempty"`
	Alias       []interface{} `json:"alias,omitempty"`
	Regions     []interface{} `json:"regions,omitempty"`
	Coordinates []float64     `json:"coordinates,omitempty"`
	Province    string        `json:"province,omitempty"`
	Timezone    string        `json:"timezone,omitempty"`
	UNLocs      []string      `json:"unlocs,omitempty"`
	Code        string        `json:"code,omitempty"`
}

type Feed interface {
	Next() (*Port, error)
}

type Store interface {
	SavePort(ctx context.Context, port *Port) error
	ListPorts(ctx context.Context, afterID string, maxItems int) ([]*Port, error)
	Close()
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) UpdatePorts(ctx context.Context, feed Feed) error {
	for {
		port, err := feed.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := s.store.SavePort(ctx, port); err != nil {
			return err
		}
	}
}

type PortsFilter struct {
	AfterID  string
	MaxItems int
}

func (s *Service) ListPorts(ctx context.Context, filter PortsFilter) ([]*Port, error) {
	ports, err := s.store.ListPorts(ctx, filter.AfterID, filter.MaxItems)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStoreOperationFailed, err.Error())
	}

	return ports, nil
}
