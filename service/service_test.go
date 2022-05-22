package service

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testPorts = []*Port{
		{
			ID: "A",
			Details: PortDetails{
				Name:        "Port A",
				City:        "City A",
				Country:     "Country A",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{1.1234, 2.1234},
				Province:    "Province A",
				Timezone:    "Timezone A",
				UNLocs:      []string{"LOCA"},
				Code:        "12345",
			},
		},
		{
			ID: "B",
			Details: PortDetails{
				Name:        "Port B",
				City:        "City B",
				Country:     "Country B",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{2.1234, 3.1234},
				Province:    "Province B",
				Timezone:    "Timezone B",
				UNLocs:      []string{"LOCA"},
				Code:        "23456",
			},
		},
		{
			ID: "C",
			Details: PortDetails{
				Name:        "Port C",
				City:        "City C",
				Country:     "Country C",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{3.1234, 4.1234},
				Province:    "Province C",
				Timezone:    "Timezone C",
				UNLocs:      []string{"LOCA"},
				Code:        "53456",
			},
		},
	}
)

func TestUpdatePorts(t *testing.T) {
	// given a service that has been instantiated with
	// a test store implementation and a feed with the test data
	store := &testStore{stored: map[string]PortDetails{}}
	feed := newTestFeed(testPorts)
	s := New(store)
	ctx := context.Background()

	// when UpdatePorts is called on the service
	assert.NoError(t, s.UpdatePorts(ctx, feed))

	// then all the ports are stored successfully
	for _, testPort := range testPorts {
		stored, ok := store.stored[testPort.ID]
		assert.True(t, ok)
		assert.EqualValues(t, testPort.Details, stored)
	}

	// and the ports can be retrieved successfully
	ports, err := s.ListPorts(ctx, PortsFilter{})
	assert.NoError(t, err)
	assert.Len(t, ports, len(store.stored))
}

type testStore struct {
	stored map[string]PortDetails
}

func (ts *testStore) SavePort(_ context.Context, port *Port) error {
	ts.stored[port.ID] = port.Details
	return nil
}

func (ts *testStore) ListPorts(_ context.Context, _ string, _ int) ([]*Port, error) {
	ports := make([]*Port, 0, len(ts.stored))
	for id, details := range ts.stored {
		ports = append(ports, &Port{
			ID:      id,
			Details: details,
		})
	}
	return ports, nil
}

func (ts *testStore) Close() {}

func newTestFeed(ports []*Port) *testFeed {
	tf := &testFeed{
		ports: make(chan *Port, len(ports)),
	}

	for _, p := range ports {
		tf.ports <- p
	}

	return tf
}

type testFeed struct {
	ports chan *Port
}

func (tf *testFeed) Next() (*Port, error) {
	select {
	case p := <-tf.ports:
		return p, nil
	default:
		return nil, io.EOF
	}
}
