package store

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zgiber/ports/service"
)

var (
	testPorts = []*service.Port{
		{
			ID: "ABCDE",
			Details: service.PortDetails{
				Name:        "Port A",
				City:        "City A",
				Country:     "Country A",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{1.1, 2.2},
				Province:    "Province A",
				Timezone:    "TZ A",
				UNLocs:      []string{"ABCDE"},
				Code:        "1234",
			},
		},
		{
			ID: "BCDEF",
			Details: service.PortDetails{
				Name:        "Port B",
				City:        "City B",
				Country:     "Country B",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{1.1, 2.2},
				Province:    "Province B",
				Timezone:    "TZ B",
				UNLocs:      []string{"BCDEF"},
				Code:        "2345",
			},
		},
		{
			ID: "CDEFG",
			Details: service.PortDetails{
				Name:        "Port C",
				City:        "City C",
				Country:     "Country C",
				Alias:       []interface{}{},
				Regions:     []interface{}{},
				Coordinates: []float64{1.1, 2.2},
				Province:    "Province C",
				Timezone:    "TZ C",
				UNLocs:      []string{"CDEFG"},
				Code:        "3456",
			},
		},
	}
)

func TestSavePort(t *testing.T) {
	path, err := os.MkdirTemp(os.TempDir(), "badger_test_*")
	if !assert.NoError(t, err) {
		return
	}

	defer os.RemoveAll(path)

	store, err := New(path)
	if !assert.NoError(t, err) {
		return
	}
	defer store.Close()

	port := testPorts[0]
	assert.NoError(t, store.SavePort(context.Background(), port))
}

func TestListPorts(t *testing.T) {
	path, err := os.MkdirTemp(os.TempDir(), "badger_test_*")
	if !assert.NoError(t, err) {
		return
	}

	defer os.RemoveAll(path)
	store, err := New(path)
	if !assert.NoError(t, err) {
		return
	}
	defer store.Close()

	for _, p := range testPorts {
		if !assert.NoError(t, store.SavePort(context.Background(), p)) {
			return
		}
	}

	testCases := []struct {
		description   string
		iterateFrom   string
		maxItems      int
		expectedPorts []*service.Port
	}{
		{
			description:   "iterate from first item",
			iterateFrom:   "",
			expectedPorts: testPorts[0:3],
		},
		{
			description:   "iterate from non existent key before second item",
			iterateFrom:   "BBBBB",
			expectedPorts: testPorts[1:3],
		},
		{
			description:   "iterate from second item",
			iterateFrom:   testPorts[1].ID,
			expectedPorts: testPorts[1:3],
		},
		{
			description:   "return limited items",
			iterateFrom:   "",
			maxItems:      2,
			expectedPorts: testPorts[0:2],
		},
	}

	for _, tc := range testCases {
		ports, err := store.ListPorts(context.Background(), tc.iterateFrom, tc.maxItems)
		assert.NoError(t, err, tc.description)
		if !assert.Len(t, ports, len(tc.expectedPorts), tc.description) {
			return
		}

		for i := 0; i < len(tc.expectedPorts); i++ {
			assert.EqualValues(t, tc.expectedPorts[i], ports[i], tc.description)
		}
	}
}
