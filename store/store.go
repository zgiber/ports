package store

import (
	"context"
	"encoding/json"
	"log"

	"github.com/zgiber/ports/service"

	"github.com/dgraph-io/badger/v3"
)

const (
	DefaultMaxItems = 1000
)

type Store struct {
	db *badger.DB
}

func New(path string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}

	return &Store{db: db}, nil
}

func (s *Store) SavePort(ctx context.Context, p *service.Port) error {
	return s.db.Update(func(txn *badger.Txn) error {
		val, _ := json.Marshal(p.Details)
		if err := txn.Set([]byte(p.ID), val); err != nil {
			return err
		}

		if ctx.Err() != nil {
			txn.Discard()
		}

		return ctx.Err()
	})
}

func (s *Store) ListPorts(ctx context.Context, afterKey string, maxItems int) ([]*service.Port, error) {
	if maxItems == 0 {
		maxItems = DefaultMaxItems
	}

	ports := make([]*service.Port, 0, maxItems)
	err := s.db.View(func(txn *badger.Txn) error {
		iterator := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()
		iterator.Seek([]byte(afterKey))

		for i := 0; i < maxItems; i++ {
			if !iterator.Valid() {
				break
			}

			if err := iterator.Item().Value(func(val []byte) error {
				portDetails := service.PortDetails{}
				err := json.Unmarshal(val, &portDetails)
				if err != nil {
					return err
				}

				ports = append(ports, &service.Port{
					ID:      string(iterator.Item().Key()),
					Details: portDetails,
				})

				iterator.Next()
				return nil
			}); err != nil {
				return err
			}
		}

		if ctx.Err() != nil {
			txn.Discard()
		}

		return ctx.Err()
	})

	return ports, err
}

func (s *Store) Close() {
	s.db.Close()
}
